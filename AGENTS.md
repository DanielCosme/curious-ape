# AGENTS.md

This document captures what an agent needs to know to work effectively in the curious-ape repository.

## Project Overview

Curious Ape is a personal habit/life tracking web application. It tracks:
- Daily habits (wake up, fitness, deep work, eat healthy)
- Sleep logs (Fitbit integration)
- Fitness logs (Hevy integration)
- Deep work time (Toggl integration)
- Deadlines

**Tech stack**: Go 1.26, SQLite (modernc.org/sqlite), bob ORM codegen, custom minimal web framework (dove), gomponents for HTML, alexedwards/scs for sessions.

**Module path**: `danicos.dev/daniel/curious-ape`

## Essential Commands

Use `mage` (task runner) for nearly everything. Aliases: `v` (version), `r` (run), `t` (test).

| Command | Description |
|---------|-------------|
| `mage` or `mage run` / `mage r` | Build + run dev server |
| `mage build` | Build dev binary to `./tmp/ape` |
| `mage build_prod` | Build production static binary |
| `mage build_kube` | Generate Kubernetes manifests via `cmd/kubernetes` |
| `mage test` / `mage t` | Run all tests via `go tool gotest ./...` |
| `mage audit` | `go mod tidy/verify`, `fmt`, `vet`, `staticcheck` |
| `mage ci` | Runs `Test` + `Audit` (matches CI) |
| `mage tools` | Install dev tools (migrate, gotest, staticcheck, bobgen, mage) |
| `mage db:open` | Open dev SQLite DB in sqlite3 REPL |
| `mage db:gen` | Regenerate bob models from `bobgen.yaml` |
| `mage migrate:new <name>` | Create new timestamped migration |
| `mage migrate:up` | Run migrations up |
| `mage migrate:down` | Run migrations down |
| `mage migrate:force <version>` | Force schema_migrations to specific version |
| `mage encrypt` / `mage decrypt` | Encrypt/decrypt secrets via `scripts/enc_dec.fish` |
| `mage enc_sops` | Encrypt SOPS secrets for Flux GitOps |

**Direct Go commands** (when needed):
- `go build -o ./tmp/web ./cmd/web` - verify server compiles (used in CI)
- `go tool gotest ./...` - tests with color output
- `go tool staticcheck -checks='inherit,-ST1001' ./cmd... ./pkg...` - lint (note: `-ST1001` excluded)

**Agent rule for builds**: Whenever you run `go build` (or any direct compile) *to test for compilation errors*, you **must** always place the binary in `./tmp/` using `-o ./tmp/<name>`. Never emit a binary to the project root (bare `go build ./cmd/web` produces `./web`, which is forbidden). The `tmp/` directory exists for dev artifacts and is the canonical location (see also `mage build`, which already targets `./tmp/ape`). This keeps the repository root clean.

**Database connection string for tools**: `sqlite3://./tmp/ape.db`

## Architecture and Code Organization

```
cmd/
  web/main.go           # Entry point: config, migrations, sessions, HTTP server, background sync
  kubernetes/main.go    # Generates k8s manifests

pkg/
  application/          # Business logic layer (App struct + methods)
    application.go      # App, Config, Environment, New()
    days.go, sync.go, habits.go, etc.
  api/                  # HTTP handlers + routing + middleware
    api.go              # API struct, auth helpers
    routes.go           # All route registration via dove
    middleware.go
  ui/                   # Server-side HTML rendering (gomponents)
    layout.go, habits.go, sleep.go, etc.
  core/                 # Domain models + repository interfaces (pure, no deps on persistence)
    day.go, habit.go, date.go, repository.go, error.go, constants.go
  persistence/          # Repository implementations using bob-generated models
    persistence.go      # Database struct wiring all repos
    days.go, habits.go, etc.
  dove/                 # Custom HTTP framework/router (see below)
  oak/                  # Custom slog wrapper with layers (see below)
  integrations/         # External API clients + sync orchestration
    sync.go             # Integrations struct, New()
    hevy/, toggl/, fitbit/, google/
  config/               # Constants only (APP_NAME, paths, env var names)
  validator/
  deployment/           # K8s secret helpers
  test/                 # Test helpers (test.go)

database/
  migrations/sqlite/    # golang-migrate .sql files (numbered)
  gen/models/           # Bob-generated (DO NOT EDIT)
  gen/dberrors/

assets/                 # Static files (css, fonts); efs.go for prod embedding
deployment/             # k8s overlays (flux, k3s), systemd unit, litestream config
magefiles/              # Mage tasks (magefile.go, database.go)
scripts/                # build.fish, enc_dec.fish, encrypt_sops.sh
```

