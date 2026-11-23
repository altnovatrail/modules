# modules — architecture

> What the component does, what it does not do and what is still to be decided about it. This is the detail of its card, the one in `spec/docs/architecture.md` whose heading names `modules`. The card stays a summary: if the detail changes what the card says, the card is edited in `spec` by the same route as the questions for `spec`, a row «В реестр spec» of [run/QUESTIONS.md](../run/QUESTIONS.md).
>
> Why it was decided so and not otherwise is in [decisions.md](decisions.md); the work plan is in [RUN.md](../RUN.md). The concepts are those of the glossary in `spec/docs/concepts.md`, the module format is the [format text](../spec/format/README.md), in English. The documents and decisions of `spec` lie in `spec/docs/`.

## 1. Purpose and boundaries

**Does:**
- holds the official modules of the first slice: olympiad maths 1–6 from MathTrail, school maths for grade 3 on a workbook and short essays in English for any age. They have two subjects — maths and writing;
- holds, beside them, a small module of maths in Russian for adults, `novatrail/check-board`: the eleven tasks of the check board of `tutor`, so that the author solves them in the app on the path a learner takes, and so that its eval set checks every new tag of `tutor` once the eval runner of epic E6 is in CI; the module is described below, with the other sets. The goal of the first slice, three modules of two subjects, does not change with it;
- holds all the modules in one repository while the format changes: all the modules move to a new format version in one pull request. They may split into a repository per module once the format is v1;
- holds each module as a folder that the [module format](../spec/format/README.md) defines:
  - at the folder's root — the module manifest `manifest.yaml`, the skill graph `skill-graph.yaml`, the catalogs `traps.yaml` and `techniques.yaml` and the solving steps `solving-steps.yaml`;
  - the tasks and the lessons — a file for each, `tasks/<id>.yaml` and `lessons/<id>.yaml`;
  - the media files — in `media/`, the eval set — in `eval/`;
  - the license text — the mandatory `LICENSE`, the license notices — `NOTICE`, if they are needed.

  A module is written in YAML 1.2.2 with the format's restrictions, and is checked by the JSON Schema 2020-12 schemas of the format version that its module manifest names. Where the folders themselves lie and what they are called is the repository's tree, below. There are no widgets in the first slice;
- keeps in the eval set synthetic samples and samples that adult volunteers and authors made for eval. The sample file and the eval set are set by the eval protocol of the [contracts text](../spec/contracts/README.md): a sample is the message `Sample` of the package `novatrail.eval.v1`, and a module's eval set is each `.json` under its `eval/`. The layout inside `eval/` is flat, `eval/<id>.json`, as the layout of the samples in `eval/` below says;
- checks the modules in CI, as the checks of CI below describe:
  - with the CLI of `module-kit` — by the format's schemas and rules: a graph without cycles, references, anchor skills, keys, with a template's solver too, fallback texts — with the leak check; it also compares a module's new version by the rules of semver with each published version below it, withdrawn ones too, by the format's rules, and a version in another format version by the places that the format keeps in all its versions. CI takes the CLI as a release binary of `module-kit` at an exact tag and, before its first run, checks its signature with `spec`'s recipe `release-verify`, rather than installing it through `go install`;
  - with the simulator `sim` at the tag that `ci.yml` pins, if what the simulator reads from the package has changed — `sim` decides this —, on the pull request's package, with invariants and absolute thresholds, without bands;
  - the files of the eval-set samples — with `spec`'s recipe `sample-check`, without a key, on every pull request, and from a fork too: with the same code that checks `spec`'s examples, as the samples in the eval set below describe;
  - from stage 2, when `tutor` and the samples page of `canvas` appear, — with the eval runner of `tutor` on the module's eval set, with a separate key of the platform at an AI provider — only for eval, not the key of the service account of `api` and not the author's own key. The job with this key, like any job with a secret, does its work only on a push to `main` and on a pull request from a branch of the repository itself. On a pull request from a fork it is not skipped but fails and names the reason, and the merge waits until a maintainer reads the change and carries it into a branch of the repository. A pull request of Dependabot also comes from a branch of the repository, but gets no Actions secrets. The key is a secret per provider in Secret Manager of `infra`, reached through Workload Identity Federation by a pool of its own, `github-eval`, only from a job of the GitHub environment `eval`; a pull request of Dependabot is evaluated after the change is carried into a branch of the repository;
- on a release, the CLI of `module-kit` builds the module packages and the module index inside the shared release workflow of `.github`, and that workflow signs them:
  - a package is the server part, the device part and the inventory. The server part is the module's canonical view: all the content, keys, solutions of the tasks and traps included, with the tier of each skill and the numbers of each instance of a template. The device part is `LICENSE`, `NOTICE` and the media files, without the keys and the texts of the module. The inventory names the module and both parts by digests, and the package digest is the digest of the inventory. The eval set is not part of the package: the eval runner takes the samples from the module's folder in CI of `modules`;
  - the index is the file `module-index.json` in each release, one for all modules and cumulative. It names each published version of each module by the digest of its package, with the tag of the release in whose files the package lies, with the anchor skills that the version teaches, and with the mark of withdrawal. Entries do not leave the index and, apart from the mark of withdrawal, do not change;
  - a release builds the index from the previous release's index and its own packages, after the number and before publication. The previous index is downloaded and checked with the recipe `release-verify` by a step of the shared `release.yml` on the runner, and the recipe `release-files` of `modules` gets the already checked file;
  - a version is withdrawn by the next release: its index marks the version's entry `withdrawn` with a reason — `wrong_key`, `answer_leak` or `unsuitable_content`. A withdrawal is not undone, and a withdrawn version stays only for a repeat of a verdict. A withdrawn newest version is fixed by a version above it by semver. A fix that is ready comes out in the same release as the withdrawal, and the withdrawal does not wait for the fix;
  - the name of a package's file in a release is derived from the index's entry — without going through the release's files and without the REST API: so the package is found by the build of the image of `api`, CI of `canvas` and a repeat of a verdict. The name itself is given under the release, below;
- delivers the packages by stages (Т-17):
  - stage 1 — the packages and the index in releases of `modules`;
  - stage 2 — at the build of the image of `api` the digests of the packages are taken from the release's index. A pull request pins the release of `modules` and its index in the configuration of `api` and serves as the anchor of trust, and its CI checks the index's signature with the recipe `release-verify`. With the index the marks of withdrawal reach `api` too. The files of the device part are served by `api`;
  - old versions of packages for a repeat of a verdict are taken from releases of `modules`.

**Does not:**
- holds no code of the platform: a canvas kind and a check kind are code of the platform, shared by all modules, and a module only chooses one and sets its parameters;
- holds no instructions to a model: the prompts belong to the platform, and a module's free text reaches the model only as a separate marked block of data;
- holds no learners' work — neither children's nor adults': a module is published under an open license and cannot be withdrawn. Learners' work from `app` is not taken into eval, and a child's sample made for eval lies outside the repositories, under a parent's separate form;
- does not hide the keys: the keys of open modules are public, like the modules themselves. `spec`'s rule on keys guards no secret, but only keeps the key out of the learner's app;
- does not own the format: the format's schemas, rules, examples and test vectors lie in `spec/format/`, and the code that executes them — the check, the package, the index, the open view, the check kinds — is `module-kit`;
- does not publish third-party modules: the recommendation of О-9 is to publish only the official ones at stages 1–3, and the registry service will appear at stage 4.

**Data and links:**
- there are no learners' data and no secrets, as in no repository. How volunteers' samples get into a module is described under the samples in the eval set, below; the consent to them, and what lies in a sample from the samples page, `spec` sets;
- at the build it depends on nothing but `spec`. CI takes the CLI of `module-kit` and `sim` at an exact tag, and from stage 2 the eval runner of `tutor`. A file of a NovaTrail release — the binary of the CLI, the previous release's index — enters CI only after `release-verify`;
- CI runs the shared checks of `.github`: the secrets search, DCO, CodeQL, Scorecard and the documents check. The merge into `main` is held by all of them but Scorecard, together with the quality gate of SonarCloud, and the ruleset of `main` allows only a merge commit;
- the packages are read by `sim` — by the types of the canonical view from `spec`, without `module-kit`, from stage 1, and the gate file of `sim` pins the release of `modules` by its tag and the digest of its `module-index.json`, and the packages by digests, — and from stage 2 by `api`, by the digest from its configuration, and by the eval runner of `tutor`. The learner model marks its state with the release of `modules`;
- from stage 2 CI of `canvas` also takes a release of `modules` at an exact tag, having checked its signature with the recipe `release-verify`. The open views of the tasks for the build of the samples page are built by the CLI of `module-kit` from the package's server part, the media files it takes from the device part, and the marks of withdrawal CI reads in the index of the newest release of `modules` and builds no page from a withdrawn version. Until `canvas` moves the build to `modules`, by a task of its own, the measurement tasks come from a release of `tutor`;
- already at stage 1 the generator of synthetic samples of `canvas` reads the module package that `module-kit build` built at the commit where the samples go and that the call gives it — the synthetic samples of the eval set; and the samples page of `canvas` opens a sample file, so that it is looked at before the push, as the samples in the eval set below describe;
- volunteers' samples come from the samples page of `canvas` by a pull request that the eval-set organizer makes; nothing reaches the platform's server on the way. For the first measurement the page has no hosting: its build is downloaded from a release of `canvas` and opened on the device, and while it works it does not go to the network: the tasks and media files are put into it at the build. Child samples are made with the same page, but without a pull request, and do not get into the repository;
- the license is Apache-2.0, for the content, for the code and for the volunteers' samples in the eval sets. The module manifest names it by the SPDX identifier `Apache-2.0`, and the module's folder carries its text as the file `LICENSE`. The notices that the licenses of the content require — for example MIT for the tasks of standalone, as the port of standalone below says — stand only in the folder's `NOTICE`: the comments of a module's files do not get into the package. `LICENSE` and `NOTICE` are UTF-8, their lines end with LF, and git does not convert their line endings: the package carries them byte for byte, by the module format. Their language is the license's, not the module's: a module of any explanation language holds the English text of Apache-2.0;
- the order of work is set by `spec`'s cross-repository plan, [COORDINATION.md](../spec/COORDINATION.md). The module of olympiad maths serves the first application, on which the author solves tasks locally in a browser, and its module manifest names `adult` in the audience too, as the first batch below says; the sets of school maths for grade 3 and of essays are small, to try the format;
- stage 1, together with `module-kit`: the module format v0 is tried at once on three modules of two subjects.

## 2. To decide when detailed

Every item is decided, and this section names only what is deliberately deferred. The layout of the modules and the volunteers' samples are items of the card in `spec`, the port of standalone's tasks is beyond it, and the layout of the samples in `eval/` is what the eval protocol left to `modules`; each item keeps its number, which the plan's archives cite, and names where below its decision is described. A new item is the next free subsection, by the rules of [RUN.md](../RUN.md): the gist, what is known, the options and the task. A decided item names where its decision is described, the decision itself moves to how it is built, below, and its why is in the log. A question that goes beyond this repository goes «в реестр spec», as a row of [run/QUESTIONS.md](../run/QUESTIONS.md).

**Deliberately deferred:**
- what is still to build — epic E6 in table 2, [run/EPICS.md](../run/EPICS.md), by its row of what to build, below: the volunteers' samples and the eval runner of `tutor` from stage 2, `report-check` on its reports, and the samples of templates, of a choice by hand, of essays and of the module of the check board;
- the module of olympiad maths in Russian is another module, by О-2: the explanations of this module are in English;
- the parameter of 12 or 24 hours of the check kind `time` — the format will give it by a minor change when content asks for it. Until then the 12-hour miss is held by the trap `twelve-hours-off` of the format example, as the port of standalone below says.

### 2.1. The layout of the modules in the repository

- **The gist:** where in the repository the modules' folders lie and what they are called; how CI finds which module has changed and what to run on `sim`; how releases come out — which modules a release builds, how the packages lie in its files and where the withdrawals lie; which tag CI of `canvas` takes. What lies inside a module's folder is set by the format.
- **Decided:** the repository's tree, the checks of CI and the release, below.

### 2.2. Volunteers' samples in the eval set

- **The gist:** how a sample that an adult volunteer made on the samples page of `canvas` gets into a module's eval set: who opens the pull request, what is checked before the merge, where the consent to publication is marked. One path for all the modules: handwriting on a sheet, the answer slots of a workbook and essays. The consent itself is `spec`'s.
- **Decided:** the samples in the eval set, below.

