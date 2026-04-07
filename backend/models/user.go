package models

import (
	"time"
	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username       string             `bson:"username" json:"username" validate:"required,min=3,max=30"`
	Email          string             `bson:"email" json:"email" validate:"required,email"`
	Password       string             `bson:"password" json:"-"`
	Avatar         string             `bson:"avatar" json:"avatar"`
	IsActive       bool               `bson:"is_active" json:"is_active"`
	LastLogin      time.Time          `bson:"last_login" json:"last_login"`
	LoginStreak    int                `bson:"login_streak" json:"login_streak"`
	TotalLogins    int                `bson:"total_logins" json:"total_logins"`
	Subscription   Subscription       `bson:"subscription" json:"subscription"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
}

type Subscription struct {
	HasSubscription bool   `bson:"has_subscription" json:"has_subscription"`
	SubscriptionID  string `bson:"subscription_id" json:"subscription_id"`
	ExpiresAt       time.Time `bson:"expires_at" json:"expires_at"`
	Type            string `bson:"type" json:"type"` // premium, vip, etc.
}

// NewUser создает нового пользователя с хешированным паролем
func NewUser(username, email, password string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	avatar := "https://i.pravatar.cc/150?u=" + username

	return &User{
		Username:    username,
		Email:       email,
		Password:    string(hashedPassword),
		Avatar:      avatar,
		IsActive:    true,
		LastLogin:   now,
		LoginStreak: 0,
		TotalLogins: 0,
		Subscription: Subscription{
			HasSubscription: false,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// ComparePassword сравнивает пароль с хешем
func (u *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

// UpdateDailyLogin обновляет информацию о ежедневном входе
func (u *User) UpdateDailyLogin() bool {
	now := time.Now()
	lastLogin := u.LastLogin
	
	// Проверяем, был ли вход сегодня
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	lastLoginDay := time.Date(lastLogin.Year(), lastLogin.Month(), lastLogin.Day(), 0, 0, 0, 0, lastLogin.Location())
	
	if today.Equal(lastLoginDay) {
		// Уже заходил сегодня
		return false
	}
	
	// Проверяем, был ли вход вчера
	yesterday := today.AddDate(0, 0, -1)
	
	if lastLoginDay.Equal(yesterday) {
		// Был вход вчера - увеличиваем стрик
		u.LoginStreak++
	} else {
		// Пропустил день - сбрасываем стрик
		u.LoginStreak = 1
	}
	
	u.LastLogin = now
	u.TotalLogins++
	u.UpdatedAt = now
	
	return true
}

// ToJSON возвращает пользователя без пароля
func (u *User) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"id":            u.ID,
		"username":      u.Username,
		"email":         u.Email,
		"avatar":        u.Avatar,
		"is_active":     u.IsActive,
		"last_login":    u.LastLogin,
		"login_streak":  u.LoginStreak,
		"total_logins":  u.TotalLogins,
		"subscription":  u.Subscription,
		"created_at":    u.CreatedAt,
		"updated_at":    u.UpdatedAt,
	}
}