**Control flow**:
1. `main.go` loads `config.json`, opens SQLite (with WAL + busy_timeout), runs migrations, sets up scs sessions, creates `application.App`, seeds admin user, starts 6h background `DaySync` goroutine, then serves `api.Routes(t)`.

2. Request: `dove.ServeHTTP` → prefix/exact route match → middleware chain → handler → `c.RenderOK(ui.XXX(state))` or redirect/error.

3. Handlers call `a.App.*` (application layer) → `a.db.*` (persistence) → bob queries → models.

4. Background sync (`DaySync`): spawns 3 goroutines for sleep/fitness/deepWork sync from integrations, then returns fresh Day.

**Data model notes**:
- `core.Date` is always normalized to UTC midnight (see `core/date.go`).
- `Day` aggregates Habits + SleepLogs + FitnessLogs + DeepWorkLogs via bob relations.
- Habits are auto-created (4 per day) when a Day is first created.
- `core.IfErrNNotFound(err)` returns true if err is NOT a not-found (inverted name, be careful).

## Custom Frameworks (Critical Knowledge)

### Dove (pkg/dove) — the HTTP layer

**Not** a standard router. Do not assume chi, gorilla, gin, or http.ServeMux patterns.

Key registration (see `pkg/api/routes.go`):
```go
d := dove.New(logHandler)
d.Use(mw1, mw2)                    // applies to endpoints registered AFTER this
d.Prefix("/assets").GET(handler)   // prefix match, checked before exact routes
d.Endpoint("/path").
    GET(h).
    POST(h).
    PUT(h).
    DELETE(h)
return d  // implements http.Handler
```

- `Use()` adds middleware that will be applied to **subsequently registered** endpoints (middleware is captured at `Endpoint()` call time).
- Middleware execution order: last-registered middleware runs first (see `endpoint.addMiddleware` — it walks `middleware` slice backwards).
- `dove.Context`: `Req *http.Request`, `Res *Response`, `Log *oak.Oak`, `StartTime`, `Ctx() context.Context`, `ParseForm()`, `Redirect(url)`, `Render(status, Renderer)`, `RenderOK(Renderer)`, `JSON(status, any)`.
- If handler returns error: logged, 500 written, no further handling.
- `Response.Before(fn)` hooks run before `WriteHeader` (used by session middleware to commit).
- `ServeStaticAssets` reads from `./assets` in dev, `assets.Assets` embedded FS in prod.

**Route order matters**: prefix routes are checked first, then exact routes.

### Oak (pkg/oak) — structured logging

Wraps `slog`. Main entrypoint sets default:
```go
logger := oak.New(oak.TintHandler(os.Stdout, oak.LevelTrace, false))
oak.SetDefault(logger)
```

Usage:
- `oak.Info("msg", "key", val, ...)`
- `l := logger.Layer("app"); l.Info(...); defer l.PopLayer()`
- `oak.FromContext(ctx)` — requires context populated by dove (see `dove/context.go:30`)
- Layers are dot-separated: `"web.app.persistence"`
- Custom levels: Trace, Notice, Warning, Fatal (in addition to Debug/Info/Error)

All production logging should go through oak, never raw `log` or `slog` directly (except in a few legacy spots).

## Configuration and Environment

- **Required at runtime**: `config.json` in working directory. Contains port, database.dsn, integrations, admin/user/guest credentials, environment.
- **Environment variable**: `APE_ENVIRONMENT` (prod/dev/test). Parsed by `application.ParseEnvironment`. Affects session cookie flags (prod → HttpOnly+Secure).
- Integrations are conditionally enabled based on presence of keys/tokens in config.
- Google integration is currently commented out in `pkg/integrations/sync.go`.

