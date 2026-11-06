// The tests of what modules-check compares with (R03, R11), without the network: _release-tags on
// the lines of git ls-remote, _entries on indexes, _unpack on archives made here, and _modules-check against releases made
// here by the real module-kit — the same version with the same digest, the same version with another,
// a new version, no package of a version the index names, a module the index does not name. The
// signature is not tried: release-verify is spec's recipe, with tests of its own on files of real
// releases (R77 spec), and _published runs it on the real index in CI. Neither is git ls-remote nor
// the download: until the first release of modules, CI takes the path with no release, and the
// first pull request after it takes the one with an index (RUN.md, T10, «Риски»).

import assert from "node:assert/strict";
import { chmodSync, existsSync, mkdirSync, readdirSync, readFileSync, symlinkSync, writeFileSync } from "node:fs";
import { join } from "node:path";
import { test } from "node:test";

import {
  EXAMPLES,
  ROOT,
  addModule,
  check,
  firstRelease,
  mustRun,
  published,
  repository,
  run,
  stubModuleKit,
  temporaryFolder,
} from "./helpers.mjs";

function releaseTags(lines) {
  return spawnTags(lines.map((tag, i) => `${String(i).padStart(40, "0")}\trefs/tags/${tag}`).join("\n") + "\n");
}

function spawnTags(input) {
  return run(ROOT, "_release-tags", [], {}, input);
}

test("the releases are the tags vX.Y.Z, newest first by the numbers of their parts", () => {
  const result = releaseTags(["v0.9.9", "v0.10.0", "v0.2.30", "v1.0.0", "v0.10.1"]);
  assert.equal(result.status, 0, result.stderr);
  assert.equal(result.stdout, "v1.0.0\nv0.10.1\nv0.10.0\nv0.9.9\nv0.2.30\n");
});

test("a tag that is not vX.Y.Z is no release of modules", () => {
  const result = releaseTags(["v1.0.0-rc.1", "v01.0.0", "v2", "contracts/v9.0.0", "v0.1.0"]);
  assert.equal(result.status, 0, result.stderr);
  assert.equal(result.stdout, "v0.1.0\n");
});

test("no tag at all is no release yet: nothing is printed", () => {
  for (const input of ["", "abc\trefs/tags/latest\n"]) {
    const result = spawnTags(input);
    assert.equal(result.status, 0, result.stderr);
    assert.equal(result.stdout, "");
  }
});

const DIGEST = `sha256:${"a".repeat(64)}`;
const OTHER = `sha256:${"b".repeat(64)}`;

function entries(t, index, id = "novatrail/a") {
  const file = join(temporaryFolder(t), "module-index.json");
  writeFileSync(file, JSON.stringify(index));
  return run(ROOT, "_entries", [file, id]);
}

test("_entries lists a module's versions, withdrawn ones too, with their release and digest", (t) => {
  const result = entries(t, {
    module_format_version: "0.3",
    modules: {
      "novatrail/a": {
        "0.1.0": { module_format_version: "0.3", package_digest: DIGEST, release: "v0.1.0", withdrawn: "wrong_key" },
        "0.2.0": { module_format_version: "0.3", package_digest: OTHER, release: "v0.2.0" },
      },
      "novatrail/b": { "9.0.0": { module_format_version: "0.3", package_digest: DIGEST, release: "v0.3.0" } },
    },
  });
  assert.equal(result.status, 0, result.stderr);
  assert.equal(result.stdout, `0.1.0 v0.1.0 ${DIGEST}\n0.2.0 v0.2.0 ${OTHER}\n`);
});

test("_entries of a module the index does not name is nothing", (t) => {
  const result = entries(t, { module_format_version: "0.3", modules: {} });
  assert.equal(result.status, 0, result.stderr);
  assert.equal(result.stdout, "");
});

