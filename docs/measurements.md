# modules — measurements

> The results of measurements: what was measured, how, on what, and the numbers. The documents that describe how things are built cite them — `measurements 2` — and hold only what was decided from them (RULES.md, The bar; R143 spec). The first measurements came here with T64 (R31). The external facts are in [sources.md](sources.md).

## 1. The tokens of the Claude Code sessions

The recipe is spec's: `just -f spec/justfile tokens-report --since … --until … --row` reads the transcripts that Claude Code keeps on this machine of this repository's sessions and prints one row of the table below, in the columns of measurements 1.9 spec; it prints numbers, paths and names alone, never what a transcript says. A row goes here at the archive of each epic, for the epic's window (RUN.md, «Как устроен план: эпики»), and once a month from spec, by R145 spec. A new analysis comes by the threshold of R145 spec, measured against this repository's reference.

| Project | Window, UTC | Sessions | Calls, main / subagents | Context of a call, k, main / subagents | Compactions | Cost at API prices |
|---|---|---|---|---|---|---|
| `-workspaces-modules` — before R143 spec, from measurements 1.9 spec | 2026-10-13T00:00 … 2026-11-12T14:27 | 16 | 3,976 / 2,287 | 351.0 / 91.7 | 6 | $586.03 |
| `-workspaces-modules` — after: the sessions of T64–T67 | 2026-11-12T20:29 … 2026-11-13T03:58 | 4 | 432 / 590 | 157.2 / 112.7 | 0 | $59.42 |

- The row before is spec's measurement of 2026-11-12, of every session of `modules` on this machine until then: E1–E8, the moves to tags of spec and the debts.
- The reference of `modules`, against which the threshold of R145 spec is measured, is its measurement after, which T124 spec takes on the sessions after the task of row 7.9 of COORDINATION.md spec, T64–T67 (measurements 1.9 spec).
- The row after is the four sessions of the task of row 7.9, one for each of T64–T67, from the start of T64 to T67's own session, which was still running when the row was taken. T64 put in the hook of the session and turned off the connectors of claude.ai, and T65–T67 began with them (R31).

## 2. The corridor of `engine` on `grade-3-math`, 2026-11-11

Whether more pages of `grade-3-math` would raise the share of the items in the corridor of `engine`, as Д170 spec asked of the next version of the module (Д8 in [DEBT.md](../DEBT.md)).

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
- **What it did not measure.** `olympiad-math` was not in it: `sim` v0.1.8 did not translate the check kind `time` (Д5 in [DEBT.md](../DEBT.md)), so `M4` and the cell `A-linked` were not read. The same measurement with both maths, on `sim` v0.1.9 or newer, is what Д8 still waits for.

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
