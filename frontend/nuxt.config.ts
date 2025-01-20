export default defineNuxtConfig({
  srcDir: "./",
  // ネストしたcomposable, componentsも自動でインポートするようにする
  // https://nuxt.com/docs/guide/directory-structure/composables#how-files-are-scanned
  imports: {
    dirs: ["composables/**", "components/**"],
  },
  modules: [
    "@nuxt/ui",
    [
      "@pinia/nuxt",
      {
        autoImports: ["defineStore", ["defineStore", "definePiniaStore"]],
      },
    ],
  ],
  devtools: { enabled: true },
  ui: {
    icons: ["lucide"],
  },
  typescript: {
    strict: true,
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
