// The tests of the recipe eval-check (R04, R25): the flat layout of each module's eval/, eval/<id>.json
// with the name the sample's ID, and the guard, which refuses every file whose name ends in .json under
// eval/ until spec's sample-check checks eval sets. The recipe calls no module-kit, so the modules here
// are bare folders, and the sample is spec's own example of an eval set's sheet, a synthetic one.

import assert from "node:assert/strict";
import { cpSync, mkdirSync, readdirSync, symlinkSync, writeFileSync } from "node:fs";
import { join } from "node:path";
import { test } from "node:test";

import { ROOT, mustRun, repository, run } from "./helpers.mjs";

const SAMPLE = join(ROOT, "spec", "contracts", "examples", "eval", "v1", "sheet-eval-set.json");
const SAMPLE_ID = "synthetic-1104";

const GUARD = /no sample enters an eval set until spec's sample-check checks eval sets \(R25\)$/;
const NOT_NAMED = /: not eval\/<id>\.json, whose <id> is 1 to 64 lowercase Latin letters, digits and hyphens \(R04\)$/;

// A repository with a bare module folder modules/<name>/ in it, and its eval/ when asked.
function withModule(t, name = "maths", { evalFolder = true } = {}) {
  const repo = repository(t);
  mkdirSync(join(repo, "modules", name), { recursive: true });
  if (evalFolder) mkdirSync(join(repo, "modules", name, "eval"));
  return repo;
}

function evalCheck(repo, env = {}) {
  return run(repo, "eval-check", [], env);
}

// The lines of stderr that name a path, each as it was written.
function refusals(result) {
  return result.stderr.split("\n").filter((line) => line.startsWith("modules/"));
}

function passes(result) {
  assert.equal(result.status, 0, result.stdout + result.stderr);
  assert.match(result.stdout, /^eval-check: every eval\/ is laid out as eval\/<id>\.json, and no module holds a sample$/m);
}

// A jq first on the PATH that cannot run, in a folder of the repository's own.
function brokenJq(repo) {
  const bin = join(repo, "bin");
  mkdirSync(bin);
  writeFileSync(join(bin, "jq"), "#!/bin/sh\nexit 127\n", { mode: 0o755 });
  return { PATH: `${bin}:${process.env.PATH}` };
}

test("without modules/, or with an empty one, there is nothing to check yet", (t) => {
  const repo = repository(t);
  let result = evalCheck(repo);
  assert.equal(result.status, 0, result.stderr);
  assert.match(result.stdout, /^modules\/: no modules yet$/m);
  mkdirSync(join(repo, "modules"));
  result = evalCheck(repo);
  assert.equal(result.status, 0, result.stderr);
  assert.match(result.stdout, /^modules\/: no modules yet$/m);
});

test("modules that is a file or a link is refused, as modules-check refuses it", (t) => {
  for (const make of [
    (repo) => writeFileSync(join(repo, "modules"), "Not a folder.\n"),
    (repo) => {
      mkdirSync(join(repo, "elsewhere", "maths", "eval"), { recursive: true });
      cpSync(SAMPLE, join(repo, "elsewhere", "maths", "eval", `${SAMPLE_ID}.json`));
      symlinkSync("elsewhere", join(repo, "modules"));
    },
  ]) {
    const repo = repository(t);
    make(repo);
    const result = evalCheck(repo);
    assert.equal(result.status, 1);
    assert.match(result.stderr, /^modules: not a folder; it holds the folders of the modules \(R03\)$/m);
  }
});

test("a module without eval/ passes, and so does an empty eval/", (t) => {
  passes(evalCheck(withModule(t, "maths", { evalFolder: false })));
  passes(evalCheck(withModule(t)));
});

test("files and folders named with a dot are no part of the layout", (t) => {
  const repo = withModule(t);
  const folder = join(repo, "modules", "maths", "eval");
  writeFileSync(join(folder, ".gitkeep"), "");
  mkdirSync(join(folder, ".draft"));
  writeFileSync(join(folder, ".draft", "notes.txt"), "Not a sample.\n");
  passes(evalCheck(repo));
});

test("the guard refuses a .json named with a dot, or in a folder named with one: it is published all the same", (t) => {
  const repo = withModule(t);
  const folder = join(repo, "modules", "maths", "eval");
  mkdirSync(join(folder, ".draft"));
  cpSync(SAMPLE, join(folder, ".draft", `${SAMPLE_ID}.json`));
  cpSync(SAMPLE, join(folder, `.${SAMPLE_ID}.json`));
  const result = evalCheck(repo);
  assert.equal(result.status, 1);
  const lines = refusals(result).sort();
  assert.equal(lines.length, 2, result.stderr);
  assert.ok(lines[0].startsWith(`modules/maths/eval/.draft/${SAMPLE_ID}.json: `), lines[0]);
  assert.ok(lines[1].startsWith(`modules/maths/eval/.${SAMPLE_ID}.json: `), lines[1]);
  for (const line of lines) assert.match(line, GUARD);
});

