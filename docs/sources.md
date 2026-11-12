# modules — источники

> Внешние факты, на которых стоят документы `docs/`: условия GitHub, стандарты, документация других программ — репозиториев NovaTrail — и тексты MathTrail. Здесь их записи — адрес, версия или дата страницы, дата прочтения и цитата на языке источника, — а [architecture.md](architecture.md) держит утверждение и ссылку на запись по её номеру, как на раздел: `(sources 2.1)` (RULES.md, The bar; R143 spec).
>
> Записи разложены по темам. До T64 они стояли пунктами в architecture 4. Там остался указатель: на тот раздел ссылаются журнал и архивы `docs/plan/`, а их не переписывают (R31). Итоги замеров — в [measurements.md](measurements.md). Файл на русском, пока его не переведёт T65.

## 1. Как вести

- **Запись** — подраздел «N.M» темы: один источник или одно чтение нескольких источников ради одного решения. В ней адрес, версия или дата страницы, дата прочтения и цитата на языке источника. Первая строка записи называет, для каких разделов и решений она прочитана. Что прочитано не из первоисточника — копия, указатель, страница, которая не открылась, — запись так и называет.
- **Ссылка.** Утверждение документа ссылается на запись по номеру, как на раздел любого документа: `(sources 4.8)`. Номер записи не меняется: новая запись встаёт в конец своей темы, новая тема — в конец файла.
- **Источник прочитан снова** — в запись добавляется пункт со своей датой, а прежний текст не правится: на нём стоят решения своего дня. Если источник изменился и утверждение с ним разошлось, это поправка к решению или открытый вопрос, а не правка записи.
- **Чего здесь нет:** цитат, которые запись журнала решений приводит сама, — они остаются в записи (RULES.md, How we work: решение не переписывается), — и цитат, которые стоят в разделах 1–3 [architecture.md](architecture.md) доводом своего пункта: они остаются там, где стоят, как в T133 `spec`.

## 2. GitHub

### 2.1. Файлы релиза GitHub

Для architecture 3.3.

