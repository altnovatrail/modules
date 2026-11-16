// The tests of what the job "Gate" gives sim's gate (R37), without the network and without sim:
// gate-inputs with no module and before the first release, on stand-ins of git and curl, and
// _gate-args against a release made here by the real module-kit, as test/published.test.mjs makes
// one. The arguments are held to the call form of sim's README at the tag ci.yml pins; that sim
// takes them, the job's first run on that tag showed (RUN.md, T78), and every run of the job holds.

import assert from "node:assert/strict";
import { cpSync, existsSync, mkdirSync, readdirSync, readFileSync, symlinkSync, writeFileSync } from "node:fs";
import { join } from "node:path";
import { test } from "node:test";

import { EXAMPLES, addModule, firstRelease, published, repository, run, stubModuleKit, stubNetwork } from "./helpers.mjs";

const args = (folder) => readFileSync(join(folder, "args"), "utf8").split("\n").slice(0, -1);

test("with no module the gate has nothing to read: args is empty, and the network is not asked", (t) => {
  const repo = repository(t);
  const network = stubNetwork(t, { tags: ["v0.1.0"], indexes: { "v0.1.0": 200 } });
  const result = run(repo, "gate-inputs", [".cache/gate"], network.env);
  assert.equal(result.status, 0, result.stderr);
  assert.equal(readFileSync(join(repo, ".cache", "gate", "args"), "utf8"), "");
  assert.deepEqual(network.calls(), []);
});

test("before the first release every package is new: -first, then the package of each module", (t) => {
  const repo = repository(t);
  addModule(repo, join(EXAMPLES, "smallest"), "beta");
  addModule(repo, join(EXAMPLES, "smallest"), "alpha");
  const network = stubNetwork(t, { tags: [] });
  const result = run(repo, "gate-inputs", [".cache/gate"], network.env);
  assert.equal(result.status, 0, result.stderr);
  const dir = join(repo, ".cache", "gate");
  assert.deepEqual(args(dir), ["-first", ".cache/gate/packages/alpha", ".cache/gate/packages/beta"]);
  for (const name of ["alpha", "beta"]) {
    const inventory = JSON.parse(readFileSync(join(dir, "packages", name, "inventory.json"), "utf8"));
    assert.equal(inventory.id, `novatrail/${name}`);
  }
});

test("a folder that is there already is refused, and nothing is written into it", (t) => {
  const repo = repository(t);
  addModule(repo, join(EXAMPLES, "smallest"), "alpha");
  mkdirSync(join(repo, "gate"));
  const result = run(repo, "gate-inputs", ["gate"], stubNetwork(t).env);
  assert.equal(result.status, 2);
  assert.match(result.stderr, /is there already/);
  assert.deepEqual(readdirSync(join(repo, "gate")), []);
});

test("published versions that cannot be read give the status 2, and no args", (t) => {
  const repo = repository(t);
  addModule(repo, join(EXAMPLES, "smallest"), "alpha");
  const result = run(repo, "gate-inputs", ["gate"], stubNetwork(t, { tags: null }).env);
  assert.equal(result.status, 2);
  assert.ok(!existsSync(join(repo, "gate", "args")));
});

// _gate-args in a made-up repository, against the published versions given, copied into
// <dir>/published/ as _published writes them.
function gateArgs(t, repo, from, env = {}) {
  cpSync(from, join(repo, "gate", "published"), { recursive: true });
  return run(repo, "_gate-args", ["gate"], env);
}

test("with a release: -index, a -published for each package unpacked, then the package of every module", (t) => {
  const repo = repository(t);
  addModule(repo, join(EXAMPLES, "smallest"), "alpha");
  addModule(repo, join(EXAMPLES, "smallest"), "beta");
  const release = published(t, [join(repo, "modules", "alpha"), join(repo, "modules", "beta")]);
  const result = gateArgs(t, repo, release);
  assert.equal(result.status, 0, result.stderr);
  const digests = readdirSync(join(release, "packages")).sort();
  assert.equal(digests.length, 2);
  assert.deepEqual(args(join(repo, "gate")), [
    "-index",
    "gate/published/module-index.json",
    ...digests.flatMap((digest) => ["-published", `gate/published/packages/${digest}`]),
    "gate/packages/alpha",
    "gate/packages/beta",
  ]);
});

test("a link or a name with a dot in modules/ is left out: modules-check refuses it", (t) => {
  const repo = repository(t);
  const alpha = addModule(repo, join(EXAMPLES, "smallest"), "alpha");
  symlinkSync(alpha, join(repo, "modules", "link"));
  addModule(repo, join(EXAMPLES, "smallest"), ".hidden", "novatrail/hidden");
  const result = gateArgs(t, repo, firstRelease(t));
  assert.equal(result.status, 0, result.stderr);
  assert.deepEqual(args(join(repo, "gate")), ["-first", "gate/packages/alpha"]);
});

test("published versions with both module-index.json and first, or neither, are refused", (t) => {
  for (const files of [["first", "module-index.json"], []]) {
    const repo = repository(t);
    addModule(repo, join(EXAMPLES, "smallest"), "alpha");
    const folder = join(repo, "gate", "published");
    mkdirSync(folder, { recursive: true });
    for (const file of files) writeFileSync(join(folder, file), "{}\n");
    const result = run(repo, "_gate-args", ["gate"]);
    assert.equal(result.status, 2);
    assert.match(result.stderr, /what the gate compares with is unsaid/);
    assert.ok(!existsSync(join(repo, "gate", "args")));
  }
});

test("a package module-kit cannot build gives the status 2, and no args", (t) => {
  const repo = repository(t);
  addModule(repo, join(EXAMPLES, "smallest"), "alpha");
  const kit = stubModuleKit(t);
  const result = gateArgs(t, repo, firstRelease(t), { ...kit.env, STUB_BUILD_STATUS: JSON.stringify({ alpha: 1 }) });
  assert.equal(result.status, 2);
  assert.match(result.stderr, /could not build modules\/alpha/);
  assert.ok(!existsSync(join(repo, "gate", "args")));
});
