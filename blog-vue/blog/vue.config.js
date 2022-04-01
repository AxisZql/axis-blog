module.exports = {
    // assetsDir: 'static',
    // parallel: false,
    // publicPath: './',
    transpileDependencies: ["vuetify"],
    devServer: {
        proxy: {
            "/api": {
                target: "http://localhost:8080",
                changeOrigin: true,
                pathRewrite: {
                    "^/api": ""
                }
            }
        },
        disableHostCheck: true
    },
    productionSourceMap: false,
    css: {
        extract: true,
        sourceMap: false
    }
};