### 2.3. The port of standalone's tasks

- **The gist:** what of MathTrail's maths goes into the module of olympiad maths and what each part turns into, which tasks go into the first batch, and on what terms their content enters the module under Apache-2.0.
- **Decided:** the port of standalone and the first batch, below.

### 2.4. The layout of the samples in `eval/`

- **The gist:** where inside a module's `eval/` the sample files lie and what they are called. The eval protocol of the [contracts text](../spec/contracts/README.md) leaves this to `modules`: a module's eval set is each file under its `eval/` at any depth whose name ends in `.json`, except files and folders starting with `.`.
- **Decided:** flat, `eval/<id>.json`, and nothing starting with `.` in `eval/` — the samples in the eval set, below.

## 3. How it is built

What the items to decide settled: the parts of the component, its data and how it is linked to other repositories — the layout of the modules, the volunteers' samples, and the port of standalone with its first batch. Why it was decided so and not otherwise is in the log. The modules are in `modules/`, and which version of each the releases carry is its module manifest and the module index of the newest release; what to build, in which epic and what it waits for is under what to build, below.

### 3.1. The repository's tree

```
modules/                     the repository root
├── modules/                 module folders, and nothing else
│   ├── olympiad-math/       novatrail/olympiad-math — olympiad maths 1–6
│   ├── grade-3-math/        novatrail/grade-3-math — school maths for grade 3 on a workbook
│   ├── english-essays/      novatrail/english-essays — short essays in English
│   └── check-board/         novatrail/check-board — the tasks of the check board of tutor, in Russian
├── withdrawals.yaml         withdrawals, from the first withdrawal on — an input of the index build
├── justfile                 CI recipes and the recipe release-files
├── test/                    tests of the justfile recipes, node --test
├── go.mod, go.sum           the Go module github.com/altnovatrail/modules: the code of the tools
├── tools/port/              the converter of standalone tasks: code, batches batches/, overlay files overlay/,
│                            lesson pictures pictures/ and the sources of instance pictures template-pictures/
├── .gitattributes           modules/** -text
├── docs/                    the component documents; the language of each is in CLAUDE.md;
│                            docs/README.md — their index, for people
├── spec/                    the submodule spec at a tag
├── .github/                 CI: calls of the shared workflows, own jobs, release.yml
├── .devcontainer/           the spec container of the kind go; in the stage toolchain-core — the CLI
│                            module-kit, in toolchain — the generator of canvas
├── run/                     the work plan's tables 2–4 and its open tasks outside an epic
├── README.md, LICENSE       the README, a description alone, and the license of the repository
└── CLAUDE.md, AGENTS.md, RUN.md
```

`dist/`, where the recipe `release-files` puts the release files, git does not keep — `.gitignore`. There is no `withdrawals.yaml` while nothing is withdrawn, as the release below says.

- **The module's folder** is `modules/<name>/`, and `<name>` is the second half of the module ID: the module ID in the folder `modules/<name>/` is `novatrail/<name>`. What lies inside is set by the module layout of the [format text](../spec/format/README.md):

  ```
  modules/<name>/
  ├── manifest.yaml, skill-graph.yaml, traps.yaml, techniques.yaml, solving-steps.yaml
  ├── tasks/<id>.yaml, lessons/<id>.yaml
  ├── media/<path>             pictures and sound; they go into the device part of the package
  ├── eval/<path>              eval-set files; they do not go into the package
  └── LICENSE, NOTICE          the text of Apache-2.0 and the license notices
  ```
- **In `modules/` — only the modules' folders:** neither other files nor files and folders starting with `.`. Each module in `modules/` on `main` comes out in a release — as its package, without `eval/`. A module's draft lives in a branch.
- **The names** name the subject and the audience, not the explanation language: the languages of the slice are О-2. A module in another explanation language is another module, by the module format, and its name will get a suffix with a BCP 47 tag, for example `olympiad-math-en`. A module ID never changes: the full IDs of its skills and tasks are kept by the learner model and the progress.
- **`LICENSE`** lies both at the root — the repository's license —, and in the folder of each module: this text the package carries to the learner's device.
- **`.gitattributes`** is `modules/** -text`: git converts the line endings of no file of the modules, neither at a commit nor at a checkout. `LICENSE`, `NOTICE` and the media files the package carries byte for byte, and CR LF would change the package digest. A `LICENSE` or `NOTICE` with CR would not pass the rule `module.plain-text` anyway.

### 3.2. The checks of CI

| Check | When | Which modules |
|---|---|---|
| the CLI of `module-kit`: the format's schemas and rules | every pull request and push to `main` | all |
| the tree: in `modules/` only the modules' folders, the ID of each is `novatrail/<folder>` | the same | all |
| the package and the comparison with what is published | the same | all |
| the simulator `sim`: the job "Gate", `just gate modules` of `sim` at its pinned tag | every pull request and push to `main` | those in which what the simulator reads from the package has changed — `sim` decides — and new ones; all — when a pull request changes the tag of `sim`; none — in a format transition |
| the layout of `eval/` and the eval set of each module: the recipe `just eval-check`, and in it `spec`'s recipe `sample-check --eval-set novatrail/<folder>`, one call per module | every pull request and push to `main`, without a key, and from a fork too | all that have `eval/`; the mode — those whose `eval/` has `.json` files, as the samples in the eval set below describe |
| the samples of `eval/` are samples of the commit's package: a synthetic one is drawn again by the generator of `canvas` and compared with the file in everything but the version and the digest of the package, for a person's sample the digest of the open view is checked by `module-kit open`; the recipe `just eval-match` | every pull request and push to `main`, without a key and without the network, and from a fork too | all whose `eval/` has samples |
| the eval runner of `tutor` | pull request, from stage 2 | with a changed `eval/` or package, new ones too; all — when a pull request changes the tag of `tutor`, as the samples in the eval set below describe |

- **All modules on every pull request.** The CLI of `module-kit` checks each folder `modules/*/` and builds its package. So a new tag of `spec` or of `module-kit` re-checks all the modules without a special rule.
- **The CLI is pinned in one place** — in the stage `toolchain-core` of `.devcontainer/Dockerfile`, under `toolchain`: the release binary of `module-kit` at an exact tag and sha256 sum.
  - The recipe `release-files` takes the same binary too, as the release below says, so CI checks the modules with the same version of the CLI that the release will build them with.
  - The binary's signature is checked by the recipe `release-verify` in CI of each pull request and push to `main` — the job "module-kit signature", the recipe `just module-kit-verify`: the release's `SHA256SUMS` and the check of its lines against the sums of the `ARG`. From then on the sum holds the binary.
  - This is how the container differs from `spec`'s container; CLAUDE.md names it too.
- **The generator of synthetic samples is pinned there too** — the file `synthetic-module_<tag>.mjs` of a release of `canvas` at an exact tag and sha256 sum, in the stage `toolchain`. The release's signature is checked by the job "canvas signature", the recipe `just canvas-verify`: the release's `SHA256SUMS` through `release-verify` and the check of the generator's line against the sum of the `ARG`. The recipe of synthetic samples `just synthetic` is a person's step in the devcontainer, not CI, as the samples in the eval set below describe; in CI the generator is called by `just eval-match`, the check of the samples of `eval/` against the commit's package.
  - The form of the call of the commands `check`, `build`, `index` and `diff` is set by `module-kit`, in its `docs/architecture.md`: the recipes of `modules` follow it and check the call with their tests.
- **The comparison with what is published.** CI takes the `module-index.json` of the newest release of `modules` and checks its signature with the recipe `release-verify`. For each module:
  1. the version is in the index — the digest of the package that CI built must match the digest of the entry. Otherwise the content was edited without raising the version, and the release would refuse by the format's rule `module-index.one-package`;
  2. the version is not in the index — it must be above each published version of the module, withdrawn ones too. The CLI of `module-kit` compares the module with the package of each published version by the format's rules of version comparison: which rules go against each version and which against the newest, the format says. CI downloads the packages by the entries, as the release below names them, and checks them by the chain of digests from the checked index;
  3. a published version is written in another format version — while the format is v0, this happens after each of its transitions. The CLI reads its package by the places that the format keeps in all versions: the inventory, the names of the package's files, the languages, the bands, the skills and `retired_skills`, the keys of the tasks, lessons, traps and techniques. The rules about skills, `retired_skills` and IDs go across a transition too, and the version part, while the format is v0, is major;
  4. a repository with no release has nothing to compare with: all its versions are new.

  How CI does this: the newest release is the greatest tag `vX.Y.Z` from `git ls-remote`, without a token, that has a `module-index.json`; a tag whose index answers 404 is a release that did not come out, and the log names it. Its `module-index.json` is checked by `release-verify`, and the archive of each version of the modules of `modules/`, withdrawn ones too, comes from the address in the entry, as the release below derives it. Only `inventory.json`, `server.json` and `device.json` are taken from the archive; an archive with a link, a special file or a path outside the folder is rejected whole. Then comes `module-kit diff --index … --published … <package>` over the package that `build` built, and without a release `--first`: that there is no release the recipe says explicitly, with the file `first`, and an index silently skipped does not pass. The chain of digests from the index to the package's files is checked by `diff`, and its exit code is the job's code.
- **Whether what the simulator reads has changed** — for the gate of `sim`. The gate is everything that the simulator reads from the package, and whether it has changed is decided by `sim` itself, with its command "has it changed". The simulator runs on the pull request's package that CI built, with invariants and absolute thresholds, without bands: the pull request's package has no base to compare them with. `sim` reads `server.json` by the types of the canonical view, without `module-kit`. A pull request that changes the tag of `sim` runs all the modules.
- **A format transition goes without `sim`.** If the newest published version of a module is written in another format version, the simulator is not run in this pull request: `sim` at its former tag will not read the package of the new version. So in a transition nothing that the simulator reads from the package is changed — neither the graph nor the bindings of the tasks — except the writing itself in the new format version, and the review of `modules` holds this. All of that is changed by a separate pull request, with the simulator, once `sim` has moved to the new format version.
- **The job "Gate".** `ci.yml` pins `sim` by `SIM_TAG` and `SIM_COMMIT`, the tag and the SHA of its commit, and the job checks that the tag is on that commit; Dependabot sees neither, and a person raises both together. In this repository's stage `just gate-inputs .cache/gate` writes the published versions, as `_published` writes them for `modules-check`, the package of every module, and the arguments of `just gate modules` by the call form of `sim`'s README: `-index` and a `-published` for each published package, or `-first`, then every package. The gate runs from the root of the checkout of `sim`, in the `toolchain` stage of `sim`'s own Dockerfile at that commit, and `sim` decides what of each package changed and runs only that. `-all` goes first when the commit the run is checked against — the pull request's base, the push's previous commit — pins another `SIM_COMMIT`, none, or cannot be read. The image of `sim`'s stage and the caches of Go of the gate are in the cache of GitHub Actions; a pull request only reads them. The job's exit is the gate's: invariants and absolute thresholds. A module whose graph no population of `sim` at `SIM_TAG` runs on, as `check-board`, `sim` names and passes over: populations on a new graph are a change of `sim`. `SIM_TAG` must read every check kind the modules use, `time` of `olympiad-math` too.
- **The jobs and the recipes.** The own jobs of `ci.yml` take the image of their stage from the cache of GitHub Actions, by the sha256 of `.devcontainer/Dockerfile`, or build the stage on a miss — `toolchain-core` for "module-kit signature", "canvas signature" and "Gate", by `TOOLCHAIN_STAGE`, `toolchain` with the generator of `canvas` for "Modules" —, and run the recipes of the `justfile` in it with the script `spec/tools/toolchain/toolchain.sh` from the submodule. Only a push to `main` fills the cache, under the key and the file of each stage that "Documents" and the jobs of "Go" take too, so there is one entry a stage for all.
  - The steps of "Modules" share one cache of Go, `.cache/go/` of the checkout, which `toolchain.sh` mounts at `/work`, so that the port is compiled once. It comes from the cache of GitHub Actions with no restore key; what the cache gave passes `go mod verify` or is deleted with a warning, and only a push to `main` saves it. Its key is the runner, the Dockerfile and the root `go.sum` alone, not `**/go.sum`, and the folder is in the checkout, not in `RUNNER_TEMP`, since `toolchain.sh` mounts `RUNNER_TEMP` read-only. The composite action `toolchain` of `.github`, which gives a caller's own job the image and the caches of Go, restored by a prefix too, keeps them in `RUNNER_TEMP`: moving "Modules" onto it is an open task of [run/TASKS.md](../run/TASKS.md) (sources 4.23).
  - "module-kit signature" is `just module-kit-verify`, "canvas signature" is `just canvas-verify`; "Modules" waits for both and runs `just test`, then `just modules-check`: the tree of `modules/`, the version of the CLI on `PATH` against `MODULE_KIT_VERSION`, `check` and `build` of each folder into a temporary folder, and the ID from the package's `inventory.json`. Then `module-kit diff` compares each package — with the published versions by the index of the newest release or, with no release, with `--first`. Without `modules/` the recipe prints `modules/: no modules yet`, exits with code 0 and does not go to the network. The refusal of one module does not stop the check of the others; the recipe's code is the greatest of the codes of what failed: 1 — the tree, an ID other than `novatrail/<folder>` and a refusal by the rules, 2 — `module-kit` could not do its work, or the published versions could not be read.
  - After it "Modules" runs `just eval-check` — the layout of each module's `eval/` and its eval set with the mode `sample-check --eval-set`, as the samples in the eval set below describe, — then `just eval-match` — each sample against the package of this commit, — the check of the lessons' pictures, `just lesson-pictures modules/olympiad-math --check`, — and the check of the pictures of the templates' instances, `just template-pictures modules/olympiad-math --check`.
