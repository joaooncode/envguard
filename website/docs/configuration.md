---
sidebar_position: 4
title: Arquivo de Configuração
---

# Arquivo de Configuração

O `envguard` suporta personalização profunda por meio de um arquivo de configuração `.envguard.yaml` (ou `.envguard.yml`) localizado na raiz do projeto ou indicado via flag `--config`.

---

## Ordem de Descoberta

Quando um comando é executado, a configuração é carregada seguindo esta prioridade:

1. Caminho explícito passado pela flag `--config <arquivo>`.
2. Arquivo `.envguard.yaml` no diretório alvo.
3. Arquivo `.envguard.yml` no diretório alvo.
4. Padrões embutidos padrão do `envguard`.

---

## Estrutura do Arquivo `.envguard.yaml`

```yaml
version: '1'

scanner:
  # Diretórios ignorados durante a varredura recursiva
  ignore_dirs:
    - 'node_modules'
    - '.git'
    - 'vendor'
    - 'dist'
    - 'build'
    - '.idea'
    - '.vscode'

detector:
  # Padrões adicionais de arquivos a serem reconhecidos como variáveis de ambiente
  custom_patterns:
    - '*.env.vault'
    - '.env.release'

  # Padrões considerados seguros (não geram alerta)
  allowlist:
    - '.env.example'
    - '.env.sample'
    - '.env.template'
    - '.env.ci'

  # Sobrescrita explícita de níveis de severidade
  severity_overrides:
    - pattern: '.env.test'
      severity: 'warning' # info, warning, high, critical
    - pattern: '.env.sandbox'
      severity: 'info'
```

---

## Gerando a Configuração Automaticamente

Para criar o arquivo inicial no seu projeto com todos os comentários explicativos:

```bash
envguard init
```
