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

### Workflow update token

The built-in `GITHUB_TOKEN` can create branches and pull requests, but GitHub rejects pushes that create or modify files under `.github/workflows` unless the token has the special Workflows repository permission. Each derived repository should configure an Actions repository secret named `TEMPLATE_UPDATE_TOKEN` before enabling automatic template updates.

Use a dedicated fine-grained personal access token or GitHub App installation token scoped to the derived repository with:

- Contents: Read and write
- Workflows: Read and write
- Pull requests: Read and write

The workflow falls back to `GITHUB_TOKEN` when the secret is absent, which is sufficient only for updates that do not modify workflow files. Never commit the token or place it in `.env` files.

## Compatibility policy

- Patch and compatible minor template updates may be proposed automatically when the declared Go and PostgreSQL compatibility ranges remain valid.
- Changes to public routes, configuration semantics, database migrations, dependency major versions, or security behavior require an explicit compatibility note.
- Breaking changes require a new documented template version and a migration note for derived repositories.

## Post-merge repository checklist

After the backend and separate frontend PRs are reviewed and merged:

- [x] Mark `ThatSoftwareCompany/template-go-api` as a GitHub Template Repository in `Settings -> General -> Template repository`.
- [ ] Mark the frontend template repository as a GitHub Template Repository in `Settings -> General -> Template repository`.
- [ ] Confirm that template generation preserves the manifest and setup-script behavior.
