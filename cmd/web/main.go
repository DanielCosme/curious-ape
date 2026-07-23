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
	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v3"
	"github.com/lmittmann/tint"
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
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		Level:      config.Global.LogLevel,
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

	router := chi.NewMux()
	router.Use(
		httplog.RequestLogger(logger, nil),
		middleware.Recoverer,
	)

	db := initDB()
	migrateDB(db)
	sessionManager := initSessionManager(db)
	errGroup, errGroupCtx := errgroup.WithContext(ctx)

	logger.Info("Application initialized",
		"env", config.Global.Environment,
		"version", version,
	)

	// Initialize application modules inside.
	err := api.SetupRoutes(errGroupCtx, router, sessionManager, db)
	if err != nil {
		return fmt.Errorf("error setting up routes: %w", err)
	}

	// TODOs
	// - Integrations Config Loading from secret.
	// - Bob DB init?
	// - Set admin password
	// - Create ticker to sync data daily

	addr := fmt.Sprintf(":%s", config.Global.Port)
	server := &http.Server{
		Addr:     addr,
		Handler:  router,
		ErrorLog: slog.NewLogLogger(slog.Default().Handler(), slog.LevelError),
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
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

func initDB() *sql.DB {
	slog.Info("Opening database", "path", config.Global.DatabasePath)
	dsn := config.Global.DatabasePath + "?_busy_timeout=5000" + "&_journal_mode=WAL"
	db, err := sql.Open("sqlite", dsn)
	exitIfErr(err)
	err = db.Ping()
	exitIfErr(err)
	slog.Info("Database connection established")
	return db
}

func migrateDB(db *sql.DB) {
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

func initSessionManager(db *sql.DB) *scs.SessionManager {
	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(db)
	sessionManager.Lifetime = 24 * time.Hour * 7 // 7 days
	sessionManager.Cookie.SameSite = http.SameSiteStrictMode
	sessionManager.Cookie.Name = "curious-ape-session"
	if config.Global.Environment == config.Prod {
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