**Secrets**: Encrypted via age/sops. See `deployment/enc/`, `scripts/enc_dec.fish`, `scripts/encrypt_sops.sh`, `.sops.yaml`. Mage tasks: `encrypt`, `decrypt`, `enc_sops`.

## Database and Code Generation

- **Driver**: `modernc.org/sqlite` (pure Go, CGO_ENABLED=0 friendly).
- **Connection flags** (hardcoded in `cmd/web/main.go:79`): `?_busy_timeout=5000&_journal_mode=WAL`.
- **Migrations**: golang-migrate. Source is `iofs` embedded from `database/migrations`. Applied at startup.
- **Bob codegen**: `bobgen.yaml` controls generation. Run `mage db:gen`. Output goes to `database/gen/models`, `database/gen/dberrors`.
  - Never hand-edit generated files.
  - `persistence/` code uses bob query builders + `SelectThenLoad` relations.

**Persistence test DBs**: Use `:memory:` + file-based migrations (see `pkg/persistence/auth_test.go:NewTestDB`).

**Application test helper**: `application_test.go:NewTestApplication(t)` — sets up full stack with temp DB.

## Testing Patterns

- Package naming: `package foo_test` (external test) for application tests; `package persistence` for persistence tests.
- Parallel: `t.Parallel()` at start of test and subtests.
- Assertions: custom `pkg/test/test.go` — `test.True(t, cond)`, `test.False`, `test.NilErr(t, err)`.
- Do not use `testing.TB.Error` directly in most cases; use the test helpers.
- For deadline validation errors, tests check `strings.Contains(err.Error(), expected)` because wrapped errors.
- Background: `go tool gotest` (from rakyll/gotest) for colored output; invoked via mage.

## UI and Rendering

- All HTML is Go code via `maragu.dev/gomponents` + `gomponents-datastar` (client-side reactivity, similar to htmx).
- Layout in `pkg/ui/layout.go`. Nav is conditional on `State.Authenticated`.
- State structs in `pkg/ui/ui.go` (e.g., `State`, `DeadlineState`).
- Datastar script loaded from CDN (config constant).
- Assets: dev serves from disk; prod uses embedded FS (`assets/efs.go`).

## Authentication and Sessions

- `scs` (alexedwards) with sqlite3store.
- Cookie: name `curious-ape-session`, 7-day lifetime, SameSite=Strict.
- Middleware chain (order in routes.go):
  1. `MiddlewareLoadCookie` (loads scs, sets up commit hook on Response.Before)
  2. `MiddlewareAuthenticateFromSession` (populates ctx with user + isAuthenticated)
  3. `MiddlewareRequireAuthentication` (redirects to /login if not authed; sets no-store cache)
- Dev gets extra `DevMiddleware` (no-store cache).
- Login form at `/login`; POST creates session; DELETE logs out.
- Admin/user/guest users are upserted at startup from config (passwords are hashed via `bcrypt` presumably in `application/users.go`).

## Integrations and Sync

- `pkg/integrations/sync.go:Integrations` holds clients.
- OAuth2 (fitbit, google): stored in `oauth_token` table; refreshed automatically via `GetHttpClient`.
- Hevy and Toggl are token/key based (no oauth flow in app).
- `DaySync` runs every 6h in background + on-demand via POST `/day/sync`.
- Sync errors are logged but do not fail the whole operation (see `sync.go:25-28`).

## Deployment and Build

- **Docker**: `Dockerfile` — multi-stage, alpine base, `CGO_ENABLED=0`, static binary, tzdata, version injected via ldflags.
- **CI**: `.gitea/workflows/ci-pipeline.yaml` — runs on act runner, installs Go 1.26, mage, tools, builds, tests, audits.
- **Release**: tag `v*` triggers docker build+push to `danicos.dev/daniel/curious-ape`.
- **Kubernetes**: kustomize overlays under `deployment/kubernetes/`. Flux GitOps. Secrets in overlays/config (SOPS encrypted).
- **Linux bare-metal**: systemd unit + litestream sidecar for backup.
- Build script: `scripts/build.fish` (fish shell). Sets version via `mage version`.
- Version is injected at build time via `-X main.version=...` (see Dockerfile + build.fish).

