# TASKS.md — tasks outside an epic of `modules`

Open tasks outside the current epic, each as its body, in Russian (R154 spec, R155 spec).

- Never read whole: a session reads only its own task — `grep -n '^#### T' run/TASKS.md` lists them, `awk '/^#### T73\./{p=1;print;next} /^#### /{p=0} p' run/TASKS.md` gives one body.
- A new task outside an epic goes here as a body, its number the one free when it is written down (RUN.md, Citations). When done, the body leaves for the archive (RUN.md rule 9).

## Tasks outside an epic

#### T73. CI на `.github` v0.1.24

- **Цель:** каждый `uses:` общих workflow `.github` в `.github/workflows/` — на последний релиз `.github`, v0.1.24, по SHA его тега: CI берёт образ `toolchain` из кеша GitHub Actions, а не строит его из сети в каждом job (Д146 spec), а Dependabot поднимает ссылки на `.github` одной группой (Д88 spec).
- **Почему:** следующая задача после перехода на `spec` v0.18.3 (T72, R34), как T68 `sim`. Сейчас `codeql.yml` и job документов `ci.yml` стоят на v0.1.10, остальные — на v0.1.15; Д146 spec называет `modules` среди репозиториев, которым переходить на `.github` v0.1.22 или новее, Д88 spec — среди тех, кому нужна группа `shared-workflows`. Два открытых pull request Dependabot — `codeql.yml` и `documents.yml` на v0.1.15 — поднимают по одной ссылке, и каждый — свой круг проверок.
- **Сделать:**
  1. `gh release verify v0.1.24 -R altnovatrail/.github`; SHA тега. `uses:` в `ci.yml`, `codeql.yml`, `scorecard.yml` и `release.yml` — на этот SHA с комментарием `# v0.1.24`. Входы и права общих workflow от v0.1.10 и v0.1.15 до v0.1.24 — по README `.github` и architecture 3.3 .github на этом теге; что изменилось, — в вызовы.
  2. `.github/dependabot.yml` — группа `shared-workflows`, `patterns: ["altnovatrail/.github/*"]`, под `github-actions`, по шаблону architecture 3.3 .github (Д88 spec).
  3. Поправка к записи журнала, которая называет версии `.github` этого репозитория, или своя запись — что взято и почему; [docs/architecture.md](../docs/architecture.md), раздел 3.3, — если он называет версию.
- **Как проверить:** CI pull request зелёный — все обязательные проверки, `codecov/patch` тоже; в логе job — образ `toolchain` из кеша; проверки части `ci` — [.claude/rules/ci.md](../.claude/rules/ci.md); `just -f spec/justfile docs-check-repo`; `/code-review high`, один круг.
- **Риски:** `.github/workflows/release.yml` — путь релиза: мерж выпускает следующий релиз `modules` без новых архивов (R17), и pull request и отчёт его называют. Pull request Dependabot на v0.1.15 после мержа закрываются сами или вручную — назвать в отчёте.
- **Разделы:** Д88 spec, Д146 spec; architecture 3.3 .github; R03, R17, R34; [.claude/rules/ci.md](../.claude/rules/ci.md).
