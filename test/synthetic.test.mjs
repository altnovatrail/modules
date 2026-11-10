// The tests of the recipe of synthetic samples (R29): _synthetic in made-up repositories, with a
// stand-in for canvas's generator whose sum the made-up Dockerfile pins, and module-kit's stand-in
// where a test holds the form of the calls or a status; then the real generator of the toolchain
// stage and the real module-kit on a copy of grade-3-math, whose sample spec's sample-check holds
// and which a second call writes again byte for byte. The tag is given, or, empty, taken from
// _newest-release on the stand-ins of git and curl: the network is no part of these tests.

import assert from "node:assert/strict";
import { createHash } from "node:crypto";
import { existsSync, mkdirSync, readFileSync, readdirSync, rmSync, symlinkSync, writeFileSync } from "node:fs";
import { join } from "node:path";
import { test } from "node:test";

import { EXAMPLES, ROOT, addModule, repository, run, stubModuleKit, stubNetwork, temporaryFolder, withSpec } from "./helpers.mjs";

// The tag of canvas that this repository's Dockerfile pins, whose generator the stage installs.
const CANVAS = /^ARG CANVAS_VERSION=(.*)$/m.exec(readFileSync(join(ROOT, ".devcontainer", "Dockerfile"), "utf8"))[1];
const GENERATOR = `synthetic-module_${CANVAS}.mjs`;

// A stand-in for the generator: each call goes into STUB_GENERATOR_LOG as a line of JSON, its
// arguments and whether the package's inventory.json was there; it prints a line for each call of
// a sample, STUB_GENERATOR_STDERR on stderr, and exits with STUB_GENERATOR_STATUS, 0 unless given.
const STUB = `import { appendFileSync, existsSync } from "node:fs";
import { join } from "node:path";
const args = process.argv.slice(2);
appendFileSync(process.env.STUB_GENERATOR_LOG, JSON.stringify({ args, built: existsSync(join(args[1], "inventory.json")) }) + "\\n");
for (const call of args.slice(4)) console.log("synthetic-" + call.split(":")[1] + "\\t" + call.split(":")[0] + "\\ta page");
if (process.env.STUB_GENERATOR_STDERR) console.error(process.env.STUB_GENERATOR_STDERR);
process.exitCode = Number(process.env.STUB_GENERATOR_STATUS || 0);
`;

function sha256(text) {
  return createHash("sha256").update(text).digest("hex");
}

// A made-up repository with a module of the format in modules/maths/, the stand-ins of module-kit
// and of the generator, and its Dockerfile pinning the stand-in's sum, or the sum given.
function withGenerator(t, { pinned = sha256(STUB) } = {}) {
  const repo = repository(t);
  addModule(repo, firstExample(), "maths");
  const dockerfile = join(repo, ".devcontainer", "Dockerfile");
  writeFileSync(dockerfile, readFileSync(dockerfile, "utf8").replace(/^ARG CANVAS_SYNTHETIC_SHA256=.*$/m, `ARG CANVAS_SYNTHETIC_SHA256=${pinned}`));
  const lib = temporaryFolder(t);
  writeFileSync(join(lib, GENERATOR), STUB);
  const log = join(lib, "calls.jsonl");
  const kit = stubModuleKit(t);
  return {
    repo,
    lib,
    kit,
    env: { ...kit.env, CANVAS_LIB: lib, STUB_GENERATOR_LOG: log },
    calls: () => (existsSync(log) ? readFileSync(log, "utf8").trim().split("\n").filter(Boolean).map((line) => JSON.parse(line)) : []),
  };
}

// The first module of the format's valid examples.
function firstExample() {
  return join(EXAMPLES, readdirSync(EXAMPLES).sort()[0]);
}

function synthetic(repo, args, env) {
  return run(repo, "_synthetic", args, env);
}

