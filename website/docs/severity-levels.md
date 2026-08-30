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

## Próximos Recursos no Roadmap

Na versão `v0.2.0`, será possível customizar padrões extras de busca, ignorar pastas específicas e definir regras corporativas através do arquivo `.envguard.yaml`.