test("_entries refuses an entry that gives no address, and an index with no object modules", (t) => {
  for (const entry of [
    { package_digest: DIGEST, release: "main" },
    { package_digest: "sha256:../../x", release: "v0.1.0" },
    { release: "v0.1.0" },
  ]) {
    const result = entries(t, { modules: { "novatrail/a": { "0.1.0": entry } } });
    assert.equal(result.status, 2, JSON.stringify(entry));
    assert.match(result.stderr, /which give no address/);
  }
  const result = entries(t, { modules: [] });
  assert.equal(result.status, 2);
  assert.match(result.stderr, /holds no object modules/);
});

// An archive made by tar in a folder of its own: entries is a function that lays out the folder,
// and members are what tar is given, in order, so that a test can make any layout of names.
function archive(t, entries, members, options = []) {
  const folder = temporaryFolder(t);
  const source = join(folder, "source");
  mkdirSync(source);
  entries(source);
  const file = join(folder, "package.tar.gz");
  mustRun("tar", ["-czf", file, ...options, "-C", source, ...members]);
  return file;
}

function packageFiles(source) {
  for (const file of ["inventory.json", "server.json", "device.json"]) writeFileSync(join(source, file), `{"${file}":1}`);
  mkdirSync(join(source, "device", "media"), { recursive: true });
  writeFileSync(join(source, "device", "LICENSE"), "Apache-2.0\n");
  writeFileSync(join(source, "device", "media", "row.png"), "png");
}

function unpack(t, file) {
  const out = join(temporaryFolder(t), "package");
  return { out, result: run(ROOT, "_unpack", [file, out]) };
}

test("_unpack takes the three files of a package and nothing else", (t) => {
  const file = archive(t, packageFiles, ["inventory.json", "server.json", "device.json", "device"]);
  const { out, result } = unpack(t, file);
  assert.equal(result.status, 0, result.stderr);
  assert.deepEqual(readdirSync(out).sort(), ["device.json", "inventory.json", "server.json"]);
  assert.equal(readFileSync(join(out, "server.json"), "utf8"), '{"server.json":1}');
});

test("_unpack takes the files with ./ before their names too", (t) => {
  const file = archive(t, packageFiles, ["."]);
  const { out, result } = unpack(t, file);
  assert.equal(result.status, 0, result.stderr);
  assert.deepEqual(readdirSync(out).sort(), ["device.json", "inventory.json", "server.json"]);
});

test("_unpack refuses an archive with a link anywhere in it, and takes nothing out", (t) => {
  for (const [name, make] of [
    ["inventory.json", (source) => symlinkSync("/etc/passwd", join(source, "inventory.json"))],
    ["device/media/row.png", (source) => symlinkSync("../../server.json", join(source, "device", "media", "row.png"))],
  ]) {
    const file = archive(
      t,
      (source) => {
        packageFiles(source);
        const target = join(source, name);
        if (existsSync(target)) mustRun("rm", [target]);
        make(source);
      },
      ["inventory.json", "server.json", "device.json", "device"],
    );
    const { out, result } = unpack(t, file);
    assert.equal(result.status, 2, name);
    assert.match(result.stderr, /a link or a special file/);
    assert.equal(existsSync(out), false);
  }
});

test("_unpack refuses a hard link", (t) => {
  const file = archive(
    t,
    (source) => {
      packageFiles(source);
      mustRun("ln", ["-f", join(source, "server.json"), join(source, "device.json")]);
    },
    ["inventory.json", "server.json", "device.json"],
  );
  const { result } = unpack(t, file);
  assert.equal(result.status, 2);
  assert.match(result.stderr, /device\.json, a link or a special file/);
});

test("_unpack refuses a path outside its folder, absolute or through ..", (t) => {
  for (const transform of ["s|^server.json$|../server.json|", "s|^server.json$|/tmp/server.json|", "s|^server.json$|a/../../server.json|"]) {
    const file = archive(t, packageFiles, ["inventory.json", "server.json", "device.json"], ["--absolute-names", `--transform=${transform}`]);
    const { out, result } = unpack(t, file);
    assert.equal(result.status, 2, transform);
    assert.match(result.stderr, /a path outside its folder/);
    assert.equal(existsSync(out), false);
  }
});

