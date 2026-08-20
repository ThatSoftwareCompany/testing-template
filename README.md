# {{PROJECT_NAME}}

Reusable Go API foundation for That Software Company. The generated application name is `{{APP_NAME}}`.

## Requirements

- Go 1.26.0 or newer.
- PostgreSQL 16 or newer for database-backed execution.
- Docker and Docker Compose for the container workflow.

## Local execution

PostgreSQL is enabled by default and `DATABASE_URL` is required in that mode. Copy `.env.example` to a local `.env` outside version control, set the database values, and export them before starting:

```bash
set -a
source .env
set +a
go run ./cmd/api
```

To run the API without PostgreSQL:

```bash
DATABASE_ENABLED=false go run ./cmd/api
```

The application validates all environment variables at startup. `APP_ENV` must be `development`, `test`, or `production`.

## Docker Compose

Set the local-only `POSTGRES_USER`, `POSTGRES_PASSWORD`, and `POSTGRES_DB` variables in your shell or a local `.env` file, then run:

```bash
docker compose up --build
```

Compose starts PostgreSQL with a healthcheck, runs development migrations on API startup, and stores PostgreSQL data in a named local volume. The production image is built separately:

```bash
docker build --target production -t example-api:local .
```

The production container runs as a non-root user. Production migrations must be executed explicitly with the `migrate` binary or `go run ./cmd/migrate`; automatic startup migrations are rejected in `production`.

## Migrations

```bash
go run ./cmd/migrate -command up
go run ./cmd/migrate -command down -steps 1
go run ./cmd/migrate -command version
```

The migration directory uses normal `.up.sql` and `.down.sql` files. The current foundation creates `error_events`.

## Operational endpoints

- `GET /__ping` checks only that the process is alive. It never queries PostgreSQL.
- `GET /api/v1/health` checks readiness. It reports `database: "disabled"`, `"up"`, or `"down"`; the latter returns HTTP 503 without exposing failure details.

Every response includes a validated `X-Correlation-ID`. Error responses include the same value in the JSON error object.

## Architecture

This is a modular monolith with a simple MVC flow:

```text
routes -> controller -> service -> repository/client
```

- `internal/modules/health` owns health transport and readiness rules.
- `internal/modules/errors` prepares the future internal error use case.
- `internal/platform` owns configuration, HTTP, logging, PostgreSQL, migrations, and safe error storage.
- `repository` is reserved for persistence.
- `client` is reserved for external APIs.
- Public routes are versioned with `/api/v1`.
- OpenAPI files live in `docs/openapi/` and are separate from controllers.

## Security

The foundation emits structured JSON logs through `log/slog`, security headers, an explicit CORS allowlist, and correlation IDs. It never logs request bodies, authorization headers, cookies, passwords, tokens, or secrets. Persisted HTTP 5xx events contain only the safe fields documented by the `error_events` migration.

In production, frontend and backend should be served under the same public origin. CORS credentials are enabled only for explicitly allowed origins; `*` is never accepted.

The internal error listing route is intentionally not registered until authentication and authorization exist.

## Tests and quality checks

```bash
go mod tidy
go mod verify
test -z "$(gofmt -l .)"
go vet ./...
go test ./...
go build ./cmd/api
go build ./cmd/migrate
```

Integration tests require PostgreSQL and use the `integration` build tag:

```bash
TEST_DATABASE_URL='postgres://USER:PASSWORD@localhost:5432/DB?sslmode=disable' go test -tags=integration ./...
```

## Setup script

The Bash setup script safely configures a generated project without arbitrary overwrites:

```bash
./scripts/setup.sh \
  --project-name "Example API" \
  --module-path "github.com/example/example-api" \
  --app-name "example-api" \
  --environment development \
  --database-enabled true \
  --architecture modular-mvc \
  --generated-from "ThatSoftwareCompany/example-api"
```

The setup script records the source commit in `template_commit` when Git metadata is available. It records `generated_from` from `--generated-from`, defaulting to the generated module path, and preserves both values on subsequent idempotent runs.

Validate the template manifest and required files with:

```bash
./scripts/validate-template.sh
```

## Template metadata and updates

`.template/manifest.json` records the source repository, template version, template commit, generated origin, compatibility, dependencies, and update policy. The update automation detects new template versions, opens PRs in derived repositories, enforces compatibility, and leaves breaking-change records and application-specific conflicts for manual review.

The generated repository also includes a scheduled and manually dispatchable template-update workflow. It looks for `vMAJOR.MINOR.PATCH` tags, applies a three-way patch from the recorded `template_commit`, checks Go and PostgreSQL compatibility, records new provenance, and opens a pull request. It never merges automatically. The repository owner must allow GitHub Actions to create pull requests and review generated changes manually.

The template maintainer must publish version tags such as `v0.1.0` before derived repositories can detect releases. The initial release tag should point to the merged template commit.

## Planned phases

- Authentication with Argon2id, Ed25519/EdDSA JWTs, approximately 15-minute access tokens, 30-day rotating/revocable refresh tokens, HttpOnly cookies, environment-specific Secure and SameSite policies, CSRF protection, authentication/authorization middleware, and securely managed Ed25519 keys.
- Authenticated access to `/api/v1/internal/errors?endpoint=<path>`.
- Google OAuth integration.
- Dependabot or Renovate, `govulncheck`, Docker image scanning, and stricter `go.sum` freshness checks.

After review and merge, the backend and frontend repositories must be marked as GitHub Template Repositories from `Settings -> General -> Template repository`. This is a post-merge checklist item, not an automated repository mutation.

## License

This template is distributed under the [Apache License 2.0](LICENSE).
