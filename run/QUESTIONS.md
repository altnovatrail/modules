# QUESTIONS.md — questions to the author of `modules`

Table 3 of [RUN.md](../RUN.md): the author's open questions that tasks leave (R153 spec, R154 spec, R160 spec).

- Never read whole: take the rows the task and its sections name — `grep -n '^| В1 ' run/QUESTIONS.md`.
- A row is a question a task leaves open (RUN.md rule 9). Record the author's answer where it belongs — the task, a document or [DEBT.md](DEBT.md) — and delete the row.
- A row whose question opens with the mark «В реестр spec» is for `spec`: its session finds it on `main` (R81 spec; RULES.md, Across repositories). It leaves once `spec` has recorded it: the task that rests on it names its number.

## Table 3. Questions to the author

Numbers are not reused: next — В42.

| ID | Question | From | Options | Recommendation |
|---|---|---|---|---|
| В41 | «В реестр spec» The deny of `.claude/settings.json` here — no session reads a file of `modules/*/eval/` but `synthetic-*.json` (the T61 amendment to R04) — is relative to the project, and every NovaTrail container now mounts `/workspaces/.epics` (R166 spec): a clone `.epics/<epic>/modules` is read by its absolute path from this container, and from spec's orchestrator container, where the task across repositories runs, with no deny at all. Where does the deny for those clones live? | T90 (review) | (a) the template's `.claude/settings.json` of `spec/repo-template/files/` and spec's own deny `Read(//workspaces/.epics/*/modules/modules/*/eval/*)` with the exception `Read(!//workspaces/.epics/*/modules/modules/*/eval/synthetic-*.json)`, and every repository takes it with its next tag; (b) `xrepo` clones `modules` sparse, without `modules/*/eval/` but `synthetic-*.json`; (c) nothing: the files are public in the repository, and the rule holds by the sessions' word alone | (a): one rule in every container, as the deny here is, and no clone of `modules` a task needs loses its files; (b) breaks `just eval-check` in the clone |
