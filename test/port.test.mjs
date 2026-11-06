// The test of the port's output against the format (R12): tools/port, run by go, writes tasks of a
// made-up standalone, with a picture of a row, into a module made of the format's examples of
// olympiad-math, and module-kit check accepts it, with no warning about the picture. The tests of
// tools/port itself are Go's, in "Go / test", which reads nothing of spec; this one needs the
// format's examples and module-kit, so it runs in the job Modules.

import assert from "node:assert/strict";
import { cpSync, mkdirSync, readFileSync, writeFileSync } from "node:fs";
import { join } from "node:path";
import { test } from "node:test";

import { ROOT, mustRun, temporaryFolder } from "./helpers.mjs";

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

// A made-up standalone with two tasks of counting.gaps in its form, texts of their own: one with a
// picture of a row, one whose wrong option is a word. It is a git repository with one commit, as a
// clone is, and the port holds it to the batch's commit.
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
  assert.match(result.stdout, /port: 2 tasks written/);

  const checked = mustRun("module-kit", ["check", module]);
  assert.match(checked.stdout, /the module breaks no rule of the format 0\.3/);
  assert.doesNotMatch(checked.stdout, /warning/);

  const task = readFileSync(join(module, "tasks", "flags-2.yaml"), "utf8");
  assert.match(task, /^traps: \["number-from-text"\]$/m, "a word among the wrong options becomes a trap of the task");

  const again = mustRun("go", [...args, "--check"], { cwd: ROOT });
  assert.match(again.stdout, /are what the port writes/);
});
