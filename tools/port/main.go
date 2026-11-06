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
// --check found a file of the module other than the port would write, and 2 when it could not do
// its work.
const (
	statusDone      = 0
	statusDiffers   = 1
	statusCannotRun = 2
)

const usage = `usage: port --standalone <clone> --batch <file> --overlay <folder> --module <folder> [--prove] [--check] [--go <command>] [--git <command>]

Writes tasks/<id>.yaml and the pictures of media/pictures/ of the module for every task of the
batch, from a clone of mathtrail-standalone at the batch's commit and the overlays. --prove first
runs standalone's bench on the batch's tasks, and writes nothing unless each passes. --check writes
nothing, and says which files of the module differ from what the port would write.`

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
	goCommand := flags.String("go", "go", "")
	gitCommand := flags.String("git", "git", "")
	if err := flags.Parse(args); err != nil || flags.NArg() != 0 ||
		*standalone == "" || *batchFile == "" || *overlayDir == "" || *module == "" {
		fmt.Fprintln(stderr, usage)
		return statusCannotRun
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
	if err := atCommit(gitCommand, standalone, batch.Standalone); err != nil {
		return nil, err
	}
	examples, err := readExamples(standalone)
	if err != nil {
		return nil, err
	}
	overlays, err := readOverlays(overlayDir, batch)
	if err != nil {
		return nil, err
	}
	var converted []Converted
	for _, id := range slices.Sorted(maps.Keys(batch.Tasks)) {
		ex, has := examples[id]
		if !has {
			return nil, fmt.Errorf("%s: %w", id, errNoExample)
		}
		topicFile := "content/examples/" + ex.Topic + ".json"
		c, err := convert(ex, batch, batch.Tasks[id], overlays[id], topicFile)
		if err != nil {
			return nil, err
		}
		converted = append(converted, c)
	}
	return converted, nil
}

// paths are the files of the converted tasks, by their paths from the module's folder, in order.
func paths(converted []Converted) ([]string, map[string][]byte) {
	files := map[string][]byte{}
	for _, c := range converted {
		maps.Copy(files, c.Files)
	}
	return slices.Sorted(maps.Keys(files)), files
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
	order, files := paths(converted)
	differ, err := stale(module, converted)
	if err != nil {
		return nil, err
	}
	for _, path := range order {
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
	info, err := os.Stat(module)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is no folder", module)
	}
	old, err := stale(module, converted)
	if err != nil {
		return err
	}
	for _, path := range old {
		if err := os.Remove(filepath.Join(module, path)); err != nil {
			return err
		}
	}
	order, files := paths(converted)
	for _, path := range order {
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
