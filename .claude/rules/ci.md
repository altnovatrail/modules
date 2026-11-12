---
paths:
  - ".github/**"
  - ".devcontainer/**"
  - "justfile"
  - "LICENSE"
  - "test/**"
---

# CI, the release and the toolchain stage

## CI

CI (`.github/workflows/ci.yml`) runs every check on every pull request and on `main`, in the `toolchain` stage of `.devcontainer/Dockerfile`, by spec's `tools/toolchain/toolchain.sh`: `just module-kit-verify` and `just canvas-verify` first, then `just test`, `just modules-check`, `just eval-check`, `just eval-match` and the picture checks; beside them the documents (`docs-check-repo`, through `.github`'s `documents.yml`), the Go checks of `tools/port` and the checks every NovaTrail repository shares — secrets, DCO, CodeQL, Scorecard —, from `.github`'s shared workflows pinned by SHA. The recipes are in the root `justfile`, and `test/*.test.mjs` are their tests: `just test`, which a task that touches the justfile or `test/` ends with green.

## The toolchain stage

The container is spec's, and differs from it by two files of NovaTrail releases at an exact tag and sum, each checked against the release's signature before its sum was written (R77 spec):
- module-kit's CLI, the binary of a release, `MODULE_KIT_VERSION` and its two sums, which every module is checked with (R10); `just module-kit-verify` checks that release's signature;
- canvas's generator of synthetic samples, `synthetic-module_<tag>.mjs`, `CANVAS_VERSION` and `CANVAS_SYNTHETIC_SHA256` (R29); `just canvas-verify` checks its signature.

Dependabot sees neither: a tag and its sums are raised by hand, together, and the container is rebuilt. The rest of the Dockerfile is spec's, and a move to a tag of spec takes its changes — Go, Node.js, just — by hand (R31).

## The release

`just release-files <version> <previous>` builds a release's files into `dist/` — the packages as archives `sha256-<hex>.tar.gz` and the cumulative `module-index.json` —, as the shared release of `.github` runs it, with the previous release's index as `previous-file`, which `release-verify` accepted first (R17; R03). A merge into `main` that changes a path of the release — `withdrawals.yaml`, `justfile`, `LICENSE`, `.github/workflows/release.yml`, or anything in `modules/` outside `modules/*/eval/` — publishes the next release after CI; its line is `RELEASE_LINE` of `release.yml` (R28).
