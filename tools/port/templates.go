package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"io"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"

	starlarkjson "go.starlark.net/lib/json"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
	"go.yaml.in/yaml/v4"
)

// The pictures of the instances of templates (amendment T52 to R12): a template's picture path
// holds placeholders, and the module ships a file for every instance, which task.known-pictures
// holds (format/README.md, section 23). Each such path has a source in a folder of its own, in the
// schema of standalone's pictures with placeholders among its strings, which the port fills in for
// each instance and draws as it draws a task's picture, into media/pictures/ of the module; or the
// source names a function of the template's solver, which builds the picture that the solver reads
// from the instance's numbers (amendment T53 to R12).
const templatePictureFolder = "media/pictures"

// placeholder is a placeholder of a parameter or a derived number, in a string of a source or in a
// picture path (format/README.md, section 23).
var placeholder = regexp.MustCompile(`\{([a-z][a-z0-9_]*)\}`)

// namedPlaceholder is a placeholder that is a whole part of a name after a -, as a picture path
// holds one (format/README.md, section 23).
var namedPlaceholder = regexp.MustCompile(`-\{[a-z][a-z0-9_]*\}`)

// functionName is a name of a function of Starlark, which a source of solver names.
var functionName = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// maxCombinations is the most combinations of parameters that a template may have
// (task.template-instances).
const maxCombinations = 10000

// templateTask is what the port reads of a task of the module that is a template: its numbers, the
// program of its solver, and the picture parts of its statement and solution. module-kit holds the rest to the format.
type templateTask struct {
	ID       string `yaml:"id"`
	Template *struct {
		Parameters map[string]struct {
			From *int64 `yaml:"from"`
			To   *int64 `yaml:"to"`
		} `yaml:"parameters"`
		Derived    map[string]string `yaml:"derived"`
		Constraint string            `yaml:"constraint"`
		Solver     struct {
			Program string `yaml:"program"`
		} `yaml:"solver"`
	} `yaml:"template"`
	Statement []picturePart `yaml:"statement"`
	Solution  []picturePart `yaml:"solution"`
}

// picturePart is a part of a statement or of a solution, as far as its picture goes.
type picturePart struct {
	Picture string `yaml:"picture"`
	Width   int    `yaml:"width"`
	Height  int    `yaml:"height"`
}

// instancesOf are the numbers of every instance of a template, parameters and derived numbers, in
// the order of the format: the combinations of the parameters as an odometer, the parameters in
// the order of their names, the first changing slowest and each through its domain upwards, and
// those for which the constraint holds (format/README.md, section 23). The expressions are
// evaluated by go.starlark.net at the version of the format's dialect, as module-kit does. The
// results of the solver are not among the numbers: the port does not call solve.
func instancesOf(task templateTask) ([]map[string]int64, error) {
	t := task.Template
	names := slices.Sorted(maps.Keys(t.Parameters))
	if len(names) == 0 {
		return nil, fmt.Errorf("%s: its template has no parameter", task.ID)
	}
	count := int64(1)
	for _, name := range names {
		domain := t.Parameters[name]
		if domain.From == nil || domain.To == nil || *domain.To < *domain.From {
			return nil, fmt.Errorf("%s: its parameter %s has no domain from and to", task.ID, name)
		}
		// The size is held to the limit before it is counted, so that it cannot overflow.
		if uint64(*domain.To-*domain.From) >= maxCombinations {
			return nil, fmt.Errorf("%s: its domains multiply past %d combinations", task.ID, maxCombinations)
		}
		count *= *domain.To - *domain.From + 1
		if count > maxCombinations {
			return nil, fmt.Errorf("%s: its domains multiply past %d combinations", task.ID, maxCombinations)
		}
	}
	thread := &starlark.Thread{Name: task.ID}
	var instances []map[string]int64
	values := make([]int64, len(names))
	for i, name := range names {
		values[i] = *t.Parameters[name].From
	}
	for {
		numbers := map[string]int64{}
		env := starlark.StringDict{}
		for i, name := range names {
			numbers[name] = values[i]
			env[name] = starlark.MakeInt64(values[i])
		}
		derivedEnv := maps.Clone(env)
		for _, name := range slices.Sorted(maps.Keys(t.Derived)) {
			value, err := starlark.EvalOptions(dialect, thread, task.ID, t.Derived[name], env)
			if err != nil {
				return nil, fmt.Errorf("%s: its derived number %s: %w", task.ID, name, err)
			}
			number, isInt := value.(starlark.Int)
			var n int64
			fits := false
			if isInt {
				n, fits = number.Int64()
			}
			if !fits {
				return nil, fmt.Errorf("%s: its derived number %s is no integer the port can hold", task.ID, name)
			}
			numbers[name] = n
			derivedEnv[name] = number
		}
		holds := true
		if t.Constraint != "" {
			value, err := starlark.EvalOptions(dialect, thread, task.ID, t.Constraint, derivedEnv)
			if err != nil {
				return nil, fmt.Errorf("%s: its constraint: %w", task.ID, err)
			}
			truth, isBool := value.(starlark.Bool)
			if !isBool {
				return nil, fmt.Errorf("%s: its constraint is no bool", task.ID)
			}
			holds = bool(truth)
		}
		if holds {
			instances = append(instances, numbers)
		}
		// The odometer: the last parameter turns fastest.
		i := len(names) - 1
		for ; i >= 0; i-- {
			if values[i] < *t.Parameters[names[i]].To {
				values[i]++
				break
			}
			values[i] = *t.Parameters[names[i]].From
		}
		if i < 0 {
			break
		}
	}
	if len(instances) == 0 {
		return nil, fmt.Errorf("%s: no combination of its template is an instance", task.ID)
	}
	return instances, nil
}

