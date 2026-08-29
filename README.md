# envguard

> **Prevent `.env` files and environment secrets from accidentally reaching Git.**

[![Go Version](https://img.shields.io/badge/go-1.22%2B-blue.svg)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)](CODE_OF_CONDUCT.md)

`envguard` é uma ferramenta open source de linha de comando (CLI) desenvolvida em **Go**, projetada para atuar como uma
camada leve de proteção entre o desenvolvedor e o Git. O foco principal é detectar, alertar e prevenir a exposição
indevida de arquivos de variáveis de ambiente (`.env`, `.env.production`, `.env.local`, etc.) em repositórios.

---

## Por que o `envguard`?

Arquivos `.env` costumam armazenar dados sensíveis: credenciais de banco de dados, chaves de API, tokens e certificados.
Um descuido na configuração do `.gitignore` ou um simples `git add .` desavisado pode comitar esses segredos no
histórico de versão.

O `envguard` resolve isso com foco específico em arquivos de ambiente:

- **Git-Aware:** Entende o estado do repositório — diferencia se um `.env` está rastreado (_tracked_), preparado
  (_staged_), ignorado ou desprotegido.
- **Rápido & Local:** Funciona 100% offline, sem envio de dados para servidores externos. Ideal para execução local,
  _precommit hooks_ e pipelines de CI/CD.
- **Seguro por Design:** Nunca imprime ou expõe valores de variáveis, ou segredos em logs, ou saídas do terminal.
- **Pronto para CI/CD:** Suporta formato JSON estruturado (`--format json`) e códigos de saída determinísticos para
  automação.

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
envguard v0.1.0

Repository: ./my-project
Scanning...

  ✗ .env             CRITICAL   tracked by Git
  ⚠ .env.local       WARNING    not covered by .gitignore
  ✓ .env.example     INFO       allowed template file

Found 2 finding(s) (1 CRITICAL, 1 WARNING)
Exit code: 1
```

### 2. Validação para CI/CD (`check`)

Ideal para pipelines e automações. Retorna código de erro (`exit code 1`) caso encontre violações bloqueantes:

```bash
envguard check
```

### 3. Saída Estruturada em JSON

```bash
envguard scan --format json
```

### 4. Verificar Versão

```bash
envguard version
```

---

## Níveis de Severidade

| Nível               | Situação                                                                | Ação Recomendada                                                       |
| :------------------ | :---------------------------------------------------------------------- | :--------------------------------------------------------------------- |
| **`CRITICAL`**      | Arquivo de ambiente rastreado (_tracked_) no histórico Git              | Remover do rastreamento (`git rm --cached`) e rotacionar credenciais   |
| **`HIGH`**          | Arquivo de ambiente adicionado para commit (_staged_)                   | Retirar da stage (`git reset HEAD <file>`) e adicionar ao `.gitignore` |
| **`WARNING`**       | Arquivo existe localmente mas **não está** no `.gitignore`              | Adicionar padrão correspondente ao `.gitignore`                        |
| **`INFO` / `SAFE`** | Arquivo protegido ou template permitido (`.env.example`, `.env.sample`) | Nenhuma ação necessária                                                |

---

## Padrões e Exceções Padrão

- **Padrões monitorados:** `.env`, `.env.*`, `*.env`
- **Exceções seguras permitidas por padrão:** `.env.example`, `.env.sample`, `.env.template`

_(O suporte a configurações personalizadas via arquivo `.envguard.yaml` está no roadmap da v0.2)_

---

## Roadmap

- [ ] **v0.1.0 (MVP):**
  - [ ] Deteção de `.env` e variantes
  - [ ] Integração Git (_tracked_, _staged_, _gitignore_)
  - [ ] Relatórios em Terminal e JSON
  - [ ] Códigos de saída para CI/CD
- [ ] **v0.2.0:**
  - [ ] `envguard init` (criação automática de `.envguard.yaml` e templates)
  - [ ] `envguard fix` (auxílio na adição automática ao `.gitignore`)
  - [ ] Instalação de _Git Precommit Hooks_
- [ ] **v0.3.0:**
  - [ ] Secret scanning básico por conteúdo & cálculo de entropia
- [ ] **v1.0.0:**
  - [ ] GitHub Action oficial
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
