// The test of the port's output against the format (R12; R20): tools/port, run by go, writes tasks
// of a made-up standalone — with a picture of a row, and with answers that are a time and a word,
// a choice written by hand — into a module made of the format's examples of olympiad-math, and
// module-kit check accepts it, with no warning about the picture. A survey of the same standalone,
// by the recipe _port-survey, writes all its tasks and counts the breaches of task.no-leaks. The tests of
// tools/port itself are Go's, in "Go / test", which reads nothing of spec; this one needs the
// format's examples and module-kit, so it runs in the job Modules.

import assert from "node:assert/strict";
import { cpSync, mkdirSync, readFileSync, writeFileSync } from "node:fs";
import { join } from "node:path";
import { test } from "node:test";

import { ROOT, mustRun, run, temporaryFolder } from "./helpers.mjs";

const PARTS = join(ROOT, "spec", "format", "examples");

// A module of the format's examples of olympiad-math: manifest, skill graph, catalogs and solving
// steps, and the repository's LICENSE.
function exampleModule(t) {
  const module = join(temporaryFolder(t), "olympiad-math");
  mkdirSync(module);
  for (const [part, file] of [
    ["manifest", "manifest.yaml"],
    ["skill-graph", "skill-graph.yaml"],
    ["traps", "traps.yaml"],
    ["techniques", "techniques.yaml"],
    ["solving-steps", "solving-steps.yaml"],
  ]) {
    cpSync(join(PARTS, part, "valid", "olympiad-math.yaml"), join(module, file));
  }
  cpSync(join(ROOT, "LICENSE"), join(module, "LICENSE"));
  return module;
}

// A made-up standalone with tasks of counting.gaps in its form, texts of their own: one with a
// picture of a row, one whose wrong option is a word, one whose answer is a time and one whose
// answer is a word of its question, and one whose explanation names its key, which no batch names.
// It is a git repository with one commit, as a clone is, and the port holds it to the batch's
// commit.
function madeUpStandalone(t) {
  const clone = temporaryFolder(t);
  const task = (id, extra) => ({
    id,
    topic: "counting.gaps",
    grade_level: "3-4",
    difficulty: 1,
    question: "Flags hang on a rope, 3 metres apart, 11 flags from the first to the last. How long is the rope between the first flag and the last?",
    options: { A: "33", B: "30", C: "11", D: "36", E: "27" },
    correct_answer: "B",
    solution: "11 flags have 10 gaps of 3 metres each: 10 × 3 = 30 metres.",
    distractors: {
      A: { trap: "off_by_one", text: "11 flags have 10 gaps, not 11." },
      C: { trap: "number_from_text", text: "That is the number of flags." },
      D: { trap: "off_by_one", text: "Count the gaps between the flags once more." },
      E: { trap: "off_by_one", text: "11 flags have 10 gaps, not 9." },
    },
    ...extra,
  });
  mkdirSync(join(clone, "content", "examples"), { recursive: true });
  writeFileSync(
    join(clone, "content", "examples", "counting.gaps.json"),
    JSON.stringify([
      task("flags-1", { picture: { kind: "row", items: [{}, {}, {}, { skip: true }, {}], gaps: "3" } }),
      task("flags-2", { options: { A: "33", B: "30", C: "none", D: "36", E: "27" } }),
      task("flags-3", {
        question: "The first flag goes up at 8:00, and a flag goes up every 10 minutes. When does the fourth go up?",
        options: { A: "8:40", B: "8:30", C: "8:00", D: "8:10", E: "9:00" },
        solution: "The fourth flag goes up 3 gaps of 10 minutes after the first: 30 minutes after 8:00.",
        distractors: {
          A: { trap: "off_by_one", text: "Four flags have 3 gaps, not 4." },
          C: { trap: "number_from_text", text: "That is when the first goes up." },
          D: { trap: "stopped_early", text: "That is the second flag." },
          E: { trap: "wrong_operation", text: "Every gap is 10 minutes, not 20." },
        },
      }),
      task("flags-4", {
        question: "Red, blue and green flags hang in this order, again and again. Of what colour is the fifth flag?",
        options: { A: "red", B: "blue", C: "green", D: "white", E: "black" },
        solution: "The colours go red, blue, green, red, and then the second colour again.",
        distractors: {
          A: { trap: "off_by_one", text: "That is the fourth flag." },
          C: { trap: "off_by_one", text: "That is the sixth flag." },
          D: { trap: "number_from_text", text: "No flag is of that colour." },
          E: { trap: "number_from_text", text: "No flag is of that colour." },
        },
        correct_answer: "B",
      }),
      task("flags-5", {
        distractors: {
          A: { trap: "off_by_one", text: "The rope is 30 metres, not 33." },
          C: { trap: "number_from_text", text: "That is the number of flags." },
          D: { trap: "off_by_one", text: "Count the gaps between the flags once more." },
          E: { trap: "off_by_one", text: "11 flags have 10 gaps, not 9." },
        },
      }),
    ]),
  );
  mkdirSync(join(clone, "content", "catalogs"));
  writeFileSync(
    join(clone, "content", "catalogs", "topics.json"),
    JSON.stringify([
      {
        id: "counting.gaps",
        name: "Gaps",
        description: "Count the gaps between things in a row.",
        grade_levels: ["1-2", "3-4", "5-6"],
        slug: "gaps",
        builds_on: [],
        site_page: true,
      },
    ]),
  );
  const env = { ...process.env, GIT_CONFIG_GLOBAL: "/dev/null", GIT_CONFIG_NOSYSTEM: "1" };
  const git = (...args) =>
    mustRun("git", ["-C", clone, "-c", "user.name=Test", "-c", "user.email=test@example.com", ...args], { env }).stdout.trim();
  git("init", "-q");
  git("add", "-A");
  git("commit", "-q", "-m", "Tasks");
  return { clone, commit: git("rev-parse", "HEAD") };
}

