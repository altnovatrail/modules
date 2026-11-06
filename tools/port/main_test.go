package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// fakeGo is a stand-in for the go command: it writes its arguments into a file and prints the
// events given, then exits with the status given, as go test -json would.
func fakeGo(t *testing.T, events []string, status int) (command, args string) {
	t.Helper()
	dir := t.TempDir()
	args = filepath.Join(dir, "args")
	eventsFile := filepath.Join(dir, "events")
	writeFile(t, eventsFile, strings.Join(events, "\n")+"\n")
	command = filepath.Join(dir, "go")
	script := "#!/bin/sh\nprintf '%s\\n' \"$@\" > '" + args + "'\ncat '" + eventsFile + "'\necho 'build output' >&2\nexit " +
		string(rune('0'+status)) + "\n"
	writeFile(t, command, script)
	if err := os.Chmod(command, 0o755); err != nil {
		t.Fatal(err)
	}
	return command, args
}

func event(action, id string) string {
	data, _ := json.Marshal(map[string]string{"Action": action, "Test": bench + "/" + id})
	return string(data)
}

func TestProve(t *testing.T) {
	ids := []string{"lamps-1", "lamps-2"}
	command, args := fakeGo(t, []string{"not json", event("run", "lamps-1"), event("pass", "lamps-1"),
		event("pass", "lamps-2"), `{"Action": "pass", "Test": "` + bench + `"}`}, 0)
	if err := prove(command, "/clone", ids); err != nil {
		t.Fatalf("prove() error = %v", err)
	}
	got, err := os.ReadFile(args)
	if err != nil {
		t.Fatal(err)
	}
	want := "-C\n/clone\ntest\n./content\n-run\n^" + bench + "$/^(lamps-1|lamps-2)$\n-count=1\n-json\n"
	if string(got) != want {
		t.Errorf("go was called with\n%s\nwant\n%s", got, want)
	}
}

func TestProveRefuses(t *testing.T) {
	ids := []string{"lamps-1", "lamps-2"}
	for name, test := range map[string]struct {
		events  []string
		status  int
		message string
	}{
		"a solver that proves another answer": {[]string{event("pass", "lamps-1"),
			`{"Action": "output", "Test": "` + bench + `/lamps-2", "Output": "    the solver proves A, the task claims C\n"}`,
			event("fail", "lamps-2")}, 1, "the solver proves A, the task claims C"},
		"a bench that does not compile": {[]string{`{"Action": "build-output", "Output": "content/x.go:1: syntax error\n"}`,
			`{"Action": "output", "Output": "FAIL content [build failed]\n"}`}, 1, "content/x.go:1: syntax error"},
		"a task the bench did not run": {[]string{event("pass", "lamps-1")}, 0, "lamps-2: the bench did not run it"},
		"go test failed though every task passed": {[]string{event("pass", "lamps-1"), event("pass", "lamps-2")}, 2,
			"go test failed"},
	} {
		command, _ := fakeGo(t, test.events, test.status)
		err := prove(command, "/clone", ids)
		if err == nil || !strings.Contains(err.Error(), test.message) || !strings.Contains(err.Error(), "build output") {
			t.Errorf("%s: error = %v, want one that says %q and what go wrote", name, err, test.message)
		}
	}
	if err := prove(filepath.Join(t.TempDir(), "no-go"), "/clone", ids); err == nil {
		t.Error("no go command: no error")
	}
}

// A made-up standalone, batch and overlays, and an empty module, for run.
func setUp(t *testing.T) (args []string, module string) {
	t.Helper()
	pictured := example("lamps-2")
	pictured.Picture = json.RawMessage(rowJSON)
	clone, commit := standaloneWith(t, map[string][]Example{"counting.gaps.json": {example("lamps-1"), pictured}})
	overlays := t.TempDir()
	writeFile(t, filepath.Join(overlays, "lamps-2.yaml"), "picture_alt: Lamps in a row.\n")
	module = t.TempDir()
	return []string{"--standalone", clone, "--batch", batchFile(t, commit, "lamps-1", "lamps-2"), "--overlay", overlays, "--module", module}, module
}

