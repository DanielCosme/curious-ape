package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"time"

	"danicos.dev/daniel/curious-ape/database/migrations"
	"danicos.dev/daniel/curious-ape/pkg/api"
	"danicos.dev/daniel/curious-ape/pkg/config"
	"danicos.dev/daniel/curious-ape/pkg/integrations"
	embbeded_nats "danicos.dev/daniel/curious-ape/pkg/nats"
	"danicos.dev/daniel/curious-ape/pkg/services/day"
	"danicos.dev/daniel/curious-ape/pkg/services/fitnesslog"
	"danicos.dev/daniel/curious-ape/pkg/services/habit"
	"danicos.dev/daniel/curious-ape/pkg/services/integration"
	"danicos.dev/daniel/curious-ape/pkg/services/sleeplog"
	"danicos.dev/daniel/curious-ape/pkg/services/user"
	"danicos.dev/daniel/curious-ape/pkg/services/worklog"
	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v3"
	"github.com/lmittmann/tint"
	"github.com/nats-io/nats.go"
	"github.com/stephenafamo/bob"
	"golang.org/x/oauth2"
	"golang.org/x/sync/errgroup"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "modernc.org/sqlite"
)

func main() {
	ctx := context.Background()
	err := run(ctx)
	if err != nil {
		slog.Error("server failure", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()

	cfg := config.Load()
	cfg.Validate() // Will panic if fails validation.

	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		Level:      cfg.LogLevel,
		TimeFormat: time.TimeOnly,
	}))
	slog.SetDefault(logger)

	version := config.Version()
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "version":
			fmt.Println(version)
			return nil
		default:
			slog.Warn("Unkouwn command received", "key", os.Args[1])
		}
	}

	ns, err := embbeded_nats.New(ctx)
	ns.WaitForServer()
	exitIfErr(err)
	nc, err := ns.Client()
	exitIfErr(err)

	db := initDB(cfg)
	bobDB := bob.NewDB(db)
	migrateDB(db)
	sessionManager := initSessionManager(cfg, db)
	errGroup, errGroupCtx := errgroup.WithContext(ctx)

	fitbitConfig := &oauth2.Config{
		ClientID:     cfg.Fitbit.ClientID,
		ClientSecret: cfg.Fitbit.ClientSecret,
		RedirectURL:  cfg.Fitbit.RedirectURL,
		Scopes:       cfg.Fitbit.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  cfg.Fitbit.AuthURL,
			TokenURL: cfg.Fitbit.TokenURL,
		},
	}
	is := integrations.New(cfg.TogglWorkspaceID, cfg.TogglAPIKey, cfg.HevyAPIKey, fitbitConfig, nil)

	userService := user.NewService(bobDB)
	integrationService := integration.NewService(is, bobDB, nc)
	handlers := api.NewHandlers(
		sessionManager,
		day.NewService(bobDB, nc),
		habit.NewService(bobDB, nc),
		integrationService,
		worklog.NewService(bobDB, nc, integrationService),
		userService,
		fitnesslog.NewService(bobDB, nc, integrationService),
		sleeplog.NewService(bobDB, nc, integrationService),
	)

	router := chi.NewMux()
	router.Use(
		httplog.RequestLogger(logger, &httplog.Options{
			Schema: &httplog.Schema{RequestRemoteIP: "ip"},
		}),
		middleware.Recoverer,
	)
	if err := api.SetupRouter(errGroupCtx, handlers, router, bobDB); err != nil {
		return fmt.Errorf("error setting up routes: %w", err)
	}
	if err := userService.SetPassword(cfg.Username, cfg.Password); err != nil {
		return fmt.Errorf("error setting up username/password: %w", err)
	}

	logger.Info("Application initialized",
		"env", cfg.Environment,
		"version", version,
	)

	addr := fmt.Sprintf(":%s", cfg.Port)
	server := &http.Server{
		Addr:     addr,
		Handler:  router,
		ErrorLog: slog.NewLogLogger(slog.Default().Handler(), slog.LevelError),
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}

	if config.IsDev() {
		go func() {
			nc.Subscribe(">", func(msg *nats.Msg) {
				slog.Debug("Event emited", "subject", msg.Subject)
			})
		}()
	}

	errGroup.Go(func() error {
		slog.Info("Server listening", "addr", server.Addr)
		err := server.ListenAndServe()
		if err == nil || err == http.ErrServerClosed {
			return nil
		}
		return fmt.Errorf("server error: %w", err)
	})

	errGroup.Go(func() error {
		<-errGroupCtx.Done()
		slog.Info("Shutting down server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("Server shutdown error: %w", err)
		}
		return nil
	})

	return errGroup.Wait()
}

func initDB(cfg *config.Config) *sql.DB {
	slog.Info("Opening database", "path", cfg.DatabasePath)
	options := "?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=synchronous(NORMAL)&_pragma=foreign_keys(ON)"
	dsn := cfg.DatabasePath + options
	db, err := sql.Open("sqlite", dsn)
	exitIfErr(err)
	err = db.Ping()
	exitIfErr(err)
	slog.Info("Database connection established")
	return db
}

func migrateDB(db *sql.DB) {
	slog.Info("Applying migrations")
	migrationsSource, err := iofs.New(migrations.Migrations, "sqlite")
	exitIfErr(err)
	migrationDriver, err := sqlite.WithInstance(db, &sqlite.Config{})
	exitIfErr(err)
	migrator, err := migrate.NewWithInstance("iofs", migrationsSource, "sqlite", migrationDriver)
	exitIfErr(err)
	err = migrator.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			v, dirty, _ := migrator.Version()
			slog.Info("Migrations up to date", "version", v, "dirty", dirty)
		} else {
			exitIfErr(err)
		}
	}
}

func initSessionManager(cfg *config.Config, db *sql.DB) *scs.SessionManager {
	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(db)
	sessionManager.Lifetime = 24 * time.Hour * 7 // 7 days
	sessionManager.Cookie.SameSite = http.SameSiteStrictMode
	sessionManager.Cookie.Name = "curious-ape-session"
	if cfg.Environment == config.Prod {
		sessionManager.Cookie.HttpOnly = true
		sessionManager.Cookie.Secure = true
	}
	return sessionManager
}

func exitIfErr(err error) {
	if err != nil {
		logFatal(err)
	}
}

func logFatal(err error) {
	slog.Error("Fatal failure", "err", err.Error(), "stack", string(debug.Stack()))
	os.Exit(1)
}
