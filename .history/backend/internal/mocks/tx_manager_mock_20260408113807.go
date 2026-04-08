package mocks

import (
    "context"

    "github.com/jackc/pgx/v5"
)

type TxManagerMock struct {
    RunInTxCalls int
    ForceError error
}

func NewTxManagerMock() *TxManagerMock {
    return &TxManagerMock{}
}

func (m *TxManagerMock) RunInTx(ctx context.Context, fn func(ctx context.Context, tx pgx.Tx) error) error {
    m.RunInTxCalls++

    if m.ForceError != nil {
        return m.ForceError
    }

    return fn(ctx, nil)
}