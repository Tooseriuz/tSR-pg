# tSR-pg

Minimal placeholder monorepo scaffold for:

- `apps/web`: Nuxt frontend placeholder
- `apps/api`: Go + Gin backend placeholder
- `shared/openapi`: OpenAPI-first contract source
- `infra/terraform`: Terraform placeholders for GCP deployment

## Workspace

- Package manager: `pnpm`
- Task runner: `Nx`

Expected commands:

```bash
pnpm install
pnpm nx run web:dev
pnpm nx run api:dev
pnpm nx run shared-openapi:generate
pnpm nx run infra:fmt
```

## Contract-first reminder

API changes should start in `shared/openapi/openapi.yaml`, then flow into generated frontend/backend DTO locations.
