export default defineNuxtConfig({
  modules: ["@nuxt/ui"],
  devtools: { enabled: true },

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