test("_unpack refuses an archive without one of the three files, or with one twice", (t) => {
  const missing = archive(t, packageFiles, ["inventory.json", "device.json"]);
  const first = unpack(t, missing);
  assert.equal(first.result.status, 2);
  assert.match(first.result.stderr, /holds no server\.json at its root/);

  // tar keeps a file named again as a hard link, unless told to store it whole.
  const twice = archive(t, packageFiles, ["inventory.json", "server.json", "device.json", "./server.json"], ["--hard-dereference"]);
  const second = unpack(t, twice);
  assert.equal(second.result.status, 2);
  assert.match(second.result.stderr, /holds server\.json twice/);
});

test("_unpack says which file an empty archive lacks", (t) => {
  const file = archive(t, () => {}, ["."], ["--exclude=."]);
  const { result } = unpack(t, file);
  assert.equal(result.status, 2);
  assert.match(result.stderr, /holds no inventory\.json at its root/);
});

test("_unpack refuses what is no archive, and a folder that is there already", (t) => {
  const folder = temporaryFolder(t);
  writeFileSync(join(folder, "not.tar.gz"), "plain text");
  const { result } = unpack(t, join(folder, "not.tar.gz"));
  assert.equal(result.status, 2);
  assert.match(result.stderr, /no gzip tar archive/);

  const file = archive(t, packageFiles, ["inventory.json", "server.json", "device.json"]);
  const there = run(ROOT, "_unpack", [file, folder]);
  assert.equal(there.status, 2);
  assert.match(there.stderr, /is there already/);
});

// A repository with the module smallest as novatrail/smallest, and a release of it as it is now.
function released(t) {
  const repo = repository(t);
  const module = addModule(repo, join(EXAMPLES, "smallest"), "smallest");
  return { repo, module, release: published(t, [module]) };
}

function edit(file, from, to) {
  const text = readFileSync(file, "utf8");
  assert.match(text, from);
  writeFileSync(file, text.replace(from, to));
}

test("the same version with the same digest passes: nothing has changed", (t) => {
  const { repo, release } = released(t);
  const result = check(t, repo, {}, release);
  assert.equal(result.status, 0, result.stdout + result.stderr);
  assert.match(result.stdout, /novatrail\/smallest 0\.1\.0: the newest published version as it was/);
});

test("the same version with another digest is refused by module-version.above", (t) => {
  const { repo, module, release } = released(t);
  edit(join(module, "manifest.yaml"), /^title: .*$/m, "title: Another title");
  const result = check(t, repo, {}, release);
  assert.equal(result.status, 1);
  assert.match(result.stdout, /^refused: module-version\.above: 0\.1\.0 is published already/m);
});

