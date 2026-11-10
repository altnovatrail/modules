package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"go.yaml.in/yaml/v4"
)

// pairsTemplate is a template of two parameters, a from 0 to 2 and b from 1 to 3, with a != b,
// whose statement shows the picture of its instance at the size given.
func pairsTemplate(width, height int) string {
	return fmt.Sprintf(`id: pairs
template:
  parameters:
    b: {from: 1, to: 3}
    a: {from: 0, to: 2}
  derived:
    sum: a + b
    half: -a // 2
    rest: -7 %% (b + 1)
  constraint: a != b
statement:
  - text: "{a} and {b}"
  - picture: pictures/pairs-{a}-{b}.png
    alt: "{a} and {b}"
    width: %d
    height: %d
`, width, height)
}

const pairsSource = `{"kind": "venn", "total": "{sum}", "sets": [{"label": "A", "count": "{a}"}, {"label": "B", "count": "{b}"}], "both": "?"}`

// templateModule is a module with the tasks given, and a folder of sources of template pictures
// with the files given, beside it.
func templateModule(t *testing.T, tasks, sources map[string]string) (folder, module string) {
	t.Helper()
	root := t.TempDir()
	folder, module = filepath.Join(root, "sources"), filepath.Join(root, "module")
	for _, dir := range []string{folder, filepath.Join(module, "tasks")} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	for name, content := range tasks {
		writeFile(t, filepath.Join(module, "tasks", name), content)
	}
	for name, content := range sources {
		writeFile(t, filepath.Join(folder, name), content)
	}
	return folder, module
}

// pairsSize is the size the port draws the pictures of pairsTemplate at: the largest of them.
func pairsSize(t *testing.T) (int, int) {
	t.Helper()
	width, height := 0, 0
	for _, numbers := range pairsInstances(t) {
		picture, err := pictureFor("pairs", []variant{{Picture: []byte(pairsSource)}}, numbers, nil)
		if err != nil {
			t.Fatal(err)
		}
		w, h := picture.size()
		width, height = max(width, w), max(height, h)
	}
	return width, height
}

func pairsInstances(t *testing.T) []map[string]int64 {
	t.Helper()
	instances, err := instancesOf(readTemplate(t, pairsTemplate(1, 1)))
	if err != nil {
		t.Fatal(err)
	}
	return instances
}

