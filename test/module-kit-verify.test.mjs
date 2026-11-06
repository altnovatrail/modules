// The tests of the sums that module-kit-verify holds to the Dockerfile (R10): made-up SHA256SUMS and
// Dockerfiles, no network. The signature itself is not tried here: release-verify is spec's recipe,
// with tests of its own on files of real releases (R77 spec), and the job "module-kit signature" of
// ci.yml runs it on the real release, on every pull request and push to main.

import assert from "node:assert/strict";
import { readFileSync, writeFileSync } from "node:fs";
import { join } from "node:path";
import { test } from "node:test";

import { ROOT, repository, run, temporaryFolder } from "./helpers.mjs";

const AMD64 = "a".repeat(64);
const ARM64 = "b".repeat(64);

function dockerfile(lines) {
  return ["FROM scratch AS toolchain", ...lines, ""].join("\n");
}

const DOCKERFILE = dockerfile([
  "ARG MODULE_KIT_VERSION=v1.2.3",
  `ARG MODULE_KIT_SHA256_AMD64=${AMD64}`,
  `ARG MODULE_KIT_SHA256_ARM64=${ARM64}`,
]);

const SUMS = [
  `${"c".repeat(64)}  module-kit_v1.2.3_darwin_amd64`,
  `${AMD64}  module-kit_v1.2.3_linux_amd64`,
  `${ARM64}  module-kit_v1.2.3_linux_arm64`,
  `${"d".repeat(64)}  LICENSE`,
  "",
].join("\n");

// _module-kit-sums of this repository's justfile, on the two files written into a new folder.
function sums(t, sumsText, dockerfileText) {
  const folder = temporaryFolder(t);
  writeFileSync(join(folder, "SHA256SUMS"), sumsText);
  writeFileSync(join(folder, "Dockerfile"), dockerfileText);
  return run(ROOT, "_module-kit-sums", [join(folder, "SHA256SUMS"), join(folder, "Dockerfile")]);
}

test("the sums of both linux binaries, as SHA256SUMS lists them, pass", (t) => {
  const result = sums(t, SUMS, DOCKERFILE);
  assert.equal(result.status, 0, result.stderr);
  assert.match(result.stdout, new RegExp(`^module-kit v1\\.2\\.3: module-kit_v1\\.2\\.3_linux_amd64 ${AMD64}$`, "m"));
  assert.match(result.stdout, new RegExp(`^module-kit v1\\.2\\.3: module-kit_v1\\.2\\.3_linux_arm64 ${ARM64}$`, "m"));
});

test("another sum of either binary is refused", (t) => {
  for (const [arch, from] of [["amd64", AMD64], ["arm64", ARM64]]) {
    const result = sums(t, SUMS.replace(from, "e".repeat(64)), DOCKERFILE);
    assert.equal(result.status, 1, arch);
    assert.match(result.stderr, new RegExp(`module-kit_v1\\.2\\.3_linux_${arch} is e{64} in SHA256SUMS`));
  }
});

test("SHA256SUMS of another tag is refused", (t) => {
  const result = sums(t, SUMS.replaceAll("v1.2.3", "v1.2.2"), DOCKERFILE);
  assert.equal(result.status, 1);
  assert.match(result.stderr, /lists module-kit_v1\.2\.3_linux_amd64 not once but 0 times/);
  assert.match(result.stderr, /lists module-kit_v1\.2\.3_linux_arm64 not once but 0 times/);
});

test("a binary SHA256SUMS lists twice is refused", (t) => {
  const result = sums(t, SUMS + `${AMD64}  module-kit_v1.2.3_linux_amd64\n`, DOCKERFILE);
  assert.equal(result.status, 1);
  assert.match(result.stderr, /lists module-kit_v1\.2\.3_linux_amd64 not once but 2 times/);
});

test("a Dockerfile without the ARG, or with it twice, is refused", (t) => {
  const missing = sums(t, SUMS, dockerfile(["ARG MODULE_KIT_VERSION=v1.2.3", `ARG MODULE_KIT_SHA256_AMD64=${AMD64}`]));
  assert.notEqual(missing.status, 0);
  assert.match(missing.stderr, /sets ARG MODULE_KIT_SHA256_ARM64 not once but 0 times/);

  const twice = sums(t, SUMS, DOCKERFILE + "ARG MODULE_KIT_VERSION=v1.2.4\n");
  assert.notEqual(twice.status, 0);
  assert.match(twice.stderr, /sets ARG MODULE_KIT_VERSION not once but 2 times/);
});

test("a sum that is not 64 lowercase hexadecimal digits is refused", (t) => {
  const upper = AMD64.toUpperCase();
  const result = sums(t, SUMS.replace(AMD64, upper), DOCKERFILE.replace(AMD64, upper));
  assert.equal(result.status, 1);
  assert.match(result.stderr, /MODULE_KIT_SHA256_AMD64 is not a sha256 of 64 lowercase hexadecimal digits/);
});

test("this repository's Dockerfile pins module-kit by a tag vX.Y.Z and two sums", () => {
  const file = join(ROOT, ".devcontainer", "Dockerfile");
  const tag = run(ROOT, "_dockerfile-arg", [file, "MODULE_KIT_VERSION"]);
  assert.equal(tag.status, 0, tag.stderr);
  assert.match(tag.stdout.trim(), /^v\d+\.\d+\.\d+$/);
  for (const arch of ["AMD64", "ARM64"]) {
    const sum = run(ROOT, "_dockerfile-arg", [file, `MODULE_KIT_SHA256_${arch}`]);
    assert.equal(sum.status, 0, sum.stderr);
    assert.match(sum.stdout.trim(), /^[0-9a-f]{64}$/);
  }
});

test("a Dockerfile that is not there, or cannot be read, is refused by its name", (t) => {
  const folder = temporaryFolder(t);
  const missing = join(folder, "Dockerfile");
  const result = run(ROOT, "_dockerfile-arg", [missing, "MODULE_KIT_VERSION"]);
  assert.equal(result.status, 1);
  assert.match(result.stderr, new RegExp(`${missing.replaceAll(".", "\\.")} is no file that can be read`));
});

test("module-kit-verify refuses a MODULE_KIT_VERSION that is not a tag vX.Y.Z, before it reads the network", (t) => {
  for (const version of ["0.2.3", "v0.2", "v0.2.3/../../x", "main"]) {
    const repo = repository(t);
    const file = join(repo, ".devcontainer", "Dockerfile");
    writeFileSync(file, readFileSync(file, "utf8").replace(/^ARG MODULE_KIT_VERSION=.*$/m, `ARG MODULE_KIT_VERSION=${version}`));
    const result = run(repo, "module-kit-verify");
    assert.equal(result.status, 1, version);
    assert.match(result.stderr, new RegExp(`MODULE_KIT_VERSION is ${version.replaceAll(".", "\\.")}, not a tag vX\\.Y\\.Z`));
  }
});
