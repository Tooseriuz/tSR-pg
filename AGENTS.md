1. This is repo designed to be monorepo that include
  - monorepo
    - nx
    - nvm
  - frontend
    - Vue/Nuxt, tailwind-css, shadcn/ui
  - backend
    - go with gin
    - hexagonal architecture
  - infrastructure
    - terraform
      - cloud run for api
      - artifact registry for docker image
      - oidc for github action authentication
      - secret manager
  - shared
    - openapi for type contracting between frontend and backend
    - type generation for frontend and backend usingo opi-codegen
    - backend and frontend should use generated types instead of create their own type
2. IMPORTANT reminders before editing
  - If related to API endpoint changes, do the openapi.yaml first then generate shared DTO types for both frontend and backend
  - Add unit tests on every services level for business logic
  - Add integretion tests on every handlers level for end to end test
  - Run lint then test then build to verify app successful
3. In design work, make use of tailwindcss and shadcn/ui