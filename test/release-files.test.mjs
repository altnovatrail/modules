// The tests of release-files (R03, R17), the recipe that .github's shared release.yml runs to build
// a release's files, on made-up repositories with the format's modules and the real module-kit: the
// first release, a second one that keeps a version, a new version, a withdrawal, an edit without a
// new version, the archive's layout and bytes, and what the recipe refuses. Each archive goes back
// through _unpack and module-kit diff, as modules-check reads a published version (R11). Neither the
// signature nor the shared release itself is tried here: the release signs the files, and its trial
// on each pull request builds them from git archive in the toolchain stage (R17).

import assert from "node:assert/strict";
import { cpSync, existsSync, mkdirSync, readdirSync, readFileSync, utimesSync, writeFileSync } from "node:fs";
import { spawnSync } from "node:child_process";
import { createHash } from "node:crypto";
import { join } from "node:path";
import { test } from "node:test";

import { EXAMPLES, JUST, PINNED, ROOT, addModule, mustRun, repository, run, stubModuleKit, temporaryFolder } from "./helpers.mjs";

// The names a file of dist/ may have, as the shared release.yml checks them.
const SAFE = /^[A-Za-z0-9_+][A-Za-z0-9._+-]*$/;

// A repository with this repository's LICENSE at its root, as the release's tree has it, and the
// modules given, each name an example of the format.
function releasable(t, names = ["smallest", "with-media-and-notice"]) {
  const repo = repository(t);
  cpSync(join(ROOT, "LICENSE"), join(repo, "LICENSE"));
  const modules = {};
  for (const name of names) modules[name] = addModule(repo, join(EXAMPLES, name), name);
  return { repo, modules };
}

// release-files of the repository, which must pass, and what it put in dist/.
function release(repo, version, previous = "", env = {}) {
  const result = run(repo, "release-files", [version, previous], env);
  assert.equal(result.status, 0, result.stdout + result.stderr);
  return { result, files: readdirSync(join(repo, "dist")).sort(), index: readIndex(join(repo, "dist", "module-index.json")) };
}

function readIndex(file) {
  return JSON.parse(readFileSync(file, "utf8"));
}

// The index of a release, kept apart from dist/, which the next run makes afresh: what the shared
// release hands the next release, read-only.
function keep(t, repo) {
  const file = join(temporaryFolder(t), "module-index.json");
  cpSync(join(repo, "dist", "module-index.json"), file);
  return file;
}

function archiveOf(digest) {
  return `${digest.replace(":", "-")}.tar.gz`;
}

function edit(file, from, to) {
  const text = readFileSync(file, "utf8");
  assert.match(text, from);
  writeFileSync(file, text.replace(from, to));
}

function sha256(file) {
  return createHash("sha256").update(readFileSync(file)).digest("hex");
}

test("the first release: the index, LICENSE and an archive of each module, and nothing else", (t) => {
  const { repo } = releasable(t);
  const { files, index } = release(repo, "v0.1.0");
  const entries = Object.entries(index.modules).flatMap(([id, versions]) =>
    Object.entries(versions).map(([version, entry]) => ({ id, version, ...entry })),
  );
  assert.deepEqual(
    entries.map(({ id, version, release }) => [id, version, release]),
    [
      ["novatrail/smallest", "0.1.0", "v0.1.0"],
      ["novatrail/with-media-and-notice", "0.1.0", "v0.1.0"],
    ],
  );
  assert.deepEqual(files, ["LICENSE", "module-index.json", ...entries.map((entry) => archiveOf(entry.package_digest))].sort());
  for (const file of files) assert.match(file, SAFE);
  assert.deepEqual(readFileSync(join(repo, "dist", "LICENSE")), readFileSync(join(ROOT, "LICENSE")));
});

