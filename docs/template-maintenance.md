# Template maintenance

The manifest is the source of truth for the template identity, version, source repository, compatibility, and generated origin.

## Implemented lifecycle safeguards

- `scripts/setup.sh` records `template_commit` from the source Git commit when available.
- `scripts/setup.sh` records `generated_from` from the derived repository identifier or module path.
- `cmd/template -command validate` validates the manifest contract.
- `scripts/validate-template.sh` verifies required template files and blocks a real `.env` file.
- CI runs the lifecycle validation and writes build outputs outside the repository root.
- `scripts/template-update.sh` applies a three-way patch between recorded and target template commits.
- `.github/workflows/template-update.yml` detects version tags and opens derived-repository PRs with least-privilege write permissions.
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

The workflow requires GitHub Actions to be allowed to create pull requests in the derived repository. It is skipped when running in the canonical template repository itself.

## Compatibility policy

- Patch and compatible minor template updates may be proposed automatically when the declared Go and PostgreSQL compatibility ranges remain valid.
- Changes to public routes, configuration semantics, database migrations, dependency major versions, or security behavior require an explicit compatibility note.
- Breaking changes require a new documented template version and a migration note for derived repositories.

## Post-merge repository checklist

After the backend and separate frontend PRs are reviewed and merged:

- [x] Mark `ThatSoftwareCompany/template-go-api` as a GitHub Template Repository in `Settings -> General -> Template repository`.
- [ ] Mark the frontend template repository as a GitHub Template Repository in `Settings -> General -> Template repository`.
- [ ] Confirm that template generation preserves the manifest and setup-script behavior.
