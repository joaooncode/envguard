---
sidebar_position: 3
title: Início Rápido
---

# Início Rápido (Quickstart)

Aprenda a instalar, executar sua primeira varredura e proteger seu repositório Git com o **`envguard`** em menos de 3 minutos.

---

## 1. Instalação Instantânea

Se você possui o Go instalado (1.22+):

```bash
go install github.com/joaooncode/envguard/cmd/envguard@latest
```

Ou no Windows com PowerShell:

```powershell
# Verifique a versão instalada
envguard version
```

> Para binários pré-compilados do Linux, macOS e Windows sem Go, consulte o [Guia de Instalação](./installation.md).

---

## 2. Execute sua Primeira Varredura

Navegue até a raiz do seu projeto com Git e execute o comando `scan`:

```bash
envguard scan
```

O `envguard` varrerá todos os arquivos e diretórios recursivamente, identificando arquivos de ambiente e cruzando o status com o Git e o `.gitignore`.

### Exemplo de Saída:

```text
[envguard] Iniciando varredura no repositório...

✗ .env                   CRITICAL  Arquivo rastreado no Git (commitado)
⚠ .env.production        WARNING   Não ignorado no .gitignore
✓ .env.example           INFO      Template padrão identificado

Status: 2 problema(s) encontrado(s).
```

---

## 3. Gerar Arquivo de Configuração

Inicialize um arquivo `.envguard.yaml` para customizar permissões e regras do projeto:

```bash
envguard init
```

Isso gerará o arquivo de configuração `.envguard.yaml` na raiz do repositório:

```yaml
version: 1
patterns:
  - '.env*'
ignore:
  - '.env.example'
  - '.env.template'
fail_on: 'warning' # Interrompe a execução com warning ou critical
```

---

## 4. Corrigir Problemas Automaticamente

Se o `envguard` detectar arquivos `.env` soltos que precisam ser adicionados ao `.gitignore` ou removidos do tracking do Git:

```bash
envguard fix
```

---

## 5. Instalar Git Hook Automático (Pre-Commit)

Garanta que nenhum arquivo de ambiente seja commitado por acidente instalando o hook nativo:

```bash
envguard hook install
```

Pronto! Toda vez que alguém executar `git commit`, o `envguard check` será executado automaticamente antes do commit ser gravado.

---

## Próximos Passos

- Explore todos os [Comandos da CLI](./commands/scan.md).
- Entenda a classificação dos [Níveis de Severidade](./severity-levels.md).
- Configure a integração com [GitHub Actions e CI/CD](./cicd-integration.md).
