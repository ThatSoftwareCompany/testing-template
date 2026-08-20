# Go API Template Implementation Plan

## Foundation

Build a runnable Go 1.26 modular monolith with PostgreSQL enabled by default and an explicit `DATABASE_ENABLED=false` mode for local and test execution. The foundation includes validated environment configuration, structured logging, correlation IDs, security middleware, liveness/readiness endpoints, graceful shutdown, PostgreSQL pooling, SQL migrations, and safe HTTP error persistence.

## Documentation and operations

Keep OpenAPI documentation separate from code and split it by module. Add development and production Docker targets, a PostgreSQL Compose service, GitHub Actions quality checks, unit/HTTP/integration tests, a template manifest, and an idempotent setup script.

The template lifecycle foundation records generated-project provenance, validates the manifest contract, verifies required template files in CI, detects version tags, checks compatibility, and prepares automatic derived-repository update PRs. Breaking-change records and application-specific conflict resolution remain manual.

## Future phases

- Authentication: Argon2id passwords, Ed25519/EdDSA JWTs, approximately 15-minute access tokens, 30-day rotating/revocable refresh tokens, HttpOnly cookies, environment-specific Secure and SameSite policy, CSRF protection, authentication/authorization middleware, and securely managed Ed25519 keys.
- Internal error access: expose `GET /api/v1/internal/errors?endpoint=<path>` only after authentication and internal permissions exist.
- Template updates: maintain version tags, review generated PRs, record incompatible changes, and resolve application-specific conflicts manually.
- Dependency and image security: add Dependabot or Renovate, `govulncheck`, Docker image scanning, and strict `go.sum` freshness checks.
- Repository administration: after review and merge, mark both backend and frontend repositories as GitHub Template Repositories.

## Explicit non-goals

Google OAuth, frontend implementation, authentication, and authorization are outside this delivery. The repository license is Apache-2.0.
