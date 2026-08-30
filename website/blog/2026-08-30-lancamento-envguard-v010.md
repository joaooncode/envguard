---
slug: lancamento-envguard-v010
title: Apresentando o envguard v0.1.0 (MVP)
authors: [joaooncode]
tags: [release, security, cli, golang]
---

Temos o prazer de anunciar o lançamento do **`envguard` v0.1.0**!

O `envguard` é uma ferramenta open source de linha de comando desenvolvida em Go para impedir que arquivos `.env` e segredos de variáveis de ambiente sejam comitados acidentalmente em repositórios Git.

<!-- truncate -->

## Destaques da versão v0.1.0

- **Detecção Git-Aware:** Reconhecimento automático de arquivos `.env` rastreados (_tracked_), na stage (_staged_) ou desprotegidos fora do `.gitignore`.
- **Varredura Ultrarrápida:** Construído em Go nativo, 100% offline e sem envio de telemetria ou segredos para a rede.
- **Saída Flexível:** Relatórios visuais no terminal e saída estruturada em JSON (`--format json`).
- **Validação para CI/CD:** Comando `envguard check` com códigos de saída determinísticos (`0` ou `1`) para bloquear pipelines inseguras.

## Como Instalar

```bash
go install github.com/joaooncode/envguard/cmd/envguard@latest
```

E execute:

```bash
envguard scan
```

Confira a [documentação completa](/docs/intro) para saber mais sobre como integrar no seu dia a dia e nos seus pipelines de CI/CD!
