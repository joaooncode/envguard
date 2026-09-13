---
sidebar_position: 4
title: Níveis de Severidade
---

# Níveis de Severidade

O `envguard` classifica cada arquivo de ambiente detectado com base no risco de exposição e no seu estado dentro da árvore de trabalho do Git.

---

## Tabela de Severidade

| Nível               | Status Git              | Descrição do Risco                                                                                                  | Ação Recomendada                                                                                                                 |
| :------------------ | :---------------------- | :------------------------------------------------------------------------------------------------------------------ | :------------------------------------------------------------------------------------------------------------------------------- |
| **`CRITICAL`**      | _Tracked_               | Arquivo de ambiente está atualmente **rastreado no histórico do Git**. O risco de vazamento de segredos é imediato. | Desrastrear o arquivo (`git rm --cached <arquivo>`), adicionar ao `.gitignore` e **rotacionar as credenciais imediatamente**.    |
| **`HIGH`**          | _Staged_                | Arquivo de ambiente foi adicionado via `git add` e está pronto para ser comitado.                                   | Retirar da stage (`git restore --staged <arquivo>` ou `git reset HEAD <arquivo>`) e garantir que a regra esteja no `.gitignore`. |
| **`WARNING`**       | _Unignored / Untracked_ | Arquivo de ambiente existe localmente, mas **não está coberto** por nenhuma regra do `.gitignore`.                  | Adicionar o padrão correspondente (ex: `.env*`) ao seu arquivo `.gitignore`.                                                     |
| **`INFO` / `SAFE`** | _Ignored / Template_    | Arquivo está corretamente ignorado ou é reconhecido como um template seguro (`.env.example`, `.env.sample`).        | Nenhuma ação necessária.                                                                                                         |

---

## Padrões Monitorados & Exceções

### Arquivos Monitorados por Padrão

- `.env`
- `.env.local`, `.env.development`, `.env.production`, `.env.staging`, `.env.test`
- Arquivos com sufixo `.env` (ex: `app.env`, `database.env`)

### Templates Seguros Permitidos

Os seguintes arquivos são reconhecidos por padrão como modelos públicos sem segredos:

- `.env.example`
- `.env.sample`
- `.env.template`
- `.env.dist`

---

## Severidade dos Secret Matches

A partir da `v0.3.0`, cada arquivo de ambiente também pode conter um ou mais **Secret Matches** — segredos reais detectados no conteúdo do arquivo pelo [Secret Scanner](./configuration.md#secret-scanning-por-conteúdo). A severidade de um Secret Match é calculada de forma **independente** da severidade do arquivo (`Finding`) que o contém:

- **Piso `CRITICAL`:** se o arquivo estiver *tracked* ou *staged* no Git.
- **Piso `HIGH`:** caso contrário (arquivo apenas presente localmente).
- **Teto:** se houver um `severity_overrides` aplicável ao arquivo, ele também limita a severidade do Secret Match (por exemplo, um `.env.test` rebaixado para `info` nunca gera um Secret Match acima de `info`).

Como um Secret Match pode ser mais severo que o próprio `Finding` (por exemplo, um segredo real dentro de um arquivo corretamente ignorado, com severidade `INFO`), o resumo final do scan e o código de saída sempre consideram a **severidade efetiva** — o maior valor entre o `Finding` e todos os seus `Secret Matches`.

---

## Próximos Recursos no Roadmap

Na versão `v1.0.0`, está prevista a GitHub Action oficial do `envguard` e pacotes para gerenciadores como Homebrew, Scoop, WinGet e AUR.
