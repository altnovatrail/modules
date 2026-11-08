package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const plusPicture = `{"kind": "lattice", "width": 3, "height": 3, "cells": [[1, 0], [0, 1], [1, 1], [2, 1], [1, 2]]}`
const stripPicture = `{"kind": "lattice", "width": 3, "height": 1, "lines": [[0, 0, 3, 0], [0, 1, 3, 1]]}`

// pictureFolder is a folder of pictures of lessons with the files given, and a module beside it.
func pictureFolder(t *testing.T, files map[string]string) (folder, module string) {
	t.Helper()
	root := t.TempDir()
	folder, module = filepath.Join(root, "pictures"), filepath.Join(root, "module")
	for _, dir := range []string{folder, module} {
		if err := os.Mkdir(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	for name, content := range files {
		writeFile(t, filepath.Join(folder, name), content)
	}
	return folder, module
}

func drawnAs(t *testing.T, raw string) []byte {
	t.Helper()
	picture, err := readPicture("test", []byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	data, _, _, err := render(picture)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func runPictures(t *testing.T, args ...string) (int, string, string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	status := run(args, &stdout, &stderr)
	return status, stdout.String(), stderr.String()
}

func TestLessonPicturesWriteAndCheck(t *testing.T) {
	folder, module := pictureFolder(t, map[string]string{"a-cross.json": plusPicture, "row-2.json": stripPicture})
	writeFile(t, filepath.Join(module, lessonFolder, "old.png"), "old")
	status, out, errs := runPictures(t, "--pictures", folder, "--module", module)
	if status != statusDone || out != "port: 2 pictures of lessons written into "+module+"\n" || errs != "" {
		t.Fatalf("run() = %d, %q, %q", status, out, errs)
	}
	for name, raw := range map[string]string{"a-cross": plusPicture, "row-2": stripPicture} {
		got, err := os.ReadFile(filepath.Join(module, lessonFolder, name+".png"))
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(got, drawnAs(t, raw)) {
			t.Errorf("%s.png is not the picture the port draws", name)
		}
	}
	if _, err := os.Stat(filepath.Join(module, lessonFolder, "old.png")); !os.IsNotExist(err) {
		t.Errorf("old.png, drawn from no file, is still there: %v", err)
	}
	status, out, _ = runPictures(t, "--pictures", folder, "--module", module, "--check")
	if status != statusDone || out != "port: the 2 pictures of lessons of "+module+" are what the port draws\n" {
		t.Errorf("--check on the pictures just drawn = %d, %q", status, out)
	}
}

func TestLessonPicturesCheckDiffers(t *testing.T) {
	folder, module := pictureFolder(t, map[string]string{"a-cross.json": plusPicture, "row-2.json": stripPicture})
	writeFile(t, filepath.Join(module, lessonFolder, "a-cross.png"), "another picture")
	writeFile(t, filepath.Join(module, lessonFolder, "stale.png"), "stale")
	status, out, errs := runPictures(t, "--pictures", folder, "--module", module, "--check")
	want := []string{"a-cross.png", "row-2.png", "stale.png"}
	var lines []string
	for _, name := range want {
		lines = append(lines, "port: "+filepath.Join(module, lessonFolder, name)+" is not what the port draws")
	}
	if status != statusDiffers || out != strings.Join(lines, "\n")+"\n" || errs != "" {
		t.Errorf("run() = %d, %q, %q", status, out, errs)
	}
	if got, _ := os.ReadFile(filepath.Join(module, lessonFolder, "a-cross.png")); string(got) != "another picture" {
		t.Errorf("--check wrote into the module")
	}
}

func TestLessonPicturesCheckWithNoFolder(t *testing.T) {
	folder, module := pictureFolder(t, map[string]string{"row-2.json": stripPicture})
	status, out, _ := runPictures(t, "--pictures", folder, "--module", module, "--check")
	if status != statusDiffers || !strings.Contains(out, "row-2.png is not what the port draws") {
		t.Errorf("run() = %d, %q", status, out)
	}
}

func TestLessonPicturesRefuse(t *testing.T) {
	for name, test := range map[string]struct {
		files   map[string]string
		message string
	}{
		"no picture":         {map[string]string{}, "holds no picture"},
		"a file of no JSON":  {map[string]string{"notes.txt": "x"}, "holds only files <name>.json"},
		"a name in capitals": {map[string]string{"Cross.json": plusPicture}, "holds only files <name>.json"},
		"a kind not drawn":   {map[string]string{"cube.json": `{"kind": "cube"}`}, `of the kind "cube"`},
		"a broken picture":   {map[string]string{"bad.json": `{"kind": "lattice", "width": 9, "height": 1}`}, "1 to 8 cells"},
		"no JSON":            {map[string]string{"bad.json": `{`}, "bad: its picture"},
	} {
		t.Run(name, func(t *testing.T) {
			folder, module := pictureFolder(t, test.files)
			status, out, errs := runPictures(t, "--pictures", folder, "--module", module)
			if status != statusCannotRun || out != "" || !strings.Contains(errs, test.message) {
				t.Errorf("run() = %d, %q, %q, want the message %q", status, out, errs, test.message)
			}
		})
	}
}

func TestLessonPicturesWithNoModule(t *testing.T) {
	folder, module := pictureFolder(t, map[string]string{"row-2.json": stripPicture})
	file := filepath.Join(module, "file")
	writeFile(t, file, "x")
	for name, target := range map[string]string{"a file": file, "nothing": filepath.Join(module, "none")} {
		for _, check := range [][]string{nil, {"--check"}} {
			status, out, errs := runPictures(t, append([]string{"--pictures", folder, "--module", target}, check...)...)
			if status != statusCannotRun || out != "" || errs == "" {
				t.Errorf("%s as the module, %v: run() = %d, %q, %q", name, check, status, out, errs)
			}
		}
	}
	status, _, errs := runPictures(t, "--pictures", filepath.Join(folder, "none"), "--module", module)
	if status != statusCannotRun || errs == "" {
		t.Errorf("no folder of pictures: run() = %d, %q", status, errs)
	}
}

func TestLessonPicturesFlags(t *testing.T) {
	folder, module := pictureFolder(t, map[string]string{"row-2.json": stripPicture})
	for name, args := range map[string][]string{
		"with a clone":    {"--standalone", "/clone"},
		"with a batch":    {"--batch", "b.yaml", "--overlay", "o", "--standalone", "/clone"},
		"with a survey":   {"--survey", "--commit", "abc", "--except", "e", "--standalone", "/clone"},
		"with --prove":    {"--prove"},
		"with no module":  {"--module", ""},
		"with a leftover": {"extra"},
	} {
		status, _, errs := runPictures(t, append([]string{"--pictures", folder, "--module", module}, args...)...)
		if status != statusCannotRun || !strings.HasPrefix(errs, "usage:") {
			t.Errorf("%s: run() = %d, %q", name, status, errs)
		}
	}
}

func TestLessonPicturesFailOnTheFileSystem(t *testing.T) {
	expect := func(name string, status int, errs string) {
		t.Helper()
		if status != statusCannotRun || errs == "" {
			t.Errorf("%s: run() = %d, %q", name, status, errs)
		}
	}
	folder, module := pictureFolder(t, map[string]string{"row-2.json": stripPicture})
	writeFile(t, filepath.Join(module, "media", "lessons"), "a file, not a folder")
	status, _, errs := runPictures(t, "--pictures", folder, "--module", module, "--check")
	expect("media/lessons as a file, --check", status, errs)
	status, _, errs = runPictures(t, "--pictures", folder, "--module", module)
	expect("media/lessons as a file", status, errs)

	folder, module = pictureFolder(t, map[string]string{"row-2.json": stripPicture})
	writeFile(t, filepath.Join(module, lessonFolder, "row-2.png", "inside"), "a folder where the picture goes")
	status, _, errs = runPictures(t, "--pictures", folder, "--module", module, "--check")
	expect("a folder where the picture goes, --check", status, errs)
	status, _, errs = runPictures(t, "--pictures", folder, "--module", module)
	expect("a folder where the picture goes", status, errs)

	folder, module = pictureFolder(t, map[string]string{"row-2.json": stripPicture})
	writeFile(t, filepath.Join(module, lessonFolder, "old", "inside"), "a folder that is no picture")
	status, _, errs = runPictures(t, "--pictures", folder, "--module", module)
	expect("a folder that is no picture", status, errs)

	folder, module = pictureFolder(t, map[string]string{"row-2.json": stripPicture})
	if err := os.Chmod(filepath.Join(folder, "row-2.json"), 0); err != nil {
		t.Fatal(err)
	}
	if _, err := os.ReadFile(filepath.Join(folder, "row-2.json")); err == nil {
		t.Skip("the file system lets this user read a file of mode 0")
	}
	status, _, errs = runPictures(t, "--pictures", folder, "--module", module)
	expect("a picture that cannot be read", status, errs)
}
