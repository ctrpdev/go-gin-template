package domain

import (
	"context"
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"-"`
	Role     string `json:"role"`
	Verified bool   `json:"verified"`
	BaseModel
}

type UserRepository interface {
	CreateUser(ctx context.Context, email, passwordHash string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	VerifyUser(ctx context.Context, id int64) error
	GetUserByID(ctx context.Context, id int64) (*User, error)
}

type UserService interface {
	Register(ctx context.Context, email, password string) (*User, error)
	Login(ctx context.Context, email, password string) (accessToken, refreshToken string, err error)
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
	Logout(ctx context.Context, userID int64, accessTokenID, refreshTokenID string) error
	VerifyAccount(ctx context.Context, id int64) error
	GetMe(ctx context.Context, id int64) (*User, error)
}