- **No filter `paths` on the workflow.** What to run is decided by the job itself. The checks of a workflow that a filter skipped hang "Pending", and a required one holds the merge (sources 2.3).

### 3.3. The release

- **Who and when.** The own `release.yml`, by the release template of `.github`, calls the shared `release.yml` of `.github` — pinned by SHA in its `uses:` — after a green CI on `main`, if something of `paths` has changed since the previous release. The number is the next `vX.Y.Z` of its line, as for any NovaTrail repository: the line is `RELEASE_LINE` of `release.yml`, and a new format version that the index carries may need a new line.
  - **A dry run on a pull request.** The same call runs on each pull request from a branch of the repository itself with `publish: false`, as in `canvas`: the files are built from `git archive` in the stage `toolchain-core`, signed and checked, but nothing is published. A pull request from a fork skips the dry run — it is given no `id-token`, and there is nothing to sign with. The dry run is not in the ruleset of `main`. The input `previous-file: module-index.json` gives the build the previous release's index (sources 4.1). The releases are immutable, and the ruleset of the tags `v*`, `v-tags`, holds the rules `update`, `deletion` and `non_fast_forward` without bypass: both are conditions of the signature check.
  - **The epic branch.** An epic that puts out one release at its end gathers the pull requests of its tasks in the epic's branch, not in `main`: a release comes out after each green CI on `main` if something of `paths` has changed, and a change of the `justfile` touches them too. The branch is added to the include of the ruleset of `main`, so that it holds the same required checks, and is removed from it before the branch is deleted, each a step outside the repository that [RUN.md](../RUN.md) takes only on the author's confirmation; the dry run of the release on its pull request publishes nothing, the button `workflow_dispatch` is not pressed on it, and the epic's last task merges it into `main`.
- **`paths`** is what the release is built from: the modules, without their eval sets, and what builds the release files from them. The list for the kind `modules` belongs to the release template of `.github`; until the template carries it, this is the list, `release.yml` carries the same, and where the template diverges from it, the template is right:

  ```
  :(glob)modules/**
  :(exclude,glob)modules/*/eval/**
  withdrawals.yaml
  justfile
  LICENSE
  .github/workflows/release.yml
  ```

  - No release comes from the eval sets, the documents, CI other than `release.yml`, the submodule `spec` and `.devcontainer/Dockerfile`. The package and the index depend only on the files of the module and the format version of its module manifest, by the module format, not on the version of the CLI in the Dockerfile.
  - The price is the whole `justfile`: a change of a CI recipe in it will give a release with the same index.
  - The magic `glob` is there so that `*` does not match `/`: without it `modules/*/eval/*` would exclude a module's `media/eval/` too.
- **The recipe `release-files VERSION PREVIOUS`** runs in the container of the stage `toolchain-core`, which the shared release builds from the Dockerfile with no cache, on the commit's tree without git, submodules, network and tokens.
  - The CLI of `module-kit` is the same binary of this stage as CI's, as the checks of CI above say.
  - `PREVIOUS` is a read-only path to the `module-index.json` of the newest release, which the shared release downloaded and checked with `release-verify` before the build. With no release before, it is an empty string, and the recipe tells the build explicitly that there is no previous index.

  The recipe builds the packages of all the modules and the index: from the previous index, its own packages, the tag `VERSION` and `withdrawals.yaml`, by the module format. The package of a version that the previous index already names does not change its entry. If the digest of such a package is different, the build refuses by `module-index.one-package`.

  How it does this:
  - the input: `VERSION` is a tag `vX.Y.Z`; `PREVIOUS` is an empty string, then `--first`, or a file, then `--previous`; `withdrawals.yaml` — `--withdrawals`, and without it `--no-withdrawals`. Anything else is code 2, and in `modules/` there are only the modules' folders, otherwise code 1;
  - the tree, the CLI on `PATH` against `MODULE_KIT_VERSION`, `check` and the ID of each module — with the recipe `_modules-check`, without the comparison with the published ones: that is the business of CI of the same commit, and what the release needs of it is held by `index`. A relative `PREVIOUS` is from the folder the recipe is called from;
  - `module-kit build` of each module into a temporary folder, then `module-kit index … --tag VERSION … --out <file> <packages>` — by the call form that `module-kit` sets; the code of `index` is the recipe's code;
  - an archive — only for a package whose entry in the new index names `VERSION`. Before the archive the recipe checks the entry's `package_digest` against the sha256 of the bytes of the package's `inventory.json`, and the package's composition — `inventory.json`, `server.json`, `device.json` and the folder `device/`, and nothing more: otherwise code 2;
  - `dist/` is built in a temporary folder and takes its place whole, at the end: the old `dist/` is removed at the start, and a run with a refusal leaves no `dist/`.

  If not one module has changed since the previous release, and the change is in the `justfile` or `release.yml`, the release carries the same index and `LICENSE`.
- **The release files** are what the recipe puts into `dist/`. The names are of letters, digits and `. _ + -`, without folders. The bundle `<file>.sigstore.json` for each is put by the shared release:

  | File | What it is |
  |---|---|
  | `module-index.json` | the module index: each published version of each module, from the first release |
  | `LICENSE` | the text of Apache-2.0 from the root |
  | `sha256-<hex>.tar.gz` | the archive of a package — for each version that this release publishes, and only for it |

  The package of a version that the index already named lies in the release that its entry names, and does not come out a second time. A release that only withdraws versions carries the index and `LICENSE`.
- **The archive of a package** is named by the package digest — the entry's `package_digest`, with `:` replaced by `-`, — with `.tar.gz`. Inside is the package as files, at the root of the archive, without a top folder, by the module format. Apart from these files and their folders, there is nothing in the archive — neither links nor other entries:

  ```
  inventory.json  server.json  device.json
  device/LICENSE  device/NOTICE  device/media/<path>
  ```

  The archive is the same byte for byte on any runner and at any time of the build: GNU tar — `--format=gnu --sort=name --owner=0 --group=0 --numeric-owner --mode=u=rwX,go=rX --mtime=@0`, the entries in the order `inventory.json`, `server.json`, `device.json`, `device/`, without `./` before the names; then `gzip -n -9`, without a name and a time in the header. `tutor` builds its archive the same way (sources 4.14).

  The address is derived from the index's entry, without going through the release's files and without the REST API: `https://github.com/altnovatrail/modules/releases/download/<release>/sha256-<hex>.tar.gz`. Nothing holds the bytes of the archive, and nothing needs to. A reader checks the inventory against the package digest, the parts against the inventory, the files against the list of the device part, and takes only the paths that this list names.
- **Withdrawals** are `withdrawals.yaml`, YAML by the module format's rules. The form is `$defs/withdrawals` of the index's schema, `module-index.schema.json`: module ID → version → withdrawal reason. The file is read by the CLI's command "build the index". For example:

  ```yaml
  novatrail/olympiad-math:
    "1.2.0": wrong_key
  ```

  - The form requires at least one entry, so there is no file while nothing is withdrawn. Then the recipe tells the build that there are no withdrawals.
  - The file is in `paths`, so a pull request with only a withdrawal gives a release.
  - Entries do not leave the file: a repeat of a withdrawal with the same reason changes nothing, and another reason would break `module-index.entries-kept`. A version that comes out in the same release can be withdrawn too.
  - If the newest version of a module is withdrawn, it is fixed by a version above it by semver, and which part to raise the rules of semver decide. A change of a published task's key is major, so the fix of a version withdrawn by `wrong_key` is the next major. A fix that is ready goes in the same pull request as the withdrawal and comes out in the same release: otherwise, until it comes, the current version would be the version below the withdrawn one, and if a new major is withdrawn, the previous major, and there would be nothing to carry the withdrawn version's evidence about its skills over with. The withdrawal does not wait for the fix.
  - A version older than the newest cannot be fixed: one line of versions lives in a module's folder, and a new version is above all the published ones, as the checks of CI above require.
- **Who takes the release.**
  - CI of `canvas` — from stage 2, at an exact tag `vX.Y.Z` of a release of `modules`: the number of the repository itself, not a module's version. By the tag it reads the index, and it takes the packages from the releases that the index's entries name. The index of the newest release it reads only for the marks of withdrawal.
  - `api` — from stage 2: a pull request pins the release and its index, as the purpose above says, and the build of the image takes the packages by the entries.
  - `sim`, which reads packages from releases of `modules`, finds them by the same layout: a package's address is derived from the index's entry.
- **Old releases** are not deleted or rewritten: releases are immutable, and a verdict is repeated on old packages. Each release but the first needs the previous one's index, which the shared release gives it.

### 3.4. Samples in the eval set

A module's eval set lies in `eval/` of its folder and is not part of the package. The samples in it are of three origins: synthetic samples, samples of adult volunteers and samples of authors.

The CONTRIBUTING of `.github` does not accept volunteers' samples (sources 4.7); they come by the path below. The sample file is set by the eval protocol of the [contracts text](../spec/contracts/README.md): the message `Sample` of the package `novatrail.eval.v1`, one JSON per sample, without the task's key, with a license for an eval-set sample and with the markup `markup`, and the crops in it are PNG in base64. A module's eval set is each `.json` under its `eval/`, and the layout inside `eval/` the protocol left to `modules`. The layout is flat: each sample is a file `eval/<id>.json`, where `<id>` is its ID, and in `eval/` there are neither folders, nor other files, nor anything starting with `.`.

- **The roles:**
  - **the volunteer** — an adult who writes a sample on the samples page of `canvas` under a signed consent form;
  - **the eval-set organizer** — a maintainer of `modules`. The organizer accepts the form, gets the file, checks it and adds it;
  - **the module's author** adds their own samples themselves, under their own form.
- **The path of a volunteer's sample:**
  1. Before the first sample the organizer accepts the signed form (П-18). The form is `volunteer-form` of the text set `eval-set` in force, together with its notice `eval-set-notice`, and the organizer takes them from the legal part of a `spec` release, the file `legal-<version>.tar.gz`. The form lies outside the repositories — in the operator's private Cloud Storage bucket, under a prefix of its own, with the organizer's access, — and agents are not let near it (П-12). The organizer acts for the operator, ALT EdTech LLC, on its instructions: the controller of the samples and the forms is the operator. In the form:
     - the name, the signature, the date and the version of the text set;
     - the declaration "I am 18 or older, and under the laws of where I live I may make this agreement";
     - the text of DCO 1.1 word for word and the certification of clause (a);
     - publication under Apache-2.0 — and the consent withdrawal, below;
     - runs at Anthropic and OpenAI in CI of `modules`: a new provider is a new version of the set `eval-set` and a row «в реестр spec»;
     - samples without names and signatures, and essays on made-up topics;
     - the term of the form.

     The organizer does not take the form of a volunteer from the EU, the EEA or the United Kingdom while the operator has no representative under Art. 27 GDPR, and the form asks the volunteer to state this.
  2. The volunteer opens the samples page and chooses the purpose «eval-набор» ("eval set"). On it they state that they are 18 or older and that the form was accepted, write the sample and save the file on their device.
  3. The file and the scan of the form reach the organizer by an upload link valid for at most 12 hours from the session, into the same bucket. This is not the platform's server and not GitHub. It is this link that the page names: «передать файл организатору eval-набора» ("pass the file to the eval-set organizer").
  4. The organizer checks the file — by the table below, before the file gets into any branch. A file that did not pass is not edited and not added, but erased at once: the volunteer writes the sample anew.
  5. The record of the form, outside the repository, gets the SHA-256 digest of the file. Under which form each sample was added the organizer knows, and the repository does not. The form and the record live as long as at least one sample under it is in the repository or in a released version of a module; a form without an added sample is erased 12 months after signing. The scores of an essay's criteria are given by two people, each on their own, and where they diverge, a third decides, and one mark per criterion goes into the markup.
  6. The organizer puts the file into the module's `eval/` as `eval/<id>.json`, by the sample's ID, in a branch of the repository itself. The file's name is the sample's ID, and the organizer does not change the content: clause (c) of the DCO holds while the file is byte for byte the one the volunteer certified. The organizer signs off the commit with their own name under clause (c). The commit's author is the organizer too: the DCO check checks the sign-off against the author (sources 4.6). Neither `Co-authored-by` nor the volunteer's data are in the commit's message, the branch's name or the pull request's text. If there are samples of several volunteers, they go in one commit, so that a commit does not show which files come from one person.
  7. The pull request passes CI and is merged with a merge commit. It does not change a package and gives no release: `eval/` is neither in the package nor in the release's `paths`. Once the file is added, the organizer erases its copy in the bucket at once: the file is public, and the copy adds nothing. A file not added is erased by the bucket's lifecycle rule 12 months after the upload, so the organizer adds a file while its form still lives.
