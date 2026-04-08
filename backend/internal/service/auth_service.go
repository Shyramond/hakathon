package service

import (
    "context"
    "errors"
    "fmt"
    "log/slog"
    "time"

    "github.com/Shyramond/hakathon/backend/internal/config"
    "github.com/Shyramond/hakathon/backend/internal/models"
    "github.com/Shyramond/hakathon/backend/internal/repository"
    "github.com/Shyramond/hakathon/backend/pkg"
)

type AuthService struct {
    userRepo     repository.UserRepository
    tokenService *TokenService
    xsollaClient *xsolla.Client
    rewardCfg    config.RewardConfig
}

type LoginResult struct {
    User        *models.User `json:"user"`
    IsNewUser   bool         `json:"is_new_user"`
    RewardGiven bool         `json:"reward_given"`
    Reward      int          `json:"reward"`
    NewBalance  int          `json:"new_balance"`
    Message     string       `json:"message"`
}

func NewAuthService(
    userRepo repository.UserRepository,
    tokenService *TokenService,
    xsollaClient *xsolla.Client,
    rewardCfg config.RewardConfig,
) *AuthService {
    return &AuthService{
        userRepo:     userRepo,
        tokenService: tokenService,
        xsollaClient: xsollaClient,
        rewardCfg:    rewardCfg,
    }
}

// Login обрабатывает вход пользователя:
// 1. Валидирует токен через Xsolla
// 2. Создаёт или находит пользователя в БД
// 3. Если прошло >= 24ч с последнего входа — начисляет случайные токены
func (s *AuthService) Login(ctx context.Context, xsollaToken string) (*LoginResult, error) {
    // 1. Валидация токена и получение профиля
    userInfo, err := s.xsollaClient.ValidateToken(ctx, xsollaToken)
    if err != nil {
        return nil, fmt.Errorf("validate xsolla token: %w", err)
    }

    slog.Info("xsolla user authenticated",
        "xsolla_id", userInfo.ID,
        "username", userInfo.Username,
    )

    // 2. Получаем или создаём пользователя
    user, isNew, err := s.userRepo.GetOrCreate(ctx, userInfo.ID, userInfo.Username, userInfo.Email)
    if err != nil {
        return nil, fmt.Errorf("get or create user: %w", err)
    }

    result := &LoginResult{
        User:      user,
        IsNewUser: isNew,
    }

    // 3. Проверяем cooldown и начисляем награду
    if s.canClaimReward(user) {
        reward, newBalance, err := s.tokenService.CreditLoginReward(ctx, user.ID)
        if err != nil {
            // Не фейлим весь логин из-за ошибки начисления
            slog.Error("failed to credit login reward", "error", err, "user_id", user.ID)
            result.Message = "Login successful, but reward could not be credited"
        } else {
            result.RewardGiven = true
            result.Reward = reward
            result.NewBalance = newBalance
            result.Message = fmt.Sprintf("Welcome! You earned %d tokens", reward)

            // Обновляем данные пользователя в результате
            user.TokenBalance = newBalance
            now := time.Now()
            user.LastLoginAt = &now
        }
    } else {
        result.NewBalance = user.TokenBalance
        nextReward := s.nextRewardTime(user)
        result.Message = fmt.Sprintf("Welcome back! Next reward available at %s", nextReward.Format(time.RFC3339))
    }

    return result, nil
}

// canClaimReward проверяет, прошло ли достаточно времени с последнего логина
func (s *AuthService) canClaimReward(user *models.User) bool {
    // Новый пользователь всегда получает награду
    if user.LastLoginAt == nil {
        return true
    }

    elapsed := time.Since(*user.LastLoginAt)
    return elapsed >= s.rewardCfg.CooldownDuration
}

// nextRewardTime возвращает время, когда можно получить следующую награду
func (s *AuthService) nextRewardTime(user *models.User) time.Time {
    if user.LastLoginAt == nil {
        return time.Now()
    }
    return user.LastLoginAt.Add(s.rewardCfg.CooldownDuration)
}

// GetUserByXsollaID — для middleware, возвращает пользователя по xsolla ID
func (s *AuthService) GetUserByXsollaID(ctx context.Context, xsollaUserID string) (*models.User, error) {
    user, err := s.userRepo.GetByXsollaID(ctx, xsollaUserID)
    if err != nil {
        if errors.Is(err, repository.ErrUserNotFound) {
            return nil, ErrUserNotFound
        }
        return nil, err
    }
    return user, nil
}