package domain

import (
	"context"
	"time"
)

type SessionRepository interface {
	StoreRefreshToken(ctx context.Context, userID int64, tokenID string, expiresAt time.Duration) error
	RevokeRefreshToken(ctx context.Context, userID int64, tokenID string) error
	IsTokenRevoked(ctx context.Context, userID int64, tokenID string) (bool, error)
	BlacklistAccessToken(ctx context.Context, tokenID string, expiresAt time.Duration) error
	IsAccessTokenBlacklisted(ctx context.Context, tokenID string) (bool, error)
}
