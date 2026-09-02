import type { SidebarsConfig } from '@docusaurus/plugin-content-docs';

const sidebars: SidebarsConfig = {
  docsSidebar: [
    {
      type: 'category',
      label: '🚀 Primeiros Passos',
      collapsed: false,
      items: ['intro', 'installation', 'quickstart'],
    },
    {
      type: 'category',
      label: '🛡️ Conceitos & Segurança',
      collapsed: false,
      items: ['severity-levels'],
    },
    {
      type: 'category',
      label: '💻 Comandos da CLI',
      collapsed: false,
      items: [
        'commands/scan',
        'commands/check',
        'commands/init',
        'commands/fix',
        'commands/hook',
        'commands/version',
      ],
    },
    {
      type: 'category',
      label: '⚙️ Configuração',
      collapsed: false,
      items: ['configuration'],
    },
    {
      type: 'category',
      label: '🔄 Automação & CI/CD',
      collapsed: false,
      items: ['cicd-integration'],
    },
    {
      type: 'category',
      label: '🤝 Comunidade',
      collapsed: true,
      items: ['contributing'],
    },
  ],
};

export default sidebars;
