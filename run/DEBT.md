# DEBT.md — technical debt of `modules`

Table 4 of [RUN.md](../RUN.md): what was put off on purpose, its risk, when to return.

- Never read whole: take the rows the task and its sections name — `grep -n '^| Д9 ' run/DEBT.md` — and those `grep -n` finds by the paths and parts it changes (RUN.md rule 2).
- A row is what a task leaves open (RUN.md rule 9) and holds only what is open: drop a closed part, delete the row once closed. A closed row's text is in the archive of the epic it closed in.

## Table 4. Technical debt

Numbers are not reused: next — Д16.

| ID | Debt | From | Risk while open | When to return |
|---|---|---|---|---|
| Д12 | Links to the gate of `sim` name its README — R37, architecture 3.2, the comment of `test/gate-inputs.test.mjs` —, which `sim` moved to its `docs/README.md` (R21 spec); row 3.14 of COORDINATION.md spec | T90 | small: the link leads to a README without the call form | the next task that touches the gate, or the move to a new tag of `sim` |