test("each archive is the package its entry names: _unpack takes it, and diff finds it as it was", (t) => {
  const { repo, modules } = releasable(t);
  const { index } = release(repo, "v0.1.0");
  const unpacked = temporaryFolder(t);
  const published = [];
  for (const versions of Object.values(index.modules)) {
    for (const entry of Object.values(versions)) {
      const folder = join(unpacked, entry.package_digest.replace(":", "-"));
      const result = run(ROOT, "_unpack", [join(repo, "dist", archiveOf(entry.package_digest)), folder]);
      assert.equal(result.status, 0, result.stderr);
      assert.equal(`sha256:${sha256(join(folder, "inventory.json"))}`, entry.package_digest);
      published.push("--published", folder);
    }
  }
  for (const module of Object.values(modules)) {
    const diff = mustRun("module-kit", ["diff", "--index", join(repo, "dist", "module-index.json"), ...published, module]);
    assert.match(diff.stdout, /the newest published version as it was/);
  }
});

test("an archive holds the package's files at its root, with no ./, no link, no owner and one time", (t) => {
  const { repo } = releasable(t, ["with-media-and-notice"]);
  const { index } = release(repo, "v0.1.0");
  const entry = index.modules["novatrail/with-media-and-notice"]["0.1.0"];
  const archive = join(repo, "dist", archiveOf(entry.package_digest));
  const names = mustRun("tar", ["-tzf", archive]).stdout.trim().split("\n");
  assert.deepEqual(names.slice(0, 4), ["inventory.json", "server.json", "device.json", "device/"]);
  assert.ok(names.includes("device/LICENSE") && names.includes("device/NOTICE"), names.join(" "));
  assert.ok(names.some((name) => name.startsWith("device/media/")), names.join(" "));
  assert.deepEqual([...names.slice(4)].sort(), names.slice(4), "the entries of device/ are sorted by name");
  for (const line of mustRun("tar", ["--numeric-owner", "--full-time", "-tvzf", archive]).stdout.trim().split("\n")) {
    assert.match(line, /^(-rw-r--r--|drwxr-xr-x) 0\/0 +\d+ 1970-01-01 00:00:00 /, line);
  }
  const header = readFileSync(archive).subarray(0, 10);
  assert.deepEqual([...header.subarray(4, 8)], [0, 0, 0, 0], "the gzip header holds no time");
  assert.equal(header[3] & 0x08, 0, "the gzip header holds no name");
});

test("the same tree gives the same bytes, whatever the times of its files and of the build", (t) => {
  const { repo, modules } = releasable(t);
  release(repo, "v0.1.0");
  const first = Object.fromEntries(readdirSync(join(repo, "dist")).map((file) => [file, sha256(join(repo, "dist", file))]));
  for (const module of Object.values(modules)) utimesSync(join(module, "manifest.yaml"), new Date(2001, 0, 1), new Date(2001, 0, 1));
  // build writes the package's files afresh, so a second later each has another time of its own.
  Atomics.wait(new Int32Array(new SharedArrayBuffer(4)), 0, 0, 1100);
  release(repo, "v0.1.0");
  const second = Object.fromEntries(readdirSync(join(repo, "dist")).map((file) => [file, sha256(join(repo, "dist", file))]));
  assert.deepEqual(second, first);
});

test("a second release keeps each version where it was published, and puts out no archive of it", (t) => {
  const { repo } = releasable(t);
  const first = release(repo, "v0.1.0");
  const { result, files, index } = release(repo, "v0.1.1", keep(t, repo));
  assert.deepEqual(index, first.index);
  assert.deepEqual(files, ["LICENSE", "module-index.json"]);
  assert.match(result.stdout, /^release-files: novatrail\/smallest 0\.1\.0 stays in v0\.1\.0, which published it$/m);
});

test("a new version gets an entry of the new release and an archive, and the old entry stays", (t) => {
  const { repo, modules } = releasable(t);
  const first = release(repo, "v0.1.0");
  edit(join(modules.smallest, "manifest.yaml"), /^title: .*$/m, "title: Another title");
  edit(join(modules.smallest, "manifest.yaml"), /^version: .*$/m, 'version: "0.1.1"');
  const { files, index } = release(repo, "v0.2.0", keep(t, repo));
  const smallest = index.modules["novatrail/smallest"];
  assert.deepEqual(smallest["0.1.0"], first.index.modules["novatrail/smallest"]["0.1.0"]);
  assert.equal(smallest["0.1.1"].release, "v0.2.0");
  assert.deepEqual(index.modules["novatrail/with-media-and-notice"], first.index.modules["novatrail/with-media-and-notice"]);
  assert.deepEqual(files, ["LICENSE", "module-index.json", archiveOf(smallest["0.1.1"].package_digest)].sort());
});