// substitution is what the format's substitution reads from left to right in a single pass: a
// doubled brace, which is one brace, or a placeholder (format/README.md, section 23).
var substitution = regexp.MustCompile(`\{\{|\}\}|\{[a-z][a-z0-9_]*\}`)

// fill puts an instance's numbers in place of the placeholders of a text, and one brace in place of
// each doubled one, as the format does; a placeholder that names no parameter and no derived number
// is an error, a result of the solver among them.
func fill(text string, numbers map[string]int64) (string, error) {
	var missing string
	filled := substitution.ReplaceAllStringFunc(text, func(found string) string {
		if found == "{{" || found == "}}" {
			return found[:1]
		}
		name := found[1 : len(found)-1]
		number, known := numbers[name]
		if !known {
			if missing == "" {
				missing = name
			}
			return found
		}
		return strconv.FormatInt(number, 10)
	})
	if missing != "" {
		return "", fmt.Errorf("{%s} names no parameter and no derived number of the template", missing)
	}
	return filled, nil
}

// sourceName is the name of the source of a template's picture: its path without pictures/, its
// placeholders and .png, so that pictures/groups-solution-{seat}-{size}.png is groups-solution.
func sourceName(path string) (string, bool) {
	name, inPictures := strings.CutPrefix(path, "pictures/")
	name, isPNG := strings.CutSuffix(name, ".png")
	if !inPictures || !isPNG || !placeholder.MatchString(name) {
		return "", false
	}
	return namedPlaceholder.ReplaceAllString(name, ""), true
}

// pathPattern is the pattern of the names of a template's pictures in media/pictures/: the path
// with a number in place of each placeholder.
func pathPattern(path string) *regexp.Regexp {
	name := strings.TrimPrefix(path, "pictures/")
	parts := placeholder.Split(name, -1)
	for i, part := range parts {
		parts[i] = regexp.QuoteMeta(part)
	}
	return regexp.MustCompile("^" + strings.Join(parts, "[0-9]+") + "$")
}

// variant is a source with a picture for some instances: when names the numbers they have, and a
// variant with no when is for every instance. A variant of solver has no picture of its own: it is
// what that function of the template's solver gives on the instance's numbers (amendment T53 to
// R12).
type variant struct {
	When    map[string]int64 `json:"when"`
	Picture json.RawMessage  `json:"picture"`
	Solver  string           `json:"-"`
}

