import { defineStore } from "pinia";
import type { IPost } from "~/types/blog";
import { BLOG_CONSTANTS } from "~/constants/blog";

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
      {
        id: "4",
        title: "Modern CSS Techniques",
        excerpt: "Exploring the latest CSS features and best practices",
        image: "https://picsum.photos/800/600?random=4",
        date: "2024-03-17",
        readTime: "8 min read",
        author: {
          name: "Sarah Wilson",
          avatar: "https://i.pravatar.cc/150?img=4",
        },
      },
      {
        id: "5",
        title: "Testing Vue Applications",
        excerpt:
          "Comprehensive guide to testing Vue.js applications with Vitest",
        image: "https://picsum.photos/800/600?random=5",
        date: "2024-03-16",
        readTime: "10 min read",
        author: {
          name: "Alex Brown",
          avatar: "https://i.pravatar.cc/150?img=5",
        },
      },
      {
        id: "6",
        title: "State Management with Pinia",
        excerpt:
          "Learn how to effectively manage state in Vue applications using Pinia",
        image: "https://picsum.photos/800/600?random=6",
        date: "2024-03-15",
        readTime: "9 min read",
        author: {
          name: "Emily Chen",
          avatar: "https://i.pravatar.cc/150?img=6",
        },
      },
      {
        id: "7",
        title: "Building Responsive Layouts with Tailwind CSS",
        excerpt: "Master responsive design using Tailwind CSS utility classes",
        image: "https://picsum.photos/800/600?random=7",
        date: "2024-03-14",
        readTime: "7 min read",
        author: {
          name: "David Lee",
          avatar: "https://i.pravatar.cc/150?img=7",
        },
      },
      {
        id: "8",
        title: "API Integration in Nuxt Applications",
        excerpt:
          "Best practices for integrating APIs in your Nuxt.js applications",
        image: "https://picsum.photos/800/600?random=8",
        date: "2024-03-13",
        readTime: "8 min read",
        author: {
          name: "Sophie Taylor",
          avatar: "https://i.pravatar.cc/150?img=8",
        },
      },
      {
        id: "9",
        title: "Optimizing Vue Application Performance",
        excerpt:
          "Tips and techniques for improving Vue.js application performance",
        image: "https://picsum.photos/800/600?random=9",
        date: "2024-03-12",
        readTime: "12 min read",
        author: {
          name: "Ryan Martinez",
          avatar: "https://i.pravatar.cc/150?img=9",
        },
      },
      {
        id: "10",
        title: "Authentication in Nuxt Applications",
        excerpt:
          "Implementing secure authentication in your Nuxt.js applications",
        image: "https://picsum.photos/800/600?random=10",
        date: "2024-03-11",
        readTime: "11 min read",
        author: {
          name: "Lisa Anderson",
          avatar: "https://i.pravatar.cc/150?img=10",
        },
      },
      {
        id: "11",
        title: "Deploying Nuxt Applications",
        excerpt:
          "Complete guide to deploying Nuxt applications to various platforms",
        image: "https://picsum.photos/800/600?random=11",
        date: "2024-03-10",
        readTime: "9 min read",
        author: {
          name: "Tom Wilson",
          avatar: "https://i.pravatar.cc/150?img=11",
        },
      },
    ] as IPost[],
    currentPage: BLOG_CONSTANTS.DEFAULT_PAGE as number,
    itemsPerPage: BLOG_CONSTANTS.ITEMS_PER_PAGE,
  });

  const allPosts = computed(() => _state.posts);

  const posts = computed(() => {
    const start = (_state.currentPage - 1) * _state.itemsPerPage;
    const end = start + _state.itemsPerPage;
    return _state.posts.slice(start, end);
  });

  const currentPage = computed({
    get: () => _state.currentPage,
    set: (value: number) => {
      _state.currentPage = value;
    },
  });

  const totalPages = computed(() =>
    Math.ceil(_state.posts.length / _state.itemsPerPage)
  );

  function setPage(page: number) {
    if (page > 0 && page <= totalPages.value) {
      _state.currentPage = page;
    }
  }

  async function fetchPosts() {
    // TODO: APIが実装されるまでは何もしない
    return Promise.resolve();
  }

  return {
    posts,
    allPosts,
    currentPage,
    totalPages,
    setPage,
    fetchPosts,
  };
});