test("a new version above the published one passes, and diff reports what changed", (t) => {
  const { repo, module, release } = released(t);
  edit(join(module, "manifest.yaml"), /^title: .*$/m, "title: Another title");
  edit(join(module, "manifest.yaml"), /^version: .*$/m, 'version: "0.1.1"');
  const result = check(t, repo, {}, release);
  assert.equal(result.status, 0, result.stdout + result.stderr);
  assert.match(result.stdout, /^novatrail\/smallest 0\.1\.1, compared with 0\.1\.0$/m);
  assert.match(result.stdout, /^server\.json #\/manifest\/title: patch: changes$/m);
});

test("no package of a version the index names: diff cannot compare, and its status 2 is the recipe's", (t) => {
  const { repo, release } = released(t);
  mustRun("rm", ["-r", join(release, "packages")]);
  const result = check(t, repo, {}, release);
  assert.equal(result.status, 2);
  assert.match(result.stderr, /no --published folder holds it/);
});

test("a module the index does not name is a first version", (t) => {
  const { repo } = released(t);
  const other = repository(t);
  const release = published(t, [addModule(other, join(EXAMPLES, "with-tasks"), "with-tasks")]);
  const result = check(t, repo, {}, release);
  assert.equal(result.status, 0, result.stdout + result.stderr);
  assert.match(result.stdout, /novatrail\/smallest 0\.1\.0: the first version/);
});

test("with no release yet, --first: the first version, with nothing to compare it with", (t) => {
  const { repo } = released(t);
  const result = check(t, repo, {}, firstRelease(t));
  assert.equal(result.status, 0, result.stdout + result.stderr);
  assert.match(result.stdout, /novatrail\/smallest 0\.1\.0: the first version, with nothing to compare it with/);
});

test("diff gets --index and the module's packages alone, in the form of architecture 3.5 module-kit", (t) => {
  const repo = repository(t);
  addModule(repo, join(EXAMPLES, "smallest"), "smallest");
  const other = repository(t);
  const both = published(t, [join(repo, "modules", "smallest"), addModule(other, join(EXAMPLES, "with-tasks"), "with-tasks")]);
  const stub = stubModuleKit(t);
  const result = check(t, repo, stub.env, both);
  assert.equal(result.status, 0, result.stderr);
  const diff = stub.calls().find((call) => call.args[0] === "diff");
  const index = JSON.parse(readFileSync(join(both, "module-index.json"), "utf8"));
  const digest = index.modules["novatrail/smallest"]["0.1.0"].package_digest;
  assert.deepEqual(diff.args.slice(0, 5), [
    "diff",
    "--index",
    join(both, "module-index.json"),
    "--published",
    join(both, "packages", digest.replace(":", "-")),
  ]);
  assert.equal(diff.args.length, 6, "the package of another module is not given");
});

test("the recipe's status is the highest of its modules', diff's among them, and no refusal stops the rest", (t) => {
  const repo = repository(t);
  for (const name of ["a", "b", "c"]) addModule(repo, join(EXAMPLES, "smallest"), name);
  const stub = stubModuleKit(t);
  const result = check(t, repo, { ...stub.env, STUB_DIFF_STATUS: JSON.stringify({ a: 2, b: 1 }) });
  assert.equal(result.status, 2);
  assert.deepEqual(
    stub.calls().filter((call) => call.args[0] === "diff").map((call) => call.args.at(-1).split("/").at(-1)),
    ["a", "b", "c"],
  );
  const refused = check(t, repo, { ...stub.env, STUB_DIFF_STATUS: JSON.stringify({ c: 1 }) });
  assert.equal(refused.status, 1);
});

test("without modules/, modules-check asks nothing of the network", (t) => {
  const repo = repository(t);
  const bin = join(temporaryFolder(t), "bin");
  mkdirSync(bin);
  for (const tool of ["git", "curl"]) {
    writeFileSync(join(bin, tool), `#!/bin/sh\necho "${tool} must not run" >&2\nexit 99\n`);
    chmodSync(join(bin, tool), 0o755);
  }
  const result = run(repo, "modules-check", [], { PATH: `${bin}:${process.env.PATH}` });
  assert.equal(result.status, 0, result.stderr);
  assert.match(result.stdout, /^modules\/: no modules yet$/m);
});

test("what to compare with must be said: neither or both of module-index.json and first are refused", (t) => {
  const { repo, release } = released(t);
  const neither = temporaryFolder(t);
  const none = check(t, repo, {}, neither);
  assert.equal(none.status, 2);
  assert.match(none.stderr, /what to compare with is unsaid/);

  writeFileSync(join(release, "first"), "");
  const both = check(t, repo, {}, release);
  assert.equal(both.status, 2);
  assert.match(both.stderr, /what to compare with is unsaid/);
});
