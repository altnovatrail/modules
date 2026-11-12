# modules — sources

> The external facts that the documents of `docs/` rest on: the terms of GitHub, standards, the documentation of other programs — NovaTrail's repositories — and the texts of MathTrail. Here lie their records — the address, the version or the date of the page, the date read and the quote in the language of the source, — and [architecture.md](architecture.md) holds the claim and a reference to the record by its number, as to a section: `(sources 2.1)` (RULES.md, The bar; R143 spec).
>
> The records are grouped by topic. Until T64 they stood as items in architecture 4. A pointer remains there: the decision log and the archives of `docs/plan/` cite that section, and those are not rewritten (R31). The results of measurements are in [measurements.md](measurements.md). The file is in English since T65 translated it.

## 1. How to keep it

- **A record** is a subsection "N.M" of a topic: one source, or one reading of several sources for the sake of one decision. It holds the address, the version or the date of the page, the date read, and the quote, in the language of the source. The first line of a record names the sections and decisions it was read for. What was not read from the primary source — a copy, a pointer, a page that did not open — the record names as such.
- **A reference.** A claim of a document refers to a record by its number, as to a section of any document: `(sources 4.8)`. The number of a record does not change: a new record goes at the end of its topic, a new topic at the end of the file.
- **The source was read again** — an item with its own date is added to the record, and the earlier text is not edited: the decisions of its day rest on it. If the source changed and the claim diverged from it, that is an amendment to a decision or an open question, not an edit of the record.
- **What is not here:** the quotes that an entry of the decision log gives itself — they stay in the entry (RULES.md, How we work: a decision is never rewritten), — and the quotes that [architecture.md](architecture.md), sections 1–3 give as the argument of their item: they stay where they stand, as in T133 of `spec`.

## 2. GitHub

### 2.1. The files of a GitHub release

For architecture 3.3.

