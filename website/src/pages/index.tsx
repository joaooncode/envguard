import { useState, type ReactNode } from 'react';
import Link from '@docusaurus/Link';
import Layout from '@theme/Layout';
import HomepageFeatures from '@site/src/components/HomepageFeatures';
import Heading from '@theme/Heading';

import styles from './index.module.css';

function GitInspectionShowcase() {
  return (
    <div className={styles.inspectionCard}>
      <div className={styles.inspectionHeader}>
        <span className={styles.headerTitle}>envguard scan --format text</span>
        <span className={styles.headerStatus}>Git Inspection</span>
      </div>
      <div className={styles.inspectionBody}>
        <div className={styles.fileRow}>
          <div>
            <span className={styles.fileName}>.env</span>
            <span className={styles.fileContext}>• Tracked in Git commit history</span>
          </div>
          <span className={styles.tagCritical}>CRITICAL</span>
        </div>
        <div className={styles.fileRow}>
          <div>
            <span className={styles.fileName}>.env.production</span>
            <span className={styles.fileContext}>• Not ignored in .gitignore</span>
          </div>
          <span className={styles.tagWarning}>WARNING</span>
        </div>
        <div className={styles.fileRow}>
          <div>
            <span className={styles.fileName}>.env.example</span>
            <span className={styles.fileContext}>• Allowlisted sample template</span>
          </div>
          <span className={styles.tagInfo}>SAFE (INFO)</span>
        </div>
      </div>
    </div>
  );
}

function InstallCommand() {
  const [copied, setCopied] = useState(false);
  const cmd = 'go install github.com/joaooncode/envguard/cmd/envguard@latest';

  const handleCopy = () => {
    navigator.clipboard.writeText(cmd);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className={styles.installPill} onClick={handleCopy} title="Clique para copiar comando">
      <span>$ {cmd}</span>
      <span className={styles.copyIcon}>{copied ? 'COPIADO ✓' : 'COPIAR'}</span>
    </div>
  );
}

function HomepageHeader() {
  return (
    <header className={styles.heroBanner}>
      <div className="container">
        <div className={styles.heroContent}>
          <div className={styles.versionBadge}>
            <span>envguard v0.2.0</span>
          </div>

          <Heading as="h1" className={styles.heroTitle}>
            Proteção de segredos .env com inteligência do Git
          </Heading>

          <p className={styles.heroSubtitle}>
            Detecte, alerte e impeça que arquivos de variáveis de ambiente e credenciais cheguem ao
            histórico do Git.
          </p>

          <InstallCommand />

          <div className={styles.buttons}>
            <Link className="button button--primary" to="/docs/quickstart">
              Início Rápido
            </Link>
            <Link className="button button--secondary" to="/docs/intro">
              Documentação
            </Link>
          </div>

          <GitInspectionShowcase />
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
