# UnlearnNLearn — Project Architecture (2026)

## Overview

This repository contains a small Go backend and an Angular frontend that together form the `UnlearnNLearn` example project. The codebase demonstrates several concurrency patterns, worker-pool implementations, and a simple web frontend to consume backend APIs.

High-level responsibilities:
- Backend: implemented in Go — provides services, data access, handlers, and example pipelines/workers.
- Frontend: implemented with Angular under `my-website` — provides UI, components, and API integration.

## Architecture

- Client (Angular): single-page application located at `my-website/frontend/src`.
- API Server (Go): entrypoints in `cmd` and `my-website/cmd` provide HTTP endpoints that are handled by the `handler` package and implemented using services in the `service` package.
- Data access: `dao` contains data access logic and models live in `models` (top-level and `dao/models.go`).
- Concurrency examples: multiple folders demonstrate worker-pools and producer-consumer patterns (`worker-pools`, `worker-pool-hardening`, `refactoring`, `producer-consumer`, `pipelines`).

### Key backend folders

- `cmd/` — backend application entrypoint(s). Look at `cmd/main.go` for server startup.
- `dao/` — data access objects and DB-related helpers.
- `handler/` — HTTP handlers, middleware, and response helpers.
- `models/` — domain models used across the backend.
- `service/` — business logic, authentication (`auth.go`), and worker/service implementations.
- `pipelines/`, `producer-consumer/`, `worker-pools/`, `worker-pool-hardening/`, `refactoring/` — example programs showing concurrency patterns and designs used for learning and experimentation.

### Frontend layout

- `my-website/` — Angular project root.
	- `my-website/frontend/src/app/components` — feature components (`home`, `blogs`, `blog-detail`, `profile`).
	- `my-website/frontend/src/app/services` — API service integration in `api.service.ts`.
	- `my-website/frontend/src/app/models` — shared TypeScript models.

## How to run (Local, Windows)

Backend (Go):

1. From repository root, build and run the backend server:

```bash
cd cmd
go run .
```

or build a binary:

```bash
go build -o unlearn-server ./cmd
.
\unlearn-server.exe
```

Frontend (Angular):

1. Install dependencies and run dev server:

```bash
cd my-website
cd frontend
npm install
npm start
# or: ng serve --open
```

2. The frontend expects the backend API to be running; configure the API base URL in `my-website/frontend/src/app/services/api.service.ts` if needed.

## Development notes

- The repo contains multiple example programs; `cmd/main.go` is the main API server used by the frontend. Other example folders are educational and can be run individually (they often contain their own `main.go`).
- Use the `handler` and `service` layers to add new API endpoints. Keep business logic in `service` and HTTP concerns in `handler`.
- Concurrency examples are intentionally separate to make them easy to run and test independently.

## Contributing

- Open an issue describing the change you want to make.
- Send a PR with a clear description and small, focused commits.

## Contact

For questions, check package comments or open an issue in this repo.
