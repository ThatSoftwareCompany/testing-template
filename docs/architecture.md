# Architecture

The template is a modular monolith. Each business module owns its HTTP controller, service rules, DTOs, persistence repository, and external clients when those responsibilities exist.

```text
HTTP middleware -> module controller -> module service -> repository/client
```

## Composition

- `cmd/api` owns startup, signal handling, dependency composition, and shutdown.
- `cmd/migrate` owns explicit schema migration commands.
- `internal/platform/config` validates environment configuration.
- `internal/platform/db` owns the `pgxpool` lifecycle.
- `internal/platform/httpserver` owns shared middleware and response contracts.
- `internal/platform/logging` creates the JSON `slog` logger.
- `internal/platform/errstore` owns safe PostgreSQL/no-op error persistence.
- `internal/platform/migrate` wraps `golang-migrate` SQL migrations.
- `internal/modules/health` owns liveness/readiness transport and rules.
- `internal/modules/errors` prepares the future authenticated error listing use case.

## Database modes

PostgreSQL is enabled by default and required through `DATABASE_URL`. Set `DATABASE_ENABLED=false` when the API must run without PostgreSQL; in that mode health reports `database: "disabled"` and error persistence is no-op.

## Public and internal boundaries

The public API is versioned under `/api/v1`. `/__ping` intentionally sits outside the API version because it is an infrastructure liveness probe. The internal error listing route is documented but not registered until authentication and authorization exist.

## Future authentication boundary

The authentication phase will add Argon2id password hashing, Ed25519/EdDSA JWTs, short-lived access tokens, rotating/revocable refresh tokens in HttpOnly cookies, environment-specific Secure/SameSite behavior, CSRF protection, and authentication/authorization middleware. Ed25519 keys must be generated and supplied securely outside the repository.
