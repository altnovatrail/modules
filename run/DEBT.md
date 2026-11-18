# DEBT.md — technical debt of `modules`

Table 4 of [RUN.md](../RUN.md): what was put off on purpose, its risk, when to return (R153 spec, R154 spec).

- Never read whole: take the rows the task and its sections name — `grep -n '^| Д9 ' run/DEBT.md` — and those `grep -n` finds by the paths and parts it changes (RUN.md rule 2).
- A row is what a task leaves open (RUN.md rule 9) and holds only what is open: drop a closed part, delete the row once closed. A closed row's text is in the archive of the epic it closed in.

## Table 4. Technical debt

Numbers are not reused: next — Д11.

| ID | Debt | From | Risk while open | When to return |
|---|---|---|---|---|
| Д10 | `docs-budgets.json` keeps its own limits for architecture 3.4, 35 KB, and 3.5, 27 KB, though both fit the 24 KB of the block `repository`; R31 records them | T88 | small: the two sections may grow past 24 KB unnoticed | a task that rewrites R31 whole (R160 spec): drop both limits with it |
