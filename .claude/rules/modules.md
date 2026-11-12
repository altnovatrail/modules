---
paths:
  - "modules/**"
---

# The modules: `modules/`

Each folder of `modules/` is one module in the module format, `spec/format/` at spec's tag (R42 spec, R45 spec); its ID is `novatrail/<folder>` (R03). The format's text is [spec/format/README.md](../../spec/format/README.md).

- `just modules-check` checks every module with module-kit's CLI of the toolchain stage (R10) and compares each with its published versions, from the newest release's module index and the archives it names (R11; the 2026-11-02 amendment to R03). A version that changes what a published one holds gets a new version number, as `module-kit diff` says.
- `just lesson-pictures <module>` draws the pictures of the module's lessons, from `tools/port/pictures/<the module's folder>/`, into its `media/lessons/` (the T41 amendment to R12), and `just template-pictures <module>` the picture of every instance of the module's templates, from `tools/port/template-pictures/<the module's folder>/`, into its `media/pictures/` (the T52 amendment to R12). CI runs both on `modules/olympiad-math` with `--check`.
- `just test` holds the recipes, and with them two tests of content: `test/templates.test.mjs` builds `olympiad-math` and checks each template against the task it came from (R22), and `test/lesson-pictures.test.mjs` checks the size of each picture of a lesson against its file (the T41 amendment to R12).
- A task that touches `modules/` ends with `just modules-check`, `just test` and the picture checks green.
- **A merge publishes a release.** A merge into `main` that changes anything in `modules/` outside `modules/*/eval/` — or `withdrawals.yaml`, `justfile`, `LICENSE` or `.github/workflows/release.yml` — puts out the next immutable release of `modules` after CI (R17): the pull request and the report name it, and after the merge `gh release view` and `release-verify` confirm it. A change of what the index carries — a new module version, a new format — waits for the author's word, and a new format may need a new release line (R28).
- A withdrawn version goes into `withdrawals.yaml` at the root, which grows only (R03); nothing is withdrawn yet, so there is no file.
- The eval sets, `modules/*/eval/`, have a rule of their own: [eval](eval.md).
