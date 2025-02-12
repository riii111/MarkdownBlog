export const CATEGORIES = [
  "All",
  "Development",
  "Design",
  "DevOps",
  "Tips & Tricks",
] as const;
export const SORT_OPTIONS = [
  { label: "Newest", value: "newest" },
  { label: "Oldest", value: "oldest" },
  { label: "Most Popular", value: "popular" },
] as const;

export const ARTICLE_CONSTANTS = {
  ITEMS_PER_PAGE: 9,
  DEFAULT_PAGE: 1,
} as const;

export const ARTICLE_MESSAGES = {
  DESCRIPTION:
    "Here, we share development tips, technical guides, and stories that inspire your next project.",
} as const;
