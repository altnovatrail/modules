// Command port carries tasks of mathtrail-standalone into a module of the format 0.3 (R05, R07,
// R12): each task of a batch becomes the file tasks/<id>.yaml of the module, by the table «что во
// что» of docs/architecture.md, section 3.5, and a picture of the kind row a PNG in media/. What a
// person writes — the givens, the unknown, the simpler case, the alt of a picture, explanations
// written anew — comes from an overlay file beside this code, so that a run writes the same bytes
// again and loses nothing. The key and the trap values are held to the solver of each task by
// standalone's own bench, run in a clone of standalone: no solver goes into the module.
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

// letters are the options of a task of standalone, in the order it shows them.
var letters = []string{"A", "B", "C", "D", "E"}

// Example is a reference task of standalone, as content/examples/<topic>.json holds it at the
// commit the batch names. Every field it may have is here, so that a field this code does not know
// stops the port rather than be dropped unsaid.
type Example struct {
	ID              string                `json:"id"`
	Topic           string                `json:"topic"`
	GradeLevel      string                `json:"grade_level"`
	Difficulty      int                   `json:"difficulty"`
	Question        string                `json:"question"`
	Picture         json.RawMessage       `json:"picture,omitempty"`
	Options         map[string]string     `json:"options"`
	CorrectAnswer   string                `json:"correct_answer"`
	Solution        string                `json:"solution"`
	SolutionTotal   string                `json:"solution_total,omitempty"`
	SolutionPicture json.RawMessage       `json:"solution_picture,omitempty"`
	Distractors     map[string]Distractor `json:"distractors"`
	// Hint is not carried: the format has no field for it, and the ladder of hints is built from
	// the traps, the givens, the techniques and the simpler case (format/README.md, section 19).
	Hint string `json:"hint,omitempty"`
	// Idea is the place of the task's idea in its topic's list of ideas, at 42b4408 on a task solved
	// from the picture. It is not carried: the ideas are not (docs/architecture.md, section 3.5).
	Idea int `json:"idea,omitempty"`
	// fromPicture says that the task's solver reads its picture, solve(options, picture): its facts
	// stand on the picture alone. It is read from content/examples/solvers/<id>.star.
	fromPicture bool
}

// readsPicture is the head of a solver that reads its task's picture: solve with a second parameter,
// solve(options, picture), however it is spaced or named.
var readsPicture = regexp.MustCompile(`(?m)^def\s+solve\s*\(\s*\w+\s*,\s*\w+`)

// solvedFromPicture says whether the solver of a task, content/examples/solvers/<id>.star of a
// clone of standalone, reads the task's picture. A task with no such file reads none: the bench,
// which runs the solver, says that the file is missing.
func solvedFromPicture(clone, id string) (bool, error) {
	data, err := os.ReadFile(filepath.Join(clone, "content", "examples", "solvers", id+".star"))
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return readsPicture.Match(data), nil
}

// Topic is a topic of standalone, as content/catalogs/topics.json holds it at the commit the batch
// names. The port reads it to hold each task to its topic, and carries none of it: the skill graph
// is the format's example, and a topic solved from the picture tells standalone in which chat its
// tasks may be given, while the module's every app shows the picture (docs/architecture.md, section
// 3.5).
type Topic struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	Description       string   `json:"description"`
	GradeLevels       []string `json:"grade_levels"`
	Slug              string   `json:"slug"`
	BuildsOn          []string `json:"builds_on"`
	SitePage          bool     `json:"site_page"`
	SolvedFromPicture bool     `json:"solved_from_picture,omitempty"`
}

// readTopics reads content/catalogs/topics.json of a clone of standalone, strictly, and gives its
// topics by ID.
func readTopics(clone string) (map[string]Topic, error) {
	file := filepath.Join(clone, "content", "catalogs", "topics.json")
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	var list []Topic
	if err := decoder.Decode(&list); err != nil {
		return nil, fmt.Errorf("%s: %w", file, err)
	}
	topics := map[string]Topic{}
	for _, topic := range list {
		if topic.ID == "" {
			return nil, fmt.Errorf("%s: a topic has no ID", file)
		}
		if _, twice := topics[topic.ID]; twice {
			return nil, fmt.Errorf("%s: the topic %s is there twice", file, topic.ID)
		}
		topics[topic.ID] = topic
	}
	return topics, nil
}

// ofTopic holds a task to its topic: the topic is one of the catalog, and a task of a topic solved
// from the picture has a picture, since its facts stand on the picture alone.
func ofTopic(ex Example, topics map[string]Topic) error {
	topic, known := topics[ex.Topic]
	if !known {
		return fmt.Errorf("%s: its topic %s is no topic of content/catalogs/topics.json", ex.ID, ex.Topic)
	}
	if topic.SolvedFromPicture && ex.Picture == nil {
		return fmt.Errorf("%s: its topic %s is solved from the picture, and the task has no picture", ex.ID, ex.Topic)
	}
	return nil
}

