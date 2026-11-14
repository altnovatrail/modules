# modules — architecture

> Status: a skeleton from `spec`'s template, 2026-10-19; checked against `spec` v0.2.9 on 2026-10-20 (R01), against `spec` v0.6.3 on 2026-11-02 (R02), against `spec` v0.6.5 the same day (R06), against `spec` v0.9.3 on 2026-11-06 (R09), against `spec` v0.9.5 on 2026-11-07 (R18), against `spec` v0.9.8 on 2026-11-08 (R19), against `spec` v0.9.9 the same day (R23) and against `spec` v0.9.10 too (R24), on 2026-11-08 against `spec` v0.9.16 (R27), on 2026-11-09 against `spec` v0.12.0 and format 0.4: the modules 0.2.0 on 0.4, CI on `module-kit` v0.3.0 (R28), and on 2026-11-11 against `spec` v0.15.0: the volunteer's form, the channel, the periods and the consent withdrawal of a sample (R30); against `spec` v0.16.3 on 2026-11-12 — the sources in [sources.md](sources.md), the measurements in [measurements.md](measurements.md), and T65–T67 translate the documents into English (R31); against `spec` v0.17.2 on 2026-11-14 — one call of `sample-check`, the container's name in its network and `legal-3.tar.gz` (R32). The detailing was finished on 2026-11-02 (R08). The first content — epic E2 — was done on 2026-11-07: three modules 0.1.0, their check in CI and the first release `v0.1.0` (R10–R17). The rest of the port — epic E3 — was done on 2026-11-08: `olympiad-math` 0.1.1 holds all 753 tasks of standalone, 19 templates and 41 lessons, and the second release `v0.1.1` (R20–R23). The eval sets — epic E4 — were done on 2026-11-08: the guard of `eval/` in CI and the requests that the first samples wait for (R25, R26). The eval-set mode in CI — epic E5 — was done on 2026-11-09: `sample-check --eval-set` in place of the guard (the 2026-11-09 amendments to R04 and R25). The content of format 0.4 — epic E7 — was done on 2026-11-10: `olympiad-math` 0.3.0 and `grade-3-math` 0.3.0, the release `v0.2.1` (the amendments T49–T55 to R12–R15, R20, R22 and R23). The first eval-set samples — epic E8 — were done on 2026-11-11: the synthetic samples of `canvas`'s generator, the check of the samples against the commit's package in CI and two samples of the author (R29; the amendments T58, T59 to R25 and T61 to R04).
>
> What the component does, what it does not do and what is still to be decided about it. This is the detail of its card — architecture 3.10 spec. The card stays a summary: if the detail changes what the card says, the card is edited in `spec` by the same route as the questions into its registers (RULES.md, Across repositories).
>
> Why it was decided so and not otherwise is in [decisions.md](decisions.md); the work plan is in [RUN.md](../RUN.md). The concepts are concepts spec, the module format is the [format text](../spec/format/README.md), in English (R42 spec). The documents and decisions of `spec` lie in `spec/docs/` and are cited with the suffix "spec".

## 1. Purpose and boundaries

**Does:**
- holds the official modules of the first slice: olympiad maths 1–6 from MathTrail, school maths for grade 3 on a workbook and short essays in English for any age (architecture 3.10 spec; R41 spec; product 10 spec). They have two subjects — maths and writing (R08 spec);
- holds, beside them, a small module of maths in Russian for adults, `novatrail/check-board`: the eleven tasks of the check board of `tutor`, so that the author solves them in the app on the path a learner takes, and so that its eval set checks every new tag of `tutor` once the eval runner of epic E6 is in CI (the 2026-11-14 amendment to R76 spec; section 3.10; R33). The goal of product 10 spec, three modules of two subjects, does not change with it;
- holds all the modules in one repository while the format changes: all the modules move to a new format version in one pull request (R18 spec; R44 spec; repositories 2.1, 4 spec). They may split into a repository per module once the format is v1;
- holds each module as a folder that the format defines (R46 spec; [format text](../spec/format/README.md), section 10; concepts 3 spec):
  - at the folder's root — the module manifest `manifest.yaml`, the skill graph `skill-graph.yaml`, the catalogs `traps.yaml` and `techniques.yaml` and the solving steps `solving-steps.yaml`;
  - the tasks and the lessons — a file for each, `tasks/<id>.yaml` and `lessons/<id>.yaml`;
  - the media files — in `media/`, the eval set — in `eval/`;
  - the license text — the mandatory `LICENSE`, the license notices — `NOTICE`, if they are needed.

  A module is written in YAML 1.2.2 with the format's restrictions, and is checked by the JSON Schema 2020-12 schemas of the format version that its module manifest names — since 2026-11-09, 0.4 (R44, R45, R101 spec; R28). Where the folders themselves lie and what they are called is section 3.1 (R03). There are no widgets in the first slice (product 10 spec);
- keeps in the eval set synthetic samples and samples that adult volunteers and authors made for eval (concepts 3 spec; R02 spec; П-16). The sample file and the eval set are set by the eval protocol: a sample is the message `Sample` of the package `novatrail.eval.v1`, and a module's eval set is each `.json` under its `eval/` (R93, R95 spec; [contracts text](../spec/contracts/README.md), sections 11 and 13). The layout inside `eval/` is flat, `eval/<id>.json` (subsection 2.4; the 2026-11-06 amendment to R04);
- checks the modules in CI (architecture 3.10 spec; repositories 4 spec):
  - with the CLI of `module-kit` — by the format's schemas and rules: a graph without cycles, references, anchor skills, keys, with a template's solver too, fallback texts — with the leak check; it also compares a module's new version by the rules of semver with each published version below it, withdrawn ones too, by the format's rules, and a version in another format version by the places that the format keeps in all its versions (architecture 3.7 spec; concepts 3 spec; R89, R101 spec; [format text](../spec/format/README.md), sections 2 and 27) — section 3.2. CI takes the CLI as a release binary of `module-kit` at an exact tag and, before its first run, checks its signature with `spec`'s recipe `release-verify`, rather than installing it through `go install` (R77 spec);
  - with the simulator `sim`, if what the simulator reads from the package has changed — `sim` decides this —, on the pull request's package, with invariants and absolute thresholds, without bands, from the first tag of `sim`. The first releases of `modules` come out before the simulator, and until its tag there is no gate, and the first tag of `sim` runs the packages already published (R76 spec; R103 spec);
  - the files of the eval-set samples — with `spec`'s recipe `sample-check`, without a key, on every pull request, and from a fork too: with the same code that checks `spec`'s examples (architecture 3.10 spec; R93 spec) — section 3.4;
  - from stage 2, when `tutor` and the samples page of `canvas` appear, — with the eval runner of `tutor` on the module's eval set, with a separate key of the platform at an AI provider — only for eval, not the key of the service account of `api` and not the author's own key (R24 spec). The job with this key, like any job with a secret, does its work only on a push to `main` and on a pull request from a branch of the repository itself. On a pull request from a fork it is not skipped but fails and names the reason, and the merge waits until a maintainer reads the change and carries it into a branch of the repository (R24 spec, the 2026-10-20 amendment). A pull request of Dependabot also comes from a branch of the repository, but gets no Actions secrets (ibid.). Where the key lies and how eval goes on a pull request of Dependabot is the open question Т-28. Its recommendation: the key in Secret Manager of `infra` through Workload Identity Federation, eval on a pull request of Dependabot after the change is carried into a branch of the repository, and while there is no code of `infra`, eval runs from the devcontainer, with the key in an environment variable, not in CI;
- on a release, the CLI of `module-kit` builds the module packages and the module index inside the shared release workflow of `.github`, and that workflow signs them (R17 spec; repositories 4 spec):
  - a package is the server part, the device part and the inventory (R62 spec; [format text](../spec/format/README.md), section 25). The server part is the module's canonical view: all the content, keys, solutions of the tasks and traps included, with the tier of each skill and the numbers of each instance of a template. The device part is `LICENSE`, `NOTICE` and the media files, without the keys and the texts of the module. The inventory names the module and both parts by digests, and the package digest is the digest of the inventory. The eval set is not part of the package: the eval runner takes the samples from the module's folder in CI of `modules`;
  - the index is the file `module-index.json` in each release, one for all modules and cumulative (R63 spec; [format text](../spec/format/README.md), section 26). It names each published version of each module by the digest of its package, with the tag of the release in whose files the package lies, with the anchor skills that the version teaches, and with the mark of withdrawal. Entries do not leave the index and, apart from the mark of withdrawal, do not change;
  - a release builds the index from the previous release's index and its own packages, after the number and before publication. The first release has no previous index. The previous index is downloaded and checked with the recipe `release-verify` by a step of the shared `release.yml` on the runner, and the recipe `release-files` of `modules` gets the already checked file (the 2026-11-01 and 2026-11-02 amendments to R63 spec; Т-35 closed);
  - a version is withdrawn by the next release: its index marks the version's entry `withdrawn` with a reason — `wrong_key`, `answer_leak` or `unsuitable_content`. A withdrawal is not undone, and a withdrawn version stays only for a repeat of a verdict (R63 spec). A withdrawn newest version is fixed by a version above it by semver. A fix that is ready comes out in the same release as the withdrawal, and the withdrawal does not wait for the fix (the second 2026-11-02 amendment to R63 spec);
  - the name of a package's file in a release is derived from the index's entry — without going through the release's files and without the REST API: so the package is found by the build of the image of `api`, CI of `canvas` and a repeat of a verdict (ibid.). The name itself is section 3.3 (R03);
- delivers the packages by stages (architecture 3.10 spec; Т-17):
  - stage 1 — the packages and the index in releases of `modules`;
  - stage 2 — at the build of the image of `api` the digests of the packages are taken from the release's index. A pull request pins the release of `modules` and its index in the configuration of `api` and serves as the anchor of trust, and its CI checks the index's signature with the recipe `release-verify` (R77 spec). With the index the marks of withdrawal reach `api` too. The files of the device part are served by `api`;
  - old versions of packages for a repeat of a verdict are taken from releases of `modules`.

**Does not:**
- holds no code of the platform: a canvas kind and a check kind are code of the platform, shared by all modules, and a module only chooses one and sets its parameters (R02 spec);
- holds no instructions to a model: the prompts belong to the platform, and a module's free text reaches the model only as a separate marked block of data (R02 spec; concepts 3 spec);
- holds no learners' work — neither children's nor adults': a module is published under an open license and cannot be withdrawn. Learners' work, if it is taken into eval, lies outside the repositories and only with a separate consent (concepts 3 spec; П-10);
- does not hide the keys: the keys of open modules are public, like the modules themselves. R05 spec guards no secret, but only keeps the key out of the learner's app (repositories 4 spec);
- does not own the format: the format's schemas, rules, examples and test vectors lie in `spec/format/` (R42, R45 spec), and the code that executes them — the check, the package, the index, the open view, the check kinds — is `module-kit` (R17 spec; architecture 3.7 spec);
- does not publish third-party modules: the recommendation of О-9 is to publish only the official ones at stages 1–3, and the registry service will appear at stage 4 (R17 spec; architecture 3.13 spec).

**Data and links:**
- there are no learners' data and no secrets, as in no repository (RULES.md, Hard rules). How volunteers' samples get into a module was decided by R04 — section 3.4. The consent to them is the open question П-16; its recommendation is a written consent and neither names nor signatures on the samples. What lies in a sample from the samples page and how one consents to it is П-18;
- at the build it depends on nothing but `spec`. CI takes the CLI of `module-kit` and `sim` at an exact tag, and from stage 2 the eval runner of `tutor` (repositories 2.1 spec, repositories 3 spec). A file of a NovaTrail release — the binary of the CLI, the previous release's index — enters CI only after `release-verify` (R77 spec; RULES.md, Hard rules);
- the shared checks of CI from `.github` are connected (T21 `spec`): the secrets search, DCO, CodeQL, Scorecard and the documents check. The merge into `main` is held by all of them but Scorecard, together with the quality gate of SonarCloud, and the ruleset of `main` allows only a merge commit (R30 spec; architecture 3.11 spec);
- the packages are read by `sim` — by the types of the canonical view from `spec`, without `module-kit`, from stage 1, and the gate file of `sim` pins the release of `modules` by its tag and the digest of its `module-index.json`, and the packages by digests (R103 spec; architecture 3.9 spec), — and from stage 2 by `api`, by the digest from its configuration, and by the eval runner of `tutor` (architecture 3.9 spec, architecture 3.10 spec). The learner model marks its state with the release of `modules` (R80 spec);
- from stage 2 CI of `canvas` also takes a release of `modules` at an exact tag, having checked its signature with the recipe `release-verify` (architecture 3.2 spec; R77 spec; Д54 in table 4 of the RUN.md of `spec`). The open views of the tasks for the build of the samples page are built by the CLI of `module-kit` from the package's server part, the media files it takes from the device part, and the marks of withdrawal CI reads in the index of the newest release of `modules` and builds no page from a withdrawn version. Until the first content of `modules`, the build takes the measurement tasks from a release of `tutor` (R25 spec, the 2026-10-21 and 2026-10-28 amendments; R63 spec; architecture 3.2 spec);
- already at stage 1 the generator of synthetic samples of `canvas` reads the module package that `module-kit build` built at the commit where the samples go and that the call gives it — the synthetic samples of the eval set; and the samples page of `canvas` opens a sample file, so that it is looked at before the push (R110 spec; architecture 3.2 and 3.10 spec; section 3.4);
- volunteers' samples come from the samples page of `canvas` by a pull request that the eval-set organizer makes (section 3.4); nothing reaches the platform's server on the way (architecture 3.2 spec). By the first measurement the page has no hosting: its build is downloaded from a release of `canvas` and opened on the device, and while it works it does not go to the network: the tasks and media files are put into it at the build. Child samples are made with the same page, but without a pull request, and do not get into the repository (R25 spec);
- the license is Apache-2.0, for the content, for the code and for the volunteers' samples in the eval sets (R20 spec). The module manifest names it by the SPDX identifier `Apache-2.0`, and the module's folder carries its text as the file `LICENSE` (concepts 3 spec; R46 spec). The notices that the licenses of the content require — for example MIT for the tasks of standalone — stand only in the folder's `NOTICE`: the comments of a module's files do not get into the package (the 2026-10-28 amendment to R46 spec; R62 spec; section 3.5). `LICENSE` and `NOTICE` are UTF-8, their lines end with LF, and git does not convert their line endings: the package carries them byte for byte ([format text](../spec/format/README.md), section 10). Their language is the license's, not the module's: a module of any explanation language holds the English text of Apache-2.0 (format 10 spec);
- the order of work is set by `spec`'s cross-repository plan — rows 3.4–3.7 of [COORDINATION.md](../spec/COORDINATION.md) (R76 spec). The first content is a first batch of olympiad maths for the first application, on which the author solves tasks locally in a browser: 20–30 tasks of standalone in three bands, without templates and solvers, with lessons, and the module manifest names `adult` in the audience too. Small sets of school maths for grade 3 and of essays are written by `modules` to check the format (T38 `spec`). The list of the batch is section 3.6 (R07). The rest of the port comes after the first batch;
- stage 1, together with `module-kit`: the module format v0 is tried at once on three modules of two subjects (repositories 6 spec; R08, R41 spec).

## 2. To decide when detailed

**Outcome 2026-11-02** (T04; R08). All the items are decided. Subsections 2.1 and 2.2 are items of the card (architecture 3.10 spec), and 2.3 is beyond it. The subsections keep their former numbers: the decision log and the tasks of [RUN.md](../RUN.md) refer to them. Of each, what remains here is the gist, the decision and the task, and what was known and what was chosen from is verbatim at the end of section 3.3, 3.4 or 3.5, beside how it is built. How it is built is section 3, and what to build is section 3.7.

**2026-11-06** (R09): the eval protocol left to `modules` the layout of the samples inside `eval/` — subsection 2.4. The author decided it the same day, and section 2 again names only what is deferred.

A new item, if one appears, is the next free subsection, 2.5 and on, by rule 2 of [RUN.md](../RUN.md): the gist, what is known, the options and the task. A decided item names the log's entry, and the decision itself moves to section 3. A question that goes beyond this repository goes «в реестр spec» (RULES.md, Across repositories).

Since v0.6.5 the card names two more items — the name of a package's file in a release and where CI takes a module's earlier versions from for version comparison (architecture 3.7 spec). Both belong to subsection 2.1, and R03 decided them (R06). Subsection 2.3 the card at v0.6.5 did not name: the module of olympiad maths comes from MathTrail (product 10 spec; R08 spec), and the port of its tasks is the first content of `modules` (repositories 6 spec). Since v0.6.6 it names the port as the item that remains, and at v0.9.3 it names the name of the archive, and where CI takes the earlier versions from, and the port with the first batch, and its "To decide when detailed" is "everything is decided" (architecture 3.10 spec; rows В12 and В15, closed — R09).

