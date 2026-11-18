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

The official NovaTrail modules of the first slice: olympiad maths for primary and lower secondary school, carried over from MathTrail, school maths for the third grade on a workbook, and short essays in English for any age. Beside them is a small module in Russian for adults: the tasks of tutor's check board, which the author solves in the app on a learner's path. A module is data, not code and not instructions to a model: a folder whose layout spec's module format sets, with its skill graph, lessons, tasks, traps and techniques in YAML, its media, its license's text, and an eval set: synthetic samples and samples that adult volunteers and authors made for eval, never learners' work.

A component of [NovaTrail](https://github.com/altnovatrail/spec), a universal open platform for adaptive learning, interactive lessons and AI tutoring. It tries the module format, together with module-kit, on modules of two unlike subjects at once.

Its place in NovaTrail: module-kit's CLI checks and builds every module in this repository's CI, and the simulator, sim, gates what each change does to adaptation; the releases carry the signed module packages and the module index, which sim reads and api pins by digest; canvas's samples page takes the tasks of a release, and volunteers' samples come back to the eval sets by pull request. Its card is architecture 3.10 spec, among the documents of [spec](https://github.com/altnovatrail/spec); its own documents, its decisions and its plan are in [docs/](docs/README.md).
