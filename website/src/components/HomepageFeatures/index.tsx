import type { ReactNode } from 'react';
import clsx from 'clsx';
import Heading from '@theme/Heading';
import styles from './styles.module.css';

type PillarItem = {
  tag: string;
  title: string;
  description: ReactNode;
};

const PillarList: PillarItem[] = [
  {
    tag: 'Detecção Git-Native',
    title: 'Consciência da Árvore Git',
    description: (
      <>
        Diferencia com precisão se arquivos de ambiente estão comitados no histórico (
        <em>tracked</em>), preparados para o próximo commit (<em>staged</em>) ou desprotegidos
        contra regras do <code>.gitignore</code>.
      </>
    ),
  },
  {
    tag: 'Segurança & Velocidade',
    title: '100% Local & Determinístico',
    description: (
      <>
        Desenvolvido em Go nativo, executa instantaneamente em pipelines e máquinas locais sem
        qualquer telemetria. Valores de variáveis e credenciais nunca são impressos em terminal ou
        logs.
      </>
    ),
  },
  {
    tag: 'Automação Contínua',
    title: 'Pre-Commit & Autocorreção',
    description: (
      <>
        Instale pre-commit hooks nativos com <code>envguard hook install</code> e execute correções
        automáticas no
        <code>.gitignore</code> e no cache do Git através de <code>envguard fix</code>.
      </>
    ),
  },
];

function PillarCard({ tag, title, description }: PillarItem) {
  return (
    <div className={clsx('col col--4', styles.pillarCol)}>
      <div className={styles.pillarCard}>
        <span className={styles.pillarTag}>{tag}</span>
        <Heading as="h3" className={styles.pillarTitle}>
          {title}
        </Heading>
        <p className={styles.pillarDescription}>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures(): ReactNode {
  return (
    <section className={styles.featuresSection}>
      <div className="container">
        <div className={styles.sectionHeader}>
          <Heading as="h2" className={styles.sectionTitle}>
            Construído para engenharia e segurança
          </Heading>
          <p className={styles.sectionSubtitle}>
            Prevenção ativa contra vazamento de segredos sem adicionar atrito ou dependências
            lentas.
          </p>
        </div>
        <div className="row">
          {PillarList.map((props, idx) => (
            <PillarCard key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
