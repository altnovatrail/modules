package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// bench is the test of standalone that runs the solver of each reference task the way a submitted
// solver is run — its sandbox, its dialect and helpers, two runs under rotated letters — and holds
// it to the task's correct answer: a subtest each, named by the task's ID (content/solver_test.go
// of standalone). The runs must agree on exactly one letter, so the solver proves the key and that
// no other option, none of the trap values, is right (R05; R12).
const bench = "TestEveryReferenceTaskProvesItsOwnAnswer"

// prove runs the bench of a clone of standalone on the tasks given, with the go command named, and
// fails unless each of them passed. The clone is standalone's own module: go fetches its modules,
// none of which enters this repository or its build.
func prove(goCommand, clone string, ids []string) error {
	quoted := make([]string, len(ids))
	for i, id := range ids {
		quoted[i] = regexp.QuoteMeta(id)
	}
	pattern := fmt.Sprintf("^%s$/^(%s)$", bench, strings.Join(quoted, "|"))
	command := exec.Command(goCommand, "-C", clone, "test", "./content", "-run", pattern, "-count=1", "-json")
	var stderr bytes.Buffer
	command.Stderr = &stderr
	out, runErr := command.Output()

	// What go says goes with a refusal: the output of a subtest that failed, and the output of the
	// build, which go test -json writes as events of their own, so that a bench that does not
	// compile says why.
	passed, failed := map[string]bool{}, map[string]bool{}
	said := map[string][]string{}
	var build []string
	scanner := bufio.NewScanner(bytes.NewReader(out))
	scanner.Buffer(make([]byte, 0, 1<<16), 1<<24)
	for scanner.Scan() {
		var event struct {
			Action string
			Test   string
			Output string
		}
		if json.Unmarshal(scanner.Bytes(), &event) != nil {
			continue
		}
		id, sub := strings.CutPrefix(event.Test, bench+"/")
		switch {
		case event.Action == "build-output", event.Action == "output" && event.Test == "":
			build = append(build, strings.TrimRight(event.Output, "\n"))
		case !sub:
		case event.Action == "pass":
			passed[id] = true
		case event.Action == "fail":
			failed[id] = true
		case event.Action == "output":
			said[id] = append(said[id], strings.TrimSpace(event.Output))
		}
	}
	var problems []string
	for _, id := range ids {
		switch {
		case failed[id]:
			problems = append(problems, id+": the solver does not prove the task's answer")
			for _, line := range said[id] {
				problems = append(problems, "    "+line)
			}
		case !passed[id]:
			problems = append(problems, id+": the bench did not run it")
		}
	}
	if runErr != nil {
		problems = append(problems, fmt.Sprintf("go test failed: %v", runErr))
	}
	if len(problems) == 0 {
		return nil
	}
	problems = append(problems, build...)
	if detail := strings.TrimSpace(stderr.String()); detail != "" {
		problems = append(problems, detail)
	}
	return fmt.Errorf("the bench of standalone:\n%s", strings.Join(problems, "\n"))
}
