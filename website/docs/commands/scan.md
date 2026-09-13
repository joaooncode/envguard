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
🛡️  envguard v0.3.0
Target: ./meu-projeto
──────────────────────────────────────────────────

Findings:
  ✗ [CRITICAL] .env
    Message:     Environment file is tracked by Git (committed in repository history).
    Suggestions:
      • Remove file from git tracking: git rm --cached .env
      • Add to .gitignore
      • Rotate any leaked credentials
    Secrets:
      • [CRITICAL] AWS_ACCESS_KEY_ID (aws) at line 3

  ⚠ [WARNING] .env.local
    Message:     Environment file exists locally and is not ignored by .gitignore.
    Suggestions:
      • Add to .gitignore

──────────────────────────────────────────────────
Summary:
  Total Findings: 2 (Critical: 1, High: 0, Warning: 1, Info: 0)
  Status:         ✗ FAILED
```

Quando o arquivo contém segredos reais detectados pelo [Secret Scanner](../configuration.md#secret-scanning-por-conteúdo), cada ocorrência é listada em `Secrets:` com sua própria severidade, chave, provedor (se aplicável) e linha — **o valor do segredo nunca é exibido**.

---

### 2. Saída em JSON Estruturado

Ideal para integração com scripts personalizados, automações ou ferramentas terceiras:

```bash
envguard scan --format json
```

#### Exemplo de Saída JSON:

```json
{
  "version": "0.3.0",
  "timestamp": "2026-09-12T12:00:00Z",
  "scanned_dir": "./meu-projeto",
  "findings": [
    {
      "path": ".env",
      "severity": "critical",
      "message": "Environment file is tracked by Git (committed in repository history).",
      "suggestions": [
        "Remove file from git tracking: git rm --cached .env",
        "Add to .gitignore",
        "Rotate any leaked credentials"
      ],
      "git_status": {
        "is_repo": true,
        "is_tracked": true,
        "is_staged": false,
        "is_ignored": false
      },
      "is_allowed": false,
      "secret_matches": [
        {
          "line": 3,
          "key": "AWS_ACCESS_KEY_ID",
          "method": "pattern",
          "provider": "aws",
          "severity": "critical"
        }
      ]
    },
    {
      "path": ".env.local",
      "severity": "warning",
      "message": "Environment file exists locally and is not ignored by .gitignore.",
      "suggestions": ["Add to .gitignore"],
      "git_status": {
        "is_repo": true,
        "is_tracked": false,
        "is_staged": false,
        "is_ignored": false
      },
      "is_allowed": false
    }
  ],
  "summary": {
    "total": 2,
    "critical": 1,
    "high": 0,
    "warning": 1,
    "info": 0,
    "passed": false
  }
}
```

---

## Códigos de Retorno (Exit Codes)

| Código | Significado                                                          |
| :----: | :------------------------------------------------------------------- |
|  `0`   | Nenhuma violação bloqueante encontrada.                              |
|  `1`   | Encontrada ao menos uma violação de severidade `CRITICAL` ou `HIGH`. |
