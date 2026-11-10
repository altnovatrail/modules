// The test of the module's templates against the tasks they came from (R22): module-kit builds the
// module, as the release does, and for each template the instance with the numbers of its reference
// task, which a line of its header gives, has that task's key, its band, its difficulty and its
// leading question. So the solver of a template still proves the answer of the task that
// standalone's template came from, as standalone's own check of its templates does. It runs the
// real module-kit, so it runs in the job Modules.

import assert from "node:assert/strict";
import { readFileSync } from "node:fs";
import { join } from "node:path";
import { test } from "node:test";

import { ROOT, mustRun, temporaryFolder } from "./helpers.mjs";

const MODULE = join(ROOT, "modules", "olympiad-math");

// The line of a template's header that names its reference task and that task's numbers, as in
// "# The numbers of reference task gaps-34-d4-4: flight 18, told_to 3, asked_from 3, asked_to 7."
const REFERENCE = /^# The numbers of reference task (\S+): (.+)\.$/m;

// The key of a task's first check of a number or a fraction, as the package holds it.
function keyOf(task) {
  for (const item of task.items) {
    for (const slot of item.answer_slots) {
      if (slot.check.kind === "number" || slot.check.kind === "fraction") return slot.check.key;
    }
  }
  return undefined;
}

test("the instance of each template with its reference task's numbers has that task's key", (t) => {
  const out = join(temporaryFolder(t), "package");
  mustRun("module-kit", ["build", MODULE, out]);
  const server = JSON.parse(readFileSync(join(out, "server.json"), "utf8"));
  // The templates are the tasks whose instances the package carries, however their files write them.
  const templates = Object.keys(server.instances);
  assert.ok(templates.length > 0, "the module has no template");
  for (const id of templates) {
    const line = REFERENCE.exec(readFileSync(join(MODULE, "tasks", `${id}.yaml`), "utf8"));
    assert.ok(line, `${id}: its header names no numbers of a reference task`);
    const [, reference, list] = line;
    const numbers = list.split(", ").map((pair) => pair.split(" "));
    const instance = server.instances[id].find((each) => numbers.every(([name, value]) => each[name] === Number(value)));
    assert.ok(instance, `${id}: no instance has the numbers of ${reference}`);
    const key = keyOf(server.tasks[id]).replace(/\{([a-z][a-z0-9_]*)\}/g, (_, name) => String(instance[name]));
    assert.equal(key, keyOf(server.tasks[reference]), `${id}: the key of the instance of ${reference}`);
    for (const field of ["band", "difficulty"]) {
      assert.equal(server.tasks[id][field], server.tasks[reference][field], `${id}: its ${field} and that of ${reference}`);
    }
    // A template whose reference task has a leading question has it too, with placeholders for
    // the reference task's numbers (amendment T51 to R22).
    const question = server.tasks[id].leading_question?.replace(/\{([a-z][a-z0-9_]*)\}/g, (_, name) => String(instance[name]));
    assert.equal(question, server.tasks[reference].leading_question, `${id}: the leading question of the instance of ${reference}`);
  }
});
