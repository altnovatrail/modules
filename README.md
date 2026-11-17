# modules

[![CI](https://github.com/altnovatrail/modules/actions/workflows/ci.yml/badge.svg)](https://github.com/altnovatrail/modules/actions/workflows/ci.yml)
[![CodeQL](https://github.com/altnovatrail/modules/actions/workflows/codeql.yml/badge.svg)](https://github.com/altnovatrail/modules/actions/workflows/codeql.yml)
[![Coverage](https://codecov.io/gh/altnovatrail/modules/branch/main/graph/badge.svg)](https://codecov.io/gh/altnovatrail/modules)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/altnovatrail/modules/badge)](https://scorecard.dev/viewer/?uri=github.com/altnovatrail/modules)
[![License](https://img.shields.io/github/license/altnovatrail/modules)](LICENSE)

[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=altnovatrail_modules&metric=code_smells)](https://sonarcloud.io/summary/new_code?id=altnovatrail_modules)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=altnovatrail_modules&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=altnovatrail_modules)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=altnovatrail_modules&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=altnovatrail_modules)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=altnovatrail_modules&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=altnovatrail_modules)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=altnovatrail_modules&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=altnovatrail_modules)

The official NovaTrail modules of the first slice: olympiad maths for grades 1–6, carried over from MathTrail, school maths for grade 3 on a workbook, and short essays in English for any age. Beside them is a small module in Russian for adults: the eleven tasks of tutor's check board, which the author solves in the app on a learner's path. A module is data, not code and not instructions to a model: a folder whose layout spec's module format sets, with its skill graph, lessons, tasks, traps and techniques in YAML, its media, its license's text, and an eval set: synthetic samples and samples that adult volunteers and authors made for eval, never learners' work. This repository's CI checks each module with the authors' CLI from module-kit, and its releases carry the signed module packages and the module index that the platform reads.

A component of [NovaTrail](https://github.com/altnovatrail/spec), a universal open platform for adaptive learning, interactive lessons and AI tutoring. It belongs to stage 1, Format and core, with spec, module-kit, learner-model, engine and sim; together with module-kit it tries the module format v0 on three modules of two unlike subjects at once.

## Documents

- [docs/architecture.md](docs/architecture.md) — the component: what it does and does not do, its data and links, and what is still to decide.
- [docs/decisions.md](docs/decisions.md) — the decision log: what was decided, why, and what was rejected.
- [docs/sources.md](docs/sources.md) — the external facts the documents rest on, each with its source, its date and its quote.
- [docs/measurements.md](docs/measurements.md) — the results of measurements: the tokens of the Claude Code sessions, and the corridor of `grade-3-math`.
- [RUN.md](RUN.md) — the work plan, in epics.
- `run/` — the plan's other tables: [the epics for later](run/EPICS.md), [the open questions](run/QUESTIONS.md), [the technical debt](run/DEBT.md) and [the tasks outside an epic](run/TASKS.md).

The working rules, for people and coding agents alike, are in [CLAUDE.md](CLAUDE.md) and [AGENTS.md](AGENTS.md). The platform's shared vocabulary, formats, contracts and decisions are in [spec](https://github.com/altnovatrail/spec), checked out here as a git submodule: clone this repository with `git clone --recurse-submodules`.