[About releases](https://docs.github.com/en/repositories/releasing-projects-on-github/about-releases), the section "Storage and bandwidth quotas", read 2026-11-02: "Up to 1000 release assets may be associated with a single release. Each file included in a release must be under 2 GiB." A release file has only a name, without folders ([Upload a release asset](https://docs.github.com/en/rest/releases/assets#upload-a-release-asset), ibid.: "GitHub renames asset filenames that have special characters, non-alphanumeric characters, and leading or trailing periods").

### 2.2. Immutable releases

For architecture 3.3.

[Managing releases in a repository](https://docs.github.com/en/repositories/releasing-projects-on-github/managing-releases-in-a-repository), read 2026-11-02: "If you have enabled immutable releases for your repository, you cannot add, replace, or delete assets after a release is published, and you cannot move or delete its tag while the release exists."

### 2.3. The filters `paths` and required checks

For architecture 3.2.

[Workflow syntax for GitHub Actions](https://docs.github.com/en/actions/reference/workflows-and-actions/workflow-syntax), `on.<push|pull_request>.paths`, read 2026-11-02: "If a workflow is skipped due to path filtering, branch filtering, or a commit message, then checks associated with that workflow will remain in a "Pending" state. A pull request that requires those checks to be successful will be blocked from merging."

### 2.4. The age of a GitHub account

For architecture 3.4; R04.

[GitHub Terms of Service](https://docs.github.com/en/site-policy/github-terms/github-terms-of-service), B. Account Terms, 3. Account Requirements, in effect since 27 April 2026, read 2026-11-02: "You must be age 13 or older."

### 2.5. What remains after a push

For architecture 3.4; R04.

[Removing sensitive data from a repository](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/removing-sensitive-data-from-a-repository), the section "About sensitive data exposure", read 2026-11-02: "If you only rewrite your history and force push it, the commits with sensitive data may still be accessible elsewhere: In any clones or forks of your repository; Directly via their SHA-1 hashes in cached views on GitHub; Through any pull requests that reference them". And ibid.: "you can permanently remove cached views and references to the sensitive data in pull requests on GitHub by contacting us through the GitHub Support portal".

## 3. Standards and libraries

### 3.1. `-text` in `.gitattributes`

For architecture 3.1 and 3.3.

[gitattributes](https://git-scm.com/docs/gitattributes), the attribute `text`, read 2026-11-02: "Unsetting the `text` attribute on a path tells Git not to attempt any end-of-line conversion upon checkin or checkout." And `text` turns on the conversion of line endings — "enables end-of-line conversion on checkin and checkout" — for any file it is set on, media files too.

### 3.2. DCO 1.1

For architecture 3.4; R04.

[Developer Certificate of Origin 1.1](https://developercertificate.org/), read 2026-11-02:
- (a) "The contribution was created in whole or in part by me and I have the right to submit it under the open source license indicated in the file";
- (b) "The contribution is based upon previous work that, to the best of my knowledge, is covered under an appropriate open source license and I have the right under that license to submit that work with modifications…";
- (c) "The contribution was provided directly to me by some other person who certified (a), (b) or (c) and I have not modified it";
- (d) "a record of the contribution (including all personal information I submit with it, including my sign-off) is maintained indefinitely".

### 3.3. PNG chunks

For architecture 3.4.

[Portable Network Graphics (PNG) Specification (Third Edition)](https://www.w3.org/TR/png-3/), W3C Recommendation 24 June 2025, read 2026-11-02:
- "A critical chunk is a chunk that is absolutely required in order to successfully decode a PNG image from a PNG datastream" (there, 11.2). There are four critical chunks — IHDR, PLTE, IDAT and IEND —, and the others are ancillary (there, 4.8.2);
- not every image needs PLTE: "This chunk shall appear for color type 3, and may appear for color types 2 and 6; it shall not appear for color types 0 and 4" (there, 11.2.2);
- among the ancillary ones are tEXt, iTXt and zTXt, "Textual information (which may be compressed) associated with the image", eXIf, "Exchangeable image file format metadata", and tIME, "The time when the PNG image was last modified" (there, 4.7, table 2).

### 3.4. The font and drawing

For architecture 3.5; R12.

`golang.org/x/image` v0.47.0: `font/gofont/goregular`, `font/opentype` and `vector`, the license BSD-3-Clause (the module's `LICENSE`), read 2026-11-06; `image/png` of Go's standard library writes an opaque image without ancillary chunks — the test `TestRowDraws` holds that there are none.

## 4. NovaTrail's repositories

### 4.1. The shared `release.yml` of `.github`

For architecture 3.3.

The release v0.1.15, commit 89d3cee, the latest on 2026-11-02, read 2026-11-02:
- `.github/workflows/release.yml`: `paths` is read as a line per pathspec and goes into `git diff --quiet "$previous" "$TARGET" -- "${paths[@]}"`. The caller's recipe runs as `just release-files` in a container of the `toolchain` stage of its `.devcontainer/Dockerfile`, on the tree from `git archive` and without tokens. In `dist/` there are only files with names matching `^[A-Za-z0-9_+][A-Za-z0-9._+-]*$`, and each file is signed by a bundle `<file>.sigstore.json` beside it;
- architecture 3.3 .github: the input `previous-file: module-index.json` for `modules` is «с `v0.1.14`», and the recipe then takes «два параметра, `release-files version previous:`». The first release gets an empty string as the second. The lists of `paths` in the template are of the kind go and of the kind typescript; there is no kind for `modules`.

### 4.2. The template of `.github` on 2026-11-06

For architecture 3.3; the 2026-11-06 amendment to R03.

`docs/architecture.md` of `.github`, `main`, commit af40f94, the latest release v0.1.15, read 2026-11-06 (`gh api`). In architecture 3.3 .github, «Вызывающие файлы», there is no list of `paths` of the kind `modules`: neither `withdrawals.yaml` nor `modules/**` occurs in the file. Of the release of `modules` the section says only the input `previous-file: module-index.json`, «с `v0.1.14`».

### 4.3. The shared `release.yml` and the template on 2026-11-07

For architecture 3.3; R17.

`.github`, the release v0.1.15, commit 89d3cee, and `main`, commit af40f94, read 2026-11-07 (`gh api`). There is no release newer than v0.1.15, and `main` is three commits of documents ahead: `.github/workflows/release.yml` is one and the same file on both. There is still no list of the kind `modules` in the template of architecture 3.3 .github, only the kind go. The inputs are `paths`, `publish` ("false builds the files, signs them and checks the signatures, and publishes nothing"), `npm-package` and `previous-file`. The build: `git archive "$TARGET"`, then `docker build --target toolchain -` from `.devcontainer/Dockerfile` without a context, then `docker run … just release-files "${arguments[@]}"`, where the second argument is `""` or `/previous/module-index.json`. Neither a token nor the step's variables go into the container, only `HOME=/tmp`, and `spec/` in the tree is empty. The recipe needs no network, but it is not closed off either: `docker run` has no `--network none`. A pull request from a fork skips the dry run, because cosign's keyless signing needs `id-token: write`.

### 4.4. The `release.yml` of `canvas`

For architecture 3.3; R17.

`.github/workflows/release.yml` on `main` of `canvas`, read 2026-11-07 (`gh api`): `pull_request:` without a list of paths of its own, the condition "`github.event_name != 'workflow_run' || github.event.workflow_run.conclusion == 'success'`" and `publish: ${{ github.event_name != 'pull_request' }}`. The dry runs on pull requests on 2026-11-02 — from the author's branch and from Dependabot — are green, but for one (the API `actions/workflows/release.yml/runs?event=pull_request`).

### 4.5. The shared workflows of `.github` v0.1.15

For architecture 3.2; R10.

`.github/workflows/`, the release v0.1.15, read 2026-11-06 (`gh api`): `documents.yml` runs `docs-check-repo`, and `go.yml` and `typescript.yml` run the caller's recipes `ci-lint`, `ci-test`, `ci-vuln` and `ci-licenses`. There is no workflow among them that runs recipes of the caller's choice. The shared `codeql.yml` takes the input `languages` with the value `javascript-typescript` both at v0.1.0, ff2a359, and at v0.1.10, f3addb7, — to which Dependabot raised the `codeql.yml` of `modules` on 2026-11-06.

### 4.6. The DCO check

For architecture 3.4.

`.github/workflows/dco.yml` of `.github`, `main`, commit af40f94, read 2026-11-02: "The sign-off has to name the author exactly, by name and address". The code checks the trailer `Signed-off-by` of each commit, but for merge commits, against the commit's author, `%an <%ae>`, and looks at who opened the pull request only for the sake of Dependabot. CI of `modules` calls this workflow at v0.1.15, commit 89d3cee, — Dependabot raised it from v0.1.0 (#8), — and the rule there is the same: the author's sign-off `%an <%ae>`, the exception for Dependabot and the skipping of merge commits (`dco.yml` at 89d3cee, read 2026-11-07, `gh api`).

### 4.7. The CONTRIBUTING of `.github`

For architecture 3.4.

`CONTRIBUTING.md`, `main` of `.github`, commit af40f94, read 2026-11-02: "No samples for eval sets yet. Handwriting, voice and photo samples for eval sets are not accepted until the project decides how the people who make them consent to their publication." `modules` has no CONTRIBUTING of its own.

### 4.8. The release of `module-kit` v0.2.3

For architecture 3.2; R10.

[altnovatrail/module-kit, release v0.2.3](https://github.com/altnovatrail/module-kit/releases/tag/v0.2.3), commit e069943, the latest on 2026-11-06, immutable (`gh release view`, `isImmutable: true`), read 2026-11-06:
- the linux binary lies in the release as it is, without an archive: `module-kit_v0.2.3_linux_amd64` and `module-kit_v0.2.3_linux_arm64`. The sums of all the files are in `SHA256SUMS`, and each file has its `<file>.sigstore.json`;
- `just -f spec/justfile release-verify module-kit v0.2.3 SHA256SUMS`: "SHA256SUMS: module-kit v0.2.3, signed by the shared release.yml of .github v0.1.15";
- `cmd/module-kit/main.go`: `--version` prints `module-kit %s` and, as the second line, `Unicode %s`;
- architecture 3.6 module-kit, «Команды `check` и `build`»: «вызов — `module-kit check <module>` и `module-kit build <module> <folder>`» (the source writes these placeholders in Russian), and for `build` — «папка не существует или пуста, иначе код 2 и ничего не тронуто».

### 4.9. The release of `module-kit` v0.3.0

For architecture 3.2; R28.

[altnovatrail/module-kit, release v0.3.0](https://github.com/altnovatrail/module-kit/releases/tag/v0.3.0), the first of line 0.3, immutable (`gh release view`, `isImmutable: true`), read 2026-11-09:
- R52 module-kit at the tag v0.3.0: «Библиотека читает формат модуля 0.4, и только его: модуль 0.3 отклоняется по `manifest.format-version`»;
- `just module-kit-verify`: "SHA256SUMS: module-kit v0.3.0, signed by the shared release.yml of .github v0.1.15"; the sums of `module-kit_v0.3.0_linux_amd64` and `…_linux_arm64` are in `ARG` of the Dockerfile;
- the forms of the call of `check`, `build`, `index` and `diff` in `docs/architecture.md` of `module-kit` are the same as at v0.2.3 (`gh api` on both tags).

### 4.10. The command `index`

For architecture 3.3; R17.

architecture 3.6 module-kit, the release v0.2.3, commit e069943, read 2026-11-07 (`gh api`; `main` matches the tag): «module-kit index (--previous <module-index.json> | --first) --tag vX.Y.Z (--withdrawals <withdrawals.yaml> | --no-withdrawals) --out <файл> <папка пакета>...». The inputs are «папки пакетов, которые написал `build`: каждая — `modulepackage.Read` по дайджесту байтов своего `inventory.json`, с `device/`». The output is «`--out`, которого ещё нет». The codes: 0 — the index is written, 1 — a refusal by the rules of the index, 2 — the command could not do its job. `module-kit` does not write the package's archive: «Архив для релиза кладёт рецепт `modules` (R03 modules)».

### 4.11. The command `diff`

For architecture 3.2; R11.

architecture 3.5 module-kit, the release v0.2.3, read 2026-11-06: «module-kit diff [--json] (--index <module-index.json> | --first) [--published <папка пакета>]... <модуль>»; `<module>`, «<модуль>» in the quote, is «папка модуля или папка пакета»; `--published` is «папки пакетов файлами …, без `device/`», and «пакет каждой версии, которую называет индекс, нужен: без него сравнить нельзя»; code 2 is «нет пакета версии, которую называет индекс». Tried on v0.2.3 the same day: the same version with the same digest — code 0, "the newest published version as it was"; with another one — code 1, `module-version.above`; a new version — 0; without the package — 2; `--first` — 0.

### 4.12. `module-kit` on the previous version

For architecture 3.3.

architecture 2.4 module-kit, `main` of `module-kit`, the tag v0.1.1, commit 3f5a4c4, read 2026-11-02: «Откуда CI `modules` берёт прошлую версию модуля для сравнения, решает T01 modules». The CLI at this commit knows only `--version` (`cmd/module-kit/main.go`).

### 4.13. The module's version and the eval set

For architecture 3.4.

T04 in the RUN.md of `module-kit`, «Правила semver для `diff`», `main` of `module-kit`, commit 2489d19, read 2026-11-02: «Перечень изменений по составу модуля — … eval-набор (concepts 3 spec) — и какая часть версии растёт от каждого». The task is not started.

### 4.14. A reproducible archive

For architecture 3.3; R17.

The recipe `release-files` of `tutor`, `main`, read 2026-11-07 (`gh api`): "tar --create --format=gnu --sort=name --owner=0 --group=0 --numeric-owner --mode=go-w --mtime=… | gzip -n -9". The `toolchain` stage of `modules` has GNU tar 1.35, and it has `--sort=name` (`tar --version` and `tar --help`, 2026-11-07).

### 4.15. Where the runner takes the samples from

For architecture 3.4.

R49 tutor, `main` of `tutor`, commit be607e5, read 2026-11-02: «eval-раннер берёт задачи из пакета модуля, а образцы — из папки модуля, ведь eval-набор в пакет не входит».

### 4.16. The request of `api` on the layout of the release

For architecture 3.3.

В6 in table 3 of the RUN.md of `api`, «в реестр spec», for T01 of `modules`; `main` of `api`, commit d409351, read 2026-11-02: a package's files are found «по записи — ID модуля, версии и тегу, — без перебора файлов релиза: по публичному адресу, без токена и REST API». The recommendation of `api` is «один архив на пакет, названный по дайджесту пакета», and to fix a withdrawn major version by a new version of the same major in the same release that withdraws it.

### 4.17. `sim` on the gate of the graph

For architecture 3.3.

T03 in the RUN.md of `sim`, `main` of `sim`, commit 067f624, read 2026-11-02: «CI `modules`: как он узнаёт, что граф изменился, и на каком пакете идёт стенд».

### 4.18. The sample file and the samples page

For architecture 3.4.

`canvas`, `main`, commit 5cd33ce, read 2026-11-02:
- R09 canvas: an eval-set sample has no writer code, and of the person the file holds «ни имени, почты, аккаунта GitHub, ни устройства, браузера, экрана, языка, ни даты, возраста и класса»; the container is «один JSON на образец», the crops are «PNG в base64»;
- R11 canvas: «eval-набор закрыт в каждой сборке, пока нет утверждённого текста по П-16»; «Версию набора текстов и «кто писал» сверяют организатор замера, раннер `tutor` и CI `modules`»; «кто открывает pull request, решает T02 modules»;
- architecture 3.5 canvas: «У образца eval-набора имя на образце ищет T02 modules»;
- architecture 3.6 canvas: on the page, the declarations «мне 18 лет или больше» and «форма согласия версии … подписана и принята организатором», and «куда передать файл: … образец eval-набора — pull request по правилам T02 modules».

### 4.19. The layout of the workbook

For architecture 3.8; R15.

[altnovatrail/canvas, release v0.1.2](https://github.com/altnovatrail/canvas/releases/tag/v0.1.2), commit 8e9c6f6 of 2026-11-02, the newest on 2026-11-07, read 2026-11-07, a clone in the session's temporary folder:
- `src/core/page.ts`: "A task whose lines do not fit is a TypeError: nothing shrinks", and `pageOf` is "a function of the items alone";
- `vectors/pages.json`: the font size 30, a box 100 × 100, margins 40, and a short item is "no wider than the column (440) less the indent (56)". A long item fills lines of 864, and a row below 1400 − 40 is a refusal `task.items[i]`;
- `pageOf` in the temporary folder with Node 24.21: all 19 cases of `vectors/pages.json` matched. The 12 pages of the set have no refusal, and the longest ends at 932 of 1360.

## 5. MathTrail: standalone and its prototype

### 5.1. The license of standalone

For architecture 3.5.

[mathtrail/mathtrail-standalone, `LICENSE`](https://github.com/mathtrail/mathtrail-standalone/blob/c6cbeb2/LICENSE), the branch `main`, commit c6cbeb2, read 2026-10-19: "MIT License", "Copyright (c) 2026 MathTrail". There are no other files named LICENSE, COPYING or NOTICE in the repository; GitHub detects its license as MIT (`gh repo view`, `licenseInfo`, 2026-10-19).

### 5.2. standalone at 89986bc

For architecture 3.5; R05.

[mathtrail/mathtrail-standalone](https://github.com/mathtrail/mathtrail-standalone/tree/89986bc), `main`, commit 89986bc of 2026-11-01. Cloned into the session's temporary folder, read 2026-11-02:
- `LICENSE`: "MIT License", "Copyright (c) 2026 MathTrail"; `CONTRIBUTING.md`: "The code and the content are MIT, and a contribution arrives under the same licence";
- the count of the parts and answers in architecture 3.5 is by jq over `content/examples/*.json`, `content/catalogs/` and `content/solvers/`. The kind of an answer is by the text of the correct option: an integer, `a/b`, a number with `%`, a time `h:mm`, "Yes" or "No", other text with digits, text without digits. An explanation "names the answer" if the key's number stands in it and the statement does not have it, — roughly by the rule of the number search of [the format text](../spec/format/README.md), section 20;
- 705 tasks: `SPEC.md` — "Grades 1–4 — 450 tasks, 442 of them the prototype's" and "Grades 5–6 — 153 new tasks", and R309 standalone added another 102. On the tasks from the prototype: "They are ported as they are (T23)… Their questions stay word for word";
- own wordings: `SPEC.md` — "Nobody else's wording, anywhere. The ideas come from open sources; the text is ours". `CONTRIBUTING.md` — "The wording is your own. Ideas from open sources are fine — a wording lifted from a book or an olympiad is not, whatever its age". `site/content/en/research.yaml` — "from a book only a task's idea and its answer are taken, and the wording, the story, the options and the traps are written anew";
- the ideas of grades 5–6, R309 standalone (`docs/decisions.md`): "The ideas come from MathNet (mathnet.mit.edu, under CC BY 4.0, the countries that assert rights to their problems keeping them) at `5-6`, in the topics it has a section for, and from the sources of the prototype's research 12 elsewhere. An idea is taken and never a wording";
- the catalogs in the format examples are the same since then: `git diff c6cbeb2 89986bc` is empty for `content/catalogs/topics.json` and `content/catalogs/traps.json`, and `git diff abf3b61 89986bc` for `site/content/en/techniques.yaml`;
- the authors: `gh api repos/mathtrail/mathtrail-standalone/contributors`, 2026-11-02 — RyazanovAlexander, 153 commits, and dependabot[bot], 12.

### 5.3. standalone's tasks by topic

For architecture 3.3.

[mathtrail/mathtrail-standalone](https://github.com/mathtrail/mathtrail-standalone/tree/89986bc/content/examples), `content/examples/`, commit 89986bc, read 2026-11-02 (`gh api`): 17 files `<topic>.json`, one per topic, with all its tasks, and the solvers are `content/examples/solvers/<task id>.star`. There are no folders by grade: the band is in the task's ID and its fields.

### 5.4. standalone's test bench

For architecture 3.5; R12.

[mathtrail/mathtrail-standalone](https://github.com/mathtrail/mathtrail-standalone/tree/89986bc01a12ff3188cf08194677bfae0464568a), commit 89986bc, read 2026-11-06 (`gh api`; a clone in a temporary folder):
- `content/solver_test.go`, `TestEveryReferenceTaskProvesItsOwnAnswer`: "It runs each solver the way a submitted one is run — the same sandbox, the same dialect and helpers, the same two runs under rotated labels", a subtest per task ID, and `prove` — "the solver proves %s, the task claims %s";
- `internal/domain/solver/verdict.go`: "The second run happens only when the first reached exactly one option";
- `go.mod`: `go 1.27.3` and `go.starlark.net v0.0.0-20260908191801-89a6a09411d5`;
- `content/schemas/task.json`, `$defs/row` and `row_item`: "Objects in a row … 2 to 12 items", `line`, `gaps`, `span`, `copies` 1–2, an object has `label`, `below` and `mark` — `dot`, `ring` or `square`, and a skip "is never first or last, and is never next to another skip";
- the run of 2026-11-06: the test bench proved the answers of all 27 tasks of the first batch, and the tasks and the converter's two pictures `row` passed `module-kit check` v0.2.3 in a module from the format examples without warnings, but for two explanations that name the key and that T12 rewrites (architecture 3.6).

### 5.5. The first batch from standalone

For architecture 3.6; R07.

[mathtrail/mathtrail-standalone](https://github.com/mathtrail/mathtrail-standalone/tree/89986bc/site/content/en/topics), commit 89986bc. The repository was cloned into the session's temporary folder and read 2026-11-02:
- `site/content/en/topics/gaps-and-boundaries.yaml`, `enumeration.yaml` and `pigeonhole-principle.yaml` — `hero.lead`, the cards `basis`, `method.steps` and three examples `examples` each, with `question`, `steps`, `answer` and `note`. At the start of each page: "an example here is worked through to that answer";
- `site/data.json`, the member `examples` — the level and the solver of each example. The three topics have an example each for `1-2`, `3-4` and `5-6`;
- `content/examples/counting.gaps.json`, `combinatorics.enumeration.json` and `pigeonhole.basic.json` — the tasks, their explanations, pictures and hints;
- the selection and the count are by a Node script with js-yaml 4.1.0 in the same folder. The kind of an answer is found by the text of the correct option, the question "Which of these…" by the words, a key in the text by the rule of the number search ([the format text](../spec/format/README.md), section 20). The script checked the rule against all the examples of this section of the format.

### 5.6. standalone at 42b4408

For architecture 3.5; the 2026-11-08 amendment to R05.

[mathtrail/mathtrail-standalone](https://github.com/mathtrail/mathtrail-standalone/tree/42b4408225e4163d8efbdfb93c299f3fc4b70afa), `main`, commit 42b4408 of 2026-11-07, the newest on 2026-11-08. Cloned into the session's temporary folder and read 2026-11-08:
- `git log 89986bc..42b4408` — seven commits "Tasks solved from the picture", T105.3–T105.12 standalone. `git diff 89986bc 42b4408` changes neither `LICENSE`, nor `CONTRIBUTING.md`, nor `content/catalogs/traps.json`, nor `site/content/en/techniques.yaml`. `go.mod` has the same `go 1.27.3` and `go.starlark.net v0.0.0-20260908191801-89a6a09411d5`, and the test `TestEveryReferenceTaskProvesItsOwnAnswer` is still in `content/solver_test.go`. The test bench was not run at 42b4408 in this session. T18 ran it on 2026-11-08 on the 586 tasks of the run of `just port-survey`, and it proved each of them (R20);
- the count in architecture 3.5 is by jq over `content/examples/*.json`, `content/catalogs/`, `content/solvers/` and `content/ideas/` at both commits. The earlier tasks were checked by ID: the JSON of each of the 705 tasks at 42b4408 equals its JSON at 89986bc. A solver "reads the picture" if it has `def solve(options, picture`;
- `content/instructions/task_writing_picture.md`: "Its facts stand on the picture alone, and without the picture it cannot be solved", and `solve(options, picture)` receives "the description you wrote, as Starlark values";
- `content/schemas/task.json`: "A description of one of seventeen kinds of picture", and among them the new `lattice`, `tiles`, `shapes` and `stairs`;
- `SPEC.md`: "A topic all of whose tasks are solved from the picture says so with `"solved_from_picture": true`", and the rule of choice sets such a topic only there, "where the card is known to be drawn".

### 5.7. standalone's prototype

For architecture 3.5.

[MathTrail/llm-taskgen-prototype](https://github.com/MathTrail/llm-taskgen-prototype), commit 0263835, which standalone's `CLAUDE.md` names, read 2026-11-02 (`gh api`). The license is Apache-2.0, `LICENSE` is its text without a filled-in copyright line. There is no `NOTICE` file at the root. The author is RyazanovAlexander, 59 commits, there are no others.
