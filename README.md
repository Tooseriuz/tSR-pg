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
pnpm nx run setup:deps
pnpm nx run web:dev
pnpm nx run api:dev
pnpm nx run mock:up
pnpm nx run mock:down
pnpm nx run shared-openapi:generate
pnpm nx run infra:fmt
```

## Local mock services

Start PostgreSQL and fake-gcs-server, then run local Go-tool Tern migrations after PostgreSQL is healthy:

```bash
pnpm nx run mock:up
```

- PostgreSQL: `localhost:5432`, database/user/password `tsr_pg`
- fake-gcs-server: `http://localhost:4443`

Run migrations again without restarting containers:

```bash
pnpm nx run mock:migrate
```

Stop the services:

```bash
pnpm nx run mock:down
```

## Contract-first reminder

API changes should start in `shared/openapi/openapi.yaml`, then flow into generated frontend/backend DTO locations.