func readTemplate(t *testing.T, text string) templateTask {
	t.Helper()
	_, module := templateModule(t, map[string]string{"pairs.yaml": text}, nil)
	var task templateTask
	data, err := os.ReadFile(filepath.Join(module, "tasks", "pairs.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if err := yaml.Unmarshal(data, &task); err != nil {
		t.Fatal(err)
	}
	return task
}

// The instances come as an odometer, the parameters in the order of their names, the first
// slowest; the derived numbers are Starlark's, with floored division and remainder; and the
// constraint leaves out the combinations it does not hold for (format/README.md, section 23).
func TestInstancesOf(t *testing.T) {
	instances := pairsInstances(t)
	var got []string
	for _, numbers := range instances {
		got = append(got, fmt.Sprintf("%d%d", numbers["a"], numbers["b"]))
	}
	want := []string{"01", "02", "03", "12", "13", "21", "23"}
	if !slices.Equal(got, want) {
		t.Errorf("instances = %v, want %v", got, want)
	}
	last := instances[len(instances)-1]
	if last["sum"] != 5 || last["half"] != -1 || last["rest"] != 1 {
		t.Errorf("the derived numbers of a 2, b 3 = %v: want sum 5, half -1, rest 1", last)
	}
}

func TestInstancesOfRefuses(t *testing.T) {
	for name, change := range map[string][2]string{
		"no domain":        {"a: {from: 0, to: 2}", "a: {from: 0}"},
		"too many":         {"a: {from: 0, to: 2}", "a: {from: 0, to: 9999}"},
		"past int64":       {"a: {from: 0, to: 2}", "a: {from: -9223372036854775808, to: 9223372036854775807}"},
		"no instance":      {"constraint: a != b", "constraint: a > 9"},
		"a constraint int": {"constraint: a != b", "constraint: a + b"},
		"no derived":       {"sum: a + b", "sum: a // 0"},
		"no parameter":     {"    b: {from: 1, to: 3}\n    a: {from: 0, to: 2}\n", ""},
	} {
		task := readTemplate(t, strings.Replace(pairsTemplate(1, 1), change[0], change[1], 1))
		if _, err := instancesOf(task); err == nil {
			t.Errorf("%s: instancesOf gives no error", name)
		}
	}
}

func TestFillAndNames(t *testing.T) {
	if got, err := fill("{a} of {b}", map[string]int64{"a": -2, "b": 10}); err != nil || got != "-2 of 10" {
		t.Errorf("fill = %q, %v", got, err)
	}
	if got, err := fill("{{a}} is {a}}}", map[string]int64{"a": 4}); err != nil || got != "{a} is 4}" {
		t.Errorf("doubled braces: fill = %q, %v", got, err)
	}
	if _, err := fill("{row}", map[string]int64{"a": 1}); err == nil || !strings.Contains(err.Error(), "{row} names no parameter and no derived number") {
		t.Errorf("a result of the solver: error = %v", err)
	}
	if name, has := sourceName("pictures/groups-solution-{seat}-{size}.png"); !has || name != "groups-solution" {
		t.Errorf("sourceName = %q, %v", name, has)
	}
	if _, has := sourceName("pictures/div-56-d3-1.png"); has {
		t.Error("a path with no placeholder has a source")
	}
	pattern := pathPattern("pictures/groups-{seat}-{size}.png")
	for name, matches := range map[string]bool{"groups-97-18.png": true, "groups-solution-97-18.png": false, "groups-97.png": false} {
		if pattern.MatchString(name) != matches {
			t.Errorf("the pattern on %s = %v", name, !matches)
		}
	}
}

func TestTemplatePicturesWriteAndCheck(t *testing.T) {
	width, height := pairsSize(t)
	folder, module := templateModule(t, map[string]string{
		"pairs.yaml": pairsTemplate(width, height),
		"plain.yaml": "id: plain\nstatement:\n  - picture: pictures/plain.png\n",
	}, map[string]string{"pairs.json": pairsSource})
	writeFile(t, filepath.Join(module, templatePictureFolder, "pairs-9-9.png"), "stale")
	writeFile(t, filepath.Join(module, templatePictureFolder, "plain.png"), "a task's own")
	status, out, errs := runPictures(t, "--template-pictures", folder, "--module", module)
	if status != statusDone || out != "port: 7 pictures of the instances of 1 templates written into "+module+"\n" || errs != "" {
		t.Fatalf("run() = %d, %q, %q", status, out, errs)
	}
	for _, name := range []string{"pairs-0-1.png", "pairs-2-3.png", "plain.png"} {
		if _, err := os.Stat(filepath.Join(module, templatePictureFolder, name)); err != nil {
			t.Errorf("%s: %v", name, err)
		}
	}
	if _, err := os.Stat(filepath.Join(module, templatePictureFolder, "pairs-9-9.png")); !os.IsNotExist(err) {
		t.Errorf("pairs-9-9.png, of no instance, is still there: %v", err)
	}
	status, out, _ = runPictures(t, "--template-pictures", folder, "--module", module, "--check")
	if status != statusDone || out != "port: the 7 pictures of the instances of 1 templates of "+module+" are what the port draws\n" {
		t.Errorf("--check on the pictures just drawn = %d, %q", status, out)
	}
	writeFile(t, filepath.Join(module, templatePictureFolder, "pairs-0-1.png"), "another picture")
	writeFile(t, filepath.Join(module, templatePictureFolder, "pairs-1-1.png"), "of no instance")
	status, out, _ = runPictures(t, "--template-pictures", folder, "--module", module, "--check")
	if status != statusDiffers || !strings.Contains(out, "pairs-0-1.png is not what the port draws") || !strings.Contains(out, "pairs-1-1.png is not what the port draws") {
		t.Errorf("--check on a changed and a stale picture = %d, %q", status, out)
	}
}

// A source may give pictures by the numbers of an instance: the first variant whose when the
// instance has, and a variant with no when for the rest.
func TestTemplatePicturesVariants(t *testing.T) {
	other := `{"kind": "venn", "total": "{sum}", "sets": [{"label": "X", "count": "{a}"}, {"label": "Y", "count": "{b}"}], "both": "{a}"}`
	instances := pairsInstances(t)
	variants := []variant{{When: map[string]int64{"a": 0}, Picture: []byte(other)}, {Picture: []byte(pairsSource)}}
	for i, want := range map[int]string{
		0: strings.NewReplacer("{sum}", "1", "{a}", "0", "{b}", "1").Replace(other),
		3: strings.NewReplacer("{sum}", "3", "{a}", "1", "{b}", "2").Replace(pairsSource),
	} {
		drawn, err := pictureFor("pairs", variants, instances[i], nil)
		if err != nil {
			t.Fatal(err)
		}
		got, _ := encodePNG(paintImage(drawn))
		if !bytes.Equal(got, drawnAs(t, want)) {
			t.Errorf("the instance %v is not drawn from its variant", instances[i])
		}
	}
	if _, err := pictureFor("pairs", []variant{{When: map[string]int64{"a": 9}, Picture: []byte(other)}}, instances[0], nil); err == nil ||
		!strings.Contains(err.Error(), "no variant of its source is for the instance") {
		t.Errorf("no variant: error = %v", err)
	}
}

func TestTemplatePicturesRefuses(t *testing.T) {
	width, height := pairsSize(t)
	for name, c := range map[string]struct {
		tasks   map[string]string
		sources map[string]string
		want    string
	}{
		"a size not the template's": {map[string]string{"pairs.yaml": pairsTemplate(width+1, height)}, map[string]string{"pairs.json": pairsSource},
			fmt.Sprintf("is %d by %d, the largest of its instances, and the template gives %d by %d", width, height, width+1, height)},
		"no source": {map[string]string{"pairs.yaml": pairsTemplate(width, height)}, map[string]string{"other.json": pairsSource}, "pairs.json"},
		"a source of no template": {map[string]string{"pairs.yaml": pairsTemplate(width, height)},
			map[string]string{"pairs.json": pairsSource, "other.json": pairsSource}, "other.json is the source of no picture of a template"},
		"a result in a source": {map[string]string{"pairs.yaml": pairsTemplate(width, height)},
			map[string]string{"pairs.json": strings.Replace(pairsSource, "{sum}", "{row}", 1)}, "{row} names no parameter and no derived number"},
		"no variant": {map[string]string{"pairs.yaml": pairsTemplate(width, height)}, map[string]string{"pairs.json": "[]"}, "holds no variant"},
		"no template": {map[string]string{"plain.yaml": "id: plain\n"}, map[string]string{"pairs.json": pairsSource},
			"has no template whose picture path holds placeholders"},
		"a member twice": {map[string]string{"pairs.yaml": pairsTemplate(width, height)},
			map[string]string{"pairs.json": strings.Replace(pairsSource, `"both": "?"`, `"both": "?", "both": "!"`, 1)}, `the member "both" is named in other letters than the schema's, or twice`},
		"one path for two instances": {map[string]string{"pairs.yaml": strings.Replace(pairsTemplate(width, height), "pairs-{a}-{b}.png", "pairs-{a}.png", 1)},
			map[string]string{"pairs.json": pairsSource}, "two of its instances name the picture"},
		"a when of no number": {map[string]string{"pairs.yaml": pairsTemplate(width, height)},
			map[string]string{"pairs.json": `[{"when": {"c": 0}, "picture": ` + pairsSource + `}]`}, "is when c, which is no parameter and no derived number"},
	} {
		folder, module := templateModule(t, c.tasks, c.sources)
		status, _, errs := runPictures(t, "--template-pictures", folder, "--module", module)
		if status != statusCannotRun || !strings.Contains(errs, c.want) {
			t.Errorf("%s: run() = %d, %q", name, status, errs)
		}
	}
	if status, _, errs := runPictures(t, "--template-pictures", "x", "--pictures", "y", "--module", "z"); status != statusCannotRun || !strings.Contains(errs, "usage:") {
		t.Errorf("both kinds of pictures at once = %d, %q", status, errs)
	}
}

// The pictures of the tasks are the port's of the batches: a file named as a picture of an instance
// that is the picture of a task is an error, not a stale file. A source's keys are filled too, as
// the marks of a grid are keyed by a cell.
func TestTemplatePicturesTaskAndKeys(t *testing.T) {
	width, height := pairsSize(t)
	folder, module := templateModule(t, map[string]string{
		"pairs.yaml":     pairsTemplate(width, height),
		"pairs-9-9.yaml": "id: pairs-9-9\n",
	}, map[string]string{"pairs.json": pairsSource})
	writeFile(t, filepath.Join(module, templatePictureFolder, "pairs-9-9.png"), "a task's own")
	if status, _, errs := runPictures(t, "--template-pictures", folder, "--module", module); status != statusCannotRun ||
		!strings.Contains(errs, "pairs-9-9.png is named as a picture of a template's instance, and it is the picture of task pairs-9-9") {
		t.Errorf("run() = %d, %q", status, errs)
	}
	tree, err := fillTree(map[string]any{"marks": map[string]any{"B{a}": "x", "C1": "{b}"}}, map[string]int64{"a": 2, "b": 3})
	if err != nil || fmt.Sprint(tree) != "map[marks:map[B2:x C1:3]]" {
		t.Errorf("fillTree = %v, %v", tree, err)
	}
	if _, err := fillTree(map[string]any{"B{a}": 1, "B2": 2}, map[string]int64{"a": 2}); err == nil || !strings.Contains(err.Error(), `two members are "B2" once filled`) {
		t.Errorf("two keys alike once filled: error = %v", err)
	}
}

// pairsFigure is the solver of pairsTemplate whose figure builds the picture of pairsSource.
const pairsFigure = `
  solver:
    program: |
      print("print writes nothing")

      def figure(numbers):
          return {"kind": "venn", "total": str(numbers["a"] + numbers["b"]), "both": "?",
                  "sets": [{"label": "A", "count": str(numbers["a"])}, {"label": "B", "count": str(numbers["b"])}]}

      def solution(numbers):
          return figure({"a": numbers["a"], "b": numbers["b"] + numbers["sum"] - numbers["sum"]})

      def solve(numbers):
          return {"row": numbers["a"]}
    results: [row]
`

// solverTemplate is pairsTemplate with the solver given.
func solverTemplate(width, height int, solver string) string {
	return strings.Replace(pairsTemplate(width, height), "  constraint: a != b\n", "  constraint: a != b\n"+solver, 1)
}

// A source of solver is the picture that a function of the template's solver builds from the
// instance's numbers, run under the dialect of the format (amendment T53 to R12).
func TestTemplatePicturesSolver(t *testing.T) {
	width, height := pairsSize(t)
	task := readTemplate(t, solverTemplate(width, height, pairsFigure))
	call := &solverCall{program: &solverProgram{id: "pairs", source: task.Template.Solver.Program}}
	for _, numbers := range pairsInstances(t) {
		built, err := pictureFor("pairs", []variant{{Solver: "figure"}}, numbers, call)
		if err != nil {
			t.Fatal(err)
		}
		filled, err := pictureFor("pairs", []variant{{Picture: []byte(pairsSource)}}, numbers, nil)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(built, filled) {
			t.Errorf("the instance %v is not drawn from what figure builds", numbers)
		}
	}
	// The whole run, on a line of squared paper, which draws faster than a venn: the statement's
	// picture gets the numbers that solve gets, and the solution's every number, sum among them,
	// which the statement does not show.
	lines := strings.Replace(pairsFigure, "      def figure(numbers):\n", `      def line(numbers):
          return {"kind": "lattice", "width": 5, "height": 1, "lines": [(0, 0, numbers["a"] + numbers["b"], 0)]}

      def lines(numbers):
          return {"kind": "lattice", "width": 5, "height": 1, "lines": [(0, 0, numbers["sum"], 0)]}

      def figure(numbers):
`, 1)
	line, err := pictureFor("pairs", []variant{{Picture: []byte(`{"kind": "lattice", "width": 5, "height": 1, "lines": [[0, 0, 1, 0]]}`)}}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	w, h := line.size()
	solution := fmt.Sprintf("solution:\n  - picture: pictures/pairs-solution-{a}-{b}.png\n    width: %d\n    height: %d\n", w, h)
	folder, module := templateModule(t, map[string]string{"pairs.yaml": solverTemplate(w, h, lines) + solution},
		map[string]string{"pairs.json": `{"solver": "line"}`, "pairs-solution.json": `{"solver": "lines"}`})
	if status, out, errs := runPictures(t, "--template-pictures", folder, "--module", module); status != statusDone || !strings.HasPrefix(out, "port: 14 pictures") || errs != "" {
		t.Errorf("run() = %d, %q, %q", status, out, errs)
	}
}

func TestTemplatePicturesSolverRefuses(t *testing.T) {
	width, height := pairsSize(t)
	program := func(body string) string {
		return "\n  solver:\n    program: |\n" + indented(body) + "    results: [row]\n"
	}
	for name, c := range map[string]struct {
		template string
		source   string
		want     string
	}{
		"no solver":                   {pairsTemplate(width, height), `{"solver": "figure"}`, "a source of its picture is the solver's figure, and its template has no solver"},
		"no such function":            {solverTemplate(width, height, pairsFigure), `{"solver": "drawing"}`, "its solver defines no function drawing"},
		"no function name":            {solverTemplate(width, height, pairsFigure), `{"solver": "a b"}`, `its solver, "a b", is no name of a function`},
		"no string":                   {solverTemplate(width, height, pairsFigure), `{"solver": 5}`, "its solver, 5, is no name of a function"},
		"a null":                      {solverTemplate(width, height, pairsFigure), `{"solver": null}`, "its solver, null, is no name of a function"},
		"a number solve does not get": {solverTemplate(width, height, program("def figure(numbers):\n    return {\"kind\": \"venn\", \"total\": str(numbers[\"sum\"])}\n")), `{"solver": "figure"}`, `key "sum" not in dict`},
		"another member":              {solverTemplate(width, height, pairsFigure), `{"solver": "figure", "when": {}}`, "a source of solver has no member but solver"},
		"a member twice":              {solverTemplate(width, height, pairsFigure), `{"solver": "figure", "solver": "solve"}`, `the member "solver" is named in other letters than the schema's, or twice`},
		"no picture":                  {solverTemplate(width, height, program("def figure(numbers):\n    return [1]\n")), `{"solver": "figure"}`, "its solver's figure gives a list, not a picture"},
		"a member misnamed":           {solverTemplate(width, height, program("def figure(numbers):\n    return {\"Kind\": \"venn\"}\n")), `{"solver": "figure"}`, `its picture: the member "Kind" is named in other letters than the schema's, or twice, in the picture that its solver's figure builds`},
		"no JSON":                     {solverTemplate(width, height, program("def figure(numbers):\n    return {\"kind\": figure}\n")), `{"solver": "figure"}`, "its solver's figure: json.encode"},
		"a load":                      {solverTemplate(width, height, program("load(\"x.star\", \"y\")\n")), `{"solver": "figure"}`, "its solver: "},
		"a changed number":            {solverTemplate(width, height, program("def figure(numbers):\n    numbers[\"a\"] = 1\n")), `{"solver": "figure"}`, "frozen"},
		"a call of itself":            {solverTemplate(width, height, program("def figure(numbers):\n    return figure(numbers)\n")), `{"solver": "figure"}`, "called recursively"},
	} {
		folder, module := templateModule(t, map[string]string{"pairs.yaml": c.template}, map[string]string{"pairs.json": c.source})
		status, _, errs := runPictures(t, "--template-pictures", folder, "--module", module)
		if status != statusCannotRun || !strings.Contains(errs, c.want) {
			t.Errorf("%s: run() = %d, %q", name, status, errs)
		}
	}
	// The budget is the template's: a solver that runs past it is refused, not left to run.
	defer func(steps uint64) { solverSteps = steps }(solverSteps)
	solverSteps = 1000
	folder, module := templateModule(t, map[string]string{"pairs.yaml": solverTemplate(width, height, program("while True:\n    pass\n"))},
		map[string]string{"pairs.json": `{"solver": "figure"}`})
	if status, _, errs := runPictures(t, "--template-pictures", folder, "--module", module); status != statusCannotRun || !strings.Contains(errs, "too many steps") {
		t.Errorf("past the budget: run() = %d, %q", status, errs)
	}
	if _, err := pictureFor("pairs", []variant{{Solver: "figure"}}, pairsInstances(t)[0], nil); err == nil || !strings.Contains(err.Error(), "no solver is given") {
		t.Errorf("a source of solver with no solver given: error = %v", err)
	}
}

// indented indents each line of a program as a block scalar of the solver holds it.
func indented(body string) string {
	var lines []string
	for line := range strings.Lines(body) {
		lines = append(lines, "      "+line)
	}
	return strings.Join(lines, "")
}
