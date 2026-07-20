# AGENTS.md

This document captures what an agent needs to know to work effectively in this Go codebase.

## Essential Commands

Use **mage** (not make) for all development tasks:

| Command | Description |
|---------|-------------|
| `mage` or `mage run` / `mage r` | Build dev binary + run |
| `mage build` | Build dev binary to `tmp/ape` |
| `mage build_prod` | Build static prod binary to `deployment/ape` |
| `mage test` / `mage t` | Run all tests via `go tool gotest` |
| `mage audit` | `mod tidy/verify`, `fmt`, `vet`, `staticcheck` |
| `mage ci` | `test` + `audit` (matches CI pipeline) |
| `mage tools` | Install required tools (migrate, gotest, staticcheck, bobgen, mage) |
| `mage encrypt` / `mage decrypt` | Encrypt/decrypt secrets via fish script |
| `mage enc_sops` | SOPS encryption for Kubernetes GitOps |
| `mage build_kube` | Generate k8s manifests from `cmd/kubernetes` |
| `mage flux_reconcile` | Trigger Flux source reconciliation |
| `mage version` / `mage v` | Print version |

CI (`.gitea/workflows/ci-pipeline.yaml`) runs `go build ./cmd/web`, then `mage test`, then `mage audit`.

## Project Structure & Architecture

**Vertical slice packages** (domain-oriented, not layered):

- `pkg/day/` — Day aggregate, month queries, auto-habit creation
- `pkg/habit/` — Empty placeholder only
- `pkg/event/` — Bus interface + in-memory impl (stubbed usage only)

**Layered packages**:

- `cmd/web/main.go` — Entry point, config, migrations, session setup, background sync ticker
- `pkg/api/` — HTTP handlers, middleware, route registration on dove
- `pkg/application/` — Use cases / orchestration (App struct)
- `pkg/persistence/` — Bob-backed repository implementations
- `pkg/core/` — Domain models, value objects (Date), repository interfaces, errors
- `pkg/dove/` — Minimal custom web framework (router + context + renderer)
- `pkg/oak/` — Structured logger wrapper over slog with layer stacking
- `pkg/ui/` — Gomponents views + UI State
- `pkg/integrations/` — External API clients (fitbit, toggl, hevy, google)
- `pkg/validator/` — Minimal validation
- `pkg/config/` — Constants (APP_NAME="ape", paths, env keys)
- `pkg/deployment/` — Secret handling helpers

**Data flow**:
1. HTTP → `dove` routes → `api` handlers
2. Handlers → `application.App` methods
3. App → `persistence.Database` (repos) or `day.App`
4. Repos → bob generated models + raw queries
5. Core types returned up the stack

**Middleware order** (see `api/routes.go`):
panic → dev (no-cache) → loadCookie → authFromSession → setUIState → requireAuth

## Key Patterns & Gotchas

### Error Handling (inverted logic)
- `core.IfErrNNotFound(err)` returns `true` when err is **NOT** the not-found error. Used as: `if core.IfErrNNotFound(err) { return nil, err }` then proceed to create.
- `catchDBErr(op, err)` in persistence: `sql.ErrNoRows` → `core.ErrRepositoryNotFound`; other errors are logged and wrapped.
- Never compare errors directly; use `errors.Is(err, core.ErrRepositoryNotFound)` or the helper.

### Date Value Type
- `core.Date` wraps `time.Time` normalized to UTC midnight.
- Always construct via `core.NewDate(t)` or `core.NewDateToday()`.
- Use `.Time()` to get `time.Time`; `.String()` produces ISO8601.
- `RangeMonth()` is exclusive of the target date (includes it at end).

### Habit Auto-Creation & Upsert Rules
- `day.App.GetOrCreate` creates 4 default habits (wake_up, fitness, deep_work, food) with `HabitStateNoInfo`.
- In `persistence/habits.go:55-66`: **non-automated habits are never overwritten by automated upserts**. Check `Automated` flag and `State != no_info`.

### Deadlines
- `DaysLeft` computed on every `Find` using `core.DaysLeft(today, end)`.
- `application.DeadlineList` mutates: negative days → recurring bumps +1 year or deletes.
- Results are sorted by `DaysLeft` after mutation.

