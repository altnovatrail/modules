---
paths:
  - "modules/*/eval/**"
---

# The eval sets: `modules/*/eval/`

A module's eval set is its folder `eval/`: synthetic samples and samples that adult volunteers and authors made, a file `eval/<id>.json` each, by the ID of the sample, with nothing named with a dot in it (R04, R25). No learner's work ever enters it.

- **No session opens a person's sample** (the T61 amendment to R04): `.claude/settings.json` denies reading any file of `modules/*/eval/` but `synthetic-*.json`. A person's sample is placed, moved or removed by its name alone, never read, and a script that would print one is not run.
- `just eval-check` checks the layout of each module's `eval/` and each module's eval set with spec's `sample-check --eval-set`, one call per module, which reads the files one at a time (R25 and its amendment R32).
- `just eval-match` holds each sample to the package of the commit: a synthetic one drawn again by canvas's generator, a person's by the digest of its task's open view, which `module-kit open` gives (the T59 amendment to R25). A pull request that changes a task's open view remakes or removes its samples in the same pull request.
- `just synthetic <module> <task>:<seed>...` writes synthetic samples into a module's `eval/` with canvas's generator of the toolchain stage, on the package `module-kit build` makes of the module, with the tag of the newest release of `modules` (R29). It refuses with "rebuild the container" when the container's `module-kit` or generator is not the one the Dockerfile pins.
- A task that touches `eval/` ends with `just eval-check` and `just eval-match` green. `eval/` is no path of a release: a change there publishes none (R17).
- How samples arrive — the volunteer's form, the organizer, the upload link, consent withdrawal and the author's own path — is [docs/architecture.md](../../docs/architecture.md), section 3.4.
