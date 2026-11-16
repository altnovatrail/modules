# modules — measurements

> The results of measurements: what was measured, how, on what, and the numbers. The documents that describe how things are built cite them — `measurements 2` — and hold only what was decided from them (RULES.md, The bar; R143 spec). The first measurements came here with T64 (R31). The external facts are in [sources.md](sources.md).

## 1. The tokens of the Claude Code sessions

The recipe is spec's: `just -f spec/justfile tokens-report --since … --until … --row` reads the transcripts that Claude Code keeps on this machine of this repository's sessions and prints one row of the table below, in the columns of measurements 1.9 spec; it prints numbers, paths and names alone, never what a transcript says. A row goes here at the archive of each epic, for the epic's window (RUN.md, Epics), and once a month from spec, by R145 spec. A new analysis comes by the threshold of R145 spec, measured against this repository's reference.

| Project | Window, UTC | Sessions | Calls, main / subagents | Context of a call, k, main / subagents | Compactions | Cost at API prices |
|---|---|---|---|---|---|---|
| `-workspaces-modules` — before R143 spec, from measurements 1.9 spec | 2026-10-13T00:00 … 2026-11-12T14:27 | 16 | 3,976 / 2,287 | 351.0 / 91.7 | 6 | $586.03 |
| `-workspaces-modules` — after: the sessions of T64–T67 | 2026-11-12T20:29 … 2026-11-13T03:58 | 4 | 432 / 590 | 157.2 / 112.7 | 0 | $59.42 |
| `-workspaces-modules` — epic E9: its breakdown and T68–T71 | 2026-11-14T00:00 … 2026-11-14T19:41 | 5 | 326 / 215 | 187.4 / 120.6 | 0 | $34.22 |

- The row before is spec's measurement of 2026-11-12, of every session of `modules` on this machine until then: E1–E8, the moves to tags of spec and the debts.
- The reference of `modules`, against which the threshold of R145 spec is measured, is its measurement after, which T124 spec takes on the sessions after the task of row 7.9 of COORDINATION.md spec, T64–T67 (measurements 1.9 spec).
- The row after is the four sessions of the task of row 7.9, one for each of T64–T67, from the start of T64 to T67's own session, which was still running when the row was taken. T64 put in the hook of the session and turned off the connectors of claude.ai, and T65–T67 began with them (R31).
- The row of epic E9 is its five sessions, from the first day of the epic, by RUN.md, «Как устроен план: эпики»: the breakdown, T68, T69, T70 and T71, whose own session was still running when the row was taken. Its context of a main call is 19% above the reference, under the 25% of the threshold of R145 spec, and `docs-budgets.json` has not been raised since T66, so no task of analysis follows (RUN.md, T71).

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

- **What it showed.** Easy pages help `G0` and hinder `G2-half`, and the full ladder gives what the module gives now: the corridor is held not by the pool but by the model's estimate — Г4 of architecture 2.9 engine, the lag `M3.lag`. The author decided on 2026-11-11 to write no pages, and spec took the request back by its answer to В39 (Д170 spec; R31).
- **What it did not measure.** `olympiad-math` was not in it: `sim` v0.1.8 did not translate the check kind `time` (Д5, closed by R37), so `M4` and the cell `A-linked` were not read. The same measurement with both maths, on `sim` v0.1.9 or newer, was done on 2026-11-16, section 5, and Д8 is closed.

## 3. One call of `sample-check` on a large eval set, 2026-11-14

Whether `eval-check` can call `sample-check --eval-set` once on a module's eval set, with no chunks of 50, since T116 `spec` (Д7, closed by R32).