[About releases](https://docs.github.com/en/repositories/releasing-projects-on-github/about-releases), раздел «Storage and bandwidth quotas», прочитано 2026-11-02: «Up to 1000 release assets may be associated with a single release. Each file included in a release must be under 2 GiB.» У файла релиза есть только имя, без папок ([Upload a release asset](https://docs.github.com/en/rest/releases/assets#upload-a-release-asset), там же: «GitHub renames asset filenames that have special characters, non-alphanumeric characters, and leading or trailing periods»).

### 2.2. Неизменяемые релизы

Для architecture 3.3.

[Managing releases in a repository](https://docs.github.com/en/repositories/releasing-projects-on-github/managing-releases-in-a-repository), прочитано 2026-11-02: «If you have enabled immutable releases for your repository, you cannot add, replace, or delete assets after a release is published, and you cannot move or delete its tag while the release exists.»

### 2.3. Фильтры `paths` и обязательные проверки

Для architecture 3.2.

[Workflow syntax for GitHub Actions](https://docs.github.com/en/actions/reference/workflows-and-actions/workflow-syntax), `on.<push|pull_request>.paths`, прочитано 2026-11-02: «If a workflow is skipped due to path filtering, branch filtering, or a commit message, then checks associated with that workflow will remain in a "Pending" state. A pull request that requires those checks to be successful will be blocked from merging.»

### 2.4. Возраст аккаунта GitHub

Для architecture 3.4; R04.

[GitHub Terms of Service](https://docs.github.com/en/site-policy/github-terms/github-terms-of-service), B. Account Terms, 3. Account Requirements, действуют с 27 апреля 2026, прочитано 2026-11-02: «You must be age 13 or older.»

### 2.5. Что остаётся после push

Для architecture 3.4; R04.

[Removing sensitive data from a repository](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/removing-sensitive-data-from-a-repository), раздел «About sensitive data exposure», прочитано 2026-11-02: «If you only rewrite your history and force push it, the commits with sensitive data may still be accessible elsewhere: In any clones or forks of your repository; Directly via their SHA-1 hashes in cached views on GitHub; Through any pull requests that reference them». И там же: «you can permanently remove cached views and references to the sensitive data in pull requests on GitHub by contacting us through the GitHub Support portal».

## 3. Стандарты и библиотеки

### 3.1. `-text` в `.gitattributes`

Для architecture 3.1 и 3.3.

[gitattributes](https://git-scm.com/docs/gitattributes), атрибут `text`, прочитано 2026-11-02: «Unsetting the `text` attribute on a path tells Git not to attempt any end-of-line conversion upon checkin or checkout.» А `text` включает перевод концов строк — «enables end-of-line conversion on checkin and checkout» — у любого файла, на который стоит, медиафайлы тоже.

### 3.2. DCO 1.1

Для architecture 3.4; R04.

[Developer Certificate of Origin 1.1](https://developercertificate.org/), прочитано 2026-11-02:
- (a) «The contribution was created in whole or in part by me and I have the right to submit it under the open source license indicated in the file»;
- (b) «The contribution is based upon previous work that, to the best of my knowledge, is covered under an appropriate open source license and I have the right under that license to submit that work with modifications…»;
- (c) «The contribution was provided directly to me by some other person who certified (a), (b) or (c) and I have not modified it»;
- (d) «a record of the contribution (including all personal information I submit with it, including my sign-off) is maintained indefinitely».

### 3.3. Чанки PNG

Для architecture 3.4.

[Portable Network Graphics (PNG) Specification (Third Edition)](https://www.w3.org/TR/png-3/), W3C Recommendation 24 June 2025, прочитано 2026-11-02:
- «A critical chunk is a chunk that is absolutely required in order to successfully decode a PNG image from a PNG datastream» (там — 11.2). Критических чанков четыре — IHDR, PLTE, IDAT и IEND, а остальные вспомогательные (там — 4.8.2);
- PLTE нужен не всякой картинке: «This chunk shall appear for color type 3, and may appear for color types 2 and 6; it shall not appear for color types 0 and 4» (там — 11.2.2);
- среди вспомогательных — tEXt, iTXt и zTXt, «Textual information (which may be compressed) associated with the image», eXIf, «Exchangeable image file format metadata», и tIME, «The time when the PNG image was last modified» (там — 4.7, таблица 2).

### 3.4. Шрифт и рисование

Для architecture 3.5; R12.

`golang.org/x/image` v0.47.0: `font/gofont/goregular`, `font/opentype` и `vector`, лицензия BSD-3-Clause (`LICENSE` модуля), прочитано 2026-11-06; `image/png` стандартной библиотеки Go пишет непрозрачное изображение без вспомогательных чанков — тест `TestRowDraws` держит, что их нет.

## 4. Репозитории NovaTrail

### 4.1. Общий `release.yml` `.github`

Для architecture 3.3.

релиз v0.1.15, коммит 89d3cee, последний на 2026-11-02, прочитано 2026-11-02:
- `.github/workflows/release.yml`: `paths` читается по строке на pathspec и идёт в `git diff --quiet "$previous" "$TARGET" -- "${paths[@]}"`. Рецепт вызывающего идёт как `just release-files` в контейнере стадии `toolchain` его `.devcontainer/Dockerfile`, на дереве из `git archive` и без токенов. В `dist/` — только файлы с именами по `^[A-Za-z0-9_+][A-Za-z0-9._+-]*$`, и каждый файл подписан bundle `<файл>.sigstore.json` рядом;
- architecture 3.3 .github: вход `previous-file: module-index.json` у `modules` — «с `v0.1.14`», а рецепт тогда берёт «два параметра, `release-files version previous:`». Первый релиз получает вторым пустую строку. Перечни `paths` в шаблоне — вида go и вида typescript; вида для `modules` нет.

### 4.2. Шаблон `.github` на 2026-11-06

Для architecture 3.3; поправка 2026-11-06 к R03.

`docs/architecture.md` `.github`, `main`, коммит af40f94, последний релиз — v0.1.15, прочитано 2026-11-06 (`gh api`). В architecture 3.3 .github, «Вызывающие файлы», перечня `paths` вида `modules` нет: ни `withdrawals.yaml`, ни `modules/**` в файле не встречаются. О релизе `modules` раздел говорит только вход `previous-file: module-index.json`, «с `v0.1.14`».

### 4.3. Общий `release.yml` и шаблон на 2026-11-07

Для architecture 3.3; R17.

`.github`, релиз v0.1.15, коммит 89d3cee, и `main`, коммит af40f94, прочитано 2026-11-07 (`gh api`). Новее v0.1.15 релиза нет, а `main` впереди на три коммита документов: `.github/workflows/release.yml` на обоих — один и тот же файл. Перечня вида `modules` в шаблоне architecture 3.3 .github по-прежнему нет, есть только вид go. Входы — `paths`, `publish` («false builds the files, signs them and checks the signatures, and publishes nothing»), `npm-package` и `previous-file`. Сборка: `git archive "$TARGET"`, затем `docker build --target toolchain -` из `.devcontainer/Dockerfile` без контекста, затем `docker run … just release-files "${arguments[@]}"`, где вторым аргументом идёт `""` или `/previous/module-index.json`. В контейнер не идут ни токен, ни переменные шага, только `HOME=/tmp`, а `spec/` в дереве пуст. Сети рецепту не нужно, но и не закрыто: у `docker run` нет `--network none`. Pull request с форка пропускает пробу, потому что для подписи cosign без ключа нужен `id-token: write`.

### 4.4. `release.yml` `canvas`

Для architecture 3.3; R17.

`.github/workflows/release.yml` на `main` `canvas`, прочитано 2026-11-07 (`gh api`): `pull_request:` без своего перечня путей, условие «`github.event_name != 'workflow_run' || github.event.workflow_run.conclusion == 'success'`» и `publish: ${{ github.event_name != 'pull_request' }}`. Пробные прогоны на pull request 2026-11-02 — из ветки автора и Dependabot — зелёные, кроме одного (API `actions/workflows/release.yml/runs?event=pull_request`).

### 4.5. Общие workflow `.github` v0.1.15

Для architecture 3.2; R10.

`.github/workflows/`, релиз v0.1.15, прочитано 2026-11-06 (`gh api`): `documents.yml` запускает `docs-check-repo`, `go.yml` и `typescript.yml` — рецепты `ci-lint`, `ci-test`, `ci-vuln` и `ci-licenses` вызывающего. Workflow, который запускает рецепты вызывающего на выбор, среди них нет. Общий `codeql.yml` принимает вход `languages` со значением `javascript-typescript` и на v0.1.0, ff2a359, и на v0.1.10, f3addb7, — на неё `codeql.yml` `modules` поднял Dependabot 2026-11-06.

### 4.6. Проверка DCO

Для architecture 3.4.

`.github/workflows/dco.yml` `.github`, `main`, коммит af40f94, прочитано 2026-11-02: «The sign-off has to name the author exactly, by name and address». Код сверяет трейлер `Signed-off-by` каждого коммита, кроме merge-коммитов, с автором коммита, `%an <%ae>`, а кто открыл pull request, смотрит только ради Dependabot. CI `modules` вызывает этот workflow на v0.1.15, коммит 89d3cee, — Dependabot поднял его с v0.1.0 (#8), — и правило там то же: подпись автора `%an <%ae>`, исключение Dependabot и пропуск merge-коммитов (`dco.yml` на 89d3cee, прочитано 2026-11-07, `gh api`).

### 4.7. CONTRIBUTING `.github`

Для architecture 3.4.

`CONTRIBUTING.md`, `main` `.github`, коммит af40f94, прочитано 2026-11-02: «No samples for eval sets yet. Handwriting, voice and photo samples for eval sets are not accepted until the project decides how the people who make them consent to their publication.» Своего CONTRIBUTING у `modules` нет.

### 4.8. Релиз `module-kit` v0.2.3

Для architecture 3.2; R10.

[altnovatrail/module-kit, релиз v0.2.3](https://github.com/altnovatrail/module-kit/releases/tag/v0.2.3), коммит e069943, последний на 2026-11-06, неизменяемый (`gh release view`, `isImmutable: true`), прочитано 2026-11-06:
- бинарник linux лежит в релизе как есть, без архива: `module-kit_v0.2.3_linux_amd64` и `module-kit_v0.2.3_linux_arm64`. Суммы всех файлов — в `SHA256SUMS`, у каждого файла — `<файл>.sigstore.json`;
- `just -f spec/justfile release-verify module-kit v0.2.3 SHA256SUMS`: «SHA256SUMS: module-kit v0.2.3, signed by the shared release.yml of .github v0.1.15»;
- `cmd/module-kit/main.go`: `--version` печатает `module-kit %s` и второй строкой `Unicode %s`;
- architecture 3.6 module-kit, «Команды `check` и `build`»: «вызов — `module-kit check <модуль>` и `module-kit build <модуль> <папка>`», а для `build` — «папка не существует или пуста, иначе код 2 и ничего не тронуто».

### 4.9. Релиз `module-kit` v0.3.0

Для architecture 3.2; R28.

[altnovatrail/module-kit, релиз v0.3.0](https://github.com/altnovatrail/module-kit/releases/tag/v0.3.0), первый линии 0.3, неизменяемый (`gh release view`, `isImmutable: true`), прочитано 2026-11-09:
- R52 module-kit на теге v0.3.0: «Библиотека читает формат модуля 0.4, и только его: модуль 0.3 отклоняется по `manifest.format-version`»;
- `just module-kit-verify`: «SHA256SUMS: module-kit v0.3.0, signed by the shared release.yml of .github v0.1.15»; суммы `module-kit_v0.3.0_linux_amd64` и `…_linux_arm64` — в `ARG` Dockerfile;
- формы вызова `check`, `build`, `index` и `diff` в `docs/architecture.md` `module-kit` — те же, что на v0.2.3 (`gh api` на обоих тегах).

### 4.10. Команда `index`

Для architecture 3.3; R17.

architecture 3.6 module-kit, релиз v0.2.3, коммит e069943, прочитано 2026-11-07 (`gh api`; `main` совпадает с тегом): «module-kit index (--previous <module-index.json> | --first) --tag vX.Y.Z (--withdrawals <withdrawals.yaml> | --no-withdrawals) --out <файл> <папка пакета>...». Входы — «папки пакетов, которые написал `build`: каждая — `modulepackage.Read` по дайджесту байтов своего `inventory.json`, с `device/`». Вывод — «`--out`, которого ещё нет». Коды: 0 — индекс записан, 1 — отказ по правилам индекса, 2 — команда не смогла сделать своё. Архив пакета `module-kit` не пишет: «Архив для релиза кладёт рецепт `modules` (R03 modules)».

### 4.11. Команда `diff`

Для architecture 3.2; R11.

architecture 3.5 module-kit, релиз v0.2.3, прочитано 2026-11-06: «module-kit diff [--json] (--index <module-index.json> | --first) [--published <папка пакета>]... <модуль>»; `<модуль>` — «папка модуля или папка пакета»; `--published` — «папки пакетов файлами …, без `device/`», и «пакет каждой версии, которую называет индекс, нужен: без него сравнить нельзя»; код 2 — «нет пакета версии, которую называет индекс». Проба на v0.2.3 в тот же день: та же версия с тем же дайджестом — код 0, «the newest published version as it was»; с другим — код 1, `module-version.above`; новая версия — 0; без пакета — 2; `--first` — 0.

### 4.12. `module-kit` о прошлой версии

Для architecture 3.3.

architecture 2.4 module-kit, `main` `module-kit`, тег v0.1.1, коммит 3f5a4c4, прочитано 2026-11-02: «Откуда CI `modules` берёт прошлую версию модуля для сравнения, решает T01 modules». CLI на этом коммите знает только `--version` (`cmd/module-kit/main.go`).

### 4.13. Версия модуля и eval-набор

Для architecture 3.4.

T04 в RUN.md `module-kit`, «Правила semver для `diff`», `main` `module-kit`, коммит 2489d19, прочитано 2026-11-02: «Перечень изменений по составу модуля — … eval-набор (concepts 3 spec) — и какая часть версии растёт от каждого». Задача не начата.

### 4.14. Воспроизводимый архив

Для architecture 3.3; R17.

рецепт `release-files` `tutor`, `main`, прочитано 2026-11-07 (`gh api`): «tar --create --format=gnu --sort=name --owner=0 --group=0 --numeric-owner --mode=go-w --mtime=… | gzip -n -9». В стадии `toolchain` `modules` — GNU tar 1.35, и `--sort=name` в нём есть (`tar --version` и `tar --help`, 2026-11-07).

### 4.15. Где раннер берёт образцы

Для architecture 3.4.

R49 tutor, `main` `tutor`, коммит be607e5, прочитано 2026-11-02: «eval-раннер берёт задачи из пакета модуля, а образцы — из папки модуля, ведь eval-набор в пакет не входит».

### 4.16. Просьба `api` к раскладке релиза

Для architecture 3.3.

В6 в таблице 3 RUN.md `api`, «в реестр spec», к T01 `modules`; `main` `api`, коммит d409351, прочитано 2026-11-02: файлы пакета находят «по записи — ID модуля, версии и тегу, — без перебора файлов релиза: по публичному адресу, без токена и REST API». Рекомендация `api` — «один архив на пакет, названный по дайджесту пакета», и исправлять снятую мажорную версию новой версией той же мажорной в том же релизе, который её снимает.

### 4.17. `sim` о гейте графа

Для architecture 3.3.

T03 в RUN.md `sim`, `main` `sim`, коммит 067f624, прочитано 2026-11-02: «CI `modules`: как он узнаёт, что граф изменился, и на каком пакете идёт стенд».

### 4.18. Файл образца и страница образцов

Для architecture 3.4.

`canvas`, `main`, коммит 5cd33ce, прочитано 2026-11-02:
- R09 canvas: у образца eval-набора нет кода писавшего, а о человеке в файле нет «ни имени, почты, аккаунта GitHub, ни устройства, браузера, экрана, языка, ни даты, возраста и класса»; контейнер — «один JSON на образец», вырезки — «PNG в base64»;
- R11 canvas: «eval-набор закрыт в каждой сборке, пока нет утверждённого текста по П-16»; «Версию набора текстов и «кто писал» сверяют организатор замера, раннер `tutor` и CI `modules`»; «кто открывает pull request, решает T02 modules»;
- architecture 3.5 canvas: «У образца eval-набора имя на образце ищет T02 modules»;
- architecture 3.6 canvas: на странице — заявления «мне 18 лет или больше» и «форма согласия версии … подписана и принята организатором», а «куда передать файл: … образец eval-набора — pull request по правилам T02 modules».

### 4.19. Раскладка рабочей тетради

Для architecture 3.8; R15.

[altnovatrail/canvas, релиз v0.1.2](https://github.com/altnovatrail/canvas/releases/tag/v0.1.2), коммит 8e9c6f6 от 2026-11-02, новейший на 2026-11-07, прочитано 2026-11-07, клон во временной папке сессии:
- `src/core/page.ts`: «A task whose lines do not fit is a TypeError: nothing shrinks», а `pageOf` — «a function of the items alone»;
- `vectors/pages.json`: кегль 30, окошко 100 × 100, поля 40, а короткий пункт — «no wider than the column (440) less the indent (56)». Длинный пункт заполняет строки по 864, а ряд ниже 1400 − 40 — отказ `task.items[i]`;
- `pageOf` во временной папке с Node 24.21: все 19 случаев `vectors/pages.json` совпали. На 12 страницах набора отказа нет, а самая длинная кончается на 932 из 1360.

## 5. MathTrail: standalone и его прототип

### 5.1. Лицензия standalone

Для architecture 3.5.

[mathtrail/mathtrail-standalone, `LICENSE`](https://github.com/mathtrail/mathtrail-standalone/blob/c6cbeb2/LICENSE), ветка `main`, коммит c6cbeb2, прочитано 2026-10-19: «MIT License», «Copyright (c) 2026 MathTrail». Других файлов с именем LICENSE, COPYING или NOTICE в репозитории нет; GitHub определяет его лицензию как MIT (`gh repo view`, `licenseInfo`, 2026-10-19).

### 5.2. standalone на 89986bc

Для architecture 3.5; R05.

[mathtrail/mathtrail-standalone](https://github.com/mathtrail/mathtrail-standalone/tree/89986bc), `main`, коммит 89986bc от 2026-11-01. Клонирован во временную папку сессии, прочитан 2026-11-02:
- `LICENSE`: «MIT License», «Copyright (c) 2026 MathTrail»; `CONTRIBUTING.md`: «The code and the content are MIT, and a contribution arrives under the same licence»;
- счёт частей и ответов в architecture 3.5 — jq по `content/examples/*.json`, `content/catalogs/` и `content/solvers/`. Вид ответа — по тексту верного варианта: целое число, `a/b`, число с `%`, время `ч:мм`, «Yes» или «No», другой текст с цифрами, текст без цифр. Пояснение «называет ответ», если в нём стоит число ключа, а в условии его нет, — приблизительно по правилу поиска числа [текста формата](../spec/format/README.md), раздел 20;
- 705 задач: `SPEC.md` — «Grades 1–4 — 450 tasks, 442 of them the prototype's» и «Grades 5–6 — 153 new tasks», а ещё 102 добавил R309 standalone. Про задачи из прототипа: «They are ported as they are (T23)… Their questions stay word for word»;
- свои формулировки: `SPEC.md` — «Nobody else's wording, anywhere. The ideas come from open sources; the text is ours». `CONTRIBUTING.md` — «The wording is your own. Ideas from open sources are fine — a wording lifted from a book or an olympiad is not, whatever its age». `site/content/en/research.yaml` — «from a book only a task's idea and its answer are taken, and the wording, the story, the options and the traps are written anew»;
- идеи 5–6 классов, R309 standalone (`docs/decisions.md`): «The ideas come from MathNet (mathnet.mit.edu, under CC BY 4.0, the countries that assert rights to their problems keeping them) at `5-6`, in the topics it has a section for, and from the sources of the prototype's research 12 elsewhere. An idea is taken and never a wording»;
- каталоги в примерах формата — с тех пор те же: `git diff c6cbeb2 89986bc` пуст для `content/catalogs/topics.json` и `content/catalogs/traps.json`, а `git diff abf3b61 89986bc` — для `site/content/en/techniques.yaml`;
- авторы: `gh api repos/mathtrail/mathtrail-standalone/contributors`, 2026-11-02 — RyazanovAlexander, 153 коммита, и dependabot[bot], 12.

### 5.3. Задачи standalone по темам

Для architecture 3.3.

[mathtrail/mathtrail-standalone](https://github.com/mathtrail/mathtrail-standalone/tree/89986bc/content/examples), `content/examples/`, коммит 89986bc, прочитано 2026-11-02 (`gh api`): 17 файлов `<тема>.json`, по одному на тему, со всеми её задачами, а решатели — `content/examples/solvers/<id задачи>.star`. Папок по классам нет: полоса — в ID задачи и её полях.

### 5.4. Стенд standalone

Для architecture 3.5; R12.

[mathtrail/mathtrail-standalone](https://github.com/mathtrail/mathtrail-standalone/tree/89986bc01a12ff3188cf08194677bfae0464568a), коммит 89986bc, прочитано 2026-11-06 (`gh api`; клон во временной папке):
- `content/solver_test.go`, `TestEveryReferenceTaskProvesItsOwnAnswer`: «It runs each solver the way a submitted one is run — the same sandbox, the same dialect and helpers, the same two runs under rotated labels», подтест на ID задачи, а `prove` — «the solver proves %s, the task claims %s»;
- `internal/domain/solver/verdict.go`: «The second run happens only when the first reached exactly one option»;
- `go.mod`: `go 1.27.3` и `go.starlark.net v0.0.0-20260908191801-89a6a09411d5`;
- `content/schemas/task.json`, `$defs/row` и `row_item`: «Objects in a row … 2 to 12 items», `line`, `gaps`, `span`, `copies` 1–2, у предмета `label`, `below` и `mark` — `dot`, `ring` или `square`, а пропуск «is never first or last, and is never next to another skip»;
- прогон 2026-11-06: стенд доказал ответы всех 27 задач первой пачки, а задачи и две картинки `row` конвертера прошли `module-kit check` v0.2.3 в модуле из примеров формата без предупреждений, кроме двух пояснений, которые называют ключ и которые переписывает T12 (architecture 3.6).

### 5.5. Первая пачка по standalone

Для architecture 3.6; R07.

[mathtrail/mathtrail-standalone](https://github.com/mathtrail/mathtrail-standalone/tree/89986bc/site/content/en/topics), коммит 89986bc. Репозиторий клонирован во временную папку сессии и прочитан 2026-11-02:
- `site/content/en/topics/gaps-and-boundaries.yaml`, `enumeration.yaml` и `pigeonhole-principle.yaml` — `hero.lead`, карточки `basis`, `method.steps` и по три примера `examples` с `question`, `steps`, `answer` и `note`. В начале каждой страницы: «an example here is worked through to that answer»;
- `site/data.json`, член `examples` — уровень и решатель каждого примера. У трёх тем — по примеру на `1-2`, `3-4` и `5-6`;
- `content/examples/counting.gaps.json`, `combinatorics.enumeration.json` и `pigeonhole.basic.json` — задачи, их пояснения, картинки и подсказки;
- отбор и счёт — скриптом Node с js-yaml 4.1.0 в той же папке. Вид ответа определяется по тексту верного варианта, вопрос «Which of these…» — по словам, ключ в тексте — по правилу поиска чисел ([текст формата](../spec/format/README.md), раздел 20). Правило скрипт сверил со всеми примерами этого раздела формата.

### 5.6. standalone на 42b4408

Для architecture 3.5; поправка 2026-11-08 к R05.

[mathtrail/mathtrail-standalone](https://github.com/mathtrail/mathtrail-standalone/tree/42b4408225e4163d8efbdfb93c299f3fc4b70afa), `main`, коммит 42b4408 от 2026-11-07, новейший на 2026-11-08. Клонирован во временную папку сессии и прочитан 2026-11-08:
- `git log 89986bc..42b4408` — семь коммитов «Tasks solved from the picture», T105.3–T105.12 standalone. `git diff 89986bc 42b4408` не меняет ни `LICENSE`, ни `CONTRIBUTING.md`, ни `content/catalogs/traps.json`, ни `site/content/en/techniques.yaml`. В `go.mod` — те же `go 1.27.3` и `go.starlark.net v0.0.0-20260908191801-89a6a09411d5`, а тест `TestEveryReferenceTaskProvesItsOwnAnswer` по-прежнему в `content/solver_test.go`. Стенд на 42b4408 в этой сессии не запускался. T18 прогнал его 2026-11-08 на 586 задачах прогона `just port-survey`, и он доказал каждую (R20);
- счёт в architecture 3.5 — jq по `content/examples/*.json`, `content/catalogs/`, `content/solvers/` и `content/ideas/` на обоих коммитах. Прежние задачи сверены по ID: JSON каждой из 705 задач на 42b4408 равен её JSON на 89986bc. Решатель «читает картинку», если в нём стоит `def solve(options, picture`;
- `content/instructions/task_writing_picture.md`: «Its facts stand on the picture alone, and without the picture it cannot be solved», а `solve(options, picture)` получает «the description you wrote, as Starlark values»;
- `content/schemas/task.json`: «A description of one of seventeen kinds of picture», и среди них новые `lattice`, `tiles`, `shapes` и `stairs`;
- `SPEC.md`: «A topic all of whose tasks are solved from the picture says so with `"solved_from_picture": true`», и такую тему правило выбора задаёт только там, «where the card is known to be drawn».

### 5.7. Прототип standalone

Для architecture 3.5.

[MathTrail/llm-taskgen-prototype](https://github.com/MathTrail/llm-taskgen-prototype), коммит 0263835, который называет `CLAUDE.md` standalone, прочитано 2026-11-02 (`gh api`). Лицензия — Apache-2.0, `LICENSE` — её текст без заполненной строки copyright. Файла `NOTICE` в корне нет. Автор — RyazanovAlexander, 59 коммитов, других нет.
