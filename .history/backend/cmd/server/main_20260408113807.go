package main

import (
    "context"
    "errors"
    "fmt"
    "log/slog"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/Shyramond/hakathon/backend/internal/config"
    "github.com/Shyramond/hakathon/backend/internal/database"
    "github.com/Shyramond/hakathon/backend/internal/repository"
    "github.com/Shyramond/hakathon/backend/internal/router"
    "github.com/Shyramond/hakathon/backend/internal/service"
    "github.com/Shyramond/hakathon/backend/pkg/xsolla"
)

func main() {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
    }))
    slog.SetDefault(logger)

    cfg, err := config.Load()
    if err != nil {
        slog.Error("failed to load config", "error", err)
        os.Exit(1)
    }

    migrationsPath := getMigrationsPath()
    if err := database.RunMigrations(cfg.Database.URL, migrationsPath); err != nil {
        slog.Error("failed to run migrations", "error", err)
        os.Exit(1)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    pool, err := database.NewPostgresPool(ctx, cfg.Database.URL)
    if err != nil {
        slog.Error("failed to connect to database", "error", err)
        os.Exit(1)
    }
    defer pool.Close()

    repos := repository.NewRepos(pool, cfg.Idempotency.TTL)

    xsollaClient := xsolla.NewClient(xsolla.Config{
        ProjectID:      cfg.Xsolla.ProjectID,
        LoginProjectID: cfg.Xsolla.LoginProjectID,
        APIKey:         cfg.Xsolla.APIKey,
        OAuth2ClientID: cfg.Xsolla.OAuth2ClientID,
        OAuth2Secret:   cfg.Xsolla.OAuth2Secret,
        Issuer:         cfg.Xsolla.Issuer,
    })

    services := service.NewServices(repos, cfg, xsollaClient)

    go startIdempotencyCleanup(repos.Idempotency)

    deps := router.Dependencies{
        Config:       cfg,
        Repos:        repos,
        Services:     services,
        XsollaClient: xsollaClient,
    }
    engine := router.Setup(deps)

    srv := &http.Server{
        Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
        Handler:      engine,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    go func() {
        slog.Info("starting server",
            "port", cfg.Server.Port,
            "gin_mode", cfg.Server.GinMode,
            "xsolla_dev_mode", xsollaClient.IsDevMode(),
        )
        if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
            slog.Error("server error", "error", err)
            os.Exit(1)
        }
    }()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    slog.Info("shutting down server...")

    shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer shutdownCancel()

    if err := srv.Shutdown(shutdownCtx); err != nil {
        slog.Error("server forced to shutdown", "error", err)
    }

    slog.Info("server stopped")
}

func getMigrationsPath() string {
    if _, err := os.Stat("/migrations"); err == nil {
        return "/migrations"
    }
    return "internal/database/migrations"
}

func startIdempotencyCleanup(repo repository.IdempotencyRepository) {
    ticker := time.NewTicker(1 * time.Hour)
    defer ticker.Stop()

    for range ticker.C {
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        if err := repo.Cleanup(ctx); err != nil {
            slog.Error("idempotency cleanup failed", "error", err)
        }
        cancel()
    }
}