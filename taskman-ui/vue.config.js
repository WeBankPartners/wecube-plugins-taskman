const CompressionPlugin = require('compression-webpack-plugin')
const packageName = require('./package.json').name
const postcssWrap = require('postcss-wrap')
const path = require('path')
const baseUrl = process.env.BASE_URL
module.exports = {
  devServer: {
    headers: {
      'Access-Control-Allow-Origin': '*'
    },
    port: 3010,
    proxy: {
      '/taskman': {
        target: baseUrl
      },
      '/wecmdb': {
        target: baseUrl
      },
      '/': {
        target: baseUrl
      }
    }
  },
  runtimeCompiler: true,
  // publicPath: '/taskman/',
  chainWebpack: config => {
    if (process.env.PLUGIN !== 'plugin') {
      // remove the old loader
      const img = config.module.rule('images')
      img.uses.clear()
      // add the new one
      img
        .use('file-loader')
        .loader('file-loader')
        .options({
          outputPath: 'img'
        })
    }

    config.when(process.env.PLUGIN === 'plugin', config => {
      config
        .entry('app')
        .clear()
        .add('./src/main-plugin.js') // 作为插件时
    })
    config.when(!process.env.PLUGIN, config => {
      config
        .entry('app')
        .clear()
        .add('./src/main.js') // 独立运行时
    })
  },
  productionSourceMap: process.env.PLUGIN !== 'plugin',
  // configureWebpack: config => {
  //   if (process.env.PLUGIN === 'plugin') {
  //     config.optimization.splitChunks = {}
  //     return
  //   }
  //   if (process.env.NODE_ENV === 'production') {
  //     return {
  //       plugins: [
  //         new CompressionPlugin({
  //           algorithm: 'gzip',
  //           test: /\.js$|\.html$|.\css/, // 匹配文件名
  //           threshold: 10240, // 对超过10k的数据压缩
  //           deleteOriginalAssets: false // 不删除源文件
  //         })
  //       ]
  //     }
  //   }
  // },
  configureWebpack: {
    output: {
      library: `${packageName}-[name]`,
      libraryTarget: 'umd', // 把微应用打包成 umd 库格式
      jsonpFunction: `webpackJsonp_${packageName}`
    },
    plugins: [
      new CompressionPlugin({
        algorithm: 'gzip',
        test: /\.js$|\.html$|.\css/, // 匹配文件名
        threshold: 10240, // 对超过10k的数据压缩
        deleteOriginalAssets: false // 不删除源文件
      })
    ]
  },
  transpileDependencies: ['detect-indent', 'redent', 'strip-indent', 'indent-string', 'crypto-random-string'],
  pluginOptions: {
    'style-resources-loader': {
      preProcessor: 'less',
      patterns: [path.resolve(__dirname, './src/assets/css/common.less')] // 引入全局样式变量
    }
  },
  css: {
    loaderOptions: {
      postcss: {
        plugins: process.env.PLUGIN === 'plugin' ? [
          postcssWrap({
            // selector: '.taskman-wrap'
          })
        ] : []
      }
    }
  }
}
