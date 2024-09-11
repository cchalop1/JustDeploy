import type { CodegenConfig } from "@graphql-codegen/cli";

const config: CodegenConfig = {
  overwrite: true,
  schema: "http://localhost:8080/graphql",
  documents: "src/**/*.tsx",
  ignoreNoDocuments: true,
  generates: {
    "./src/graphql/": {
      preset: "client",
      config: {
        documentMode: "string",
      },
    },
    "./schema.graphql": {
      plugins: ["schema-ast"],
      config: {
        includeDirectives: true,
      },
    },
  },
};

export default config;