## Gotchas and Non-Obvious Patterns

1. **Date normalization**: Always use `core.NewDate(t)` or `core.NewDateToday()`. Raw `time.Time` will cause comparison bugs.

2. **Inverted error helper**: `core.IfErrNNotFound(err)` returns `true` when err is **not** a not-found. Used to distinguish "no row" from real errors.

3. **dove middleware timing**: `Use()` calls only affect endpoints registered **after** the Use call. Reordering routes can silently change behavior.

4. **dove error handling is minimal**: Handler returning error → immediate 500, error logged. No custom error types or user-facing messages beyond that.

5. **Session writes via Before hook**: Cookie commit happens in `Response.Before`, which fires on first `WriteHeader`. If you write a response without going through Response methods, sessions won't persist.

6. **Bob generated code**: If you change DB schema, you must: edit migration + run `mage migrate:up` (or against dev DB) + `mage db:gen`. Then update persistence layer to match new relations.

7. **Config file is not optional**: `main.go` does `os.ReadFile("config.json")` with no fallback. Missing file = fatal.

8. **Environment affects cookie security**: In prod, session cookie gets HttpOnly+Secure. Dev and test do not.

9. **Background sync is fire-and-forget**: The 6h ticker in main.go spawns `DaySync` but does not wait for or surface errors to the user.

10. **Static asset embedding**: `assets/efs.go` exists. In prod, `ServeStaticAssets` uses `assets.Assets.ReadFile`. Changing assets requires rebuild.

11. **WAL mode + busy timeout**: Hardcoded in main. If you open SQLite elsewhere (tests, mage db:open), you may see locking if you don't use the same flags.

12. **Oak context propagation**: `oak.FromContext(ctx)` only works if dove set it up. Outside HTTP handlers, you may need to use a passed `*oak.Oak` or `oak.NewDefault()`.

13. **Unused code warnings**: There are pre-existing gopls diagnostics (unused params in days.go:28, main.go:240, unused vars, unused funcs in ui). Do not "fix" them unless asked; they may be intentional or deferred cleanup.

14. **No standard log package usage**: Prefer oak everywhere. Some internal bob/migrate code may log via slog, but app code should not.

15. **Migrations are iofs-embedded**: At runtime they come from the compiled binary via `database/migrations/efs.go`. Adding a migration requires rebuild to be visible in a running binary.

16. **Build output location (agents)**: Any direct `go build` performed to test compilation **must** use `-o ./tmp/...`. Binaries must never land in the project root.

## What Not to Do

- Do not edit files under `database/gen/`.
- Do not assume a popular router or logger; check dove/oak first.
- Do not add new external HTTP frameworks or ORMs without explicit discussion.
- Do not bypass `application.App` methods from handlers; keep business logic in the application layer.
- Do not hardcode secrets or tokens in source.
- When using `go build` (or similar) to check for compilation errors, never place the resulting binary in the project root. Always use `-o ./tmp/<name>`.

## Quick File Reference for Common Tasks

| Task | Start Here |
|------|------------|
| Add a new page/route | `pkg/api/routes.go`, add handler in appropriate `api/*.go`, render via `pkg/ui/` |
| Add business logic | `pkg/application/*.go`, call from API handler |
| Add DB table/column | New migration in `database/migrations/sqlite/`, update bobgen if needed, implement repo in `persistence/`, model in `core/` |
| Add integration | `pkg/integrations/<name>/`, wire in `integrations/sync.go`, handle in `application/integrations.go` |
| Change auth behavior | `pkg/api/middleware.go`, `application/users.go` |
| Debug logging | Use `oak.FromContext(ctx).Layer("foo").Debug(...)` or pass logger down |
| Write a test | Look at `application_test.go` (NewTestApplication) or `persistence/*_test.go` (NewTestDB) |
| Touch CI/build | `magefiles/`, `.gitea/workflows/`, `Dockerfile`, `scripts/build.fish` |