**Deliberately deferred:**
- the code and the files of the modules — epic E8 in table 1 of [RUN.md](../RUN.md), with tasks T58–T63, and epic E6 in table 2, by section 3.7; the first content, epic E2, was done on 2026-11-07 (R10–R17), the rest of the port, epic E3, on 2026-11-08 (R20–R23), the eval sets, epics E4 and E5, on 2026-11-08 and 2026-11-09 (R25, R26; the 2026-11-09 amendments to R04 and R25), and the content of format 0.4, epic E7, on 2026-11-10, in the archive [plan/E7.md](plan/E7.md);
- the checks of the eval-set samples in CI — with the recipe `sample-check` by the eval protocol (R93–R95 spec), in epics E4–E6 and E8 — E4 and E5 in the archives [plan/E4.md](plan/E4.md) and [plan/E5.md](plan/E5.md), E6 and E8 in table 2 of [RUN.md](../RUN.md): the guard and the layout — E4, done (R25), the eval-set mode in CI — E5, done in T47 (the 2026-11-09 amendment to R25), the first samples, synthetic ones and the authors' samples, — E8, the volunteers' samples and the runner — E6. The text of the form, the channel of the file and the keeping of the forms are epic E4 in the RUN.md of `spec`, and the eval runner comes from stage 2 (section 3.4);
- the gate of `sim` — from the first tag of `sim` (section 3.2). What it includes was decided by R103 spec: everything that the simulator reads from the package, without bands. The change of CI for the gate is Д5 in table 4 of [RUN.md](../RUN.md);
- the module of olympiad maths in Russian is another module, by О-2: the explanations of this module are in English (product 10 spec; section 3.1);
- the parameter of 12 or 24 hours of the check kind `time` — the format will give it by a minor change when content asks for it (R111 spec). Until then the 12-hour miss is held by a trap of its own, `twelve-hours-off`, and its addition to a format example is the row В35 «в реестр spec» in table 3 of [RUN.md](../RUN.md) (the T50 amendment to R23; section 3.7). The rest of what format 0.4 gave was taken by epic E7 — the archive [plan/E7.md](plan/E7.md).

### 2.1. The layout of the modules in the repository

- **The gist:** where in the repository the modules' folders lie and what they are called; how CI finds which module has changed and whose graph to run on `sim`; how releases come out — which modules a release builds, how the packages lie in its files and where the withdrawals lie; which tag CI of `canvas` takes. What lies inside a module's folder is set by the format. It is decided here: it is an item of the card (architecture 3.10 spec).
- **Decided:** R03 — sections 3.1–3.3. What goes beyond `modules` went «в реестр spec», as rows В2–В4 and В12 of table 3 of [RUN.md](../RUN.md). В4 was closed by the move to v0.6.5 (R06), and В2, В3 and В12 by the move to v0.9.3 (R09). The 2026-11-02 amendment to R03 — by `spec` v0.6.5 (R06): a withdrawn version is fixed by a version above it by semver, and CI compares the new version with each published one. The 2026-11-06 amendment to R03 — by `spec` v0.9.3 (R09): the comparison across a format transition too, the form of the withdrawals `$defs/withdrawals` and the list of `paths` in the template of `.github`. The 2026-11-07 amendment to R03 — by `spec` v0.9.5 (R18): the gate of `sim` is not "whose graph" but everything that the simulator reads from the package (R103 spec). The 2026-11-08 amendment to R03 — by `spec` v0.9.8 (R19): in a format transition nothing that the simulator reads is changed, except the writing itself.
- **How it was decided:** what was known and what was chosen from is at the end of section 3.3, verbatim.
- **Task:** T01.

### 2.2. Volunteers' samples in the eval set

- **The gist:** how a sample that an adult volunteer made on the samples page of `canvas` gets into a module's eval set: who opens the pull request, what is checked before the merge, where the consent to publication is marked. There are three modules, and one path for all of them: handwriting on a sheet, the answer slots of a workbook and essays. It is decided here: it is an item of the card (architecture 3.10 spec). The consent itself is П-16, a question of `spec`'s register: its text goes to the lawyer together with П-1, and П-16 is decided by epic E4 «Приватность до первого ребёнка» in the RUN.md of `spec`. So the path of a sample is decided here, and what the task finds out about consent goes «в реестр spec».
- **Decided:** R04 — section 3.4. What goes beyond `modules` went «в реестр spec», as rows В6–В9 of table 3 of [RUN.md](../RUN.md), and all of them were closed by the move to v0.9.3 (R09). The 2026-11-06 amendment to R04 — by the eval protocol: the sample file, the check `sample-check` and the eval set, and the layout of `eval/` is subsection 2.4.
- **How it was decided:** what was known and what was chosen from is at the end of section 3.4, verbatim.
- **Task:** T02.

### 2.3. The port of standalone's tasks

- **The gist:** what of MathTrail's maths goes into the module of olympiad maths and what each part turns into, which tasks go into the first batch, and on what terms their content enters the module under Apache-2.0. The card at v0.6.5 did not name this item, and since v0.6.6 it does (the opening of section 2). The module of olympiad maths that it names comes from MathTrail (architecture 3.10 spec; product 10 spec; R08 spec), and the port of its tasks stands first in the order in which `modules` appears (repositories 6 spec).
- **Decided:** R05 — section 3.5, and the first batch — R07, section 3.6. What goes beyond `modules` went «в реестр spec», as rows В10 and В11 of table 3 of [RUN.md](../RUN.md), from the first batch — as rows В13 and В14, and from the outcome of the detailing — as row В15. All of them were closed by the move to v0.9.3 (R09).
- **How it was decided:** what was known and what was chosen from is at the end of section 3.5, verbatim.
- **Task:** T03, the first batch — T07.

### 2.4. The layout of the samples in `eval/`

- **The gist:** where inside a module's `eval/` the sample files lie and what they are called. The eval protocol left this to `modules`: a module's eval set is each file under its `eval/` at any depth whose name ends in `.json`, except files and folders starting with `.`, and "Where the samples lie within `eval/` is for the repository that holds the modules to decide" ([contracts text](../spec/contracts/README.md), section 13; R95 spec). The format lets any path into `eval/` ([format text](../spec/format/README.md), section 10). The item appeared with `spec` v0.9.3 (R09).
- **Decided:** flat, `eval/<id>.json` — the 2026-11-06 amendment to R04, section 3.4; since 2026-11-09 there is nothing starting with `.` in `eval/` either (the 2026-11-09 amendment to R04).
- **How it was decided:** what was known and what was chosen from is at the end of section 3.4, verbatim.
- **Task:** the move to `spec` v0.9.3 (R09), the author's answer of 2026-11-06.

## 3. How it is built

What is decided by section 2 moves here: what parts the component consists of, its data and how it is linked to other repositories. Why it was decided so and not otherwise is in the log. Here is what was decided by subsections 2.1 (R03), 2.2 (R04) and 2.3 (R05 and R07). Since T12 the repository has its first module — the first batch `novatrail/olympiad-math` 0.1.0 (R13), since T13 also its lessons (R14), since T14 the set of school maths for grade 3 `novatrail/grade-3-math` 0.1.0 (R15), and since T15 the set of short essays `novatrail/english-essays` 0.1.0 (R16), and since T69 the module of the check board of `tutor`, `novatrail/check-board` 0.1.0 (R33). There is a `.gitattributes` at the root, and since T09 the CLI of `module-kit` in the `toolchain` stage, the recipes of CI in the `justfile` and their tests (R10). Since T16 — `release.yml`, the recipe `release-files` and the first release `v0.1.0` (R17). What to build, in which epic and what it waits for is section 3.7.

### 3.1. The repository's tree

```
modules/                     the repository root
├── modules/                 module folders, and nothing else
│   ├── olympiad-math/       novatrail/olympiad-math — olympiad maths 1–6
│   ├── grade-3-math/        novatrail/grade-3-math — school maths for grade 3 on a workbook
│   ├── english-essays/      novatrail/english-essays — short essays in English
│   └── check-board/         novatrail/check-board — the tasks of the check board of tutor, in Russian
├── withdrawals.yaml         withdrawals, from the first withdrawal on — an input of the index build (section 3.3)
├── justfile                 CI recipes and the recipe release-files
├── test/                    tests of the justfile recipes, node --test
├── go.mod, go.sum           the Go module github.com/altnovatrail/modules: the code of the tools (R12)
├── tools/port/              the converter of standalone tasks: code, batches batches/, overlay files overlay/,
│                            lesson pictures pictures/ and the sources of instance pictures template-pictures/
├── .gitattributes           modules/** -text
├── docs/                    the component documents; the language of each is in CLAUDE.md
├── spec/                    the submodule spec at a tag
├── .github/                 CI: calls of the shared workflows, own jobs, release.yml
├── .devcontainer/           the spec container; in the toolchain stage — the CLI module-kit
├── README.md, LICENSE       the README and the license of the repository
└── CLAUDE.md, AGENTS.md, RUN.md
```

Since T09 the tree has the `justfile` with the recipes of CI, its tests in `test/` and the CLI in the `toolchain` stage (R10), and since T11 the converter of standalone's tasks `tools/port` and `go.mod` (R12). Since T12 `modules/` has `olympiad-math/` — the first batch (R13), since T13 with nine lessons (R14), since T14 `grade-3-math/`, the set of school maths for grade 3 (R15), and since T15 `english-essays/`, the set of short essays (R16), and since T69 `check-board/`, the module of the check board of `tutor` (R33). Since T16 there are `release.yml` and the recipe `release-files` (R17); `dist/`, where the recipe puts the release files, git does not keep — `.gitignore`. There is no `withdrawals.yaml` yet: it appears with the first withdrawal (section 3.3).

- **The module's folder** is `modules/<name>/`, and `<name>` is the second half of the module ID: the module ID in the folder `modules/<name>/` is `novatrail/<name>`. What lies inside is set by the module layout ([format text](../spec/format/README.md), section 10; R46 spec):

  ```
  modules/<name>/
  ├── manifest.yaml, skill-graph.yaml, traps.yaml, techniques.yaml, solving-steps.yaml
  ├── tasks/<id>.yaml, lessons/<id>.yaml
  ├── media/<path>             pictures and sound; they go into the device part of the package
  ├── eval/<path>              eval-set files; they do not go into the package
  └── LICENSE, NOTICE          the text of Apache-2.0 and the license notices
  ```
- **In `modules/` — only the modules' folders:** neither other files nor files and folders starting with `.`. Each module in `modules/` on `main` comes out in a release — as its package, without `eval/` (R62 spec). A module's draft lives in a branch.
- **The names** name the subject and the audience, not the explanation language: the languages of the slice are О-2. A module in another explanation language is another module ([format text](../spec/format/README.md), section 10), and its name will get a suffix with a BCP 47 tag, for example `olympiad-math-en`. A module ID never changes: the full IDs of its skills and tasks are kept by the learner model and the progress.
- **`LICENSE`** lies both at the root — the repository's license (Д22 in table 4 of the RUN.md of `spec`), — and in the folder of each module: this text the package carries to the learner's device (the 2026-10-28 amendment to R46 spec).
- **`.gitattributes`** is `modules/** -text`: git converts the line endings of no file of the modules, neither at a commit nor at a checkout. `LICENSE`, `NOTICE` and the media files the package carries byte for byte, and CR LF would change the package digest (Д44 in table 4 of the RUN.md of `spec`). A `LICENSE` or `NOTICE` with CR would not pass the rule `module.plain-text` anyway.

### 3.2. The checks of CI

| Check | When | Which modules |
|---|---|---|
| the CLI of `module-kit`: the format's schemas and rules | every pull request and push to `main` | all |
| the tree: in `modules/` only the modules' folders, the ID of each is `novatrail/<folder>` | the same | all |
| the package and the comparison with what is published | the same | all |
| the simulator `sim` | from the first tag of `sim` (R76 spec; R103 spec) | those in which what the simulator reads from the package has changed — `sim` decides — and new ones; all — when a pull request changes the tag of `sim`; none — in a format transition |
| the layout of `eval/` and the eval set of each module: the recipe `just eval-check`, and in it `spec`'s recipe `sample-check --eval-set novatrail/<folder>`, one call per module (R04; R25 and its amendments T58 and R32) | every pull request and push to `main`, without a key, and from a fork too | all that have `eval/`; the mode — those whose `eval/` has `.json` files (section 3.4) |
| the samples of `eval/` are samples of the commit's package: a synthetic one is drawn again by the generator of `canvas` and compared with the file in everything but the version and the digest of the package, for a person's sample the digest of the open view is checked by `module-kit open`; the recipe `just eval-match` (the T59 amendment to R25) | every pull request and push to `main`, without a key and without the network, and from a fork too | all whose `eval/` has samples |
| the eval runner of `tutor` | pull request, from stage 2 (R24 spec) | with a changed `eval/` or package, new ones too; all — when a pull request changes the tag of `tutor` (section 3.4) |

- **All modules on every pull request.** The CLI of `module-kit` checks each folder `modules/*/` and builds its package. So a new tag of `spec` or of `module-kit` re-checks all the modules without a special rule.
- **The CLI is pinned in one place** — in the `toolchain` stage of `.devcontainer/Dockerfile`: the release binary of `module-kit` at an exact tag and sha256 sum (R34 spec).
  - The recipe `release-files` takes the same binary too (section 3.3), so CI checks the modules with the same version of the CLI that the release will build them with.
  - The binary's signature is checked by the recipe `release-verify` in CI of each pull request and push to `main` — the job "module-kit signature", the recipe `just module-kit-verify`: the release's `SHA256SUMS` and the check of its lines against the sums of the `ARG`. From then on the sum holds the binary (R77 spec; RULES.md, Hard rules; the second 2026-11-06 amendment to R03).
  - This is how the container differs from `spec`'s container — R10 (CLAUDE.md, Environment).
- **The generator of synthetic samples is pinned there too** — the file `synthetic-module_<tag>.mjs` of a release of `canvas` at an exact tag and sha256 sum (R29; R68 `canvas`). The release's signature is checked by the job "canvas signature", the recipe `just canvas-verify`: the release's `SHA256SUMS` through `release-verify` and the check of the generator's line against the sum of the `ARG`. The recipe of synthetic samples `just synthetic` is a person's step in the devcontainer, not CI (section 3.4); in CI the generator is called by `just eval-match`, the check of the samples of `eval/` against the commit's package (the T59 amendment to R25).
  - The form of the call of the commands `check`, `build`, `index` and `diff` is kept by architecture 3.6 module-kit: the recipes of `modules` cite it and check the call with their tests (R102 spec; the 2026-11-06 amendment to R03).
- **The comparison with what is published.** CI takes the `module-index.json` of the newest release of `modules` and checks its signature with the recipe `release-verify`. For each module:
  1. the version is in the index — the digest of the package that CI built must match the digest of the entry. Otherwise the content was edited without raising the version, and the release would refuse by `module-index.one-package` ([format text](../spec/format/README.md), section 26);
  2. the version is not in the index — it must be above each published version of the module, withdrawn ones too. The CLI of `module-kit` compares the module with the package of each published version by the format's rules of version comparison: which rules go against each version and which against the newest, the format says (architecture 3.7 spec; [format text](../spec/format/README.md), section 27; the 2026-11-02 amendment to R03). CI downloads the packages by the entries (section 3.3) and checks them by the chain of digests from the checked index;
  3. a published version is written in another format version — while the format is v0, this happens after each of its transitions (R44 spec). The CLI reads its package by the places that the format keeps in all versions: the inventory, the names of the package's files, the languages, the bands, the skills and `retired_skills`, the keys of the tasks, lessons, traps and techniques. The rules about skills, `retired_skills` and IDs go across a transition too, and the version part, while the format is v0, is major ([format text](../spec/format/README.md), sections 2 and 27; R89 spec; the 2026-11-06 amendment to R03);
  4. there is no index yet, before the first release, — all the versions are new, and there is nothing to compare with.

  How CI does this (R11): the newest release is the greatest tag `vX.Y.Z` from `git ls-remote`, without a token, that has a `module-index.json`; a tag whose index answers 404 is a release that did not come out, and the log names it. Its `module-index.json` is checked by `release-verify`, and the archive of each version of the modules of `modules/`, withdrawn ones too, comes from the address in the entry (section 3.3). Only `inventory.json`, `server.json` and `device.json` are taken from the archive; an archive with a link, a special file or a path outside the folder is rejected whole. Then comes `module-kit diff --index … --published … <package>` over the package that `build` built, and without a release `--first`: that there is no release the recipe says explicitly, with the file `first`, and an index silently skipped does not pass. The chain of digests from the index to the package's files is checked by `diff`, and its exit code is the job's code.
