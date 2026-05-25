export const blogTopics = [
  {
    slug: 'architecture-notes',
    name: 'Architecture Notes',
    description:
      'Short field notes on service boundaries, generated contracts, and the tradeoffs that shape maintainable systems.',
    detail:
      'Placeholder detail copy for future architecture writing. This space will hold context, diagrams, tradeoffs, and practical notes from backend and shared contract work.',
  },
  {
    slug: 'frontend-journal',
    name: 'Frontend Journal',
    description:
      'Practical writing on Nuxt interfaces, component structure, and the small details that make product screens feel composed.',
    detail:
      'Placeholder detail copy for future frontend writing. This page will collect interface decisions, component notes, layout patterns, and delivery notes from product UI work.',
  },
  {
    slug: 'infrastructure-logs',
    name: 'Infrastructure Logs',
    description:
      'A working record of Cloud Run, Terraform, release paths, and the operational choices behind production software.',
    detail:
      'Placeholder detail copy for future infrastructure writing. This page will outline deployment paths, environment choices, operational notes, and reliability decisions.',
  },
  {
    slug: 'engineering-essays',
    name: 'Engineering Essays',
    description:
      'Clear reflections on debugging, code review, team habits, and the discipline of building software with fewer surprises.',
    detail:
      'Placeholder detail copy for future engineering essays. This page will hold longer notes on engineering judgment, collaboration, debugging, and review habits.',
  },
]

export type BlogTopic = (typeof blogTopics)[number]
