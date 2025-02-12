export default defineAppConfig({
  ui: {
    primary: "emerald",
    gray: "zinc",
    input: {
      default: {
        rounded: "rounded-md",
        variant: "bordered",
        class: [
          // ライトモード
          "bg-white text-gray-900 placeholder-gray-400 border-gray-300",
          "focus:ring-2 focus:ring-emerald-500",
        ].join(" "),
      },
      // エラー時のスタイル
      error: {
        class: "bg-red-100",
      },
      search: {
        variant: "outline",
        class: "bg-transparent text-white placeholder-white/70 border-white/30",
      },
    },
    modal: {
      rounded: "rounded-xl",
      background: "bg-white",
    },
    button: {
      default: {
        rounded: "rounded-md",
      },
    },
  },
});
