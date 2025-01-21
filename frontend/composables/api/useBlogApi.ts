import type { IPost } from "~/types/blog";
import { useCoreApi } from "./useCoreApi";

export const useBlogApi = () => {
  const fetchPosts = async (page = 1) => {
    return await useCoreApi<IPost[]>("/posts", {
      method: "GET",
      params: { page },
    });
  };

  return {
    fetchPosts,
  };
};
