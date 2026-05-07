import { defineConfig } from "orval";

export default defineConfig({
  kora: {
    input: "../api-spec/openapi.yaml",
    output: {
      mode: "tags-split",
      target: "src/lib/api/generated",
      client: "react-query",
      override: {
        mutator: {
          path: "src/lib/api/mutator.ts",
          name: "customFetch",
        },
      },
    },
  },
});
