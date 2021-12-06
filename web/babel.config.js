module.exports = {
    module: {
        rules: [
            {
                test: /\.js$/,
                use: {
                    loader: "babel-loader",
                    options: {
                        presets: ["@babel/preset-env"]
                    }
                }
            }
        ]
    },
    presets: [
        [
            "@babel/preset-env",
            {
                "corejs": "3",
                "useBuiltIns": "usage"
            }
        ],
        "@babel/preset-typescript"
    ]
}