test("the port's tasks and its picture of a row pass module-kit check, with no warning", (t) => {
  const { clone, commit } = madeUpStandalone(t);
  const module = exampleModule(t);
  const work = temporaryFolder(t);
  const batch = join(work, "batch.yaml");
  writeFileSync(
    batch,
    [
      `standalone: ${commit}`,
      "canvas: {kind: handwriting_sheet, grid: squares}",
      "skills:",
      "  counting-gaps: {techniques: [draw, small-numbers]}",
      "tasks:",
      "  flags-1: {units: [metre, metres, m]}",
      "  flags-2: {units: [metre, metres, m]}",
      "  flags-3: {}",
      "  flags-4: {}",
      "",
    ].join("\n"),
  );
  const overlays = join(work, "overlay");
  mkdirSync(overlays);
  writeFileSync(
    join(overlays, "flags-1.yaml"),
    [
      "givens: [The flags are 3 metres apart., There are 11 flags.]",
      "unknown: How long the rope is between the first flag and the last.",
      "simpler_case: Try 3 flags first.",
      "picture_alt: Flags in a row on a rope, with a stretch of the rope left out, and every gap drawn marked 3.",
      "",
    ].join("\n"),
  );
  const args = ["run", "./tools/port", "--standalone", clone, "--batch", batch, "--overlay", overlays, "--module", module];
  const result = mustRun("go", args, { cwd: ROOT });
  assert.match(result.stdout, /port: 4 tasks written/);

  const checked = mustRun("module-kit", ["check", module]);
  assert.match(checked.stdout, /the module breaks no rule of the format 0\.4/);
  assert.doesNotMatch(checked.stdout, /warning/);

  const task = readFileSync(join(module, "tasks", "flags-2.yaml"), "utf8");
  assert.match(task, /^traps: \["number-from-text"\]$/m, "a word among the wrong options becomes a trap of the task");
  for (const [id, key, notation] of [["flags-3", "8:30", ""], ["flags-4", "blue", ' {10}notation: "word"\n']]) {
    const byHand = readFileSync(join(module, "tasks", `${id}.yaml`), "utf8");
    assert.match(byHand, new RegExp(`binding: "anywhere"\n {8}check:\n {10}kind: "choice"\n${notation} {10}key: "${key}"\n`),
      `${id} is a choice written by hand${notation ? ", with notation word" : ""}`);
    assert.match(byHand, new RegExp(`key: "${key}"\n {10}options: \\[`), `${id} keeps its options in its check`);
  }

  const again = mustRun("go", [...args, "--check"], { cwd: ROOT });
  assert.match(again.stdout, /are what the port writes/);
});

test("a survey writes every task no batch names, and counts the breaches of task.no-leaks", (t) => {
  const { clone, commit } = madeUpStandalone(t);
  const module = join(temporaryFolder(t), "survey");
  const result = run(ROOT, "_port-survey", [clone, commit, module]);
  assert.equal(result.status, 0, result.stderr);
  assert.match(result.stdout, /: 5 — number 3, fraction 0, choice by hand 2, choice tapped 0\n/);
  assert.match(result.stdout, /port: refused: 0\n/);
  assert.match(result.stdout, /: task\.no-leaks 1, for the overlays to mend\n/);
  assert.match(readFileSync(join(module, "tasks", "flags-1.yaml"), "utf8"), /alt: "Survey: a picture of standalone, with no alt written\."/);

  const again = run(ROOT, "_port-survey", [clone, commit, module]);
  assert.equal(again.status, 2);
  assert.match(again.stderr, /is there already/);

  const elsewhere = run(ROOT, "_port-survey", [clone, "0".repeat(40), join(temporaryFolder(t), "survey")]);
  assert.equal(elsewhere.status, 2, "a port that could not run is no refusal of a task");
  assert.match(elsewhere.stderr, /and the batch is of 0{40}/);
  assert.doesNotMatch(elsewhere.stdout, /module-kit check/);
});
