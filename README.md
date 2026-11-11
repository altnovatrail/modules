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

The official NovaTrail modules of the first slice: olympiad maths for grades 1–6, carried over from MathTrail, school maths for grade 3 on a workbook, and short essays in English for any age. A module is data, not code and not instructions to a model: a folder whose layout spec's module format sets, with its skill graph, lessons, tasks, traps and techniques in YAML, its media, its license's text, and an eval set: synthetic samples and samples that adult volunteers and authors made for eval, never learners' work. This repository's CI checks each module with the authors' CLI from module-kit, and its releases carry the signed module packages and the module index that the platform reads.

A component of [NovaTrail](https://github.com/altnovatrail/spec/tree/v0.15.0), a universal open platform for adaptive learning, interactive lessons and AI tutoring. It belongs to stage 1, Format and core, with spec, module-kit, learner-model, engine and sim; together with module-kit it tries the module format v0 on three modules of two unlike subjects at once.

## Status

Detailing is done: the first epic of [RUN.md](RUN.md) decided what the component's card left open. The second, the first content, is done, and so is the third, the rest of the port of olympiad maths. `novatrail/olympiad-math` 0.3.0 holds all 753 tasks of mathtrail-standalone in 19 skills on three bands, 22 templates with solvers and 41 lessons. Beside it is a small set of school maths for grade 3 on a workbook, `novatrail/grade-3-math` 0.3.0: twelve pages of four skills and their four lessons, each with a worked example on the workbook. And a small set of short essays in English for any age, `novatrail/english-essays` 0.2.0: eight tasks of six skills on a sheet of lines, three choices of a first sentence and five essays with a rubric, and their six lessons. The module format v0 is in spec, at version 0.4, which all three modules are written in, and CI checks every module in `modules/` with module-kit's CLI, a binary of its release whose signature CI verifies. A release builds the packages of the three modules and the cumulative module index, which the shared release workflow of `.github` signs: `release.yml` runs it after CI on `main`, and tries it on every pull request without publishing; the first release, `v0.1.0`, is out, the second, `v0.1.1`, carries `olympiad-math` 0.1.1, and `v0.2.0`, the first of the line 0.2, carries the three modules 0.2.0 in format 0.4, and `v0.2.1` carries `olympiad-math` 0.3.0 and `grade-3-math` 0.3.0. The fourth epic, eval sets, is done too, and so is the fifth, the eval-set mode in CI: CI holds each module's `eval/` to its layout, `eval/<id>.json`, and checks its eval set with spec's `sample-check --eval-set`; the eighth epic, the first samples, is done too: `just synthetic` writes a module's synthetic samples with canvas's generator, the signed file of its release v0.3.6, on the package of the commit, and CI holds every sample to that package, drawing a synthetic one again and matching a person's by its task's open view; canvas has its samples page with a view of a sample's file, and the first synthetic samples are in: one for each pair of skill and band, 41 in `olympiad-math` and 4 in `grade-3-math`, and two of the author's own samples of `grade-3-math`, written on the samples page of canvas v0.3.11. The samples of volunteers, and the eval runner of `tutor` in CI, wait for stage 2. The seventh epic, done too, gave the modules what format 0.4 allows, in the release `v0.2.1`: a word written by hand reads any one word, so that a word none of the options is a wrong answer, which makes `olympiad-math` 0.3.0, and a time on a clock is checked as a time, with a trap of its own for an afternoon time written as a clock face shows it; 311 tasks and 8 templates have a leading question, standalone's hint, which a hint asked for gives before the first wrong answer; five templates show the picture of their reference task with the numbers of each instance, and three more templates come from standalone's solvers that read a picture, which their own solver now builds from the numbers of each instance: 961 pictures that `just template-pictures` draws. The steps of the method at the end of each of the 41 lessons are a numbered list, whose numbers the app draws, and two of the four solving steps say how to go through them.

## Documents

The documents are in Russian.

- [docs/architecture.md](docs/architecture.md) — the component: what it does and does not do, its data and links, and what is still to decide.
- [docs/decisions.md](docs/decisions.md) — the decision log: what was decided, why, and what was rejected.
- [RUN.md](RUN.md) — the work plan, in epics.

The working rules, for people and coding agents alike, are in [CLAUDE.md](CLAUDE.md) and [AGENTS.md](AGENTS.md). The platform's shared vocabulary, formats, contracts and decisions are in [spec](https://github.com/altnovatrail/spec/tree/v0.15.0), checked out here as a git submodule at `v0.12.0`: clone this repository with `git clone --recurse-submodules`.

## License

[Apache-2.0](LICENSE), for the modules' content and code alike. Each module's manifest names its license by the SPDX identifier `Apache-2.0`, and each module's folder carries the license's text as its own `LICENSE`, since the module's package takes it to every reader.
