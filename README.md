# modules

[![CI](https://github.com/altnovatrail/modules/actions/workflows/ci.yml/badge.svg)](https://github.com/altnovatrail/modules/actions/workflows/ci.yml)
[![CodeQL](https://github.com/altnovatrail/modules/actions/workflows/codeql.yml/badge.svg)](https://github.com/altnovatrail/modules/actions/workflows/codeql.yml)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/altnovatrail/modules/badge)](https://scorecard.dev/viewer/?uri=github.com/altnovatrail/modules)
[![License](https://img.shields.io/github/license/altnovatrail/modules)](LICENSE)

[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=altnovatrail_modules&metric=code_smells)](https://sonarcloud.io/summary/new_code?id=altnovatrail_modules)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=altnovatrail_modules&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=altnovatrail_modules)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=altnovatrail_modules&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=altnovatrail_modules)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=altnovatrail_modules&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=altnovatrail_modules)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=altnovatrail_modules&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=altnovatrail_modules)

The official NovaTrail modules of the first slice: olympiad maths for grades 1–6, carried over from MathTrail, school maths for grade 3 on a workbook, and short essays in English for any age. A module is data, not code and not instructions to a model: a folder whose layout spec's module format sets, with its skill graph, lessons, tasks, traps and techniques in YAML, its media, its license's text, and an eval set: synthetic samples and samples that adult volunteers and authors made for eval, never learners' work. This repository's CI checks each module with the authors' CLI from module-kit, and its releases carry the signed module packages and the module index that the platform reads.

A component of [NovaTrail](https://github.com/altnovatrail/spec/tree/v0.6.3), a universal open platform for adaptive learning, interactive lessons and AI tutoring. It belongs to stage 1, Format and core, with spec, module-kit, learner-model, engine and sim; together with module-kit it tries the module format v0 on three modules of two unlike subjects at once.

## Status

Detailing. The documents come first: the first epic of [RUN.md](RUN.md) decides what the component's card left open, and the code follows in the order RUN.md sets. There are no modules here yet. The module format v0 is in spec, at version 0.2, and the first content — a first batch of olympiad maths for the first application, where the author solves tasks in a browser — comes with the epic after this one, once module-kit's CLI can check it.

## Documents

The documents are in Russian.

- [docs/architecture.md](docs/architecture.md) — the component: what it does and does not do, its data and links, and what is still to decide.
- [docs/decisions.md](docs/decisions.md) — the decision log: what was decided, why, and what was rejected.
- [RUN.md](RUN.md) — the work plan, in epics.

The working rules, for people and coding agents alike, are in [CLAUDE.md](CLAUDE.md) and [AGENTS.md](AGENTS.md). The platform's shared vocabulary, formats, contracts and decisions are in [spec](https://github.com/altnovatrail/spec/tree/v0.6.3), checked out here as a git submodule at `v0.6.3`: clone this repository with `git clone --recurse-submodules`.

## License

[Apache-2.0](LICENSE), for the modules' content and code alike. Each module's manifest names its license by the SPDX identifier `Apache-2.0`, and each module's folder carries the license's text as its own `LICENSE`, since the module's package takes it to every reader.
