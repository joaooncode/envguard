# envguard

> **Prevent `.env` files and environment secrets from accidentally reaching Git.**

[![Go Version](https://img.shields.io/badge/go-1.22%2B-blue.svg)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)](CODE_OF_CONDUCT.md)

`envguard` é uma ferramenta open source de linha de comando (CLI) desenvolvida em **Go**, projetada para atuar como uma camada leve de proteção entre o desenvolvedor e o Git. O foco principal é detectar, alertar, remediar e prevenir a exposição indevida de arquivos de variáveis de ambiente (`.env`, `.env.production`, `.env.local`, etc.) em repositórios.

---

## Por que o `envguard`?

Arquivos `.env` costumam armazenar dados sensíveis: credenciais de banco de dados, chaves de API, tokens e certificados. Um descuido na configuração do `.gitignore` ou um simples `git add .` desavisado pode comitar esses segredos no histórico de versão.

O `envguard` resolve isso com foco específico em arquivos de ambiente:

- **Git-Aware:** Entende o estado do repositório — diferencia se um `.env` está rastreado (_tracked_), preparado (_staged_), ignorado ou desprotegido.
- **Rápido & Local:** Funciona 100% offline, sem envio de dados para servidores externos. Ideal para execução local, _pre-commit hooks_ e pipelines de CI/CD.
- **Remediação Automática:** Adiciona padrões ausentes ao `.gitignore` automaticamente (`envguard fix`).
- **Hooks Nativos & Framework Pre-commit:** Instalação direta em `.git/hooks/pre-commit` e suporte ao framework Python `pre-commit`.
- **Seguro por Design:** Nunca imprime ou expõe valores de variáveis ou segredos em logs e saídas do terminal.
- **Pronto para CI/CD:** Suporta formato JSON estruturado (`--format json`) e códigos de saída determinísticos para automação.

---

## Instalação

### Via `go install` (Requer Go 1.22+)

```bash
go install github.com/joaooncode/envguard/cmd/envguard@latest
```

### Compilando do código-fonte

```bash
# Clone o repositório
git clone https://github.com/joaooncode/envguard.git
cd envguard

# Compile o binário
go build -o envguard ./cmd/envguard
```

_(Distribuição futura via Homebrew, Scoop, WinGet e GitHub Releases)_

---

## Uso & Comandos

### 1. Varredura Local (`scan`)

Analisa o diretório atual em busca de arquivos de ambiente e valida o estado no Git:

```bash
envguard scan
```

Exemplo de saída no terminal:

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

### 2. Validação para CI/CD (`check`)

Ideal para pipelines e automações. Retorna código de erro (`exit code 1`) caso encontre violações bloqueantes:

```bash
envguard check
```

### 3. Git Pre-Commit Hooks (`hook`)

Instala ou executa inspeções ultrarrápidas (<10ms) focadas exclusivamente em arquivos preparados para commit (`staged`):

```bash
# Instalar o hook nativo em .git/hooks/pre-commit
envguard hook install

# Executar checagem de stage (bloqueia commits com .env não permitidos)
envguard hook run

# Desinstalar o hook nativo
envguard hook uninstall
```

#### Integração com Python `pre-commit`:

Adicione ao seu `.pre-commit-config.yaml`:

```yaml
repos:
  - repo: https://github.com/joaooncode/envguard
    rev: v0.2.0
    hooks:
      - id: envguard
```

### 4. Remediação Automática (`fix`)

Adiciona automaticamente padrões correspondentes para arquivos desprotegidos (`WARNING`) no `.gitignore` da raiz, preservando comentários e formatação existente:

```bash
# Aplicar correções no .gitignore
envguard fix

# Simular alterações propostas sem modificar arquivos
envguard fix --dry-run

# Executar em diretório específico
envguard fix --path ./meu-projeto
```

### 5. Inicialização de Configuração e Templates (`init`)

Gera o arquivo de configuração `.envguard.yaml` documentado e, opcionalmente, cria templates `.env.example` sanitizados a partir de variáveis locais:

```bash
# Inicializar .envguard.yaml padrão
envguard init

# Inicializar configuração e gerar template .env.example sanitizado
envguard init --template

# Inicializar a partir de arquivo de origem específico
envguard init --template-from .env.production
```

### 6. Saída Estruturada em JSON

```bash
envguard scan --format json
```

### 7. Verificar Versão

```bash
envguard version
```

---

## Arquivo de Configuração (`.envguard.yaml`)

O `envguard` pode ser personalizado criando um arquivo `.envguard.yaml` na raiz do repositório:

```yaml
version: '1'

scanner:
  ignore_dirs:
    - 'node_modules'
    - '.git'
    - 'vendor'

detector:
  custom_patterns:
    - '*.env.vault'
  allowlist:
    - '.env.example'
    - '.env.sample'
    - '.env.template'
  severity_overrides:
    - pattern: '.env.test'
      severity: 'warning'
```

---

## Níveis de Severidade

| Nível               | Situação                                                                | Ação Recomendada                                                             |
| :------------------ | :---------------------------------------------------------------------- | :--------------------------------------------------------------------------- |
| **`CRITICAL`**      | Arquivo de ambiente rastreado (_tracked_) no histórico Git              | Remover do rastreamento (`git rm --cached`) e rotacionar credenciais         |
| **`HIGH`**          | Arquivo de ambiente adicionado para commit (_staged_)                   | Retirar da stage (`git restore --staged <file>`) e adicionar ao `.gitignore` |
| **`WARNING`**       | Arquivo existe localmente mas **não está** no `.gitignore`              | Executar `envguard fix` ou adicionar padrão ao `.gitignore`                  |
| **`INFO` / `SAFE`** | Arquivo protegido ou template permitido (`.env.example`, `.env.sample`) | Nenhuma ação necessária                                                      |

---

## Padrões e Exceções Padrão

- **Padrões monitorados:** `.env`, `.env.*`, `*.env`
- **Exceções seguras permitidas por padrão:** `.env.example`, `.env.sample`, `.env.template`

---

## Roadmap

- [x] **v0.1.0 (MVP):**
  - [x] Deteção de `.env` e variantes
  - [x] Integração Git (_tracked_, _staged_, _gitignore_)
  - [x] Relatórios em Terminal e JSON
  - [x] Códigos de saída para CI/CD
- [x] **v0.2.0:**
  - [x] `envguard init` (criação automática de `.envguard.yaml` e templates)
  - [x] `envguard fix` (auxílio na adição automática ao `.gitignore`)
  - [x] Suporte a arquivo de configuração `.envguard.yaml` e flag `--config`
  - [x] Instalação de _Git Pre-commit Hooks_ nativos e suporte a Python `pre-commit`
- [ ] **v0.3.0:**
  - [ ] Secret scanning básico por conteúdo & cálculo de entropia
  - [ ] Deteção de padrões comuns de chaves (AWS, Stripe, GitHub, etc.)
- [ ] **v1.0.0:**
  - [ ] GitHub Action oficial do envguard
  - [ ] Pacotes para Homebrew, Scoop, WinGet e AUR

---

## Como Contribuir

Contribuições são super bem-vindas! Como um projeto Open Source mantido pela comunidade:

1. Faça um **Fork** do projeto.
2. Crie uma branch para sua funcionalidade/correção: `git checkout -b feat/minha-feature`.
3. Commit as suas alterações seguindo [Conventional Commits](https://www.conventionalcommits.org/):
   `git commit -m "feat: adiciona nova funcionalidade"`.
4. Envie para a sua branch: `git push origin feat/minha-feature`.
5. Abra um **Pull Request**.

Por favor, leia o nosso [Código de Conduta](CODE_OF_CONDUCT.md) antes de interagir na comunidade.

---

## Licença

Distribuído sob a licença **MIT**. Veja o arquivo [`LICENSE`](LICENSE) para mais detalhes.