- **The path of an author's sample** — under their own form and without an organizer: the author is their own organizer. The samples page of `canvas` for the purpose «eval-набор» ("eval set") requires declaring that the consent form was signed and accepted by the organizer, and this purpose has no declarations of the author such as the measurement draft has. So the author signs the same form as a volunteer and accepts it themselves:
  1. **The form.** Before the first sample the author signs the volunteer's form of the text set in force (step 1 of the volunteer's path) and puts it where the volunteers' forms lie, with a record of the form (step 5). The form is the same proof of the license and of clause (a) of the DCO as a volunteer's.
  2. **The page.** The author takes the file `samples-page_v<X.Y.Z>.html` from a release of `canvas` in which the purpose «eval-набор» ("eval set") is open, and checks its signature: `just -f spec/justfile release-verify canvas v<X.Y.Z> samples-page_v<X.Y.Z>.html`, with the `SHA256SUMS` and the bundle of the same release beside it. The author opens the page from the file on the computer, in Firefox or Chrome: it does not go to the network. The page offers the tasks of the modules of the release of `modules` that its build names.
  3. **«Начало» ("Start").** The purpose is «eval-набор» ("eval set"). The notice shows the version of the text set — it must be the version of the signed form. Then come the page's two declarations: «мне 18 лет или больше, и по законам места, где я живу, я вправе заключить этот договор» and «форма согласия версии … подписана и принята организатором» ("the consent form of version … is signed and accepted by the organizer").
  4. **Writing and markup.** The task is chosen after «Начать» ("Begin"): for the eval set only the tasks of the modules, in a group per module, and for a template the instances by the seeds that the build of the page names. «Кто писал» ("Who wrote") is «взрослый своим почерком» ("an adult in their own hand") or «взрослый, подражая детскому» ("an adult imitating a child's"). After writing comes «Разметить» ("Mark up"): on a sheet, the parts of the solution with texts and flags, the first error and the final answer; on an essay, its text as written; on a workbook, each box. Nothing personal is written on the sheet: neither a name, nor a signature, nor a school, nor contacts.
  5. **Saving.** The page shows the crops — what the model will see — and asks to confirm that there is nothing personal on the sheet. The file lands in the browser's downloads, which must be a folder not synced with a cloud: the author's, `~/Downloads` on Linux, is not. The link to the organizer that the page names the author does not need: the author adds the file themselves.
  6. **The check — before any branch.** The author opens the file in the page's section «Открыть образец» ("Open a sample") and goes by the table of checks below: in the strokes — with the checkbox «показать стёртое» ("show erased"), which also shows what was erased, which lies in the file too, — in the crops and in the markup in words the author looks for a name, a signature, initials, a school and contacts, and in an essay, real names of people and a topic that is not made up (П-20). A file that did not pass is not edited: the author erases it and writes the sample anew. The rows about the fields, the schema and the PNG are held by CI.
  7. **The name and the record.** The page saves the file as `<id>.json`, for any purpose, and it is put into the module as `modules/<module>/eval/<id>.json`, where `<id>` is the field `id` of the file: `jq -r .id <file>` in the author's terminal shows it. The ID of a person's sample is 12 characters of `a–z` and `2–7`. The ID of a person's sample does not start with `synthetic-`: only synthetic samples are named so, and such a file a session of Claude Code may read — `just eval-match` rejects such a sample. Into the record of the form goes the SHA-256 digest of the file, `sha256sum <file>`.
  8. **The commit and the pull request.** A branch of the repository itself, `git commit -s` with the author's own name under clause (a) of the DCO. The sign-off is tied to the sample forever — clause (d), — and the author accepts this by signing. In the commit's message, the branch's name and the pull request's text — nothing about the sample's content. CI — `just eval-check` and `just eval-match`, — then a merge with a merge commit, as in step 7 of the volunteer's path.

  A session of Claude Code in this repository does not open the files of people's samples: `.claude/settings.json` forbids it to read `modules/*/eval/*`, except `synthetic-*.json`. The ban covers its tools for reading files and the Bash commands that name the file — `cat`, `head`, `sed` and others, — but not a program that opens files by itself, so a person's sample is in the working copy only while the author adds it. The session's work is CI's output: how many samples are in which module and whether they passed.

  An author without write access adds their samples as a volunteer.
- **Synthetic samples** are drawn by code from a seed: there is no person behind them, and they need no form. They are written by the recipe `just synthetic <module> <task>:<seed>...` in the devcontainer: it builds the module package with `module-kit build` at this commit and calls the generator of `canvas` of the exact version from the `toolchain` stage, and the generator writes `eval/synthetic-<seed>.json` with the source `modules` and the tag of the newest release of `modules`; a task that the generator does not draw the recipe names with the generator's reason. A synthetic sample, like any file that something other than the samples page of `canvas` wrote, writes `canvas_version` `0.0.0` and `text_set_version` `none`, and `sample-check` requires this of synthetic samples. CI draws each synthetic sample again with the same generator, with the tag from the file, and compares it with the whole file, except the version and the digest of the package, which the contracts text does not compare: a change of another task of the module does not touch the synthetic sample, and a change of its task is a refusal until the sample is redrawn. Which samples lie in `eval/` is the author's rule: a sample per pair "the task's primary skill × band", the task is the first by ID in the pair of those that the generator draws, the seeds are 1…N in the order of the pairs, skill, then band. `english-essays` has no synthetic samples — the generator does not draw essays. A new sample takes the module's next free seed. It is added by the author or the organizer, like their own sample: by a commit from a branch of the repository, with their own DCO sign-off. Of the table, the rows about the fields, the schema and the PNG apply to it, and the rows about the form and the names do not.
- **The checks before the merge:**

  | What | Who and when | By what |
  |---|---|---|
  | the form is accepted, and the version of the text set in the file is its version and is approved | the organizer, and for their own sample the author, before the push | the record of the form |
  | the purpose is the eval set; who wrote is an adult, in their own hand or imitating a child's, not a child; there is no writer code; the license is Apache-2.0; there is markup; the task is of its own module | the organizer, before the push, and CI on every pull request, and from a fork too: the purpose, the markup, the source `modules` with the module ID and distinct IDs are held by the eval-set mode `sample-check --eval-set` in `just eval-check`, and who wrote, the writer code and the license by the rules of the sample file for the eval-set purpose in the [contracts text](../spec/contracts/README.md). That the task is in the module's package and the digest of its open view is the one in the sample — `just eval-match` | the fields of the file |
  | on the strokes and the crops there is no name, signature, initials, school or contacts, and in the markup neither: the texts of the parts and the essay's truth are published with the sample | the organizer, and for their own sample the author, before the push, by eye | a view of the sample without the network and without AI — the samples page of `canvas`, the section «Открыть образец» ("Open a sample"), which shows the strokes of the final geometry, the crops and the texts of the markup, changing and saving nothing |
  | an essay has no real names of people, and its topic is made up (П-20) | the same | the same |
  | the file has only the protocol's fields: an unknown field is a refusal. The PNG of the crops has only IHDR, IDAT and IEND, not one ancillary chunk: neither the text chunks tEXt, iTXt and zTXt, nor the metadata eXIf, nor the time tIME, nor others. The file's name says nothing of the person | the organizer, before the push, and CI | the eval protocol; the PNG specification (sources 3.3) |
  | the file lies as `eval/<id>.json`, and `<id>` is its ID; in `eval/` there are neither folders, nor other files, nor anything starting with `.` | CI, on every pull request and push to `main`, and from a fork too | the recipe of `modules` `just eval-check`: a rule of the layout, not of the sample's validity |
  | the file is by the eval protocol, by the rules of the sample file, each file separately | CI, on every pull request and push to `main`, and from a fork too | `spec`'s recipe `just -f spec/justfile sample-check --eval-set novatrail/<folder> <files>` from `just eval-check`, without a key, with the same code as `spec`'s examples |
  | the sample's task is a task of its own module, and the digest of its open view is the one in the sample; a synthetic sample is the one that the generator of `canvas` draws on the commit's package, in everything but the version and the digest of the package | CI, on every pull request and push to `main`, without a key and without the network, and from a fork too | `just eval-match`: `module-kit open` and the generator of `canvas` from the `toolchain` stage over the commit's package |
  | the eval runner reads the samples and runs them, and takes the sample's task from the commit's package | CI, from stage 2 | the eval runner of `tutor`, with the platform's key |

  - The checks without a key CI runs on any pull request, and from a fork too.
  - By hand — before the push: a sample that got into a branch is already published, and only GitHub Support can remove it from a pull request (sources 2.5). Before the check a sample may carry a name, so neither AI nor an agent looks at it (П-12).
- **What the eval set does not accept:**
  - a sample in a pull request from a volunteer or from a fork: such a pull request is not merged and is closed;
  - a sample that failed the check — and one corrected by the organizer too: erased strokes change the sample on which the markup stands;
  - learners' work and child samples — never.

  If a name is found only after the merge, the sample leaves `eval/` with the next pull request. The organizer asks GitHub Support to remove the cached views and the pull request's references to the commit (sources 2.5) and tells the volunteer about all of it. Clones and forks cannot be removed this way — that is why the checks stand before the push.
- **Consent withdrawal.** A volunteer or an author withdraws a sample by a letter to the organizer and to the operator, by the name and the date of the form. The organizer carries out the consent withdrawal within 30 days:
  - a sample not yet added is erased;
  - an added one the organizer removes from the module's `eval/` with the next pull request, with their own sign-off and, as when adding, without the person's data in the commit, the branch and the pull request. From that commit the sample is neither in the eval set nor in the module's tests. The pull request does not change the package: `eval/` is not part of the package.

  What has already come out stays: the git history, the releases that came out and others' copies under Apache-2.0. The record of the form marks the consent withdrawal, and the form lives as long as a sample by it is in a released version.
- **The eval runner in CI** — from stage 2, only on a pull request:
  - runs the modules whose `eval/` has changed against the pull request's base: an eval set has no published version;
  - runs the modules whose package has changed against the published one — an edit of the package changes the verdict — and new modules;
  - runs all the modules when a pull request changes the tag of `tutor`: with it the prompts, the models and the runner change. A move to a new version of the format changes the packages of all modules, but the same pull request also changes the tag of `tutor`, so the runner reads the packages of the new version;
  - does not run on a push to `main`: a branch is up to date before the merge, and the same tree lands on `main`;
  - gives a report — the message `Report` of the package `novatrail.eval.v1`, without the content of the work: only the IDs of the samples, places, verdicts and marks. CI of `modules` reads it and checks it with `spec`'s recipe `report-check`, and writes no check of a report of its own; the log of the run with the models' answers is `tutor`'s own. The run and the published report follow the terms of the provider whose models the report names, and the open log of CI publishes what it prints.

  The price of a run is named before the first launch. A job with a key fails on a pull request from a fork, and the samples come from a branch of the repository anyway.
- **What it waits for** — epic E6 of [run/EPICS.md](../run/EPICS.md): the volunteers' samples, and the eval runner in CI with what it prints and for which providers.

### 3.5. The port of standalone

