# DEBT.md — technical debt of `modules`

Table 4 of [RUN.md](../RUN.md): what was put off on purpose, its risk, when to return (R153 spec, R154 spec).

- Never read whole: take the rows the task and its sections name — `grep -n '^| Д9 ' run/DEBT.md` — and those `grep -n` finds by the paths and parts it changes (RUN.md rule 2).
- A row comes from a task's «Открыто» (RUN.md rule 9) and holds only what is open: drop a closed part, delete the row once closed. A closed row's text is in the archive of the epic it closed in.
- "Д5 in table 4 of RUN.md" or "in DEBT.md" in older texts is the same row (R31; R34).

## Table 4. Technical debt

Numbers are not reused: next — Д10.

| ID | Debt | From | Risk while open | When to return |
|---|---|---|---|---|
| Д9 | Д152 spec: the reminder "Do not write your name or the real name of anyone you know." in the statements of five `english-essays` tasks leaves in a new module version once `app` shows it itself: once the release of `app` with T47 `app`, which shows it above "Check" of every handwriting sheet whose work goes to AI whole, is deployed — not merged (R127 spec, its 2026-11-17 amendment; R39) | move to v0.15.0 (R30); move to v0.18.14 (R39) | None: the reminder stays in the statement until `app` shows it | Once that release of `app` is deployed — a new version of `english-essays`; the row leaves with it |
