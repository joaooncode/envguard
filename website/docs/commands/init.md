---
sidebar_position: 3
title: envguard init
---

# envguard init

O comando **`envguard init`** inicializa o arquivo de configuração do projeto (`.envguard.yaml`) e permite gerar templates seguros de variáveis de ambiente (`.env.example`) através de higienização automatizada.

---

## Uso

```bash
envguard init [flags]
```

---

## Flags Disponíveis

| Flag               | Tipo     | Padrão  | Descrição                                                                   |
| :----------------- | :------- | :------ | :-------------------------------------------------------------------------- |
| `-p`, `--path`     | `string` | `"."`   | Diretório alvo para criar os arquivos de configuração e template.           |
| `-f`, `--force`    | `bool`   | `false` | Sobrescreve arquivos existentes (`.envguard.yaml` ou `.env.example`).       |
| `-t`, `--template` | `bool`   | `false` | Gera automaticamente um template seguro `.env.example`.                     |
| `--template-from`  | `string` | `""`    | Caminho do arquivo `.env` de origem a ser sanitizado para criar o template. |

---

## Exemplos de Uso

### 1. Inicializar configuração padrão

Cria o arquivo `.envguard.yaml` com comentários explicativos e padrões de segurança:

```bash
envguard init
```

### 2. Inicializar configuração e gerar `.env.example`

Se já existir um arquivo `.env` local, ele é automaticamente lido e higienizado (valores removidos, mantendo chaves e comentários):

```bash
envguard init --template
```

### 3. Sanitizar um arquivo `.env` específico

Gera o `.env.example` a partir de um arquivo com outro nome (ex.: `.env.production` ou `.env.local`):

```bash
envguard init --template-from .env.production
```

### 4. Forçar sobrescrita em caso de arquivos existentes

```bash
envguard init --path ./meu-projeto --template --force
```

---

## Como Funciona a Sanitização

A sanitização do `envguard init` garante que nenhum segredo seja exposto:

1. Preserva comentários (`# ...`), quebras de linha e estrutura do arquivo.
2. Mantém os nomes das variáveis (ex.: `DATABASE_URL=`, `JWT_SECRET=`).
3. Limpa com segurança todos os valores atribuídos.
