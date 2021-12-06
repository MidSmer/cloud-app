const path = require('path');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const {CleanWebpackPlugin} = require('clean-webpack-plugin');
const {VueLoaderPlugin} = require('vue-loader/dist/index');
const CopyWebpackPlugin = require('copy-webpack-plugin');

module.exports = {
    entry: {
        'index': './src/index.js'
    },
    output: {
        filename: 'static/[contenthash].js',
        path: path.resolve(__dirname, '..', 'deploy', 'public'),
        publicPath: '/',
    },
    plugins: [
        new MiniCssExtractPlugin({
            filename: ({ chunk }) => {
                if (chunk.name === 'tinymce_content') {
                    return 'static/tinymce_content.css'
                } else {
                    return 'static/[contenthash].css'
                }
            }
        }),
        new HtmlWebpackPlugin({
            template: './public/index.html',
            filename: 'index.html',
            title: 'web'
        }),
        new CleanWebpackPlugin(),
        new VueLoaderPlugin()
    ],
    module: {
        rules: [
            {
                test: /\.tsx?$/,
                exclude: /node_modules/,
                use: ['babel-loader']
            },
            {
                test: /\.vue$/,
                use: [
                    'vue-loader'
                ],
            },
            {
                test: /\.css$/,
                use: [
                    MiniCssExtractPlugin.loader,
                    'css-loader'
                ]
            },
            {
                test: /\.less$/,
                use: [
                    'style-loader',
                    'css-loader',
                    'less-loader'
                ]
            },
            {
                test: /\.(jpg|png|jpeg|gif|bmp)$/,
                use: {
                    loader: 'url-loader',
                    options: {
                        limit: 1024,
                        fallback: {
                            loader: 'file-loader',
                            options: {
                                name: '[name].[ext]'
                            }
                        }
                    }
                }
            },
            {
                test: /\.(mp4|ogg|mp3|wav)$/,
                use: {
                    loader: 'url-loader',
                    options: {
                        limit: 1024,
                        fallback: {
                            loader: 'file-loader',
                            options: {
                                name: '[name].[ext]'
                            }
                        }
                    }
                }
            }
        ]
    },
    optimization: {
        splitChunks: {
            chunks: 'all',
            cacheGroups: {
                tinymceVendor: {
                    test: /[\\/]node_modules[\\/](tinymce)[\\/](.*js)|[\\/]plugins[\\/]/,
                    name: 'tinymce'
                },
                tinymceContent: {
                    test: /[\\/]node_modules[\\/](tinymce)[\\/](.*content.css)/,
                    name: 'tinymce_content'
                },
            },
        }
    },
}
