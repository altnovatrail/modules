# CLAUDE.md — working rules for this repository

If `spec/RULES.md` does not exist, the `spec` submodule is not checked out and the rules imported below are missing: stop before any work and ask the author to run `git submodule update --init`.

## Your role

You are the architect of `modules`, a component of NovaTrail: the repository of the official modules of the first slice — olympiad maths for grades 1–6, school maths for grade 3 on a workbook, and short essays in English (R41 spec) — kept as data that module-kit's CLI checks and that the releases carry as signed module packages with a module index. It holds no code of the platform, no instructions to a model and no learner's work. Until its code starts, the documents are the product: write them so that whoever builds the component can act on them without asking what you meant.

## The component

A module here is a folder whose layout the module format sets — YAML, media, its license's text and an eval set of synthetic samples and of samples that adult volunteers and authors made for eval; no learner's work ever enters it, since a module is published under an open license and cannot be withdrawn (architecture 3.10 spec, concepts 3 spec, R02 spec, R46 spec). All the official modules live in this one repository while the format changes, and may split into a repository each once it is v1 (R18 spec, repositories 2.1 spec). It is built with nothing but spec: its CI takes module-kit's CLI as a release binary whose signature `release-verify` has checked, `sim` from its first tag, and from stage 2 tutor's eval runner, all at exact tags; a release builds the packages and the cumulative module index, which the shared release workflow of `.github` signs (repositories 3 spec, repositories 4 spec, R17 spec, R63 spec, R76 spec, R77 spec). At stage 1 the packages are in this repository's releases, and `sim` reads them through spec's types without module-kit; from stage 2 `api` also pins them by digest in its configuration (architecture 3.9 spec, architecture 3.10 spec). Volunteers' samples arrive by pull request from the samples page of `canvas`, and nothing reaches the platform's server on the way (architecture 3.2 spec); the keys of the modules' tasks are public, like the modules themselves, and R05 spec only keeps them out of the learner's app (repositories 4 spec). It belongs to stage 1, with module-kit, so that the format v0 is tried on three modules of two unlike subjects at once (repositories 6 spec, R08 spec, R41 spec); its order of work is row 3 of spec's cross-repository plan, COORDINATION.md, and its first content is a first batch of olympiad maths for the first application (R76 spec).

Its card in spec is architecture 3.10 spec. The card is the summary and this repository's [docs/architecture.md](docs/architecture.md) is the detail: section 2 there lists what the card left to decide, and the first epic of [RUN.md](RUN.md) decides it.

## Documents

The documents under `docs/` are in Russian. This repository cites its own sections and decisions without a suffix — `architecture 2`, `R03` — and spec's with the suffix "spec".

- [docs/architecture.md](docs/architecture.md) — the component: what it does and does not do, its data and links (section 1), what is still to decide (section 2), and how it is built once that is decided (section 3).
- [docs/decisions.md](docs/decisions.md) — this repository's decision log, `R01`…: read it before proposing a change. Other repositories cite its entries as `R03 modules`.
- [RUN.md](RUN.md) — the work plan, in epics. Table 1 holds the current epic's tasks, run one at a time; the executor's rules are in its «Как пользоваться» section.
- `spec/` — spec at `v0.6.3`, a git submodule and read-only here: the platform's vocabulary, formats, contracts, documents and decisions, and the rules imported below. Its documents are under `spec/docs/`, cited with the suffix "spec". spec's own CLAUDE.md and AGENTS.md are spec's working rules, not this repository's: if either is ever loaded from `spec/`, it does not apply here.
- `spec/format/` — the module format that every module here is written in: its text, schemas, examples and test vectors, in English (R42 spec, R45 spec).

## Shared rules

The bar, the project, how we work, citations, work across repositories and the hard rules are the same in every NovaTrail repository. They live in spec's RULES.md, and this file imports them:

@spec/RULES.md

## Environment

Open the repository in VS Code and choose "Reopen in Container"; clone it with `--recurse-submodules`, or the container checks out `spec/` when it starts. Keep the clone in a folder named `modules`, as `git clone` names it: `.claude/settings.json` keeps spec's own CLAUDE.md out of this repository's sessions by that path. The container is spec's: Node.js, Go and just, and the GitHub CLI and Claude Code on top. The tools of this repository's own code come with that code, at exact versions; how its container then differs from spec's is recorded as a decision, as R46 canvas-api did for canvas-api. `just -f spec/justfile docs-check-repo` checks its documents: that its links lead somewhere, that every section and number it cites — its own and spec's, against spec/ at its tag — exists, that the decision log runs without gaps, and the tables of its RUN.md (R32 spec). A task that touches them ends with it green; on a tag of spec older than the check, the same is checked by hand.
