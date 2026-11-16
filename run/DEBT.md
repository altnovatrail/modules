# DEBT.md — technical debt of `modules`

Table 4 of [RUN.md](../RUN.md): what was put off on purpose, its risk, when to return (R153 spec, R154 spec).

- Never read whole: take the rows the task and its sections name — `grep -n '^| Д5 ' run/DEBT.md` — and those `grep -n` finds by the paths and parts it changes (RUN.md rule 2).
- A row comes from a task's «Открыто» (RUN.md rule 9) and holds only what is open: drop a closed part, delete the row once closed. A closed row's text is in the archive of the epic it closed in.
- "Д5 in table 4 of RUN.md" or "in DEBT.md" in older texts is the same row (R31; R34).

## Table 4. Technical debt

Numbers are not reused: next — Д10.

| ID | Debt | From | Risk while open | When to return |
|---|---|---|---|---|
| Д5 | No `sim` gate in CI: every release from the first, `v0.1.0`, went out without it (R76 spec). Call form, `sim`'s README: `just gate modules (-index <module-index.json> \| -first) [-published <folder>]... [-all] <package folder>...`, from a checkout of `sim` at tag and SHA, in the `toolchain` stage of its devcontainer. `sim` v0.1.8 (d957fc7, 2026-11-11) did not translate check kind `time` and refused `olympiad-math` 0.3.0 whole (30 tasks with `time`; architecture 3.9 `sim`); without that module the gate is red too: `M4` on `G0` needs both maths. Author 2026-11-11: wait for a tag that translates `time`, no always-red job. `sim` v0.1.9 (2026-11-11; row 4.10 of [COORDINATION.md](../spec/COORDINATION.md) spec) translates it, v0.1.10 (2026-11-12) is newer: the job waits only for its task (R31). Scope — R103 spec, Т-16 closed: all the stand reads from a package, on the pull request's package, invariants and absolute thresholds, no bands; a pull request that changes the `sim` tag runs all modules; a format move runs without the stand (Д111 spec; 2026-11-07 amendment to R03). A move changes nothing the stand reads: that goes in a pull request of its own, with the stand (2026-11-08 amendment to R03) | T04 | A graph or tasks that break adaptation go out in a release without the stand. `sim`'s first tag runs the published packages (R76 spec), but nobody starts the CI change without this row | A task of its own, on `sim` v0.1.10 or newer: a CI job by section 3.2 of [docs/architecture.md](../docs/architecture.md) and the call form of `sim`'s README at that tag. Recipe `modules-check` builds the pull request's packages; recipe `_published` gives the index and the published packages — `-index` and `-published`; `-all` when the pull request changes the `sim` tag. Leaves once the job is merged |
| Д9 | Д152 spec: the reminder "Do not write your name or the real name of anyone you know." in the statements of five `english-essays` tasks leaves in a new module version once `app` shows it itself (R127 spec) | move to v0.15.0 (R30) | None: the reminder stays in the statement until `app` shows it | Once `app` shows the reminder itself — a new version of `english-essays`; the row leaves with it |
