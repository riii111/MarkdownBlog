import { defineStore } from "pinia";
import type { IPost } from "~/types/blog";

export const useBlogStore = defineStore("blog", () => {
  const _state = reactive({
    posts: [
      {
        id: "1",
        title: "Getting Started with Nuxt 3",
        excerpt: "Learn how to build modern web applications with Nuxt 3",
        image: "https://picsum.photos/800/600?random=1",
        date: "2024-03-20",
        readTime: "5 min read",
        author: {
          name: "John Doe",
          avatar: "https://i.pravatar.cc/150?img=1",
        },
      },
      {
        id: "2",
        title: "Vue 3 Composition API Best Practices",
        excerpt: "Discover the best practices for using Vue 3 Composition API",
        image: "https://picsum.photos/800/600?random=2",
        date: "2024-03-19",
        readTime: "7 min read",
        author: {
          name: "Jane Smith",
          avatar: "https://i.pravatar.cc/150?img=2",
        },
      },
      {
        id: "3",
        title: "TypeScript Tips and Tricks",
        excerpt: "Advanced TypeScript techniques for better type safety",
        image: "https://picsum.photos/800/600?random=3",
        date: "2024-03-18",
        readTime: "6 min read",
        author: {
          name: "Mike Johnson",
          avatar: "https://i.pravatar.cc/150?img=3",
        },
      },
    ] as IPost[],
    currentPage: 1,
    totalPages: 1,
  });

  const posts = computed(() => _state.posts);

  async function fetchPosts() {
    // TODO: APIが実装されるまでは何もしない
    return Promise.resolve();
  }

  return {
    posts,
    currentPage: computed(() => _state.currentPage),
    totalPages: computed(() => _state.totalPages),
    fetchPosts,
  };
});
