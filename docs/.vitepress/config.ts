import { mkdirSync, writeFileSync } from 'node:fs'
import { dirname, join } from 'node:path'
import { defineConfig } from 'vitepress'

const base = '/vox/'

// Pages of the previous site layout, kept alive as redirects.
const redirects: Record<string, string> = {
  'docs/usage': 'guide/routing',
  'docs/middleware': 'guide/middleware',
  'docs/context': 'guide/context',
  'docs/request': 'guide/request',
  'docs/response': 'guide/response',
  'docs/run': 'guide/running',
  'docs/pprof': 'middlewares/pprof',
}

export default defineConfig({
  lang: 'en-US',
  title: 'Vox',
  description: 'A Go web framework for humans, heavily inspired by Koa.',
  base,
  cleanUrls: true,

  themeConfig: {
    nav: [
      { text: 'Guide', link: '/getting-started', activeMatch: '^/(getting-started|guide/)' },
      { text: 'Middleware', link: '/middlewares/', activeMatch: '^/middlewares/' },
      { text: 'Recipes', link: '/recipes' },
      { text: 'API Reference', link: 'https://pkg.go.dev/github.com/aisk/vox' },
    ],

    sidebar: [
      {
        text: 'Introduction',
        items: [
          { text: 'Getting Started', link: '/getting-started' },
        ],
      },
      {
        text: 'Guide',
        items: [
          { text: 'Overview', link: '/guide/' },
          { text: 'Routing', link: '/guide/routing' },
          { text: 'Route Handlers', link: '/guide/handlers' },
          { text: 'Middleware', link: '/guide/middleware' },
          { text: 'Context', link: '/guide/context' },
          { text: 'Request', link: '/guide/request' },
          { text: 'Response', link: '/guide/response' },
          { text: 'Error Handling', link: '/guide/errors' },
          { text: 'Configuration', link: '/guide/configuration' },
          { text: 'Running', link: '/guide/running' },
          { text: 'Testing', link: '/guide/testing' },
        ],
      },
      {
        text: 'Bundled Middleware',
        items: [
          { text: 'Overview', link: '/middlewares/' },
          { text: 'Static Files', link: '/middlewares/static' },
          { text: 'Pprof', link: '/middlewares/pprof' },
        ],
      },
      {
        text: 'Cookbook',
        items: [
          { text: 'Recipes', link: '/recipes' },
        ],
      },
    ],

    outline: { level: [2, 3] },

    search: { provider: 'local' },

    socialLinks: [
      { icon: 'github', link: 'https://github.com/aisk/vox' },
    ],

    editLink: {
      pattern: 'https://github.com/aisk/vox/edit/master/docs/:path',
      text: 'Edit this page on GitHub',
    },

    footer: {
      message: 'Released under the MIT License.',
      copyright: 'Copyright © 2016-2026 aisk',
    },
  },

  buildEnd(site) {
    for (const [from, to] of Object.entries(redirects)) {
      const target = base + to
      const file = join(site.outDir, from + '.html')
      mkdirSync(dirname(file), { recursive: true })
      writeFileSync(file, `<!doctype html>
<meta charset="utf-8">
<title>Redirecting…</title>
<link rel="canonical" href="${target}">
<meta http-equiv="refresh" content="0; url=${target}">
<meta name="robots" content="noindex">
<a href="${target}">Click here if you are not redirected.</a>
`)
    }
  },
})
