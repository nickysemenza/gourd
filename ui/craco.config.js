// https://dev.to/nicolasrannou/wasm-in-create-react-app-4-in-5mn-without-ejecting-cf6
module.exports = {
  webpack: {
    configure: (webpackConfig) => {
      // https://github.com/rustwasm/wasm-pack/issues/835#issuecomment-772591665
      webpackConfig.module.rules.push({
        test: /\.wasm$/,
        type: "webassembly/sync",
      });
      if (webpackConfig.experiments === undefined) {
        webpackConfig.experiments = {};
      }
      webpackConfig.experiments["syncWebAssembly"] = true;
      // https://github.com/facebook/create-react-app/pull/11752
      webpackConfig.ignoreWarnings = [/Failed to parse source map/];
      return webpackConfig;
    },
  },
};