Olympiad maths `novatrail/olympiad-math` comes from standalone, each batch of `tools/port/batches/` at the commit its field `standalone` names: the first batch at 89986bc (sources 5.2), the rest of the port at 42b4408, which also has the tasks that are solved from the picture (sources 5.6).
- the first batch — as the first batch below says;
- the rest of the port — the other batches of `tools/port/batches/`, at 42b4408;
- the port's code — the converter `tools/port`: the batch's tasks by this table, with an answer of any kind, the check of the key and the trap values by standalone's test bench, the pictures as PNG and the overlay files, which hold what a person writes. It draws every kind of picture the tasks have, those solved from the picture among them, and their `alt` by the rule below; only `flags` is not drawn: no task has it. The recipe `just port <module>` clones standalone at the batch's commit into a temporary folder and writes its tasks, and the recipe `just port-survey <folder>` runs all the tasks that no batch names into a temporary module.
- the pictures of lessons — the data `tools/port/pictures/<module folder>/<name>.json` in standalone's schema of pictures, which the converter draws into `media/lessons/<name>.png`; the recipe `just lesson-pictures <module>`, and CI checks them with `--check`.
- the pictures of the instances of templates — the sources `tools/port/template-pictures/<module folder>/<picture path without placeholders>.json` in standalone's schema of pictures, with placeholders in the strings and, where the picture changes not only by numbers, with variants `when`, or the function of the template's solver, `{"solver": "figure"}`; the converter computes the template's instances itself — `go.starlark.net` at the version of the format's dialect — and draws the picture of each into `media/pictures/`, all of one size, which the template sets; the recipe `just template-pictures <module>`, and CI checks them with `--check`.

- **What is in standalone** at 42b4408 (sources 5.6). Its 705 tasks of 89986bc are the same there, JSON for JSON, and the 48 it added are solved from the picture: their data stand only on the picture, and the solver reads the picture's description, `solve(options, picture)`.

  | Part | Where | How many |
  |---|---|---|
  | topics and their links `builds_on` | `content/catalogs/topics.json` | 19: ten in all three bands, two in 3–4 and 5–6, seven only in 5–6. `counting.figures` and `geometry.cubes_paper` have `solved_from_picture: true` and an empty `builds_on` |
  | traps | `content/catalogs/traps.json` | 20, each with a description and advice to an adult |
  | techniques and Pólya's steps | `site/content/en/techniques.yaml`, `site/data.json` | 11 techniques in three groups |
  | tasks | `content/examples/<topic>.json` | 753: 214 in band 1–2, 267 in 3–4, 272 in 5–6. All fixed, all in English. Of the 48 solved from the picture, 30 are in the two topics above, five per band, and 18 in five others: `arithmetic.tricks` 2, `combinatorics.enumeration` 6, `counting.gaps` 4, `geometry.grid` 4, `parity.alternation` 2. Each of these 18 has the field `idea` — the place of its idea in the topic's list of ideas |
  | the tasks' solvers | `content/examples/solvers/<id>.star` | 753, one per task |
  | templates | `content/solvers/<topic>/*.star` | 39. The solver of six reads the picture: `shapes-rows`, `sticks-square`, `squares-on-lines`, `triangles-on-lines`, `fold-net` and `fold-sheet` |
  | pictures | the task's field `picture` — a JSON description that standalone's card draws; the kinds — `content/pictures/` | in 152 tasks, 15 kinds of 17; a clock in 36, `lattice` in 38. A solution picture in 83 |
  | hints | the task's field `hint` | in 311 tasks |
  | the site's topic pages | `site/content/<language>/topics/` | 19: the 17 of 89986bc, in each of 10 languages, `en` and `ru` among them (sources 5.2), and `counting-figures` and `cubes-and-paper` |
  | ideas | `content/ideas/<topic>.json` | 634 |
  | excluded skills | `content/catalogs/skills.json` | 25 |

  Every task's answer is a choice of five options. Each of the four wrong options has a trap and an explanation. By the text of the correct option, the answers are these:
  - an integer — 544, all 48 tasks solved from the picture among them. Of them at least 17 have a question that points to the options: "Which of these numbers…";
  - a number with a unit — 12: "8 km", "16 cents";
  - other text with digits — 44: dates, "9 hours 30 minutes", "bag 3", "5 with 12 left over";
  - a fraction — 7;
  - a percentage — 5;
  - a time on a clock — 30;
  - "yes" or "no" — 1;
  - words — 110.

  One explanation of a task solved from the picture names the key — that of `gaps-12-d2-6`, option 6: "The steps were not divided by 3." The three in this task is a jump on the picture, not a number of the statement.
- **The `alt` of a task solved from the picture.** A task is solved from the picture if its standalone solver declares `solve(options, picture)`. Its data stand only on the picture, and a screen reader reads the `alt`. So:
  - the `alt` says in words every fact that the solver reads: the size of the field, every line, wall and cell, every number and label, every stair and jump;
  - the `alt` names neither the key nor a step of the solution;
  - the converter checks that the `alt` names every text that the picture writes — a label, a number on a tile, the total of a row of figures — and refuses if not. A label of capital letters it looks for case-sensitively. The geometry is held by proofreading;
  - for any task, the `alt` of the statement's picture does not name the key unless the question names it or the picture itself writes it: otherwise the leak check would count the key as shown and would not look for it anywhere. This is a check of the converter too. If the picture writes the key, `leak_phrases` carry its giving away.
- **What becomes what**, in the module's files by the [format text](../spec/format/README.md):

  | In standalone | In the module |
  |---|---|
  | topic | a skill of `skill-graph.yaml` with the assessability `verdict`, as in the format example `olympiad-math`. The ID is the topic with `.` and `_` replaced by `-`: `counting.gaps` → `counting-gaps` |
  | `builds_on` | the skill's `prerequisites`. As in the example, `number-divisibility` also has the anchor skill `core.division-with-remainder`, and `fractions-parts` — `core.fractions-of-a-whole`. The tiers are computed from the graph, and their names are Foundations, Rules and Hard problems |
  | the topics' grades | `audience.grades` of the module manifest, 1–6, and the bands `bands`: `grades-1-2`, `grades-3-4` and `grades-5-6` |
  | trap | an entry of `traps.yaml`, as in the format example, `twelve-hours-off` among them — the value of a 12-hour dial for the kind `time`. The ID — with `_` replaced by `-`; `advice` → `trap_advice` |
  | a technique of the site | an entry of `techniques.yaml`, as in the format example: `title` and `description`. The techniques' numbers and groups are not carried over — the format has none |
  | Pólya's steps | `solving-steps.yaml`, as in the format example, with the descriptions of two steps of the four. The example takes them from the coach of standalone's prototype, and standalone's page of techniques opens with the same four steps from Pólya's book "How to Solve It" |
  | task ID `gaps-12-d1-1` | the same ID and the file `tasks/gaps-12-d1-1.yaml` |
  | `topic`, `grade_level`, `difficulty` | `skills.primary` — the topic's skill; `band` — the band; `difficulty` — the same, 1–5 within the band. In `supporting` — an anchor skill, if the task is on fractions or divisibility |
  | `question` | `statement` — the text as it is. It has no formulas, and the format lets the signs × and ÷ and the fractions `a/b` through |
  | the correct option is an integer, a number with a unit or a fraction | one item with an answer slot `anywhere` and the check `number`, for a number with a unit — with the batch's `units`, and for a fraction — with the check `fraction`. The answer is written by hand. A percentage is a number with the unit `%`, "44%": the check `number` of the [format text](../spec/format/README.md) drops it after the number, with a space and without. The exception is a question that points to the options: without them it is not posed, and the task gets `choice` by tapping |
  | the correct option is a time on a clock `H:MM` | an answer slot `anywhere` with the check `time`, without `options`: the learner writes any time, the hours and minutes are compared as numbers, and each wrong option is a trap value. A key after 12:59 and before 1:00 also has the value `twelve-hours-off`: the same time on the dial, 9:10 for 21:10 and 12:15 for 0:15. There are 30 such tasks, all `clk-*`, and five of them have the dial's value |
  | the correct option is one word, other than "yes" and "no", that stands in the question | an answer slot `anywhere` with the check `choice`, which has its own `options`, all five options, and `notation: word`: the learner writes the answer by hand and does not see the options, any one word is read, a word not among the five options is `not_yet` with a hint rung, and a full stop at the end is dropped. There are 44 such choices — 34 in `ord-*`, 7 in `par-*` and 3 days of the week in `cal-*`. The paragraph `answer_form` of the overlay may say in the statement how to write the answer. If the wrong options are arguments, not answers, the batch sets `tapped: true` on the task, and the choice is by tapping |
  | the correct option is something else: a word that is not in the question, a date, "yes" or "no", a quantity of several parts, a phrase, a number with words for which the batch gives no `units` | an answer slot `choice` with five options, by tapping: without the options such an answer cannot be written, and the checks of format v0 do not judge dates, quantities of several parts and text |
  | a wrong option with a trap and an explanation | a trap value `trap_values`, where the trap is with `-` instead of `_`, and the text is the option's explanation. The value is written like the key: for `number` and `fraction` — the option's number in the kind's notation, without a unit, no longer than 16 characters, and for `choice` — the option's text as it is. A non-numeric wrong option of a task with a numeric answer cannot become a trap value, and its trap goes into the task's `traps` — AI will suggest it |
  | an explanation that names the correct answer | is rewritten in `trap_texts` of the overlay. The key of a number and of a fraction is looked for by `task.no-leaks`; a key that also stands in the statement it does not look for. Among the 705 tasks of 89986bc, with an integer key there are 159 such explanations, in 119 tasks, and with a fraction or a number with a unit 5 more; among the 48 solved from the picture, one. The key of a tapped choice the check never looks for, and a word that is written by hand stands in the question, so the explanations of a choice — and `givens`, `unknown` and `simpler_case` of the overlay — the converter checks itself: it finds the key's phrase in the normal form of the leak check of the [format text](../spec/format/README.md), and refuses |
  | `solution`, `solution_total` | `solution`, and the total as its last line |
  | `picture`, `solution_picture` | a PNG in `media/`, drawn from the description: with `alt`, width and height, on an opaque light background. The solution picture — in `solution` |
  | the task's solver | does not go into the module: the key of a fixed task is the author's. The port runs the solver once and checks the key and the trap values against it |
  | template | a task `tasks/<template name>.yaml` with `template` — parameters, derived numbers, a constraint and a solver — written by hand. The parameters are the numbers of the reference task's statement, and its numbers are one of the instances; the solver is rewritten for the format's dialect: without uppercase constants, without the helpers `permutations`, `product` and `gcd`, without `sum` and `**` and without recursion. A template whose solver reads the picture is carried over only if the picture can be built from the numbers it shows: the solver's function `figure` builds it from the parameters, `solve` computes on it, and the port draws the same one — so `shapes-rows`, `squares-on-lines` and `sticks-square`; `fold-net`, `fold-sheet` and `triangles-on-lines` stay tasks. The reference task's picture is the picture of each instance: a path with placeholders of the parameters, a file for each instance, which `just template-pictures` draws, for all the instances |
  | the site's topic page | a lesson `lessons/<id>.yaml`. The gist of the topic and the method become an explanation, the examples worked examples, and their steps parts of the demo. An example that solves a task of the module does not go into the lesson. The language of the lesson is the module's language, and the page is taken in it |

- **What is not carried over:**

  | In standalone | What becomes of it |
  |---|---|
  | the task's hint, `hint` | the leading question `leading_question` at the root of the task, after `unknown`, word for word: `api` gives it on "Hint" before the first wrong answer, and `task.no-leaks` looks for keys in it. 311 tasks have a hint — 22 in band 1–2, 17 in 3–4 and 272 in 5–6 — and 53 of them are not questions but instructions, which also break the task down. The key of a choice and of a time in it the converter looks for, and the key of a number and of a fraction `task.no-leaks`, except a key that the statement shows: a question with it shows no more than the statement. Where the question names the key, the field `leading_question` of the overlay rewrites it. A template whose reference task carries a hint keeps it with placeholders |
  | the techniques' numbers and groups | the format has none; the group is seen from the technique's description |
  | the kinds of pictures, 17 | a picture is a file, so the kinds are not needed |
  | ideas, 634 | there is no generation on the fly |
  | the task's field `idea` and the topic's `solved_from_picture` | the first is the place of the idea in its list, and the ideas are not carried over. The second tells standalone in which chat the task may be given: only where the card is drawn. In a module the picture is part of the statement, and every app shows it |
  | the excluded skills, 25, the limits of readability and the excluded genres | these are rules for whoever writes tasks. standalone's texts already follow them, and new tasks need them in the guide for authors |
  | the page's parts `home` and `ask` | they are for a parent, not for a lesson |
  | the site's groups of topics, 6 | the graph has none |
  | the solvers of the site's examples | they prove the examples' answers; in a lesson an example carries the author's answer |
