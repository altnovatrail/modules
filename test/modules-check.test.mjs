// The tests of the recipe _modules-check (R03, R10, R11) before the first release: the tree of
// modules/, and every module in it checked, built and compared by module-kit — the real one of the
// toolchain stage on the PATH, and a stand-in where the test holds the calls to the forms of
// module-kit's architecture 3.5 and 3.6 (R102 spec). Against published versions — published.test.mjs.

import assert from "node:assert/strict";
import { mkdirSync, readdirSync, renameSync, symlinkSync, writeFileSync } from "node:fs";
import { basename, isAbsolute, join } from "node:path";
import { test } from "node:test";

import { EXAMPLES, PINNED, VECTORS, addModule, check, repository, stubModuleKit } from "./helpers.mjs";

// Every path under folder, so that a test can see that a run left nothing behind.
function listing(folder) {
  return readdirSync(folder, { recursive: true }).sort();
}

test("without modules/ there is nothing to check yet, and the recipe says so", (t) => {
  const repo = repository(t);
  const result = check(t, repo);
  assert.equal(result.status, 0, result.stderr);
  assert.match(result.stdout, /^modules\/: no modules yet$/m);
});

test("an empty modules/, which git never holds, is the tree with no module", (t) => {
  const repo = repository(t);
  mkdirSync(join(repo, "modules"));
  const result = check(t, repo);
  assert.equal(result.status, 0, result.stderr);
  assert.match(result.stdout, /^modules\/: no modules yet$/m);
});

test("each valid example of the format, as a module of its own in modules/, passes", (t) => {
  const repo = repository(t);
  const examples = readdirSync(EXAMPLES);
  assert.ok(examples.length > 0, `no examples in ${EXAMPLES}`);
  for (const example of examples) addModule(repo, join(EXAMPLES, example), example);
  const result = check(t, repo);
  assert.equal(result.status, 0, result.stdout + result.stderr);
  for (const example of examples) {
    assert.match(result.stdout, new RegExp(`^modules/${example}: the module breaks no rule`, "m"));
  }
});

test("a run leaves nothing behind in the tree: build writes into a temporary folder", (t) => {
  const repo = repository(t);
  addModule(repo, join(EXAMPLES, "with-media-and-notice"), "with-media-and-notice");
  const before = listing(repo);
  const result = check(t, repo);
  assert.equal(result.status, 0, result.stdout + result.stderr);
  assert.deepEqual(listing(repo), before);
});

test("a file in modules/ is refused", (t) => {
  const repo = repository(t);
  addModule(repo, join(EXAMPLES, "smallest"), "smallest");
  writeFileSync(join(repo, "modules", "README.md"), "Not a module.\n");
  const result = check(t, repo);
  assert.equal(result.status, 1);
  assert.match(result.stderr, /^modules\/README\.md: not a folder/m);
});

test("a folder named with a dot is refused, and so is a file", (t) => {
  for (const make of [
    (repo) => addModule(repo, join(EXAMPLES, "smallest"), ".draft", "novatrail/draft"),
    (repo) => writeFileSync(join(repo, "modules", ".keep"), ""),
  ]) {
    const repo = repository(t);
    addModule(repo, join(EXAMPLES, "smallest"), "smallest");
    make(repo);
    const result = check(t, repo);
    assert.equal(result.status, 1);
    assert.match(result.stderr, /^modules\/\.(draft|keep): .*nothing named with a dot/m);
  }
});

test("a link in modules/ is refused, even to a module", (t) => {
  const repo = repository(t);
  addModule(repo, join(EXAMPLES, "smallest"), "smallest");
  symlinkSync("smallest", join(repo, "modules", "alias"));
  const result = check(t, repo);
  assert.equal(result.status, 1);
  assert.match(result.stderr, /^modules\/alias: not a folder/m);
});

test("modules that is a file and not a folder is refused", (t) => {
  const repo = repository(t);
  writeFileSync(join(repo, "modules"), "");
  const result = check(t, repo);
  assert.equal(result.status, 1);
  assert.match(result.stderr, /^modules: not a folder/m);
});

test("a module whose ID is not novatrail/<folder> is refused", (t) => {
  for (const id of ["novatrail/other", "example/smallest"]) {
    const repo = repository(t);
    addModule(repo, join(EXAMPLES, "smallest"), "smallest", id);
    const result = check(t, repo);
    assert.equal(result.status, 1, id);
    assert.match(
      result.stderr,
      new RegExp(`^modules/smallest: the module's ID is ${id}, and the module in modules/smallest/ is novatrail/smallest`, "m"),
    );
  }
});

test("a module that module-kit check refuses is refused, by the ID of the rule it breaks", (t) => {
  const repo = repository(t);
  addModule(repo, join(VECTORS, "module.license", "no-license"), "no-license");
  const result = check(t, repo);
  assert.equal(result.status, 1);
  assert.match(result.stdout, /^modules\/no-license: module\.license: /m);
});

