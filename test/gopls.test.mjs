// The lock file of gopls in .devcontainer/gopls/ and GOPLS_VERSION of the Dockerfile are those of
// the go kind of spec's repo-template at the tag the submodule pins (R34). CI never builds the
// devcontainer stage, so nothing else would see them drift apart before a Rebuild Container.

import assert from "node:assert/strict";
import { readFileSync } from "node:fs";
import { join } from "node:path";
import { test } from "node:test";

import { ROOT } from "./helpers.mjs";

const KIND = join(ROOT, "spec", "repo-template", "kinds", "go");

function goplsVersion(dockerfile) {
  const match = /^ARG GOPLS_VERSION=(.*)$/m.exec(readFileSync(dockerfile, "utf8"));
  assert.ok(match, `${dockerfile} has no ARG GOPLS_VERSION`);
  return match[1];
}

test("the lock file of gopls is the template's, byte for byte", () => {
  for (const file of ["go.mod", "go.sum"]) {
    assert.deepEqual(
      readFileSync(join(ROOT, ".devcontainer", "gopls", file)),
      readFileSync(join(KIND, "gopls", file)),
      `.devcontainer/gopls/${file} differs from spec/repo-template/kinds/go/gopls/${file}`,
    );
  }
});

test("GOPLS_VERSION is the template's, and the lock file requires that gopls", () => {
  const version = goplsVersion(join(ROOT, ".devcontainer", "Dockerfile"));
  assert.equal(version, goplsVersion(join(KIND, "Dockerfile")));
  const lock = readFileSync(join(ROOT, ".devcontainer", "gopls", "go.mod"), "utf8");
  assert.match(lock, new RegExp(`^\\s*golang\\.org/x/tools/gopls v${version.replaceAll(".", "\\.")}\\b`, "m"));
});
