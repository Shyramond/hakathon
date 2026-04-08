package service

import (
    "context"
    "fmt"
    "log/slog"
    "math/rand"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5"
    "github.com/Shyramond/hakathon/backend/internal/config"
    "github.com/Shyramond/hakathon/backend/internal/models"
    "github.com/Shyramond/hakathon/backend/internal/repository"
)

type TokenService struct {
    userRepo    repository.UserRepository
    txRepo      repository.TransactionRepository
    txManager   repository.TxManager
    rewardCfg   config.RewardConfig
}

func NewTokenService(
    userRepo repository.UserRepository,
    txRepo repository.TransactionRepository,
    txManager repository.TxManager,
    rewardCfg config.RewardConfig,
) *TokenService {
    return &TokenService{
        userRepo:  userRepo,
        txRepo:    txRepo,
        txManager: txManager,
        rewardCfg: rewardCfg,
    }
}


func (s *TokenService) GenerateLoginReward() int {
    if s.rewardCfg.Min >= s.rewardCfg.Max {
        return s.rewardCfg.Min
    }
    return s.rewardCfg.Min + rand.Intn(s.rewardCfg.Max-s.rewardCfg.Min+1)
}

func (s *TokenService) CreditLoginReward(ctx context.Context, userID uuid.UUID) (int, int, error) {
    var reward int
    var newBalance int

    err := s.txManager.RunInTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
        user, err := s.userRepo.GetByIDTx(ctx, tx, userID)
        if err != nil {
            return fmt.Errorf("get user: %w", err)
        }

        reward = s.GenerateLoginReward()
        newBalance = user.TokenBalance + reward

        if err := s.userRepo.UpdateBalanceTx(ctx, tx, userID, newBalance); err != nil {
            return fmt.Errorf("update balance: %w", err)
        }

        if err := s.userRepo.UpdateLastLoginTx(ctx, tx, userID); err != nil {
            return fmt.Errorf("update last login: %w", err)
        }

        t := &models.TokenTransaction{
            UserID:      userID,
            Amount:      reward,
            Type:        models.TransactionLoginReward,
            Description: fmt.Sprintf("Login reward: +%d tokens", reward),
        }
        if err := s.txRepo.CreateTx(ctx, tx, t); err != nil {
            return fmt.Errorf("create transaction: %w", err)
        }

        slog.Info("login reward credited",
            "user_id", userID,
            "reward", reward,
            "new_balance", newBalance,
        )

        return nil
    })

    if err != nil {
        return 0, 0, err
    }

    return reward, newBalance, nil
}

// DeductTokens списывает токены (транзакционно, с проверкой баланса).
func (s *TokenService) DeductTokens(ctx context.Context, tx pgx.Tx, userID uuid.UUID, amount int, txType models.TransactionType, referenceID *uuid.UUID, description string) (int, error) {
    user, err := s.userRepo.GetByIDTx(ctx, tx, userID)
    if err != nil {
        return 0, fmt.Errorf("get user: %w", err)
    }

    if user.TokenBalance < amount {
        return 0, ErrInsufficientBalance
    }

    newBalance := user.TokenBalance - amount

    if err := s.userRepo.UpdateBalanceTx(ctx, tx, userID, newBalance); err != nil {
        return 0, fmt.Errorf("update balance: %w", err)
    }

    t := &models.TokenTransaction{
        UserID:      userID,
        Amount:      -amount,
        Type:        txType,
        ReferenceID: referenceID,
        Description: description,
    }
    if err := s.txRepo.CreateTx(ctx, tx, t); err != nil {
        return 0, fmt.Errorf("create transaction: %w", err)
    }

    return newBalance, nil
}

func (s *TokenService) GetBalance(ctx context.Context, userID uuid.UUID) (int, error) {
    user, err := s.userRepo.GetByID(ctx, userID)
    if err != nil {
        return 0, err
    }
    return user.TokenBalance, nil
}