- **The license.** standalone's content is under MIT (sources 5.1), and the module carries its notice in `NOTICE`. `NOTICE` follows the format's package vector `with-solver`:
  - a line saying that the module holds standalone's content, and the text of MIT with "Copyright (c) 2026 MathTrail";
  - a line saying that part of the tasks and Pólya's steps came from standalone's prototype under Apache-2.0 and are changed here. Apache-2.0, whose text is the repository's `LICENSE`, requires it: "You must cause any modified files to carry prominent notices stating that You changed the files". The copy of Apache-2.0 is carried by the module's `LICENSE`, as in every module.

  Neither standalone nor the prototype names the copyright holder, except for the word "MathTrail", though they have one author on GitHub. So `NOTICE` meets the terms of both licenses, and does not rest on there being one author.
- **Provenance** — by sources 5.2 and 5.7:
  - 442 tasks of grades 1–4 standalone took word for word from its prototype. The prototype is under Apache-2.0, without `NOTICE` and without a filled-in copyright line;
  - the wording of the other tasks is its own. The ideas for grades 5–6 are taken from MathNet under CC BY 4.0 and from books and collections, but the wording is not. standalone's tasks have no field "source";
  - a task of the module has no field of provenance either. For О-10 `spec` recorded: others' sets for this module were checked by the records of standalone and the prototype, there was no comparison of the wording with the sources, and the module does not wait for О-10. For the lawyer at О-10 — attribution for an idea or a kind of task, and wording that a model wrote. The sources of the ideas can be named voluntarily only outside the module: `NOTICE` holds only what the licenses require (О-10).

### 3.6. The first batch

The first batch of `novatrail/olympiad-math` is 27 tasks of standalone across three bands and 9 lessons, for the first application, where the author solves tasks locally in a browser. The tasks and lessons follow the table of the port above and this section; the batch is `tools/port/batches/first.yaml`. The list is counted by standalone 89986bc (sources 5.5).

- **Skills** — `counting-gaps`, `combinatorics-enumeration` and `pigeonhole-basic`. These three topics of standalone run in all three bands, their answers are integers almost everywhere, and the site's topic page has an example for each level. Tasks with a number answer, without "Which of these…" and without a picture, they have 52, 64 and 63 of 65.
- **The rule of choice:**
  1. a fitting task has an integer as its correct option, with a unit or without, like "50 m", and its question does not point to the options. It has no picture, except a picture `row` under point 5. An answer with a unit is written the same as without it: the key is a number, and the unit is in `units`, as in the port above;
  2. for each skill and band three tasks are taken — of difficulty 1, 3 and 5;
  3. at each difficulty the first fitting task by ID is taken whose explanations of the wrong options do not name the key, and if there is no such one — the first fitting one. A task is not taken if its key:
     - stands in the text of its lesson — `lesson.no-leaks`;
     - stands in the statement of an earlier task of the lesson;
     - equals the key of an earlier task of the lesson: that task's solution ends with its key, and a lesson must not name tasks the texts of one of which give the answer of another.

     If there is no fitting task at this difficulty, the nearest difficulty is taken, higher first;
  4. the key is looked for in a text by the format's rule of finding numbers: by value, so `7` is in "7 rungs" but not in "17";
  5. for `counting-gaps` in bands 1–2 and 3–4 one task of three has a picture `row`: at the first difficulty where such a one fits, it is taken instead of a task without a picture. There are 10 fitting tasks with `row`, all in these two bands.

  The difficulty shifted for one task: for `counting-gaps` 1–2 the keys of the fitting tasks of difficulty 1 stand in the lesson's text, and `gaps-12-d2-1` was taken.

  In three lessons a number from the solution or from an option's explanation of an earlier task equals the key of a later one:
  - `enum-56-d1-1` → `enum-56-d3-2`, 16 — in the solution and in the explanations of options 36 and 576;
  - `pig-12-d3-1` → `pig-12-d5-1`, 10 — in the solution and in the explanation of option 10;
  - `pig-56-d1-1` → `pig-56-d3-2`, 19 — in the solution and in the explanation of option 19;
  - `pig-56-d3-2` → `pig-56-d5-1`, 11 — in the explanation of option 11.

  This is a coincidence of numbers, not an answer, and the rule does not look at these texts: otherwise for `combinatorics-enumeration` 5–6 both tasks would have moved to difficulty 4. An explanation of such a pair does not read as the answer of the next task; where it would, its number is changed.

  The two paragraphs above are about the first batch. For the lessons of the rest of the port the rule is the same without points 1 and 5: any task of the skill and band fits, except a template, with an answer of any kind and with any picture. A number key is looked for by value, by point 4, and a key with a word in it — as a phrase, case-insensitively, without quotes and a sign at the end, and with word boundaries. A task whose solution gives the answer of a later task of the lesson not by a number but by the course of solution the rule skips at the executor's word. If no fitting task remains at a later difficulty, the rule takes at an earlier one the next in its order, and the lesson's tasks go by difficulty. The lessons `counting-figures` and `geometry-cubes-paper` rewrite the cards and examples of the pages that solved tasks of the module, and show the pictures `media/lessons/` that the converter draws from `tools/port/pictures/<module folder>/`. The preference of point 3 changes nothing there: the port rewrites the explanations that name the key. The difficulty shifts and the pairs with coinciding numbers of each such lesson are recorded in «Сделано» of its task, in the epic archives `docs/plan/`.
- **The graph** — three skills with the assessability `verdict`, as in the format example `olympiad-math`:

  | Skill | Prerequisites | Tier |
  |---|---|---|
  | `counting-gaps` | — | 1, Foundations |
  | `combinatorics-enumeration` | — | 1, Foundations |
  | `pigeonhole-basic` | `combinatorics-enumeration` | 2, Rules |

  The tier names are `[Foundations, Rules, Hard problems]`, as in the example. The graph is not a chain, and the module is not sequential. There are no anchor skills among the three. The rest of the port adds the other 16 skills, and the module's graph is its `skill-graph.yaml`.
- **The module manifest** — like the format example `olympiad-math.yaml`, except for `audience.age_categories` — `[under_13, adult]` — and the version, which the module's `manifest.yaml` sets. The grades are 1–6, the bands `grades-1-2`, `grades-3-4` and `grades-5-6`, the explanation language `en`, `rung_templates` from the example, and there are no numbers in them.
- **The catalogs** — from the format examples:
  - `traps.yaml` — the batch's tasks name 10 traps of it: `answered-other-question`, `best-case-not-worst`, `double-count`, `ignored-condition`, `ignored-not`, `missed-case`, `number-from-text`, `off-by-one`, `stopped-early`, `wrong-operation`. None of them names skills, and there are no numbers in their advice;
  - `techniques.yaml` — the batch's tasks name 4 techniques of it: `draw`, `small-numbers`, `systematic`, `pigeonhole`;
  - `solving-steps.yaml` — all four steps of the example.
- **A task of the batch** — fixed, with the same ID and the file `tasks/<id>.yaml`:
  - one item: an answer slot `anywhere`, the check `number`, the key as a string and `units`. The question of every task of the batch counts objects or a measure, so every one has `units`: the word that the question counts, in the singular and the plural, and for a measure also the abbreviation. Without them "14 metres" by hand would not be in the number's notation;
  - the wrong options are all numbers, and all become trap values `trap_values`: the value is the option's number without a unit, and the text is the option's explanation. `traps` from the port is empty;
  - techniques — by skill: those that standalone's site links with the topic (`site/data.json`, the member `techniques`), and for `counting-gaps` also `draw` — the method of its page "Draw a small case and count on it" and the techniques of the task `gaps-ladder-1` of the format example. In total: `counting-gaps` — `draw` and `small-numbers`, `combinatorics-enumeration` — `systematic`, `pigeonhole-basic` — `pigeonhole`;
  - every task has the item's `givens`, `unknown` and `simpler_case`. Without them the hint before the first wrong answer and rung 3 are left without words. The overlay holds them;
  - standalone's hints of nine tasks of band 5–6 are questions that break the task down, that is, a path to the answer: they are leading questions `leading_question`, given on "Hint" before the first error;
  - `solution` — standalone's solution, the total as the last line. The port checks the key and the trap values against the task's solver;
  - the key stands in no statement, so `task.no-leaks` looks for it in all the task's texts.

- **The list of tasks** — by skills and bands, in the order of the lessons. A trap value is a wrong option and its trap:

  | ID | Skill | Band | Difficulty | Key | `units` | Trap values | Special |
  |---|---|---|---|---|---|---|---|
  | `gaps-12-d2-1` | `counting-gaps` | `grades-1-2` | 2 | 14 | metre, metres, m | 7 `stopped-early`, 8 `number-from-text`, 10 `wrong-operation`, 16 `off-by-one` | picture `row`; difficulty 2 instead of 1 |
  | `gaps-12-d3-1` | `counting-gaps` | `grades-1-2` | 3 | 10 | gap, gaps | 5 `stopped-early`, 6 `number-from-text`, 11 `ignored-condition`, 12 `off-by-one` | — |
  | `gaps-12-d5-5` | `counting-gaps` | `grades-1-2` | 5 | 50 | step, steps | 26 `wrong-operation`, 30 `stopped-early`, 40 `wrong-operation`, 60 `off-by-one` | — |
  | `gaps-34-d1-1` | `counting-gaps` | `grades-3-4` | 1 | 100 | metre, metres, m | 21 `number-from-text`, 26 `wrong-operation`, 95 `off-by-one`, 105 `off-by-one` | picture `row` |
  | `gaps-34-d3-1` | `counting-gaps` | `grades-3-4` | 3 | 90 | second, seconds, s | 40 `wrong-operation`, 75 `off-by-one`, 100 `off-by-one`, 300 `wrong-operation` | — |
  | `gaps-34-d5-1` | `counting-gaps` | `grades-3-4` | 5 | 33 | minute, minutes, min | 25 `ignored-condition`, 35 `off-by-one`, 37 `off-by-one`, 42 `off-by-one` | — |
  | `gaps-56-d1-2` | `counting-gaps` | `grades-5-6` | 1 | 50 | metre, metres, m | 55 `off-by-one`, 105 `ignored-condition`, 45 `off-by-one`, 10 `stopped-early` | standalone's answer "50 m" |
  | `gaps-56-d3-1` | `counting-gaps` | `grades-5-6` | 3 | 39 | tear, tears | 40 `off-by-one`, 28 `wrong-operation`, 13 `wrong-operation`, 11 `ignored-condition` | — |
  | `gaps-56-d5-1` | `counting-gaps` | `grades-5-6` | 5 | 120 | part, parts | 130 `off-by-one`, 108 `off-by-one`, 21 `wrong-operation`, 240 `ignored-condition` | — |
  | `enum-12-d1-3` | `combinatorics-enumeration` | `grades-1-2` | 1 | 4 | number, numbers | 2 `ignored-condition`, 3 `missed-case`, 5 `off-by-one`, 6 `double-count` | the explanation of option 5 is rewritten |
  | `enum-12-d3-1` | `combinatorics-enumeration` | `grades-1-2` | 3 | 10 | call, calls | 5 `number-from-text`, 9 `missed-case`, 20 `double-count`, 25 `wrong-operation` | — |
  | `enum-12-d5-2` | `combinatorics-enumeration` | `grades-1-2` | 5 | 8 | outfit, outfits | 5 `wrong-operation`, 6 `wrong-operation`, 7 `off-by-one`, 9 `ignored-not` | — |
  | `enum-34-d1-1` | `combinatorics-enumeration` | `grades-3-4` | 1 | 15 | handshake, handshakes | 6 `number-from-text`, 12 `wrong-operation`, 30 `double-count`, 36 `wrong-operation` | — |
  | `enum-34-d3-1` | `combinatorics-enumeration` | `grades-3-4` | 3 | 45 | number, numbers | 36 `missed-case`, 44 `missed-case`, 54 `ignored-condition`, 90 `answered-other-question` | — |
  | `enum-34-d5-2` | `combinatorics-enumeration` | `grades-3-4` | 5 | 105 | way, ways | 21 `stopped-early`, 35 `stopped-early`, 70 `wrong-operation`, 210 `double-count` | — |
  | `enum-56-d1-1` | `combinatorics-enumeration` | `grades-5-6` | 1 | 52 | swap, swaps | 36 `missed-case`, 100 `ignored-condition`, 576 `wrong-operation`, 20 `wrong-operation` | — |
  | `enum-56-d3-2` | `combinatorics-enumeration` | `grades-5-6` | 3 | 16 | way, ways | 24 `ignored-condition`, 54 `ignored-condition`, 8 `wrong-operation`, 81 `ignored-condition` | — |
  | `enum-56-d5-2` | `combinatorics-enumeration` | `grades-5-6` | 5 | 30 | triangle, triangles | 18 `missed-case`, 35 `ignored-condition`, 60 `double-count`, 31 `missed-case` | — |
  | `pig-12-d1-2` | `pigeonhole-basic` | `grades-1-2` | 1 | 4 | apple, apples | 1 `best-case-not-worst`, 2 `best-case-not-worst`, 3 `stopped-early`, 6 `wrong-operation` | — |
  | `pig-12-d3-1` | `pigeonhole-basic` | `grades-1-2` | 3 | 11 | ball, balls | 3 `best-case-not-worst`, 6 `wrong-operation`, 10 `stopped-early`, 15 `ignored-condition` | — |
  | `pig-12-d5-1` | `pigeonhole-basic` | `grades-1-2` | 5 | 10 | ball, balls | 3 `best-case-not-worst`, 6 `wrong-operation`, 9 `stopped-early`, 12 `ignored-condition` | — |
  | `pig-34-d1-2` | `pigeonhole-basic` | `grades-3-4` | 1 | 3 | pupil, pupils | 2 `best-case-not-worst`, 12 `number-from-text`, 13 `wrong-operation`, 25 `number-from-text` | — |
  | `pig-34-d3-1` | `pigeonhole-basic` | `grades-3-4` | 3 | 7 | ball, balls | 3 `best-case-not-worst`, 6 `stopped-early`, 9 `wrong-operation`, 13 `answered-other-question` | — |
  | `pig-34-d5-1` | `pigeonhole-basic` | `grades-3-4` | 5 | 17 | ball, balls | 5 `best-case-not-worst`, 16 `stopped-early`, 20 `wrong-operation`, 21 `ignored-condition` | — |
  | `pig-56-d1-1` | `pigeonhole-basic` | `grades-5-6` | 1 | 20 | person, people | 19 `off-by-one`, 21 `ignored-condition`, 1 `best-case-not-worst`, 610 `answered-other-question` | — |
  | `pig-56-d3-2` | `pigeonhole-basic` | `grades-5-6` | 3 | 19 | pupil, pupils | 18 `off-by-one`, 11 `answered-other-question`, 2 `best-case-not-worst`, 20 `ignored-condition` | — |
  | `pig-56-d5-1` | `pigeonhole-basic` | `grades-5-6` | 5 | 11 | girl, girls | 10 `missed-case`, 21 `number-from-text`, 20 `ignored-condition`, 12 `ignored-condition` | the explanation of option 10 is rewritten |

