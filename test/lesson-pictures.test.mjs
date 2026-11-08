// The size of each picture of a lesson against its file (amendment T41 to R12): the port draws a
// picture of media/lessons/ at a scale of 2, as it draws a task's, so a part of a lesson that shows
// it says half the size of its PNG in CSS px. The port writes the size of a task's picture itself;
// a lesson is written by hand, and this test keeps its sizes true when a picture is drawn again.

import assert from "node:assert/strict";
import { readdirSync, readFileSync } from "node:fs";
import { join } from "node:path";
import { test } from "node:test";

import { ROOT } from "./helpers.mjs";

const MODULE = join(ROOT, "modules", "olympiad-math");

// A picture part of a lesson as the lessons write it: its path, its alt, its width and its height,
// one line each.
const PART = /^ +- picture: "(lessons\/[^"]+)"\n +alt: .*\n +width: (\d+)\n +height: (\d+)$/gm;

test("every picture of a lesson is shown at half the size of its PNG", () => {
  const lessons = readdirSync(join(MODULE, "lessons")).filter((name) => name.endsWith(".yaml"));
  let shown = 0;
  for (const name of lessons) {
    const text = readFileSync(join(MODULE, "lessons", name), "utf8");
    const parts = [...text.matchAll(PART)];
    assert.equal(parts.length, (text.match(/picture: "lessons\//g) ?? []).length, `${name}: a picture part in another layout`);
    for (const [, path, width, height] of parts) {
      const png = readFileSync(join(MODULE, "media", path));
      assert.deepEqual([Number(width), Number(height)], [png.readUInt32BE(16) / 2, png.readUInt32BE(20) / 2], `${name}: ${path}`);
      shown++;
    }
  }
  assert.ok(shown > 0, "no lesson shows a picture of media/lessons/");
});