### Database
- SQLite only. DSN modified at runtime: `?_busy_timeout=5000&_journal_mode=WAL`.
- Migrations: `database/migrations/sqlite/` (golang-migrate).
- Bob codegen: `bobgen.yaml` → `database/gen/models/`. Run bobgen manually when schema changes.
- Schema migrations table and `sessions` are excluded from bob.

### Testing
- In-memory SQLite + **real migrations** per test (duplicated `NewTestDB` helpers).
- Migration path is relative: `file://../../database/migrations/sqlite` — run tests from package dir or it fails.
- Use `pkg/test` helpers: `test.True`, `test.False`, `test.NilErr` (no testify).
- Tests use `t.Parallel()`.
- `application_test.go` has `NewTestApplication(t)` helper; `persistence/auth_test.go` has its own.

### Logging
- Never use `log` or `slog` directly. Use `oak`.
- `oak.SetDefault` in main; per-request via `oak.FromContext(ctx)` or `oak.FromContextWithLayer`.
- Layers stack with dots: `logger.Layer("app").Layer("sync")` → `layer=app.sync`.
- `Oak` is mutable (Layer/PopLayer/ClearLayers); clone via context for isolation.

### Assets & Caching
- Dev: read from `./assets/` on disk, `Cache-Control: no-cache, must-revalidate`.
- Prod: embedded via `assets/efs.go`, `Cache-Control: public, max-age=86400, immutable`.
- Strong ETag (sha256 content hash) + 304 support in both modes.

### Configuration
- `config.json` (gitignored). Must exist; read unconditionally.
- `APE_ENVIRONMENT` env var required (dev/prod/test). Parsed via `application.ParseEnvironment`.
- Admin/user/guest credentials from config file, set at startup via `app.SetPassword`.
- Version injected at build: `-X main.version=...`. Binary supports `ape version` CLI.

### Background Sync
- `main.go:175-180`: goroutine ticker every 6 hours calls `app.DaySync(date)`.
- `DaySync` fans out sleep/fitness/deepWork sync concurrently via channels (errors logged, not returned).

### Event Bus
- `pkg/event` defines `Bus` interface with `Publish`/`PublishAsync`/`Subscribe`.
- Currently only a stub subscription in `main.go` for "day-created". Not wired to real flows.

### Integrations
- Fitbit/Google: oauth2 with refresh; tokens in `oauth_tokens` table.
- Toggl/Hevy: static token/key.
- Google is commented out in `integrations/sync.go:36`.
- `application.IntegrationGet` returns status + info; may return auth URL for oauth flows.

### UI Framework
- `maragu.dev/gomponents` for HTML-as-Go.
- `gomponents-datastar` for client-side reactivity (see `ds.On("click", ...)`).
- State passed via context: `ui.StateWithContext` / `StateFromContext`.
- All pages require auth except `/login`.

## Non-Obvious Details

- `pkg/habit/habit.go` exists but `App` is empty — habit logic lives in `day` + `core` + `persistence/habits.go`.
- `persistence.Database` embeds both interface types (`core.DayRepository`) and concrete types (`Users`, `Auths`).
- `day.App` has its own `NewDays`/`NewHabits` — slight duplication with `persistence.New`.
- Delete on deadlines uses wrong table name in query: `models.DeleteWhere.Days.ID` (copy-paste bug, but exists).
- `runDataMigration` in main is entirely commented out.
- No stdlib `http.ServeMux`; all routing through dove's `Endpoint`/`Prefix`/`GET`/`POST` etc.
- Session cookie name hardcoded: `"curious-ape-session"`. 7-day lifetime.
- Prod sessions set `HttpOnly` + `Secure`; dev does not.

## When Changing Things

- Schema change → update migration + `bobgen.yaml` + regenerate models.
- New integration → add to `core.Integration` consts, `integrations/sync.go`, `application/integrations.go`.
- New route → add in `api/routes.go`, implement handler, add to UI nav in `ui/layout.go` if needed.
- New background job → consider wiring to event bus instead of ad-hoc goroutines (future intent).
- Tests touching DB → ensure migration path resolves (run from package or adjust relative path).
