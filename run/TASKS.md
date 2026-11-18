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
  3. Поправка к R36 и к R41; правило `ci`, `docs/architecture.md`, раздел 3.2; запись источника о действии, если её не хватает.
- **Как проверить:** проверки pull request зелёные, «Modules» — с действием; правило `ci`: `just test`; `just -f spec/justfile docs-check-repo`; `/code-review high`, один круг.
- **Разделы:** R36, R41; R49 `.github`, поправка T75 к R41 `.github`; sources 4.23.
