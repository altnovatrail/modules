# TASKS.md — tasks outside an epic of `modules`

Open tasks outside the current epic, each as its body, in Russian.

- Never read whole: a session reads only its own task — `grep -n '^#### T' run/TASKS.md` lists them, `awk '/^#### T76\./{p=1;print;next} /^#### /{p=0} p' run/TASKS.md` gives one body.
- A new task outside an epic goes here as a body, its number the one free when it is written down (RUN.md, Citations). When done, the body leaves for the archive (RUN.md rule 9).

## Tasks outside an epic

#### T83. Job «Modules» — на составном действии `toolchain` `.github`

- **Цель:** собственный job «Modules» `ci.yml` берёт образ стадии `toolchain` и кеши Go составным действием `.github/actions/toolchain` `.github` v0.1.33 или новее — без своей копии этих шагов, с восстановлением кешей Go по префиксу без хеша `go.sum`, как `go.yml`.
- **Почему:** R41: `.github` v0.1.32 дал действие для собственных заданий вызывающего (R49 `.github`), v0.1.33 — восстановление по префиксу (строка 2.7 COORDINATION.md spec); шаги «Modules» — копия правил `go.yml` v0.1.31 (поправка T75 к R36), и перевод на действие убирает её. Отдельной задачей: это меняет то, что запускает CI.
- **Сделать:**
  1. «Modules»: шаги образа и кешей Go — двумя вызовами действия, `altnovatrail/.github/.github/actions/toolchain@<sha> # vX.Y.Z` по SHA релиза, с `check: modules`: до работы и после неё с `save: "true"`; Dependabot поднимает его группой `shared-workflows`.
  2. Шаги рецептов: `toolchain.sh` монтирует `RUNNER_TEMP` только на чтение, а кеши Go действия лежат в `${{ runner.temp }}/go-cache` — каждый шаг запускает `docker run` с этой папкой на `/go-cache` и `GOMODCACHE`, `GOCACHE`, как в комментарии `action.yml`, или другой способ, который сессия выберет и запишет; `.cache/go/` уходит.
  3. R36 и R41 — переписать на месте, каждую целиком (R160 spec), по строгому правилу ссылок `spec/RULES.md` — поправок в них уже нет; правило `ci`, `docs/architecture.md`, раздел 3.2; запись источника о действии, если её не хватает.
- **Как проверить:** проверки pull request зелёные, «Modules» — с действием; правило `ci`: `just test`; `just -f spec/justfile docs-check-repo`; `/code-review high`, один круг.
- **Разделы:** R36, R41; R49 `.github`, поправка T75 к R41 `.github`; sources 4.23.

#### T91. Журнал без переходов на тег (R160 spec)

- **Цель:** журнал `docs/decisions.md` держит решения в силе, а не хронику переходов на теги `spec`: ни одна запись не говорит «`modules` stands on `spec` v…» и «`.devcontainer/` — заготовка v…», и ни одна не ссылается на поправку `spec`, которую тег свернул в её запись.
- **Почему:** R160 spec — переход на тег — задача и её коммит, не запись журнала; запись с датированными поправками переписывается целиком, когда задача её касается. На v0.19.4 записи-переходы — R01, R02, R06, R09, R18, R19, R23, R24, R27, R28, R30, R31, R32, R34, R38, R39, R41 и их поправки — называют старые теги и заготовки, а R39 говорит, что `.devcontainer/` — заготовка v0.18.14, хотя T90 взял заготовку v0.19.4. Отдельной задачей: это больше одного ревью (правило 11), T90 — переход.
- **Сделать:**
  1. При планировании — разбить по правилу 11 на задачи одного ревью каждая, по группам записей.
  2. Каждую запись-переход — переписать на месте целиком (R160 spec), по строгому правилу ссылок `spec/RULES.md`: что в ней ещё решение в силе — отличия `.devcontainer/` от заготовки — R10, R29 и R38 —, раскладка плана, пределы `docs-budgets.json` — остаётся в той записи, что его держит, без номера тега и истории; что было только переходом — уходит, и записи, которые на неё ссылаются, — тоже переписать.
  3. Д10 и Д11 закрыла T103: датированных поправок в журнале не осталось, R31 называет пределы по `docs-budgets.json`; закрыть строки, которые задача закрывает.
  4. R36 и R41 переписывает T83: здесь — только если T83 ещё не сделана, и тогда без её содержания.
- **Как проверить:** `just -f spec/justfile docs-check-repo`; `grep -n 'stands on .spec. v' docs/decisions.md` — пусто; `/code-review high`, один круг.
- **Разделы:** R160 spec; RULES.md spec — The bar, Across repositories; R10, R29, R31, R34, R38, R39, R41.
