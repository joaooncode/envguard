import { themes as prismThemes } from 'prism-react-renderer';
import type { Config } from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';

const config: Config = {
  title: 'envguard',
  tagline: 'Prevent .env files and environment secrets from accidentally reaching Git.',
  favicon: 'img/favicon.ico',

  // Configurações para GitHub Pages
  url: 'https://joaooncode.github.io',
  baseUrl: '/envguard/',
  organizationName: 'joaooncode',
  projectName: 'envguard',
  trailingSlash: false,

  onBrokenLinks: 'throw',

  i18n: {
    defaultLocale: 'pt-BR',
    locales: ['pt-BR'],
  },

  presets: [
    [
      'classic',
      {
        docs: {
          sidebarPath: './sidebars.ts',
          editUrl: 'https://github.com/joaooncode/envguard/tree/main/website/',
        },
        blog: {
          showReadingTime: true,
          feedOptions: {
            type: ['rss', 'atom'],
            xslt: true,
          },
          editUrl: 'https://github.com/joaooncode/envguard/tree/main/website/',
        },
        theme: {
          customCss: './src/css/custom.css',
        },
      } satisfies Preset.Options,
    ],
  ],

  themes: [
    [
      require.resolve('@easyops-cn/docusaurus-search-local'),
      {
        hashed: true,
        language: ['en'],
        indexDocs: true,
        indexBlog: true,
        indexPages: true,
        docsRouteBasePath: '/docs',
        highlightSearchTermsOnTargetPage: true,
        searchResultLimits: 8,
        searchResultContextMaxLength: 60,
      },
    ],
  ],

  themeConfig: {
    image: 'img/docusaurus-social-card.jpg',
    colorMode: {
      defaultMode: 'dark',
      respectPrefersColorScheme: true,
    },
    navbar: {
      title: 'envguard',
      logo: {
        alt: 'envguard Logo',
        src: 'img/logo.svg',
      },
      items: [
        {
          type: 'docSidebar',
          sidebarId: 'docsSidebar',
          position: 'left',
          label: 'Documentação',
        },
        { to: '/docs/quickstart', label: 'Início Rápido', position: 'left' },
        { to: '/blog', label: 'Blog & Versões', position: 'left' },
        {
          href: 'https://github.com/joaooncode/envguard',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {
          title: 'Começando',
          items: [
            { label: 'Visão Geral', to: '/docs/intro' },
            { label: 'Instalação', to: '/docs/installation' },
            { label: 'Início Rápido', to: '/docs/quickstart' },
          ],
        },
        {
          title: 'Comandos & Regras',
          items: [
            { label: 'Varredura (scan)', to: '/docs/commands/scan' },
            { label: 'Verificação (check)', to: '/docs/commands/check' },
            { label: 'Níveis de Severidade', to: '/docs/severity-levels' },
            { label: 'Configuração (.envguard.yaml)', to: '/docs/configuration' },
          ],
        },
        {
          title: 'Automação & Comunidade',
          items: [
            { label: 'CI/CD & Git Hooks', to: '/docs/cicd-integration' },
            { label: 'Guia de Contribuição', to: '/docs/contributing' },
            { label: 'GitHub Repository', href: 'https://github.com/joaooncode/envguard' },
            { label: 'Reportar Problema', href: 'https://github.com/joaooncode/envguard/issues' },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} envguard. Criado por João Vitor. Distribuído sob a licença MIT.`,
    },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
      additionalLanguages: ['bash', 'json', 'yaml', 'go', 'powershell'],
    },
  } satisfies Preset.ThemeConfig,
};

export default config;