- **Lessons** — a lesson per skill and band, ID `<skill>-grades-<band>`, and the lesson's skill is its skill. The title is the title of the topic page and the band's words, without numbers: "Gaps and boundaries: first steps" in band 1–2, "…: going further" in 3–4 and "…: olympiad level" in 5–6. The band's title in the module manifest, "Grades 1–2", does not go into the lesson's title: `lesson.no-leaks` checks the numbers in it against the keys of the lesson's tasks. The segments in order:
  1. an explanation: the page's `hero.lead` and cards `basis` — the gist of the topic and its rules;
  2. a worked example of its own level: the statement — `question`, a handwriting sheet with a grid, the demo — `steps` and `answer`;
  3. an explanation: the example's note, `note`, if it has one;
  4. the lesson's three tasks in increasing difficulty;
  5. an explanation: the paragraph "How to solve: …" — the page's `method.title` — and its steps `method.steps` as a part `numbered_list`.

  There are no pictures in the lessons. No example solves a task of standalone: the tasks nearest in wording ask something else and have another answer. The lessons of `pigeonhole-basic` come after the lessons of `combinatorics-enumeration` by the graph. The format puts every lesson of `pigeonhole-basic` after every lesson of `combinatorics-enumeration`, after the lesson of 5–6 too: a lesson has no band. A learner of band 1–2 needs the lesson of `pigeonhole-basic` after the lesson of `combinatorics-enumeration` of their own band, not after the lesson of 5–6. Which of a skill's three lessons to give a learner, how to read the order of lessons of different bands and how the lessons of `counting-gaps` stand beside the lessons of `combinatorics-enumeration` is decided by the engine's rule. By the engine's rule, a lesson's band is the set of the bands of its tasks, a lesson fits a skill when the band of its recommended task is in this set, and a lesson of another band holds no order.

  | Lesson | Skill | The page's example | The example's note | Tasks in order |
  |---|---|---|---|---|
  | `counting-gaps-grades-1-2` | `counting-gaps` | `gaps-and-boundaries`: "The bus route", `bus-stops` — 15 minutes | "Why not 18?" | `gaps-12-d2-1`, `gaps-12-d3-1`, `gaps-12-d5-5` |
  | `counting-gaps-grades-3-4` | `counting-gaps` | `gaps-and-boundaries`: "Floor 9 and floor 3", `floor-nine` — 4 times as many | "Why not 3 times?" | `gaps-34-d1-1`, `gaps-34-d3-1`, `gaps-34-d5-1` |
  | `counting-gaps-grades-5-6` | `counting-gaps` | `gaps-and-boundaries`: "How many logs?", `woodcutters` — 5 logs | "Does it matter how the cuts were shared?" | `gaps-56-d1-2`, `gaps-56-d3-1`, `gaps-56-d5-1` |
  | `combinatorics-enumeration-grades-1-2` | `combinatorics-enumeration` | `enumeration`: "Flags of two stripes", `two-stripes` — 6 flags | — | `enum-12-d1-3`, `enum-12-d3-1`, `enum-12-d5-2` |
  | `combinatorics-enumeration-grades-3-4` | `combinatorics-enumeration` | `enumeration`: "Digits three apart", `three-apart` — 13 numbers | "The case most often lost" | `enum-34-d1-1`, `enum-34-d3-1`, `enum-34-d5-2` |
  | `combinatorics-enumeration-grades-5-6` | `combinatorics-enumeration` | `enumeration`: "No two 2s side by side", `no-two-twos` — 8 numbers | "A count that grows" | `enum-56-d1-1`, `enum-56-d3-2`, `enum-56-d5-2` |
  | `pigeonhole-basic-grades-1-2` | `pigeonhole-basic` | `pigeonhole-principle`: "Sweets in a jar", `lemon-sweets` — 9 sweets | — | `pig-12-d1-2`, `pig-12-d3-1`, `pig-12-d5-1` |
  | `pigeonhole-basic-grades-3-4` | `pigeonhole-basic` | `pigeonhole-principle`: "Two black socks", `two-black-socks` — 10 socks | "Why not 9?" | `pig-34-d1-2`, `pig-34-d3-1`, `pig-34-d5-1` |
  | `pigeonhole-basic-grades-5-6` | `pigeonhole-basic` | `pigeonhole-principle`: "Sums in a grid", `grid-sums` — Yes, always | "The holes are not always boxes" | `pig-56-d1-1`, `pig-56-d3-2`, `pig-56-d5-1` |

- **Pictures** — `row` for `gaps-12-d2-1` and `gaps-34-d1-1`: a row of posts with a gap and a label of the interval, 2 and 5 metres. This is a PNG in `media/`, drawn by the converter from standalone's description — with `alt`, width and height, on a light opaque background. They have no solution pictures.
- **Language** — English, `languages.explanation: en`. Such is the language of standalone's tasks, of the pages `en` and of the format examples. The same module in Russian would be another module, with its own name. The module's language is named in spec's open question О-2.

### 3.7. What to build

The rest of this part describes how the component is built; the table is what of it is not built yet, by the epics of table 2, [run/EPICS.md](../run/EPICS.md), which give its order. What waits for the format or for content, not for an epic, is in what is to decide when detailed, above. The first application waits for nothing more of `modules`.

| What | Where | Epic | What it waits for |
|---|---|---|---|
| volunteers' samples and the eval runner of `tutor` in CI, `report-check` on its reports; the samples of templates, of a choice by hand and of essays, and people's samples of the module of the check board | the samples in the eval set; the module of the check board | E6 | stage 2, with the eval runner of `tutor` in a release — the releases of `tutor` carry no binaries `tutor-eval` yet —, and volunteers |

### 3.8. The set of school maths for grade 3

The set `novatrail/grade-3-math` is a small module on a workbook, by which `modules` tries the format on a second canvas.

- **The manifest, the graph and the catalogs** are those of the format examples `school-math-grade-3` word for word, but for the ID, `novatrail/grade-3-math`. The audience is `under_13`, grade 3, one band `grade-3`, the explanation language is `en`, without `dependencies`.
- **Skills** — the example's four skills, closed under prerequisites:

  | Skill | Prerequisite | Tier | Teaches the anchor skill |
  |---|---|---|---|
  | `times-tables` | — | 1 | `core.multiplication-table` |
  | `division-with-remainder` | `times-tables` | 2 | `core.division-with-remainder` |
  | `fractions-of-a-whole` | `times-tables` | 2 | `core.fractions-of-a-whole` |
  | `comparing-fractions` | `fractions-of-a-whole` | 3 | — |

  The graph is not a chain, so the module is not sequential.
- **Traps** — 11:
  - three from the example word for word;
  - `wrong-operation` from the example, but with other words: lines without a story catch it too;
  - `reversed-comparison` from the example, but with the skill `comparing-fractions`: the example's skill `numbers-to-1000` is not in the set (`traps.known-skills`);
  - six of its own — `took-away`, `remainder-left-out`, `remainder-not-added`, `top-and-bottom-swapped`, `one-part-only` and `greater-bottom-greater-fraction`.

  **Techniques** — three from the example and one of its own, `draw-the-parts`. **Solving steps** — the example's two steps.
- **Tasks** — 12 pages `tasks/<skill>-d<difficulty>.yaml`, difficulties 1, 3 and 5, 6–10 lines each:

  | Skill | Lines | Check kinds |
  |---|---|---|
  | `times-tables` | a box anywhere in the line: a factor, the dividend, the product, the quotient; at d5, stories with `units` | `number` |
  | `division-with-remainder` | "47 ÷ 6 = [ ] r [ ]" — two answer slots in a line; at d5, stories with two boxes, a remainder that needs room, and a missing dividend | `number` |
  | `fractions-of-a-whole` | a fraction by equal parts, `form: any`; a fraction of a number; an irreducible fraction, `form: reduced` | `fraction`, `number` |
  | `comparing-fractions` | a sign between two fractions with the same denominator or the same numerator; at d5, the greater or the smaller of two fractions | `choice` by hand — `<`, `>` and `=` —, `fraction` |

  - The page's statement names the form of the answer: "in lowest terms", "Write <, > or =".
  - The solution is a paragraph per line. The lines have trap values with texts without numbers, and one error is one trap value: in a division line only the quotient's answer slot carries it.
  - `simpler_case` only in the stories, on smaller numbers. Neither what is written in it nor what follows from it is the key of its line, and the fractions do not name the page's keys in words either.
  - A line does not repeat a fact from the text of its page or from the solution of an earlier page of the same lesson.
  - There are no `givens` and `unknown`: the solving steps do not give them.
  - Three pages per skill are enough for the corridor of `engine` too: more pages do not lift it (measurements 2 and 5), and every other set of pages measured made the gate of `modules` red (measurements 5).
  - There are no pictures, `predicate` or templates.
- **Lessons** — a lesson per skill, `lessons/<skill>.yaml`: an explanation, a worked example on the workbook, then three pages by increasing difficulty. The texts of the lessons, traps and techniques have no digits: the pages' keys are small numbers, and the leak check searches for them in all these texts. There are no signs of comparison there either, though the check does not search for them: the page's statement names all three signs, and a key that stands in a text seen together with the item is not searched for. After the explanation comes a worked example on the workbook: two or three lines with boxes that are not among the lines of the lesson's pages, and a demo that fills each box; the statement and the texts of the demo are without digits.
- **The page fits.** The rule `task.workbook-page` holds each page to the format's layout `page.yaml`, and each of the 12 pages takes from 3 to 7 rows by the layout of `canvas` (sources 4.19).
- **`LICENSE`** — the text of Apache-2.0. There is no `NOTICE`: the texts are its own.

