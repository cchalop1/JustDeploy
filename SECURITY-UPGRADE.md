# Security and reliability upgrade

This branch changes administrative access and deployment semantics. Back up before installing it on an existing server. No production secrets have been rotated and no production data has been migrated by the code-editing session.

## Before upgrading

1. Stop JustDeploy before copying its old state, or use SQLite's online backup mechanism. Copy the whole previous data directory, including any SQLite WAL/SHM files, to a protected backup location. Test opening a backup on an isolated instance.
2. Export logical backups of the application databases using each engine's backup tool. A copy of JustDeploy's SQLite database does not back up Postgres, Redis, Mongo or object storage.
3. Record container mount names and destinations using `docker inspect`. Existing anonymous volumes and bind mounts are adopted on replacement. Volumes are never deleted by service deletion. A legacy database without a reusable mount is deliberately refused for replacement: migrate its data into a named volume first. Container rollback cannot undo database schema changes.
4. Revoke GitHub tokens and rotate database/application credentials that were printed by older versions. Review journal/container logs and copies of Git workspaces for exposed credentials. The branch prevents new leaks; it cannot revoke previously copied secrets.

## Data directory

`JUSTDEPLOY_DATA_DIR` overrides the historical `$HOME/.config/justdeploy` location. The runtime image sets it to `/app/data`, matching the documented mount. Before replacing an old Docker container, copy `/root/.config/justdeploy` out of that container into the persistent host directory. Do not delete the old container until recovery has been checked.

Native installations use `/var/lib/justdeploy`; the installer copies a legacy root configuration there if the new directory does not exist. Other HOME directories require explicit migration. Data directories use 0700 and SQLite uses 0600. Backups contain secrets and need the same access controls.

SQLite receives an additive `runtime_settings` column and a webhook delivery table; existing rows remain. Malformed historical GitHub JSON tags are accepted on read. Take a backup before migrating: old binaries do not understand the new runtime settings.

## Administration and first setup

The native API binds to `127.0.0.1` by default. Set `JUSTDEPLOY_LISTEN_ADDR` to change this. Remote cleartext HTTP is rejected. Use either:

- Direct TLS with `JUSTDEPLOY_TLS_CERT` and `JUSTDEPLOY_TLS_KEY`.
- An HTTPS reverse proxy. Set `JUSTDEPLOY_PUBLIC_URL=https://deploy.example.com` and `JUSTDEPLOY_TRUSTED_PROXIES` to the narrow CIDR of the proxy peer (comma separated). Forward `X-Forwarded-Proto: https`. Forwarded headers from other peers are ignored.

For Docker, publish the administration port on the host loopback only (`-p 127.0.0.1:5915:5915`), and configure the actual proxy peer CIDR as seen by the container. Do not use a trust-all CIDR. For local setup, an SSH tunnel to the loopback API can be used.

The first administrator account is created from the web interface using a valid email and a password of 12–72 bytes. Anyone who can reach the unconfigured instance can create this account, so complete initial setup before making the administration interface publicly accessible. Creation is atomic; a second request cannot overwrite the account. Existing administrator accounts are preserved.

JWTs expire after eight hours; browser storage is per-tab. `POST /api/logout`, authenticated with the current bearer token, revokes all sessions. Upgrading signs out browsers that used the old persistent storage.

## Applications, storage and networking

Each service gets a private connection group. Services imported from one supported Compose file share a group. Use the settings form to explicitly put an application and its database in the same group, then deploy those changed services. Configure connection environment variables explicitly; database passwords are no longer copied into every application.

Existing services stay on their old Docker network until redeployed. Plan the network transition for related services together; check connectivity in an isolated copy first. Public applications also join the proxy network and can reach other public applications; connection groups protect private dependencies, not all network traffic between exposed applications.

New known database/model images receive stable named volumes; existing mounts are reused. Every container has CPU, memory and PID limits. Image builds happen before stopping an existing application. During replacement, the previous container is retained and restarted if startup/health checks fail. This is rollback with brief downtime, not zero-downtime deployment. Without an image HEALTHCHECK, readiness checks only that the process stays alive for two seconds.

Global deploy starts new/changed services, with dependencies before repository applications. Already-running unchanged services are retained. `POST /api/deploy/redeploy/:id` refreshes and deploys only that service. Authenticated push webhooks match the full repository name and configured branch. Duplicate deliveries are retained for seven days; failed jobs release their delivery ID. Jobs are serialized in-process and have bounded build time; they are not a durable job queue. GitHub failed deliveries must be retried/redelivered explicitly.

Traefik's insecure API is disabled. `JUSTDEPLOY_DOCKER_PROXY` can point Traefik at a restricted socket proxy; otherwise its provider requires a raw Docker socket. JustDeploy itself also requires privileged Docker access. Treat both components as infrastructure administrators; only trusted code/images should be built on this host.

## Supported Compose import

Recognized filenames: `compose.yaml`, `compose.yml`, `docker-compose.yaml`, `docker-compose.yml`.

The importer supports one application build at the repository root, image dependencies, explicit environment values, list-form commands, numeric short-form ports, short-list dependencies, and named-volume mounts (`name:/absolute/path`). Image services without ports are valid. Host ports on dependencies are not published. Dependency DNS uses the Compose service name.

Advanced build contexts, multiple builds, variable interpolation, bind mounts, external volumes, custom networks, secrets, privilege settings and other unsupported options fail explicitly. This is a restricted importer, not a replacement for the full Docker Compose specification. Build-required values must be supplied through supported build mechanisms; `.env`, private keys and `.git` are excluded from contexts.

Repositories use their actual default branch on import. Dockerfile variants are detected; multiple variants require an explicit selection. Nixpacks is optional and must be installed from a verified versioned release for repositories with no Dockerfile. The runtime no longer executes an unversioned installer script.

## Verification

Run `go test -race ./...`, `go vet ./...`, and in `web`, `bun install --frozen-lockfile`, `bun run lint`, `bun run build`, `bun audit`. The CI also scans Go dependencies and builds/scans the runtime image. The default Go Docker tests use an HTTP simulator. The opt-in `TestDockerIntegrationPersistenceAndRollback` creates and removes only its own randomly named containers, volume and network when `JUSTDEPLOY_DOCKER_TESTS=1` is set.

Before a production rollout, verify with an isolated Docker daemon: create data, redeploy, confirm that data remains, fail a new application version and observe rollback, check public/private visibility, and restore a database backup. Rebuilding and scanning application images and renewing external credentials remain deployment operations.
