---
sidebar_position: 2
title: envguard check
---

# `envguard check`

O comando `check` é otimizado para **CI/CD e automações**. Ele executa uma validação rigorosa e retorna código de saída `1` caso identifique qualquer risco iminente ou arquivo de ambiente exposto.

---

## Sintaxe

```bash
envguard check [flags]
```

---

## Quando usar `check` vs `scan`?

- **`envguard scan`:** Focado na experiência do desenvolvedor local, gerando relatórios visuais e informativos.
- **`envguard check`:** Focado em validações estritas de pipelines (GitHub Actions, GitLab CI, pre-commit hooks), falhando imediatamente caso uma regra de segurança seja quebrada.

---

## Exemplo de Uso em CI

```bash
envguard check
```

Se todos os arquivos estiverem devidamente ignorados e nenhum `.env` estiver rastreado/staged, o comando encerra silenciosamente ou com resumo de sucesso (`Exit code 0`).

Se houver violações (`CRITICAL` ou `HIGH`), o pipeline será interrompido com `Exit code 1`.

---

## Exemplo de Workflow GitHub Actions

```yaml
- name: Check Environment Security
  run: |
    go install github.com/joaooncode/envguard/cmd/envguard@latest
    envguard check
```
