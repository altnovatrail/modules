// The tests of eval-match (amendment T59 to R25): each sample of a module's eval set held to the
// package of the commit. The real module-kit and the real generator of the toolchain stage — by
// CANVAS_LIB where the container has not been rebuilt since T58 — on copies of grade-3-math and of
// the format's example with a template; a stand-in for the generator where a test holds the calls
// and what is compared. The samples of people here are made up from synthetic ones, by their
// fields alone: no sample of a person is ever read (П-12).

import assert from "node:assert/strict";
import { mkdirSync, readFileSync, rmSync, writeFileSync } from "node:fs";
import { join } from "node:path";
import { test } from "node:test";

import { EXAMPLES, ROOT, addModule, mustRun, pinGenerator, repository, run, stubModuleKit, temporaryFolder, withSpec } from "./helpers.mjs";

const STAGE = process.env.CANVAS_LIB ? { CANVAS_LIB: process.env.CANVAS_LIB } : {};

function evalMatch(repo, env = STAGE) {
  return run(repo, "eval-match", [], env);
}

function read(path) {
  return JSON.parse(readFileSync(path, "utf8"));
}

function write(path, value) {
  mkdirSync(join(path, ".."), { recursive: true });
  writeFileSync(path, `${JSON.stringify(value, null, 2)}\n`);
}

// A person's sample made up of a synthetic one: the same task, another origin and ID.
function asPerson(sample, id, task = {}) {
  const { synthetic, ...rest } = sample;
  return { ...rest, id, person: {}, task: { ...rest.task, ...task } };
}

// The lines of stderr that name a file of eval/.
function refusals(result) {
  return result.stderr.split("\n").filter((line) => line.startsWith("modules/"));
}

// A repository with a copy of grade-3-math and two of its synthetic samples, drawn by the stage's
// generator on its package.
function withGrade3(t) {
  const repo = withSpec(repository(t));
  mkdirSync(join(repo, "modules"));
  // The module without its own eval set, whose samples would meet the test's.
  rmSync(join(addModule(repo, join(ROOT, "modules", "grade-3-math"), "grade-3-math"), "eval"), { recursive: true, force: true });
  const made = run(repo, "_synthetic", ["grade-3-math", "v0.2.1", "times-tables-d1:1", "comparing-fractions-d1:2"], STAGE);
  assert.equal(made.status, 0, made.stderr);
  return { repo, eval: join(repo, "modules", "grade-3-math", "eval") };
}

test("synthetic samples of the package, and a person's sample of an open view it has, pass", (t) => {
  const { repo, eval: folder } = withGrade3(t);
  write(join(folder, "person-a.json"), asPerson(read(join(folder, "synthetic-1.json")), "person-a"));
  const result = evalMatch(repo);
  assert.equal(result.status, 0, result.stderr);
  assert.match(result.stdout, /^_match-synthetic: 2 synthetic samples of novatrail\/grade-3-math are drawn again the same/m);
  assert.match(result.stdout, /^_match-people: 1 samples of people of novatrail\/grade-3-math are of open views/m);
  assert.match(result.stdout, /^eval-match: every sample of every eval set is one of the package of this commit$/m);
});

test("a change of a task's open view refuses its samples, and only theirs: the package's digest is not compared", (t) => {
  const { repo, eval: folder } = withGrade3(t);
  write(join(folder, "person-a.json"), asPerson(read(join(folder, "synthetic-2.json")), "person-a"));
  write(join(folder, "person-b.json"), asPerson(read(join(folder, "synthetic-1.json")), "person-b"));
  const task = join(repo, "modules", "grade-3-math", "tasks", "comparing-fractions-d1.yaml");
  const text = readFileSync(task, "utf8");
  const changed = text.replace("Write <, > or = between the two fractions.", "Write <, > or = between the two fractions on each line.");
  assert.notEqual(changed, text, "the statement the test changes is there");
  writeFileSync(task, changed);
  const result = evalMatch(repo);
  assert.equal(result.status, 1, result.stdout + result.stderr);
  assert.deepEqual(refusals(result), [
    "modules/grade-3-math/eval/synthetic-2.json: is not what canvas's generator draws of comparing-fractions-d1:2 on this commit's package, with the tag v0.2.1; re-make it: just synthetic grade-3-math comparing-fractions-d1:2",
    "modules/grade-3-math/eval/person-a.json: the open view of comparing-fractions-d1 is not the one the sample was made on; remove the sample (amendment 2026-11-06 to R04)",
  ]);
  assert.match(result.stdout, /^_match-synthetic: 1 synthetic samples/m, "the sample of the other task is drawn the same");
  assert.match(result.stdout, /^_match-people: 1 samples of people/m);
});

