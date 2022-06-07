import { defineConfig } from "cypress";

export default defineConfig({
  viewportWidth: 1600,
  viewportHeight: 660,
  e2e: {
    baseUrl: "http://localhost:3000",
    setupNodeEvents(on, config) {
      // bind to the event we care about
      // `on` is used to hook into various events Cypress emits
      // `config` is the resolved Cypress config

      require("@cypress/code-coverage/task")(on, config);
      const options = {
        printLogsToFile: "always",
        printLogsToConsole: "onFail",
        outputRoot: config.projectRoot + "/cypress/logs/",
        outputTarget: {
          "out.txt": "txt",
          "out.json": "json",
        },
      };
      require("cypress-terminal-report/src/installLogsPrinter")(on, options);

      // IMPORTANT to return the config object
      // with the any changed environment variables
      return config;
    },
  },
  retries: {
    runMode: 2,
    openMode: 0,
  },
  chromeWebSecurity: false,
  projectId: "hf26uf",
});
