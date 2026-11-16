# EPICS.md — epics for later of `modules`

Table 2 of [RUN.md](../RUN.md) (R153 spec, R154 spec).

- A row has a result and a basis, no tasks yet. Row order is the order epics are taken, until the author decides otherwise. What each builds and waits for — [docs/architecture.md](../docs/architecture.md), section 3.7.
- Read by the choice and planning of the next epic, not by each task. Once an epic is split into tasks, they go to table 1 and the body of RUN.md, and its row leaves (RUN.md, Epics).

## Table 2. Epics for later

| ID | Epic | Result | Basis |
|---|---|---|---|
| E6 | Volunteers' samples and the eval runner | A volunteer's sample by the form — organizer, the form's record outside the repository, the file's channel — and `tutor`'s eval runner in CI with a platform key of its own, by section 3.4 of [docs/architecture.md](../docs/architecture.md); the runner's reports CI reads are checked by `report-check` (second 2026-11-08 amendment to R95 spec; R27); decided what CI prints of a report and at which providers, by their terms ([contracts text](../spec/contracts/README.md), section 14; R26) | R04; R24 spec, R95 spec, R110 spec; П-16, П-18. Eval keys in CI — environment `eval` and a job by R120 spec (Д149 spec; R28). Also the samples of templates, handwritten choice and essays that T62 of E8 did not give ([docs/plan/E8.md](../docs/plan/E8.md)), and those of the board module that T70 of E9 did not give ([docs/plan/E9.md](../docs/plan/E9.md)). Waits for stage 2 with `tutor`'s eval runner released — no `tutor-eval` binaries in `tutor`'s releases yet (architecture 3.2 tutor, read 2026-11-11) — and for volunteers. Form, channel, periods and withdrawal are in since `spec` v0.15.0 (R30); "open file" mode — T87 `canvas`; the words "give the file to the eval-set organizer" — T20 `canvas` |
