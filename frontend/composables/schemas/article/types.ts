export interface IAuthor {
  name: string;
  avatar: string;
}

export interface IArticle {
  id: string;
  title: string;
  excerpt: string;
  image: string;
  date: string;
  readTime: string;
  author: IAuthor;
  content?: string;
}

export type SortOptionType = "newest" | "oldest" | "popular";