test("a synthetic sample is drawn again with the tag it names, and a sample of a task the package lacks is refused", (t) => {
  const { repo, eval: folder } = withGrade3(t);
  const sample = read(join(folder, "synthetic-1.json"));
  write(join(folder, "synthetic-1.json"), { ...sample, task: { ...sample.task, modules: { ...sample.task.modules, tag: "v0.1.0" } } });
  write(join(folder, "person-a.json"), asPerson(sample, "person-a", { task: "no-such-task" }));
  const result = evalMatch(repo);
  assert.equal(result.status, 1, result.stdout + result.stderr);
  const lines = refusals(result);
  assert.equal(lines.length, 1, result.stderr);
  assert.equal(lines[0], "modules/grade-3-math/eval/person-a.json: the package of this commit has no open view of no-such-task; remove the sample (amendment 2026-11-06 to R04)");
  assert.match(result.stdout, /^_match-synthetic: 2 synthetic samples/m, "synthetic-1 with the tag v0.1.0 is drawn again with v0.1.0");
});

test("a person's sample of a template is held to the instance of its seed", (t) => {
  const repo = withSpec(repository(t));
  mkdirSync(join(repo, "modules"));
  const module = addModule(repo, join(EXAMPLES, "with-templates"), "maths");
  const work = temporaryFolder(t);
  mustRun("module-kit", ["build", module, join(work, "package")]);
  mustRun("module-kit", ["index", "--first", "--tag", "v0.0.0", "--no-withdrawals", "--out", join(work, "index.json"), join(work, "package")]);
  // Seeds past 2^53 too, which a comparison of numbers would take for one another.
  const big = ["18446744073709551614", "18446744073709551615"];
  mustRun("module-kit", ["open", "--index", join(work, "index.json"), "--latest", join(work, "index.json"), "--seed", "3", ...big.flatMap((seed) => ["--seed", seed]), "--out", join(work, "open"), join(work, "package")]);
  const views = read(join(work, "open", "open-views.json")).modules[0].tasks;
  const of = (seed) => views.find((view) => view.task.id === "missing-factor" && view.seed === seed).open_view_digest;
  const three = of("3");
  assert.notEqual(of(big[0]), of(big[1]), "the two instances differ");
  const sample = (id, task) => ({ id, person: {}, task: { modules: { tag: "v0.1.0", module: "novatrail/maths", version: "0.1.0", package_digest: `sha256:${"0".repeat(64)}` }, task: "missing-factor", ...task } });
  const folder = join(module, "eval");
  write(join(folder, "person-a.json"), sample("person-a", { seed: "3", open_view_digest: three }));
  write(join(folder, "person-b.json"), sample("person-b", { seed: "4", open_view_digest: three }));
  write(join(folder, "person-c.json"), sample("person-c", { open_view_digest: three }));
  write(join(folder, "person-d.json"), sample("person-d", { seed: big[0], open_view_digest: of(big[0]) }));
  write(join(folder, "person-e.json"), sample("person-e", { seed: big[1], open_view_digest: of(big[1]) }));
  const result = evalMatch(repo);
  assert.equal(result.status, 1, result.stdout + result.stderr);
  assert.deepEqual(refusals(result), [
    "modules/maths/eval/person-b.json: the open view of missing-factor with the seed 4 is not the one the sample was made on; remove the sample (amendment 2026-11-06 to R04)",
    "modules/maths/eval/person-c.json: the package of this commit has no open view of missing-factor; remove the sample (amendment 2026-11-06 to R04)",
  ]);
  assert.match(result.stdout, /^_match-people: 3 samples of people of novatrail\/maths/m);
});