test("withdrawals.yaml withdraws a version, and a release that only withdraws holds the index and LICENSE", (t) => {
  const { repo } = releasable(t);
  release(repo, "v0.1.0");
  const previous = keep(t, repo);
  writeFileSync(join(repo, "withdrawals.yaml"), 'novatrail/smallest:\n  "0.1.0": wrong_key\n');
  const { files, index } = release(repo, "v0.1.1", previous);
  assert.equal(index.modules["novatrail/smallest"]["0.1.0"].withdrawn, "wrong_key");
  assert.equal(index.modules["novatrail/smallest"]["0.1.0"].release, "v0.1.0");
  assert.deepEqual(files, ["LICENSE", "module-index.json"]);
});

test("an edit without a new version is refused by module-index.one-package, and leaves no dist/", (t) => {
  const { repo, modules } = releasable(t);
  release(repo, "v0.1.0");
  const previous = keep(t, repo);
  edit(join(modules.smallest, "manifest.yaml"), /^title: .*$/m, "title: Another title");
  const result = run(repo, "release-files", ["v0.1.1", previous]);
  assert.equal(result.status, 1, result.stdout + result.stderr);
  assert.match(result.stdout, /^module-index\.one-package: novatrail\/smallest 0\.1\.0: /m);
  assert.equal(existsSync(join(repo, "dist")), false);
});

test("dist/ is made afresh: nothing of an earlier run stays in it", (t) => {
  const { repo } = releasable(t);
  mkdirSync(join(repo, "dist"));
  writeFileSync(join(repo, "dist", "stale.txt"), "from before\n");
  const { files } = release(repo, "v0.1.0");
  assert.equal(files.includes("stale.txt"), false);
});

test("a relative previous index is from where the recipe is called, not from the justfile", (t) => {
  const { repo } = releasable(t, ["smallest"]);
  release(repo, "v0.1.0");
  const here = temporaryFolder(t);
  cpSync(join(repo, "dist", "module-index.json"), join(here, "previous.json"));
  const result = spawnSync(JUST, ["-f", join(repo, "justfile"), "release-files", "v0.1.1", "previous.json"], {
    cwd: here,
    encoding: "utf8",
  });
  assert.equal(result.status, 0, result.stdout + result.stderr);
  assert.deepEqual(readdirSync(join(repo, "dist")).sort(), ["LICENSE", "module-index.json"]);
});

test("a refused call leaves no dist/ of an earlier run behind", (t) => {
  const { repo } = releasable(t, ["smallest"]);
  release(repo, "v0.1.0");
  const result = run(repo, "release-files", ["0.1.1", ""]);
  assert.equal(result.status, 2);
  assert.equal(existsSync(join(repo, "dist")), false);
});

test("a version that is no tag vX.Y.Z is refused before anything is built", (t) => {
  const { repo } = releasable(t);
  const stub = stubModuleKit(t);
  for (const version of ["0.1.0", "v0.1", "v01.1.0", "v0.1.0-rc.1", ""]) {
    const result = run(repo, "release-files", [version, ""], stub.env);
    assert.equal(result.status, 2, version);
    assert.match(result.stderr, /the version is a tag vX\.Y\.Z/);
  }
  assert.deepEqual(stub.calls(), []);
});

test("a previous index that is not there is refused: only an empty string says there is none", (t) => {
  const { repo } = releasable(t);
  const stub = stubModuleKit(t);
  const result = run(repo, "release-files", ["v0.1.1", join(repo, "no-such-index.json")], stub.env);
  assert.equal(result.status, 2);
  assert.match(result.stderr, /is no file that can be read/);
  assert.deepEqual(stub.calls(), []);
});

