import type { ReactNode } from 'react';
import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import HomepageFeatures from '@site/src/components/HomepageFeatures';
import Heading from '@theme/Heading';

import styles from './index.module.css';

function TerminalPreview() {
  return (
    <div className={styles.terminalWrapper}>
      <div className={styles.terminalHeader}>
        <div className={styles.terminalDots}>
          <span className={clsx(styles.dot, styles.dotRed)} />
          <span className={clsx(styles.dot, styles.dotYellow)} />
          <span className={clsx(styles.dot, styles.dotGreen)} />
        </div>
        <span className={styles.terminalTitle}>envguard terminal • zsh / powershell</span>
      </div>
      <div className={styles.terminalBody}>
        <div>
          <span className={styles.termPrompt}>$ </span>
          <span className={styles.termCmd}>envguard scan</span>
        </div>
        <div className={styles.termComment}>
          [envguard] Varrendo repositório em busca de arquivos de ambiente sensíveis...
        </div>
        <br />
        <div>
          <span className={styles.termCrit}>✗ .env</span>{' '}
          &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;{' '}
          <span className={styles.termCrit}>CRITICAL</span> &nbsp; Arquivo rastreado no Git
          (commitado)
        </div>
        <div>
          <span className={styles.termWarn}>⚠ .env.production</span> &nbsp;&nbsp;&nbsp;{' '}
          <span className={styles.termWarn}>WARNING</span> &nbsp;&nbsp; Não ignorado no .gitignore
        </div>
        <div>
          <span className={styles.termInfo}>✓ .env.example</span>{' '}
          &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; <span className={styles.termInfo}>INFO</span>{' '}
          &nbsp;&nbsp;&nbsp;&nbsp;&nbsp; Template padrão identificado
        </div>
        <br />
        <div>
          <span className={styles.termCrit}>[ERRO]</span> Encontrado 1 problema CRITICAL e 1
          WARNING.
        </div>
        <div>
          <span className={styles.termComment}>
            Dica: execute `envguard fix` para adicionar ao .gitignore e remover do cache Git.
          </span>
        </div>
      </div>
    </div>
  );
}

function HomepageHeader() {
  const { siteConfig } = useDocusaurusContext();
  return (
    <header className={styles.heroBanner}>
      <div className="container">
        <div className={styles.heroContent}>
          <div className={styles.versionBadge}>
            <span className={styles.badgeGlowDot} />
            <span>v0.2.0 • Proteção Git-Aware para Segredos</span>
          </div>

          <Heading as="h1" className={styles.heroTitle}>
            <span className={styles.gradientText}>Segurança de .env</span>
            <br />
            com Consciência do Git
          </Heading>

          <p className="hero__subtitle">
            Evite que arquivos <code>.env</code> e credenciais de ambiente vazem ou cheguem ao
            histórico de versionamento do Git.
          </p>

          <div
            style={{
              margin: '1.25rem 0 2rem',
              display: 'flex',
              justifyContent: 'center',
              gap: '0.6rem',
              flexWrap: 'wrap',
            }}
          >
            <img
              alt="Go Version"
              src="https://img.shields.io/badge/Go-1.22%2B-6A42C2?style=flat-square&logo=go&logoColor=white"
            />
            <img
              alt="License: MIT"
              src="https://img.shields.io/badge/License-MIT-8B5DFF?style=flat-square"
            />
            <img
              alt="100% Offline"
              src="https://img.shields.io/badge/Privacy-100%25%20Offline-green?style=flat-square"
            />
          </div>

          <div className={styles.buttons}>
            <Link className="button button--primary button--lg" to="/docs/quickstart">
              ⚡ Início Rápido (3 min)
            </Link>
            <Link className="button button--secondary button--lg" to="/docs/intro">
              📖 Explorar Documentação
            </Link>
          </div>

          <TerminalPreview />
        </div>
      </div>
    </header>
  );
}

export default function Home(): ReactNode {
  return (
    <Layout
      title="Documentação Oficial"
      description="Proteja seus arquivos .env e variáveis de ambiente sensíveis antes que cheguem ao Git."
    >
      <HomepageHeader />
      <main>
        <HomepageFeatures />
      </main>
    </Layout>
  );
}
