# AGENTS.md

## Repository Overview

This repository is a monorepo structured with Nx.

### Stack

#### Monorepo
- Nx
- Node.js managed with nvm

#### Frontend
- Vue / Nuxt
- Tailwind CSS
- shadcn/ui

#### Backend
- Go
- Gin
- Hexagonal Architecture

#### Infrastructure
- Terraform
  - Cloud Run for API deployment
  - Artifact Registry for Docker images
  - OIDC authentication for GitHub Actions
  - Secret Manager

#### Shared
- OpenAPI is the source of truth for API contracts
- Shared DTO/types must be generated from OpenAPI
- Frontend and backend must use generated types instead of defining duplicate DTOs manually
- Use `oapi-codegen` for backend/frontend type generation

---

# Development Rules

## API Changes

If a change affects API request/response structures or endpoints:

1. Update `openapi.yaml` first
2. Regenerate shared types
3. Use generated types in both frontend and backend
4. Do not manually redefine generated DTOs

---

## Testing Rules

Write only the tests explicitly requested in the prompt.

Do not add speculative tests or unnecessary test coverage.

### Backend
- Unit tests:
  - Service/business logic layer
- Integration tests:
  - Handler/API layer end-to-end behavior

---

## Validation Before Finishing

Before completing any task:

1. Run lint
2. Run tests
3. Run build

Do not run build on backend go

Ensure all commands succeed before considering the task complete.

Do not try to fix it without permission, report to user instead.

Do not run development servers unless explicitly requested.

---

# Architecture Rules

## Backend
- Follow hexagonal architecture
- Keep business logic independent from transport/framework layers
- Avoid leaking Gin-specific logic into domain/services

## Shared Types
- Generated types are the single source of truth
- Avoid duplicate type definitions across apps

---

# Frontend Design Rules

For UI/design-related tasks:
- Use Tailwind CSS
- Use shadcn/ui components where appropriate
- Follow styles and patterns from:
  - `./agents/skills/design-taste-frontend`
- create new component instead of nested div within one file

Prefer:
- consistent spacing
- accessible UI
- responsive layouts
- reusable components

Avoid:
- inline styles
- unnecessary custom CSS
- duplicated UI patterns

---

# General Coding Guidelines

- Prefer simple and maintainable solutions
- Avoid premature abstraction
- Avoid introducing new dependencies unless necessary, or if it's better option to add new dependencies then ask user.
- Keep files focused and cohesive
- Follow existing project conventions before introducing new patterns