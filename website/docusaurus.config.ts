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

  themeConfig: {
    image: 'img/docusaurus-social-card.jpg',
    colorMode: {
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
          sidebarId: 'tutorialSidebar',
          position: 'left',
          label: 'Documentação',
        },
        { to: '/blog', label: 'Blog / Notas de Versão', position: 'left' },
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
          title: 'Documentação',
          items: [
            { label: 'Visão Geral', to: '/docs/intro' },
            { label: 'Instalação', to: '/docs/installation' },
            { label: 'Comandos CLI', to: '/docs/commands/scan' },
            { label: 'Níveis de Severidade', to: '/docs/severity-levels' },
          ],
        },
        {
          title: 'Integrações',
          items: [
            { label: 'CI/CD & Git Hooks', to: '/docs/cicd-integration' },
            { label: 'Guia de Contribuição', to: '/docs/contributing' },
          ],
        },
        {
          title: 'Comunidade',
          items: [
            { label: 'Repositório GitHub', href: 'https://github.com/joaooncode/envguard' },
            {
              label: 'Reportar Problema (Issue)',
              href: 'https://github.com/joaooncode/envguard/issues',
            },
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
