# DEBT.md — technical debt of `modules`

Table 4 of [RUN.md](../RUN.md): what was put off on purpose, its risk, when to return.

- Never read whole: take the rows the task and its sections name — `grep -n '^| Д9 ' run/DEBT.md` — and those `grep -n` finds by the paths and parts it changes (RUN.md rule 2).
- A row is what a task leaves open (RUN.md rule 9) and holds only what is open: drop a closed part, delete the row once closed. A closed row's text is in the archive of the epic it closed in.

## Table 4. Technical debt

Numbers are not reused: next — Д16.

| ID | Debt | From | Risk while open | When to return |
|---|---|---|---|---|
| Д12 | Links to the gate of `sim` name its README — R37, architecture 3.2, the comment of `test/gate-inputs.test.mjs` —, which `sim` moved to its `docs/README.md` (R21 spec); row 3.14 of COORDINATION.md spec | T90 | small: the link leads to a README without the call form | the next task that touches the gate, or the move to a new tag of `sim` |
| Д14 | The deny of `.claude/settings.json` that keeps sessions off a person's sample covers this checkout alone, not a clone `/workspaces/.epics/<epic>/modules` (R166 spec): `spec` took the author's answer to В41, option (a), as its T248 — the template's `.claude/settings.json` and spec's own deny for those clones —, which the template up to v0.19.17 does not carry yet | T93 (review) | a session of a task across repositories can read a person's sample in a clone by its absolute path; the files are public in the repository, and no session reads them by the rules alone | the move to the tag of `spec` whose template carries T248's deny: take the template's `.claude/settings.json` |