// readSource reads a source of a template's picture: one picture, a list of variants, or
// {"solver": "<function>"}, the function of the template's solver that builds the picture.
func readSource(file string) ([]variant, error) {
	raw, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	// A source of solver is told by its member, whatever its value, so that a value of the wrong
	// type is refused as such and not read as a picture.
	var members map[string]json.RawMessage
	if json.Unmarshal(raw, &members) == nil && members["solver"] != nil {
		if err := exactKeys(raw, ""); err != nil {
			return nil, fmt.Errorf("%s: %w", file, err)
		}
		var function string
		if len(members) != 1 {
			return nil, fmt.Errorf("%s: a source of solver has no member but solver", file)
		}
		if err := json.Unmarshal(members["solver"], &function); err != nil || !functionName.MatchString(function) {
			return nil, fmt.Errorf("%s: its solver, %s, is no name of a function", file, members["solver"])
		}
		return []variant{{Solver: function}}, nil
	}
	if trimmed := strings.TrimSpace(string(raw)); strings.HasPrefix(trimmed, "[") {
		var variants []variant
		decoder := json.NewDecoder(strings.NewReader(trimmed))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&variants); err != nil {
			return nil, fmt.Errorf("%s: %w", file, err)
		}
		if len(variants) == 0 {
			return nil, fmt.Errorf("%s holds no variant", file)
		}
		for _, v := range variants {
			if err := sourceKeys(v.Picture); err != nil {
				return nil, fmt.Errorf("%s: %w", file, err)
			}
		}
		return variants, nil
	}
	if err := sourceKeys(raw); err != nil {
		return nil, fmt.Errorf("%s: %w", file, err)
	}
	return []variant{{Picture: raw}}, nil
}

// sourceKeys holds a picture of a source to the names of its members as readPicture holds a picture:
// the source is decoded and written again before it is drawn, which would keep only the last of a
// member named twice.
func sourceKeys(raw json.RawMessage) error {
	var kind struct {
		Kind string `json:"kind"`
	}
	if err := json.Unmarshal(raw, &kind); err != nil {
		return err
	}
	return exactKeys(raw, keyedByData[kind.Kind])
}

// solverSteps is the budget of steps of a template's solver, the top level and every call together
// (format/README.md, section 24).
var solverSteps uint64 = 100000000

// solverProgram is the solver of a template, as the port runs it to build its pictures: once, on one
// thread under the dialect of the format, and then a function of it on each instance, as module-kit
// runs solve (format/README.md, section 24).
type solverProgram struct {
	id      string
	source  string
	thread  *starlark.Thread
	globals starlark.StringDict
	err     error
}

// dialect are the options of the format's dialect: those that standalone's solvers run with
// (format/README.md, section 24).
var dialect = &syntax.FileOptions{Set: true, While: true, TopLevelControl: true, GlobalReassign: true}

// solverCall is a function of a template's solver to call for a picture: on the numbers named shown,
// or on all the instance's numbers when shown is nil.
type solverCall struct {
	program *solverProgram
	shown   map[string]bool
}

// figure is what the function of the solver gives on an instance's numbers, those named shown or
// all its parameters and derived numbers when shown is nil, in a frozen dict by their names in
// order, as JSON: a picture in the schema of standalone's pictures, as the solver reads it.
func (p *solverProgram) figure(function string, numbers map[string]int64, shown map[string]bool) (json.RawMessage, error) {
	if p.thread == nil {
		p.thread = &starlark.Thread{Name: p.id}
		// print writes nothing in the dialect of the format (format/README.md, section 24).
		p.thread.Print = func(*starlark.Thread, string) {}
		p.thread.SetMaxExecutionSteps(solverSteps)
		if p.source == "" {
			p.err = fmt.Errorf("%s: a source of its picture is the solver's %s, and its template has no solver", p.id, function)
		} else if p.globals, p.err = starlark.ExecFileOptions(dialect, p.thread, p.id, p.source, nil); p.err != nil {
			p.err = fmt.Errorf("%s: its solver: %w", p.id, p.err)
		}
	}
	if p.err != nil {
		return nil, p.err
	}
	fn, isFunction := p.globals[function].(*starlark.Function)
	if !isFunction {
		return nil, fmt.Errorf("%s: its solver defines no function %s", p.id, function)
	}
	dict := starlark.NewDict(len(numbers))
	for _, name := range slices.Sorted(maps.Keys(numbers)) {
		if shown != nil && !shown[name] {
			continue
		}
		if err := dict.SetKey(starlark.String(name), starlark.MakeInt64(numbers[name])); err != nil {
			return nil, err
		}
	}
	dict.Freeze()
	value, err := starlark.Call(p.thread, fn, starlark.Tuple{dict}, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: its solver's %s on %v: %w", p.id, function, numbers, err)
	}
	if _, isDict := value.(*starlark.Dict); !isDict {
		return nil, fmt.Errorf("%s: its solver's %s gives a %s, not a picture", p.id, function, value.Type())
	}
	encoded, err := starlark.Call(p.thread, starlarkjson.Module.Members["encode"], starlark.Tuple{value}, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: its solver's %s: %w", p.id, function, err)
	}
	return json.RawMessage(encoded.(starlark.String)), nil
}

