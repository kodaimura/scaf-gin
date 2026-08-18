import assert from "node:assert/strict";
import test from "node:test";

import { resetTokenFromMessage } from "./support.mjs";

test("plain reset tokens beginning with hexadecimal characters remain intact", () => {
  const token = "d6yE9tdwZvEtPC5SpzkifVJJ2gD8uWAM1_9eG7-eYk";
  assert.equal(
    resetTokenFromMessage({
      Content: {
        Headers: {},
        Body: `http://localhost:3000/reset-password?token=${token}`,
      },
    }),
    token,
  );
});

test("quoted-printable reset links are decoded before extraction", () => {
  assert.equal(
    resetTokenFromMessage({
      Content: {
        Headers: { "Content-Transfer-Encoding": ["quoted-printable"] },
        Body: "http://localhost:3000/reset-password?token=3Dabc_123-xyz",
      },
    }),
    "abc_123-xyz",
  );
});
