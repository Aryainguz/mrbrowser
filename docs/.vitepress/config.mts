import { defineConfig } from 'vitepress'

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "Mr. Browser",
  description: "Intent-driven, self-healing browser automation",
  appearance: 'dark',
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [
      { text: 'Home', link: '/' },
      { text: 'Guide', link: '/guide/getting-started' },
      { text: 'API Reference', link: '/sdk/python' }
    ],

    sidebar: [
      {
        text: 'Introduction',
        items: [
          { text: 'Getting Started', link: '/guide/getting-started' },
          { text: 'Core Concepts', link: '/guide/core-concepts' }
        ]
      },
      {
        text: 'SDKs',
        items: [
          { text: 'Python SDK', link: '/sdk/python' },
          { text: 'TypeScript SDK', link: '/sdk/typescript' }
        ]
      }
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/mrbrowser/mrbrowser' }
    ],

    footer: {
      message: 'Released under the MIT License.',
      copyright: 'Copyright © 2026 Mr. Browser Contributors'
    }
  }
})
