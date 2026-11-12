---
paths:
  - "tools/port/**"
  - "go.mod"
  - "go.sum"
---

# The port: `tools/port/`

The Go code of `modules` is `tools/port`, the port of standalone's tasks into the module format, a Go module with its `go.mod` at the root (R12). Its Go is the `GO_VERSION` of the toolchain stage.

- `just ci-lint`, `just ci-test`, `just ci-vuln` and `just ci-licenses` are its checks; CI runs them through `.github`'s shared `go.yml`, and a task that touches the Go code ends with all four green. Code comes with its tests, and the merge waits for `codecov/patch`, at least 90% of the changed lines covered (R48 spec).
- `just port <module>` writes the tasks of a batch, `tools/port/batches/`, into a module (R12), and `just port-survey <folder>` writes every task no batch names yet into a new module of the format's examples and says what became of each (R20).
- `tools/port/pictures/<the module's folder>/` and `tools/port/template-pictures/<the module's folder>/` are what `just lesson-pictures` and `just template-pictures` draw into a module's `media/`: a change there runs them, and CI runs them with `--check` ([modules](modules.md)).
- `test/port.test.mjs` runs the port on a made-up standalone and holds its output to `module-kit check`, so a change of the port runs `just test` too (R12; R20).
- standalone is read from a clone in a temporary directory, never in the repository (RULES.md spec, Citations).
