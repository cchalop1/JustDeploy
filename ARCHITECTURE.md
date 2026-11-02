# JustDeploy – Project Architecture

**JustDeploy** is a self-hosted PaaS (Platform-as-a-Service) that helps developers deploy their applications more easily on their own VPS or server. It provides a clean and modern developer experience, without vendor lock-in.

---

## 🧱 Architecture Overview

The project is split into two main parts:

1. **Frontend (Web)**
2. **Backend (Go CLI/API)**

Both are tightly integrated, with a unified build process using a `Makefile`.

---

## 1. Web (Frontend)

### 📦 Tech Stack

- **React** (UI)
- **Vite** (build tool)
- **Bun** (JS runtime for install/build)
- **TailwindCSS** (for styling)
- **HTTP**: REST API communication with the backend

### 📁 Location

- Source code: `web/`
- Final build output: `internal/web/dist/` (copied here by the Makefile)

### 🔧 Responsibilities

- Provide a user-friendly interface to:
  - Create/manage apps
  - Configure services (DB, env, etc.)
  - Trigger deployments
  - View logs and environments
- Communicates with backend over HTTP (REST)

---

## 2. Backend (Go)

### 🛠️ Tech Stack

- **Language**: Go
- **Web Framework**: [Echo](https://echo.labstack.com/)
- **Binary**: Compiled CLI named `justdeploy`

### 📁 Structure

- Main entry point: `cmd/just-deploy/main.go`
- Web assets: served from `internal/web/dist/`
- Versioning: injected via `ldflags` from Git tags

### 🔧 Responsibilities

- Serve the frontend (as static assets)
- Expose a REST API for the frontend
- **Communicate with the Docker socket** to:
  - **Build** containers
  - **Run** containers
  - **Stop** and **delete** containers
- Manage app deployments, configs, and logs
- Handle environment variables and CLI integration

---

## 3. Build System

The project uses a `Makefile` to streamline development and CI builds.

## Security and persistence

State is stored in SQLite (`JUSTDEPLOY_DATA_DIR/database.db`, defaulting to the legacy HOME config directory outside Docker). Service runtime metadata stores the connection group, branch, Dockerfile, data volumes and settings revision. Additive schema migration preserves service rows.

All mutating deployments share an operation mutex. Global deployment preserves unchanged running services; targeted redeployment fetches only the selected repository. Image construction happens before replacement, data mounts are retained, and the previous container is kept until startup checks pass. The previous container is restarted on failure.

Administrative remote access requires TLS (direct or through explicitly trusted reverse proxies). First setup requires a one-time server token. GitHub push deliveries are authenticated over the raw payload and deduplicated persistently. See [SECURITY-UPGRADE.md](SECURITY-UPGRADE.md) for migration, operational limits and verification.
