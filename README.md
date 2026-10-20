# modules

The official NovaTrail modules: olympiad maths for grades 1–6, carried over from MathTrail, and Spanish A1 for adults, the second module that spec recommends and has yet to settle. A module is data, not code and not instructions to a model: its skill graph, lessons, tasks, traps and techniques in YAML, its media, and an eval set: synthetic samples and samples that adult volunteers and authors made for eval, never learners' work. This repository's CI checks each module with the authors' CLI from module-kit, and its releases carry the signed module packages and the module index that the platform reads.

A component of [NovaTrail](https://github.com/altnovatrail/spec/tree/v0.2.9), a universal open platform for adaptive learning, interactive lessons and AI tutoring. It belongs to stage 1, Format and core, with spec, module-kit, learner-model, engine and sim; together with module-kit it tries the module format v0 on two unlike subjects at once.

## Status

Detailing. The documents come first: the first epic of [RUN.md](RUN.md) decides what the component's card left open, and the code follows in the order RUN.md sets. There are no modules here yet: they come with the module format v0, which spec has still to define.

## Documents

The documents are in Russian.

- [docs/architecture.md](docs/architecture.md) — the component: what it does and does not do, its data and links, and what is still to decide.
- [docs/decisions.md](docs/decisions.md) — the decision log: what was decided, why, and what was rejected.
- [RUN.md](RUN.md) — the work plan, in epics.

The working rules, for people and coding agents alike, are in [CLAUDE.md](CLAUDE.md) and [AGENTS.md](AGENTS.md). The platform's shared vocabulary, formats, contracts and decisions are in [spec](https://github.com/altnovatrail/spec/tree/v0.2.9), checked out here as a git submodule at `v0.2.9`: clone this repository with `git clone --recurse-submodules`.

## License

[Apache-2.0](LICENSE), for the modules' content and code alike. Each module's manifest names its license by the SPDX identifier `Apache-2.0`.
