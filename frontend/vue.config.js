module.exports = {
  devServer: {
    proxy: {
      "^/": {
        target: "http://localhost:8000",
        ws: false,
        changeOrigin: true
      }
    }
  },
  transpileDependencies: ["vuetify"]
};
