---
sidebar_position: 1
title: envguard scan
---

# `envguard scan`

O comando `scan` realiza uma varredura completa no repositório atual em busca de arquivos de variáveis de ambiente (`.env`, `.env.*`, `*.env`), cruzando os achados com o status do Git.

---

## Sintaxe

```bash
envguard scan [flags]
```

---

## Flags Disponíveis

| Flag               | Tipo     | Descrição                                                                 | Padrão  |
| :----------------- | :------- | :------------------------------------------------------------------------ | :------ |
| `-p`, `--path`     | `string` | Diretório alvo para a varredura                                           | `"."`   |
| `-f`, `--format`   | `string` | Formato da saída (`text`, `terminal` ou `json`)                           | `text`  |
| `-s`, `--severity` | `string` | Nível mínimo de severidade (`all`, `info`, `warning`, `high`, `critical`) | `all`   |
| `-c`, `--config`   | `string` | Caminho para arquivo de configuração customizado                          | `""`    |
| `--no-color`       | `bool`   | Desativa cores ANSI na saída do terminal                                  | `false` |
| `--help`, `-h`     | `bool`   | Exibe a ajuda do comando                                                  | `false` |

---

## Exemplos de Uso

### 1. Varredura Padrão (Texto Formatado)

```bash
envguard scan
```

#### Exemplo de Saída:

```text
🛡️  envguard v0.2.0
Target: ./meu-projeto
──────────────────────────────────────────────────

Findings:
  ✗ [CRITICAL] .env
    Message:     Environment file is tracked by Git (committed in repository history).
    Suggestions:
      • Remove file from git tracking: git rm --cached .env
      • Add to .gitignore
      • Rotate any leaked credentials

  ⚠ [WARNING] .env.local
    Message:     Environment file exists locally and is not ignored by .gitignore.
    Suggestions:
      • Add to .gitignore

──────────────────────────────────────────────────
Summary:
  Total Findings: 2 (Critical: 1, High: 0, Warning: 1, Info: 0)
  Status:         ✗ FAILED
```

---

### 2. Saída em JSON Estruturado

Ideal para integração com scripts personalizados, automações ou ferramentas terceiras:

```bash
envguard scan --format json
```

#### Exemplo de Saída JSON:

```json
{
  "version": "v0.1.0",
  "repository": "./meu-projeto",
  "total_findings": 2,
  "critical": 1,
  "warning": 1,
  "findings": [
    {
      "path": ".env",
      "severity": "CRITICAL",
      "status": "tracked",
      "message": "File is tracked by Git version control"
    },
    {
      "path": ".env.local",
      "severity": "WARNING",
      "status": "unignored",
      "message": "File is not ignored by .gitignore"
    },
    {
      "path": ".env.example",
      "severity": "INFO",
      "status": "allowed_template",
      "message": "Recognized as safe environment template"
    }
  ]
}
```

---

## Códigos de Retorno (Exit Codes)

| Código | Significado                                                          |
| :----: | :------------------------------------------------------------------- |
|  `0`   | Nenhuma violação bloqueante encontrada.                              |
|  `1`   | Encontrada ao menos uma violação de severidade `CRITICAL` ou `HIGH`. |
