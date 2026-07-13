# Repository Guidelines

## Project Structure & Module Organization

`cmd/server/main.go` wires configuration, SQLite storage, services, and routes. Under `internal/`, `api/public/upload` and `api/public/download` own their handlers; `api/http_utils` contains shared HTTP helpers; `service/uploader` holds business logic; and `storage/sqlite` persists files. Domain types are in `internal/entity/upload` and configuration is in `internal/config`.

Keep package scope narrow: each package should have one responsibility. Put dependency interfaces only in the consumer's `contract.go`; generated GoMock files live in the adjacent `mocks/` directory. API documentation is `schema.yaml`, while the visual design reference is under `design/`.

## Build, Test, and Development Commands

Run `make help` to see all commands:

```sh
make env     # create .env from .env-example if needed
make run     # start the server
make api     # run curl-based checks against the running server
```

Use `make stop` or `make restart` to manage the configured port. Before submitting, run `make fmt`, `make test`, `make lint`, and `make build`. Run `make mocks` after changing a contract; do not edit generated files.

## Coding Style & Naming Conventions

Use idiomatic Go and tabs; `make fmt` runs `gofmt`, and `make lint` runs `go vet`. Keep names conventional and clear. Injected dependencies use no more than two words and end in `-er`, for example `FileGetter`, `FileSaver`, or `FileUploader`; use the same name for the interface, struct field, and related code.

Add short English comments above reusable functions, in long functions (more than 30 lines), and around non-obvious domain logic. Explain intent, such as security limits or compatibility behavior; mark temporary decisions with `TODO`.

## Testing Guidelines

Write table-driven Go tests in `*_test.go` beside the tested package. Use `go.uber.org/mock/gomock` and only generated mocks from `mocks/`; do not add manual stubs. Cover success, validation, and dependency-error paths. Update `api.sh` when a public endpoint behavior changes, then run `make api` with the server running.

## Commit & Pull Request Guidelines

Make focused commits using Conventional Commit-style prefixes: `feat(api):`, `fix(api):`, `build:`, `docs:`, or `refactor:`. Do not mix generated files, database files, or unrelated formatting with a feature. Pull requests should summarize behavior, list verification commands, reference an issue when available, and update `schema.yaml` when the public API changes.

## Configuration & Data

Keep secrets out of Git. Required local settings are `SERVER_HOST`, `SERVER_PORT`, and `DATASOURCE_PATH` in `.env`. Treat `sqlite/main.db` as local runtime data unless a change to it is explicitly required.
