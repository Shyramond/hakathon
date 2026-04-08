package service

import (
    "context"
    "testing"
    "time"

    "github.com/google/uuid"
    "github.com/Shyramond/hakathon/backend/internal/config"
    "github.com/Shyramond/hakathon/backend/internal/mocks"
    "github.com/Shyramond/hakathon/backend/internal/models"
    "github.com/Shyramond/hakathon/backend/pkg"
)

type authTestEnv struct {
    authService  *AuthService
    userRepo     *mocks.UserRepoMock
    txRepo       *mocks.TransactionRepoMock
    txMgr        *mocks.TxManagerMock
    xsollaClient *mocks.XsollaClientMock
}

func newAuthTestEnv() *authTestEnv {
    userRepo := mocks.NewUserRepoMock()
    txRepo := mocks.NewTransactionRepoMock()
    txMgr := mocks.NewTxManagerMock()
    xsollaClient := mocks.NewXsollaClientMock()

    rewardCfg := config.RewardConfig{
        Min:              10,
        Max:              50,
        CooldownDuration: 24 * time.Hour,
    }

    tokenService := NewTokenService(userRepo, txRepo, txMgr, rewardCfg)

    // Создаём реальный xsolla.Client в dev mode для тестов
    realClient := xsolla.NewClient(xsolla.Config{
        ProjectID:      "XXXXXX",
        OAuth2ClientID: "XXXXXX",
    })

    authService := NewAuthService(userRepo, tokenService, realClient, rewardCfg)

    return &authTestEnv{
        authService:  authService,
        userRepo:     userRepo,
        txRepo:       txRepo,
        txMgr:        txMgr,
        xsollaClient: xsollaClient,
    }
}

func TestLogin_NewUser(t *testing.T) {
    env := newAuthTestEnv()

    result, err := env.authService.Login(context.Background(), "test_user_123")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if !result.IsNewUser {
        t.Error("expected IsNewUser to be true")
    }

    if !result.RewardGiven {
        t.Error("expected reward to be given to new user")
    }

    if result.Reward < 10 || result.Reward > 50 {
        t.Errorf("reward %d out of range [10, 50]", result.Reward)
    }

    if result.User == nil {
        t.Fatal("expected user to be non-nil")
    }
}

func TestLogin_ExistingUser_CooldownNotExpired(t *testing.T) {
    env := newAuthTestEnv()

    // Создаём пользователя с недавним логином
    recentLogin := time.Now().Add(-1 * time.Hour) // 1 час назад
    user := &models.User{
        ID:           uuid.New(),
        XsollaUserID: "existing_user",
        Username:     "dev_existing_user",
        Email:        "[email protected]",
        TokenBalance: 100,
        LastLoginAt:  &recentLogin,
    }
    env.userRepo.AddUser(user)

    result, err := env.authService.Login(context.Background(), "existing_user")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if result.IsNewUser {
        t.Error("expected IsNewUser to be false")
    }

    if result.RewardGiven {
        t.Error("expected no reward due to cooldown")
    }

    if result.NewBalance != 100 {
        t.Errorf("expected balance 100, got %d", result.NewBalance)
    }
}

func TestLogin_ExistingUser_CooldownExpired(t *testing.T) {
    env := newAuthTestEnv()

    // Логин > 24ч назад
    oldLogin := time.Now().Add(-25 * time.Hour)
    user := &models.User{
        ID:           uuid.New(),
        XsollaUserID: "old_user",
        Username:     "dev_old_user",
        Email:        "[email protected]",
        TokenBalance: 50,
        LastLoginAt:  &oldLogin,
    }
    env.userRepo.AddUser(user)

    result, err := env.authService.Login(context.Background(), "old_user")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if result.RewardGiven != true {
        t.Error("expected reward after cooldown expired")
    }

    if result.Reward < 10 || result.Reward > 50 {
        t.Errorf("reward %d out of range", result.Reward)
    }

    if result.NewBalance != 50+result.Reward {
        t.Errorf("expected balance %d, got %d", 50+result.Reward, result.NewBalance)
    }
}

func TestCanClaimReward_TableDriven(t *testing.T) {
    env := newAuthTestEnv()

    now := time.Now()
    hour1Ago := now.Add(-1 * time.Hour)
    hour23Ago := now.Add(-23 * time.Hour)
    hour25Ago := now.Add(-25 * time.Hour)
    day3Ago := now.Add(-72 * time.Hour)

    tests := []struct {
        name     string
        user     *models.User
        expected bool
    }{
        {
            name:     "nil last_login (new user)",
            user:     &models.User{LastLoginAt: nil},
            expected: true,
        },
        {
            name:     "1 hour ago",
            user:     &models.User{LastLoginAt: &hour1Ago},
            expected: false,
        },
        {
            name:     "23 hours ago",
            user:     &models.User{LastLoginAt: &hour23Ago},
            expected: false,
        },
        {
            name:     "25 hours ago",
            user:     &models.User{LastLoginAt: &hour25Ago},
            expected: true,
        },
        {
            name:     "3 days ago",
            user:     &models.User{LastLoginAt: &day3Ago},
            expected: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := env.authService.canClaimReward(tt.user)
            if result != tt.expected {
                t.Errorf("expected %v, got %v", tt.expected, result)
            }
        })
    }
}