test("the package is built first, and the generator gets it, the tag, the module's eval/ and the calls", (t) => {
  const { repo, env, kit, calls } = withGenerator(t);
  const result = synthetic(repo, ["maths", "v0.2.1", "a-task:1", "b-task:20"], env);
  assert.equal(result.status, 0, result.stderr);
  const builds = kit.calls();
  assert.equal(builds.length, 1);
  assert.deepEqual(builds[0].args.slice(0, 2), ["build", "modules/maths"]);
  assert.equal(builds[0].existed, false, "build's folder is not there yet");
  const [call, ...more] = calls();
  assert.deepEqual(more, []);
  assert.deepEqual(call.args, ["--module", builds[0].args[2], "v0.2.1", "modules/maths/eval", "a-task:1", "b-task:20"]);
  assert.equal(call.built, true, "the generator gets the package build wrote");
  assert.equal(result.stdout.split("\n").filter((line) => line.startsWith("synthetic-")).join("\n"), "synthetic-1\ta-task\ta page\nsynthetic-20\tb-task\ta page");
  assert.equal(existsSync(builds[0].args[2]), false, "the package goes when the recipe ends");
});

test("the generator's refusal and its reasons are the recipe's, and so is a call out of its form", (t) => {
  for (const status of [1, 2]) {
    const { repo, env } = withGenerator(t);
    const reason = "ages:3: ages is a template: its instance by a seed is module-kit's";
    const result = synthetic(repo, ["maths", "v0.2.1", "ages:3"], { ...env, STUB_GENERATOR_STATUS: String(status), STUB_GENERATOR_STDERR: reason });
    assert.equal(result.status, status, result.stderr);
    assert.match(result.stderr, new RegExp(`^${reason}$`, "m"));
  }
});

