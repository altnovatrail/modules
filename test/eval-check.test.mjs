// The tests of the recipe eval-check (R04, R25): the flat layout of each module's eval/, eval/<id>.json
// with the name the sample's ID and nothing named with a dot, and each module's eval set held by spec's
// sample-check --eval-set — the real one, from the spec/ of this repository at its tag, and a stand-in
// where a test holds the form of the call or a status sample-check gives. The recipe calls no
// module-kit, so the modules here are bare folders, and the samples are spec's own examples, synthetic
// ones, made samples of a module's task where a test needs one.

import assert from "node:assert/strict";
import { cpSync, mkdirSync, readFileSync, readdirSync, symlinkSync, writeFileSync } from "node:fs";
import { join } from "node:path";
import { test } from "node:test";

import { ROOT, mustRun, repository, run, stubSpec, withSpec } from "./helpers.mjs";

const EXAMPLES = join(ROOT, "spec", "contracts", "examples", "eval", "v1");
const SAMPLE = join(EXAMPLES, "sheet-eval-set.json");
const SAMPLE_ID = "synthetic-1104";

const NOT_NAMED = /: not eval\/<id>\.json, whose <id> is 1 to 64 lowercase Latin letters, digits and hyphens \(R04\)$/;
const DOT = "named with a dot; eval/ holds the files eval/<id>.json alone, and nothing named with a dot (R04)";
const LAID_OUT = /^eval-check: every eval\/ is laid out as eval\/<id>\.json, and no module holds a sample$/m;
const HELD = /^eval-check: every eval\/ is laid out as eval\/<id>\.json, and sample-check holds the eval set of every module that has one$/m;

// An example of spec's sample file, by its case.
function example(name) {
  return JSON.parse(readFileSync(join(EXAMPLES, `${name}.json`), "utf8"));
}

// A sample made a sample of a task of the module given: its source modules in place of tutor, as
// spec's own tests of the eval-set mode make one.
function ofModule(sample, module = "novatrail/maths") {
  const { tutor, ...task } = sample.task;
  return { ...sample, task: { modules: { tag: "v0.1.1", module, version: "0.1.1", package_digest: `sha256:${"0".repeat(64)}` }, ...task } };
}

// A repository with spec/ and a bare module folder modules/<name>/ in it, and its eval/ when asked.
function withModule(t, name = "maths", { evalFolder = true } = {}) {
  const repo = withSpec(repository(t));
  mkdirSync(join(repo, "modules", name), { recursive: true });
  if (evalFolder) mkdirSync(join(repo, "modules", name, "eval"));
  return repo;
}

// A sample written as a file of a module's eval/.
function put(repo, module, file, sample) {
  const path = join(repo, "modules", module, "eval", file);
  mkdirSync(join(path, ".."), { recursive: true });
  writeFileSync(path, typeof sample === "string" ? sample : `${JSON.stringify(sample, null, 2)}\n`);
}

function evalCheck(repo, env = {}) {
  return run(repo, "eval-check", [], env);
}

// The lines of stderr that name a path: the breaches of the layout, each as it was written.
function refusals(result) {
  return result.stderr.split("\n").filter((line) => line.startsWith("modules/"));
}

// What sample-check printed: the lines that name a file, file:line: message, and its verdicts.
function checked(result) {
  const lines = result.stdout.split("\n");
  return { breaches: lines.filter((line) => line.startsWith("modules/")), verdicts: lines.filter((line) => line.startsWith("sample-check: ")) };
}

