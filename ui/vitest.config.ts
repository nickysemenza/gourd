import { defineConfig } from "vitest/config";
import react from "@vitejs/plugin-react-swc";
import wasm from "vite-plugin-wasm";
import topLevelAwait from "vite-plugin-top-level-await";

export default defineConfig({
  plugins: [react(), wasm(), topLevelAwait()],
  test: {
    coverage: {
      reporter: ["text", "json-summary", "json", "html"],
    },
  },
});
