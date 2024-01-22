// const { defineConfig } = require('@vue/cli-service')
// module.exports = defineConfig({
//   transpileDependencies: true
// })

// module.exports = {
//   devServer: {
//     port: 8081
    
//   }
// };

const { defineConfig } = require('@vue/cli-service');

module.exports = defineConfig({
  transpileDependencies: true,
  devServer: {
    port: 8081,
    allowedHosts: 'all'  // Разрешает запросы со всех хостов
  }
});