test("a file whose sample cannot be held to the package is refused, and nothing is built for it alone", (t) => {
  const repo = repository(t);
  const kit = stubModuleKit(t);
  const folder = join(repo, "modules", "maths", "eval");
  const task = (fields) => ({ modules: { tag: "v0.1.0", module: "novatrail/maths" }, task: "a-task", ...fields });
  write(join(folder, "empty.json"), {});
  write(join(folder, "other.json"), { id: "other", person: {}, task: { ...task(), modules: { module: "novatrail/other" } } });
  write(join(folder, "nobody.json"), { id: "nobody", task: task() });
  write(join(folder, "seedless.json"), { id: "seedless", synthetic: {}, task: task() });
  write(join(folder, "untagged.json"), { id: "untagged", synthetic: { seed: "1" }, task: task({ modules: { module: "novatrail/maths" } }) });
  write(join(folder, "bad-task.json"), { id: "bad-task", person: {}, task: task({ task: "A task" }) });
  write(join(folder, "deep", "nested.json"), { id: "nested", person: {}, task: task() });
  write(join(folder, "synthetic-9.json"), { id: "synthetic-9", person: {}, task: task() });
  writeFileSync(join(folder, "two.json"), "{}\n{}\n");
  const result = evalMatch(repo, kit.env);
  assert.equal(result.status, 1, result.stdout + result.stderr);
  assert.deepEqual(refusals(result), [
    "modules/maths/eval/bad-task.json: names no task ID and seed that can be held to the package",
    "modules/maths/eval/deep/nested.json: outside the layout eval/<id>.json, which eval-check refuses, so it is not held to the package",
    "modules/maths/eval/empty.json: holds no sample whose task and origin can be read, so it cannot be held to the package",
    "modules/maths/eval/nobody.json: neither a person's sample nor a synthetic one, so it cannot be held to the package",
    "modules/maths/eval/other.json: a sample of the eval set of novatrail/maths is of a task of that module",
    "modules/maths/eval/seedless.json: a synthetic sample with no seed, or no tag vX.Y.Z of modules, cannot be drawn again",
    "modules/maths/eval/synthetic-9.json: the ID of a person's sample begins with synthetic-, which only a synthetic sample's does, and such a file a Claude Code session may read (amendment T61 to R04)",
    "modules/maths/eval/two.json: holds no sample whose task and origin can be read, so it cannot be held to the package",
    "modules/maths/eval/untagged.json: a synthetic sample with no seed, or no tag vX.Y.Z of modules, cannot be drawn again",
  ]);
  assert.deepEqual(kit.calls(), [], "no sample to hold, so no package is built");
});

test("without a sample nothing is built, and without modules/ there is nothing to hold", (t) => {
  const repo = repository(t);
  const kit = stubModuleKit(t);
  mkdirSync(join(repo, "modules", "maths", "eval"), { recursive: true });
  writeFileSync(join(repo, "modules", "maths", "eval", "notes.txt"), "Not a sample.\n");
  mkdirSync(join(repo, "modules", "essays"));
  const result = evalMatch(repo, kit.env);
  assert.equal(result.status, 0, result.stderr);
  assert.equal(result.stdout, "eval-match: no module holds a sample\n");
  assert.deepEqual(kit.calls(), []);

  const bare = evalMatch(repository(t), kit.env);
  assert.equal(bare.status, 0, bare.stderr);
  assert.equal(bare.stdout, "eval-match: no modules/, so no sample\n");
});

test("a module-kit other than the pinned one holds nothing, with the status 2", (t) => {
  const repo = repository(t);
  const kit = stubModuleKit(t);
  write(join(repo, "modules", "maths", "eval", "a.json"), { id: "a", person: {}, task: { modules: { module: "novatrail/maths" }, task: "a-task" } });
  const result = evalMatch(repo, { ...kit.env, STUB_VERSION: "v0.0.1" });
  assert.equal(result.status, 2);
  assert.match(result.stderr, /says "module-kit v0\.0\.1", and \.devcontainer\/Dockerfile pins module-kit .*: rebuild the container/);
  assert.deepEqual(kit.calls(), []);
});

