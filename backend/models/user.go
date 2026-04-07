package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username    string             `bson:"username" json:"username" validate:"required,min=3,max=30"`
	Email       string             `bson:"email" json:"email" validate:"required,email"`
	Password    string             `bson:"password" json:"-"`
	Avatar      string             `bson:"avatar" json:"avatar"`
	IsActive    bool               `bson:"is_active" json:"is_active"`
	LastLogin   time.Time          `bson:"last_login" json:"last_login"`
	TotalLogins int                `bson:"total_logins" json:"total_logins"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
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
		TotalLogins: 0,
		CreatedAt:   now,
		UpdatedAt:   now,
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

	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	lastLoginDay := time.Date(lastLogin.Year(), lastLogin.Month(), lastLogin.Day(), 0, 0, 0, 0, lastLogin.Location())

	if today.Equal(lastLoginDay) {
		return false
	}

	u.LastLogin = now
	u.TotalLogins++
	u.UpdatedAt = now

	return true
}

// ToJSON возвращает пользователя без пароля
func (u *User) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"id":           u.ID,
		"username":     u.Username,
		"email":        u.Email,
		"avatar":       u.Avatar,
		"is_active":    u.IsActive,
		"last_login":   u.LastLogin,
		"total_logins": u.TotalLogins,
		"created_at":   u.CreatedAt,
		"updated_at":   u.UpdatedAt,
	}
}
