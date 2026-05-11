# Repository Guidelines

## Project Structure & Module Organization

This repository is primarily a Go 1.22 backend service under `backend/`. The service entrypoint is `backend/main.go`, with Cobra command wiring in `backend/command/`. HTTP routing and middleware live in `backend/router/`, `backend/api/`, `backend/middleware/`, and `backend/initialize/`. Shared helpers are in `backend/utils/`, `backend/common/`, `backend/constants/`, and `backend/validators/`.

Database and generated GORM code are split across `backend/dao/` and `backend/model/`; files ending in `.gen.go` should be treated as generated artifacts. RPC definitions live in `backend/proto/`. Environment-specific configs are `backend/application.dev.yml`, `backend/application.test.yml`, and `backend/application.prod.yml`. `frontend/` currently exists but contains no tracked application files.

## Build, Test, and Development Commands

Run commands from `backend/` unless noted:

```sh
make runDev
```

Starts the service with `go run main.go -envString dev`.

```sh
make runProd
```

Starts locally with production config.

```sh
go test ./...
```

Runs all Go tests. Use this before opening a pull request.

```sh
make build
```

Builds a Linux amd64 binary named `gin-admin-api`. `make buildMac` builds a Darwin amd64 binary.

```sh
protoc --go_out=. --go-grpc_out=. ./proto/account/account.proto
```

Regenerates Go protobuf outputs after editing `backend/proto/account/account.proto`.

## Coding Style & Naming Conventions

Use standard Go formatting: run `gofmt` or `go fmt ./...` before committing. Keep package names short and lowercase. Exported identifiers use PascalCase; unexported identifiers use camelCase. Follow the existing directory conventions for DTOs, VOs, DAOs, middleware, and generated models rather than introducing new layering. Do not hand-edit `.gen.go` or `.pb.go` files; update the generator input instead.

## Testing Guidelines

The project uses Go’s built-in testing framework. Add tests as `*_test.go` files next to the package being tested. Current tests are sparse, so prioritize coverage for validators, utility functions, middleware behavior, and any changed business logic. Prefer table-driven tests for validators and helpers. Run `go test ./...` from `backend/` before submitting changes.

## Commit & Pull Request Guidelines

Git history is not available in this checkout, so use concise imperative commit subjects, for example `fix token validation` or `add message retry tests`. Keep generated-code commits separate from handwritten logic when possible. Pull requests should include a short summary, test results, configuration or migration notes, and linked issue IDs when available. Include screenshots only for future frontend changes.

## Security & Configuration Tips

Do not commit secrets, private credentials, or environment-specific production values. Review changes to `application.*.yml`, Kafka certificates, Redis/MySQL settings, and auth middleware carefully. Keep local-only overrides outside the repository or document them in the PR without exposing sensitive values.
