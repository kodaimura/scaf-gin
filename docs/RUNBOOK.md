# Deployment Runbook

This runbook defines checks, execution boundaries, verification, and recovery
steps shared by every deployment environment. Complete the project-specific
values before the first deployment and keep them current.

## Project-specific values

Document references and commands, but never credentials or secret values.

| Item | Value |
| --- | --- |
| Environments | `<development, staging, production, ...>` |
| Deployment command | `<command or script>` |
| Authentication command | `<command or documented login procedure>` |
| Health check | `<URL or command>` |
| Critical flows | `<smoke tests or manual checks>` |
| Logs and metrics | `<location or command>` |
| Data backup and recovery | `<document or procedure>` |
| Previous release lookup | `<artifact, image, or revision lookup>` |
| Incident contact | `<team or contact method>` |

## Authorization and scope

- Identify the target environment and exact revision before running anything.
- An implementation, review, or commit request does not authorize deployment.
- Deploy only when the environment and revision are unambiguous.
- Use approved credentials without printing or storing them in this repository.
- Record the operator, environment, revision, and start time.

## Before deployment

1. Confirm the working tree and exact revision:

   ```sh
   git status --short
   git rev-parse HEAD
   ```

2. Confirm the revision is available from the expected remote and CI passed.
3. Run the scaffold checks:

   ```sh
   make check
   make test_e2e
   make build_prod
   ```

4. Review dependency audit results when dependencies changed or the release
   requires a security review.
5. Confirm required environment variables and secrets exist without displaying
   their values.
6. Review pending database migrations for locking, deletion, data conversion, and backward compatibility.
7. Confirm backup and recovery before an irreversible migration.
8. Identify the previous known-good artifact or revision.

## Deployment

1. Authenticate with the documented project-specific procedure.
2. Run the deployment command with an explicit environment and revision.
3. Stop on a non-zero exit status; preserve sanitized output.
4. Record the deployed artifact or revision and completion time.

Do not invent a replacement command when the documented command is missing or
fails.

## Verification

1. Confirm the running revision when the platform exposes it.
2. Confirm the health check succeeds.
3. Verify the documented critical flows.
4. Review startup logs, application errors, and available service metrics.
5. Report the environment, revision, checks, and result.

## Stop conditions

Stop and request a decision when the environment, revision, credentials, or
target is ambiguous; required checks failed; deployment failed; health or
critical flows fail; or continuing requires a destructive undocumented action.

## Rollback and incident record

Prefer the previous known-good artifact. Do not automatically reverse database
changes or perform destructive recovery. After rollback, repeat verification.

Record the environment, requested and actual revision, timing, failed checks,
sanitized logs, user and data impact, mitigation, and follow-up work.
