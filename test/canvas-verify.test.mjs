// The tests of the sum that canvas-verify holds to the Dockerfile (R29): made-up SHA256SUMS and
// Dockerfiles, no network. The signature itself is not tried here: release-verify is spec's recipe,
// with tests of its own on files of real releases (R77 spec), and the job "canvas signature" of
// ci.yml runs it on the real release, on every pull request and push to main.

import assert from "node:assert/strict";
import { readFileSync, writeFileSync } from "node:fs";
import { join } from "node:path";
import { test } from "node:test";

import { ROOT, repository, run, temporaryFolder } from "./helpers.mjs";

const SUM = "a".repeat(64);

function dockerfile(lines) {
  return ["FROM scratch AS toolchain", ...lines, ""].join("\n");
}

const DOCKERFILE = dockerfile(["ARG CANVAS_VERSION=v1.2.3", `ARG CANVAS_SYNTHETIC_SHA256=${SUM}`]);

const SUMS = [
  `${"b".repeat(64)}  altnovatrail-canvas-1.2.3.tgz`,
  `${"c".repeat(64)}  samples-page_v1.2.3.html`,
  `${SUM}  synthetic-module_v1.2.3.mjs`,
  `${"d".repeat(64)}  LICENSE`,
  "",
].join("\n");

// _canvas-sum of this repository's justfile, on the two files written into a new folder.
function sum(t, sumsText, dockerfileText) {
  const folder = temporaryFolder(t);
  writeFileSync(join(folder, "SHA256SUMS"), sumsText);
  writeFileSync(join(folder, "Dockerfile"), dockerfileText);
  return run(ROOT, "_canvas-sum", [join(folder, "SHA256SUMS"), join(folder, "Dockerfile")]);
}

test("the sum of the generator, as SHA256SUMS lists it, passes", (t) => {
  const result = sum(t, SUMS, DOCKERFILE);
  assert.equal(result.status, 0, result.stderr);
  assert.equal(result.stdout, `canvas v1.2.3: synthetic-module_v1.2.3.mjs ${SUM}\n`);
});

test("another sum of the generator is refused", (t) => {
  const result = sum(t, SUMS.replace(`${SUM}  synthetic`, `${"e".repeat(64)}  synthetic`), DOCKERFILE);
  assert.equal(result.status, 1);
  assert.match(result.stderr, new RegExp(`synthetic-module_v1\\.2\\.3\\.mjs is e{64} in SHA256SUMS, and CANVAS_SYNTHETIC_SHA256 is a{64}`));
});

test("SHA256SUMS of another tag is refused", (t) => {
  const result = sum(t, SUMS.replaceAll("1.2.3", "1.2.2"), DOCKERFILE);
  assert.equal(result.status, 1);
  assert.match(result.stderr, /lists synthetic-module_v1\.2\.3\.mjs not once but 0 times/);
});

test("a generator SHA256SUMS lists twice is refused", (t) => {
  const result = sum(t, SUMS + `${SUM}  synthetic-module_v1.2.3.mjs\n`, DOCKERFILE);
  assert.equal(result.status, 1);
  assert.match(result.stderr, /lists synthetic-module_v1\.2\.3\.mjs not once but 2 times/);
});

test("a Dockerfile without the ARG, or with it twice, is refused", (t) => {
  const missing = sum(t, SUMS, dockerfile(["ARG CANVAS_VERSION=v1.2.3"]));
  assert.notEqual(missing.status, 0);
  assert.match(missing.stderr, /sets ARG CANVAS_SYNTHETIC_SHA256 not once but 0 times/);

  const twice = sum(t, SUMS, DOCKERFILE + "ARG CANVAS_VERSION=v1.2.4\n");
  assert.notEqual(twice.status, 0);
  assert.match(twice.stderr, /sets ARG CANVAS_VERSION not once but 2 times/);
});

test("a sum that is not 64 lowercase hexadecimal digits is refused", (t) => {
  const upper = SUM.toUpperCase();
  const result = sum(t, SUMS.replace(SUM, upper), DOCKERFILE.replace(SUM, upper));
  assert.equal(result.status, 1);
  assert.match(result.stderr, /CANVAS_SYNTHETIC_SHA256 is not a sha256 of 64 lowercase hexadecimal digits/);
});

test("this repository's Dockerfile pins canvas's generator by a tag vX.Y.Z and a sum", () => {
  const file = join(ROOT, ".devcontainer", "Dockerfile");
  const tag = run(ROOT, "_dockerfile-arg", [file, "CANVAS_VERSION"]);
  assert.equal(tag.status, 0, tag.stderr);
  assert.match(tag.stdout.trim(), /^v\d+\.\d+\.\d+$/);
  const pinned = run(ROOT, "_dockerfile-arg", [file, "CANVAS_SYNTHETIC_SHA256"]);
  assert.equal(pinned.status, 0, pinned.stderr);
  assert.match(pinned.stdout.trim(), /^[0-9a-f]{64}$/);
});

test("canvas-verify refuses a CANVAS_VERSION that is not a tag vX.Y.Z, before it reads the network", (t) => {
  for (const version of ["0.3.6", "v0.3", "v0.3.6/../../x", "main"]) {
    const repo = repository(t);
    const file = join(repo, ".devcontainer", "Dockerfile");
    writeFileSync(file, readFileSync(file, "utf8").replace(/^ARG CANVAS_VERSION=.*$/m, `ARG CANVAS_VERSION=${version}`));
    const result = run(repo, "canvas-verify");
    assert.equal(result.status, 1, version);
    assert.match(result.stderr, new RegExp(`CANVAS_VERSION is ${version.replaceAll(".", "\\.")}, not a tag vX\\.Y\\.Z`));
  }
});
