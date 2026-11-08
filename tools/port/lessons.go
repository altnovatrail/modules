package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// The pictures of the lessons (amendment T41 to R12): a folder of pictures in the schema of
// standalone's tasks, one <name>.json each, which the port draws as it draws a task's picture, into
// media/lessons/<name>.png of the module. They come from no clone: a lesson is written by hand, and
// its pictures with it, from the data of the drawings of standalone's site where the lesson keeps
// them (amendment T41 to R14).
const lessonFolder = "media/lessons"

// pictureName is the name of a picture of a lesson: lowercase letters and digits in words joined by
// hyphens, as a local ID is.
var pictureName = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// drawLessonPictures draws every picture of the folder, by the path of its file in the module.
func drawLessonPictures(folder string) (map[string][]byte, error) {
	entries, err := os.ReadDir(folder)
	if err != nil {
		return nil, err
	}
	files := map[string][]byte{}
	for _, entry := range entries {
		name, isJSON := strings.CutSuffix(entry.Name(), ".json")
		if entry.IsDir() || !isJSON || !pictureName.MatchString(name) {
			return nil, fmt.Errorf("%s: a folder of pictures of lessons holds only files <name>.json, the name in lowercase letters, digits and hyphens", filepath.Join(folder, entry.Name()))
		}
		raw, err := os.ReadFile(filepath.Join(folder, entry.Name()))
		if err != nil {
			return nil, err
		}
		picture, err := readPicture(name, raw)
		if err != nil {
			return nil, err
		}
		data, _, _, err := render(picture)
		if err != nil {
			return nil, err
		}
		files[lessonFolder+"/"+name+".png"] = data
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("%s holds no picture", folder)
	}
	return files, nil
}

// staleLessonPictures names the files of the module's media/lessons/ that the port does not draw
// now, from the module's folder, which must be there.
func staleLessonPictures(module string, files map[string][]byte) ([]string, error) {
	if err := isFolder(module); err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(filepath.Join(module, lessonFolder))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var old []string
	for _, entry := range entries {
		path := lessonFolder + "/" + entry.Name()
		if _, draws := files[path]; !draws {
			old = append(old, path)
		}
	}
	return old, nil
}

// runLessonPictures draws the pictures of the folder into the module, or with check compares them.
func runLessonPictures(folder, module string, check bool, stdout, stderr io.Writer) int {
	files, err := drawLessonPictures(folder)
	var old, differ []string
	if err == nil {
		old, err = staleLessonPictures(module, files)
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
		fmt.Fprintf(stdout, "port: %d pictures of lessons written into %s\n", len(files), module)
		return statusDone
	}
	for _, path := range differ {
		fmt.Fprintf(stdout, "port: %s is not what the port draws\n", filepath.Join(module, path))
	}
	if len(differ) > 0 {
		return statusDiffers
	}
	fmt.Fprintf(stdout, "port: the %d pictures of lessons of %s are what the port draws\n", len(files), module)
	return statusDone
}