test("module-kit is called in the forms of its architecture 3.5 and 3.6: check, build into a new folder, diff", (t) => {
  const repo = repository(t);
  const stub = stubModuleKit(t);
  addModule(repo, join(EXAMPLES, "smallest"), "a");
  addModule(repo, join(EXAMPLES, "smallest"), "b");
  const result = check(t, repo, stub.env);
  assert.equal(result.status, 0, result.stderr);
  const calls = stub.calls();
  assert.deepEqual(
    calls.map((call) => call.args.slice(0, 2)),
    [
      ["check", "modules/a"], ["build", "modules/a"], ["diff", "--first"],
      ["check", "modules/b"], ["build", "modules/b"], ["diff", "--first"],
    ],
  );
  const built = calls.filter((call) => call.args[0] === "build").map((call) => call.args[2]);
  assert.deepEqual(
    calls.filter((call) => call.args[0] === "diff").map((call) => call.args),
    built.map((folder) => ["diff", "--first", folder]),
    "diff compares the package build wrote, and before the first release with --first",
  );
  for (const call of calls) {
    if (call.args[0] === "check") {
      assert.equal(call.args.length, 2);
      continue;
    }
    if (call.args[0] === "diff") continue;
    assert.equal(call.args.length, 3);
    assert.ok(isAbsolute(call.args[2]), call.args[2]);
    assert.ok(!call.args[2].startsWith(repo), `build writes into ${call.args[2]}, inside the tree`);
    assert.equal(basename(call.args[2]), basename(call.args[1]));
    assert.equal(call.existed, false, "build's folder must not be there yet");
  }
});

test("a refusal does not stop the check of the other modules, and fails the recipe", (t) => {
  const repo = repository(t);
  const stub = stubModuleKit(t);
  addModule(repo, join(EXAMPLES, "smallest"), "a");
  addModule(repo, join(EXAMPLES, "smallest"), "b");
  const result = check(t, repo, { ...stub.env, STUB_CHECK_STATUS: JSON.stringify({ a: 1 }) });
  assert.equal(result.status, 1);
  assert.deepEqual(
    stub.calls().map((call) => call.args.slice(0, 2)),
    [["check", "modules/a"], ["check", "modules/b"], ["build", "modules/b"], ["diff", "--first"]],
  );
});

test("module-kit that cannot do its work, status 2, gives the recipe its status", (t) => {
  const repo = repository(t);
  const stub = stubModuleKit(t);
  addModule(repo, join(EXAMPLES, "smallest"), "a");
  const result = check(t, repo, { ...stub.env, STUB_CHECK_STATUS: JSON.stringify({ a: 2 }) });
  assert.equal(result.status, 2);
  assert.deepEqual(stub.calls().map((call) => call.args[0]), ["check"]);
});

test("modules that is a link, even to a folder, is refused", (t) => {
  const repo = repository(t);
  mkdirSync(join(repo, "elsewhere"));
  addModule(repo, join(EXAMPLES, "smallest"), "smallest");
  // modules/ moves aside, and a link takes its name.
  renameSync(join(repo, "modules"), join(repo, "elsewhere", "modules"));
  symlinkSync(join("elsewhere", "modules"), join(repo, "modules"));
  const result = check(t, repo);
  assert.equal(result.status, 1);
  assert.match(result.stderr, /^modules: not a folder/m);
});

test("a build that fails after a passing check fails the recipe, and the other modules are still checked", (t) => {
  const repo = repository(t);
  const stub = stubModuleKit(t);
  addModule(repo, join(EXAMPLES, "smallest"), "a");
  addModule(repo, join(EXAMPLES, "smallest"), "b");
  const result = check(t, repo, { ...stub.env, STUB_BUILD_STATUS: JSON.stringify({ a: 2 }) });
  assert.equal(result.status, 2);
  assert.deepEqual(
    stub.calls().map((call) => call.args.slice(0, 2)),
    [["check", "modules/a"], ["build", "modules/a"], ["check", "modules/b"], ["build", "modules/b"], ["diff", "--first"]],
  );
});

test("a package whose inventory.json names no ID is refused", (t) => {
  const repo = repository(t);
  const stub = stubModuleKit(t);
  addModule(repo, join(EXAMPLES, "smallest"), "a");
  const result = check(t, repo, { ...stub.env, STUB_NO_ID: JSON.stringify(["a"]) });
  assert.equal(result.status, 2);
  assert.match(result.stderr, /^modules\/a: the package module-kit built names no ID in its inventory\.json$/m);
});

test("a module-kit on the PATH other than the one the Dockerfile pins checks nothing", (t) => {
  const repo = repository(t);
  const stub = stubModuleKit(t);
  addModule(repo, join(EXAMPLES, "smallest"), "a");
  const result = check(t, repo, { ...stub.env, STUB_VERSION: "v0.0.1" });
  assert.equal(result.status, 1);
  assert.match(result.stderr, new RegExp(`says "module-kit v0\\.0\\.1", and \\.devcontainer/Dockerfile pins module-kit ${PINNED.replaceAll(".", "\\.")}`));
  assert.deepEqual(stub.calls(), []);
});
