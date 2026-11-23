# modules — measurements

> The results of measurements: what was measured, how, on what, and the numbers. The documents that describe how things are built cite them — `measurements 2` — and hold only what was decided from them. The external facts are in [sources.md](sources.md).

## 1. The tokens of the Claude Code sessions

The recipe is spec's: `just -f spec/justfile tokens-report --since … --until … --row` reads the transcripts that Claude Code keeps on this machine of this repository's sessions and prints one row of the table below, in the columns of measurements 1.9 spec; it prints numbers, paths and names alone, never what a transcript says. A row goes here when an epic goes to its archive, for the epic's window, and once a month from spec. A new analysis comes once a row passes spec's threshold, measured against this repository's reference.

| Project | Window, UTC | Sessions | Calls, main / subagents | Context of a call, k, main / subagents | Compactions | Cost at API prices | Context of a main call among its session's first 80, k |
|---|---|---|---|---|---|---|---|
| `-workspaces-modules` — before the documents were rebuilt to save tokens, from measurements 1.9 spec | 2026-10-13T00:00 … 2026-11-12T14:27 | 16 | 3,976 / 2,287 | 351.0 / 91.7 | 6 | $586.03 | — |
| `-workspaces-modules` — after: the four sessions of the rebuilding | 2026-11-12T20:29 … 2026-11-13T03:58 | 4 | 432 / 590 | 157.2 / 112.7 | 0 | $59.42 | — |
| `-workspaces-modules` — epic E9: its breakdown and its first four tasks | 2026-11-14T00:00 … 2026-11-14T19:41 | 5 | 326 / 215 | 187.4 / 120.6 | 0 | $34.22 | — |

- The last column, the context of a main call among the first 80 of its session, came with spec's recipe on 2026-11-18, and the rows before it have "—" there (measurements 1.9 spec).
- The row before is spec's measurement of 2026-11-12, of every session of `modules` on this machine until then: E1–E8, the moves to tags of spec and the debts.
- The reference of `modules`, against which spec's threshold is measured, is its row after, which spec takes on the sessions after the rebuilding (measurements 1.9 spec).
- The row after is the four sessions of the rebuilding, one for each of its tasks, from the start of the first to the last one's own session, which was still running when the row was taken. The first put in the hook of the session and turned off the connectors of claude.ai, and the other three began with them.
- The row of epic E9 is its five sessions, from the first day of the epic, as the plan's rules for epics have it: the breakdown and four tasks, the last of which was still running when the row was taken. Its context of a main call is 19% above the reference, under the 25% of spec's threshold, and `docs-budgets.json` had not been raised since the rebuilding, so no task of analysis follows.

## 2. The corridor of `engine` on `grade-3-math`, 2026-11-11

Whether more pages of `grade-3-math` would raise the share of the items in the corridor of `engine`, as Д170 spec asked of the next version of the module (Д8, closed by section 5).

- **How.** The gate of `sim` v0.1.8, commit d957fc7, built from a checkout of its tag in a temporary folder, `just gate modules -first <package folder>`, on trial copies of the module in that folder, where the new pages were copies of the pages of difficulty 1, 3 and 5 with another difficulty: from a page the simulator reads its difficulty, its skills and its boxes. No page was written into `modules/`.
- **What.** The share of the items in the corridor, `E1.corridor`, on the populations `G0` and `G2-half`.
- **The numbers:**

  | Pages | `G0` | `G2-half` |
  |---|---|---|
  | now: 12 pages of difficulty 1, 3 and 5 | 0.197 | 0.194 |
  | with d2 | 0.216 | 0.181 |
  | with d2 and d4 | 0.172 | 0.166 |
  | with a second d1 | 0.224 | 0.161 |
  | two pages for each difficulty from 1 to 5, 32 pages | 0.197 | 0.202 |

- **What it showed.** Easy pages help `G0` and hinder `G2-half`, and the full ladder gives what the module gives now: the corridor is held not by the pool but by the model's estimate, the lag `M3.lag` of `engine`. So no pages are written, and spec took its request back (Д170 spec).
- **What it did not measure.** `olympiad-math` was not in it: `sim` v0.1.8 did not translate the check kind `time` (Д5, closed), so `M4` and the cell `A-linked` were not read. The same measurement with both maths, on `sim` v0.1.9 or newer, was done on 2026-11-16, section 5, and Д8 is closed.

