// when running locally, allow the frontend to proxy results to the golang API
// No CORS so u can live test basically

module.exports = {
  devServer: {
    proxy: {
      '^/': {
        target: 'http://localhost:8000',
        ws: false,
        changeOrigin: true
      }
    }
  }
}
