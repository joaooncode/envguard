---
sidebar_position: 5
title: envguard hook
---

# envguard hook

O comando **`envguard hook`** gerencia hooks de pré-commit locais no Git e realiza inspeções ultrarrápidas (`<10ms`) direcionadas exclusivamente aos arquivos que estão em stage (`git add`).

---

## Subcomandos

```bash
envguard hook <subcommand> [flags]
```

| Subcomando  | Descrição                                                                              |
| :---------- | :------------------------------------------------------------------------------------- |
| `install`   | Instala o script executável de pré-commit em `.git/hooks/pre-commit`.                  |
| `run`       | Executa a inspeção dos arquivos em stage e bloqueia commits com variáveis de ambiente. |
| `uninstall` | Remove com segurança o script instalado pelo envguard de `.git/hooks/pre-commit`.      |
| `help`      | Exibe a ajuda dos comandos de hook.                                                    |

---

## Flags Disponíveis

| Flag             | Subcomandos            | Tipo     | Padrão  | Descrição                                                 |
| :--------------- | :--------------------- | :------- | :------ | :-------------------------------------------------------- |
| `-p`, `--path`   | Todos                  | `string` | `"."`   | Diretório raiz do repositório Git.                        |
| `-f`, `--force`  | `install`, `uninstall` | `bool`   | `false` | Força a instalação ou remoção sobre hooks pré-existentes. |
| `-c`, `--config` | `run`                  | `string` | `""`    | Caminho para um arquivo de configuração customizado.      |
| `--no-color`     | `run`                  | `bool`   | `false` | Desativa cores ANSI na saída do terminal.                 |

---

## Exemplos de Uso

### 1. Instalar o Hook Nativo

Instala o script de hook no repositório local:

```bash
envguard hook install
```

Caso já exista um hook customizado ou de terceiros, o envguard evita sobrescrever acidentalmente. Utilize `--force` caso deseje substituir:

```bash
envguard hook install --force
```

### 2. Executar Inspeção em Stage (`hook run`)

Chamado automaticamente pelo Git antes de cada commit:

```bash
envguard hook run
```

Se nenhum arquivo sensível estiver em stage:

```text
✓ No unprotected environment files staged for commit.
```

Se um arquivo `.env` não permitido for adicionado para commit:

```text
🚨 Git Pre-Commit Check Failed!
Found 1 unprotected environment file(s) staged for commit:

  ✗ [STAGED] .env.production
    Message:     Environment file is staged for commit in Git index.
    Suggestions:
      • Unstage file: git restore --staged .env.production
      • Add to .gitignore

Commit blocked to prevent sensitive credentials from reaching Git.
```

### 3. Desinstalar o Hook Nativo

```bash
envguard hook uninstall
```

---

## Integração com o Framework Python `pre-commit`

O `envguard` possui suporte nativo ao framework popular [pre-commit](https://pre-commit.com/). Basta adicionar o arquivo `.pre-commit-config.yaml` no seu projeto:

```yaml
repos:
  - repo: https://github.com/joaooncode/envguard
    rev: v0.2.0
    hooks:
      - id: envguard
```
