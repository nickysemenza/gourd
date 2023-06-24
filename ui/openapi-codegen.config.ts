import {
  generateSchemaTypes,
  generateReactQueryComponents,
} from "@openapi-codegen/typescript";
import { defineConfig } from "@openapi-codegen/cli";
export default defineConfig({
  gourdApi: {
    from: {
      relativePath: "../internal/api/openapi.yaml",
      source: "file",
    },
    outputDir: "src/api/react-query",
    to: async (context) => {
      const filenamePrefix = "gourdApi";
      const { schemasFiles } = await generateSchemaTypes(context, {
        filenamePrefix,
      });
      await generateReactQueryComponents(context, {
        filenamePrefix,
        schemasFiles,
      });
    },
  },
});
