import type { ReactNode } from 'react';
import clsx from 'clsx';
import Heading from '@theme/Heading';
import styles from './styles.module.css';

type FeatureItem = {
  icon: string;
  badge?: string;
  title: string;
  description: ReactNode;
};

const FeatureList: FeatureItem[] = [
  {
    icon: '🛡️',
    badge: 'Detecção Ativa',
    title: 'Git-Aware & Inteligente',
    description: (
      <>
        O <code>envguard</code> analisa a árvore de trabalho do Git e diferencia se arquivos <code>.env</code>{' '}
        estão comitados (<em>tracked</em>), preparados (<em>staged</em>) ou desprotegidos pelo <code>.gitignore</code>.
      </>
    ),
  },
  {
    icon: '⚡',
    badge: 'Zero Telemetria',
    title: 'Ultrarrápido & 100% Local',
    description: (
      <>
        Construído em Go nativo, executa instantaneamente e 100% offline. Nunca envia dados para servidores
        remotos e jamais expõe o conteúdo dos seus segredos.
      </>
    ),
  },
  {
    icon: '🤖',
    badge: 'Pre-Commit & CI',
    title: 'Pronto para CI/CD & Hooks',
    description: (
      <>
        Instale pre-commit hooks nativos com <code>envguard hook install</code> e utilize saídas JSON determinísticas
        para automação no GitHub Actions e pipelines.
      </>
    ),
  },
  {
    icon: '🛠️',
    badge: 'Autocorreção',
    title: 'Remediação com 1 Comando',
    description: (
      <>
        Com o comando <code>envguard fix</code>, você remove arquivos sensíveis do cache do Git e adiciona as regras
        necessárias no <code>.gitignore</code> automaticamente.
      </>
    ),
  },
  {
    icon: '📊',
    badge: 'Diagnóstico Claro',
    title: 'Níveis de Severidade',
    description: (
      <>
        Classificação precisa dos alertas em <strong>CRITICAL</strong>, <strong>WARNING</strong> e <strong>INFO</strong>,
        com suporte a flags como <code>--fail-on</code> para bloquear builds arriscados.
      </>
    ),
  },
  {
    icon: '⚙️',
    badge: 'Flexível',
    title: 'Configuração Simplificada',
    description: (
      <>
        Personalize padrões de busca, exceções de templates e regras customizadas através do arquivo central{' '}
        <code>.envguard.yaml</code> gerado via <code>envguard init</code>.
      </>
    ),
  },
];

function Feature({ icon, badge, title, description }: FeatureItem) {
  return (
    <div className={clsx('col col--4', styles.featureCol)}>
      <div className={styles.featureCard}>
        <div className={styles.featureHeader}>
          <span className={styles.featureIcon}>{icon}</span>
          {badge && <span className={styles.featureBadge}>{badge}</span>}
        </div>
        <Heading as="h3" className={styles.featureTitle}>
          {title}
        </Heading>
        <p className={styles.featureDescription}>{description}</p>
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
            Recursos projetados para proteger seu fluxo de trabalho
          </Heading>
          <p className={styles.sectionSubtitle}>
            Simplicidade, velocidade e segurança máxima sem atrapalhar a produtividade do time de engenharia.
          </p>
        </div>
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}