## 3. One call of `sample-check` on a large eval set, 2026-11-14

Whether `eval-check` can call `sample-check --eval-set` once on a module's eval set, with no chunks of 50, on `spec` v0.17.2 (Д7, closed).

- **How.** A model of the case of 2026-11-10, when one call on 631 synthetic samples of `olympiad-math`, 322 MB, stopped Node with 134. In a temporary folder the 41 synthetic samples of `olympiad-math` were copied in turn into 631 files, 323 MB, each with an ID of its own by `jq --arg id … '.id = $id'`, indented as the originals. On it ran the recipe `eval-check`, with one call, on `spec` v0.17.2, and `sample-check` of `spec` v0.16.3, from an archive of its tag, called once on the same files. Node v24.21.0, in the devcontainer.
- **The numbers:**

  | `sample-check` | Calls | Result | Time |
  |---|---|---|---|
  | `spec` v0.16.3 | one, 631 files | 134: "JavaScript heap out of memory" | 81 s, to the crash |
  | `spec` v0.17.2, in `eval-check` | one, 631 files | ok — 631 samples | 97 s, the whole recipe |

  On 2026-11-10, `eval-check` with chunks of 50 took 2 min 48 s on the 631 real synthetic samples.
- **What it showed.** The model gives the crash of 2026-11-10 on the old tag, and the new tag holds on the samples of `modules`, not only on the model of `spec`, 631 copies of an example essay (Д197 spec). One call is enough.
- **What it did not measure.** The synthetic samples of every task: the model repeats 41 of them. The heap limit of CI's runner either, which may differ from the devcontainer's.

## 4. The CI of `modules` before and after the parallel drawing, 2026-11-16

What drawing in parallel — the pictures in the tests, and `template-pictures` by a pool over the cores — and one cache of Go for the steps of "Modules" gave the jobs they change, "Go / test" and "Modules".

- **How.** The times of the jobs and of their steps from the GitHub API, `actions/runs/<id>/jobs`, read only: run 43808095168 of `main` at 84a0098, before, and run 43812154856 of pull request #110 at 338b57e, after, each one run on `ubuntu-24.04`. Locally, in the devcontainer — 20 cores, of them 4 fast and 16 slow, so a local time is only a direction —: a CPU profile of `drawTemplatePictures` on `olympiad-math` by `go test -cpuprofile`, and `go test -race ./tools/port` pinned to the 4 fast cores by `taskset -c 0-3`.
- **The numbers:**

  | Job or step | Before, s | After, s |
  |---|---|---|
  | "Go / test" | 342 | 170 |
  | — its step `ci-test` | 295 | 138 |
  | "Modules" | 152 | 119 |
  | — "Tests of the recipes and of the templates" | 27 | 33 |
  | — "The pictures of the lessons" | 10 | 0 |
  | — "The pictures of the instances of the templates" | 41 | 18 |
  | `go test -race ./tools/port`, locally on 4 cores | 379 | 214 |

  The profile before the pool: of 33.2 s of samples, `drawInstances` 32.1 s — `encodePNG` 17.2 s and `paintImage` 14.8 s —, and no text function among the nodes above 0.17 s. Under `-race`, the longest subtest is `TestKinds21Draw/11`, the stairs of 16 steps: 80 s alone, 6.7 s without `-race`, 95% of it in `vector.Rasterizer.Draw`, since each line, disc and rectangle of the canvas is composited over the whole image.
- **What it showed.** "Go / test" halves, and is no longer the longest path of the run: "module-kit signature" and "Modules" after it took about 185 s. The steps that run the port after the tests no longer compile it; the tests' step compiles it once, cold.
- **What it did not measure.** The whole run: in pull request #110 the jobs waited for runners up to about 3 minutes, as the orchestrator ran other repositories at the same time. One run each, not the medians of three of measurements 3.1 spec.

## 5. The corridor of `engine` on `grade-3-math`, with both maths, 2026-11-16

The measurement of section 2 again, on a tag of `sim` that translates the check kind `time`, so that `olympiad-math` runs beside `grade-3-math`, with `M4` and the cell `A-linked` (Д8, closed by it).

