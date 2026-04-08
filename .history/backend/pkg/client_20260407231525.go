package xsolla

import (
    "context"
    "encoding/json"
    "fmt"
    "log/slog"
    "net/http"
    "strings"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

type Client struct {
    ProjectID      string
    LoginProjectID string
    APIKey         string
    OAuth2ClientID string
    OAuth2Secret   string
    Issuer         string

    jwksCache  *JWKSCache
    httpClient *http.Client

    // devMode — если true, пропускаем реальную валидацию JWT
    devMode bool
}

type Config struct {
    ProjectID      string
    LoginProjectID string
    APIKey         string
    OAuth2ClientID string
    OAuth2Secret   string
    Issuer         string
}

func NewClient(cfg Config) *Client {
    issuer := cfg.Issuer
    if issuer == "" {
        issuer = "https://login.xsolla.com"
    }

    jwksURL := fmt.Sprintf("%s/api/oauth2/keys", issuer)

    devMode := cfg.ProjectID == "XXXXXX" || cfg.OAuth2ClientID == "XXXXXX"
    if devMode {
        slog.Warn("xsolla client running in DEV MODE — JWT validation disabled")
    }

    return &Client{
        ProjectID:      cfg.ProjectID,
        LoginProjectID: cfg.LoginProjectID,
        APIKey:         cfg.APIKey,
        OAuth2ClientID: cfg.OAuth2ClientID,
        OAuth2Secret:   cfg.OAuth2Secret,
        Issuer:         issuer,
        jwksCache:      NewJWKSCache(jwksURL),
        httpClient: &http.Client{
            Timeout: 10 * time.Second,
        },
        devMode: devMode,
    }
}

func (c *Client) IsDevMode() bool {
    return c.devMode
}


func (c *Client) ValidateToken(ctx context.Context, tokenString string) (*UserInfo, error) {
    if c.devMode {
        return c.devValidateToken(tokenString)
    }

    return c.prodValidateToken(ctx, tokenString)
}

func (c *Client) devValidateToken(tokenString string) (*UserInfo, error) {
    parser := jwt.NewParser(jwt.WithoutClaimsValidation())
    token, _, err := parser.ParseUnverified(tokenString, &XsollaClaims{})
    if err == nil {
        if claims, ok := token.Claims.(*XsollaClaims); ok {
            return &UserInfo{
                ID:       claims.Subject,
                Email:    claims.Email,
                Username: claims.Username,
                Nickname: claims.Nickname,
            }, nil
        }
    }

    slog.Debug("dev mode: using token string as user ID", "token", tokenString[:min(len(tokenString), 20)])
    return &UserInfo{
        ID:       tokenString,
        Email:    fmt.Sprintf("%s@dev.local", tokenString),
        Username: fmt.Sprintf("dev_%s", tokenString),
    }, nil
}

func (c *Client) prodValidateToken(ctx context.Context, tokenString string) (*UserInfo, error) {
    token, err := jwt.ParseWithClaims(tokenString, &XsollaClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }

        kid, ok := token.Header["kid"].(string)
        if !ok {
            return nil, fmt.Errorf("missing kid in token header")
        }

        key, err := c.jwksCache.GetKey(ctx, kid)
        if err != nil {
            return nil, fmt.Errorf("get signing key: %w", err)
        }

        return key, nil
    },
        jwt.WithIssuer(c.Issuer),
        jwt.WithExpirationRequired(),
        jwt.WithValidMethods([]string{"RS256"}),
    )

    if err != nil {
        return nil, fmt.Errorf("invalid token: %w", err)
    }

    claims, ok := token.Claims.(*XsollaClaims)
    if !ok || !token.Valid {
        return nil, fmt.Errorf("invalid token claims")
    }

    userInfo := &UserInfo{
        ID:       claims.Subject,
        Email:    claims.Email,
        Username: claims.Username,
        Nickname: claims.Nickname,
    }

    if userInfo.Username == "" {
        userInfo.Username = userInfo.Nickname
    }
    if userInfo.Username == "" && userInfo.Email != "" {
        parts := strings.SplitN(userInfo.Email, "@", 2)
        userInfo.Username = parts[0]
    }

    return userInfo, nil
}

// GetUserProfile получает профиль через Xsolla Login API (server-to-server).
// GET https://login.xsolla.com/api/users/me
func (c *Client) GetUserProfile(ctx context.Context, accessToken string) (*UserInfo, error) {
    if c.devMode {
        return c.devValidateToken(accessToken)
    }

    url := fmt.Sprintf("%s/api/users/me", c.Issuer)

    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    if err != nil {
        return nil, fmt.Errorf("create request: %w", err)
    }
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
    req.Header.Set("Content-Type", "application/json")

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("fetch user profile: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("xsolla API returned status %d", resp.StatusCode)
    }

    var profile struct {
        ID       string `json:"id"`
        Email    string `json:"email"`
        Username string `json:"username"`
        Nickname string `json:"nickname"`
        Name     string `json:"name"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
        return nil, fmt.Errorf("decode profile: %w", err)
    }

    username := profile.Username
    if username == "" {
        username = profile.Nickname
    }
    if username == "" {
        username = profile.Name
    }

    return &UserInfo{
        ID:       profile.ID,
        Email:    profile.Email,
        Username: username,
        Nickname: profile.Nickname,
    }, nil
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}