func runPort(args ...string) (int, string, string) {
	var stdout, stderr bytes.Buffer
	status := run(args, &stdout, &stderr)
	return status, stdout.String(), stderr.String()
}

func TestRunWritesThenChecks(t *testing.T) {
	args, module := setUp(t)
	status, stdout, stderr := runPort(args...)
	if status != statusDone || !strings.Contains(stdout, "2 tasks written") {
		t.Fatalf("run() = %d, %s%s", status, stdout, stderr)
	}
	for _, file := range []string{"tasks/lamps-1.yaml", "tasks/lamps-2.yaml", "media/pictures/lamps-2.png"} {
		if _, err := os.Stat(filepath.Join(module, file)); err != nil {
			t.Errorf("%s is not written: %v", file, err)
		}
	}

	if status, stdout, _ := runPort(append(args, "--check")...); status != statusDone || !strings.Contains(stdout, "are what the port writes") {
		t.Errorf("--check after a run: %d, %s", status, stdout)
	}

	writeFile(t, filepath.Join(module, "tasks", "lamps-1.yaml"), "id: lamps-1\n")
	if err := os.Remove(filepath.Join(module, "media", "pictures", "lamps-2.png")); err != nil {
		t.Fatal(err)
	}
	status, stdout, _ = runPort(append(args, "--check")...)
	if status != statusDiffers || !strings.Contains(stdout, "lamps-1.yaml is not what the port writes") ||
		!strings.Contains(stdout, "lamps-2.png is not what the port writes") {
		t.Errorf("--check after an edit: %d, %s", status, stdout)
	}
}

// A picture the port wrote once, and would not write now, is stale: --check names it, and a run
// removes it.
func TestRunStale(t *testing.T) {
	args, module := setUp(t)
	if status, _, stderr := runPort(args...); status != statusDone {
		t.Fatalf("run() = %d, %s", status, stderr)
	}
	old := filepath.Join(module, "media", "pictures", "lamps-1-solution.png")
	writeFile(t, old, "an old picture")
	status, stdout, _ := runPort(append(args, "--check")...)
	if status != statusDiffers || !strings.Contains(stdout, "lamps-1-solution.png is not what the port writes") {
		t.Errorf("--check with a stale picture: %d, %s", status, stdout)
	}
	if status, _, stderr := runPort(args...); status != statusDone {
		t.Fatalf("run() = %d, %s", status, stderr)
	}
	if _, err := os.Stat(old); !os.IsNotExist(err) {
		t.Errorf("the stale picture stays: %v", err)
	}
	if status, _, _ := runPort(append(args, "--check")...); status != statusDone {
		t.Errorf("--check after the run: %d", status)
	}
}

func TestRunProves(t *testing.T) {
	args, _ := setUp(t)
	command, _ := fakeGo(t, []string{event("pass", "lamps-1"), event("pass", "lamps-2")}, 0)
	status, stdout, stderr := runPort(append(args, "--prove", "--go", command)...)
	if status != statusDone || !strings.Contains(stdout, "proves the answers of all 2 tasks") {
		t.Errorf("run() = %d, %s%s", status, stdout, stderr)
	}

	failing, _ := fakeGo(t, []string{event("pass", "lamps-1"), event("fail", "lamps-2")}, 1)
	status, stdout, stderr = runPort(append(args, "--prove", "--go", failing)...)
	if status != statusCannotRun || !strings.Contains(stderr, "lamps-2: the solver does not prove") || strings.Contains(stdout, "written") {
		t.Errorf("a failing bench: %d, %s%s", status, stdout, stderr)
	}
}