- **Whether what the simulator reads has changed** — for the gate of `sim`. The gate is not one graph but everything that the simulator reads from the package, and whether it has changed is decided by `sim` itself, with its command "has it changed" (R103 spec; Т-16 closed). The simulator runs on the pull request's package that CI built, with invariants and absolute thresholds, without bands: the pull request's package has no base to compare them with. `sim` reads `server.json` by the types of the canonical view, without `module-kit` (architecture 3.9 spec). A pull request that changes the tag of `sim` runs all the modules (repositories 4 spec).
- **A format transition goes without `sim`.** If the newest published version of a module is written in another format version, the simulator is not run in this pull request. `sim` at its former tag will not read the package of the new version. So in a transition nothing that the simulator reads from the package is changed — neither the graph nor the bindings of the tasks — except the writing itself in the new format version, and the review of `modules` holds this. All of that is changed by a separate pull request after step 4 of repositories 4 spec, with the simulator (the 2026-11-07 amendment to R103 spec; the answer to В23).
- **Until the first tag of `sim` there is no gate** (R76 spec): the first release `v0.1.0` came out without it, and the first tag will run the packages already published. The change of CI for the gate is Д5 in table 4 of [RUN.md](../RUN.md), a job by the form of the call that `sim` will give with its first tag. The first tag came out on 2026-11-09, but the newest, v0.1.8, does not translate the check kind `time` and refuses the package `olympiad-math` 0.3.0, so the job waits for a tag of `sim` that translates it (T42 of the RUN.md of `sim`; the author's decision of 2026-11-11, Д5). The former answer — a gate only by the graph, by the canonical text of the member `skill_graph`, until T03 of `sim` answers — was replaced by R103 spec (the 2026-11-07 amendment to R03).
- **The jobs and the recipes** (R10; R29). The own jobs of `ci.yml` build the `toolchain` stage and run the recipes of the `justfile` in it with the script `spec/tools/toolchain/toolchain.sh` from the submodule. "module-kit signature" is `just module-kit-verify`, "canvas signature" is `just canvas-verify`; "Modules" waits for both and runs `just test`, then `just modules-check`: the tree of `modules/`, the version of the CLI on `PATH` against `MODULE_KIT_VERSION`, `check` and `build` of each folder into a temporary folder, and the ID from the package's `inventory.json`. Then `module-kit diff` compares each package — with the published versions by the index of the newest release or, before the first release, with `--first` (R11). Without `modules/` the recipe prints `modules/: no modules yet` and exits with code 0. Without modules the recipe does not go to the network. The refusal of one module does not stop the check of the others; the recipe's code is the greatest of the codes of what failed: 1 — the tree, R03 and a refusal by the rules, 2 — `module-kit` could not do its work, or the published versions could not be read. After it "Modules" runs `just eval-check` — the layout of each module's `eval/` and its eval set with the mode `sample-check --eval-set` (R25; section 3.4), — then `just eval-match` — each sample against the package of this commit (the T59 amendment to R25), — the check of the lessons' pictures, `just lesson-pictures modules/olympiad-math --check` (the T41 amendment to R12), — and the check of the pictures of the templates' instances, `just template-pictures modules/olympiad-math --check` (the T52 amendment to R12).
- **No filter `paths` on the workflow.** What to run is decided by the job itself. The checks of a workflow that a filter skipped hang "Pending", and a required one holds the merge (sources 2.3).

### 3.3. The release

- **Who and when.** The own `release.yml`, by the template of architecture 3.3 .github, calls the shared `release.yml` of `.github` after a green CI on `main`, if something of `paths` has changed since the previous release. The number is the next `vX.Y.Z` of its line, as for any NovaTrail repository (R31 spec): the line is `RELEASE_LINE`: `"0.1"` from the first release, `v0.1.0` (R17), and `"0.2"` from the move to format 0.4, whose first release is `v0.2.0` (R28). Epic E7 comes out as one release of the same line, `v0.2.1`: the merge of the branch `e7` into `main` after the author's confirmation, with the new packages `olympiad-math` 0.3.0 and `grade-3-math` 0.3.0 and the former entries of the index in their places (decision 1 of epic E7; RUN.md, T56). The shared `release.yml` is v0.1.15, by SHA 89d3cee.
  - **A dry run on a pull request.** The same call runs on each pull request from a branch of the repository itself with `publish: false`, as in `canvas`: the files are built from `git archive` in the `toolchain` stage, signed and checked, but nothing is published. A pull request from a fork skips the dry run — it is given no `id-token`, and there is nothing to sign with. The dry run is not in the ruleset of `main` (R17). The input `previous-file: module-index.json` stands from the first release, and the shared `release.yml` for it is no lower than v0.1.14 (sources 4.1). Immutable releases and the ruleset of the tags `v*` — a condition of the signature check — stand since 2026-11-06: the ruleset `v-tags`, 26725460, with the rules `update`, `deletion` and `non_fast_forward`, without bypass (R77 spec). The third condition — the call of the shared `release.yml` — stands since T16 (R17), and Д3 was closed by the first release, `v0.1.0`, whose files `release-verify` accepted.
  - **The epic branch.** An epic that puts out one release at its end gathers the pull requests of its tasks in the epic's branch, not in `main`: a release comes out after each green CI on `main` if something of `paths` has changed, and a change of the `justfile` touches them too. The branch of E3 is `e3`: it stands in the include of the ruleset of `main` and holds the same required checks, the dry run of the release on its pull request publishes nothing, the button `workflow_dispatch` is not pressed on it, and T42 merges it into `main` (R20).
- **`paths`** is what the release is built from (the 2026-10-22 amendment to R31 spec). The list for the kind `modules` is kept by the template of architecture 3.3 .github, and this section refers to it (the 2026-11-01 amendment to R31 spec; the 2026-11-06 amendment to R03). On `main` of `.github` (af40f94, 2026-11-06, and the same on 2026-11-07) the template does not carry it yet, and `.github` will bring it in with its move to a tag (Д80 in table 4 of the RUN.md of `spec`). Until then, here is the list that the template will get, and where it diverges from the template, the template is right. The same list is carried by `release.yml` since T16 (R17):

  ```
  :(glob)modules/**
  :(exclude,glob)modules/*/eval/**
  withdrawals.yaml
  justfile
  LICENSE
  .github/workflows/release.yml
  ```

  - R03 gives the rule by which the list is made, and the dry run, and the reason why it has no Dockerfile is carried by the template too.
  - No release comes from the eval sets, the documents, CI other than `release.yml`, the submodule `spec` and `.devcontainer/Dockerfile`. The package and the index depend only on the files of the module and the format version of its module manifest ([format text](../spec/format/README.md), section 25), not on the version of the CLI in the Dockerfile.
  - The price is the whole `justfile`: a change of a CI recipe in it will give a release with the same index.
  - The magic `glob` is there so that `*` does not match `/`: without it `modules/*/eval/*` would exclude a module's `media/eval/` too.
- **The recipe `release-files VERSION PREVIOUS`** runs in the container of the `toolchain` stage, on the commit's tree without git, submodules, network and tokens.
  - The CLI of `module-kit` is the same binary of this stage as CI's (section 3.2).
  - `PREVIOUS` is a read-only path to the `module-index.json` of the newest release, which the shared release downloaded and checked with `release-verify` before the build. For the first release it is an empty string, and the recipe tells the build explicitly that there is no previous index (the 2026-11-02 amendment to R63 spec).

  The recipe builds the packages of all the modules and the index: from the previous index, its own packages, the tag `VERSION` and `withdrawals.yaml` ([format text](../spec/format/README.md), section 26). The package of a version that the previous index already names does not change its entry. If the digest of such a package is different, the build refuses by `module-index.one-package`.

  How it does this (R17):
  - the input: `VERSION` is a tag `vX.Y.Z`; `PREVIOUS` is an empty string, then `--first`, or a file, then `--previous`; `withdrawals.yaml` — `--withdrawals`, and without it `--no-withdrawals`. Anything else is code 2, and in `modules/` there are only the modules' folders, otherwise code 1;
  - the tree, the CLI on `PATH` against `MODULE_KIT_VERSION`, `check` and the ID of each module — with the recipe `_modules-check`, without the comparison with the published ones: that is the business of CI of the same commit, and what the release needs of it is held by `index`. A relative `PREVIOUS` is from the folder the recipe is called from;
  - `module-kit build` of each module into a temporary folder, then `module-kit index … --tag VERSION … --out <file> <packages>` — by the form of architecture 3.6 module-kit; the code of `index` is the recipe's code;
  - an archive — only for a package whose entry in the new index names `VERSION`. Before the archive the recipe checks the entry's `package_digest` against the sha256 of the bytes of the package's `inventory.json`, and the package's composition — `inventory.json`, `server.json`, `device.json` and the folder `device/`, and nothing more: otherwise code 2;
  - `dist/` is built in a temporary folder and takes its place whole, at the end: the old `dist/` is removed at the start, and a run with a refusal leaves no `dist/`.

  If not one module has changed since the previous release, and the change is in the `justfile` or `release.yml`, the release carries the same index and `LICENSE` (R03).
- **The release files** are what the recipe puts into `dist/`. The names are of letters, digits and `. _ + -`, without folders. The bundle `<file>.sigstore.json` for each is put by the shared release:

  | File | What it is |
  |---|---|
  | `module-index.json` | the module index: each published version of each module, from the first release |
  | `LICENSE` | the text of Apache-2.0 from the root |
  | `sha256-<hex>.tar.gz` | the archive of a package — for each version that this release publishes, and only for it |

  The package of a version that the index already named lies in the release that its entry names, and does not come out a second time. A release that only withdraws versions carries the index and `LICENSE`.
- **The archive of a package** is named by the package digest — the entry's `package_digest`, with `:` replaced by `-`, — with `.tar.gz`. Inside is the package as files, at the root of the archive, without a top folder ([format text](../spec/format/README.md), section 25). Apart from these files and their folders, there is nothing in the archive — neither links nor other entries:

  ```
  inventory.json  server.json  device.json
  device/LICENSE  device/NOTICE  device/media/<path>
  ```

  The archive is the same byte for byte on any runner and at any time of the build (R17): GNU tar — `--format=gnu --sort=name --owner=0 --group=0 --numeric-owner --mode=u=rwX,go=rX --mtime=@0`, the entries in the order `inventory.json`, `server.json`, `device.json`, `device/`, without `./` before the names; then `gzip -n -9`, without a name and a time in the header. `tutor` builds its archive the same way (sources 4.14).

  The address is derived from the index's entry, without going through the release's files and without the REST API: `https://github.com/altnovatrail/modules/releases/download/<release>/sha256-<hex>.tar.gz`. Nothing holds the bytes of the archive, and nothing needs to. A reader checks the inventory against the package digest, the parts against the inventory, the files against the list of the device part, and takes only the paths that this list names.
- **Withdrawals** are `withdrawals.yaml`, YAML by the format's section 4. The form is `$defs/withdrawals` of the index's schema, `module-index.schema.json`: module ID → version → withdrawal reason ([format text](../spec/format/README.md), section 26; the 2026-11-06 amendment to R03). The file is read by the CLI's command "build the index" (architecture 3.7 spec). For example:

  ```yaml
  novatrail/olympiad-math:
    "1.2.0": wrong_key
  ```

  - The form requires at least one entry, so there is no file while nothing is withdrawn. Then the recipe tells the build that there are no withdrawals.
  - The file is in `paths`, so a pull request with only a withdrawal gives a release (Д45 in table 4 of the RUN.md of `spec`).
  - Entries do not leave the file: a repeat of a withdrawal with the same reason changes nothing, and another reason would break `module-index.entries-kept`. A version that comes out in the same release can be withdrawn too (R63 spec).
  - If the newest version of a module is withdrawn, it is fixed by a version above it by semver, and which part to raise the rules of semver decide. A change of a published task's key is major (concepts 3 spec), so the fix of a version withdrawn by `wrong_key` is the next major. A fix that is ready goes in the same pull request as the withdrawal and comes out in the same release: otherwise, until it comes, the current version would be the version below the withdrawn one, and if a new major is withdrawn, the previous major, and there would be nothing to carry the withdrawn version's evidence about its skills over with. The withdrawal does not wait for the fix (the second 2026-11-02 amendment to R63 spec; the 2026-11-02 amendment to R03).
  - A version older than the newest cannot be fixed: one line of versions lives in a module's folder, and a new version is above all the published ones (section 3.2).
- **Who takes the release.**
  - CI of `canvas` — from stage 2, at an exact tag `vX.Y.Z` of a release of `modules`: the number of the repository itself, not a module's version. By the tag it reads the index, and it takes the packages from the releases that the index's entries name. The index of the newest release it reads only for the marks of withdrawal (the 2026-10-28 amendment to R25 spec).
  - `api` — from stage 2: a pull request pins the release and its index (section 1), and the build of the image takes the packages by the entries.
  - `sim`, which reads packages from releases of `modules` (architecture 3.9 spec), finds them by the same layout: a package's address is derived from the index's entry.
- **Old releases** are not deleted or rewritten: releases are immutable, and a verdict is repeated on old packages (architecture 3.10 spec). The second release needs the first one's index — the shared release gives it, since v0.1.14 (row 3.7 of [COORDINATION.md](../spec/COORDINATION.md)).

**How it was decided** — subsection 2.1 before the outcome of the detailing, verbatim (R08). "What is known" was written before the decision, so "`modules` decides" there speaks of the day it was written, and the answer is in sections 3.1–3.3. The gate of `sim` by the graph, which "What is known" and the options speak of, is replaced: the gate is everything that the simulator reads from the package (R103 spec; the 2026-11-07 amendment to R03).

- **What is known:**
  - all the modules are in one repository while the format changes; they can be split once the format is v1 (R18 spec; Т-2; repositories 2.1 spec);
  - what is inside a module's folder and what its files are called is set by the layout table ([format text](../spec/format/README.md), section 10; R46 spec). Where the folder itself lies and what it is called the format does not say: `modules` decides this;
  - a module ID is `<publisher>/<name>`, both in kebab-case starting with a letter: `novatrail/olympiad-math` (R46 spec; concepts 3 spec). The official modules are published by the project. The format gives no names for school maths and essays: the modules themselves name them;
  - a module's version is semver `X.Y.Z` in the module manifest (R46 spec). The tag of a release of `modules` is the number `vX.Y.Z` of the repository itself, as for any NovaTrail repository (R31 spec), not a module's version: the index names for each version of a module the tag of the release in whose files its package lies (R63 spec);
  - a release of another repository is caused only by the files its release is built from: the template of `paths` of its `release.yml` is their list by the kind of repository, and it is kept by the template of architecture 3.3 `.github` (the 2026-10-22 and 2026-11-01 amendments to R31 spec);
  - the index is one and cumulative, and a package whose version is already in the index does not change its entry. So a release may build all the modules or only the changed ones — `modules` decides this (R63 spec). How the packages lie in the release files `modules` decides too (R62, R63 spec);
  - the withdrawals are an input of the index's build. The file of withdrawals must be in the release's `paths`, otherwise a pull request with only a withdrawal will not give a release (Д45 in table 4 of the RUN.md of `spec`);
  - the first release builds the first index. The second needs the previous one's index: the shared release downloads and checks it, and how the file reaches the recipe `release-files` is decided by T41 of `.github` (the 2026-11-02 amendment to R63 spec; row 3.7 of [COORDINATION.md](../spec/COORDINATION.md));
  - git must not convert a module's `LICENSE` and `NOTICE` to CR LF: otherwise the package digest will change (Д44 in table 4 of the RUN.md of `spec`; [format text](../spec/format/README.md), section 10);
  - `sim` runs if the skill graph has changed — from the first tag of `sim` (R76 spec; Т-16). The modules move to a new format version in one pull request: `sim` runs in it only if the graph has changed, and the graph is changed by a separate pull request when `sim` has moved to a new tag (repositories 4 spec);
  - CI of `canvas` also takes a release of `modules` at an exact tag, for the samples page, and reads the marks of withdrawal in the index of the newest release (section 1; R25, R63 spec);
  - in MathTrail the content of the maths lies in standalone, in `content/` and `site/` (repositories 7 spec).
- **Options** (T01; the first in each item was chosen; the sources are sources 2.1, 2.2, 2.3, 3.1, 4.1, 4.12, 4.16, 4.17 and 5.3):
  - where the modules' folders lie: `modules/<name>/`, where `<name>` is the second half of the ID; `<name>/` at the root, beside `docs/` and `spec/`; `modules/novatrail/<name>/`, a path by the full ID. The format's index vectors lay out a release's modules as `modules/<name>/` ([format text](../spec/format/README.md), section 9). Standalone gives no precedent: its tasks lie by topic — a file `content/examples/<topic>.json` per topic, with all its tasks, — and the format requires a file per task (R46 spec);
  - the names: `olympiad-math`, `grade-3-math` and `english-essays`; `olympiad-math`, `school-math-3` and `short-essays`; `olympiad-math`, `school-math-grade-3` and `essays-in-english`; names with the explanation language, though the languages of the slice are О-2, and it is open;
  - how the packages lie in the release files: an archive `sha256-<hex>.tar.gz` per package, only the versions that the release publishes itself; an archive named by ID and version, `novatrail+olympiad-math+1.2.0.tar.gz`; an archive by digest, but the packages of all the current versions in each release; a package as loose files; one archive for all the packages of a release. The files of a GitHub release are a flat list, no more than 1000, with safe names. The request of `api`: the name of a package's file is derived from the index's entry, best of all from the digest (В6 in table 3 of the RUN.md of `api`);
  - whether a release builds all the modules or only the changed ones: all, and in the files only the new versions; only the changed ones. The recipe `release-files` gets the commit's tree without git (architecture 3.1 .github);
  - what CI compares a module with — by semver and for the gate of `sim`: with the newest published version, by the index of the newest release and its package; with the module's folder in the pull request's base; with the module's folder at the tag of the release that published the previous version. `module-kit` left this choice to `modules` (architecture 2.4 module-kit), and `sim` asks how CI learns that the graph has changed (T03 in the RUN.md of `sim`);
  - where the withdrawals lie: `withdrawals.yaml` at the root; a file in the module's folder, which would break `module.layout`;
  - the release's `paths`: the list by the rule of the 2026-10-22 amendment to R31 spec, without `.devcontainer/Dockerfile`; the same list with the Dockerfile, as for the kind go. There is no list for `modules` in the template of architecture 3.3 .github: it has only the kinds go and typescript;
  - `.gitattributes`: `modules/** -text`; `text eol=lf` on `LICENSE` and `NOTICE`.

### 3.4. Samples in the eval set

A module's eval set lies in `eval/` of its folder and is not part of the package (section 3.1; R62 spec). The samples in it are of three origins: synthetic samples, samples of adult volunteers and samples of authors (concepts 3 spec). Why the path is such is R04.

For now there are no volunteers' samples, and there cannot be any. The samples page will open the purpose «eval-набор» ("eval set") only when epic E4 in the RUN.md of `spec` approves the text set — the notice and the consent form (R11 canvas). The CONTRIBUTING of `.github` does not accept such samples yet (sources 4.7). The sample file is set by the eval protocol: the message `Sample` of the package `novatrail.eval.v1`, one JSON per sample, without the task's key, with a license for an eval-set sample and with the markup `markup`, and the crops in it are PNG in base64 (R93, R95 spec; [contracts text](../spec/contracts/README.md), section 11). A module's eval set is each `.json` under its `eval/`, and the layout inside `eval/` the protocol left to `modules` ([contracts text](../spec/contracts/README.md), section 13). The layout is flat: each sample is a file `eval/<id>.json`, where `<id>` is its ID, and in `eval/` there are neither folders, nor other files, nor anything starting with `.` (subsection 2.4; the 2026-11-06 and 2026-11-09 amendments to R04).

- **The roles:**
  - **the volunteer** — an adult who writes a sample on the samples page of `canvas` under a signed consent form;
  - **the eval-set organizer** — a maintainer of `modules`. The organizer accepts the form, gets the file, checks it and adds it;
  - **the module's author** adds their own samples themselves, under their own form (the T61 amendment to R04).
- **The path of a volunteer's sample:**
  1. Before the first sample the organizer accepts the signed form (R11 canvas; П-18). The form is `volunteer-form` of version 1, from the text set `eval-set-1` together with the notice `eval-set-notice` 1. The set is in force since 2026-11-10, and the organizer takes it from the file of `spec`'s release `legal-<version>.tar.gz`: `legal-2.tar.gz` since `spec` v0.15.0, and from the line 0.17 on `legal-3.tar.gz`, the legal part of version 3, with the same set (R130, R133, R139 spec; Д162 and Д207 in table 4 of the RUN.md of `spec`, Д232 in its `DEBT.md`). The form lies outside the repositories — in the bucket of R125 spec, under a prefix of its own, with the organizer's access, — and agents are not let near it (П-12; R126 spec). The organizer acts for the operator, ALT EdTech LLC, on its instructions: the controller of the samples and the forms is the operator (R137, R139 spec). In the form:
     - the name, the signature, the date and the version of the text set (R133 spec);
     - the declaration "I am 18 or older, and under the laws of where I live I may make this agreement" (R124 spec);
     - the text of DCO 1.1 word for word and the certification of clause (a);
     - publication under Apache-2.0 — and the consent withdrawal, below (R139 spec);
     - runs at Anthropic and OpenAI in CI of `modules`: a new provider is a new version of the set `eval-set` and a row «в реестр spec» (R133 spec; Д172 in table 4 of the RUN.md of `spec`);
     - samples without names and signatures, and essays on made-up topics (R127 spec);
     - the term of the form.

     The organizer does not take the form of a volunteer from the EU, the EEA or the United Kingdom while the operator has no representative under Art. 27 GDPR, and the form asks the volunteer to state this (R139 spec).
  2. The volunteer opens the samples page and chooses the purpose «eval-набор» ("eval set"). On it they state that they are 18 or older and that the form was accepted, write the sample and save the file on their device (R09 and R11 canvas; architecture 3.6 canvas).
  3. The file and the scan of the form reach the organizer by an upload link valid for at most 12 hours from the session, into the bucket of R125 spec (R126 spec; the 2026-11-13 amendment to R125 spec). This is not the platform's server (architecture 3.2 spec) and not GitHub. It is this link that the page names: «передать файл организатору eval-набора» ("pass the file to the eval-set organizer") (T20 `canvas`).
  4. The organizer checks the file — by the table below, before the file gets into any branch. A file that did not pass is not edited and not added, but erased at once: the volunteer writes the sample anew (R126 spec).
  5. The record of the form, outside the repository, gets the SHA-256 digest of the file. Under which form each sample was added the organizer knows, and the repository does not. The form and the record live as long as at least one sample under it is in the repository or in a released version of a module; a form without an added sample is erased 12 months after signing (R126, R139 spec). The scores of an essay's criteria are given by two people, each on their own, and where they diverge, a third decides, and one mark per criterion goes into the markup (Д172 in table 4 of the RUN.md of `spec`).
  6. The organizer puts the file into the module's `eval/` as `eval/<id>.json`, by the sample's ID, in a branch of the repository itself. The file's name is the sample's ID, and the organizer does not change the content: clause (c) of the DCO holds while the file is byte for byte the one the volunteer certified. The organizer signs off the commit with their own name under clause (c). The commit's author is the organizer too: the DCO check checks the sign-off against the author (sources 4.6). Neither `Co-authored-by` nor the volunteer's data are in the commit's message, the branch's name or the pull request's text. If there are samples of several volunteers, they go in one commit, so that a commit does not show which files come from one person.
  7. The pull request passes CI and is merged with a merge commit (R30 spec). It does not change a package and gives no release: `eval/` is neither in the package nor in `paths` (section 3.3). Once the file is added, the organizer erases its copy in the bucket at once: the file is public, and the copy adds nothing. A file not added is erased by the bucket's lifecycle rule 12 months after the upload, so the organizer adds a file while its form still lives (the 2026-11-13 amendment to R126 spec; Д239 in the `DEBT.md` of `spec`).
- **The path of an author's sample** — under their own form and without an organizer: the author is their own organizer (the T61 amendment to R04). The samples page of `canvas` for the purpose «eval-набор» ("eval set") requires declaring that the consent form was signed and accepted by the organizer, and this purpose has no declarations of the author such as the measurement draft has (architecture 3.6 canvas). So the author signs the same form as a volunteer and accepts it themselves:
  1. **The form.** Before the first sample the author signs the volunteer's form of the eval set — `volunteer-form` 1 from `legal-2.tar.gz` of `spec`'s release v0.15.0 or from `legal-3.tar.gz` of the line 0.17 on, or a newer form, if the purpose opens with it, — and puts it where the volunteers' forms lie, with a record of the form (steps 1 and 5 of the volunteer's path). The form is the same proof of the license and of clause (a) of the DCO as a volunteer's.
  2. **The page.** The author takes the file `samples-page_v<X.Y.Z>.html` from a release of `canvas` in which the purpose «eval-набор» ("eval set") is open, and checks its signature: `just -f spec/justfile release-verify canvas v<X.Y.Z> samples-page_v<X.Y.Z>.html`, with the `SHA256SUMS` and the bundle of the same release beside it. The author opens the page from the file on the computer, in Firefox or Chrome: it does not go to the network (architecture 3.6 canvas). The purpose is open since `canvas` v0.3.11, 2026-11-11: it has the set `eval-set-1` and the tasks of three modules from the release `modules` v0.2.4 (T101 and T102 `canvas`).
  3. **«Начало» ("Start").** The purpose is «eval-набор» ("eval set"). The notice shows the version of the text set, now `eval-set-1`, — it must be the version of the signed form. Then come the page's two declarations: «мне 18 лет или больше, и по законам места, где я живу, я вправе заключить этот договор» and «форма согласия версии … подписана и принята организатором» ("the consent form of version … is signed and accepted by the organizer").
  4. **Writing and markup.** The task is chosen after «Начать» ("Begin"): for the eval set only the tasks of the modules, in a group per module, and for a template the instances by the seeds 1–5 that the build of the page named (T101 `canvas`). «Кто писал» ("Who wrote") is «взрослый своим почерком» ("an adult in their own hand") or «взрослый, подражая детскому» ("an adult imitating a child's"). After writing comes «Разметить» ("Mark up"): on a sheet, the parts of the solution with texts and flags, the first error and the final answer; on an essay, its text as written; on a workbook, each box. Nothing personal is written on the sheet: neither a name, nor a signature, nor a school, nor contacts.
  5. **Saving.** The page shows the crops — what the model will see — and asks to confirm that there is nothing personal on the sheet. The file lands in the browser's downloads: for the author on Linux this is `~/Downloads`, and the folder is not synced with a cloud (T29 `canvas`). The link to the organizer that the page names the author does not need: the author adds the file themselves.
  6. **The check — before any branch.** The author opens the file in the page's section «Открыть образец» ("Open a sample") (T87 `canvas`) and goes by the table of checks below: in the strokes — with the checkbox «показать стёртое» ("show erased"), which also shows what was erased, which lies in the file too, — in the crops and in the markup in words the author looks for a name, a signature, initials, a school and contacts, and in an essay, real names of people and a topic that is not made up (П-20). A file that did not pass is not edited: the author erases it and writes the sample anew. The rows about the fields, the schema and the PNG are held by CI.
  7. **The name and the record.** The page saves the file as `<id>.json`, for any purpose (T101 `canvas`), and it is put into the module as `modules/<module>/eval/<id>.json`, where `<id>` is the field `id` of the file: `jq -r .id <file>` in the author's terminal shows it. The ID of a person's sample is 12 characters of `a–z` and `2–7` (`canvas`'s answer of 2026-11-11). The ID of a person's sample does not start with `synthetic-`: only synthetic samples are named so, and such a file a session of Claude Code may read — `just eval-match` rejects such a sample. Into the record of the form goes the SHA-256 digest of the file, `sha256sum <file>`.
  8. **The commit and the pull request.** A branch of the repository itself, `git commit -s` with the author's own name under clause (a) of the DCO. The sign-off is tied to the sample forever — clause (d), — and the author accepts this by signing. In the commit's message, the branch's name and the pull request's text — nothing about the sample's content. CI — `just eval-check` and `just eval-match`, — then a merge with a merge commit, as in step 7 of the volunteer's path.

  A session of Claude Code in this repository does not open the files of people's samples: `.claude/settings.json` forbids it to read `modules/*/eval/*`, except `synthetic-*.json`. The ban covers its tools for reading files and the Bash commands that name the file — `cat`, `head`, `sed` and others, — but not a program that opens files by itself, so a person's sample is in the working copy only while the author adds it (the T61 amendment to R04). The session's work is CI's output: how many samples are in which module and whether they passed.

  An author without write access adds their samples as a volunteer.
- **Synthetic samples** are drawn by code from a seed (R38 spec; R09 canvas): there is no person behind them, and they need no form. They are written by the recipe `just synthetic <module> <task>:<seed>...` in the devcontainer: it builds the module package with `module-kit build` at this commit and calls the generator of `canvas` of the exact version from the `toolchain` stage, and the generator writes `eval/synthetic-<seed>.json` with the source `modules` and the tag of the newest release of `modules`; a task that the generator does not draw the recipe names with the generator's reason (R29; R110 spec). CI draws each synthetic sample again with the same generator, with the tag from the file, and compares it with the whole file, except the version and the digest of the package, which the contracts text does not compare: a change of another task of the module does not touch the synthetic sample, and a change of its task is a refusal until the sample is redrawn (the T59 amendment to R25). What lies in `eval/` since T60, 2026-11-11, is by the author's decision: a sample per pair "the task's primary skill × band", the task is the first by ID in the pair of those that the generator draws, the seeds are 1…N in the order of the pairs, skill, then band. This is 41 samples of `olympiad-math` (21 MB), 4 of `grade-3-math` (2.1 MB) and, since T69, 3 of `check-board` (1.4 MB, R33); `english-essays` has no synthetic samples — the generator does not draw essays. A new sample takes the module's next free seed. It is added by the author or the organizer, like their own sample: by a commit from a branch of the repository, with their own DCO sign-off. Of the table, the rows about the fields, the schema and the PNG apply to it, and the rows about the form and the names do not.
- **The checks before the merge:**

  | What | Who and when | By what |
  |---|---|---|
  | the form is accepted, and the version of the text set in the file is its version and is approved | the organizer, and for their own sample the author, before the push | the record of the form |
  | the purpose is the eval set; who wrote is an adult, in their own hand or imitating a child's, not a child; there is no writer code; the license is Apache-2.0; there is markup; the task is of its own module | the organizer, before the push, and CI on every pull request, and from a fork too: the purpose, the markup, the source `modules` with the module ID and distinct IDs are held by the eval-set mode `sample-check --eval-set` in `just eval-check`, and who wrote, the writer code and the license by the rules of the sample file for the eval-set purpose ([contracts text](../spec/contracts/README.md), sections 11 and 13; the 2026-11-06 amendment to R93 spec; the 2026-11-09 amendment to R25). That the task is in the module's package and the digest of its open view is the one in the sample — `just eval-match` (the T59 amendment to R25) | the fields of the file (R93, R95 spec) |
  | on the strokes and the crops there is no name, signature, initials, school or contacts, and in the markup neither: the texts of the parts and the essay's truth are published with the sample | the organizer, and for their own sample the author, before the push, by eye | a view of the sample without the network and without AI ([contracts text](../spec/contracts/README.md), section 11) — the samples page of `canvas`, the section «Открыть образец» ("Open a sample") (R26; R110 spec; T87 `canvas`) |
  | an essay has no real names of people, and its topic is made up (П-20) | the same | the same |
  | the file has only the protocol's fields: an unknown field is a refusal. The PNG of the crops has only IHDR, IDAT and IEND, not one ancillary chunk: neither the text chunks tEXt, iTXt and zTXt, nor the metadata eXIf, nor the time tIME, nor others. The file's name says nothing of the person | the organizer, before the push, and CI | the eval protocol (R93 spec); the PNG specification (sources 3.3) |
  | the file lies as `eval/<id>.json`, and `<id>` is its ID; in `eval/` there are neither folders, nor other files, nor anything starting with `.` | CI, on every pull request and push to `main`, and from a fork too | the recipe of `modules` `just eval-check`: a rule of the layout, not of the sample's validity (the 2026-11-06 and 2026-11-09 amendments to R04; R25) |
  | the file is by the eval protocol, by the rules of the sample file ([contracts text](../spec/contracts/README.md), section 11), each file separately | CI, on every pull request and push to `main`, and from a fork too | `spec`'s recipe `just -f spec/justfile sample-check --eval-set novatrail/<folder> <files>` from `just eval-check`, without a key, with the same code as `spec`'s examples (R93 spec; the 2026-11-06 amendment to R04; the 2026-11-09 amendment to R25) |
  | the sample's task is a task of its own module, and the digest of its open view is the one in the sample; a synthetic sample is the one that the generator of `canvas` draws on the commit's package, in everything but the version and the digest of the package | CI, on every pull request and push to `main`, without a key and without the network, and from a fork too | `just eval-match`: `module-kit open` and the generator of `canvas` from the `toolchain` stage over the commit's package (the T59 amendment to R25) |
  | the eval runner reads the samples and runs them, and takes the sample's task from the commit's package | CI, from stage 2 | the eval runner of `tutor`, with the platform's key (R24 spec; R95 spec) |

  - The checks without a key CI runs on any pull request, and from a fork too.
  - By hand — before the push: a sample that got into a branch is already published, and only GitHub Support can remove it from a pull request (sources 2.5). Before the check a sample may carry a name, so neither AI nor an agent looks at it (П-12).
- **What the eval set does not accept:**
  - a sample in a pull request from a volunteer or from a fork: such a pull request is not merged and is closed;
  - a sample that failed the check — and one corrected by the organizer too: erased strokes change the sample on which the markup stands;
  - learners' work and child samples — never (concepts 3 spec; R25 spec).

  If a name is found only after the merge, the sample leaves `eval/` with the next pull request. The organizer asks GitHub Support to remove the cached views and the pull request's references to the commit (sources 2.5) and tells the volunteer about all of it. Clones and forks cannot be removed this way — that is why the checks stand before the push.
- **Consent withdrawal** (R139 spec; Д194 in table 4 of the RUN.md of `spec`). A volunteer or an author withdraws a sample by a letter to the organizer and to the operator, by the name and the date of the form. The organizer carries out the consent withdrawal within 30 days:
  - a sample not yet added is erased;
  - an added one the organizer removes from the module's `eval/` with the next pull request, with their own sign-off and, as when adding, without the person's data in the commit, the branch and the pull request. From that commit the sample is neither in the eval set nor in the module's tests. The pull request does not change the package: `eval/` is not part of the package (section 3.1).

  What has already come out stays: the git history, the releases that came out and others' copies under Apache-2.0 (R139 spec). The record of the form marks the consent withdrawal, and the form lives as long as a sample by it is in a released version (R126 spec).
- **The eval runner in CI** — from stage 2 (R24 spec), only on a pull request:
  - runs the modules whose `eval/` has changed against the pull request's base: an eval set has no published version;
  - runs the modules whose package has changed against the published one (section 3.2) — an edit of the package changes the verdict — and new modules;
  - runs all the modules when a pull request changes the tag of `tutor`: with it the prompts, the models and the runner change. A move to a new version of the format changes the packages of all modules, but the same pull request also changes the tag of `tutor` (repositories 4 spec), so the runner reads the packages of the new version;
  - does not run on a push to `main`: a branch is up to date before the merge (R30 spec), and the same tree lands on `main`;
  - gives a report — the message `Report` of the package `novatrail.eval.v1`, without the content of the work: only the IDs of the samples, places, verdicts and marks. CI of `modules` reads it, and the log of the run with the models' answers is `tutor`'s own ([contracts text](../spec/contracts/README.md), section 14; R95 spec). The run and the published report follow the terms of the provider whose models the report names, and the open log of CI publishes what it prints ([contracts text](../spec/contracts/README.md), section 14; R95 spec; R26). What CI prints, and for which providers, epic E6 decides together with the runner.

  The price of a run is named before the first launch (R24 spec). A job with a key fails on a pull request from a fork (R24 spec), and the samples come from a branch of the repository anyway.
- **What it waits for:**
  - the recipe `report-check` for the eval runner's reports that CI reads came out in `spec` v0.9.12, and `modules` writes no check of a report of its own (the second 2026-11-08 amendment to R95 spec; Д115 in table 4 of the RUN.md of `spec`; R27). The reports will come with the runner, in E6;
  - synthetic samples — `canvas`: the generator of `canvas` so far draws samples only for tasks of `tutor`, with the source `tutor` (R49 canvas), while an eval-set sample of a module is a task of its own module with the source `modules`. `spec` decided at the request of `modules`: the generator on the module package that `module-kit build` built on the commit where the samples land, with the tag of the newest release of `modules` at the source, so a pull request that changes a task's open view remakes its synthetic samples in the same pull request (R110 spec; R26). `canvas` made the generator in task T86 of its epic E5, with `canvas` v0.3.2: sheets with one answer slot `anywhere` and workbook pages, without templates, tapped choices, essays and zones (Д128 in table 4 of the RUN.md of `spec`; В72 in the RUN.md of `canvas`). `canvas` released the generator's file with a signature and the command of an open view's digest at the request of `modules`, in v0.3.6 (R68 `canvas`), the recipe `just synthetic` on it exists (R29), and CI checks the samples against the commit's package — `just eval-match` (the T59 amendment to R25) —, and the first synthetic samples lie in `eval/` since T60 of epic E8 of [RUN.md](../RUN.md). A synthetic sample, like any file that something other than the samples page of `canvas` wrote, writes `canvas_version` `0.0.0` and `text_set_version` `none`, and `sample-check` requires this of synthetic samples (the 2026-11-07 amendment to R93 spec; Д110 in table 4 of the RUN.md of `spec`);
  - from `spec` — nothing: the text of the form, where the forms and their records lie, how long they are kept, the channel, the rejected file and the consent withdrawal were decided by R126, R133 and R139 spec, and the set `eval-set-1` is in force since v0.15.0 (R30);
  - from `canvas` — nothing: the declarations on the page, the words «передать файл организатору eval-набора» ("pass the file to the eval-set organizer") (T20 of its epic E5), the release of the samples page (T75), «Открыть образец» ("Open a sample") (T87), the text sets from the release of `spec` (T89), the modules' tasks in the build (T101, `canvas` v0.3.10) and the move to `spec` v0.15.0 with `eval-set-1` (T102, `canvas` v0.3.11, 2026-11-11) are done, and the purpose «eval-набор» ("eval set") is open. T62 of epic E8 waits for the author's samples, and epic E6 for the volunteers' samples;
  - `canvas` — viewing a sample without the network and without AI: the samples page opens a sample file and shows the strokes of the final geometry, the crops and the texts of the markup, changing and saving nothing. Epic E4 chose it (R26), `spec` recorded it (R110 spec), and `canvas` made it in task T87 of its epic E5, and it is in the page of `canvas` v0.3.4. The authors' samples need it in E8, the volunteers' samples in E6.

**How it was decided** — subsection 2.2 before the outcome of the detailing, word for word (R08). "What is known" was written before the decision, so "`modules` decides" there speaks of the day it was written, and the answer is in this section.

- **What is known:**
  - a sample is made with the same canvas as the learners'. It is saved as a file on the volunteer's device and reaches the eval set through a pull request to `modules`; nothing goes to the platform's server (architecture 3.2 spec);
  - the samples page for the first measurement is decided: without hosting, from a release of `canvas`, and it makes child samples without a pull request — section 1 (R25 spec). The argument of R25 spec against hosting assumes that the volunteer opens the pull request themselves, on GitHub;
  - the recommendation of П-16 is written consent to publication under the module's license, with a warning that it cannot be withdrawn; no names and no signatures on the samples. The license of the samples is the module's, Apache-2.0 (R20 spec);
  - the eval sets of the first slice are of volunteers' samples: on them the reading thresholds and the target of assessment by a rubric that `tutor` set on synthetic samples and the author's sheets are repeated (the 2026-10-28 addition to П-16; R67, R68, R71 spec). If an own reading model on the device is to learn on the samples, the consent names that too (ibid.);
  - what lies in a sample from the samples page and how one consents to it is the open question П-18. It proceeds from what `canvas` has already decided (R09 canvas, R11 canvas): an eval-set sample, which is published forever, holds the final geometry of the visible strokes, without time, pressure and the order of writing. The page takes no consent, and the organizer accepts the signed form before the first sample;
  - for the module of essays a sample is a whole text by hand. What the module may ask a learner under 18 to write about is П-20, and the volunteers' samples for its eval set are published forever too (ibid.);
  - learners' work does not enter a module. If it is taken into eval, it lies outside the repositories and only with a separate consent (concepts 3 spec; П-10). Development agents get no access to it, and the volunteers' samples lie in the modules openly (П-12);
  - the eval set lies in `eval/` of the module's folder and is not part of the package: the eval runner takes the samples from the module's folder in CI of `modules` ([format text](../spec/format/README.md), section 10; R62 spec). The format of the samples is by the eval protocol of `spec` (R15 spec; architecture 3.2 spec). The protocol and the format of the samples are epic E3 in the RUN.md of `spec`, row 1.6 of [COORDINATION.md](../spec/COORDINATION.md);
  - a pull request is a contribution, and each commit of a contribution is signed off under the DCO. The sign-off is public forever, so CONTRIBUTING allows signing off with a GitHub noreply address (R22 spec). The sign-off of each commit of a pull request, except merge commits, is checked by the shared DCO check, and it holds the merge. The ruleset of `main` allows only a merge commit, so a pull request's commits land on `main` as they are (R30 spec);
  - the eval runner of `tutor` runs a module's eval set in CI from stage 2 (architecture 3.10 spec) — with the platform's key at the AI provider, and on a pull request from a fork its job fails until the change is moved to a branch of the repository (section 1; R24 spec; Т-28).
- **Options** (T02; the first in each point was chosen; the sources: sources 2.4, 2.5, 3.2, 3.3, 4.6, 4.7, 4.13, 4.15 and 4.18):
  - who adds a volunteer's sample:
    - the eval-set organizer — a maintainer of `modules` — with their own commit under clause (c) of the DCO, and the volunteer certifies clause (a) in the form;
    - the volunteer themselves, from their own GitHub account;
    - both ways.

    Clause (c) of the DCO works only if the volunteer has certified clause (a) themselves, and the DCO check checks the sign-off against the commit's author. A GitHub account is allowed from age 13, and a push publishes at once: only GitHub Support can remove a commit from a pull request that refers to it. П-18 asks how the DCO sign-off will tie a sample to an account, and names the option "the maintainer adds them themselves, under clause (c)". `canvas` left this choice to T02 (R11 canvas);
  - where the consent lies and what of it is in the repository:
    - the form outside the repositories, nothing in the repository, and the record of the form names the digests of the samples added by it;
    - the form outside the repositories, and the sample in git carries a random number of the form.

    An eval-set sample has no writer code (R09 canvas);
  - how to check that the volunteer is an adult:
    - the declaration "I am 18 or older" in the signed form and on the samples page (R11 canvas; П-6);
    - the same, and a document that the organizer looks at;
  - the authors' samples:
    - the author adds their sample themselves, with their own DCO sign-off under clause (a), without a form;
    - by a form, like a volunteer's sample;
  - who looks for a name on a sample, and when — `canvas` left this to T02 (architecture 3.5 canvas):
    - the organizer, by hand, before the first push;
    - the review of the pull request;
    - AI or an agent;
  - which modules the eval runner runs — T01 left this to T02 (R03):
    - those with a changed `eval/` or package, and all of them at a new tag of `tutor`;
    - all of them on every pull request and on a push to `main`.

**How it was decided** — subsection 2.4, word for word (R09).

- **What is known:**
  - the ID of a sample is 1–64 lowercase Latin letters, digits and hyphens, "so that it can name the sample's file". For a person's sample it is random, and for a synthetic one it is derived from the seed, so one seed gives one file. Two samples of one set do not share an ID ([contracts text](../spec/contracts/README.md), sections 11 and 13);
  - the file's name must say nothing about the person, and the organizer gives it the name without changing the content (R04; section 3.4);
  - the samples of several volunteers go in one commit, so that the commit does not show which files are from one person (R04);
  - a pull request that changes a task's open view remakes or removes its samples in the same pull request: the runner rejects a sample whose open view digest did not match (R95 spec; the 2026-11-06 amendment to R04), and until the runner `just eval-match` holds this (the T59 amendment to R25);
  - there are no samples in `eval/`, and there will be none before the volunteer's form, except synthetic samples and the authors' samples (section 3.4).
- **Options:**
  - flat, by the sample's ID — `eval/<id>.json`. The name is derived from the file and says nothing about the person, and there is nothing to break the layout with. A task's samples are found by the task's field in the file or by the runner's refusal, which names them. So row В8, «в реестр spec», recommended before the protocol (R09);
  - by task — `eval/<task id>/<id>.json`. The samples of a task whose open view a pull request changes are seen by the folder, but the path repeats a field of the file and may diverge from it, and then one more check of CI is needed;
  - by origin — `eval/synthetic/`, `eval/volunteers/`, `eval/authors/`. The path repeats the field «кто писал» ("who wrote") and divides people's samples into groups that are easier to match against one another.
- **Recommendation:** flat, `eval/<id>.json`, and CI checks that the file's name is its ID: the path repeats no field of the file but the ID, and there is one check. The samples of a task whose open view a pull request changes CI names anyway — by the runner's refusal.

### 3.5. The port of standalone

Olympiad maths `novatrail/olympiad-math` comes from standalone at commit 89986bc, 2026-11-01 (sources 5.2). Why so — R05. The rest of the port takes commit 42b4408, 2026-11-07, which also has the tasks that are solved from the picture. It has 753 tasks (the 2026-11-08 amendment to R05).
- the first batch — section 3.6 (R07);
- the rest of the port — row 3.7 of [COORDINATION.md](../spec/COORDINATION.md);
- the port's code — the converter `tools/port`, since T11 (R12): the batch's tasks by this table, the check of the key and the trap values by standalone's test bench, the pictures `row` as PNG and the overlay files, which hold what a person writes. The recipe `just port <module>` clones standalone at the batch's commit into a temporary folder and writes its tasks. Since T18 the converter reads 42b4408 as well and writes an answer of any kind, and the recipe `just port-survey <folder>` runs all the tasks that no batch names into a temporary module (R20). Since T19 and T20 it also draws `clock`, `table`, `bars`, `number_line`, `calendar`, `containers`, `grid`, `piles`, `venn`, `balance` and `ring` (the 2026-11-08 amendments to R12). Since T21 — `lattice`, `tiles`, `shapes` and `stairs` too, the kinds of the tasks that are solved from the picture, and the `alt` of such tasks — by the rule below (R21). Only `flags` is not drawn: no task has it.
- the pictures of lessons — the data `tools/port/pictures/<module folder>/<name>.json` in standalone's schema of pictures, which the converter draws into `media/lessons/<name>.png`; the recipe `just lesson-pictures <module>`, and CI checks them with `--check` (the T41 amendment to R12).
- the pictures of the instances of templates — the sources `tools/port/template-pictures/<module folder>/<picture path without placeholders>.json` in standalone's schema of pictures, with placeholders in the strings and, where the picture changes not only by numbers, with variants `when`, or the function of the template's solver, `{"solver": "figure"}` (the T53 amendment to R12); the converter computes the template's instances itself — `go.starlark.net` at the version of the format's dialect — and draws the picture of each into `media/pictures/`, all of one size, which the template sets; the recipe `just template-pictures <module>`, and CI checks them with `--check` (the T52 amendment to R12).

- **What is in standalone** (a count by jq, 2026-11-02):

  | Part | Where | How many |
  |---|---|---|
  | topics and their links `builds_on` | `content/catalogs/topics.json` | 17: eight in all three bands, two in 3–4 and 5–6, seven only in 5–6 |
  | traps | `content/catalogs/traps.json` | 20, each with a description and advice to an adult |
  | techniques and Pólya's steps | `site/content/en/techniques.yaml`, `site/data.json` | 11 techniques in three groups |
  | tasks | `content/examples/<topic>.json` | 705: 200 in band 1–2, 250 in 3–4, 255 in 5–6. All fixed, all in English |
  | the tasks' solvers | `content/examples/solvers/<id>.star` | 705, one per task |
  | templates | `content/solvers/<topic>/*.star` | 33 |
  | pictures | the task's field `picture` — a JSON description that standalone's card draws; the kinds — `content/pictures/` | in 104 tasks, 11 kinds of 13; a clock in 36. A solution picture in 35 |
  | hints | the task's field `hint` | in 263 tasks: in all 255 of band 5–6 and in 8 calendar tasks |
  | the site's topic pages | `site/content/<language>/topics/` | 17 in each of 10 languages, `en` and `ru` among them |
  | ideas | `content/ideas/<topic>.json` | 552 |
  | excluded skills | `content/catalogs/skills.json` | 25 |

  Every task's answer is a choice of five options. Each of the four wrong options has a trap and an explanation. By the text of the correct option, the answers are these:
  - an integer — 496. Of them at least 17 have a question that points to the options: "Which of these numbers…";
  - a number with a unit — 12: "8 km", "16 cents";
  - other text with digits — 44: dates, "9 hours 30 minutes", "bag 3", "5 with 12 left over";
  - a fraction — 7;
  - a percentage — 5;
  - a time on a clock — 30;
  - "yes" or "no" — 1;
  - words — 110.
- **At commit 42b4408**, 2026-11-07, — what the rest of the port takes (`git diff 89986bc 42b4408`, a count by jq, 2026-11-08; the 2026-11-08 amendment to R05). Standalone added tasks that are solved from the picture. The data of such a task stand only on the picture, and its solver reads the picture's description, `solve(options, picture)`. The rest has not changed: the traps, the techniques, the excluded skills and the JSON of all 705 earlier tasks are the same as at 89986bc.

  | Part | At 89986bc | At 42b4408 |
  |---|---|---|
  | topics | 17 | 19: the new `counting.figures` and `geometry.cubes_paper`, both in three bands, `builds_on` empty. Both have `solved_from_picture: true` |
  | tasks and their solvers | 705 | 753: 214 in band 1–2, 267 in 3–4, 272 in 5–6. Of the 48 new ones, 30 are in the two new topics, five per band. Another 18 are in five earlier ones: `arithmetic.tricks` 2, `combinatorics.enumeration` 6, `counting.gaps` 4, `geometry.grid` 4, `parity.alternation` 2. Each of these 18 has the field `idea` — the place of its idea in the topic's list of ideas |
  | templates | 33 | 39. The solver of each of the six new ones reads the picture: `shapes-rows`, `sticks-square`, `squares-on-lines`, `triangles-on-lines`, `fold-net` and `fold-sheet` |
  | pictures | in 104 tasks, 11 kinds of 13. A solution picture in 35 | in 152 tasks, 15 kinds of 17. A solution picture in 83. The new kinds are `lattice` (in 38 new tasks), `tiles`, `shapes` and `stairs` (in two each); two more have `piles`, and two `row` |
  | hints | in 263 tasks | in 311: in every new task, in all three bands |
  | the site's topic pages | 17 | 19: the new `counting-figures` and `cubes-and-paper` |
  | ideas | 552 | 634 |

  The correct option of all 48 new tasks is an integer, and none has the question "Which of these…". One explanation names the key — that of `gaps-12-d2-6`, option 6: "The steps were not divided by 3." The three in this task is a jump on the picture, not a number of the statement.
- **The `alt` of a task solved from the picture** (R21). A task is solved from the picture if its standalone solver declares `solve(options, picture)` — at 42b4408 there are 48 such. Its data stand only on the picture, and a screen reader reads the `alt`. So:
  - the `alt` says in words every fact that the solver reads: the size of the field, every line, wall and cell, every number and label, every stair and jump;
  - the `alt` names neither the key nor a step of the solution;
  - the converter checks that the `alt` names every text that the picture writes — a label, a number on a tile, the total of a row of figures — and refuses if not. A label of capital letters it looks for case-sensitively. The geometry is held by proofreading;
  - for any task, the `alt` of the statement's picture does not name the key unless the question names it or the picture itself writes it: otherwise the leak check would count the key as shown and would not look for it anywhere ([format text](../spec/format/README.md), section 20). This is a check of the converter too. If the picture writes the key, `leak_phrases` carry its giving away.
- **What becomes what** ([format text](../spec/format/README.md), sections 10–24):

  | In standalone | In the module |
  |---|---|
  | topic | a skill of `skill-graph.yaml` with the assessability `verdict`, as in the format example `olympiad-math`. The ID is the topic with `.` and `_` replaced by `-`: `counting.gaps` → `counting-gaps` |
  | `builds_on` | the skill's `prerequisites`. As in the example, `number-divisibility` also has the anchor skill `core.division-with-remainder`, and `fractions-parts` — `core.fractions-of-a-whole` (R50 spec). The tiers are computed from the graph, and their names are Foundations, Rules and Hard problems (R49 spec) |
  | the topics' grades | `audience.grades` of the module manifest, 1–6, and the bands `bands`: `grades-1-2`, `grades-3-4` and `grades-5-6` |
  | trap | an entry of `traps.yaml`, as in the format example, and its own `twelve-hours-off` — the value of a 12-hour dial for the kind `time`, until the example takes it (the T50 amendment to R23; В35 in table 3 of [RUN.md](../RUN.md)). The ID — with `_` replaced by `-`; `advice` → `trap_advice` |
  | a technique of the site | an entry of `techniques.yaml`, as in the format example: `title` and `description`. The techniques' numbers and groups are not carried over — the format has none |
  | Pólya's steps | `solving-steps.yaml`, as in the format example, with the descriptions of two steps of the four (the T54 amendment to R13). The example takes them from the coach of standalone's prototype, and standalone's page of techniques opens with the same four steps from Pólya's book "How to Solve It" |
  | task ID `gaps-12-d1-1` | the same ID and the file `tasks/gaps-12-d1-1.yaml` |
  | `topic`, `grade_level`, `difficulty` | `skills.primary` — the topic's skill; `band` — the band; `difficulty` — the same, 1–5 within the band. In `supporting` — an anchor skill, if the task is on fractions or divisibility |
  | `question` | `statement` — the text as it is. It has no formulas, and the format lets the signs × and ÷ and the fractions `a/b` through |
  | the correct option is an integer, a number with a unit or a fraction | one item with an answer slot `anywhere` and the check `number`, for a number with a unit — with the batch's `units`, and for a fraction — with the check `fraction`. The answer is written by hand. A percentage is a number with the unit `%`, "44%": step 6 of the [format text](../spec/format/README.md), section 15, drops it after the number, with a space and without (R20). The exception is a question that points to the options: without them it is not posed, and the task gets `choice` by tapping |
  | the correct option is a time on a clock `H:MM` | an answer slot `anywhere` with the check `time`, without `options`: the learner writes any time, the hours and minutes are compared as numbers, and each wrong option is a trap value ([format text](../spec/format/README.md), section 15). A key after 12:59 and before 1:00 also has the value `twelve-hours-off`: the same time on the dial, 9:10 for 21:10 and 12:15 for 0:15. There are 30 such tasks, all `clk-*`, and five of them have the dial's value (the T50 amendment to R20) |
  | the correct option is one word, other than "yes" and "no", that stands in the question | an answer slot `anywhere` with the check `choice`, which has its own `options`, all five options, and `notation: word`: the learner writes the answer by hand and does not see the options, any one word is read, a word not among the five options is `not_yet` with a hint rung, and a full stop at the end is dropped ([format text](../spec/format/README.md), section 15). There are 44 such choices — 34 in `ord-*`, 7 in `par-*` and 3 days of the week in `cal-*` (the T49 amendment to R20). The paragraph `answer_form` of the overlay may say in the statement how to write the answer (R20). If the wrong options are arguments, not answers, the batch sets `tapped: true` on the task, and the choice is by tapping (the 2026-11-08 amendment to R20) |
  | the correct option is something else: a word that is not in the question, a date, "yes" or "no", a quantity of several parts, a phrase, a number with words for which the batch gives no `units` | an answer slot `choice` with five options, by tapping: without the options such an answer cannot be written, and the checks of format v0 do not judge dates, quantities of several parts and text (R20) |
  | a wrong option with a trap and an explanation | a trap value `trap_values`, where the trap is with `-` instead of `_`, and the text is the option's explanation. The value is written like the key: for `number` and `fraction` — the option's number in the kind's notation, without a unit, no longer than 16 characters, and for `choice` — the option's text as it is. A non-numeric wrong option of a task with a numeric answer cannot become a trap value, and its trap goes into the task's `traps` — AI will suggest it |
  | an explanation that names the correct answer | is rewritten in `trap_texts` of the overlay. The key of a number and of a fraction is looked for by `task.no-leaks`; a key that also stands in the statement it does not look for. With an integer key there are 159 such explanations, in 119 tasks, and with a fraction or a number with a unit 5 more. The key of a tapped choice the check never looks for, and a word that is written by hand stands in the question, so the explanations of a choice — and `givens`, `unknown` and `simpler_case` of the overlay — the converter checks itself: it finds the key's phrase in the normal form of the [format text](../spec/format/README.md), section 20, and refuses (R20). The run of 2026-11-08 found 94 such explanations in 57 tasks, and `task.no-leaks` 128 violations in 101 tasks |
  | `solution`, `solution_total` | `solution`, and the total as its last line |
  | `picture`, `solution_picture` | a PNG in `media/`, drawn from the description: with `alt`, width and height, on an opaque light background (R51 spec). The solution picture — in `solution` |
  | the task's solver | does not go into the module: the key of a fixed task is the author's (concepts 9.4 spec). The port runs the solver once and checks the key and the trap values against it |
  | template | a task `tasks/<template name>.yaml` with `template` — parameters, derived numbers, a constraint and a solver (R60, R61 spec) — written by hand by R22. The parameters are the numbers of the reference task's statement, and its numbers are one of the instances; the solver is rewritten for the format's dialect: without uppercase constants, without the helpers `permutations`, `product` and `gcd`, without `sum` and `**` and without recursion. A template whose solver reads the picture is carried over only if the picture can be built from the numbers it shows: the solver's function `figure` builds it from the parameters, `solve` computes on it, and the port draws the same one — so `shapes-rows`, `squares-on-lines` and `sticks-square`, 52 instances; `fold-net`, `fold-sheet` and `triangles-on-lines` stayed tasks (the T53 amendments to R12 and R22). The reference task's picture is the picture of each instance: a path with placeholders of the parameters, a file for each instance, which `just template-pictures` draws. So `for-every`, `split`, `groups` and `two-groups` have a picture of the statement, and `changes` and `groups` a picture and the last line of the solution, 857 files for 483 instances: all the instances, by the author's decision (the T52 amendment to R22) |
  | the site's topic page | a lesson `lessons/<id>.yaml` (R58 spec). The gist of the topic and the method become an explanation, the examples worked examples, and their steps parts of the demo. An example that solves a task of the module does not go into the lesson ([format text](../spec/format/README.md), section 21). The language of the lesson is the module's language, and the page is taken in it |

- **What is not carried over:**

  | In standalone | What becomes of it |
  |---|---|
  | the task's hint, `hint` | the leading question `leading_question` at the root of the task, after `unknown`, word for word: `api` gives it on "Hint" before the first wrong answer, and `task.no-leaks` looks for keys in it ([format text](../spec/format/README.md), sections 14 and 19; R108 spec). 311 tasks have a hint — 22 in band 1–2, 17 in 3–4 and 272 in 5–6 — and 53 of them are not questions but instructions, which also break the task down. The key of a choice and of a time in it the converter looks for, and the key of a number and of a fraction `task.no-leaks`, except a key that the statement shows: a question with it shows no more than the statement. Where the question names the key, the field `leading_question` of the overlay rewrites it — five are rewritten so (the T51 amendment to R12). Eight templates whose reference task carries a hint keep it with placeholders (the T51 amendment to R22). The request for the field went as row В31 «в реестр spec» from the outcome of epic E3, T43 (section 3.7) |
  | the techniques' numbers and groups | the format has none; the group is seen from the technique's description |
  | the kinds of pictures, 13 | a picture is a file (R51 spec), so the kinds are not needed |
  | ideas, 552, at 42b4408 — 634 | there is no generation on the fly (R60 spec) |
  | the task's field `idea` and the topic's `solved_from_picture`, at 42b4408 | the first is the place of the idea in its list, and the ideas are not carried over. The second tells standalone in which chat the task may be given: only where the card is drawn. In a module the picture is part of the statement, and every app shows it |
  | the excluded skills, 25, the limits of readability and the excluded genres | these are rules for whoever writes tasks. standalone's texts already follow them, and new tasks need them in the guide for authors |
  | the page's parts `home` and `ask` | they are for a parent, not for a lesson |
  | the site's groups of topics, 6 | the graph has none |
  | the solvers of the site's examples | they prove the examples' answers; in a lesson an example carries the author's answer |
- **The license.** standalone's content is under MIT (sources 5.1), and the module carries its notice in `NOTICE` (R12 spec; the 2026-10-28 amendment to R46 spec). `NOTICE` follows the format's package vector `with-solver`:
  - a line saying that the module holds standalone's content, and the text of MIT with "Copyright (c) 2026 MathTrail";
  - a line saying that part of the tasks and Pólya's steps came from standalone's prototype under Apache-2.0 and are changed here. This is section 4(b) of Apache-2.0: "You must cause any modified files to carry prominent notices stating that You changed the files". The copy of Apache-2.0 is carried by the module's `LICENSE`, as in every module (section 3.1).

  Neither standalone nor the prototype names the copyright holder, except for the word "MathTrail", though they have one author on GitHub. So `NOTICE` meets the terms of both licenses, and does not rest on there being one author.
- **Provenance** — by sources 5.2 and 5.7:
  - 442 tasks of grades 1–4 standalone took word for word from its prototype. The prototype is under Apache-2.0, without `NOTICE` and without a filled-in copyright line;
  - the wording of the other tasks is its own. The ideas for grades 5–6 are taken from MathNet under CC BY 4.0 and from books and collections, but the wording is not. standalone's tasks have no field "source";
  - a task of the module has no field of provenance either (concepts 9.4 spec). For О-10 `spec` recorded: others' sets for this module were checked by the records of standalone and the prototype, there was no comparison of the wording with the sources, and the module no longer waits for О-10. For the lawyer at О-10 — attribution for an idea or a kind of task, and wording that a model wrote. The sources of the ideas can be named voluntarily only outside the module: `NOTICE` holds only what the licenses require (the additions to О-10 and product 14 spec; [format text](../spec/format/README.md), section 10; row В11, closed — R09).

**How it was decided** — subsection 2.3 before the outcome of the detailing, word for word (R08). "What is known" was written before the decision, so "`modules` decides" there speaks of the day it was written, and the answer is in this section and in section 3.6.

- **What is known:**
  - carried over are 17 topics and the links `builds_on`, 20 traps, 11 techniques and 603 reference tasks with solvers — from `content/` and `site/` of standalone. The site's topic pages — the gist of the topic, the technique's steps and examples worked through to the answer — are the source of explanations and worked examples for the lessons (repositories 7 spec; R58 spec);
  - first comes the first batch: 20–30 tasks of standalone across three bands, without templates and solvers, with lessons, with pictures as files of `media/`, and the module manifest names `adult` in the audience too. `modules` decides the batch's make-up and the slot bindings of its tasks. The rest of the port comes after it, row 3.7 of [COORDINATION.md](../spec/COORDINATION.md) (R76 spec);
  - how the parts of standalone fit format 0.2:
    - topics are skills of the graph `skill-graph.yaml`. A tier is computed from the graph, and the author gives the tiers their names; a grade is the label `audience.grades` of the module manifest (R49 spec; О-18 decided). The graph rests on school skills not through skills of school maths but through anchor skills of the shared taxonomy, which school maths teaches (R50 spec);
    - tasks are `tasks/<id>.yaml`. standalone's three bands are the bands `bands` of the module manifest, the choice of five options is the slot binding `choice`, and pictures are files of `media/` with a description and a size (R51 spec), on an opaque light background (the second 2026-11-01 amendment to R51 spec). A check is at each answer slot (R56 spec);
    - traps and techniques are `traps.yaml` and `techniques.yaml`; standalone's traps and techniques already stand in the format examples (R57 spec; repositories 7 spec). There too are the solving steps, the tasks' solutions and the leak phrases;
    - lessons are `lessons/<id>.yaml`, and the graph sets the order of the lessons (R58 spec; О-19 decided);
    - standalone's IDs — `logic.knights_liars`, `off_by_one` — are carried over by replacing `.` and `_` with `-`, with a check for collisions (R46 spec; Д22 in table 4 of the RUN.md of `spec`);
  - templates and solvers are for the rest of the port, not for the first batch: a template is parameters, derived numbers and a constraint in a subset of Starlark (R60 spec), a template's solver is Starlark with the dialect and the step budget of standalone's solvers (R61 spec), and the generator and the places of randomness come from `spec` (R59 spec; Т-18 closed). The format example `number-thought-of-template` is a solver of standalone (repositories 7 spec);
  - the license: MathTrail's code under MIT is carried over freely, and the MIT notice is kept (R12 spec). The repository standalone is under MIT: at its root a single `LICENSE`, "MIT License, Copyright (c) 2026 MathTrail" (sources 5.1). The MIT notice in a module stands only in its `NOTICE`, and the package carries it in the device part (the 2026-10-28 amendment to R46 spec; R62 spec). It is not decided whether the module keeps this notice or MathTrail's copyright holder releases the content under Apache-2.0 (R46 spec). Nor is it found out where standalone's tasks come from and whose wording they have: others' sets of tasks, MathNet among them, are checked before use, and standalone's rule is "ideas — yes, wordings — our own" (О-10);
  - the port fits format 0.2. What the format does not express is a request to the next version 0.y: the requests of `modules` and `module-kit` will bring it after the first content of the three modules, with T38 of `spec` (Д60 in table 4 of the RUN.md of `spec`);
  - MathTrail's repositories are cloned into a temporary folder, not into the repository (RULES.md, Citations).
- **Options** (T03; the first in each point was chosen; the sources: sources 5.1, 5.2 and 5.7). "603" above is the number before R309 standalone, and at 89986bc there are 705 tasks (section 3.5).
  - how standalone's content enters the module: the MIT notice in `NOTICE`; a re-release under Apache-2.0 by its copyright holder. standalone has one human author, and the prototype, from which 442 tasks were taken, is under Apache-2.0;
  - the task ID: the same as in standalone — `gaps-12-d1-1` already fits; a new one, by meaning;
  - a number answer: by hand, an answer slot `anywhere` with the check `number`, and the wrong options as trap values; a tapped choice of five, as in standalone (Д32 in table 4 of the RUN.md of `spec`). R76 spec requires handwriting from day one;
  - the task's solver: check the key against it at the port and not take it into the module; take it into the module — the format has nothing to hold it with, as it has a solver only in a template (R61 spec);
  - what the format does not express: the fate of each point in a table, and as a row «в реестр spec» only what the first batch or the port needs; a row for each point;
  - whether to split the task: T03 — what becomes what and the license, and the first batch a new T07; one task;
  - the tag of `spec`: T03 at v0.6.3, and the move to v0.6.5 a new T08 before the first batch; the move before T03; the move in the content epic.

### 3.6. The first batch

The first batch of `novatrail/olympiad-math` is 27 tasks of standalone across three bands and 9 lessons, for the first application, where the author solves tasks locally in a browser (R76 spec). Why so — R07. The tasks and lessons follow the table of section 3.5 and this section. The content epic writes the module's files, row 3.6 of [COORDINATION.md](../spec/COORDINATION.md). The list is counted by standalone 89986bc and the text of format 0.2 (sources 5.5), and the batch is written in format 0.3: it changes Starlark, version comparison and the index, and the batch has neither templates nor predicates, so the list is the same (R101 spec; R09).

- **Skills** — `counting-gaps`, `combinatorics-enumeration` and `pigeonhole-basic`. These three topics of standalone run in all three bands, their answers are integers almost everywhere, and the site's topic page has an example for each level. Tasks with a number answer, without "Which of these…" and without a picture, they have 52, 64 and 63 of 65.
- **The rule of choice:**
  1. a fitting task has an integer as its correct option, with a unit or without, like "50 m", and its question does not point to the options. It has no picture, except a picture `row` under point 5. An answer with a unit is written the same as without it: the key is a number, and the unit is in `units` (section 3.5);
  2. for each skill and band three tasks are taken — of difficulty 1, 3 and 5;
  3. at each difficulty the first fitting task by ID is taken whose explanations of the wrong options do not name the key, and if there is no such one — the first fitting one. A task is not taken if its key:
     - stands in the text of its lesson — `lesson.no-leaks`;
     - stands in the statement of an earlier task of the lesson;
     - equals the key of an earlier task of the lesson: that task's solution ends with its key, and a lesson must not name tasks the texts of one of which give the answer of another ([format text](../spec/format/README.md), section 20).

     If there is no fitting task at this difficulty, the nearest difficulty is taken, higher first;
  4. the key is looked for in a text by the format's rule of finding numbers ([format text](../spec/format/README.md), section 20): by value, so `7` is in "7 rungs" but not in "17";
  5. for `counting-gaps` in bands 1–2 and 3–4 one task of three has a picture `row`: at the first difficulty where such a one fits, it is taken instead of a task without a picture. There are 10 fitting tasks with `row`, all in these two bands.

  The difficulty shifted for one task: for `counting-gaps` 1–2 the keys of the fitting tasks of difficulty 1 stand in the lesson's text, and `gaps-12-d2-1` was taken.

  In three lessons a number from the solution or from an option's explanation of an earlier task equals the key of a later one:
  - `enum-56-d1-1` → `enum-56-d3-2`, 16 — in the solution and in the explanations of options 36 and 576;
  - `pig-12-d3-1` → `pig-12-d5-1`, 10 — in the solution and in the explanation of option 10;
  - `pig-56-d1-1` → `pig-56-d3-2`, 19 — in the solution and in the explanation of option 19;
  - `pig-56-d3-2` → `pig-56-d5-1`, 11 — in the explanation of option 11.

  This is a coincidence of numbers, not an answer, and the rule does not look at these texts: otherwise for `combinatorics-enumeration` 5–6 both tasks would have moved to difficulty 4. The content epic rereads these pairs when it writes the solutions and explanations, and changes the number in an explanation if it reads as the answer of the next task.

  The two paragraphs above are about the first batch. In epic E3 the rule is the same without points 1 and 5: any task of the skill and band fits, except a template, with an answer of any kind and with any picture. A number key is looked for by value, by point 4, and a key with a word in it — as a phrase, case-insensitively, without quotes and a sign at the end, and with word boundaries. A task whose solution gives the answer of a later task of the lesson not by a number but by the course of solution the rule skips at the executor's word (the T39 amendment to R14). If no fitting task remains at a later difficulty, the rule takes at an earlier one the next in its order, and the lesson's tasks go by difficulty. The lessons `counting-figures` and `geometry-cubes-paper` rewrite the cards and examples of the pages that solved tasks of the module, and show the pictures `media/lessons/` that the converter draws from `tools/port/pictures/<module folder>/` (the T41 amendments to R12 and R14). The preference of point 3 changes nothing in E3: the port has already rewritten the explanations that name the key. The difficulty shifts and the pairs with coinciding numbers of each lesson of E3 are recorded in «Сделано» of its task in [RUN.md](../RUN.md) (the 2026-11-08 amendment to R14).
- **The graph** — three skills with the assessability `verdict`, as in the format example `olympiad-math`:

  | Skill | Prerequisites | Tier |
  |---|---|---|
  | `counting-gaps` | — | 1, Foundations |
  | `combinatorics-enumeration` | — | 1, Foundations |
  | `pigeonhole-basic` | `combinatorics-enumeration` | 2, Rules |

  The tier names are `[Foundations, Rules, Hard problems]`, as in the example: the third names nothing until the port adds tier 3 ([format text](../spec/format/README.md), section 12). The graph is not a chain, and the module is not sequential (the 2026-11-02 amendment to R58 spec). There are no anchor skills. The port adds the other 14 skills.
- **The module manifest** — like the format example `olympiad-math.yaml`, except for three fields: `audience.age_categories` — `[under_13, adult]` (R76 spec), no `dependencies`, and the version the content epic writes. The grades are 1–6, the bands `grades-1-2`, `grades-3-4` and `grades-5-6`, the explanation language `en`, `rung_templates` from the example, and there are no numbers in them.
- **The catalogs** — from the format examples:
  - `traps.yaml` — 10 traps that the batch's tasks name: `answered-other-question`, `best-case-not-worst`, `double-count`, `ignored-condition`, `ignored-not`, `missed-case`, `number-from-text`, `off-by-one`, `stopped-early`, `wrong-operation`. None of them names skills, and there are no numbers in their advice;
  - `techniques.yaml` — 4 techniques: `draw`, `small-numbers`, `systematic`, `pigeonhole`;
  - `solving-steps.yaml` — all four steps of the example.
- **A task of the batch** — fixed, with the same ID and the file `tasks/<id>.yaml` (R05):
  - one item: an answer slot `anywhere`, the check `number`, the key as a string and `units`. The question of every task of the batch counts objects or a measure, so every one has `units`: the word that the question counts, in the singular and the plural, and for a measure also the abbreviation. Without them "14 metres" by hand would not be in the number's notation ([format text](../spec/format/README.md), section 15);
  - the wrong options are all numbers, and all become trap values `trap_values`: the value is the option's number without a unit, and the text is the option's explanation. `traps` from the port is empty;
  - techniques — by skill: those that standalone's site links with the topic (`site/data.json`, the member `techniques`), and for `counting-gaps` also `draw` — the method of its page "Draw a small case and count on it" and the techniques of the task `gaps-ladder-1` of the format example. In total: `counting-gaps` — `draw` and `small-numbers`, `combinatorics-enumeration` — `systematic`, `pigeonhole-basic` — `pigeonhole`;
  - every task has the item's `givens`, `unknown` and `simpler_case`. Without them the hint before the first wrong answer and rung 3 are left without words (R82 spec; [format text](../spec/format/README.md), section 19). The content epic writes them;
  - standalone's hints of nine tasks of band 5–6 are questions that break the task down, that is, a path to the answer. In format 0.3 they had no place, and the batch did not take them, and from 0.4 they are leading questions `leading_question`, given on "Hint" before the first error (the T51 amendment to R07; [format text](../spec/format/README.md), section 14);
  - `solution` — standalone's solution, the total as the last line (R05). The port checks the key and the trap values against the task's solver;
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

- **Lessons** — a lesson per skill and band, ID `<skill>-grades-<band>`, and the lesson's skill is its skill. The title is the title of the topic page and the band's words, without numbers ([format text](../spec/format/README.md), section 21): "Gaps and boundaries: first steps" in band 1–2, "…: going further" in 3–4 and "…: olympiad level" in 5–6. The band's title in the module manifest, "Grades 1–2", does not go into the lesson's title: `lesson.no-leaks` checks the numbers in it against the keys of the lesson's tasks. The segments in order:
  1. an explanation: the page's `hero.lead` and cards `basis` — the gist of the topic and its rules;
  2. a worked example of its own level: the statement — `question`, a handwriting sheet with a grid, the demo — `steps` and `answer`;
  3. an explanation: the example's note, `note`, if it has one;
  4. the lesson's three tasks in increasing difficulty;
  5. an explanation: the paragraph "How to solve: …" — the page's `method.title` — and its steps `method.steps` as a part `numbered_list` (the T54 amendment to R14).

  There are no pictures in the lessons. No example solves a task of standalone: the tasks nearest in wording ask something else and have another answer. The lessons of `pigeonhole-basic` come after the lessons of `combinatorics-enumeration` by the graph. The format puts every lesson of `pigeonhole-basic` after every lesson of `combinatorics-enumeration`, after the lesson of 5–6 too: a lesson has no band ([format text](../spec/format/README.md), section 21). A learner of band 1–2 needs the lesson of `pigeonhole-basic` after the lesson of `combinatorics-enumeration` of their own band, not after the lesson of 5–6. Which of a skill's three lessons to give a learner, how to read the order of lessons of different bands and how the lessons of `counting-gaps` stand beside the lessons of `combinatorics-enumeration` is decided by the engine's rule (R58 spec). The card architecture 3.8 spec names how the batch's lessons are built. The engine accepted the proposal of `modules` to T01 of `engine`: a lesson's band is the set of the bands of its tasks, a lesson fits a skill when the band of its recommended task is in this set, and a lesson of another band holds no order (R05 engine; R107 spec; row В14, closed — R09).

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

- **Pictures** — `row` for `gaps-12-d2-1` and `gaps-34-d1-1`: a row of posts with a gap and a label of the interval, 2 and 5 metres. This is a PNG in `media/`, drawn by the code of the content epic from standalone's description — with `alt`, width and height, on a light opaque background (R51 spec). They have no solution pictures.
- **Language** — English, `languages.explanation: en`. Such is the language of standalone's tasks, of the pages `en` and of the format examples. The same module in Russian would be another module, with its own name (section 3.1). The module's language is recorded in product 10 spec and in the addition to О-2 (row В13, closed — R09).
- **What the content epic writes:** the module's files, the converter of tasks, the picture `row`, `givens`, `unknown`, `simpler_case`, the lessons and the rewritten explanations. It also runs the solvers of the batch's tasks once (R05).

### 3.7. What to build

The detailing was finished on 2026-11-02 (T04; R08). Sections 3.1–3.6 and 3.8–3.10 describe what to build, and the epics of [RUN.md](../RUN.md) in what order. The author breaks each epic into tasks when its turn comes: E2 «Первое содержание» was done on 2026-11-07, by tasks T09–T17, and lies in the archive [plan/E2.md](plan/E2.md); E3 «Остальной порт» was done on 2026-11-08, by tasks T18–T43, and lies in the archive [plan/E3.md](plan/E3.md); E4 «Eval-наборы» was done on 2026-11-08, by tasks T44–T46, and lies in the archive [plan/E4.md](plan/E4.md); E5 «Режим eval-набора в CI» was done on 2026-11-09, by tasks T47–T48, and lies in the archive [plan/E5.md](plan/E5.md); E7 «Содержание формата 0.4» was done on 2026-11-10, by tasks T49–T57, and lies in the archive [plan/E7.md](plan/E7.md); E8 «Первые образцы eval-набора» was taken on 2026-11-10 and broken into tasks T58–T63; E9 «Модуль доски» was taken on 2026-11-14 and broken into tasks T68–T71, and E6 waits in table 2.

| What | Section | Epic | What it waits for |
|---|---|---|---|
| the CLI of `module-kit` in the `toolchain` stage of `.devcontainer/Dockerfile`: the binary of a release at an exact tag and sum, the signature by the recipe `release-verify` in CI | 3.2 | E2 | done in T09: `v0.2.3` (R10); since 2026-11-09, `v0.3.0`, format 0.4 (R28) |
| CI: the format's schemas and rules, the tree `modules/` and the package | 3.2 | E2 | done in T09 (R10) |
| CI: the comparison with what is published | 3.2 | E2 | done in T10 (R11); with a real index since the first release, `v0.1.0`, checked in T16 |
| the converter of standalone's tasks and the picture `row` | 3.5 | E2 | done in T11 (R12) |
| the module of olympiad maths — the first batch in format 0.3: the manifest, the graph, the catalogs, 27 tasks, `LICENSE`, `NOTICE` and the picture `row` | 3.5, 3.6 | E2 | done in T12, version 0.1.0 (R13) |
| the lessons of the first batch, 9 | 3.6 | E2 | done in T13 (R14) |
| the set of school maths for grade 3: 4 skills, 12 workbook pages, 4 lessons | 3.8 | E2 | done in T14, version 0.1.0 (R15) |
| the set of short essays: 6 skills, 8 tasks, 6 lessons | 3.9 | E2 | done in T15, version 0.1.0 (R16) |
| `release.yml`, the recipe `release-files` and the first release — the packages and the index | 3.3 | E2 | done in T16 (R17): the shared `release.yml` v0.1.15, line 0.1, a dry run on a pull request. The first release is `v0.1.0`, 2026-11-07, accepted by `release-verify` |
| the rest of the port: 726 tasks of standalone at 42b4408, with 48 tasks solved from the picture, 39 templates with solvers, 16 skills, the pictures of the other kinds | 3.5 | E3: T18–T37 | done in T18–T37: 753 tasks, 19 skills, 19 templates, pictures of 16 kinds (R20–R22) |
| the lessons of the other skills, 32: a lesson per skill and band, from the topic pages of the site, as in R14 | 3.5, 3.6 | E3: T38–T41 | done in T38–T41, with pictures for two skills (the T38, T39 and T41 amendments to R14) |
| the second release — the packages and the index on top of `v0.1.0` | 3.3 | E3: T42 | done in T42: `v0.1.1`, 2026-11-08, accepted by `release-verify`; the only new archive is `olympiad-math` 0.1.1 |
| the gate of `sim` in CI | 3.2 | without an epic of its own — an edit of CI, Д5 in table 4 of [RUN.md](../RUN.md) | a tag of `sim` that translates `time`: v0.1.8 refuses `olympiad-math` 0.3.0 (Д5). What the gate covers is decided: everything that the simulator reads from a package, without bands (R103 spec) |
| the guard of `eval/`: the layout `eval/<id>.json` and a refusal to any sample until the eval-set mode | 3.4 | E4: T44 | done in T44: the recipe `just eval-check` and a step of the job Modules (R25) |
| what to look at a sample with before a push, and where the synthetic samples for a module's tasks come from | 3.4 | E4: T45 | done in T45: chosen are viewing on the samples page of `canvas` and the synthetic samples of the generator of `canvas` on the commit's package, the requests В32 and В33 (R26) |
| the eval-set mode `sample-check --eval-set` in CI instead of the guard; nothing in `eval/` starts with `.`, at any depth | 3.4 | E5: T47 | done in T47: `just eval-check` calls the mode on each module's eval set (the 2026-11-09 amendments to R04 and R25) |
| the first samples — synthetic samples of `canvas` and samples of authors | 3.4 | E8: T58–T63 | done on 2026-11-11, the archive [docs/plan/E8.md](plan/E8.md): the recipe `just synthetic` on the generator of `canvas` v0.3.6 (R29); the check `just eval-match` of each sample against the commit's package (the T59 amendment to R25); 41 synthetic samples of `olympiad-math` and 4 of `grade-3-math` (T60); the procedure for an author's sample and the ban on a session reading people's samples (the T61 amendment to R04); 2 author's samples of `grade-3-math` (T62). The samples of templates, of a choice by hand and of essays went to E6 |
| format 0.4: the word flag of a choice by hand, the kind `time`, the pictures of templates' instances, leading questions, a list of steps in a lesson's explanation and a worked example on a workbook; the descriptions of the solving steps from the format example | 3.5, 3.6, 3.8 | E7: T49–T57 | done in T49–T56: `olympiad-math` 0.3.0 — 44 words, 30 times of day, 311 leading questions, 961 pictures of instances in 8 templates, 3 of them new, lists of steps in 41 lessons; `grade-3-math` 0.3.0 — a worked example on the workbook in 4 lessons. The release `v0.2.1`, 2026-11-10, accepted by `release-verify`; the only new archives are these two modules (the T49–T55 amendments; R28) |
| a parameter of 12 or 24 hours for the kind `time` | 3.5 | without an epic — waits for the format | a minor edit of the format at the content's request (R111 spec). Until it, the 12-hour miss is held by a trap of its own, `twelve-hours-off` (the T50 amendment to R23), and its entry into the format example is В35 in table 3 of [RUN.md](../RUN.md) |
| volunteers' samples and the eval runner of `tutor`, `report-check` on its reports; the samples of templates, of a choice by hand and of essays | 3.4 | E6 | the runner of `tutor` in a release, of stage 2 (R24 spec): there are no binaries `tutor-eval` in the releases of `tutor` yet (architecture 3.2 tutor), — and volunteers. The form, the channel, the periods and the consent withdrawal are there since `spec` v0.15.0 (R30), the viewing of a sample is T87 `canvas`, the words on handing over the file are T20 `canvas`. `report-check` came out in `spec` v0.9.12 (R27) |
| the module of the check board of `tutor`: 4 skills, 11 tasks on a sheet, 1 lesson, 3 synthetic samples | 3.10 | E9: T69 | done in T69, version 0.1.0 (R33); its merge puts out the release `v0.2.7` |
| the author's samples of the module of the check board | 3.4, 3.10 | E9: T70 | the samples page of `canvas` with the tasks of the release that carries the module — row 5.10 of [COORDINATION.md](../spec/COORDINATION.md) |

The first application waited from `modules` only for epic E2 — the batch, its check in CI and the first release (R76 spec; row 3.6 of [COORDINATION.md](../spec/COORDINATION.md)), — and all of it is there since the release `v0.1.0` of 2026-11-07 (R17).

**What the first content asked of the format** — for T38 and epic E6 of `spec` (Д60 in table 4 of the RUN.md of `spec`). Each request went as a row of table 3 of [RUN.md](../RUN.md) «в реестр spec», with options and a recommendation, and `spec` answered all of them: В17 and В18 in v0.9.5 (table 5 of the RUN.md of `spec`; R18), В19–В22 in v0.9.8 (Д113 in table 4 of the RUN.md of `spec`; R19):

| Request | Module | The answer of `spec` |
|---|---|---|
| a rule for a workbook page — the data: 12 pages, 3–7 rows by the layout of `canvas` v0.1.2 — and a worked example on a workbook | `grade-3-math` (R15) | epic E6 of `spec`: the page rule and the worked example with the canvas `workbook` in one version of the format; until it, the workbook is without worked examples (В19) |
| a lesson's skill that is itself a prerequisite of another of its skills does not put the lesson in order | `english-essays` (R16) | [the format text](../spec/format/README.md), section 21, and the example `essays-for-adults.yaml` with `skills: [endings]`, since `spec` v0.9.8: a lesson that must come after the lesson of a skill names only the skills that wait for it (В20) |
| a lesson's explanation has no headings, lists or tables, and a step's number is read as a number | `olympiad-math` (R14) | epic E6 of `spec`: an optional part "list of steps", by a minor edit; until it, the steps of a method are paragraphs (В21) |
| the presentation style by age: one text serves both a child and an adult | `english-essays` (R16), `olympiad-math` (R14) | О-23, supplemented on 2026-11-07: a module asks for no texts of its own by style yet (В22) |
| a place for a reminder not to write names before the work is submitted | `english-essays` (R16) | П-20, supplemented on 2026-11-07 (В18) |
| the explanation languages of the slice's modules | all three (R15) | О-2, supplemented on 2026-11-07 (В17) |

The options and recommendations were in the rows of table 3; the rows left with the moves to `spec` v0.9.5 (R18) and v0.9.8 (R19), and their text remains in the git history of RUN.md.

**What the rest of the port asked of the format** — epic E3, T43. Each request went as a row of table 3 of [RUN.md](../RUN.md) «в реестр spec», and `spec` answered all of them by the decision R108 spec — format 0.4 in epic E6 of `spec` (Д121 in table 4 of the RUN.md of `spec`). The rows left with the move to `spec` v0.9.16 (R27), the rules came with a tag of line 0.10, to which `modules` moved in R28, and the modules will take them by epic E7 in table 2 of [RUN.md](../RUN.md).

| Request | Module | The answer of `spec` |
|---|---|---|
| a choice by hand where a value not among the options gives `not_yet` and a rung of a hint, not `unreadable` | `olympiad-math` (R20) | a flag on a choice by hand, only for options that are words: a word not among the options is "not yet", and a full stop at the end is dropped (В28) |
| a check of time — or time as a choice by hand by В28, and the days of the week, months and dates by tapping | `olympiad-math` (R20) | the check kind `time`, `H:MM`, equality by minutes; dates, days of the week and months by tapping, as they were. Decided otherwise than `modules` recommended (В29) |
| the picture of a template's instance: a path with a placeholder, a file per instance | `olympiad-math` (R22) | placeholders of parameters and derived numbers in the path, and the file of each instance is held by `task.known-pictures` (В30) |
| a task's hint before the first error — the field `hint` | `olympiad-math` (R07, R12) | the closed field `leading_question`, which `api` gives on "Hint" before the first error (В31; the 2026-11-08 amendment to R82 spec) |

The epic's rows that `spec` had already answered — В24–В27: the count of the port of 42b4408 and the trap catalog of the example `olympiad-math` — left with the move to `spec` v0.9.9 (R23).

Not requests, but what format 0.3 corrected in the content, and it is recorded in the decisions: the texts of the whole module and of the workbook's lessons have no digits — the leak check searches for them ([the format text](../spec/format/README.md), section 20; R15); a lesson's example does not solve another lesson's task — review holds this, not `lesson.no-leaks` (R14); the first batch has no hints of standalone — they are a path to the answer (R07). Requests of its own — a check of time and dates and the hints of band 5–6 — the port E3 brought: the requests В29 and В31 in the table above.

**What the eval sets asked of `spec` and `canvas`** — epic E4, T45 (R26). Both went as rows «в реестр spec», and В32 also asked `spec` to decide which tag a synthetic sample on the commit's package names. E4 asked nothing of the format: the guard holds the layout of R04 and waited for the eval-set mode, which came out in `spec` v0.9.13 (the 2026-11-06 amendment to R93 spec; R27). `spec` answered both rows by the recommendations with the decision R110 spec — the tag at the source of the synthetic samples is that of the newest release of `modules` —, and the rows left with the move to v0.9.16 (R27). `canvas` does both by tasks of its own (Д128 in table 4 of the RUN.md of `spec`):

| Request | Who needs it | The row and the answer of `spec` |
|---|---|---|
| synthetic samples of the generator of `canvas` for a module's tasks, on the commit's package, with the source `modules` | E8 | В32; R110 spec, by the recommendation |
| viewing a sample without the network and without AI on the samples page of `canvas`, by the mode "open a file" | E8 for authors' samples, E6 for volunteers' samples | В33; R110 spec, by the recommendation |

### 3.8. The set of school maths for grade 3

The set `novatrail/grade-3-math` is a small module on a workbook, by which `modules` tries the format on a second canvas, before T38 of `spec` (R41 spec). The first application does not wait for it: «другие модули» are not part of it (the COORDINATION.md of `spec`). Why so is R15.

- **The manifest, the graph and the catalogs** are those of the format examples `school-math-grade-3` word for word, but for the ID, `novatrail/grade-3-math` (R03). The audience is `under_13`, grade 3, one band `grade-3`, the explanation language is `en`, without `dependencies`.
- **Skills** — the example's four skills, closed under prerequisites:

  | Skill | Prerequisite | Tier | Teaches the anchor skill |
  |---|---|---|---|
  | `times-tables` | — | 1 | `core.multiplication-table` |
  | `division-with-remainder` | `times-tables` | 2 | `core.division-with-remainder` |
  | `fractions-of-a-whole` | `times-tables` | 2 | `core.fractions-of-a-whole` |
  | `comparing-fractions` | `fractions-of-a-whole` | 3 | — |

  The graph is not a chain, so the module is not sequential ([the format text](../spec/format/README.md), section 21).
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
  - Three pages per skill are enough for the corridor of `engine` too: the measurement of 2026-11-11 on `sim` v0.1.8 — up to 32 pages, two for each difficulty from 1 to 5 — did not raise the share of the items in the corridor, and the set of pages stays (Д8 in table 4, В39 in table 3 of [RUN.md](../RUN.md)).
  - There are no pictures, `predicate` or templates.
- **Lessons** — a lesson per skill, `lessons/<skill>.yaml`: an explanation, a worked example on the workbook, then three pages by increasing difficulty. The texts of the lessons, traps and techniques have no digits: the pages' keys are small numbers, and the leak check searches for them in all these texts ([the format text](../spec/format/README.md), section 20). There are no signs of comparison there either, though the check does not search for them: the page's statement names all three signs, and a key that stands in a text seen together with the item is not searched for (ibid.; R15). After the explanation comes a worked example on the workbook: two or three lines with boxes that are not among the lines of the lesson's pages, and a demo that fills each box; the statement and the texts of the demo are without digits ([the format text](../spec/format/README.md), section 21; the T55 amendment to R15).
- **The page fits.** Each of the 12 pages went through the layout of `canvas` v0.1.2 — its `pageOf` — and takes from 3 to 7 rows. In format 0.3 there was no page rule, only a guideline — up to 10 rows; since 0.4 the rule `task.workbook-page` holds it by the format's layout `page.yaml`, and all 12 pages pass it (`module-kit` v0.3.0, 2026-11-09; [the format text](../spec/format/README.md), section 14; R28). The source is sources 4.19.
- **`LICENSE`** — the text of Apache-2.0. There is no `NOTICE`: the texts are its own.

### 3.9. The set of short essays

The set `novatrail/english-essays` is a small module of writing on a handwriting sheet, by which `modules` tries the format on a second subject and on a path without a key — feedback by a rubric —, before T38 of `spec` (R41 spec). The first application does not wait for it: «другие модули» are not part of it (the COORDINATION.md of `spec`). Why so is R16.

- **The manifest, the graph, the techniques and the solving steps** are those of the format examples `short-essays` word for word, but for the ID, `novatrail/english-essays` (R03); so are the traps, and one of its own, below. The audience is all four age categories, without grades; the bands are `first-steps` and `whole-essays`, also without grades; the explanation language is `en`, without `dependencies`. There is an item with a verdict, so the manifest has the example's `rung_templates`.
- **Skills** — the example's six skills:

  | Skill | Prerequisite | Tier | How it is assessed |
  |---|---|---|---|
  | `reading-model-essays` | — | 1 | `none`: only a supporting skill and a lesson of explanations alone |
  | `choosing-a-first-sentence` | — | 1 | `verdict` |
  | `writing-a-first-sentence` | `choosing-a-first-sentence` | 2 | `rubric` |
  | `paragraphs` | `writing-a-first-sentence` | 3 | `rubric` |
  | `endings` | `paragraphs` | 4 | `rubric` |
  | `improving-an-essay` | `paragraphs` | 4 | `rubric` |

  The graph is not a chain: tier 1 has two skills, and `endings` and `improving-an-essay` both wait for `paragraphs`. So the module is not sequential, and where the graph does not order the lessons, the engine's rule chooses — by the lesson's ID (R05 engine; architecture 3.8 spec) —, below, "Lessons".
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
  - An essay has `binding: sheet` and a rubric of three criteria on structure and content, without capital letters, spelling and punctuation (R71 spec).
  - The statement of each essay ends with the reminder "Do not write your name or the real name of anyone you know." (П-20). Whether the platform shows it itself is decided together with П-20, before the beta with children, and until the answer the reminder stands in the statements (the 2026-11-07 supplement to П-20; R18).
  - The topics for those under 18 are made up and general, without family, friends, home and school; a topic about the learner's life is only in the task for adults (П-20; privacy 3.2 spec).
  - There are no pictures.
- **Lessons** — six, `lessons/<id>.yaml`:
  - `reading-model-essays` — explanations alone, from the example;
  - `first-sentences` — the example's explanations, the three choice tasks and `talking-umbrella-1`;
  - `paragraphs` and `endings` — an explanation and an essay;
  - `improving-an-essay` — an explanation, a worked example on a lined sheet that corrects its own draft, and `dragon-draft-1`;
  - `writing-about-a-day` — the format example `essays-for-adults.yaml` word for word, `adult` only, with the skill `endings`. The module narrowed the lesson's skills to `endings` before the example did (R16): with `paragraphs`, the lesson would not wait for the lesson `paragraphs`. Since `spec` v0.9.8 the example has `skills: [endings]` too ([the format text](../spec/format/README.md), section 21; R19).

  The order of the lessons by the graph: `first-sentences`, then `paragraphs`, then `endings`, `improving-an-essay` and `writing-about-a-day`, which the graph does not order among themselves. The graph links `reading-model-essays` to nothing. Where the graph does not order, the engine's rule chooses: lessons whose skills are not all `verdict` go by the graph, and where the graph gives no order, by the lesson's ID (R05 engine; architecture 3.8 spec). If the rule takes the ID among the lessons that the graph already lets in, the order is `first-sentences`, `paragraphs`, `endings`, `improving-an-essay`, `reading-model-essays`, `writing-about-a-day`; how exactly it combines the graph and the ID, R05 engine holds. Each lesson has a skill that is not `verdict`, so the engine neither adds segments to them nor skips any: they go in the author's order ([the format text](../spec/format/README.md), section 21). The choice tasks go in a row in the lesson, and outside a lesson the engine chooses among the three. The texts of the lessons have no `leak_phrases` of the choice tasks.
- **`LICENSE`** — the text of Apache-2.0. There is no `NOTICE`: the texts are its own and those of the format examples of `spec`, under the same Apache-2.0.

### 3.10. The module of the check board

The module `novatrail/check-board` is a small module of maths in Russian for adults: the eleven tasks of the check board of `tutor` — `just board`, `internal/board/tasks.go` of `tutor` (R19 tutor) —, so that the author solves them in the app again and again, on the path a learner takes: the reading by `tutor`, the verdict of `module-kit`, the hint ladder and the tutor's words. Its eval set checks every new tag of `tutor` once the eval runner of epic E6 is in CI (section 3.4). It stands beside the three modules of the first slice, not among them, and the first slice waits for none of it (the 2026-11-14 amendment to R76 spec; rows 3.10, 5.10 and 6.17 of [COORDINATION.md](../spec/COORDINATION.md)). Why so is R33.

- **The source** — the banks `eval/bank.json` and `eval/free-bank.json` of `tutor` at `v0.1.35`, its latest tag on 2026-11-14; the board, the banks and `internal/task` are as at `v0.1.34`. Each task is the task of the format that `internal/task.Sheet` of `tutor` builds of its bank's entry (R68 tutor), so its open view — the statement, the canvas and the items — is the JSON of the task in `eval/open-tasks.json` of `tutor`, the same for all eleven. The closed fields that `Sheet` fills with stand-ins — the skill, the band `measurement` and the difficulty 1 — are the module's own, below.
- **The manifest**: the ID `novatrail/check-board`, the title «Задачи на листе», the subject `mathematics`, the audience `adult` alone, without grades, one band `board` «Задачи доски», the explanation language `ru`, without `dependencies`. The rung templates are those of the other three modules, in Russian and without digits: «Какой шаг твоей работы дал этот ответ?» and «Прочитай условие ещё раз, медленно. Есть ли в нём то, чего твоя работа ещё не использует?».
- **Skills** — four, without prerequisites and without anchor skills:

  | Skill | Title | How it is assessed | Tasks |
  |---|---|---|---|
  | `linear-equation` | «Линейное уравнение» | `verdict` | `equation-1`, `equation-2` |
  | `calculation` | «Цепочка вычислений» | `verdict` | `chain-1`, `chain-2` |
  | `free-solution` | «Свободное решение» | `verdict` | `posts` and the six of `eval/free-bank.json` |
  | `writing-on-a-sheet` | «Решение на листе» | `none` | a supporting skill of every task |

  The first three are the kinds of the banks, by which `Sheet` names a task's skill. The fourth is practised and not assessed: it keeps the lesson in the board's order, below. All four are in tier 1, so the module is not sequential.
- **Traps** — 22, those that the eleven tasks name, each once: the banks keep one ID and one description for a mistake. The description is the bank's, in impersonal words: the bank's names a gender, «перенёс», and the learner's progress shows a description ([format text](../spec/format/README.md), section 16). AI reads it to suggest the trap from the work. The advice is the module's own, for the learner: addressed with «ты», with no verb in the past tense singular, which would name a gender, and without digits. There are no techniques and no solving steps: `Sheet` gives a task neither, and the format lets a module do without them.
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
  - The leak phrases are the bank's, without its English ones: the measurement worded hints in both languages (R05 tutor), and the module's texts are in Russian ([format text](../spec/format/README.md), section 10). After them stands the key in words, «семь», which the leak check does not read as a number ([format text](../spec/format/README.md), section 20). In `posts` the key, 5, stands in the statement, so the leak check does not search for it, and «5 столбов» stands among the leak phrases in its place.
  - The difficulty is the number of actions of the bank's solution: its lines with a calculation, without the restated equation, the drawing and «Ответ:».
  - There are no `givens`, `unknown`, `leading_question`, `simpler_case`, pictures or templates: a hint asked for before the first wrong answer has nothing to give, and the third rung gives the fallback of the second ([format text](../spec/format/README.md), section 19).
- **The lesson** — one, `lessons/board.yaml`, «Задачи доски по порядку»: the eleven tasks in the board's order, without an explanation. Its skills are all four, and since `writing-on-a-sheet` is not `verdict`, the engine adds no segment to the lesson and skips none ([format text](../spec/format/README.md), section 21). The engine of the first slice does neither in any lesson (architecture 3.8 spec), but epic E7 of `engine` would give a task of consolidation at the end of a segment of a lesson whose skills are all `verdict` (row 4.14 of [COORDINATION.md](../spec/COORDINATION.md)), and the board is to be the same path every time.
- **The eval set** — 3 synthetic samples, one per pair of a primary skill and the band, of the task first by ID, as T60 decided (section 3.4): `chain-1`, `album` and `equation-1`, seeds 1–3, 1.4 MB. They name the tag `v0.2.6`, the newest release when they were drawn (R110 spec), whose index does not hold the module yet: of a module's eval set the runner takes the task from the package of the commit, and the tag says only what the sample was made on ([the text of the contracts](../spec/contracts/README.md), section 13). The author's samples are T70 of [RUN.md](../RUN.md), once the samples page of `canvas` carries the release with the module (row 5.10).
- **`LICENSE`** — the English text of Apache-2.0, like the other modules'. There is no `NOTICE`: the tasks are those of `tutor`, under the same Apache-2.0.

## 4. Sources

Since T64 the records of sources are in [sources.md](sources.md), by topic, and the claims of sections 1–3 cite their record: `sources 5.1`. A pointer remains here: the decision log and the archives of `docs/plan/` cite this section, and they are not rewritten (R31). An item that stood here is the record with the same title: The files of a GitHub release — sources 2.1; Immutable releases — sources 2.2; The filters `paths` and required checks — sources 2.3; The age of a GitHub account — sources 2.4; What remains after a push — sources 2.5; `-text` in `.gitattributes` — sources 3.1; DCO 1.1 — sources 3.2; PNG chunks — sources 3.3; The font and drawing — sources 3.4; The shared `release.yml` of `.github` — sources 4.1; The template of `.github` on 2026-11-06 — sources 4.2; The shared `release.yml` and the template on 2026-11-07 — sources 4.3; The `release.yml` of `canvas` — sources 4.4; The shared workflows of `.github` v0.1.15 — sources 4.5; The DCO check — sources 4.6; The CONTRIBUTING of `.github` — sources 4.7; The release of `module-kit` v0.2.3 — sources 4.8; The release of `module-kit` v0.3.0 — sources 4.9; The command `index` — sources 4.10; The command `diff` — sources 4.11; `module-kit` on the previous version — sources 4.12; The module's version and the eval set — sources 4.13; A reproducible archive — sources 4.14; Where the runner takes the samples from — sources 4.15; The request of `api` on the layout of the release — sources 4.16; `sim` on the gate of the graph — sources 4.17; The sample file and the samples page — sources 4.18; The layout of the workbook — sources 4.19; The license of standalone — sources 5.1; standalone at 89986bc — sources 5.2; standalone's tasks by topic — sources 5.3; standalone's test bench — sources 5.4; The first batch from standalone — sources 5.5; standalone at 42b4408 — sources 5.6; standalone's prototype — sources 5.7.
