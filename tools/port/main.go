package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"slices"
)

// The exit statuses: 0 when the port wrote the tasks or found them as it would write them, 1 when
// --check found a file of the module other than the port would write, or a survey refused a task,
// and 2 when it could not do its work.
const (
	statusDone      = 0
	statusDiffers   = 1
	statusCannotRun = 2
)

const usage = `usage: port --standalone <clone> --batch <file> --overlay <folder> --module <folder> [--prove] [--check] [--go <command>] [--git <command>]
       port --standalone <clone> --survey --commit <sha> --except <folder> --module <folder> [--prove] [--go <command>] [--git <command>]
       port --pictures <folder> --module <folder> [--check]
       port --template-pictures <folder> --module <folder> [--check]

Writes tasks/<id>.yaml and the pictures of media/pictures/ of the module for every task of the
batch, from a clone of mathtrail-standalone at the batch's commit and the overlays. --prove first
runs standalone's bench on the batch's tasks, and writes nothing unless each passes. --check writes
nothing, and says which files of the module differ from what the port would write.

--survey writes every task of the clone at the commit that no batch of the folder --except names,
with no overlays, into a module of its own, and says what became of each: written, by the kind of
its answer, waiting for a picture of a kind not drawn yet, or refused. --prove runs the bench on the
tasks written.

--pictures draws every picture of the folder, <name>.json in the schema of standalone's pictures,
into media/lessons/<name>.png of the module, with no clone, and removes the module's other files of
media/lessons/. --check writes nothing, and says which of them differ from what the port would draw.

--template-pictures draws the picture of every instance of every template of the module whose
picture path holds placeholders, from the source <name>.json of the folder, the path without its
placeholders, into media/pictures/ of the module, and removes the module's other files of those
paths. --check writes nothing, and says which of them differ from what the port would draw.`

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

// run is the command, apart from the process: its arguments, and where it writes.
func run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("port", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	standalone := flags.String("standalone", "", "")
	batchFile := flags.String("batch", "", "")
	overlayDir := flags.String("overlay", "", "")
	module := flags.String("module", "", "")
	check := flags.Bool("check", false, "")
	proveFirst := flags.Bool("prove", false, "")
	surveying := flags.Bool("survey", false, "")
	commitOf := flags.String("commit", "", "")
	except := flags.String("except", "", "")
	pictures := flags.String("pictures", "", "")
	templatePictures := flags.String("template-pictures", "", "")
	goCommand := flags.String("go", "go", "")
	gitCommand := flags.String("git", "git", "")
	err := flags.Parse(args)
	ordinary := *batchFile != "" && *overlayDir != "" && *commitOf == "" && *except == "" && *pictures == "" && *templatePictures == ""
	surveyed := *commitOf != "" && *except != "" && *batchFile == "" && *overlayDir == "" && !*check && *pictures == "" && *templatePictures == ""
	alone := *standalone == "" && *batchFile == "" && *overlayDir == "" && *commitOf == "" && *except == "" && !*surveying && !*proveFirst
	drawing := *pictures != "" && *templatePictures == "" && alone
	drawingInstances := *templatePictures != "" && *pictures == "" && alone
	if err != nil || flags.NArg() != 0 || *module == "" || (!drawing && !drawingInstances && (*standalone == "" ||
		(*surveying && !surveyed) || (!*surveying && !ordinary))) {
		fmt.Fprintln(stderr, usage)
		return statusCannotRun
	}
	if drawing {
		return runLessonPictures(*pictures, *module, *check, stdout, stderr)
	}
	if drawingInstances {
		return runTemplatePictures(*templatePictures, *module, *check, stdout, stderr)
	}
	if *surveying {
		return runSurvey(survey{
			gitCommand: *gitCommand, goCommand: *goCommand, standalone: *standalone,
			commit: *commitOf, except: *except, module: *module, prove: *proveFirst,
		}, stdout, stderr)
	}
	converted, err := portAll(*gitCommand, *standalone, *batchFile, *overlayDir)
	if err != nil {
		fmt.Fprintf(stderr, "port: %v\n", err)
		return statusCannotRun
	}
	if *proveFirst {
		ids := make([]string, len(converted))
		for i, c := range converted {
			ids[i] = c.ID
		}
		if err := prove(*goCommand, *standalone, ids); err != nil {
			fmt.Fprintf(stderr, "port: %v\n", err)
			return statusCannotRun
		}
		fmt.Fprintf(stdout, "port: the bench of standalone proves the answers of all %d tasks\n", len(ids))
	}
	for _, c := range converted {
		for _, note := range c.Notes {
			fmt.Fprintf(stdout, "port: %s: %s\n", c.ID, note)
		}
	}
	if *check {
		differ, err := differing(*module, converted)
		if err != nil {
			fmt.Fprintf(stderr, "port: %v\n", err)
			return statusCannotRun
		}
		for _, path := range differ {
			fmt.Fprintf(stdout, "port: %s is not what the port writes\n", filepath.Join(*module, path))
		}
		if len(differ) > 0 {
			return statusDiffers
		}
		fmt.Fprintf(stdout, "port: the %d tasks of %s are what the port writes\n", len(converted), *module)
		return statusDone
	}
	if err := write(*module, converted); err != nil {
		fmt.Fprintf(stderr, "port: %v\n", err)
		return statusCannotRun
	}
	fmt.Fprintf(stdout, "port: %d tasks written into %s\n", len(converted), *module)
	return statusDone
}

