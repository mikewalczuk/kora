import { defineConfig } from "orval";

export default defineConfig({
  kora: {
    input: "../openapi.bundled.yaml",
    output: {
      mode: "tags-split",
      target: "src/lib/api/generated",
      client: "react-query",
      tsconfig: "./tsconfig.app.json",
      override: {
        mutator: {
          path: "src/lib/api/mutator.ts",
          name: "customFetch",
        },
      },
    },
  },
});