// A stand-in for the generator that logs its calls into STUB_GENERATOR_LOG and writes, for each
// <task>:<seed>, synthetic-<seed>.json: its task, tag and seed, and a version and package digest of
// its own, which the comparison leaves out.
const DRAWER = `import { appendFileSync, mkdirSync, writeFileSync } from "node:fs";
import { join } from "node:path";
const [, pack, tag, folder, ...calls] = process.argv.slice(2);
appendFileSync(process.env.STUB_GENERATOR_LOG, JSON.stringify({ tag, calls }) + "\\n");
mkdirSync(folder, { recursive: true });
for (const call of calls) {
  const [task, seed] = call.split(":");
  const sample = { id: "synthetic-" + seed, synthetic: { seed }, task: { modules: { tag, module: "novatrail/maths", version: "9.9.9", package_digest: "sha256:drawn" }, task } };
  writeFileSync(join(folder, "synthetic-" + seed + ".json"), JSON.stringify(sample));
}
`;

test("the generator is called on 50 samples at a time, a tag each, and only the package's version and digest are left out", (t) => {
  const repo = repository(t);
  const lib = pinGenerator(t, repo, DRAWER);
  const log = join(lib, "calls.jsonl");
  const rows = join(temporaryFolder(t), "rows");
  const lines = [];
  for (let seed = 0; seed < 123; seed++) {
    const tag = seed < 3 ? "v0.1.0" : "v0.2.1";
    const file = `modules/maths/eval/synthetic-${seed}.json`;
    // The file names another version and digest than the generator draws, which pass; seed 7 names
    // another task, which does not.
    const task = seed === 7 ? "b-task" : "a-task";
    write(join(repo, file), { id: `synthetic-${seed}`, synthetic: { seed: String(seed) }, task: { modules: { tag, module: "novatrail/maths", version: "0.1.0", package_digest: "sha256:made" }, task } });
    lines.push([file, "a-task", String(seed), tag].join("\x1f"));
  }
  writeFileSync(rows, `${lines.join("\n")}\n`);
  const result = run(repo, "_match-synthetic", ["maths", join(repo, "package"), rows], { CANVAS_LIB: lib, STUB_GENERATOR_LOG: log });
  assert.equal(result.status, 1, result.stderr);
  const calls = readFileSync(log, "utf8").trim().split("\n").map((line) => JSON.parse(line));
  assert.deepEqual(calls.map(({ tag, calls }) => [tag, calls.length]), [["v0.1.0", 3], ["v0.2.1", 50], ["v0.2.1", 50], ["v0.2.1", 20]]);
  assert.deepEqual(refusals(result), [
    "modules/maths/eval/synthetic-7.json: is not what canvas's generator draws of a-task:7 on this commit's package, with the tag v0.2.1; re-make it: just synthetic maths a-task:7",
  ]);
  assert.match(result.stdout, /^_match-synthetic: 122 synthetic samples of novatrail\/maths are drawn again the same/m);
});

test("the generator's refusal of a call is the recipe's, with its reason, for the sample of that call", (t) => {
  const repo = repository(t);
  const lib = pinGenerator(t, repo, `console.error("a-task:1: novatrail/maths has no task a-task"); process.exitCode = 1;\n`);
  const rows = join(temporaryFolder(t), "rows");
  writeFileSync(rows, `${["modules/maths/eval/synthetic-1.json", "a-task", "1", "v0.2.1"].join("\x1f")}\n`);
  const result = run(repo, "_match-synthetic", ["maths", join(repo, "package"), rows], { CANVAS_LIB: lib });
  assert.equal(result.status, 1);
  assert.match(result.stderr, /^a-task:1: novatrail\/maths has no task a-task$/m);
  assert.deepEqual(refusals(result), [
    "modules/maths/eval/synthetic-1.json: canvas's generator draws a-task:1 no more on this commit's package of novatrail/maths, with the tag v0.2.1, for the reason above; re-make or remove it",
  ]);
});