test("what is not a module folder in modules/ is modules-check's to refuse, not eval-check's", (t) => {
  const repo = withModule(t);
  writeFileSync(join(repo, "modules", "README.md"), "Not a module.\n");
  mkdirSync(join(repo, "modules", ".draft", "eval"), { recursive: true });
  cpSync(SAMPLE, join(repo, "modules", ".draft", "eval", `${SAMPLE_ID}.json`));
  symlinkSync(".draft", join(repo, "modules", "alias"));
  passes(evalCheck(repo));
});

test("a folder, another file, a link and a special file in eval/ are refused", (t) => {
  const repo = withModule(t);
  const folder = join(repo, "modules", "maths", "eval");
  mkdirSync(join(folder, "tasks"));
  mkdirSync(join(folder, "folder.json"));
  writeFileSync(join(folder, "README.md"), "Not a sample.\n");
  symlinkSync("README.md", join(folder, "alias"));
  mustRun("mkfifo", [join(folder, "pipe")]);
  const result = evalCheck(repo);
  assert.equal(result.status, 1);
  const other = "a folder, a link or a special file; eval/ holds the files eval/<id>.json alone (R04)";
  assert.deepEqual(refusals(result).sort(), [
    "modules/maths/eval/README.md: not eval/<id>.json, whose <id> is 1 to 64 lowercase Latin letters, digits and hyphens (R04)",
    `modules/maths/eval/alias: ${other}`,
    `modules/maths/eval/folder.json: ${other}`,
    `modules/maths/eval/pipe: ${other}`,
    `modules/maths/eval/tasks: ${other}`,
  ]);
});

test("a sample in a folder of eval/ breaks the layout, and the guard refuses it too", (t) => {
  const repo = withModule(t);
  mkdirSync(join(repo, "modules", "maths", "eval", "tasks"));
  cpSync(SAMPLE, join(repo, "modules", "maths", "eval", "tasks", `${SAMPLE_ID}.json`));
  const result = evalCheck(repo);
  assert.equal(result.status, 1);
  const lines = refusals(result);
  assert.equal(lines.length, 2, result.stderr);
  assert.ok(lines[0].startsWith("modules/maths/eval/tasks: a folder, a link or a special file"), lines[0]);
  assert.ok(lines[1].startsWith(`modules/maths/eval/tasks/${SAMPLE_ID}.json: `), lines[1]);
  assert.match(lines[1], GUARD);
});

test("eval that is a file or a link is refused", (t) => {
  for (const make of [
    (module) => writeFileSync(join(module, "eval"), "Not a folder.\n"),
    (module) => {
      mkdirSync(join(module, "samples"));
      symlinkSync("samples", join(module, "eval"));
    },
  ]) {
    const repo = withModule(t, "maths", { evalFolder: false });
    make(join(repo, "modules", "maths"));
    const result = evalCheck(repo);
    assert.equal(result.status, 1);
    assert.deepEqual(refusals(result), [
      "modules/maths/eval: not a folder; eval/ is the folder of the module's eval set (R04)",
    ]);
  }
});

test("a sample laid out right is refused by the guard alone", (t) => {
  const repo = withModule(t);
  cpSync(SAMPLE, join(repo, "modules", "maths", "eval", `${SAMPLE_ID}.json`));
  const result = evalCheck(repo);
  assert.equal(result.status, 1);
  const lines = refusals(result);
  assert.equal(lines.length, 1, result.stderr);
  assert.ok(lines[0].startsWith(`modules/maths/eval/${SAMPLE_ID}.json: `), lines[0]);
  assert.match(lines[0], GUARD);
});

test("a sample whose name is not its ID breaks the layout, and the guard refuses it too", (t) => {
  const repo = withModule(t);
  cpSync(SAMPLE, join(repo, "modules", "maths", "eval", "other-id.json"));
  const result = evalCheck(repo);
  assert.equal(result.status, 1);
  const lines = refusals(result);
  assert.equal(lines.length, 2, result.stderr);
  assert.equal(
    lines[0],
    `modules/maths/eval/other-id.json: its sample's ID is "${SAMPLE_ID}", and a sample lies in eval/<its ID>.json (R04)`,
  );
  assert.match(lines[1], GUARD);
});