test("a generator other than the one the Dockerfile pins, or none, writes nothing: rebuild the container", (t) => {
  const other = withGenerator(t, { pinned: "f".repeat(64) });
  const result = synthetic(other.repo, ["maths", "v0.2.1", "a-task:1"], other.env);
  assert.equal(result.status, 2);
  assert.match(result.stderr, new RegExp(`${GENERATOR.replaceAll(".", "\\.")} is [0-9a-f]{64}, and \\.devcontainer/Dockerfile pins f{64}: rebuild the container`));
  assert.deepEqual(other.calls(), []);
  assert.deepEqual(other.kit.calls(), [], "nothing is built");

  const none = withGenerator(t);
  rmSync(join(none.lib, GENERATOR));
  const missing = synthetic(none.repo, ["maths", "v0.2.1", "a-task:1"], none.env);
  assert.equal(missing.status, 2);
  assert.match(missing.stderr, /is not there; the toolchain stage installs canvas's generator .*: rebuild the container/);
  assert.deepEqual(none.kit.calls(), []);
});

test("a module-kit on the PATH other than the one the Dockerfile pins builds nothing", (t) => {
  const { repo, env, kit, calls } = withGenerator(t);
  const result = synthetic(repo, ["maths", "v0.2.1", "a-task:1"], { ...env, STUB_VERSION: "v0.0.1" });
  assert.equal(result.status, 2);
  assert.match(result.stderr, /says "module-kit v0\.0\.1", and \.devcontainer\/Dockerfile pins module-kit .*: rebuild the container/);
  assert.deepEqual(kit.calls(), []);
  assert.deepEqual(calls(), []);
});

test("module-kit's refusal of the module is the recipe's, and the generator is not called", (t) => {
  const { repo, env, calls } = withGenerator(t);
  const result = synthetic(repo, ["maths", "v0.2.1", "a-task:1"], { ...env, STUB_BUILD_STATUS: JSON.stringify({ maths: 1 }) });
  assert.equal(result.status, 1);
  assert.deepEqual(calls(), []);
});

test("a package whose ID is not novatrail/<module> is refused, and one that names none cannot be used", (t) => {
  const named = withGenerator(t);
  addModule(named.repo, firstExample(), "other", "novatrail/maths");
  // The real module-kit of the toolchain stage builds it, so that the ID is the module's own.
  const result = synthetic(named.repo, ["other", "v0.2.1", "a-task:1"], { ...named.env, PATH: process.env.PATH });
  assert.equal(result.status, 1, result.stderr);
  assert.match(result.stderr, /the module's ID is novatrail\/maths, and the module in modules\/other\/ is novatrail\/other \(R03\)/);
  assert.deepEqual(named.calls(), []);

  const unnamed = withGenerator(t);
  const none = synthetic(unnamed.repo, ["maths", "v0.2.1", "a-task:1"], { ...unnamed.env, STUB_NO_ID: JSON.stringify(["maths"]) });
  assert.equal(none.status, 2);
  assert.match(none.stderr, /names no ID in its inventory\.json/);
  assert.deepEqual(unnamed.calls(), []);
});

test("no sample to write, no module folder, an eval/ that is no folder, or no tag vX.Y.Z is out of form, with nothing built", (t) => {
  const { repo, env, kit, calls } = withGenerator(t);
  writeFileSync(join(repo, "modules", "file"), "");
  addModule(repo, firstExample(), "linked");
  symlinkSync(temporaryFolder(t), join(repo, "modules", "linked", "eval"));
  addModule(repo, firstExample(), "flat");
  writeFileSync(join(repo, "modules", "flat", "eval"), "");
  for (const [args, message] of [
    [["maths", "v0.2.1"], /name the samples to write/],
    [["", "v0.2.1", "a-task:1"], /modules\/ is no folder of a module in modules\/ \(R03\)/],
    [["linked", "v0.2.1", "a-task:1"], /modules\/linked\/eval is a link or no folder/],
    [["flat", "v0.2.1", "a-task:1"], /modules\/flat\/eval is a link or no folder/],
    [["absent", "v0.2.1", "a-task:1"], /modules\/absent is no folder of a module in modules\/ \(R03\)/],
    [["file", "v0.2.1", "a-task:1"], /modules\/file is no folder of a module/],
    [["../modules/maths", "v0.2.1", "a-task:1"], /is no folder of a module/],
    [["maths", "0.2.1", "a-task:1"], /0\.2\.1 is no tag vX\.Y\.Z of a release of modules/],
  ]) {
    const result = synthetic(repo, args, env);
    assert.equal(result.status, 2, args.join(" "));
    assert.match(result.stderr, message);
  }
  assert.deepEqual(kit.calls(), []);
  assert.deepEqual(calls(), []);
});

test("with no tag given, the tag is the newest release's, asked of the network after every local check", (t) => {
  const { repo, env, calls } = withGenerator(t);
  const network = stubNetwork(t, { tags: ["v0.2.1", "v0.3.0"], indexes: { "v0.2.1": 200 } });
  const both = { ...env, PATH: `${network.env.PATH.split(":")[0]}:${env.PATH}` };
  const result = run(repo, "synthetic", ["maths", "a-task:1"], both);
  assert.equal(result.status, 0, result.stderr);
  assert.equal(calls()[0].args[2], "v0.2.1");

  const before = network.calls().length;
  const refused = run(repo, "synthetic", ["absent", "a-task:1"], both);
  assert.equal(refused.status, 2);
  assert.equal(network.calls().length, before, "a module that is not there asks nothing of the network");

  const none = stubNetwork(t, { tags: [] });
  const unreleased = run(repo, "synthetic", ["maths", "a-task:1"], { ...env, PATH: `${none.env.PATH.split(":")[0]}:${env.PATH}` });
  assert.equal(unreleased.status, 2);
  assert.match(unreleased.stderr, /modules has no release yet/);
});

test("synthetic with no sample to write is out of form before it asks the network", (t) => {
  const result = run(repository(t), "synthetic", ["maths"]);
  assert.equal(result.status, 2);
  assert.match(result.stderr, /name the module and the samples to write/);
});

test("the generator of the stage writes a sample of grade-3-math that sample-check holds, and the same bytes again", (t) => {
  const repo = withSpec(repository(t));
  mkdirSync(join(repo, "modules"));
  addModule(repo, join(ROOT, "modules", "grade-3-math"), "grade-3-math");
  const env = process.env.CANVAS_LIB ? { CANVAS_LIB: process.env.CANVAS_LIB } : {};
  const first = synthetic(repo, ["grade-3-math", "v0.2.1", "times-tables-d1:5"], env);
  assert.equal(first.status, 0, first.stderr);
  assert.match(first.stdout, /^synthetic-5\ttimes-tables-d1\t/m);
  const file = join(repo, "modules", "grade-3-math", "eval", "synthetic-5.json");
  const bytes = readFileSync(file);
  const again = synthetic(repo, ["grade-3-math", "v0.2.1", "times-tables-d1:5"], env);
  assert.equal(again.status, 0, again.stderr);
  assert.ok(readFileSync(file).equals(bytes), "the same task and seed give the same file");
  const sample = JSON.parse(bytes);
  assert.equal(sample.task.modules.module, "novatrail/grade-3-math");
  assert.equal(sample.task.modules.tag, "v0.2.1");
  const checked = run(repo, "eval-check");
  assert.equal(checked.status, 0, checked.stdout + checked.stderr);
});
