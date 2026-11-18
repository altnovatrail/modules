# TASKS.md — tasks outside an epic of `modules`

Open tasks outside the current epic, each as its body, in Russian (R154 spec, R155 spec).

- Never read whole: a session reads only its own task — `grep -n '^#### T' run/TASKS.md` lists them, `awk '/^#### T76\./{p=1;print;next} /^#### /{p=0} p' run/TASKS.md` gives one body.
- A new task outside an epic goes here as a body, its number the one free when it is written down (RUN.md, Citations). When done, the body leaves for the archive (RUN.md rule 9).

## Tasks outside an epic


#### T82. Переход на тег `spec` v0.18.27

- **Цель:** `modules` стоит на последнем теге `spec`, v0.18.27, — его правила, план и документы по этому тегу.
- **Почему:** R156 spec — параллельная работа волны 4 (слово автора 2026-11-18: пока `app` идёт эпиком E5, остальные репозитории не простаивают); переход — задача и решение журнала (RULES.md spec, Across repositories). Задача вне эпика или в текущем эпике — по правилам этого RUN.md.
- **Сделать:**
  1. `gh release verify v0.18.27 -R altnovatrail/spec`; подмодуль `spec/` — на v0.18.27; где репозиторий берёт модуль Go или пакет npm `spec` — та же версия, пакет npm — через `release-verify-npm` (R77 spec); тег — везде, где этот репозиторий его называет.
  2. Что теги после v0.18.14 до v0.18.27 несут этому репозиторию: `git -C spec log --oneline v0.18.14..v0.18.27`; строки таблицы 4 `spec` с именем репозитория в обратных кавычках — `grep -n '`modules`' spec/run/DEBT.md`; карточка этого репозитория в architecture 3 spec; RULES.md spec. Что из них берёт эта задача, а что — задача своя, решает сессия; отложенное — строкой таблицы 4.
  3. Заготовка `spec/repo-template/` (repositories 2.2 spec) на клоне `spec` v0.18.27 во временной папке: где она изменилась после v0.18.14 — README, `docs/README.md`, `.claude/`, CLAUDE.md, правила RUN.md —, взять, с отличиями, которые держат решения этого журнала.
  4. Общие workflow и действие `toolchain` `.github` — с v0.1.31 на релиз v0.1.33 (по SHA его тега, как сейчас): v0.1.33 восстанавливает кеши Go по префиксу без хеша `go.sum` (строка 2.7 COORDINATION.md spec, T75 `.github`), а v0.1.32 — составное действие `toolchain` для своих заданий (T69 `.github`). Проверки этого pull request — на новом SHA, зелёные.
  5. Решение журнала о переходе — что тег принёс и что взято —, поправки к записям, которые он меняет; строки таблиц 3 и 4, которые тег закрывает. Больше одного ревью — разбить по правилу 11 и назвать следующие задачи в отчёте.
- **Как проверить:** `just -f spec/justfile docs-check-repo` на новом теге — зелёный; проверки частей, которые задача меняет, — правило 6 этого RUN.md; `/code-review high`, один круг.
- **Разделы:** R156 spec, R76 spec; RULES.md spec; repositories 2.2 spec; карточка этого репозитория в architecture 3 spec.
- **Выяснилось при планировании:** между v0.18.14 и v0.18.27 `.devcontainer/` менялся только в стадии `orchestrator`, которой у копий нет (T192 spec), а `repo-template/kinds/*/Dockerfile` — нет: пересборка контейнера не нужна, `.devcontainer/` задача не трогает.

#### T83. Job «Modules» — на составном действии `toolchain` `.github`

- **Цель:** собственный job «Modules» `ci.yml` берёт образ стадии `toolchain` и кеши Go составным действием `.github/actions/toolchain` `.github` v0.1.33 или новее — без своей копии этих шагов, с восстановлением кешей Go по префиксу без хеша `go.sum`, как `go.yml`.
- **Почему:** R41: `.github` v0.1.32 дал действие для собственных заданий вызывающего (R49 `.github`), v0.1.33 — восстановление по префиксу (строка 2.7 COORDINATION.md spec); шаги «Modules» — копия правил `go.yml` v0.1.31 (поправка T75 к R36), и перевод на действие убирает её. Отдельной задачей: это меняет то, что запускает CI.
- **Сделать:**
  1. «Modules»: шаги образа и кешей Go — двумя вызовами действия, `altnovatrail/.github/.github/actions/toolchain@<sha> # vX.Y.Z` по SHA релиза, с `check: modules`: до работы и после неё с `save: "true"`; Dependabot поднимает его группой `shared-workflows`.
  2. Шаги рецептов: `toolchain.sh` монтирует `RUNNER_TEMP` только на чтение, а кеши Go действия лежат в `${{ runner.temp }}/go-cache` — каждый шаг запускает `docker run` с этой папкой на `/go-cache` и `GOMODCACHE`, `GOCACHE`, как в комментарии `action.yml`, или другой способ, который сессия выберет и запишет; `.cache/go/` уходит.
  3. Поправка к R36 и к R41; правило `ci`, `docs/architecture.md`, раздел 3.2; запись источника о действии, если её не хватает.
- **Как проверить:** проверки pull request зелёные, «Modules» — с действием; правило `ci`: `just test`; `just -f spec/justfile docs-check-repo`; `/code-review high`, один круг.
- **Разделы:** R36, R41; R49 `.github`, поправка T75 к R41 `.github`; sources 4.23.
