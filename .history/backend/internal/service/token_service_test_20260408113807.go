package service

import (
    "context"
    "testing"
    "time"

    "github.com/google/uuid"
    "github.com/Shyramond/hakathon/backend/internal/config"
    "github.com/Shyramond/hakathon/backend/internal/mocks"
    "github.com/Shyramond/hakathon/backend/internal/models"
)

func newTestTokenService() (*TokenService, *mocks.UserRepoMock, *mocks.TransactionRepoMock, *mocks.TxManagerMock) {
    userRepo := mocks.NewUserRepoMock()
    txRepo := mocks.NewTransactionRepoMock()
    txMgr := mocks.NewTxManagerMock()

    svc := NewTokenService(userRepo, txRepo, txMgr, config.RewardConfig{
        Min:              10,
        Max:              50,
        CooldownDuration: 24 * time.Hour,
    })

    return svc, userRepo, txRepo, txMgr
}

func TestGenerateLoginReward_InRange(t *testing.T) {
    svc, _, _, _ := newTestTokenService()

    for i := 0; i < 100; i++ {
        reward := svc.GenerateLoginReward()
        if reward < 10 || reward > 50 {
            t.Errorf("reward %d out of range [10, 50]", reward)
        }
    }
}

func TestGenerateLoginReward_EqualMinMax(t *testing.T) {
    svc := NewTokenService(nil, nil, nil, config.RewardConfig{Min: 25, Max: 25})

    for i := 0; i < 50; i++ {
        if r := svc.GenerateLoginReward(); r != 25 {
            t.Errorf("expected 25, got %d", r)
        }
    }
}

func TestCreditLoginReward_Success(t *testing.T) {
    svc, userRepo, txRepo, txMgr := newTestTokenService()

    userID := uuid.New()
    userRepo.AddUser(&models.User{
        ID:           userID,
        XsollaUserID: "xid_1",
        TokenBalance: 100,
    })

    reward, newBalance, err := svc.CreditLoginReward(context.Background(), userID)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if reward < 10 || reward > 50 {
        t.Errorf("reward %d out of range", reward)
    }

    if newBalance != 100+reward {
        t.Errorf("expected balance %d, got %d", 100+reward, newBalance)
    }

    if txRepo.CreateCalls != 1 {
        t.Errorf("expected 1 transaction create call, got %d", txRepo.CreateCalls)
    }

    if txMgr.RunInTxCalls != 1 {
        t.Errorf("expected 1 RunInTx call, got %d", txMgr.RunInTxCalls)
    }

    updated := userRepo.GetUser(userID)
    if updated.TokenBalance != newBalance {
        t.Errorf("expected stored balance %d, got %d", newBalance, updated.TokenBalance)
    }
}

func TestCreditLoginReward_UserNotFound(t *testing.T) {
    svc, _, _, _ := newTestTokenService()

    _, _, err := svc.CreditLoginReward(context.Background(), uuid.New())
    if err == nil {
        t.Fatal("expected error for non-existent user")
    }
}

func TestGetBalance_Success(t *testing.T) {
    svc, userRepo, _, _ := newTestTokenService()

    userID := uuid.New()
    userRepo.AddUser(&models.User{
        ID:           userID,
        TokenBalance: 42,
    })

    balance, err := svc.GetBalance(context.Background(), userID)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if balance != 42 {
        t.Errorf("expected 42, got %d", balance)
    }
}