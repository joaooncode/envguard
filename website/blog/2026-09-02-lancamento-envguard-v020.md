---
slug: lancamento-envguard-v020
title: Lançamento do envguard v0.2.0 - Hooks, Configuração, Init e Fix
authors: [joaooncode]
tags: [release, security, cli, golang, git-hooks, automation]
---

É com muito orgulho que anunciamos o lançamento do **`envguard` v0.2.0**! 🎉

Esta versão consolida a automação e remediação do ecossistema do `envguard`, trazendo recursos nativos para proteger o fluxo de desenvolvimento antes mesmo da execução de commits.

<!-- truncate -->

## Principais Novidades da v0.2.0

### 1. Suporte a Git Pre-commit Hooks (`envguard hook`)

- **Instalação Nativa (`envguard hook install`):** Cria o script executável em `.git/hooks/pre-commit` com detecção de assinatura e proteção contra sobrescrita acidental.
- **Inspeção Instantânea (`envguard hook run`):** Avalia exclusivamente os arquivos preparados para commit (`staged`) em menos de 10ms, bloqueando o commit se houver variáveis de ambiente desprotegidas.
- **Integração Oficial com Python `pre-commit`:** Arquivo `.pre-commit-hooks.yaml` incluído na raiz do repositório para adoção imediata.

### 2. Remediação Automática com `envguard fix`

- Analisa o repositório e insere automaticamente os arquivos de ambiente desprotegidos no `.gitignore` da raiz.
- Suporte a simulação não-destrutiva com `--dry-run` para pré-visualizar as alterações antes de aplicar no disco.

### 3. Inicialização e Sanitização de Templates (`envguard init`)

- Gera o arquivo `.envguard.yaml` completo e comentado.
- Suporte à geração de templates `.env.example` através da higienização automática (`--template` e `--template-from`), removendo valores sensíveis e preservando chaves e comentários.

### 4. Arquivo de Configuração `.envguard.yaml`

- Controle total de `ignore_dirs`, `custom_patterns`, `allowlist` e `severity_overrides` por projeto ou via flag `--config`.

---

## Como Atualizar

```bash
go install github.com/joaooncode/envguard/cmd/envguard@latest
```

Verifique a versão instalada:

```bash
envguard version
# envguard v0.2.0
```

Confira a [documentação oficial](/docs/intro) para explorar todos os novos comandos e guias de integração!
