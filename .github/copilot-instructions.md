# Project Guidelines

## Architecture

This project uses a strict 4-layer architecture. Never bypass layers.

```
Handler → Service → Repository Interface → Repository Implementation
```

- **Handlers** (`internal/handler/`): HTTP concerns only — decode, validate, delegate to service, encode response. No business logic.
- **Services** (`internal/service/`): Business logic. Depend on repository **interfaces**, never implementations.
- **Repository Interfaces** (`internal/repository/`): Contracts only. One interface per entity.
- **Repository Implementations** (`internal/repository_impl/`): PostgreSQL via `pgxpool`. Always use parameterized queries (`$1, $2`).

Dependency wiring lives in `internal/server/app.go`.

## Code Style

- Use Go standard library where possible — no unnecessary third-party packages.
- Use `log/slog` for all logging (structured, key-value pairs).
- Use `go-chi/chi/v5` for routing. Each resource gets its own router file in `internal/router/`.
- Use `go-playground/validator/v10` for request validation via struct tags.
- Pointer fields (`*string`, `*float64`) in update DTOs for partial updates.
- JSON tags on all exported struct fields. Use `json:"-"` to hide sensitive fields (e.g., `PasswordHash`).

## Logging

Use `log/slog` with structured key-value fields. Follow these severity rules:

| Level        | Use When                                                                        |
| ------------ | ------------------------------------------------------------------------------- |
| `slog.Info`  | Request entry, successful completions                                           |
| `slog.Warn`  | Client errors (400s): bad input, validation failures, parse errors              |
| `slog.Error` | Server errors (500s): unexpected failures, DB errors, token generation failures |

Always include these structured fields:

- `"handler"` — fully qualified handler name (e.g., `"AssetHandler.Create"`)
- `"method"`, `"path"` — on request entry
- `"error"` — on any failure
- Entity IDs (`"userID"`, `"assetID"`, etc.) — when available

## Error Handling

- Use `util.AppError` for all domain errors. Create via `NewNotFoundError()`, `NewBadRequestError()`, `NewInternalError()`.
- Services return `*AppError`; handlers call `handleError(w, err, "HandlerName.Method")` which logs and responds.
- Never expose raw internal errors to clients — wrap them in `AppError`.

## Request/Response

- Decode request bodies with `json.NewDecoder(r.Body).Decode()`.
- Validate with `util.Validate.Struct()` immediately after decode.
- Respond via `util.SendResponse(w, statusCode, data)` and `util.SendErrorResponse(w, statusCode, message)`.
- Success responses: `map[string]any{"message": "...", "entity": entity}`.
- Parse path params with `r.PathValue()` and `parseIntegerID()` helper.
- Pagination via query params `limit` (default 50, max 200) and `offset` (default 0).

## Authentication

- JWT auth via middleware (`internal/middleware/jwt_middleware.go`), applied globally.
- Public routes: `/swagger/*`, `/api/users/*`.
- Extract user ID in handlers via `util.GetUserIDFromContext(r.Context())`.
- Tokens generated with `util.GenerateToken(userID, email)` — 24-hour expiry.

## Database

- PostgreSQL via `jackc/pgx/v5` with connection pooling (max 20, min 5).
- Connection string from `DATABASE_URL` env var.
- All queries use parameterized placeholders (`$1, $2`), never string interpolation.
- Migrations in `db-migration/` with timestamp-prefixed filenames.

## Swagger

- Annotations on every handler method using `swaggo/swag` comment format.
- Generate with: `go run github.com/swaggo/swag/cmd/swag@v1.16.6 init -g cmd/mf_stock_tracker/main.go -o docs`
- Always include `@Security BearerAuth` on protected endpoints.

## Conventions

- One file per entity per layer (e.g., `asset_handler.go`, `asset_service.go`, `asset_repository.go`).
- Constructor functions follow `NewXxx(deps) *Xxx` pattern.
- Custom validators registered in `internal/util/validator.go` (e.g., `instrument_type`, `txn_type`).
- Soft delete for users (`is_active` flag), hard delete for other entities.
