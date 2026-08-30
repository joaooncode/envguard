---
sidebar_position: 6
title: Como Contribuir
---

# Guia de Contribuição

Obrigado pelo interesse em contribuir com o **`envguard`**! Como um projeto Open Source mantido pela comunidade, toda ajuda é bem-vinda: novas features, correções de bugs, melhorias de documentação e ideias.

---

## Ambiente de Desenvolvimento

### Pré-requisitos

- **Go** 1.22 ou superior
- **Node.js** 20+ (para hooks de commit e documentação Docusaurus)
- **Git**

### Configurando o Repositório Localmente

```bash
# 1. Clone o repositório
git clone https://github.com/joaooncode/envguard.git
cd envguard

# 2. Instale as ferramentas de linting e hooks
npm install

# 3. Compile e teste o projeto
go test ./...
go build ./cmd/envguard
```

---

## Padrão de Commits (Conventional Commits)

Este projeto adota a convenção [Conventional Commits](https://www.conventionalcommits.org/). Todos os commits são validados automaticamente via Commitlint e Husky.

### Estrutura do Commit

```text
<tipo>(<escopo opcional>): <descrição>
```

### Tipos Aceitos

- `feat`: Nova funcionalidade
- `fix`: Correção de bug
- `refactor`: Refatoração de código sem mudança de comportamento
- `docs`: Alterações na documentação
- `test`: Adição ou ajuste de testes automatizados
- `style`: Formatação ou estilo de código
- `perf`: Melhoria de desempenho
- `chore`: Atualização de dependências ou tarefas de build

---

## Fluxo de Pull Request

1. Crie uma branch para a sua modificação:
   ```bash
   git checkout -b feat/minha-melhoria
   ```
2. Faça suas alterações e adicione testes correspondentes se aplicável.
3. Garanta que a formatação e testes estejam passando:
   ```bash
   npm run lint:go
   go test ./...
   ```
4. Commit suas alterações:
   ```bash
   git commit -m "feat(cli): adiciona nova flag de verbose"
   ```
5. Envie a branch para o seu fork e abra o **Pull Request** no repositório oficial.
