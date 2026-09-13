---
slug: lancamento-envguard-v030
title: Lançamento do envguard v0.3.0 - Secret Scanning por Conteúdo
authors: [joaooncode]
tags: [release, security, cli, golang, secret-scanning]
---

É com muito orgulho que anunciamos o lançamento do **`envguard` v0.3.0**! 🎉

Esta versão vai além do nome e do status Git dos arquivos de ambiente: agora o `envguard` também inspeciona o **conteúdo** desses arquivos em busca de segredos reais, sem nunca expor os valores encontrados.

<!-- truncate -->

## Principais Novidades da v0.3.0

### 1. Detecção de Padrões de Provedores Conhecidos

Um novo motor de assinaturas regex identifica credenciais reais dentro dos arquivos de ambiente escaneados:

- Chaves de acesso **AWS** (`AKIA...`)
- Chaves secretas **Stripe** (`sk_live_`/`sk_test_`)
- Tokens de acesso pessoal **GitHub** (`ghp_`, `gho_`, etc.)
- Chaves privadas **PEM** (RSA, EC, OpenSSH, DSA)
- Tokens genéricos **Bearer**

### 2. Detecção por Entropia de Shannon (opt-in)

Além das assinaturas conhecidas, o `envguard` pode calcular a entropia de cada valor `CHAVE=VALOR` para sinalizar segredos aleatórios que não correspondem a nenhum provedor catalogado — chaves criptográficas, hashes e tokens gerados aleatoriamente. Por gerar mais falsos positivos, esse modo é **desativado por padrão** e habilitado via `.envguard.yaml`.

### 3. Severidade Independente por Secret Match

Cada segredo encontrado (`Secret Match`) recebe sua própria severidade, calculada com piso `CRITICAL` (arquivo _tracked_/_staged_) ou `HIGH` (arquivo apenas local), e respeitando qualquer `severity_overrides` já configurado como teto. A severidade efetiva de um arquivo passa a considerar o maior valor entre o `Finding` e seus `Secret Matches`.

### 4. Controle Fino via `.envguard.yaml`

```yaml
detector:
  entropy_scan: false
  secret_providers: [] # vazio = todos: aws, stripe, github, pem, bearer-token
  secret_ignore: [] # ex: ["TEST_*"] para ignorar chaves de fixtures
```

### 5. Seguro por Design

Como em toda a ferramenta, **o valor do segredo nunca é impresso** em nenhuma saída (terminal ou JSON) — apenas linha, nome da chave, método de detecção e provedor.

---

## Como Atualizar

```bash
go install github.com/joaooncode/envguard/cmd/envguard@latest
```

Verifique a versão instalada:

```bash
envguard version
# envguard v0.3.0
```

Confira a [documentação de configuração](/docs/configuration#secret-scanning-por-conteúdo) para habilitar a detecção por entropia e personalizar os provedores monitorados!
