package mocks

import (
    "context"
    "fmt"

    "github.com/Shyramond/hakathon/backend/pkg"
)

type XsollaClientMock struct {
    ValidTokens map[string]*xsolla.UserInfo
    ShouldFail  bool
    FailError   error
}

func NewXsollaClientMock() *XsollaClientMock {
    return &XsollaClientMock{
        ValidTokens: make(map[string]*xsolla.UserInfo),
    }
}

func (m *XsollaClientMock) AddValidToken(token string, info *xsolla.UserInfo) {
    m.ValidTokens[token] = info
}

func (m *XsollaClientMock) ValidateToken(ctx context.Context, token string) (*xsolla.UserInfo, error) {
    if m.ShouldFail {
        if m.FailError != nil {
            return nil, m.FailError
        }
        return nil, fmt.Errorf("mock validation failed")
    }

    info, ok := m.ValidTokens[token]
    if !ok {
        return nil, fmt.Errorf("invalid token")
    }
    return info, nil
}