package config

import (
    "fmt"
    "os"
    "strconv"
    "time"
)

type Config struct {
    Server    ServerConfig
    Database  DatabaseConfig
    Xsolla    XsollaConfig
    Reward    RewardConfig
    RateLimit RateLimitConfig
    Idempotency IdempotencyConfig
}

type ServerConfig struct {
    Port    string
    GinMode string
}

type DatabaseConfig struct {
    URL string
}

type XsollaConfig struct {
    ProjectID        string
    LoginProjectID   string
    APIKey           string
    OAuth2ClientID   string
    OAuth2Secret     string
    Issuer           string
}

type RewardConfig struct {
    Min            int
    Max            int
    CooldownHours  int
    CooldownDuration time.Duration
}

type RateLimitConfig struct {
    RPS   float64
    Burst int
}

type IdempotencyConfig struct {
    TTL time.Duration
}

func Load() (*Config, error) {
    cfg := &Config{
        Server: ServerConfig{
            Port:    getEnvOrDefault("SERVER_PORT", "8080"),
            GinMode: getEnvOrDefault("GIN_MODE", "debug"),
        },
        Database: DatabaseConfig{
            URL: getEnvOrDefault("DATABASE_URL", "postgres://hakathon:hakathon_secret@localhost:5432/hakathon?sslmode=disable"),
        },
        Xsolla: XsollaConfig{
            ProjectID:      getEnvOrDefault("XSOLLA_PROJECT_ID", "XXXXXX"),
            LoginProjectID: getEnvOrDefault("XSOLLA_LOGIN_PROJECT_ID", "XXXXXX"),
            APIKey:         getEnvOrDefault("XSOLLA_API_KEY", "XXXXXX"),
            OAuth2ClientID: getEnvOrDefault("XSOLLA_OAUTH2_CLIENT_ID", "XXXXXX"),
            OAuth2Secret:   getEnvOrDefault("XSOLLA_OAUTH2_CLIENT_SECRET", "XXXXXX"),
            Issuer:         getEnvOrDefault("XSOLLA_ISSUER", "https://login.xsolla.com"),
        },
    }

    rewardMin, err := strconv.Atoi(getEnvOrDefault("LOGIN_REWARD_MIN", "10"))
    if err != nil {
        return nil, fmt.Errorf("invalid LOGIN_REWARD_MIN: %w", err)
    }
    rewardMax, err := strconv.Atoi(getEnvOrDefault("LOGIN_REWARD_MAX", "50"))
    if err != nil {
        return nil, fmt.Errorf("invalid LOGIN_REWARD_MAX: %w", err)
    }
    cooldownH, err := strconv.Atoi(getEnvOrDefault("LOGIN_REWARD_COOLDOWN_HOURS", "24"))
    if err != nil {
        return nil, fmt.Errorf("invalid LOGIN_REWARD_COOLDOWN_HOURS: %w", err)
    }

    cfg.Reward = RewardConfig{
        Min:              rewardMin,
        Max:              rewardMax,
        CooldownHours:    cooldownH,
        CooldownDuration: time.Duration(cooldownH) * time.Hour,
    }

    rps, _ := strconv.ParseFloat(getEnvOrDefault("RATE_LIMIT_RPS", "20"), 64)
    burst, _ := strconv.Atoi(getEnvOrDefault("RATE_LIMIT_BURST", "40"))
    cfg.RateLimit = RateLimitConfig{RPS: rps, Burst: burst}

    ttlSec, _ := strconv.Atoi(getEnvOrDefault("IDEMPOTENCY_TTL_SECONDS", "86400"))
    cfg.Idempotency = IdempotencyConfig{TTL: time.Duration(ttlSec) * time.Second}

    return cfg, nil
}

func getEnvOrDefault(key, fallback string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return fallback
}