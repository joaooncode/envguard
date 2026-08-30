import type { ReactNode } from 'react';
import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import HomepageFeatures from '@site/src/components/HomepageFeatures';
import Heading from '@theme/Heading';

import styles from './index.module.css';

function HomepageHeader() {
  const { siteConfig } = useDocusaurusContext();
  return (
    <header className={clsx('hero hero--primary', styles.heroBanner)}>
      <div className="container">
        <Heading as="h1" className="hero__title">
          {siteConfig.title}
        </Heading>
        <p className="hero__subtitle">{siteConfig.tagline}</p>

        <div
          style={{
            margin: '1.5rem 0',
            display: 'flex',
            justifyContent: 'center',
            gap: '0.5rem',
            flexWrap: 'wrap',
          }}
        >
          <img alt="Go Version" src="https://img.shields.io/badge/go-1.22%2B-blue.svg" />
          <img alt="License: MIT" src="https://img.shields.io/badge/License-MIT-yellow.svg" />
          <img alt="PRs Welcome" src="https://img.shields.io/badge/PRs-welcome-brightgreen.svg" />
        </div>

        <div
          className={styles.buttons}
          style={{ display: 'flex', gap: '1rem', justifyContent: 'center' }}
        >
          <Link className="button button--secondary button--lg" to="/docs/intro">
            Comece Agora
          </Link>
          <Link
            className="button button--outline button--secondary button--lg"
            to="/docs/commands/scan"
            style={{ color: 'white', borderColor: 'rgba(255,255,255,0.8)' }}
          >
            Ver Comandos CLI
          </Link>
        </div>
      </div>
    </header>
  );
}

export default function Home(): ReactNode {
  const { siteConfig } = useDocusaurusContext();
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
