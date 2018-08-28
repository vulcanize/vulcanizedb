module.exports = {
  entry: "./src/index.ts",
  mode: "production",
  devtool: "source-map",
  target: 'node',
  output: {
    filename: "vulcanize-postgraphile-server.js",
    path: __dirname + "/build/dist/",
    publicPath: "build/dist/"
  },
  resolve: {
    extensions: [".ts", ".tsx", ".mjs", ".js", ".json", ".css", ".png"]
  },
  module: {
    rules: [
      { test: /\.ts$/, loader: "awesome-typescript-loader" },
      { enforce: "pre", test: /\.js$/, loader: "source-map-loader" }
    ]
  }
};
