import type {SidebarsConfig} from '@docusaurus/plugin-content-docs';

const sidebars: SidebarsConfig = {
  docsSidebar: [
    'intro',
    'quickstart',
    'architecture',
    {
      type: 'category',
      label: 'Concepts',
      items: ['concepts/workflows', 'concepts/jobs', 'concepts/workers'],
    },
    {
      type: 'category',
      label: 'Guides',
      items: [
        'guides/creating-workflows',
        'guides/submitting-jobs',
        'guides/configuring-retries',
      ],
    },
    {
      type: 'category',
      label: 'API Reference',
      items: ['api/rest-reference'],
    },
  ],
};

export default sidebars;
