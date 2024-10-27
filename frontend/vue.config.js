const { defineConfig } = require('@vue/cli-service');

module.exports = defineConfig({
  chainWebpack: (config) => {
    config.module
      .rule('glb')
      .test(/\.glb$/) // Wyszukuje pliki z rozszerzeniem .glb
      .use('file-loader')
      .loader('file-loader')
      .options({
        name: 'assets/mcieky_mouse_2.glb.', // Użyj oryginalnej nazwy pliku i hasha dla unikalności
        outputPath: 'assets/', // Folder docelowy dla plików GLB
        publicPath: 'assets/', // Ścieżka publiczna do plików GLB
      });
  },
});