// pictureFor is the picture of a source for an instance: the first variant whose when the instance's
// numbers have, with the numbers in place of the placeholders of its strings; or what the function
// of the solver that the source names gives on them.
func pictureFor(id string, variants []variant, numbers map[string]int64, call *solverCall) (drawing, error) {
	for _, v := range variants {
		if v.Solver != "" {
			if call == nil {
				return nil, fmt.Errorf("%s: a source of its picture is the solver's %s, and no solver is given", id, v.Solver)
			}
			raw, err := call.program.figure(v.Solver, numbers, call.shown)
			if err != nil {
				return nil, err
			}
			picture, err := readPicture(id, raw)
			if err != nil {
				return nil, fmt.Errorf("%w, in the picture that its solver's %s builds", err, v.Solver)
			}
			return picture, nil
		}
		matches := true
		for name, want := range v.When {
			have, known := numbers[name]
			if !known {
				return nil, fmt.Errorf("%s: a variant of its source is when %s, which is no parameter and no derived number of the template", id, name)
			}
			if have != want {
				matches = false
			}
		}
		if !matches {
			continue
		}
		var tree any
		decoder := json.NewDecoder(strings.NewReader(string(v.Picture)))
		decoder.UseNumber()
		if err := decoder.Decode(&tree); err != nil {
			return nil, fmt.Errorf("%s: %w", id, err)
		}
		filled, err := fillTree(tree, numbers)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", id, err)
		}
		raw, err := json.Marshal(filled)
		if err != nil {
			return nil, err
		}
		return readPicture(id, raw)
	}
	return nil, fmt.Errorf("%s: no variant of its source is for the instance %v", id, numbers)
}

// fillTree puts the numbers in place of the placeholders of every string of a decoded source.
func fillTree(tree any, numbers map[string]int64) (any, error) {
	switch node := tree.(type) {
	case string:
		return fill(node, numbers)
	case []any:
		for i, each := range node {
			filled, err := fillTree(each, numbers)
			if err != nil {
				return nil, err
			}
			node[i] = filled
		}
	case map[string]any:
		// The keys too: the marks of a grid and of a calendar are keyed by data.
		filledNode := map[string]any{}
		for key, each := range node {
			filledKey, err := fill(key, numbers)
			if err != nil {
				return nil, err
			}
			filled, err := fillTree(each, numbers)
			if err != nil {
				return nil, err
			}
			if _, twice := filledNode[filledKey]; twice {
				return nil, fmt.Errorf("two members are %q once filled", filledKey)
			}
			filledNode[filledKey] = filled
		}
		return filledNode, nil
	}
	return tree, nil
}

