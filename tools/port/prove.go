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

// bench is the test of standalone, content/solver_test.go, that runs the solver of each reference
// task as a submitted solver is run — its sandbox and dialect, two runs under rotated letters — in a
// subtest named by the task's ID. The runs must agree on exactly one letter, so the solver proves
// the key and that no trap value is right. The one trap value that is no option of
// standalone, a time twelve hours off the key as a clock face shows it, no solver sees: it is never
// right.
const bench = "TestEveryReferenceTaskProvesItsOwnAnswer"

// prove fails unless the bench of the clone passed for each task given. go fetches the clone's own
// modules, none of which enters this repository or its build.
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

	// A refusal carries the output of a subtest that failed and of the build, which go test -json
	// writes as events of their own, so that a bench that does not compile says why.
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
