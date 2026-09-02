---
sidebar_position: 4
title: envguard fix
---

# envguard fix

O comando **`envguard fix`** realiza a remediação automática de arquivos de ambiente desprotegidos (`WARNING`), inserindo os padrões correspondentes no `.gitignore` da raiz do repositório de forma não destrutiva.

---

## Uso

```bash
envguard fix [flags]
```

---

## Flags Disponíveis

| Flag              | Tipo     | Padrão  | Descrição                                                                    |
| :---------------- | :------- | :------ | :--------------------------------------------------------------------------- |
| `-p`, `--path`    | `string` | `"."`   | Diretório alvo para varredura e remediação.                                  |
| `-d`, `--dry-run` | `bool`   | `false` | Pré-visualiza as alterações propostas no `.gitignore` sem modificar o disco. |
| `-c`, `--config`  | `string` | `""`    | Caminho para um arquivo de configuração customizado.                         |
| `--no-color`      | `bool`   | `false` | Desativa cores ANSI na saída do terminal.                                    |

---

## Exemplos de Uso

### 1. Aplicar correções automáticas no `.gitignore`

```bash
envguard fix
```

Saída de exemplo:

```text
✓ Successfully updated .gitignore with 2 rule(s):
  + .env
  + /services/api/.env.local
```

### 2. Modo Dry-Run (Simulação)

Permite inspecionar quais regras seriam adicionadas antes de alterar qualquer arquivo:

```bash
envguard fix --dry-run
```

Saída de exemplo:

```text
🔍 Dry run mode: changes will not be written to disk

Proposed .gitignore additions:
  + .env
  + .env.local
```

### 3. Alertas para Arquivos Rastreados (`CRITICAL`)

Se um arquivo já foi commitado no Git, o `.gitignore` não é suficiente para remover o histórico. O `envguard fix` detecta a situação, exibe instruções práticas de remoção do cache (`git rm --cached <arquivo>`) e retorna `exit code 1` alertando sobre a pendência.

---

## Características de Segurança

- **Formatação Não Destrutiva:** Novas regras são adicionadas sob o cabeçalho `# Added by envguard`, preservando comentários existentes, indentação e quebras de linha.
- **Prevenção de Duplicatas:** Verifica regras já presentes no `.gitignore` e não insere padrões repetidos.
- **Resolução de Caminhos Relativos:** Arquivos localizados em subpastas (ex: `packages/backend/.env`) são mapeados corretamente em relação à raiz do repositório.
