# DEBT.md — technical debt of `modules`

Table 4 of [RUN.md](../RUN.md): what was put off on purpose, its risk, when to return (R153 spec, R154 spec).

- Never read whole: take the rows the task and its sections name — `grep -n '^| Д9 ' run/DEBT.md` — and those `grep -n` finds by the paths and parts it changes (RUN.md rule 2).
- A row is what a task leaves open (RUN.md rule 9) and holds only what is open: drop a closed part, delete the row once closed. A closed row's text is in the archive of the epic it closed in.

## Table 4. Technical debt

Numbers are not reused: next — Д13.

| ID | Debt | From | Risk while open | When to return |
|---|---|---|---|---|
| Д10 | `docs-budgets.json` keeps its own limits for architecture 3.4, 35 KB, and 3.5, 27 KB, though both fit the 24 KB of the block `repository`; R31 records them | T88 | small: the two sections may grow past 24 KB unnoticed | a task that rewrites R31 whole (R160 spec): drop both limits with it |
| Д11 | The log cites amendments that `spec` v0.19.4 folded into their entries (R160 spec): the 2026-11-14 and 2026-11-15 amendments to R76 spec in R03, R32 and R36, the 2026-11-08 amendment to R82 spec in R07, R12 and R27, the 2026-11-16 amendment to R156 spec in R39; and, already at v0.18.34, the 2026-11-08 amendment to R32 spec in R23, the 2026-11-15 amendment to R143 spec in R34 and the 2026-11-17 amendment to R145 spec in R41. Each entry but R41 has dated amendments, so a fix rewrites it whole | T90 | small: a session that looks for the amendment in spec's entry finds none; the entry holds its text | T91, or a task that rewrites the entry whole (R160 spec) before it, as T83 does R36 and R41: cite the entry without the amendment |
| Д12 | Links to the gate of `sim` name its README — R37, architecture 3.2, the comment of `test/gate-inputs.test.mjs` —, which `sim` moved to its `docs/README.md` (R21 spec); row 3.14 of COORDINATION.md spec | T90 | small: the link leads to a README without the call form | the next task that touches the gate, or the move to a new tag of `sim` |
