const CompressionPlugin = require('compression-webpack-plugin')
const postcssWrap = require('postcss-wrap')
const path = require('path')
const baseUrl = process.env.VUE_APP_API || 'http://127.0.0.1/'
module.exports = {
  devServer: {
    // hot: true,
    // inline: true,
    open: true,
    port: 3000,
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
  publicPath: '/',
  // publicPath: '/',
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
  configureWebpack: config => {
    if (process.env.PLUGIN === 'plugin') {
      config.optimization.splitChunks = {}
      // config.optimization = {
      //   runtimeChunk: 'single',
      //   splitChunks: {
      //     chunks: 'all',
      //     minSize: 20000, // 允许新拆出 chunk 的最小体积
      //     maxSize: 500000, // 设置chunk的最大体积为500KB
      //     automaticNameDelimiter: '-',
      //     cacheGroups: {
      //       defaultVendors: {
      //         test: /[\\/]node_modules[\\/]/,
      //         priority: -10
      //       },
      //       default: {
      //         minChunks: 2,
      //         priority: -20,
      //         reuseExistingChunk: true
      //       }
      //     }
      //   }
      // }
      return
    }
    if (process.env.NODE_ENV === 'production') {
      return {
        plugins: [
          new CompressionPlugin({
            algorithm: 'gzip',
            test: /\.js$|\.html$|.\css/, // 匹配文件名
            threshold: 10240, // 对超过10k的数据压缩
            deleteOriginalAssets: false // 不删除源文件
          })
        ]
      }
    }
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
      },
      less: {
        lessOptions: { // 将 javascriptEnabled 放在 lessOptions 内部
          javascriptEnabled: true
        }
      }
    }
  }

}