// drawTemplatePictures draws the picture of every instance of every template of the module whose
// picture path holds placeholders, from the sources of the folder, by the path of each file in the
// module; and gives the pattern of the names of each such picture, to find the stale ones. Every
// instance of a picture is drawn at the size the template gives it, the largest drawing of them
// all, each on white in the middle; a source no template names is an error.
func drawTemplatePictures(folder, module string) (map[string][]byte, []*regexp.Regexp, int, error) {
	tasks, err := filepath.Glob(filepath.Join(module, "tasks", "*.yaml"))
	if err != nil {
		return nil, nil, 0, err
	}
	files := map[string][]byte{}
	var patterns []*regexp.Regexp
	used := map[string]bool{}
	templates := 0
	for _, file := range tasks {
		data, err := os.ReadFile(file)
		if err != nil {
			return nil, nil, 0, err
		}
		var task templateTask
		if err := yaml.Unmarshal(data, &task); err != nil {
			return nil, nil, 0, fmt.Errorf("%s: %w", file, err)
		}
		if task.Template == nil {
			continue
		}
		var parts []picturePart
		inStatement := map[string]bool{}
		for i, part := range slices.Concat(task.Statement, task.Solution) {
			if _, has := sourceName(part.Picture); has {
				parts = append(parts, part)
				inStatement[part.Picture] = i < len(task.Statement)
			}
		}
		if len(parts) == 0 {
			continue
		}
		templates++
		instances, err := instancesOf(task)
		if err != nil {
			return nil, nil, 0, err
		}
		program := &solverProgram{id: task.ID, source: task.Template.Solver.Program}
		shown, err := shownNames(data)
		if err != nil {
			return nil, nil, 0, fmt.Errorf("%s: %w", file, err)
		}
		for _, part := range parts {
			name, _ := sourceName(part.Picture)
			variants, err := readSource(filepath.Join(folder, name+".json"))
			if err != nil {
				return nil, nil, 0, fmt.Errorf("%s: its picture %s: %w", task.ID, part.Picture, err)
			}
			used[name+".json"] = true
			// The picture of the statement is built from the numbers that solve gets, so that it is
			// the picture that solve reads; that of the solution, which solve does not read, from all.
			call := &solverCall{program: program}
			if inStatement[part.Picture] {
				call.shown = shown
			}
			drawn, err := drawInstances(task.ID, part, variants, instances, call)
			if err != nil {
				return nil, nil, 0, err
			}
			maps.Copy(files, drawn)
			patterns = append(patterns, pathPattern(part.Picture))
		}
	}
	if templates == 0 {
		return nil, nil, 0, fmt.Errorf("%s has no template whose picture path holds placeholders", module)
	}
	entries, err := os.ReadDir(folder)
	if err != nil {
		return nil, nil, 0, err
	}
	for _, entry := range entries {
		if !used[entry.Name()] {
			return nil, nil, 0, fmt.Errorf("%s is the source of no picture of a template of %s", filepath.Join(folder, entry.Name()), module)
		}
	}
	return files, patterns, templates, nil
}

// drawInstances draws a picture part of a template for every instance, at one size: the largest
// drawing, which the part's width and height must give, each drawing in the middle of a white
// image of that size. Each instance names a path of its own: a path that leaves out a number the
// picture shows would put two instances in one file.
func drawInstances(id string, part picturePart, variants []variant, instances []map[string]int64, call *solverCall) (map[string][]byte, error) {
	drawings := map[string]drawing{}
	var paths []string
	width, height := 0, 0
	for _, numbers := range instances {
		path, err := fill(part.Picture, numbers)
		if err != nil {
			return nil, fmt.Errorf("%s: its picture path %s: %w", id, part.Picture, err)
		}
		if _, twice := drawings[path]; twice {
			return nil, fmt.Errorf("%s: two of its instances name the picture %s: its path names too few of their numbers", id, path)
		}
		picture, err := pictureFor(id, variants, numbers, call)
		if err != nil {
			return nil, err
		}
		w, h := picture.size()
		width, height = max(width, w), max(height, h)
		drawings[path] = picture
		paths = append(paths, path)
	}
	if part.Width != width || part.Height != height {
		return nil, fmt.Errorf("%s: its picture %s is %d by %d, the largest of its instances, and the template gives %d by %d", id, part.Picture, width, height, part.Width, part.Height)
	}
	files := map[string][]byte{}
	for _, path := range paths {
		img := paintImage(drawings[path])
		padded := whiteImage(width, height)
		offset := image.Pt((padded.Bounds().Dx()-img.Bounds().Dx())/2, (padded.Bounds().Dy()-img.Bounds().Dy())/2)
		draw.Draw(padded, img.Bounds().Add(offset), img, image.Point{}, draw.Src)
		data, err := encodePNG(padded)
		if err != nil {
			return nil, err
		}
		files["media/"+path] = data
	}
	return files, nil
}

