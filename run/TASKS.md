# TASKS.md — tasks outside an epic of `modules`

Open tasks outside the current epic, each as its body, in Russian (R154 spec, R155 spec).

- Never read whole: a session reads only its own task — `grep -n '^#### T' run/TASKS.md` lists them, `awk '/^#### T76\./{p=1;print;next} /^#### /{p=0} p' run/TASKS.md` gives one body.
- A new task outside an epic goes here as a body, its number the one free when it is written down (RUN.md, Citations). When done, the body leaves for the archive (RUN.md rule 9).

## Tasks outside an epic

#### T76. Переход на тег `spec` строки 7.14

- **Цель:** `modules` стоит на теге `spec` строки 7.14 COORDINATION.md spec (T166 spec): стадия `toolchain-core` в `.devcontainer/Dockerfile` по заготовке, apt — только из источников снимка; job, которым по релизу `.github` строки 7.13 хватает `toolchain-core`, — в ней.
- **Почему:** строка 7.17 COORDINATION.md spec, последняя часть её задач: переход на тег `spec` со стадией `toolchain-core` меняет `.devcontainer/`, и его делает последняя из них (T74). Переход — задача и решение журнала (RULES.md spec, Across repositories).
- **Зависит от:** тег `spec` строки 7.14; T75; для шага 3 — релиз `.github` с T66 `.github`, последней частью строки 7.13.
- **Сделать:**
  1. Как T72: `gh release verify <тег> -R altnovatrail/spec`; подмодуль `spec/` — на тег; тег — в CLAUDE.md, README, начале RUN.md и статусе [docs/architecture.md](../docs/architecture.md); что теги после v0.18.3 несут этому репозиторию — `git -C spec log --oneline v0.18.3..<тег>`, строки таблицы 4 `spec` с `` `modules` ``, карточка architecture 3.10 spec, RULES.md spec.
  2. Заготовка `just new-repo-draft modules go 3.10 <тег>` на клоне `spec` во временной папке: `.devcontainer/` — по ней, со стадией `toolchain-core`, с отличиями, которые держат решения журнала, — `module-kit` (R10) и генератор `canvas` (R29).
  3. Вызовы общих workflow, которые T75 оставила до `toolchain-core`, — здесь. T75 перевела все `uses:` на `.github` v0.1.30 (8311cbf), где есть T60–T65; T66 `.github` — `release.yml` строит стадию `toolchain-core` Dockerfile вызывающего без кеша, а проверки без браузера и Go-тестов (lint, vuln, licenses, документы) идут в ней, у Dockerfile без стадии — как прежде, — в v0.1.30 нет: ей нужен переход `.github` на тег `spec` строки 7.14 (T171 spec — ночью `.devcontainer/` не меняется). Если её релиз вышел — `gh release verify`, `uses:` в `ci.yml`, `codeql.yml`, `scorecard.yml` и `release.yml` — на его SHA, входы и права — по README и architecture 3.3 `.github` на нём, и как его проверки узнают стадию — по Dockerfile или входу; если нет — Dependabot, группа `shared-workflows`, или своя задача после него. Свои job `ci.yml` — в `toolchain-core`, если она держит их инструменты (`module-kit`, генератор `canvas`); ключ и файл её образа — по T61 `.github`, сохраняет только push в `main` (поправка T75 к R35); что взять — решает сессия, и называет в отчёте.
  4. Решение журнала о переходе, поправки к записям, которые он меняет; правило `ci`; строки таблиц 3 и 4, которые тег закрывает.
- **Как проверить:** `just -f spec/justfile docs-check-repo` на новом теге; `just test`; проверки частей, которые задача меняет, — правило 6 RUN.md; CI pull request зелёный; время изменённых job до и после — в отчёте; `/code-review high`, один круг.
- **Мой шаг:** задача меняет `.devcontainer/`: после мержа автор пересобирает контейнер, и оркестратор ждёт этого перед следующей задачей.
- **Разделы:** строки 7.14 и 7.17 COORDINATION.md spec; R156 spec, R76 spec; RULES.md spec; repositories 2.2 spec; R10, R29, R34, R36.

