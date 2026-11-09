// What the tests of the recipes share: the just that runs them, made-up repositories in temporary
// folders that hold a copy of the justfile and the Dockerfile, the modules of spec's format to fill
// them with, published versions as _published writes them, a stand-in for module-kit that writes
// down how it was called, and spec for a made-up repository: a link to this repository's spec/, or a
// stand-in for its justfile.

import { spawnSync } from "node:child_process";
import { chmodSync, cpSync, mkdirSync, mkdtempSync, readFileSync, rmSync, symlinkSync, writeFileSync } from "node:fs";
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

// The spec of a made-up repository: a link to this repository's spec/, so that a recipe there runs
// `just -f spec/justfile …` as it does here, from the submodule at its tag.
export function withSpec(repo) {
  symlinkSync(join(ROOT, "spec"), join(repo, "spec"));
  return repo;
}

// A stand-in for spec's justfile in a made-up repository: its _tools exits with the status tools
// gives, and its sample-check writes each call's arguments into the log it returns, a line of JSON,
// and exits with the status given.
export function stubSpec(repo, status = 0, tools = 0) {
  const log = join(repo, "sample-check.jsonl");
  mkdirSync(join(repo, "spec"));
  writeFileSync(
    join(repo, "spec", "justfile"),
    `_tools:
    exit ${tools}

[positional-arguments]
sample-check *args:
    #!/usr/bin/env bash
    node -e 'require("node:fs").appendFileSync(process.argv[1], JSON.stringify(process.argv.slice(2)) + "\\n")' ${JSON.stringify(log)} "$@"
    exit ${status}
`,
  );
  return () => {
    try {
      return readFileSync(log, "utf8").trim().split("\n").filter(Boolean).map((line) => JSON.parse(line));
    } catch (error) {
      if (error.code === "ENOENT") return [];
      throw error;
    }
  };
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
// process's, and input on its stdin.
export function run(folder, recipe, args = [], env = {}, input = "") {
  return spawnSync(JUST, ["-f", join(folder, "justfile"), recipe, ...args], {
    cwd: folder,
    encoding: "utf8",
    env: { ...process.env, ...env },
    input,
  });
}

// What _modules-check compares with before the first release of modules: a folder with the file
// first, as _published writes it.
export function firstRelease(t) {
  const folder = temporaryFolder(t);
  writeFileSync(join(folder, "first"), "modules has no release yet, so every version is a first one\n");
  return folder;
}

// _modules-check of the repository, against the published versions given, or none yet.
export function check(t, repo, env = {}, published = firstRelease(t)) {
  return run(repo, "_modules-check", [published], env);
}

// The published versions of a made-up release, as _published writes them: the packages that
// module-kit build makes of the modules given, the index that module-kit index builds of them for
// the tag, and packages/sha256-<hex>/ with the three files of each package. A real index, but not a
// signed one: release-verify is no part of these tests.
export function published(t, modules, tag = "v0.1.0") {
  const folder = temporaryFolder(t);
  const built = [];
  mkdirSync(join(folder, "built"));
  for (const [i, module] of modules.entries()) {
    const out = join(folder, "built", String(i));
    mustRun("module-kit", ["build", module, out]);
    built.push(out);
  }
  mustRun("module-kit", ["index", "--first", "--tag", tag, "--no-withdrawals", "--out", join(folder, "module-index.json"), ...built]);
  const index = JSON.parse(readFileSync(join(folder, "module-index.json"), "utf8"));
  for (const out of built) {
    const inventory = JSON.parse(readFileSync(join(out, "inventory.json"), "utf8"));
    const entry = index.modules[inventory.id][inventory.version];
    const target = join(folder, "packages", entry.package_digest.replace(":", "-"));
    mkdirSync(target, { recursive: true });
    for (const file of ["inventory.json", "server.json", "device.json"]) cpSync(join(out, file), join(target, file));
  }
  rmSync(join(folder, "built"), { recursive: true });
  return folder;
}

// A command that must succeed, for the making of fixtures.
export function mustRun(command, args, options = {}) {
  const result = spawnSync(command, args, { encoding: "utf8", ...options });
  if (result.status !== 0) throw new Error(`${command} ${args.join(" ")}: ${result.stdout}${result.stderr}`);
  return result;
}

// A stand-in for module-kit, first on the PATH it returns. Each call but --version goes into the log
// as a line of JSON: its arguments, and for build whether the folder was there already. --version
// says the tag STUB_VERSION names, the pinned one unless given. check, build and diff exit with the
// status that STUB_CHECK_STATUS, STUB_BUILD_STATUS and STUB_DIFF_STATUS name for the module's folder
// name — the last argument's for diff —, 0 otherwise; index writes nothing and exits with
// STUB_INDEX_STATUS, 2 unless given; build
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
const statuses = (name, of = args[1]) => JSON.parse(process.env[name] || "{}")[basename(of)] ?? 0;
if (args[0] === "diff") process.exit(statuses("STUB_DIFF_STATUS", args[args.length - 1]));
if (args[0] === "index") process.exit(Number(process.env.STUB_INDEX_STATUS || 2));
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