- **How.** The gate of `sim` v0.1.23, commit 61ed988, on `learner-model` v0.1.8 and `engine` v0.1.6, the commits of its `gate.json`. It was built from a checkout of that tag and SHA in a temporary folder, as the README of `sim` has the CI of `modules` take it, and no file of a release entered; it ran as `just gate modules -first <olympiad-math> <grade-3-math>`, 250 learners a cell. It ran on the package of `olympiad-math` 0.3.0 that `module-kit build` makes of `main` at 98429a0 and on trial copies of `grade-3-math` 0.3.0 in that folder, without `eval/`. A new page is a copy of the module's page of difficulty 1 (for d1 and d2), 3 (for d3 and d4) or 5 (for d5) of the same skill, with another ID and difficulty. No page was written into `modules/`. A run took 8 to 9 minutes on the 20 cores of the devcontainer.
- **What.** The share of the items in the corridor, `E1.corridor`, on `G0` and `G2-half`; the share of items too hard, `E1.hard`, and false "mastered", `M4.false.all`, on `G0`; the cells of `olympiad-math` and `A-linked`. The thresholds are those `sim` sets for the gate of `modules`: `E1.corridor` at least 0.184, `E1.hard` at most 0.4954, `M4` at most 0.20. In `gate.json` of this tag the module as it is has `E1.corridor` 0.194 ± 0.019 on `G0` and 0.194 ± 0.009 on `G2-half`.
- **The numbers:**

  | Pages of `grade-3-math` | `G0` corridor | `G2-half` corridor | `G0` too hard | `G0` `M4` | Gate |
  |---|---|---|---|---|---|
  | now: 12 pages of difficulty 1, 3 and 5 | 0.194 | 0.194 | 0.485 | 0.059 | green |
  | with d2, 16 | 0.210 | 0.182 | 0.444 | 0.045 | red: `G2-half` corridor |
  | with d2 and d4, 20 | 0.165 | 0.166 | 0.553 | 0.054 | red: both corridors, too hard |
  | with a second d1, 16 | 0.232 | 0.160 | 0.378 | 0.067 | red: `G2-half` corridor |
  | two pages for each difficulty from 1 to 5, 40 | 0.193 | 0.211 | 0.4957 | 0.029 | red: too hard, by 0.00035 |

  In every run the cells of `olympiad-math` and of `A-linked` were the same to the bit. `olympiad-math` had a corridor of 0.257 on `G0` and 0.265 on `G2-half`, and `M4` 0.023 on `G0`. On `A-linked` the gate of `modules` reads only its invariants, `E2.stuck` and `E2.nothing`, both 0 in every run.
- **What it showed.** The finding of section 2 holds with both maths. Easy pages help `G0` and hinder `G2-half`. A partial ladder lowers the corridor of the main population below the gate's threshold. The full ladder, more than three times the pages, leaves `G0` where it was and raises `G2-half` by 1.7 points, to 0.211: about two half-widths of `gate.json`. The corridor is held by the model's estimate, not by the pool, and the set of 12 pages stays. A new page of `grade-3-math` is also a matter for the gate: every set of pages but the module's own made the gate of `modules` red here (Д5).
- **What it did not measure.** `olympiad-math` with other pages, and `E1` on `A-linked`, which the gate does not read. The full ladder of section 2 had 32 pages and this one 40: section 2 does not say how its pages were spread, so the two ladders are compared each with its own "now", not with each other.

## 6. The CI of `modules` before and after the cache of Go across runs, 2026-11-16

What `.github` v0.1.30 and the cache of Go of "Modules" across runs gave the jobs they change.