- **How.** A model of the case of 2026-11-10, when one call on 631 synthetic samples of `olympiad-math`, 322 MB, stopped Node with 134 (the T58 amendment to R25). In a temporary folder the 41 synthetic samples of `olympiad-math` were copied in turn into 631 files, 323 MB, each with an ID of its own by `jq --arg id … '.id = $id'`, indented as the originals. On it ran the recipe `eval-check` of T68 with `spec` v0.17.2, and `sample-check` of `spec` v0.16.3, from an archive of its tag, called once on the same files. Node v24.21.0, in the devcontainer.
- **The numbers:**

  | `sample-check` | Calls | Result | Time |
  |---|---|---|---|
  | `spec` v0.16.3 | one, 631 files | 134: "JavaScript heap out of memory" | 81 s, to the crash |
  | `spec` v0.17.2, in `eval-check` | one, 631 files | ok — 631 samples | 97 s, the whole recipe |

  On 2026-11-10, `eval-check` with chunks of 50 took 2 min 48 s on the 631 real synthetic samples (the T59 amendment to R25).
- **What it showed.** The model gives the crash of 2026-11-10 on the old tag, and T116 holds on the samples of `modules`, not only on the model of `spec`, 631 copies of an example essay (Д197 spec). One call is enough.
- **What it did not measure.** The synthetic samples of every task: the model repeats 41 of them. The heap limit of CI's runner either, which may differ from the devcontainer's.

## 4. The CI of `modules` before and after T74, 2026-11-16

What row 7.17 of COORDINATION.md spec gave the jobs it changes, "Go / test" and "Modules" (R36).

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

The measurement of section 2 again, on a tag of `sim` that translates the check kind `time`, so that `olympiad-math` runs beside `grade-3-math`, with `M4` and the cell `A-linked` (Д8, closed by T77; R30).

- **How.** The gate of `sim` v0.1.23, commit 61ed988, on `learner-model` v0.1.8 and `engine` v0.1.6, the commits of its `gate.json`. It was built from a checkout of that tag and SHA in a temporary folder, as the README of `sim` has the CI of `modules` take it, and no file of a release entered; it ran as `just gate modules -first <olympiad-math> <grade-3-math>`, 250 learners a cell. It ran on the package of `olympiad-math` 0.3.0 that `module-kit build` makes of `main` at 98429a0 and on trial copies of `grade-3-math` 0.3.0 in that folder, without `eval/`. A new page is a copy of the module's page of difficulty 1 (for d1 and d2), 3 (for d3 and d4) or 5 (for d5) of the same skill, with another ID and difficulty. No page was written into `modules/`. A run took 8 to 9 minutes on the 20 cores of the devcontainer.
- **What.** The share of the items in the corridor, `E1.corridor`, on `G0` and `G2-half`; the share of items too hard, `E1.hard`, and false "mastered", `M4.false.all`, on `G0`; the cells of `olympiad-math` and `A-linked`. The thresholds are those of the gate of `modules` (architecture 3.9 `sim`): `E1.corridor` at least 0.184, `E1.hard` at most 0.4954, `M4` at most 0.20. In `gate.json` of this tag the module as it is has `E1.corridor` 0.194 ± 0.019 on `G0` and 0.194 ± 0.009 on `G2-half`.
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

## 6. The CI of `modules` before and after T75, 2026-11-16

What `.github` v0.1.30 and the cache of Go of "Modules" across runs gave the jobs T75 changes (the T75 amendments to R35 and R36).

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
- **What it showed.** With a hit, the tests' step no longer compiles the port and its modules cold: 22 s to 13 s; the restore and `go mod verify` cost 7 s, so "Modules" gains about 6 s net. "Toolchain image" loads the same entry in all three, and its 13–30 s is the runner's disk. "Go / test" missed on the pull request, as `main` holds no entry of `go.yml` yet: its 271 s against 161 s is the spread of runners, not the cache, as T74's 170 s on pull request #110 shows (measurements 4).
- **What it did not measure.** The hit of `go.yml`'s caches and of "Modules" on `main`: the first push to `main` after the merge saves them, and the runs after it read them. One run each, not the medians of three of measurements 3.1 spec.
