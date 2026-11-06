// What the tests of the recipes share: the just that runs them, made-up repositories in temporary
// folders that hold a copy of the justfile and the Dockerfile, the modules of spec's format to fill them with, and a
// stand-in for module-kit that writes down how it was called.

import { spawnSync } from "node:child_process";
import { chmodSync, cpSync, mkdirSync, mkdtempSync, readFileSync, rmSync, writeFileSync } from "node:fs";
import { tmpdir } from "node:os";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

// The just of `just test`, which passes itself on, or the one on the PATH.
export const JUST = process.env.JUST || "just";

export const ROOT = join(dirname(fileURLToPath(import.meta.url)), "..");

// The format's modules at the tag of the spec submodule: each example breaks no rule, and each
// vector of a rule breaks that rule alone.
export const EXAMPLES = join(ROOT, "spec", "format", "examples", "module", "valid");
export const VECTORS = join(ROOT, "spec", "format", "vectors", "rules");

// A new folder that the test's after() removes.
export function temporaryFolder(t) {
  const folder = mkdtempSync(join(tmpdir(), "modules-test-"));
  t.after(() => rmSync(folder, { recursive: true, force: true }));
  return folder;
}

// The tag of module-kit that this repository's Dockerfile pins.
export const PINNED = /^ARG MODULE_KIT_VERSION=(.*)$/m.exec(
  readFileSync(join(ROOT, ".devcontainer", "Dockerfile"), "utf8"),
)[1];

// A made-up repository: a temporary folder with this repository's justfile and Dockerfile in it,
// so that a recipe run there reads that folder's modules/.
export function repository(t) {
  const folder = temporaryFolder(t);
  cpSync(join(ROOT, "justfile"), join(folder, "justfile"));
  mkdirSync(join(folder, ".devcontainer"));
  cpSync(join(ROOT, ".devcontainer", "Dockerfile"), join(folder, ".devcontainer", "Dockerfile"));
  return folder;
}

// A module of the format copied into modules/<name>/ of a repository, with its manifest's ID made
// novatrail/<name>, or the ID given.
export function addModule(repo, source, name, id = `novatrail/${name}`) {
  const folder = join(repo, "modules", name);
  cpSync(source, folder, { recursive: true });
  const manifest = join(folder, "manifest.yaml");
  const text = readFileSync(manifest, "utf8");
  if (!/^id: .*$/m.test(text)) throw new Error(`${source}: its manifest has no line id:`);
  writeFileSync(manifest, text.replace(/^id: .*$/m, `id: ${id}`));
  return folder;
}

// A recipe of the justfile in folder, run from there, with the environment given on top of this
// process's.
export function run(folder, recipe, args = [], env = {}) {
  return spawnSync(JUST, ["-f", join(folder, "justfile"), recipe, ...args], {
    cwd: folder,
    encoding: "utf8",
    env: { ...process.env, ...env },
  });
}

// A stand-in for module-kit, first on the PATH it returns. Each call but --version goes into the log
// as a line of JSON: its arguments, and for build whether the folder was there already. --version
// says the tag STUB_VERSION names, the pinned one unless given. check and build exit with the status
// that STUB_CHECK_STATUS and STUB_BUILD_STATUS name for the module's folder name, 0 otherwise; build
// writes an inventory.json that names novatrail/<folder name>, as a package of that module would,
// and no ID for a folder name STUB_NO_ID lists.
export function stubModuleKit(t) {
  const folder = temporaryFolder(t);
  const bin = join(folder, "bin");
  const log = join(folder, "calls.jsonl");
  mkdirSync(bin);
  const stub = join(bin, "module-kit");
  writeFileSync(
    stub,
    `#!/usr/bin/env node
const { appendFileSync, existsSync, mkdirSync, writeFileSync } = require("node:fs");
const { basename, join } = require("node:path");
const args = process.argv.slice(2);
if (args.length === 1 && args[0] === "--version") {
  process.stdout.write("module-kit " + process.env.STUB_VERSION + "\\nUnicode 17.0.0\\n");
  process.exit(0);
}
const call = { args };
if (args[0] === "build") call.existed = existsSync(args[2]);
appendFileSync(process.env.STUB_LOG, JSON.stringify(call) + "\\n");
const statuses = (name) => JSON.parse(process.env[name] || "{}")[basename(args[1])] ?? 0;
if (args[0] === "check") process.exit(statuses("STUB_CHECK_STATUS"));
if (args[0] === "build") {
  const status = statuses("STUB_BUILD_STATUS");
  if (status !== 0) process.exit(status);
  const noId = JSON.parse(process.env.STUB_NO_ID || "[]").includes(basename(args[1]));
  mkdirSync(args[2], { recursive: true });
  writeFileSync(join(args[2], "inventory.json"), JSON.stringify(noId ? {} : { id: "novatrail/" + basename(args[1]) }));
  process.exit(0);
}
process.exit(2);
`,
  );
  chmodSync(stub, 0o755);
  return {
    env: { PATH: `${bin}:${process.env.PATH}`, STUB_LOG: log, STUB_VERSION: PINNED },
    calls: () => {
      try {
        return readFileSync(log, "utf8").trim().split("\n").filter(Boolean).map((line) => JSON.parse(line));
      } catch (error) {
        if (error.code === "ENOENT") return [];
        throw error;
      }
    },
  };
}