// portAll converts every task of a batch, in the order of their IDs, or fails on the first that
// cannot be.
func portAll(gitCommand, standalone, batchFile, overlayDir string) ([]Converted, error) {
	batch, err := readBatch(batchFile)
	if err != nil {
		return nil, err
	}
	examples, topics, err := readClone(gitCommand, standalone, batch.Standalone)
	if err != nil {
		return nil, err
	}
	overlays, err := readOverlays(overlayDir, batch)
	if err != nil {
		return nil, err
	}
	batch.overlays = filepath.ToSlash(filepath.Clean(overlayDir))
	var converted []Converted
	for _, id := range slices.Sorted(maps.Keys(batch.Tasks)) {
		ex, has := examples[id]
		if !has {
			return nil, fmt.Errorf("%s: %w", id, errNoExample)
		}
		c, err := convertOf(standalone, ex, topics, batch, batch.Tasks[id], overlays[id])
		if err != nil {
			return nil, err
		}
		converted = append(converted, c)
	}
	return converted, nil
}

// readClone holds a clone of standalone to a commit, and reads its tasks and its topics.
func readClone(gitCommand, standalone, commit string) (map[string]Example, map[string]Topic, error) {
	if err := atCommit(gitCommand, standalone, commit); err != nil {
		return nil, nil, err
	}
	examples, err := readExamples(standalone)
	if err != nil {
		return nil, nil, err
	}
	topics, err := readTopics(standalone)
	if err != nil {
		return nil, nil, err
	}
	return examples, topics, nil
}

// convertOf holds a task to its topic, reads whether its solver reads the picture, and converts it.
func convertOf(clone string, ex Example, topics map[string]Topic, batch Batch, task BatchTask, overlay Overlay) (Converted, error) {
	if err := ofTopic(ex, topics); err != nil {
		return Converted{}, err
	}
	var err error
	if ex.fromPicture, err = solvedFromPicture(clone, ex.ID); err != nil {
		return Converted{}, err
	}
	return convert(ex, batch, task, overlay, "content/examples/"+ex.Topic+".json")
}

// paths are the files of the converted tasks, by their paths from the module's folder.
func paths(converted []Converted) map[string][]byte {
	files := map[string][]byte{}
	for _, c := range converted {
		maps.Copy(files, c.Files)
	}
	return files
}

// owned are the files the port may write for a task: the task, and the pictures of its statement
// and of its solution.
func owned(id string) []string {
	return []string{"tasks/" + id + ".yaml", "media/pictures/" + id + ".png", "media/pictures/" + id + "-solution.png"}
}

// stale names the files of the module that the port wrote for a task of the batch once, and would
// not write now: a picture the task no longer has. Tasks of other batches are not the port's to
// judge here.
func stale(module string, converted []Converted) ([]string, error) {
	var old []string
	for _, c := range converted {
		for _, path := range owned(c.ID) {
			if _, writes := c.Files[path]; writes {
				continue
			}
			_, err := os.Lstat(filepath.Join(module, path))
			if err == nil {
				old = append(old, path)
			} else if !errors.Is(err, fs.ErrNotExist) {
				return nil, err
			}
		}
	}
	return old, nil
}

// differing names the files of the module that are not what the port writes, missing and stale
// ones too.
func differing(module string, converted []Converted) ([]string, error) {
	old, err := stale(module, converted)
	if err != nil {
		return nil, err
	}
	files := paths(converted)
	return filesDiffer(module, files, old)
}

// filesDiffer names the files of the module that are not the files given, missing ones too, and the
// stale ones, in order.
func filesDiffer(module string, files map[string][]byte, old []string) ([]string, error) {
	differ := slices.Clone(old)
	for _, path := range slices.Sorted(maps.Keys(files)) {
		have, err := os.ReadFile(filepath.Join(module, path))
		if errors.Is(err, fs.ErrNotExist) {
			differ = append(differ, path)
			continue
		}
		if err != nil {
			return nil, err
		}
		if !bytes.Equal(have, files[path]) {
			differ = append(differ, path)
		}
	}
	slices.Sort(differ)
	return differ, nil
}

// write writes the files of the converted tasks into the module's folder, which must be there, and
// removes the stale ones.
func write(module string, converted []Converted) error {
	if err := isFolder(module); err != nil {
		return err
	}
	old, err := stale(module, converted)
	if err != nil {
		return err
	}
	files := paths(converted)
	return placeFiles(module, files, old)
}

// isFolder refuses a module that is no folder.
func isFolder(module string) error {
	info, err := os.Stat(module)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is no folder", module)
	}
	return nil
}

// placeFiles removes the stale files of the module and writes the files given, with their folders.
func placeFiles(module string, files map[string][]byte, old []string) error {
	for _, path := range old {
		if err := os.Remove(filepath.Join(module, path)); err != nil {
			return err
		}
	}
	for _, path := range slices.Sorted(maps.Keys(files)) {
		target := filepath.Join(module, path)
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(target, files[path], 0o644); err != nil {
			return err
		}
	}
	return nil
}
