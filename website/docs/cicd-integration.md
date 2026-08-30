---
sidebar_position: 5
title: Integração CI/CD & Git Hooks
---

# Integração CI/CD & Git Hooks

Automatizar a verificação com `envguard` impede que segredos entrem no repositório antes mesmo de um Pull Request ser aprovado ou antes do desenvolvedor concluir um commit.

---

## 1. GitHub Actions

Adicione um job simples ao seu arquivo `.github/workflows/security.yml`:

```yaml
name: Security & Env Check

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main, develop]

jobs:
  envguard:
    name: Check Env Files
    runs-on: ubuntu-latest
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'

      - name: Install envguard
        run: go install github.com/joaooncode/envguard/cmd/envguard@latest

      - name: Run envguard check
        run: envguard check
```

---

## 2. Pre-commit Hook (Husky / Git Hooks)

Você pode impedir commits acidentais de arquivos `.env` diretamente na máquina do desenvolvedor.

### Com Husky (Node.js)

No arquivo `.husky/pre-commit`:

```bash
#!/usr/bin/env sh
. "$(dirname -- "$0")/_/husky.sh"

envguard check
```

### Git Hook Nativo (`.git/hooks/pre-commit`)

```bash
#!/bin/sh
if command -v envguard >/dev/null 2>&1; then
  envguard check
  if [ $? -ne 0 ]; then
    echo "envguard bloqueou o commit devido a arquivos de ambiente desprotegidos!"
    exit 1
  fi
fi
```

---

## 3. GitLab CI/CD

No seu `.gitlab-ci.yml`:

```yaml
envguard_security:
  stage: test
  image: golang:1.22-alpine
  script:
    - go install github.com/joaooncode/envguard/cmd/envguard@latest
    - envguard check
```
