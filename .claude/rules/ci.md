---
paths:
  - ".github/**"
  - ".devcontainer/**"
  - "justfile"
  - "LICENSE"
  - "test/**"
---

# CI, the release and the toolchain stages

## CI

CI (`.github/workflows/ci.yml`) runs every check on every pull request and on `main`, in a stage of `.devcontainer/Dockerfile`, by spec's `tools/toolchain/toolchain.sh`: `just module-kit-verify` and `just canvas-verify` first, then `just test`, `just modules-check`, `just eval-check`, `just eval-match`, the picture checks and the gate of `sim`; beside them the documents (`docs-check-repo`, through `.github`'s `documents.yml`), the Go checks of `tools/port` and the checks every NovaTrail repository shares — secrets, DCO, CodeQL, Scorecard —, from `.github`'s shared workflows, pinned by SHA in `.github/workflows/`, which Dependabot raises as one group, `shared-workflows`. The signatures and "Gate" run in `toolchain-core`, by `TOOLCHAIN_STAGE=toolchain-core`, as do the release and, by `.github`'s own reading of the Dockerfile, "Documents", "Go / lint", "Go / vuln" and "Go / licenses"; "Modules" and "Go / test" run in `toolchain`. Every job takes the image of its stage from the cache of GitHub Actions by the sha256 of the Dockerfile — `novatrail-<stage>.tar` under the key `<stage>-…`, as `.github`'s workflows name it — and builds it only on a miss, and only a push to `main` saves it; the release builds the core always. The steps of "Modules" share one cache of Go, `.cache/go/` of the checkout, by `env GOCACHE=… GOMODCACHE=…` before the recipe: a new step that runs `go` takes the same prefix. The folder comes from the cache of GitHub Actions by the Dockerfile and `go.sum`, passes `go mod verify` first, and only a push to `main` saves it. The job "Gate" runs the gate of `sim` at `SIM_TAG` and `SIM_COMMIT` of `ci.yml`, raised together by hand, on what `just gate-inputs` writes, in the stage of `sim`'s own Dockerfile, with `-all` when the pin moves. The recipes are in the root `justfile`, and `test/*.test.mjs` are their tests: `just test`, which a task that touches the justfile or `test/` ends with green.

## The toolchain stages

The container is spec's, and differs from it by two files of NovaTrail releases at an exact tag and sum, each checked against the release's signature before its sum was written:
- module-kit's CLI, the binary of a release, `MODULE_KIT_VERSION` and its two sums, which every module is checked with; `just module-kit-verify` checks that release's signature. It is in `toolchain-core`, by the `RUN` of Node, just, Go and buf: `release-files` runs it;
- canvas's generator of synthetic samples, `synthetic-module_<tag>.mjs`, `CANVAS_VERSION` and `CANVAS_SYNTHETIC_SHA256`; `just canvas-verify` checks its signature. It is in `toolchain`, a `RUN` of its own over the core: only `synthetic`, `eval-match` and the tests run it.

A tool a release or a check of the core runs goes into `toolchain-core`, one only the tests need into `toolchain`; a tool in the wrong stage turns the release or the check red on its first run.

Dependabot sees neither: a tag and its sums are raised by hand, together, and the container is rebuilt. The rest of the Dockerfile is spec's, with what the go kind of its repo-template adds in the `devcontainer` stage: gopls, built from the lock file `.devcontainer/gopls/` at `GOPLS_VERSION`, which `test/gopls.test.mjs` holds to the template's at the pinned tag. A move to a tag of spec takes its changes — Go, Node.js, just, gopls — by hand.

## The release

`just release-files <version> <previous>` builds a release's files into `dist/` — the packages as archives `sha256-<hex>.tar.gz` and the cumulative `module-index.json` —, as the shared release of `.github` runs it, with the previous release's index as `previous-file`, which `release-verify` accepted first. A merge into `main` that changes a path of the release — `withdrawals.yaml`, `justfile`, `LICENSE`, `.github/workflows/release.yml`, or anything in `modules/` outside `modules/*/eval/` — publishes the next release after CI; its line is `RELEASE_LINE` of `release.yml`.
