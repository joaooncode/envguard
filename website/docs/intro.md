---
sidebar_position: 1
title: Visão Geral
---

# envguard

> **Prevent `.env` files and environment secrets from accidentally reaching Git.**

O **`envguard`** é uma ferramenta open source de linha de comando (CLI) desenvolvida em **Go**, projetada para atuar como uma camada leve de proteção entre o desenvolvedor e o Git. O foco principal é detectar, alertar e prevenir a exposição indevida de arquivos de variáveis de ambiente (`.env`, `.env.production`, `.env.local`, etc.) em repositórios de código.

---

## Por que o `envguard`?

Arquivos `.env` frequentemente armazenam dados sensíveis: credenciais de banco de dados, chaves de API, segredos JWT, tokens e certificados privados. Um descuido na configuração do `.gitignore` ou um simples `git add .` desavisado pode comitar esses segredos permanentemente no histórico de versão.

O `envguard` resolve isso com foco específico em arquivos de ambiente:

- **Git-Aware:** Entende o estado do repositório — diferencia se um `.env` está rastreado (_tracked_), preparado para commit (_staged_), ignorado ou desprotegido.
- **Rápido & Local:** Desenvolvido em Go nativo, funciona 100% offline, sem envio de dados para servidores externos. Ideal para execução local, _pre-commit hooks_ e pipelines de CI/CD.
- **Seguro por Design:** Nunca imprime ou expõe valores de variáveis ou segredos em logs e saídas do terminal.
- **Pronto para CI/CD:** Suporta formato JSON estruturado (`--format json`) e códigos de saída determinísticos para automação.

---

## Como Funciona

```text
[Desenvolvedor / CI]
        │
        ▼
   envguard scan
        │
   ┌────┴───────────────────────────┐
   │ 1. Varre diretórios por .env   │
   │ 2. Consulta o status do Git    │
   │ 3. Checa regras do .gitignore  │
   └────┬───────────────────────────┘
        ▼
[Relatório com Níveis de Severidade]
 ✗ .env         CRITICAL  (Rastreado no Git)
 ⚠ .env.local   WARNING   (Não protegido pelo .gitignore)
 ✓ .env.example INFO      (Template permitido)
```

---

## Próximos Passos

- [Instalação](./installation.md): Saiba como instalar o binário localmente.
- [Comandos](./commands/scan.md): Conheça os comandos `scan`, `check`, `init`, `fix`, `hook` e `version`.
- [Arquivo de Configuração](./configuration.md): Personalize regras, exceções e severidades via `.envguard.yaml`.
- [Níveis de Severidade](./severity-levels.md): Entenda os alertas e as ações recomendadas.
- [Integração CI/CD](./cicd-integration.md): Automatize a proteção no GitHub Actions e pre-commit hooks.
