# TASKS.md — tasks outside an epic of `modules`

Open tasks outside the current epic, each as its body, in Russian (R154 spec, R155 spec).

- Never read whole: a session reads only its own task — `grep -n '^#### T' run/TASKS.md` lists them, `awk '/^#### T76\./{p=1;print;next} /^#### /{p=0} p' run/TASKS.md` gives one body.
- A new task outside an epic goes here as a body, its number the one free when it is written down (RUN.md, Citations). When done, the body leaves for the archive (RUN.md rule 9).

## Tasks outside an epic


#### T80. Переход на тег `spec` v0.18.14

- **Цель:** `modules` стоит на последнем теге `spec`, v0.18.14, — его правила, план и документы по этому тегу, — и остальные задачи эпика прогона COORDINATION.md (слово автора 2026-11-17) `spec` идут здесь на нём.
- **Почему:** R156 spec — параллельная работа волны 4 (слово автора 2026-11-17): пока `app` идёт эпиком E5, остальные репозитории переходят на последний тег `spec` и закрывают строки таблицы 4 `spec`, которые ждут этого перехода; переход — задача и решение журнала (RULES.md spec, Across repositories). Задача вне эпика или в текущем эпике — по правилам этого RUN.md.
- **Сделать:**
  1. `gh release verify v0.18.14 -R altnovatrail/spec`; подмодуль `spec/` — на v0.18.14; где репозиторий берёт модуль Go или пакет npm `spec` — та же версия, пакет npm — через `release-verify-npm` (R77 spec); тег — везде, где этот репозиторий его называет: CLAUDE.md, README, начало RUN.md, статус `docs/architecture.md`.
  2. Что теги после v0.18.4 до v0.18.14 несут этому репозиторию: `git -C spec log --oneline v0.18.4..v0.18.14`; строки таблицы 4 `spec` с именем репозитория в обратных кавычках — `grep -n '`modules`' spec/run/DEBT.md`; карточка этого репозитория в architecture 3 spec; RULES.md spec. Что из них берёт эта задача, а что — задача своя, решает сессия; отложенное — строкой таблицы 4.
  3. Заготовка `spec/repo-template/` (repositories 2.2 spec): `just new-repo-draft modules <вид> <карточка> v0.18.14` (вид — по `repo-template/kinds.json`, карточка — номер раздела architecture 3 spec) на клоне `spec` v0.18.14 во временной папке, как T56 `sim`. `.devcontainer/`, `.claude/`, CLAUDE.md, правила RUN.md и раскладка таблиц плана — по ней, где она изменилась, с отличиями, которые держат решения этого журнала.
  4. Решение журнала о переходе — что тег принёс и что взято —, поправки к записям, которые он меняет; строки таблиц 3 и 4, которые тег закрывает.
  5. Переход больше одного ревью — разбить по правилу 11: эта задача — первая часть, остальные — следующими задачами, и назвать их в отчёте.
- **Как проверить:** `just -f spec/justfile docs-check-repo` на новом теге — зелёный; проверки частей, которые задача меняет, — правило 6 этого RUN.md; `/code-review high`, один круг.
- **Мой шаг:** если задача меняет `.devcontainer/` — назвать это в отчёте: после мержа автор пересобирает контейнер, и оркестратор ждёт этого перед следующей задачей. Между v0.18.4 и v0.18.14 менялись `.devcontainer/Dockerfile`, `devcontainer.json` и `repo-template/kinds/go/Dockerfile` (итог E14 `spec`: хук `tools/session`, `shutdownAction`): изменение `.devcontainer/` согласовано автором заранее.
- **Разделы:** R156 spec, R76 spec; RULES.md spec; repositories 2.2 spec; карточка этого репозитория в architecture 3 spec.
