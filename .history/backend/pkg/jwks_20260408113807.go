package xsolla

import (
    "context"
    "crypto/rsa"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "math/big"
    "net/http"
    "sync"
    "time"
)

type JWKS struct {
    Keys []JWK `json:"keys"`
}

type JWK struct {
    Kty string `json:"kty"`
    Use string `json:"use"`
    Kid string `json:"kid"`
    Alg string `json:"alg"`
    N   string `json:"n"`
    E   string `json:"e"`
}

type JWKSCache struct {
    mu         sync.RWMutex
    keys       map[string]*rsa.PublicKey
    jwksURL    string
    httpClient *http.Client
    lastFetch  time.Time
    ttl        time.Duration
}

func NewJWKSCache(jwksURL string) *JWKSCache {
    return &JWKSCache{
        keys:    make(map[string]*rsa.PublicKey),
        jwksURL: jwksURL,
        httpClient: &http.Client{
            Timeout: 10 * time.Second,
        },
        ttl: 1 * time.Hour, // кешируем на 1 час
    }
}

func (c *JWKSCache) GetKey(ctx context.Context, kid string) (*rsa.PublicKey, error) {
    c.mu.RLock()
    key, exists := c.keys[kid]
    expired := time.Since(c.lastFetch) > c.ttl
    c.mu.RUnlock()

    if exists && !expired {
        return key, nil
    }

    if err := c.refresh(ctx); err != nil {
        if exists {
            return key, nil
        }
        return nil, fmt.Errorf("refresh JWKS: %w", err)
    }

    c.mu.RLock()
    key, exists = c.keys[kid]
    c.mu.RUnlock()

    if !exists {
        return nil, fmt.Errorf("key %s not found in JWKS", kid)
    }

    return key, nil
}

func (c *JWKSCache) refresh(ctx context.Context) error {
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.jwksURL, nil)
    if err != nil {
        return fmt.Errorf("create request: %w", err)
    }

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return fmt.Errorf("fetch JWKS: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("JWKS endpoint returned %d", resp.StatusCode)
    }

    var jwks JWKS
    if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
        return fmt.Errorf("decode JWKS: %w", err)
    }

    keys := make(map[string]*rsa.PublicKey)
    for _, jwk := range jwks.Keys {
        if jwk.Kty != "RSA" {
            continue
        }

        pubKey, err := jwkToRSAPublicKey(jwk)
        if err != nil {
            continue
        }

        keys[jwk.Kid] = pubKey
    }

    c.mu.Lock()
    c.keys = keys
    c.lastFetch = time.Now()
    c.mu.Unlock()

    return nil
}

func jwkToRSAPublicKey(jwk JWK) (*rsa.PublicKey, error) {
    nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
    if err != nil {
        return nil, fmt.Errorf("decode N: %w", err)
    }

    eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
    if err != nil {
        return nil, fmt.Errorf("decode E: %w", err)
    }

    n := new(big.Int).SetBytes(nBytes)
    e := 0
    for _, b := range eBytes {
        e = e<<8 + int(b)
    }

    return &rsa.PublicKey{
        N: n,
        E: e,
    }, nil
}