# Template maintenance

The manifest is the source of truth for the template identity, version, source repository, compatibility, and generated origin.

## Implemented lifecycle safeguards

- `scripts/setup.sh` resolves `template_commit` from the published source tag and never uses the generated repository's own commit.
- `scripts/setup.sh` records `generated_from` from the derived repository identifier or module path.
- `cmd/template -command validate` validates the manifest contract.
- `scripts/validate-template.sh` verifies required template files and blocks a real `.env` file.
- CI runs the lifecycle validation and writes build outputs outside the repository root.
- `scripts/template-update.sh` applies a three-way patch between recorded and target template commits.
- `.github/workflows/template-update.yml` detects version tags and opens derived-repository PRs with least-privilege write permissions.
- The workflow accepts an optional `TEMPLATE_UPDATE_TOKEN` secret for updates that modify `.github/workflows` files.
- Compatibility changes are rejected automatically; file deletions remain a manual migration.

## Update workflow

The derived-repository workflow performs these steps:

1. Detect the newest `vMAJOR.MINOR.PATCH` tag from `ThatSoftwareCompany/template-go-api`.
2. Compare the generated repository's recorded `template_commit` and compatibility fields.
3. Apply a three-way patch from the recorded commit to the tagged commit.
4. Refuse incompatible Go/PostgreSQL changes and template file deletions.
5. Record the new `template_version` and `template_commit`.
6. Create an update pull request in the derived repository.
7. Leave application-specific conflicts for manual resolution; it must never merge generated pull requests automatically.

If a generated repository contains its own Git commit in `template_commit`, the workflow resolves the source commit from the matching release tag and opens a provenance-repair pull request.

The workflow requires GitHub Actions to be allowed to create pull requests in the derived repository. It is skipped when running in the canonical template repository itself.

## File ownership and extension points

Template-managed files contain reusable runtime, security, CI, Docker, migration, and lifecycle behavior. Generated repositories should not modify them to add product functionality. This includes `cmd/api`, `internal/platform`, `internal/modules/health`, `internal/modules/errors`, `.github/workflows`, `scripts`, `.template`, and the root Docker, migration, and CI files.

Product-specific code belongs in new business modules under `internal/modules/<business-module>/`. Register those modules in `internal/app/routes.go`, which is an application-owned extension point intentionally preserved by `scripts/template-update.sh`. The template-provided `/__ping` and `/api/v1/health` routes are operational routes and remain active automatically; they do not need to be copied or re-registered by the generated project.

The exception is maintenance of the canonical template itself. Template maintainers may change managed files when implementing a deliberate template, security, test, documentation, or lifecycle change, with the corresponding version, validation, and review updates.

## Derived repository onboarding checklist

Complete this checklist after generating a repository from the template and before running the scheduled or manual update workflow:

- [ ] Create a dedicated fine-grained personal access token or GitHub App installation token scoped to the derived repository.
- [ ] Grant Contents: Read and write, Workflows: Read and write, and Pull requests: Read and write.
- [ ] Add the token as the repository Actions secret `TEMPLATE_UPDATE_TOKEN` under `Settings -> Secrets and variables -> Actions`.
- [ ] In `Settings -> Actions -> General`, enable read and write workflow permissions.
- [ ] Enable GitHub Actions pull request creation if the organization policy exposes that option.
- [ ] Run the workflow manually and verify that it creates an update branch and pull request.

The template cannot complete updates that modify `.github/workflows` with only the built-in `GITHUB_TOKEN`. Configure `TEMPLATE_UPDATE_TOKEN` before the first update to avoid a push error related to missing `workflows` permission.

### Workflow update token details

The built-in `GITHUB_TOKEN` can create branches and pull requests, but GitHub rejects pushes that create or modify files under `.github/workflows` unless the token has the special Workflows repository permission. Each derived repository should configure an Actions repository secret named `TEMPLATE_UPDATE_TOKEN` before enabling automatic template updates.

Use a dedicated fine-grained personal access token or GitHub App installation token scoped to the derived repository with:

- Contents: Read and write
- Workflows: Read and write
- Pull requests: Read and write

The workflow falls back to `GITHUB_TOKEN` when the secret is absent, which is sufficient only for updates that do not modify workflow files. Never commit the token or place it in `.env` files.

If a generated repository already contains a manual token edit and the update reports a conflict in `.github/workflows/template-update.yml`, resolve the file by keeping the latest template detection/provenance logic and the `TEMPLATE_UPDATE_TOKEN` checkout and `GH_TOKEN` settings. Run `go test ./...`, `go vet ./...`, and `git diff --check` before committing the resolved update.

## Compatibility policy

- Patch and compatible minor template updates may be proposed automatically when the declared Go and PostgreSQL compatibility ranges remain valid.
- Changes to public routes, configuration semantics, database migrations, dependency major versions, or security behavior require an explicit compatibility note.
- Breaking changes require a new documented template version and a migration note for derived repositories.

## Post-merge repository checklist

After the backend and separate frontend PRs are reviewed and merged:

- [x] Mark `ThatSoftwareCompany/template-go-api` as a GitHub Template Repository in `Settings -> General -> Template repository`.
- [ ] Mark the frontend template repository as a GitHub Template Repository in `Settings -> General -> Template repository`.
- [ ] Confirm that template generation preserves the manifest and setup-script behavior.
