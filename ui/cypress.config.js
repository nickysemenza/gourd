import { defineConfig } from "cypress";
import vitePreprocessor from "cypress-vite";

export default defineConfig({
  e2e: {
    baseUrl: "http://localhost:5173/",

    supportFile: false,

    viewportWidth: 1600,
    viewportHeight: 660,
    setupNodeEvents(on) {
      on("file:preprocessor", vitePreprocessor());
    },
  },
  retries: {
    runMode: 2,
    openMode: 0,
  },
  chromeWebSecurity: false,
  projectId: "hf26uf",
});