test("a tree with no module, or with anything but module folders in modules/, is refused", (t) => {
  const empty = repository(t);
  cpSync(join(ROOT, "LICENSE"), join(empty, "LICENSE"));
  const none = run(empty, "release-files", ["v0.1.0", ""]);
  assert.equal(none.status, 2);
  assert.match(none.stderr, /modules\/ holds no module/);
  const { repo } = releasable(t, ["smallest"]);
  writeFileSync(join(repo, "modules", "README.md"), "a file\n");
  const file = run(repo, "release-files", ["v0.1.0", ""]);
  assert.equal(file.status, 1);
  assert.match(file.stderr, /^modules\/README\.md: not a folder; modules\/ holds the folders of the modules alone \(R03\)$/m);
});

test("a module that check refuses, or whose ID is not novatrail/<folder>, is not released", (t) => {
  const repo = repository(t);
  cpSync(join(ROOT, "LICENSE"), join(repo, "LICENSE"));
  addModule(repo, join(EXAMPLES, "smallest"), "smallest", "novatrail/another");
  const result = run(repo, "release-files", ["v0.1.0", ""]);
  assert.equal(result.status, 1);
  assert.match(result.stderr, /the module's ID is novatrail\/another, and the module in modules\/smallest\/ is novatrail\/smallest/);
  assert.equal(existsSync(join(repo, "dist")), false);
});

test("a module-kit other than the one the Dockerfile pins builds nothing", (t) => {
  const { repo } = releasable(t, ["smallest"]);
  const stub = stubModuleKit(t);
  const result = run(repo, "release-files", ["v0.1.0", ""], { ...stub.env, STUB_VERSION: "v0.0.1" });
  assert.equal(result.status, 1);
  assert.match(result.stderr, new RegExp(`^modules-check: the module-kit on the PATH says "module-kit v0\\.0\\.1", and \\.devcontainer/Dockerfile pins module-kit ${PINNED.replaceAll(".", "\\.")}`, "m"));
  assert.deepEqual(stub.calls(), []);
});

test("index is called in the form of architecture 3.6 module-kit, and its refusal is the recipe's", (t) => {
  const { repo } = releasable(t, ["smallest"]);
  const stub = stubModuleKit(t);
  const previous = join(temporaryFolder(t), "module-index.json");
  writeFileSync(previous, "{}");
  writeFileSync(join(repo, "withdrawals.yaml"), 'novatrail/smallest:\n  "0.1.0": wrong_key\n');
  const later = run(repo, "release-files", ["v0.1.1", previous], { ...stub.env, STUB_INDEX_STATUS: "1" });
  assert.equal(later.status, 1);
  // _modules-check first — check, build and diff --first —, then the release's own build and index.
  const calls = stub.calls();
  assert.deepEqual(calls.map((call) => call.args[0]), ["check", "build", "diff", "build", "index"]);
  assert.equal(calls[2].args[1], "--first");
  const [build, index] = calls.slice(3);
  assert.deepEqual(build.args.slice(0, 2), ["build", "modules/smallest"]);
  assert.equal(build.existed, false);
  assert.deepEqual(index.args.slice(0, 7), ["index", "--previous", previous, "--tag", "v0.1.1", "--withdrawals", "withdrawals.yaml"]);
  assert.equal(index.args[7], "--out");
  assert.match(index.args[8], /\/dist\/module-index\.json$/);
  assert.deepEqual(index.args.slice(9), [build.args[2]]);
  assert.equal(existsSync(join(repo, "dist")), false);

  const first = stubModuleKit(t);
  mustRun("rm", [join(repo, "withdrawals.yaml")]);
  const result = run(repo, "release-files", ["v0.1.0", ""], { ...first.env, STUB_INDEX_STATUS: "2" });
  assert.equal(result.status, 2);
  assert.deepEqual(first.calls()[4].args.slice(0, 5), ["index", "--first", "--tag", "v0.1.0", "--no-withdrawals"]);
});

test("release-files takes two arguments, as the shared release.yml with previous-file gives them", (t) => {
  const { repo } = releasable(t, ["smallest"]);
  const one = run(repo, "release-files", ["v0.1.0"]);
  assert.notEqual(one.status, 0);
  assert.match(one.stderr, /got 1 positional argument but takes 2/);
});