function laidOut(result) {
  assert.equal(result.status, 0, result.stdout + result.stderr);
  assert.match(result.stdout, LAID_OUT);
  assert.deepEqual(checked(result).verdicts, [], "sample-check is not called when no module holds a sample");
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

test("a module without eval/ passes, and so does an empty eval/, with no call of sample-check", (t) => {
  laidOut(evalCheck(withModule(t, "maths", { evalFolder: false })));
  laidOut(evalCheck(withModule(t)));
});

test("a sample of the module's own task, laid out right, is held by sample-check and passes", (t) => {
  const repo = withModule(t);
  put(repo, "maths", `${SAMPLE_ID}.json`, ofModule(example("sheet-eval-set")));
  const result = evalCheck(repo);
  assert.equal(result.status, 0, result.stdout + result.stderr);
  assert.deepEqual(checked(result).verdicts, ["sample-check: ok — 1 sample of the eval set of novatrail/maths"]);
  assert.match(result.stdout, HELD);
});

test("sample-check refuses what is no sample of the module's eval set, and so does the recipe", (t) => {
  const repo = withModule(t);
  // A task of tutor's, a task of another module, a sample of measurement, and two samples of one ID.
  put(repo, "maths", `${SAMPLE_ID}.json`, example("sheet-eval-set"));
  put(repo, "maths", "synthetic-9002.json", { ...ofModule(example("sheet-eval-set"), "novatrail/other"), id: "synthetic-9002" });
  put(repo, "maths", "synthetic-1103.json", ofModule(example("sheet-measurement")));
  put(repo, "maths", "twice.json", { ...ofModule(example("sheet-eval-set")), id: "twice" });
  put(repo, "maths", "twice-2.json", { ...ofModule(example("sheet-eval-set")), id: "twice" });
  const result = evalCheck(repo);
  assert.equal(result.status, 1);
  const { breaches, verdicts } = checked(result);
  assert.deepEqual(verdicts, ["sample-check: 4 errors"], result.stdout);
  const of = (file) => breaches.filter((line) => line.startsWith(`modules/maths/eval/${file}:`)).join("\n");
  assert.match(of(`${SAMPLE_ID}.json`), /is of a task of that module: its source is modules/);
  assert.match(of("synthetic-9002.json"), /is of a task of that module, not of "novatrail\/other"/);
  assert.match(of("synthetic-1103.json"), /is of the purpose PURPOSE_EVAL_SET, not "PURPOSE_MEASUREMENT"/);
  assert.match(breaches.join("\n"), /no two samples of an eval set share an ID/);
  // twice-2.json is named not by its ID: the layout says so too.
  assert.deepEqual(refusals(result), ['modules/maths/eval/twice-2.json: its sample\'s ID is "twice", and a sample lies in eval/<its ID>.json (R04)']);
});

test("nothing in eval/ is named with a dot, at any depth: the repository publishes it all the same", (t) => {
  const repo = withModule(t);
  const folder = join(repo, "modules", "maths", "eval");
  writeFileSync(join(folder, ".gitkeep"), "");
  put(repo, "maths", `.drafts/${SAMPLE_ID}.json`, ofModule(example("sheet-eval-set")));
  put(repo, "maths", `.${SAMPLE_ID}.json`, ofModule(example("sheet-eval-set")));
  const result = evalCheck(repo);
  assert.equal(result.status, 1);
  assert.deepEqual(refusals(result).sort(), [
    `modules/maths/eval/.drafts: ${DOT}`,
    `modules/maths/eval/.gitkeep: ${DOT}`,
    `modules/maths/eval/.${SAMPLE_ID}.json: ${DOT}`,
  ]);
  // What is named with a dot is no part of the eval set (contracts/README.md in spec, section 13),
  // so sample-check does not get it.
  assert.deepEqual(checked(result).verdicts, []);
});

test("a name with a dot below a folder of eval/ is refused too, and kept from sample-check", (t) => {
  const repo = repository(t);
  const calls = stubSpec(repo);
  put(repo, "maths", "tasks/.x.json", { id: "x" });
  put(repo, "maths", "tasks/y.json", { id: "y" });
  const result = evalCheck(repo);
  assert.equal(result.status, 1);
  assert.deepEqual(refusals(result).sort(), [
    `modules/maths/eval/tasks/.x.json: ${DOT}`,
    "modules/maths/eval/tasks: a folder, a link or a special file; eval/ holds the files eval/<id>.json alone (R04)",
  ]);
  assert.deepEqual(calls(), [["--eval-set", "novatrail/maths", "modules/maths/eval/tasks/y.json"]]);
});

test("what is not a module folder in modules/ is modules-check's to refuse, not eval-check's", (t) => {
  const repo = withModule(t);
  writeFileSync(join(repo, "modules", "README.md"), "Not a module.\n");
  mkdirSync(join(repo, "modules", ".draft", "eval"), { recursive: true });
  cpSync(SAMPLE, join(repo, "modules", ".draft", "eval", `${SAMPLE_ID}.json`));
  symlinkSync(".draft", join(repo, "modules", "alias"));
  laidOut(evalCheck(repo));
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
  assert.deepEqual(checked(result).verdicts, [], "no regular .json file, so no call of sample-check");
});

test("a link named <id>.json is refused by the layout and never given to sample-check", (t) => {
  const repo = withModule(t);
  put(repo, "maths", "../real.json", ofModule(example("sheet-eval-set")));
  symlinkSync("../real.json", join(repo, "modules", "maths", "eval", `${SAMPLE_ID}.json`));
  const result = evalCheck(repo);
  assert.equal(result.status, 1);
  assert.deepEqual(refusals(result), [`modules/maths/eval/${SAMPLE_ID}.json: a folder, a link or a special file; eval/ holds the files eval/<id>.json alone (R04)`]);
  assert.deepEqual(checked(result).verdicts, []);
});

test("a sample in a folder of eval/ breaks the layout, and sample-check holds it all the same", (t) => {
  const repo = withModule(t);
  put(repo, "maths", `tasks/${SAMPLE_ID}.json`, ofModule(example("sheet-eval-set")));
  const result = evalCheck(repo);
  assert.equal(result.status, 1);
  assert.deepEqual(refusals(result), ["modules/maths/eval/tasks: a folder, a link or a special file; eval/ holds the files eval/<id>.json alone (R04)"]);
  assert.deepEqual(checked(result).verdicts, ["sample-check: ok — 1 sample of the eval set of novatrail/maths"]);
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

test("a sample whose name is not its ID breaks the layout, though sample-check holds it", (t) => {
  const repo = withModule(t);
  put(repo, "maths", "other-id.json", ofModule(example("sheet-eval-set")));
  const result = evalCheck(repo);
  assert.equal(result.status, 1);
  assert.deepEqual(refusals(result), [
    `modules/maths/eval/other-id.json: its sample's ID is "${SAMPLE_ID}", and a sample lies in eval/<its ID>.json (R04)`,
  ]);
  assert.deepEqual(checked(result).verdicts, ["sample-check: ok — 1 sample of the eval set of novatrail/maths"]);
});

test("an ID that is the name but for a newline or a NUL at its end is another ID", (t) => {
  for (const [id, shown] of [
    ["x\n", String.raw`"x\n"`],
    ["x\u0000", String.raw`"x\u0000"`],
  ]) {
    const repo = withModule(t);
    put(repo, "maths", "x.json", JSON.stringify({ id }));
    const result = evalCheck(repo);
    assert.equal(result.status, 1);
    assert.deepEqual(refusals(result), [`modules/maths/eval/x.json: its sample's ID is ${shown}, and a sample lies in eval/<its ID>.json (R04)`]);
  }
});

test("a name that is no ID breaks the layout", (t) => {
  for (const name of [`${SAMPLE_ID.toUpperCase()}.json`, "a_b.json", "é.json", `${"a".repeat(65)}.json`, `${SAMPLE_ID}.JSON`]) {
    const repo = withModule(t);
    put(repo, "maths", name, ofModule(example("sheet-eval-set")));
    // The locale the caller runs in does not widen [a-z]: the recipe sets its own.
    const result = evalCheck(repo, { LC_ALL: "C.UTF-8" });
    assert.equal(result.status, 1);
    const lines = refusals(result);
    assert.equal(lines.length, 1, result.stderr);
    assert.match(lines[0], NOT_NAMED);
    // A name that does not end in .json, letter case and all, is no sample (contracts/README.md in
    // spec, section 13), so sample-check does not get it; one that does, it holds.
    assert.deepEqual(checked(result).verdicts, name.endsWith(".json") ? ["sample-check: ok — 1 sample of the eval set of novatrail/maths"] : []);
  }
});

test("a file with no one JSON object holding a string id breaks the layout, and sample-check refuses it too", (t) => {
  for (const content of ["{", "", `{"id": "x"}\n{"id": "x"}\n`, `[{"id": "x"}]`, `{"id": 1}`, `{"identifier": "x"}`]) {
    const repo = withModule(t);
    put(repo, "maths", "x.json", content);
    const result = evalCheck(repo);
    assert.equal(result.status, 1, content);
    assert.deepEqual(refusals(result), [
      "modules/maths/eval/x.json: holds no one JSON object with a string id, so its name cannot be held to its sample's ID (R04)",
    ]);
    assert.notEqual(checked(result).breaches.length, 0, content + result.stdout);
  }
});

test("every module is checked: a breach in one does not stop the rest", (t) => {
  const repo = withModule(t, "alpha");
  mkdirSync(join(repo, "modules", "beta", "eval"), { recursive: true });
  writeFileSync(join(repo, "modules", "alpha", "eval", "notes.txt"), "Not a sample.\n");
  put(repo, "gamma", `${SAMPLE_ID}.json`, ofModule(example("sheet-eval-set"), "novatrail/gamma"));
  const result = evalCheck(repo);
  assert.equal(result.status, 1);
  const lines = refusals(result);
  assert.equal(lines.length, 1, result.stderr);
  assert.ok(lines[0].startsWith("modules/alpha/eval/notes.txt: not eval/<id>.json"), lines[0]);
  assert.deepEqual(checked(result).verdicts, ["sample-check: ok — 1 sample of the eval set of novatrail/gamma"]);
});

test("sample-check is called once a module, with its ID and the files of its eval set, sorted", (t) => {
  const repo = repository(t);
  const calls = stubSpec(repo);
  for (const module of ["alpha", "beta", "gamma"]) mkdirSync(join(repo, "modules", module, "eval"), { recursive: true });
  put(repo, "alpha", "b.json", { id: "b" });
  put(repo, "alpha", "a.json", { id: "a" });
  put(repo, "gamma", "c.json", { id: "c" });
  const result = evalCheck(repo);
  assert.equal(result.status, 0, result.stdout + result.stderr);
  assert.deepEqual(calls(), [
    ["--eval-set", "novatrail/alpha", "modules/alpha/eval/a.json", "modules/alpha/eval/b.json"],
    ["--eval-set", "novatrail/gamma", "modules/gamma/eval/c.json"],
  ]);
  assert.match(result.stdout, HELD);
});

test("the status of sample-check is the recipe's: 1 when it refuses, 2 when it cannot check", (t) => {
  for (const [given, expected] of [
    [1, 1],
    [2, 2],
    [7, 2],
  ]) {
    const repo = repository(t);
    stubSpec(repo, given);
    put(repo, "maths", "a.json", { id: "a" });
    const result = evalCheck(repo);
    assert.equal(result.status, expected, `sample-check exits ${given}: ${result.stdout}${result.stderr}`);
    assert.doesNotMatch(result.stdout, /^eval-check: every eval/m);
  }
});

test("the status is the highest of what failed: a breach of the layout and a check that could not run give 2", (t) => {
  const repo = repository(t);
  const calls = stubSpec(repo, 2);
  put(repo, "alpha", "a.json", { id: "a" });
  put(repo, "alpha", "notes.txt", "Not a sample.\n");
  put(repo, "beta", "b.json", { id: "b" });
  const result = evalCheck(repo);
  assert.equal(result.status, 2, result.stdout + result.stderr);
  assert.equal(refusals(result).length, 1, result.stderr);
  assert.equal(calls().length, 2, "a breach in alpha does not stop the check of beta");
});

test("spec's tools are installed before the first eval set, and if they cannot be, the status is 2, not a refusal", (t) => {
  const repo = repository(t);
  const calls = stubSpec(repo, 0, 1);
  put(repo, "alpha", "a.json", { id: "a" });
  put(repo, "beta", "b.json", { id: "b" });
  const result = evalCheck(repo);
  assert.equal(result.status, 2, result.stdout + result.stderr);
  assert.match(result.stderr, /^eval-check: the tools of spec's sample-check cannot be installed, so no eval set is checked$/m);
  assert.equal((result.stderr.match(/cannot be installed/g) ?? []).length, 1, "the install is tried once");
  assert.deepEqual(calls(), []);
});

test("jq is needed only when a file is to be read: without it the recipe says so with the status 2", (t) => {
  const quiet = withModule(t);
  laidOut(evalCheck(quiet, brokenJq(quiet)));
  const repo = withModule(t);
  put(repo, "maths", `${SAMPLE_ID}.json`, ofModule(example("sheet-eval-set")));
  const result = evalCheck(repo, brokenJq(repo));
  assert.equal(result.status, 2);
  assert.match(result.stderr, /^eval-check: jq cannot be run; the toolchain stage has it$/m);
});

test("a run leaves the tree as it was", (t) => {
  const repo = withModule(t);
  put(repo, "maths", "other-id.json", ofModule(example("sheet-eval-set")));
  const before = readdirSync(repo, { recursive: true }).sort();
  evalCheck(repo);
  assert.deepEqual(readdirSync(repo, { recursive: true }).sort(), before);
});
