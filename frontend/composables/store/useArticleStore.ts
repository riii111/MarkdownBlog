import { defineStore } from "pinia";
import { ARTICLE_CONSTANTS } from "~/constants/article";

export const useArticleStore = defineStore("article", () => {
  const _state = reactive({
    articles: [] as readonly IArticle[],
    currentPage: ARTICLE_CONSTANTS.DEFAULT_PAGE as number,
    itemsPerPage: ARTICLE_CONSTANTS.ITEMS_PER_PAGE,
    isLoading: false,
  });

  const allArticles = computed(() => _state.articles);

  const articles = computed(() => {
    const start = (_state.currentPage - 1) * _state.itemsPerPage;
    const end = start + _state.itemsPerPage;
    return _state.articles.slice(start, end);
  });

  const currentPage = computed({
    get: () => _state.currentPage,
    set: (value: number) => {
      _state.currentPage = value;
    },
  });

  const totalPages = computed(() =>
    Math.ceil(_state.articles.length / _state.itemsPerPage)
  );

  const totalItems = computed(() => _state.articles.length);

  const isLoading = computed(() => _state.isLoading);

  function setPage(page: number) {
    if (page > 0 && page <= totalPages.value) {
      _state.currentPage = page;
    }
  }

  async function fetchArticles() {
    _state.isLoading = true;
    try {
      await new Promise((resolve) => setTimeout(resolve, 1500));
      _state.articles = [
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
          content: `
# Getting Started with Nuxt 3

Nuxt 3 represents a complete rewrite of the framework, offering developers a more powerful and flexible way to build modern web applications.

## Key Features

### Server-Side Rendering (SSR)

Nuxt 3 comes with built-in server-side rendering capabilities, offering:
- Better SEO performance
- Faster initial page loads
- Improved core web vitals

### Auto Imports

One of the most convenient features in Nuxt 3 is the auto-import system:
- Components
- Composables
- Vue APIs

\`\`\`typescript
// No imports needed!
const { $fetch } = useNuxtApp()
const route = useRoute()
\`\`\`

### File-based Routing

Nuxt 3 uses a file-based routing system:

\`\`\`
pages/
  index.vue         # /
  about.vue        # /about
  posts/
    [id].vue       # /posts/:id
    index.vue      # /posts
\`\`\`
          `,
        },
        {
          id: "2",
          title: "Vue 3 Composition API Best Practices",
          excerpt:
            "Discover the best practices for using Vue 3 Composition API",
          image: "https://picsum.photos/800/600?random=2",
          date: "2024-03-19",
          readTime: "7 min read",
          author: {
            name: "Jane Smith",
            avatar: "https://i.pravatar.cc/150?img=2",
          },
          content: `
# Vue 3 Composition API Best Practices

The Composition API in Vue 3 provides a more flexible and powerful way to manage component logic compared to the Options API.

## Key Features

### Reactive State Management

The Composition API allows you to manage reactive state using reactive variables and functions.

\`\`\`typescript
const { count, increment } = useCounter()
\`\`\`

### Lifecycle Hooks

You can use lifecycle hooks to run code at specific points in a component's lifecycle.

\`\`\`typescript
onMounted(() => {
  // Code to run when the component is mounted
})
\`\`\`

### Custom Composables

You can create custom composables to encapsulate reusable logic.

\`\`\`typescript
function useCounter() {
  const count = ref(0)
  const increment = () => count.value++
  return { count, increment }
}
\`\`\`
          `,
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
          content: `
# TypeScript Tips and Tricks

TypeScript is a statically typed superset of JavaScript that compiles to plain JavaScript.

## Key Features

### Type Inference

TypeScript can infer types based on the code, reducing the need for explicit type annotations.

\`\`\`typescript
const x = 10; // x is inferred to be of type number
\`\`\`

### Interfaces and Types

You can define interfaces and types to describe the shape of objects.

\`\`\`typescript
interface Person {
  name: string;
  age: number;
}

const person: Person = { name: "John", age: 30 };
\`\`\`

### Generics

Generics allow you to create reusable components and functions that can work with multiple types.

\`\`\`typescript
function identity<T>(arg: T): T {
  return arg;
}

const result = identity<number>(10);
\`\`\`
          `,
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
          content: `
# Modern CSS Techniques

CSS is a powerful language used for describing the presentation of a document written in HTML or XML.

## Key Features

### Flexbox

Flexbox is a layout model that allows you to align items in a container in a flexible way.

\`\`\`css
.container {
  display: flex;
}
\`\`\`

### Grid

Grid is a two-dimensional layout system that allows you to align items in rows and columns.

\`\`\`css
.container {
  display: grid;
}
\`\`\`

### CSS Variables

CSS variables allow you to store reusable values in your stylesheets.

\`\`\`css
:root {
  --primary-color: #3498db;
}

.button {
  background-color: var(--primary-color);
}
\`\`\`
          `,
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
          content: `
# Testing Vue Applications

Testing is an essential part of software development. It helps ensure that your application works as expected and provides a reliable user experience.

## Key Features

### Unit Testing

Unit testing involves testing individual components or functions in isolation.

\`\`\`typescript
// Example of a unit test
test('increment function', () => {
  expect(increment(0)).toBe(1);
  expect(increment(1)).toBe(2);
});
\`\`\`

### Integration Testing

Integration testing involves testing multiple components or services together.

\`\`\`typescript
// Example of an integration test
test('increment function', async () => {
  const { increment } = useCounter()
  await increment()
  expect(count.value).toBe(1)
})
\`\`\`

### End-to-End Testing

End-to-end testing involves testing the entire application from start to finish.

\`\`\`typescript
// Example of an end-to-end test
test('increment function', async () => {
  const { increment } = useCounter()
  await increment()
  expect(count.value).toBe(1)
})
\`\`\`
          `,
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
          content: `
# State Management with Pinia

Pinia is a state management library for Vue 3 that provides a more flexible and powerful way to manage state compared to Vuex.

## Key Features

### Reactive State Management

Pinia allows you to manage reactive state using reactive variables and functions.

\`\`\`typescript
const { count, increment } = useCounter()
\`\`\`

### Lifecycle Hooks

You can use lifecycle hooks to run code at specific points in a component's lifecycle.

\`\`\`typescript
onMounted(() => {
  // Code to run when the component is mounted
})
\`\`\`

### Custom Composables

You can create custom composables to encapsulate reusable logic.

\`\`\`typescript
function useCounter() {
  const count = ref(0)
  const increment = () => count.value++
  return { count, increment }
}
\`\`\`
          `,
        },
        {
          id: "7",
          title: "Building Responsive Layouts with Tailwind CSS",
          excerpt:
            "Master responsive design using Tailwind CSS utility classes",
          image: "https://picsum.photos/800/600?random=7",
          date: "2024-03-14",
          readTime: "7 min read",
          author: {
            name: "David Lee",
            avatar: "https://i.pravatar.cc/150?img=7",
          },
          content: `
# Building Responsive Layouts with Tailwind CSS

Tailwind CSS is a utility-first CSS framework that allows you to build responsive layouts quickly and efficiently.

## Key Features

### Utility Classes

Tailwind CSS provides a wide range of utility classes that you can use to style your components.

\`\`\`css
.container {
  @apply bg-gray-100 p-4 rounded-lg;
}
\`\`\`

### Responsive Design

Tailwind CSS allows you to create responsive designs that work across different screen sizes.

\`\`\`css
.container {
  @apply w-full md:w-1/2 lg:w-1/3;
}
\`\`\`

### Customization

You can customize Tailwind CSS to fit your project's needs by adding your own custom styles.

\`\`\`css
.container {
  @apply bg-gray-100 p-4 rounded-lg;
}
\`\`\`
          `,
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
          content: `
# API Integration in Nuxt Applications

Integrating APIs into your Nuxt.js applications is a common requirement.

## Key Features

### REST API

Nuxt.js supports REST APIs out of the box. You can use the \`@nuxtjs/axios\` module to make HTTP requests.

\`\`\`typescript
const { data } = await useFetch('/api/posts')
\`\`\`

### GraphQL

Nuxt.js also supports GraphQL. You can use the \`@nuxtjs/graphql\` module to make GraphQL requests.

\`\`\`typescript
const { data } = await useFetch('/api/graphql', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    query: \`
      query {
        posts {
          id
          title
          excerpt
        }
      }
    \`
  })
})
\`\`\`
          `,
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
          content: `
# Optimizing Vue Application Performance

Performance optimization is a crucial aspect of building high-quality applications.

## Key Features

### Lazy Loading

Lazy loading allows you to load components and resources only when they are needed.

\`\`\`typescript
const { data } = await useFetch('/api/posts')
\`\`\`

### Code Splitting

Code splitting allows you to split your code into smaller chunks, which can improve load times.

\`\`\`typescript
const { data } = await useFetch('/api/posts')
\`\`\`

### Caching

Caching allows you to store frequently accessed data in memory, which can improve load times.

\`\`\`typescript
const { data } = await useFetch('/api/posts')
\`\`\`
          `,
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
          content: `
# Authentication in Nuxt Applications

Implementing secure authentication in your Nuxt.js applications is a crucial requirement.

## Key Features

### JWT Authentication

JSON Web Tokens (JWT) are a popular method for implementing authentication.

\`\`\`typescript
const { data } = await useFetch('/api/login', {
  method: 'POST',
  body: {
    email: 'john@example.com',
    password: 'password123'
  },
})
\`\`\`

### OAuth Authentication

OAuth is a popular method for implementing authentication.

\`\`\`typescript
const { data } = await useFetch('/api/login', {
  method: 'POST',
  body: {
    email: 'john@example.com',
    password: 'password123'
  },
})
\`\`\`
          `,
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
          content: `
# Deploying Nuxt Applications

Deploying Nuxt applications to various platforms is a crucial requirement.

## Key Features

### Server Deployment

You can deploy your Nuxt application to various server platforms such as Vercel, Netlify, or any custom server.

\`\`\`
pages/
  index.vue         # /
  about.vue        # /about
  posts/
    [id].vue       # /posts/:id
    index.vue      # /posts
\`\`\`

### Container Deployment

You can deploy your Nuxt application as a container using Docker.

\`\`\`
docker run -p 3000:3000 nuxt-app
\`\`\`
          `,
        },
      ] as const;
    } finally {
      _state.isLoading = false;
    }
  }

  return {
    articles,
    allArticles,
    currentPage,
    totalPages,
    totalItems,
    isLoading,
    setPage,
    fetchArticles,
  };
});