// staleTemplatePictures names the files of the module's media/pictures/ that are named as a
// picture of a template's instance and that the port does not draw now. A file of that name that
// is the picture of a task, tasks/<its name>.yaml with no template, is an error, not a stale file:
// the pictures of the tasks are the port's of the batches (R12).
func staleTemplatePictures(module string, files map[string][]byte, patterns []*regexp.Regexp) ([]string, error) {
	entries, err := os.ReadDir(filepath.Join(module, templatePictureFolder))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var old []string
	for _, entry := range entries {
		path := templatePictureFolder + "/" + entry.Name()
		if _, draws := files[path]; draws {
			continue
		}
		for _, pattern := range patterns {
			if !pattern.MatchString(entry.Name()) {
				continue
			}
			task := strings.TrimSuffix(strings.TrimSuffix(entry.Name(), ".png"), "-solution")
			if _, err := os.Stat(filepath.Join(module, "tasks", task+".yaml")); err == nil {
				return nil, fmt.Errorf("%s is named as a picture of a template's instance, and it is the picture of task %s", filepath.Join(module, path), task)
			}
			old = append(old, path)
			break
		}
	}
	return old, nil
}

// runTemplatePictures draws the pictures of the instances of the module's templates, or with check
// compares them.
func runTemplatePictures(folder, module string, check bool, stdout, stderr io.Writer) int {
	err := isFolder(module)
	var files map[string][]byte
	var patterns []*regexp.Regexp
	var templates int
	if err == nil {
		files, patterns, templates, err = drawTemplatePictures(folder, module)
	}
	var old, differ []string
	if err == nil {
		old, err = staleTemplatePictures(module, files, patterns)
	}
	if err == nil && check {
		differ, err = filesDiffer(module, files, old)
	} else if err == nil {
		err = placeFiles(module, files, old)
	}
	if err != nil {
		fmt.Fprintf(stderr, "port: %v\n", err)
		return statusCannotRun
	}
	if !check {
		fmt.Fprintf(stdout, "port: %d pictures of the instances of %d templates written into %s\n", len(files), templates, module)
		return statusDone
	}
	for _, path := range differ {
		fmt.Fprintf(stdout, "port: %s is not what the port draws\n", filepath.Join(module, path))
	}
	if len(differ) > 0 {
		return statusDiffers
	}
	fmt.Fprintf(stdout, "port: the %d pictures of the instances of %d templates of %s are what the port draws\n", len(files), templates, module)
	return statusDone
}

// shownNames are the names of the numbers that solve of a template gets: those whose placeholder
// stands in the open fields that the format names — the statement, its paragraphs and pictures,
// the givens, the unknown, the text of an item, the text pieces of its line and the options of its
// answer slots bound choice (format/README.md, section 24).
func shownNames(data []byte) (map[string]bool, error) {
	var task struct {
		Statement any `yaml:"statement"`
		Givens    any `yaml:"givens"`
		Unknown   any `yaml:"unknown"`
		Items     []struct {
			Text any `yaml:"text"`
			Line []struct {
				Text any `yaml:"text"`
			} `yaml:"line"`
			AnswerSlots []struct {
				Binding string `yaml:"binding"`
				Options any    `yaml:"options"`
			} `yaml:"answer_slots"`
		} `yaml:"items"`
	}
	if err := yaml.Unmarshal(data, &task); err != nil {
		return nil, err
	}
	open := []any{task.Statement, task.Givens, task.Unknown}
	for _, item := range task.Items {
		open = append(open, item.Text)
		for _, piece := range item.Line {
			open = append(open, piece.Text)
		}
		for _, slot := range item.AnswerSlots {
			if slot.Binding == "choice" {
				open = append(open, slot.Options)
			}
		}
	}
	shown := map[string]bool{}
	var walk func(node any)
	walk = func(node any) {
		switch node := node.(type) {
		case string:
			for _, found := range substitution.FindAllString(node, -1) {
				if found != "{{" && found != "}}" {
					shown[found[1:len(found)-1]] = true
				}
			}
		case []any:
			for _, each := range node {
				walk(each)
			}
		case map[string]any:
			for _, each := range node {
				walk(each)
			}
		}
	}
	walk(open)
	return shown, nil
}
