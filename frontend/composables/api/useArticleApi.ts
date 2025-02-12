export const useArticleApi = () => {
  const fetchArticles = async (page = 1) => {
    return await useCoreApi<IArticle[]>("/articles", {
      method: "GET",
      params: { page },
    });
  };

  return {
    fetchArticles,
  };
};
