package xsolla

import "github.com/golang-jwt/jwt/v5"

// XsollaClaims — структура JWT claims от Xsolla Login
type XsollaClaims struct {
    jwt.RegisteredClaims
    Email      string `json:"email,omitempty"`
    Username   string `json:"username,omitempty"`
    Nickname   string `json:"nickname,omitempty"`
    XsollaSubType string `json:"xsolla_login_project_id,omitempty"`
    Type       string `json:"type,omitempty"`
    // Xsolla может класть дополнительные поля
    Groups     []GroupClaim `json:"groups,omitempty"`
    Ismaster   bool         `json:"is_master,omitempty"`
    PromoEmailAgreement *int `json:"promo_email_agreement,omitempty"`
}

type GroupClaim struct {
    ID        int    `json:"id"`
    Name      string `json:"name"`
    IsDefault bool   `json:"is_default"`
}

// UserInfo — нормализованный профиль пользователя
type UserInfo struct {
    ID       string `json:"id"`
    Email    string `json:"email"`
    Username string `json:"username"`
    Nickname string `json:"nickname"`
}

// TokenResponse — ответ от OAuth2 token endpoint
type TokenResponse struct {
    AccessToken  string `json:"access_token"`
    TokenType    string `json:"token_type"`
    ExpiresIn    int    `json:"expires_in"`
    RefreshToken string `json:"refresh_token,omitempty"`
}