func TestRunRefuses(t *testing.T) {
	args, module := setUp(t)
	if status, _, stderr := runPort(); status != statusCannotRun || !strings.Contains(stderr, "usage: port") {
		t.Errorf("no arguments: %d, %s", status, stderr)
	}
	if status, _, stderr := runPort(append(args, "extra")...); status != statusCannotRun || !strings.Contains(stderr, "usage") {
		t.Errorf("an extra argument: %d, %s", status, stderr)
	}

	missing := append([]string(nil), args...)
	missing[3] = batchFile(t, git(t, args[1], "rev-parse", "HEAD"), "lamps-1", "lamps-2", "lamps-9")
	if status, _, stderr := runPort(missing...); status != statusCannotRun || !strings.Contains(stderr, "lamps-9: standalone has no such task") {
		t.Errorf("a task standalone lacks: %d, %s", status, stderr)
	}

	noAlt := append([]string(nil), args...)
	noAlt[3] = batchFile(t, git(t, args[1], "rev-parse", "HEAD"), "lamps-2")
	noAlt[5] = t.TempDir()
	if status, _, stderr := runPort(noAlt...); status != statusCannotRun || !strings.Contains(stderr, "no alt") {
		t.Errorf("a picture without alt: %d, %s", status, stderr)
	}

	badOverlay := append([]string(nil), args...)
	badOverlay[5] = filepath.Join(t.TempDir(), "file")
	writeFile(t, badOverlay[5], "")
	if status, _, _ := runPort(badOverlay...); status != statusCannotRun {
		t.Errorf("an overlay folder that is a file: %d", status)
	}

	badClone := append([]string(nil), args...)
	badClone[1] = t.TempDir()
	if status, _, stderr := runPort(badClone...); status != statusCannotRun || !strings.Contains(stderr, "no clone of standalone that git reads") {
		t.Errorf("no clone: %d, %s", status, stderr)
	}

	otherCommit := append([]string(nil), args...)
	otherCommit[3] = batchFile(t, testCommit, "lamps-1", "lamps-2")
	if status, _, stderr := runPort(otherCommit...); status != statusCannotRun || !strings.Contains(stderr, "and the batch is of "+testCommit) {
		t.Errorf("a clone at another commit: %d, %s", status, stderr)
	}

	writeFile(t, filepath.Join(args[1], "content", "examples", "extra.json"), "[]")
	if status, _, stderr := runPort(args...); status != statusCannotRun || !strings.Contains(stderr, "has changes of its own") {
		t.Errorf("a clone with changes: %d, %s", status, stderr)
	}
	if err := os.Remove(filepath.Join(args[1], "content", "examples", "extra.json")); err != nil {
		t.Fatal(err)
	}

	badBatch := append([]string(nil), args...)
	badBatch[3] = filepath.Join(t.TempDir(), "none.yaml")
	if status, _, _ := runPort(badBatch...); status != statusCannotRun {
		t.Errorf("no batch: %d", status)
	}

	file := filepath.Join(module, "a-file")
	writeFile(t, file, "")
	notDir := append([]string(nil), args...)
	notDir[7] = file
	if status, _, stderr := runPort(notDir...); status != statusCannotRun || !strings.Contains(stderr, "is no folder") {
		t.Errorf("a module that is a file: %d, %s", status, stderr)
	}
	if status, _, _ := runPort(append(notDir, "--check")...); status != statusCannotRun {
		t.Errorf("--check of a module that is a file: %d", status)
	}
	gone := append([]string(nil), args...)
	gone[7] = filepath.Join(module, "none")
	if status, _, _ := runPort(gone...); status != statusCannotRun {
		t.Errorf("no module: %d", status)
	}
	blocked := append([]string(nil), args...)
	blocked[7] = t.TempDir()
	writeFile(t, filepath.Join(blocked[7], "tasks"), "")
	if status, _, _ := runPort(blocked...); status != statusCannotRun {
		t.Errorf("tasks/ that is a file: %d", status)
	}
}