test("an ID that is the name but for a newline or a NUL at its end is another ID", (t) => {
  for (const [id, shown] of [
    ["x\n", String.raw`"x\n"`],
    ["x\u0000", String.raw`"x\u0000"`],
  ]) {
    const repo = withModule(t);
    writeFileSync(join(repo, "modules", "maths", "eval", "x.json"), JSON.stringify({ id }));
    const result = evalCheck(repo);
    assert.equal(result.status, 1);
    const lines = refusals(result);
    assert.equal(lines.length, 2, result.stderr);
    assert.equal(lines[0], `modules/maths/eval/x.json: its sample's ID is ${shown}, and a sample lies in eval/<its ID>.json (R04)`);
    assert.match(lines[1], GUARD);
  }
});

test("a name that is no ID breaks the layout, and the guard refuses it too", (t) => {
  for (const name of [`${SAMPLE_ID.toUpperCase()}.json`, "a_b.json", "é.json", `${"a".repeat(65)}.json`, `${SAMPLE_ID}.JSON`]) {
    const repo = withModule(t);
    cpSync(SAMPLE, join(repo, "modules", "maths", "eval", name));
    // The locale the caller runs in does not widen [a-z]: the recipe sets its own.
    const result = evalCheck(repo, { LC_ALL: "C.UTF-8" });
    assert.equal(result.status, 1);
    const lines = refusals(result);
    assert.match(lines[0], NOT_NAMED);
    // A name that does not end in .json, letter case and all, is no sample (contracts/README.md in
    // spec, section 13), so the guard does not name it.
    assert.equal(lines.length, name.endsWith(".json") ? 2 : 1, result.stderr);
    if (name.endsWith(".json")) assert.match(lines[1], GUARD);
  }
});

test("a file with no one JSON object holding a string id breaks the layout, and the guard refuses it too", (t) => {
  for (const content of ["{", "", `{"id": "x"}\n{"id": "x"}\n`, `[{"id": "x"}]`, `{"id": 1}`, `{"identifier": "x"}`]) {
    const repo = withModule(t);
    writeFileSync(join(repo, "modules", "maths", "eval", "x.json"), content);
    const result = evalCheck(repo);
    assert.equal(result.status, 1, content);
    const lines = refusals(result);
    assert.equal(lines.length, 2, content + result.stderr);
    assert.equal(
      lines[0],
      "modules/maths/eval/x.json: holds no one JSON object with a string id, so its name cannot be held to its sample's ID (R04)",
    );
    assert.match(lines[1], GUARD);
  }
});

test("every module is checked: a breach in one does not stop the rest", (t) => {
  const repo = withModule(t, "alpha");
  mkdirSync(join(repo, "modules", "beta", "eval"), { recursive: true });
  mkdirSync(join(repo, "modules", "gamma", "eval"), { recursive: true });
  writeFileSync(join(repo, "modules", "alpha", "eval", "notes.txt"), "Not a sample.\n");
  cpSync(SAMPLE, join(repo, "modules", "gamma", "eval", `${SAMPLE_ID}.json`));
  const result = evalCheck(repo);
  assert.equal(result.status, 1);
  const lines = refusals(result);
  assert.equal(lines.length, 2, result.stderr);
  assert.ok(lines[0].startsWith("modules/alpha/eval/notes.txt: not eval/<id>.json"), lines[0]);
  assert.ok(lines[1].startsWith(`modules/gamma/eval/${SAMPLE_ID}.json: `), lines[1]);
  assert.match(lines[1], GUARD);
});

test("jq is needed only when a file is to be read: without it the recipe says so with the status 2", (t) => {
  const quiet = withModule(t);
  writeFileSync(join(quiet, "modules", "maths", "eval", ".gitkeep"), "");
  passes(evalCheck(quiet, brokenJq(quiet)));
  const repo = withModule(t);
  cpSync(SAMPLE, join(repo, "modules", "maths", "eval", `${SAMPLE_ID}.json`));
  const result = evalCheck(repo, brokenJq(repo));
  assert.equal(result.status, 2);
  assert.match(result.stderr, /^eval-check: jq cannot be run; the toolchain stage has it$/m);
});

test("a run leaves the tree as it was", (t) => {
  const repo = withModule(t);
  cpSync(SAMPLE, join(repo, "modules", "maths", "eval", "other-id.json"));
  const before = readdirSync(repo, { recursive: true }).sort();
  evalCheck(repo);
  assert.deepEqual(readdirSync(repo, { recursive: true }).sort(), before);
});
