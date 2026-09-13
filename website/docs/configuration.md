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

  # Habilita a heurística de entropia do Secret Scanner (opt-in, desativado por padrão)
  entropy_scan: false

  # Restringe o Secret Scanner a provedores específicos (vazio = todos os embutidos)
  secret_providers:
    - 'aws'
    - 'stripe'
    - 'github'
    - 'pem'
    - 'bearer-token'

  # Ignora Secret Matches cuja chave seja igual ou faça glob-match com estas entradas
  secret_ignore:
    - 'EXAMPLE_KEY'
    - 'TEST_*'
```

---

## Secret Scanning por Conteúdo

O `detector` também controla o **Secret Scanner**, componente que inspeciona o _conteúdo_ dos arquivos de ambiente encontrados (independente da checagem de nome/status Git):

| Chave              | Tipo       | Descrição                                                                                                                                                               | Padrão       |
| :----------------- | :--------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :----------- |
| `entropy_scan`     | `bool`     | Habilita a heurística de Shannon entropy para sinalizar valores aleatórios que não correspondam a nenhuma assinatura conhecida. Opt-in por gerar mais falsos positivos. | `false`      |
| `secret_providers` | `[]string` | Restringe as assinaturas de padrões a provedores específicos. Lista vazia habilita todos os embutidos: `aws`, `stripe`, `github`, `pem`, `bearer-token`.                | `[]` (todos) |
| `secret_ignore`    | `[]string` | Suprime Secret Matches cuja chave seja igual ou corresponda (glob) a uma destas entradas — útil para chaves de exemplo ou fixtures de teste.                            | `[]`         |

Cada Secret Match encontrado é reportado com número da linha, nome da chave e provedor (quando aplicável) — **o valor do segredo nunca é exibido em nenhuma saída**. Veja [Níveis de Severidade](./severity-levels.md) para entender como a severidade de um Secret Match é calculada.

---

## Gerando a Configuração Automaticamente

Para criar o arquivo inicial no seu projeto com todos os comentários explicativos:

```bash
envguard init
```
