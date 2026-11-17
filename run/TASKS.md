# TASKS.md — tasks outside an epic of `modules`

Open tasks outside the current epic, each as its body, in Russian (R154 spec, R155 spec).

- Never read whole: a session reads only its own task — `grep -n '^#### T' run/TASKS.md` lists them, `awk '/^#### T76\./{p=1;print;next} /^#### /{p=0} p' run/TASKS.md` gives one body.
- A new task outside an epic goes here as a body, its number the one free when it is written down (RUN.md, Citations). When done, the body leaves for the archive (RUN.md rule 9).

## Tasks outside an epic


#### T81. english-essays: новая версия без фразы про имена (Д152 spec)

- **Цель:** новая версия модуля `english-essays`, в пяти условиях задач которой нет фразы «Do not write your name or the real name of anyone you know.»; релиз `modules`.
- **Почему:** Д152 spec (ответ T95, R127 spec): фразу снять в новой версии модуля, как только `app` показывает её сам перед каждой отправкой ИИ работы целиком — это T47 `app`, выложен на `novatrail` с `app` v0.1.19 (T73 `infra`, 2026-11-17). R76 spec, R156 spec; параллельная работа волны 4 (слово автора 2026-11-17).
- **Сделать:** новая версия `english-essays` по правилам версий формата и `modules` (R44 spec): пять условий без фразы, остальное как есть; проверки модуля; релиз по правилам репозитория; строку Д152 — в отчёте как закрытую с версией (её удаляет spec при разборе).
- **Как проверить:** проверки `modules` зелёные; релиз вышел.
- **Разделы:** Д152 spec; R127 spec; R44 spec.
