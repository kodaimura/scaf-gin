# Contributing

Keep changes focused, reviewable, and covered by tests. Never commit credentials,
tokens, personal data, or a production `.env` file.

## Before opening a pull request

Run the same checks used by CI:

```sh
make check
make test_e2e
make build_prod
```

When behavior changes, add or update tests at the same level as the change.
Update `.env.example` and documentation when configuration changes.
Include a numbered SQL migration for database schema changes, and verify both the migration and application behavior.

Dependency audit findings should be reviewed according to exploitability,
available fixes, and upgrade cost. They do not automatically block unrelated
pull requests.

## Pull requests

- Explain the reason for the change, not only the implementation.
- Keep unrelated changes in separate pull requests.
- Describe contract, data, configuration, and deployment impact.
- Resolve review comments and make sure CI passes before merging.
- Complete the relevant items in the pull request template.

Report vulnerabilities privately by following [SECURITY.md](SECURITY.md), not
through a public issue or pull request.

