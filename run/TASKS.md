# TASKS.md — tasks outside an epic of `modules`

Open tasks outside the current epic, each as its body, in Russian (R154 spec, R155 spec).

- Never read whole: a session reads only its own task — `grep -n '^#### T' run/TASKS.md` lists them, `awk '/^#### T76\./{p=1;print;next} /^#### /{p=0} p' run/TASKS.md` gives one body.
- A new task outside an epic goes here as a body, its number the one free when it is written down (RUN.md, Citations). When done, the body leaves for the archive (RUN.md rule 9).

## Tasks outside an epic

#### T83. Job «Modules» — на составном действии `toolchain` `.github`

- **Цель:** собственный job «Modules» `ci.yml` берёт образ стадии `toolchain` и кеши Go составным действием `.github/actions/toolchain` `.github` v0.1.33 или новее — без своей копии этих шагов, с восстановлением кешей Go по префиксу без хеша `go.sum`, как `go.yml`.
- **Почему:** R41: `.github` v0.1.32 дал действие для собственных заданий вызывающего (R49 `.github`), v0.1.33 — восстановление по префиксу (строка 2.7 COORDINATION.md spec); шаги «Modules» — копия правил `go.yml` v0.1.31 (поправка T75 к R36), и перевод на действие убирает её. Отдельной задачей: это меняет то, что запускает CI.
- **Сделать:**
  1. «Modules»: шаги образа и кешей Go — двумя вызовами действия, `altnovatrail/.github/.github/actions/toolchain@<sha> # vX.Y.Z` по SHA релиза, с `check: modules`: до работы и после неё с `save: "true"`; Dependabot поднимает его группой `shared-workflows`.
  2. Шаги рецептов: `toolchain.sh` монтирует `RUNNER_TEMP` только на чтение, а кеши Go действия лежат в `${{ runner.temp }}/go-cache` — каждый шаг запускает `docker run` с этой папкой на `/go-cache` и `GOMODCACHE`, `GOCACHE`, как в комментарии `action.yml`, или другой способ, который сессия выберет и запишет; `.cache/go/` уходит.
  3. R36 и R41 — переписать на месте, каждую целиком, с её поправками (R160 spec); правило `ci`, `docs/architecture.md`, раздел 3.2; запись источника о действии, если её не хватает.
- **Как проверить:** проверки pull request зелёные, «Modules» — с действием; правило `ci`: `just test`; `just -f spec/justfile docs-check-repo`; `/code-review high`, один круг.
- **Разделы:** R36, R41; R49 `.github`, поправка T75 к R41 `.github`; sources 4.23.

#### T85. Документы без истории (R160 spec)

- **Цель:** документы `docs/`, README.md, CLAUDE.md и правила `.claude/rules/` держат текущее состояние, а не историю (R160 spec, «One home», «A done task»).
- **Почему:** вторая часть T84 (строка 7.28 COORDINATION.md `spec`), отделённая по правилу 11 RUN.md: `docs/architecture.md` — 173 КБ; её абзац «Status», хронику тегов, T84 уже убрала.
- **Зависит от:** T84.
- **Сделать:**
  - `docs/architecture.md`: без хроник версий, статусов и повышений; документ называет файл, где задано значение — подмодуль `spec`, `ARG` Dockerfile, `uses:` workflow, `version` манифеста, — а не значение; устройство описано как есть, решение — номером, без пересказа.
  - README.md, `docs/README.md`, CLAUDE.md, `.claude/rules/` — так же; строки CLAUDE.md и RUN.md с тегом `spec` остаются (R160 spec, «Decisions»).
  - Ссылка на поправку записи `spec`, которую тег переписал на месте (R143, R145, R154, R155 spec и другие — по журналу `spec` на теге подмодуля), — на запись.
  - Журнал: переписать на месте только записи, которые эта задача меняет (R160 spec); их свой предел в `entries` `docs-budgets.json` — убрать, если запись стала не длиннее 12 КБ.
  - Вне правила: `docs/sources.md`, `docs/measurements.md`, архивы `docs/plan/`.
  - Больше одного ревью — разделить по правилу 11 RUN.md.
- **Как проверить:** `just -f spec/justfile docs-check-repo`; `/code-review high`.
- **Разделы:** R160 spec, R21 spec, R145 spec; RULES.md spec — The bar.

#### T86. Комментарии кода без истории (R160 spec)

- **Цель:** комментарии кода — `tools/port/`, `test/`, `justfile`, `.github/`, `.devcontainer/` — по R160 spec, «Comments».
- **Почему:** третья часть T84 (строка 7.28 COORDINATION.md `spec`), отделённая по правилу 11 RUN.md: около 2400 строк комментариев.
- **Зависит от:** T84.
- **Сделать:**
  - Комментарий говорит, чего не скажут код и имена: почему, ограничение, ловушка, контракт экспортируемого API; без пересказа кода, истории, номеров задач и разделов документов; номер решения — якорем в конце, `(R25)`, и он называет ту запись, что даёт причину: свой журнал без суффикса, журнал `spec` — `R48 spec`.
  - Комментарии модулей `modules/` — не здесь: их правка меняет пакеты и выпускает релиз (T87).
  - Больше одного ревью — разделить по правилу 11 RUN.md.
- **Как проверить:** `just -f spec/justfile docs-check-repo` (проверяет номера в комментариях); проверки изменённых частей по правилам `.claude/rules/` — `port`: проверки Go, `ci`: `just test`; `/code-review high`.
- **Разделы:** R160 spec, R32 spec.

#### T87. Комментарии модулей без истории (R160 spec)

- **Цель:** комментарии YAML модулей `modules/*/` — по R160 spec, «Comments», как T86 для кода.
- **Почему:** четвёртая часть T84 (строка 7.28 COORDINATION.md `spec`), отдельно от T86: около 3600 строк комментариев в опубликованных модулях; правка меняет пакет, и мерж выпускает релиз (R17).
- **Зависит от:** T84.
- **Сделать:**
  - Комментарии — по R160 spec, «Comments»; номер решения — якорем в конце, и он называет ту запись, что даёт причину.
  - Новая версия каждого изменённого модуля — по правилу `modules`, если `just modules-check` её требует; релиз — в отчёте.
  - Больше одного ревью — разделить по правилу 11 RUN.md.
- **Как проверить:** `just -f spec/justfile docs-check-repo`; проверки правила `modules`; `/code-review high`.
- **Разделы:** R160 spec, R32 spec; R17.
