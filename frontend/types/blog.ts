export interface IAuthor {
  name: string;
  avatar: string;
}

export interface IPost {
  id: string;
  title: string;
  excerpt: string;
  image: string;
  date: string;
  readTime: string;
  author: IAuthor;
}

export type SortOptionType = "newest" | "oldest" | "popular";
