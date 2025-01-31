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
          // ダークモード
          "dark:bg-gray-800 dark:text-white dark:border-gray-700 dark:placeholder-gray-500",
          // フォーカス時
          "focus:ring-2 focus:ring-emerald-500 dark:focus:ring-emerald-400",
          // オートコンプリート時
          "[&:-webkit-autofill]:dark:bg-gray-700 [&:-webkit-autofill]:dark:text-gray-100",
          "[&:-webkit-autofill]:dark:!shadow-[inset_0_0_0_1000px_rgb(55,65,81)]",
        ].join(" "),
      },
      // エラー時のスタイル
      error: {
        class: "bg-red-100 dark:bg-red-900/50 dark:text-red-100",
      },
      // オートコンプリート時のスタイル
      // TODO: オートコンプリート時のスタイルを調整する
      // autocomplete: {
      //   class: 'bg-gray-200 dark:bg-gray-700 dark:text-black',
      // },
      // 特定のケース（例：検索フィールド）用のバリエーション
      search: {
        variant: "outline",
        class: "bg-transparent text-white placeholder-white/70 border-white/30",
      },
    },
    modal: {
      rounded: "rounded-xl",
      background: "bg-white dark:bg-gray-800",
    },
    button: {
      default: {
        rounded: "rounded-md",
      },
    },
  },
});
