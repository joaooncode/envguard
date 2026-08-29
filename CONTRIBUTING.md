# 🤝 Guia de Contribuição — envguard

Primeiramente, muito obrigado por dedicar seu tempo para contribuir com o **envguard**! 🎉  
Projetos de código aberto são feitos por e para a comunidade, e toda a contribuição — seja corrigindo um bug, sugerindo
uma funcionalidade, melhorando a documentação ou escrevendo testes — é muito bem-vinda.

---

## Sumário

- [Código de Conduta](#-código-de-conduta)
- [Como Contribuir?](#-como-contribuir)
  - [Reportando Bugs](#reportando-bugs)
  - [Sugerindo Melhorias](#sugerindo-melhorias)
  - [Contribuindo com Código](#contribuindo-com-código)
- [Ambiente de Desenvolvimento](#-ambiente-de-desenvolvimento)
- [Padrões de Código e Qualidade](#-padrões-de-código-e-qualidade)
- [Convenção de Commits (Conventional Commits)](#-convenção-de-commits-conventional-commits)
- [Fluxo de Pull Requests](#-fluxo-de-pull-requests)

---

## Código de Conduta

Ao participar deste projeto, você concorda em cumprir o nosso [Código de Conduta](CODE_OF_CONDUCT.md). Esperamos que todos
os participantes mantenham um ambiente respeitoso, inclusivo e seguro para todos.

---

## Como Contribuir?

### Reportando Bugs

Se você encontrou um comportamento inesperado ou erro:

1. Verifique na aba [Issues](https://github.com/joaooncode/envguard/issues) se o problema já foi reportado.
2. Caso não exista, abra uma nova issue com:

- Uma descrição clara do problema.
- Passos para reproduzir o erro.
- Sistema operacional e versão do Go/envguard.
- Logs ou mensagens de erro relevantes (certifique-se de **não incluir credenciais reais** nos logs!).

### Sugerindo Melhorias

Ideias para novas funcionalidades são muito bem-vindas!

1. Abra uma issue detalhando a motivação da sugestão e como ela beneficiaria os usuários.
2. Descreva cenários de uso ou propostas de interface na CLI.

### Contribuindo com Código

Se deseja resolver uma issue ou adicionar uma funcionalidade:

1. Comente na issue em que deseja trabalhar para alinharmos e evitarmos retrabalho com outros contribuidores.
2. Siga as diretrizes de desenvolvimento abaixo.

---

## Ambiente de Desenvolvimento

### Pré-requisitos

- [Go](https://go.dev/dl/) versão **1.22** ou superior.
- [Git](https://git-scm.com/) instalado.

### Configurando o Repositório Localmente

```bash
# 1. Faça um Fork do repositório no GitHub para sua conta

# 2. Clone o seu fork
git clone https://github.com/<seu-usuario>/envguard.git
cd envguard

# 3. Adicione o repositório upstream oficial
git remote add upstream https://github.com/joaooncode/envguard.git

# 4. Baixe as dependências do módulo Go
go mod download
```

### Comandos Úteis

```bash
# Executar a CLI em desenvolvimento
go run ./cmd/envguard scan

# Executar a suíte de testes automatizados
go test ./...

# Executar testes com cobertura
go test -cover ./...

# Compilar o binário localmente
go build -o envguard ./cmd/envguard
```

---

## Padrões de Código e Qualidade

Para manter a consistência em todo o projeto:

1. **Formatação Oficial do Go:**

- Execute sempre o `gofmt` antes de comitar:
  ```bash
  gofmt -s -w .
  ```

2. **EditorConfig & Prettier:**

- O projeto possui [`.editorconfig`](.editorconfig) e [`.prettierrc`](.prettierrc) configurados.
- Mantenha **Tabs** para arquivos Go e **2 espaços** para Markdown, YAML e JSON.

3. **Segurança em Primeiro Lugar:**

- **Regra fundamental:** O `envguard` nunca deve imprimir nem expor valores de variáveis de ambiente ou segredos
  detectados.

---

## Convenção de Commits (Conventional Commits)

Utilizamos o padrão **[Conventional Commits](https://www.conventionalcommits.org/)** para manter o histórico claro e
possibilitar a automação de releases.

Formato: `<tipo>(<escopo opcional>): <descrição no imperativo>`

### Tipos permitidos:

- `feat`: Nova funcionalidade.
- `fix`: Correção de bug.
- `refactor`: Refatoração sem alteração de comportamento externo.
- `perf`: Melhoria de desempenho.
- `docs`: Alterações ou adições na documentação.
- `style`: Formatação, linting ou imports (sem alteração de lógica).
- `test`: Criação ou atualização de testes automatizados.
- `chore`: Atualização de dependências, builds ou configurações gerais.

### Exemplos:

```bash
git commit -m "feat(scanner): implementar detecção recursiva de arquivos .env"
git commit -m "fix(git): tratar repositórios sem commits iniciais"
git commit -m "docs: adicionar guia de instalação no README"
```

---

## Fluxo de Pull Requests

1. **Mantenha sua branch sincronizada com a upstream:**
   ```bash
   git checkout main
   git fetch upstream
   git merge upstream/main
   ```
2. **Crie uma branch local descritiva para a sua tarefa:**
   ```bash
   git checkout -b feat/minha-nova-feature
   ```
3. **Implemente as alterações com testes:**

- Adicione testes unitários para o código novo/corrigido.
- Certifique-se de que `go test ./...` está passando.

4. **Envie para o seu fork:**
   ```bash
   git push origin feat/minha-nova-feature
   ```
5. **Abra o Pull Request:**

- Abra o PR apontando para a branch `main` do repositório oficial.
- Preencha os campos do nosso **Pull Request Template** (`.github/pull_request_template.md`) com clareza.
- Marque a checklist pré-envio.

---

Muito obrigado por ajudar a tornar o **envguard** cada vez melhor e mais seguro!