- **How.** The times of the jobs and of their steps from the GitHub API, `actions/runs/<id>/jobs`, read only, each on `ubuntu-24.04`: run 43828584293 of `main` at d829b71, before; run 43855989803 of pull request #113, after — its first attempt at e228c8c, a miss of the cache of Go, and a second attempt of "Modules" alone, a hit. For the hit, the trial commit e228c8c let the pull request save the entry, which only a push to `main` saves; the next commit reverted it, and the entry was deleted (`gh cache delete`).
- **The numbers:**

  | Job or step | Before, `main`, s | After, miss, s | After, hit, s |
  |---|---|---|---|
  | "Modules" | 96 | 102 | 90 |
  | — "Toolchain image" | 27 | 30 | 13 |
  | — "Go caches from the cache" and "Go caches" | — | 0 | 7 |
  | — "Tests of the recipes and of the templates" | 22 | 22 | 13 |
  | — "The pictures of the instances of the templates" | 10 | 11 | 13 |
  | — "Go caches into the cache" | — | 3 | 0 |
  | "Go / test" | 271 | 161 | — |
  | — its step `ci-test` | 234 | 119 | — |
  | "Go / toolchain", new | — | 12 | — |

  The entry of "Modules" is 47.34 MiB. Before the change, `gh cache list` of `modules` on 2026-11-16 showed four entries of the key `toolchain-Linux-X64-c3e120a…`, 392.55–392.95 MiB each, saved from `main` and from pull requests, of 3.02 GB of the cache in all (`actions/cache/usage`). On the hit the log reads "all modules verified" and "Go caches: go-Linux-X64-modules-c3e120a…-7d612fa…, its modules verified."
- **What it showed.** With a hit, the tests' step no longer compiles the port and its modules cold: 22 s to 13 s; the restore and `go mod verify` cost 7 s, so "Modules" gains about 6 s net. "Toolchain image" loads the same entry in all three, and its 13–30 s is the runner's disk. "Go / test" missed on the pull request, as `main` holds no entry of `go.yml` yet: its 271 s against 161 s is the spread of runners, not the cache, as the 170 s of the parallel drawing on pull request #110 show (measurements 4).
- **What it did not measure.** The hit of `go.yml`'s caches and of "Modules" on `main`: the first push to `main` after the merge saves them, and the runs after it read them. One run each, not the medians of three of measurements 3.1 spec.

## 7. The CI of `modules` before and after the stage `toolchain-core`, 2026-11-16

What the stage `toolchain-core` and `.github` v0.1.31 gave the jobs they change.

- **How.** The times of the jobs and of their steps from the GitHub API, `actions/runs/<id>/jobs`, read only, each on `ubuntu-24.04`: before — run 43859371391 of CI and run 43860304299 of Release, of the push to `main` at 694480f, the merge of the cache of Go across runs, whose Dockerfile the cache held; after — run 43906415932 of CI and run 43906416263 of Release, of pull request #114 at 946fcee, whose new Dockerfile no entry held, so that every job built its stage. Which stage each job took, from its log: "The checks without a browser run in the stage toolchain-core." in "Documents", "ci-lint runs in the stage toolchain-core." in "Go / lint", `TOOLCHAIN_STAGE: toolchain-core` in the signatures and "Gate", the `RUN` of the generator of `canvas` in the build of "Modules" alone, and none in that of "Release".
- **The numbers:**

  | Job or step | Before, `main`, hit, s | After, pull request, miss, s |
  |---|---|---|
  | "module-kit signature" | 51 | 38 |
  | — its image | 32 | 22 |
  | "canvas signature" | 49 | 39 |
  | — its image | 31 | 19 |
  | "Gate" | 80 | 74 |
  | — its image | 30 | 22 |
  | "Modules", in `toolchain` | 100 | 89 |
  | — its image | 29 | 24 |
  | "Documents / docs-check-repo" | 41 | 42 |
  | "Go / lint", "Go / vuln", "Go / licenses" | 48, 61, 56 | 63, 67, 61 |
  | "Go / toolchain" | 10 | 88 |
  | "Go / test", in `toolchain` | 186 | 232 |
  | "Release / release" | 70 | 71 |
  | — "Work out the version" | 24 | 34 |
  | — "Build the files" | 5 | 5 |

- **What it showed.** The stages are where the decision log puts them, and every check passed in its stage. The time a job spends on its image is the runner's disk and network either way, 19–34 s: the core is the same downloads as `toolchain` but for one file of 412 KiB, the generator of `canvas` (422,240 bytes), so here the core saves no measurable time, unlike the browsers of the kind `typescript` (measurements 3.1 spec). "Go / toolchain" builds both stages on a miss and hands each on as an artifact, 88 s against 10 s of a hit.
- **What it did not measure.** The hit of both entries on `main`: the first push to `main` after the merge saves `toolchain-…` and `toolchain-core-…`, and the runs after it read them. One run each, not the medians of three of measurements 3.1 spec.
