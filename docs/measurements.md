# modules — measurements

> The results of measurements: what was measured, how, on what, and the numbers. The documents that describe how things are built cite them — `measurements 2` — and hold only what was decided from them (RULES.md, The bar; R143 spec). The first measurements came here with T64 (R31). The external facts are in [sources.md](sources.md).

## 1. The tokens of the Claude Code sessions

The recipe is spec's: `just -f spec/justfile tokens-report --since … --until … --row` reads the transcripts that Claude Code keeps on this machine of this repository's sessions and prints one row of the table below, in the columns of measurements 1.9 spec; it prints numbers, paths and names alone, never what a transcript says. A row goes here at the archive of each epic, for the epic's window (RUN.md, «Как устроен план: эпики»), and once a month from spec, by R145 spec. A new analysis comes by the threshold of R145 spec, measured against this repository's reference.

| Project | Window, UTC | Sessions | Calls, main / subagents | Context of a call, k, main / subagents | Compactions | Cost at API prices |
|---|---|---|---|---|---|---|
| `-workspaces-modules` — before R143 spec, from measurements 1.9 spec | 2026-10-13T00:00 … 2026-11-12T14:27 | 16 | 3,976 / 2,287 | 351.0 / 91.7 | 6 | $586.03 |

- The row before is spec's measurement of 2026-11-12, of every session of `modules` on this machine until then: E1–E8, the moves to tags of spec and the debts.
- The reference of `modules`, against which the threshold of R145 spec is measured, is its measurement after, which T124 spec takes on the sessions after the task of row 7.9 of COORDINATION.md spec, T64–T67 (measurements 1.9 spec). T67 adds the row of the task's own sessions.

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