### 3.9. The set of short essays

The set `novatrail/english-essays` is a small module of writing on a handwriting sheet, by which `modules` tries the format on a second subject and on a path without a key — feedback by a rubric.

- **The manifest, the graph, the techniques and the solving steps** are those of the format examples `short-essays` word for word, but for the ID, `novatrail/english-essays`; so are the traps, and one of its own, below. The audience is all four age categories, without grades; the bands are `first-steps` and `whole-essays`, also without grades; the explanation language is `en`, without `dependencies`. There is an item with a verdict, so the manifest has the example's `rung_templates`.
- **Skills** — the example's six skills:

  | Skill | Prerequisite | Tier | How it is assessed |
  |---|---|---|---|
  | `reading-model-essays` | — | 1 | `none`: only a supporting skill and a lesson of explanations alone |
  | `choosing-a-first-sentence` | — | 1 | `verdict` |
  | `writing-a-first-sentence` | `choosing-a-first-sentence` | 2 | `rubric` |
  | `paragraphs` | `writing-a-first-sentence` | 3 | `rubric` |
  | `endings` | `paragraphs` | 4 | `rubric` |
  | `improving-an-essay` | `paragraphs` | 4 | `rubric` |

  The graph is not a chain: tier 1 has two skills, and `endings` and `improving-an-essay` both wait for `paragraphs`. So the module is not sequential, and where the graph does not order the lessons, the engine's rule chooses, as the lessons below say.
- **Traps** — three: `flat-opening` and `announcing-the-topic` from the example word for word, and one of its own, `giving-it-all-away`. Only the wrong options of the choice tasks catch them: a rubric has no verdict. **Techniques** — the example's two, `start-in-the-middle` and `read-it-aloud`. **Solving steps** — the example's three steps, without events.
- **Tasks** — 8, all on a lined sheet:

  | Skill | Tasks | Check kind | Band, difficulty |
  |---|---|---|---|
  | `choosing-a-first-sentence` | `best-first-sentence-1`, `-2` and `-3`: which of three first sentences invites one to read on | `choice` by tapping, `leak_phrases` | `first-steps`, 1, 3 and 5 |
  | `writing-a-first-sentence` | `talking-umbrella-1`: three first sentences about a talking umbrella, one of them circled | `rubric` | `first-steps`, 2 |
  | `paragraphs` | `weather-shop-1`: a day in a weather shop, a paragraph per part of the day | `rubric` | `whole-essays`, 1 |
  | `paragraphs`, `endings` | `undiscovered-island-1`: the first morning on an island that nobody has found | `rubric` | `whole-essays`, 2 |
  | `improving-an-essay` | `dragon-draft-1`: a draft with repetitions from the statement — anew and clearer | `rubric` | `whole-essays`, 2 |
  | `paragraphs`, `endings` | `a-day-you-remember-1`: a day of one's own life — `adult` only | `rubric` | `whole-essays`, 3 |

  - A choice task has `solution`, `leak_phrases` with the key, and a trap value on each wrong option. The trap value has a text everywhere but `announcing-the-topic` in `best-first-sentence-1`, the option written word for word as in the example: after it the learner reads the trap advice.
  - An essay has `binding: sheet` and a rubric of three criteria on structure and content, without capital letters, spelling and punctuation.
  - No statement holds the reminder "Do not write your name or the real name of anyone you know.": no reminder stands before "Check" of work that goes to AI whole either, and `app` reminds the learner on the camera of a task from paper and in the adult's Help.
  - The topics for those under 18 are made up and general, without family, friends, home and school; a topic about the learner's life is only in the task for adults (П-20).
  - There are no pictures.
- **Lessons** — six, `lessons/<id>.yaml`:
  - `reading-model-essays` — explanations alone, from the example;
  - `first-sentences` — the example's explanations, the three choice tasks and `talking-umbrella-1`;
  - `paragraphs` and `endings` — an explanation and an essay;
  - `improving-an-essay` — an explanation, a worked example on a lined sheet that corrects its own draft, and `dragon-draft-1`;
  - `writing-about-a-day` — the format example `essays-for-adults.yaml` word for word, `adult` only, with the skill `endings` alone: a lesson that must come after the lesson of a skill names only the skills that wait for it, and with `paragraphs` it would not wait for the lesson `paragraphs`.

  The order of the lessons by the graph: `first-sentences`, then `paragraphs`, then `endings`, `improving-an-essay` and `writing-about-a-day`, which the graph does not order among themselves. The graph links `reading-model-essays` to nothing. Where the graph does not order, the engine's rule chooses: lessons whose skills are not all `verdict` go by the graph, and where the graph gives no order, by the lesson's tier, the lowest tier of its skills, then by the lesson's ID. Each lesson has a skill that is not `verdict`, so the engine neither adds segments to them nor skips any: they go in the author's order. The choice tasks go in a row in the lesson, and outside a lesson the engine chooses among the three. The texts of the lessons have no `leak_phrases` of the choice tasks.
- **`LICENSE`** — the text of Apache-2.0. There is no `NOTICE`: the texts are its own and those of the format examples of `spec`, under the same Apache-2.0.

### 3.10. The module of the check board

The module `novatrail/check-board` is a small module of maths in Russian for adults: the eleven tasks of the check board of `tutor` — `just board`, `internal/board/tasks.go` of `tutor` —, so that the author solves them in the app again and again, on the path a learner takes: the reading by `tutor`, the verdict of `module-kit`, the hint ladder and the tutor's words. Its eval set checks every new tag of `tutor` once the eval runner of epic E6 is in CI. It stands beside the three modules of the first slice, not among them, and the first slice waits for none of it.

- **The source** — the banks `eval/bank.json` and `eval/free-bank.json` of `tutor` at the tag that the decision log names for the module. Each task is the task of the format that `internal/task.Sheet` of `tutor` builds of its bank's entry, so its open view — the statement, the canvas and the items — is the JSON of the task in `eval/open-tasks.json` of `tutor`, the same for all eleven. The closed fields that `Sheet` fills with stand-ins — the skill, the band `measurement` and the difficulty 1 — are the module's own, below.
- **The manifest**: the ID `novatrail/check-board`, the title «Задачи на листе», the subject `mathematics`, the audience `adult` alone, without grades, one band `board` «Задачи доски», the explanation language `ru`, without `dependencies`. The rung templates are those of the other three modules, in Russian and without digits: «Какой шаг твоей работы дал этот ответ?» and «Прочитай условие ещё раз, медленно. Есть ли в нём то, чего твоя работа ещё не использует?».
- **Skills** — four, without prerequisites and without anchor skills:

  | Skill | Title | How it is assessed | Tasks |
  |---|---|---|---|
  | `linear-equation` | «Линейное уравнение» | `verdict` | `equation-1`, `equation-2` |
  | `calculation` | «Цепочка вычислений» | `verdict` | `chain-1`, `chain-2` |
  | `free-solution` | «Свободное решение» | `verdict` | `posts` and the six of `eval/free-bank.json` |
  | `writing-on-a-sheet` | «Решение на листе» | `none` | a supporting skill of every task |

  The first three are the kinds of the banks, by which `Sheet` names a task's skill. The fourth is practised and not assessed: it keeps the lesson in the board's order, below. All four are in tier 1, so the module is not sequential.
- **Traps** — 22, those that the eleven tasks name, each once: the banks keep one ID and one description for a mistake. The description is the bank's, in impersonal words: the bank's names a gender, «перенёс», and the learner's progress shows a description. AI reads it to suggest the trap from the work. The advice is the module's own, for the learner: addressed with «ты», with no verb in the past tense singular, which would name a gender, and without digits. There are no techniques and no solving steps: `Sheet` gives a task neither, and the format lets a module do without them.
- **Tasks** — 11, `tasks/<id>.yaml` by the ID of the bank, all on a sheet of the default grid:

  | Task | Skill | Key | Difficulty |
  |---|---|---|---|
  | `equation-1` | `linear-equation` | 2 | 2 |
  | `equation-2` | `linear-equation` | 8 | 2 |
  | `chain-1` | `calculation` | 38 | 3 |
  | `chain-2` | `calculation` | 11 | 2 |
  | `posts` | `free-solution` | 5 | 2 |
  | `grasshopper` | `free-solution` | 10 | 2 |
  | `backwards` | `free-solution` | 20 | 2 |
  | `houses` | `free-solution` | 9 | 1 |
  | `stamps` | `free-solution` | 7 | 2 |
  | `pages` | `free-solution` | 141 | 3 |
  | `album` | `free-solution` | 154 | 4 |

  - The statement is one paragraph, and the solution a paragraph per line of the bank. One item, one answer slot `anywhere`, the check `number` with the value of the bank's key, as `sample.KeyOf` of `tutor` takes it: `x = 2` is `2`. The reading of `tutor` gives the number of the final answer without the unknown and the `=` before it and without the words after it (`internal/reading/reading.go` of `tutor`), so the check has no `units`. Every key is an integer, so the decimal comma (О-2) does not arise.
  - The bank's traps are the task's `traps`, which AI may suggest, without trap values: the bank gives no wrong value of a trap. The bank's `false_keys` are a device of the measurement and do not enter the module.
  - The leak phrases are the bank's, without its English ones: the measurement worded hints in both languages, and the module's texts are in Russian. After them stands the key in words, «семь», which the leak check does not read as a number. In `posts` the key, 5, stands in the statement, so the leak check does not search for it, and «5 столбов» stands among the leak phrases in its place.
  - The difficulty is the number of actions of the bank's solution: its lines with a calculation, without the restated equation, the drawing and «Ответ:».
  - There are no `givens`, `unknown`, `leading_question`, `simpler_case`, pictures or templates: a hint asked for before the first wrong answer has nothing to give, and the third rung gives the fallback of the second.
- **The lesson** — one, `lessons/board.yaml`, «Задачи доски по порядку»: the eleven tasks in the board's order, without an explanation. Its skills are all four, and since `writing-on-a-sheet` is not `verdict`, the engine adds no segment to the lesson and skips none. The engine of the first slice does neither in any lesson; the skill keeps the board the same path every time beyond it too, whatever an engine does in a lesson whose skills are all `verdict`.
- **The eval set** — 3 synthetic samples, one per pair of a primary skill and the band, of the task first by ID, as for every module: `chain-1`, `album` and `equation-1`, seeds 1–3. They name the newest release of `modules` when they were drawn, whose index did not hold the module: of a module's eval set the runner takes the task from the package of the commit, and the tag says only what the sample was made on. Its people's samples wait for epic E6.
- **`LICENSE`** — the English text of Apache-2.0, like the other modules'. There is no `NOTICE`: the tasks are those of `tutor`, under the same Apache-2.0.

## 4. Sources

The records of sources are in [sources.md](sources.md), by topic, and a claim above cites its record: `sources 5.1`. The archives of `docs/plan/` cite this section by the title of an item; the item is the record of the same title: The files of a GitHub release — sources 2.1; Immutable releases — sources 2.2; The filters `paths` and required checks — sources 2.3; The age of a GitHub account — sources 2.4; What remains after a push — sources 2.5; `-text` in `.gitattributes` — sources 3.1; DCO 1.1 — sources 3.2; PNG chunks — sources 3.3; The font and drawing — sources 3.4; The shared `release.yml` of `.github` — sources 4.1; The template of `.github` on 2026-11-06 — sources 4.2; The shared `release.yml` and the template on 2026-11-07 — sources 4.3; The `release.yml` of `canvas` — sources 4.4; The shared workflows of `.github` v0.1.15 — sources 4.5; The DCO check — sources 4.6; The CONTRIBUTING of `.github` — sources 4.7; The release of `module-kit` v0.2.3 — sources 4.8; The release of `module-kit` v0.3.0 — sources 4.9; The command `index` — sources 4.10; The command `diff` — sources 4.11; `module-kit` on the previous version — sources 4.12; The module's version and the eval set — sources 4.13; A reproducible archive — sources 4.14; Where the runner takes the samples from — sources 4.15; The request of `api` on the layout of the release — sources 4.16; `sim` on the gate of the graph — sources 4.17; The sample file and the samples page — sources 4.18; The layout of the workbook — sources 4.19; The license of standalone — sources 5.1; standalone at 89986bc — sources 5.2; standalone's tasks by topic — sources 5.3; standalone's test bench — sources 5.4; The first batch from standalone — sources 5.5; standalone at 42b4408 — sources 5.6; standalone's prototype — sources 5.7.