// Distractor is a wrong option of a task of standalone: the trap it comes from and what the child
// reads after choosing it.
type Distractor struct {
	Trap string `json:"trap"`
	Text string `json:"text"`
}

var (
	taskID     = regexp.MustCompile(`^[a-z][a-z0-9]*(-[a-z0-9]+)*$`)
	gradeLevel = regexp.MustCompile(`^[1-6]-[1-6]$`)
)

// readExamples reads every file content/examples/*.json of a clone of standalone, strictly, and
// gives its tasks by ID.
func readExamples(clone string) (map[string]Example, error) {
	files, err := filepath.Glob(filepath.Join(clone, "content", "examples", "*.json"))
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("%s holds no content/examples/*.json: is it a clone of mathtrail-standalone?", clone)
	}
	examples := map[string]Example{}
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}
		decoder := json.NewDecoder(bytes.NewReader(data))
		decoder.DisallowUnknownFields()
		var tasks []Example
		if err := decoder.Decode(&tasks); err != nil {
			return nil, fmt.Errorf("%s: %w", file, err)
		}
		for _, task := range tasks {
			if err := task.valid(); err != nil {
				return nil, fmt.Errorf("%s: %w", file, err)
			}
			if _, twice := examples[task.ID]; twice {
				return nil, fmt.Errorf("%s: the task %s is there twice", file, task.ID)
			}
			examples[task.ID] = task
		}
	}
	return examples, nil
}

// valid holds a task to the form standalone gives every reference task: five options, one of them
// right, and a distractor for each of the other four.
func (e Example) valid() error {
	var problems []string
	if !taskID.MatchString(e.ID) {
		problems = append(problems, fmt.Sprintf("the ID %q is no ID the format takes", e.ID))
	}
	if e.Topic == "" || e.Question == "" || e.Solution == "" {
		problems = append(problems, "it lacks its topic, question or solution")
	}
	if !gradeLevel.MatchString(e.GradeLevel) {
		problems = append(problems, fmt.Sprintf("its grade level %q is no band", e.GradeLevel))
	}
	if e.Difficulty < 1 || e.Difficulty > 5 {
		problems = append(problems, fmt.Sprintf("its difficulty %d is not 1 to 5", e.Difficulty))
	}
	if e.Idea < 0 {
		problems = append(problems, fmt.Sprintf("its idea %d is no place in a list", e.Idea))
	}
	if len(e.Options) != len(letters) {
		problems = append(problems, fmt.Sprintf("it has %d options, not %d", len(e.Options), len(letters)))
	}
	for _, letter := range letters {
		if strings.TrimSpace(e.Options[letter]) == "" {
			problems = append(problems, fmt.Sprintf("its option %s is empty or missing", letter))
		}
	}
	if !slices.Contains(letters, e.CorrectAnswer) {
		problems = append(problems, fmt.Sprintf("its correct answer %q is no option", e.CorrectAnswer))
	}
	for _, letter := range letters {
		distractor, has := e.Distractors[letter]
		switch {
		case letter == e.CorrectAnswer && has:
			problems = append(problems, fmt.Sprintf("its correct option %s has a distractor", letter))
		case letter != e.CorrectAnswer && (!has || distractor.Trap == "" || distractor.Text == ""):
			problems = append(problems, fmt.Sprintf("its wrong option %s has no trap and text", letter))
		}
	}
	if len(e.Distractors) > len(letters)-1 {
		problems = append(problems, "it has distractors of options it does not have")
	}
	if len(problems) > 0 {
		return fmt.Errorf("the task %s: %s", e.ID, strings.Join(problems, "; "))
	}
	return nil
}

// atCommit holds a clone of standalone to the commit a batch names: its HEAD is that commit, and
// git sees no change in its working tree, so that the tasks, their headers and the bench all come
// from it. git runs with none of the caller's settings, as the command named, as the go command
// of the bench is.
func atCommit(gitCommand, clone, commit string) error {
	git := func(args ...string) (string, error) {
		command := exec.Command(gitCommand, append([]string{"-C", clone}, args...)...)
		command.Env = append(os.Environ(), "GIT_CONFIG_GLOBAL=/dev/null", "GIT_CONFIG_NOSYSTEM=1")
		out, err := command.Output()
		return strings.TrimSpace(string(out)), err
	}
	head, err := git("rev-parse", "--verify", "HEAD")
	if err != nil {
		return fmt.Errorf("%s is no clone of standalone that git reads: %w", clone, err)
	}
	if head != commit {
		return fmt.Errorf("%s is at %s, and the batch is of %s: check that commit out", clone, head, commit)
	}
	changed, err := git("status", "--porcelain")
	if err != nil {
		return fmt.Errorf("%s: git status: %w", clone, err)
	}
	if changed != "" {
		return fmt.Errorf("%s has changes of its own, so its tasks are not those of %s", clone, commit)
	}
	return nil
}

// errNoExample is a task of a batch that standalone does not have.
var errNoExample = errors.New("standalone has no such task")