test("a batch the generator refuses is drawn again a sample at a time, each with its own reason", (t) => {
  const repo = repository(t);
  const lib = pinGenerator(t, repo, `const calls = process.argv.slice(6);
if (calls.some((call) => call.startsWith("gone:"))) { console.error("gone:2: novatrail/maths has no task gone"); process.exitCode = 1; }
else { await import(${JSON.stringify("data:text/javascript," + encodeURIComponent(DRAWER))}); }
`);
  const rows = join(temporaryFolder(t), "rows");
  const lines = [];
  for (const [seed, task] of [[1, "a-task"], [2, "gone"], [3, "a-task"]]) {
    const file = `modules/maths/eval/synthetic-${seed}.json`;
    write(join(repo, file), { id: `synthetic-${seed}`, synthetic: { seed: String(seed) }, task: { modules: { tag: "v0.2.1", module: "novatrail/maths" }, task } });
    lines.push([file, task, String(seed), "v0.2.1"].join("\x1f"));
  }
  writeFileSync(rows, `${lines.join("\n")}\n`);
  const log = join(lib, "calls.jsonl");
  const result = run(repo, "_match-synthetic", ["maths", join(repo, "package"), rows], { CANVAS_LIB: lib, STUB_GENERATOR_LOG: log });
  assert.equal(result.status, 1, result.stderr);
  assert.deepEqual(refusals(result), [
    "modules/maths/eval/synthetic-2.json: canvas's generator draws gone:2 no more on this commit's package of novatrail/maths, with the tag v0.2.1, for the reason above; re-make or remove it",
  ]);
  assert.match(result.stderr, /^gone:2: novatrail\/maths has no task gone$/m);
  assert.match(result.stdout, /^_match-synthetic: 2 synthetic samples of novatrail\/maths are drawn again the same/m);
});

test("module-kit that cannot build a module, or open its package, gives 2, and the other modules are still held", (t) => {
  const kit = stubModuleKit(t);
  const person = (module) => ({ id: "a", person: {}, task: { modules: { module: `novatrail/${module}` }, task: "a-task", open_view_digest: `sha256:${"0".repeat(64)}` } });
  const repo = repository(t);
  write(join(repo, "modules", "alpha", "eval", "a.json"), person("alpha"));
  write(join(repo, "modules", "beta", "eval", "a.json"), person("beta"));
  const result = evalMatch(repo, { ...kit.env, STUB_BUILD_STATUS: JSON.stringify({ alpha: 2 }) });
  assert.equal(result.status, 2, result.stdout + result.stderr);
  assert.match(result.stderr, /^eval-match: module-kit cannot build modules\/alpha, so its samples are not held to its package$/m);
  // The stand-in's index exits 2, so beta's package is built and cannot be opened.
  assert.match(result.stderr, /^_match-people: module-kit cannot open the package of novatrail\/beta, so its samples of people are not held to it$/m);
  assert.deepEqual(kit.calls().map(({ args }) => [args[0], args[1]]), [["build", "modules/alpha"], ["build", "modules/beta"], ["index", "--first"]]);
});

test("jq that cannot run holds nothing, with the status 2, and refuses no sample", (t) => {
  const repo = repository(t);
  write(join(repo, "modules", "maths", "eval", "a.json"), { id: "a", person: {}, task: { modules: { module: "novatrail/maths" }, task: "a-task" } });
  const bin = join(repo, "bin");
  mkdirSync(bin);
  writeFileSync(join(bin, "jq"), "#!/bin/sh\nexit 127\n", { mode: 0o755 });
  const result = evalMatch(repo, { PATH: `${bin}:${process.env.PATH}` });
  assert.equal(result.status, 2);
  assert.match(result.stderr, /^eval-match: jq cannot be run; the toolchain stage has it$/m);
  assert.deepEqual(refusals(result), []);
});
