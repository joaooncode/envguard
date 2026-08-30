import type { ReactNode } from 'react';
import clsx from 'clsx';
import Heading from '@theme/Heading';
import styles from './styles.module.css';

type FeatureItem = {
  title: string;
  description: ReactNode;
};

const FeatureList: FeatureItem[] = [
  {
    title: 'Git-Aware & Inteligente',
    description: (
      <>
        O <code>envguard</code> analisa a árvore do Git e identifica se arquivos <code>.env</code>{' '}
        estão rastreados, preparados para commit (staged), ignorados ou desprotegidos.
      </>
    ),
  },
  {
    title: 'Rápido & 100% Local',
    description: (
      <>
        Construído em Go nativo, roda instantaneamente e 100% offline. Nunca envia dados para a
        internet e nunca expõe valores de variáveis.
      </>
    ),
  },
  {
    title: 'Pronto para CI/CD & Hooks',
    description: (
      <>
        Códigos de saída determinísticos e saída estruturada em JSON (<code>--format json</code>)
        para integração fácil com GitHub Actions e pre-commit hooks.
      </>
    ),
  },
];

function Feature({ title, description }: FeatureItem) {
  return (
    <div className={clsx('col col--4')}>
      <div className="text--center padding-horiz--md" style={{ paddingTop: '1.5rem' }}>
        <Heading as="h3">{title}</Heading>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures(): ReactNode {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
