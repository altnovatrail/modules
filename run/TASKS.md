# TASKS.md — tasks outside an epic of `modules`

Open tasks outside the current epic, each as its body, in Russian (R154 spec, R155 spec).

- Never read whole: a session reads only its own task — `grep -n '^#### T' run/TASKS.md` lists them, `awk '/^#### T75\./{p=1;print;next} /^#### /{p=0} p' run/TASKS.md` gives one body.
- A new task outside an epic goes here as a body, its number the one free when it is written down (RUN.md, Citations). When done, the body leaves for the archive (RUN.md rule 9).

## Tasks outside an epic

#### T75. CI на релизе `.github` строки 7.13

- **Цель:** каждый `uses:` общих workflow `.github` — на релиз `.github` строки 7.13 COORDINATION.md spec, по SHA его тега: кеш Go в `go.yml` и параллельная сверка файлов релиза; кеш Go job «Modules» — папка `.cache/go/` (R36) — живёт между прогонами, в кеше GitHub Actions, по тем же правилам, что кеш Go `go.yml` этого релиза.
- **Почему:** строка 7.17 COORDINATION.md spec, часть, которую T74 не могла сделать: «переход на `.github` строки 7.13» и кеш, который от неё зависит («Зависит от: 7.13 — для кеша»). T74 сделала остальное (R36).
- **Зависит от:** релиз `.github` строки 7.13. Пока его нет — задача ждёт.
- **Сделать:**
  1. `gh release verify <тег> -R altnovatrail/.github`; SHA тега. `uses:` в `ci.yml`, `codeql.yml`, `scorecard.yml` и `release.yml` — на этот SHA с комментарием версии, как T73 (R35). Входы и права общих workflow от v0.1.24 до этого тега — по README `.github` и architecture 3.3 .github на нём; что изменилось, — в вызовы.
  2. Job «Modules»: `.cache/go/` — в кеш GitHub Actions по правилам кеша Go `go.yml` этого релиза (строка 7.13): ключ — Dockerfile и `go.sum`, pull request только читает, сохраняет push в `main`, после восстановления — `go mod verify`; `release.yml` кеша не берёт. Go пишет `GOMODCACHE` только для чтения: если папку чистит checkout или шаг, — `GOFLAGS=-modcacherw`. Если релиз даёт для этого общий шаг или action — взять его, а не свои шаги.
  3. Если вызовы этого релиза требуют от Dockerfile стадии `toolchain-core` (строка 7.14), они остаются T76 — `.devcontainer/` эта задача не меняет; остальное — здесь.
  4. Поправки к R35 и R36 — что взято и почему; правило `ci`; [docs/architecture.md](../docs/architecture.md), разделы 3.2 и 3.3.
- **Как проверить:** CI pull request зелёный, все обязательные проверки и `codecov/patch`; в логе job «Modules» и «Go / test» — кеш Go из кеша; время изменённых job до и после — в отчёте, по прогонам `main` и pull request; проверки части `ci` — [.claude/rules/ci.md](../.claude/rules/ci.md); `just -f spec/justfile docs-check-repo`; `/code-review high`, один круг.
- **Риски:** `.github/workflows/release.yml` — путь релиза: мерж выпускает следующий релиз `modules` без новых архивов (R17), и pull request и отчёт его называют.
- **Разделы:** строки 7.13 и 7.17 COORDINATION.md spec; measurements 3.1 spec; R76 spec; R17, R35, R36; architecture 3.3 .github.

#### T76. Переход на тег `spec` строки 7.14

- **Цель:** `modules` стоит на теге `spec` строки 7.14 COORDINATION.md spec (T166 spec): стадия `toolchain-core` в `.devcontainer/Dockerfile` по заготовке, apt — только из источников снимка; job, которым по релизу `.github` строки 7.13 хватает `toolchain-core`, — в ней.
- **Почему:** строка 7.17 COORDINATION.md spec, последняя часть её задач: переход на тег `spec` со стадией `toolchain-core` меняет `.devcontainer/`, и его делает последняя из них (T74). Переход — задача и решение журнала (RULES.md spec, Across repositories).
- **Зависит от:** тег `spec` строки 7.14; T75.
- **Сделать:**
  1. Как T72: `gh release verify <тег> -R altnovatrail/spec`; подмодуль `spec/` — на тег; тег — в CLAUDE.md, README, начале RUN.md и статусе [docs/architecture.md](../docs/architecture.md); что теги после v0.18.3 несут этому репозиторию — `git -C spec log --oneline v0.18.3..<тег>`, строки таблицы 4 `spec` с `` `modules` ``, карточка architecture 3.10 spec, RULES.md spec.
  2. Заготовка `just new-repo-draft modules go 3.10 <тег>` на клоне `spec` во временной папке: `.devcontainer/` — по ней, со стадией `toolchain-core`, с отличиями, которые держат решения журнала, — `module-kit` (R10) и генератор `canvas` (R29).
  3. Вызовы общих workflow, которые T75 оставила до `toolchain-core`, — здесь; свои job `ci.yml` — в `toolchain-core`, если она держит их инструменты; что взять — решает сессия, и называет в отчёте.
  4. Решение журнала о переходе, поправки к записям, которые он меняет; правило `ci`; строки таблиц 3 и 4, которые тег закрывает.
- **Как проверить:** `just -f spec/justfile docs-check-repo` на новом теге; `just test`; проверки частей, которые задача меняет, — правило 6 RUN.md; CI pull request зелёный; время изменённых job до и после — в отчёте; `/code-review high`, один круг.
- **Мой шаг:** задача меняет `.devcontainer/`: после мержа автор пересобирает контейнер, и оркестратор ждёт этого перед следующей задачей.
- **Разделы:** строки 7.14 и 7.17 COORDINATION.md spec; R156 spec, R76 spec; RULES.md spec; repositories 2.2 spec; R10, R29, R34, R36.

