export default defineNuxtConfig({
  srcDir: "./",
  components: {
    dirs: [
      {
        path: "~/components",
        pathPrefix: false,
      },
    ],
  },
  // ネストしたcomposableも自動でインポートするように設定
  // https://nuxt.com/docs/guide/directory-structure/composables#how-files-are-scanned
  imports: {
    dirs: ["composables/**"],
  },
  modules: [
    "@nuxt/ui",
    [
      "@pinia/nuxt",
      {
        autoImports: ["defineStore", ["defineStore", "definePiniaStore"]],
      },
    ],
    "@nuxt/image",
    "@nuxtjs/color-mode",
  ],
  devtools: { enabled: true },
  ui: {
    icons: ["lucide"],
  },
  typescript: {
    strict: true,
  },
  colorMode: {
    classSuffix: "",
  },

  app: {
    head: {
      title: "MD Writer",
      meta: [
        { charset: "utf-8" },
        { name: "viewport", content: "width=device-width, initial-scale=1" },
      ],
    },
  },

  compatibilityDate: "2025-01-17",
});
