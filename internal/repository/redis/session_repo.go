package redis

import (
	"context"
	"fmt"
	"time"

	"api/internal/domain"

	goredis "github.com/redis/go-redis/v9"
)

type sessionRepository struct {
	client *goredis.Client
}

func NewSessionRepository(client *goredis.Client) domain.SessionRepository {
	return &sessionRepository{client: client}
}

func (r *sessionRepository) StoreRefreshToken(ctx context.Context, userID int64, tokenID string, expiresAt time.Duration) error {
	key := fmt.Sprintf("refresh_token:%d:%s", userID, tokenID)
	return r.client.Set(ctx, key, "active", expiresAt).Err()
}

func (r *sessionRepository) RevokeRefreshToken(ctx context.Context, userID int64, tokenID string) error {
	key := fmt.Sprintf("refresh_token:%d:%s", userID, tokenID)
	return r.client.Del(ctx, key).Err()
}

func (r *sessionRepository) IsTokenRevoked(ctx context.Context, userID int64, tokenID string) (bool, error) {
	key := fmt.Sprintf("refresh_token:%d:%s", userID, tokenID)
	err := r.client.Get(ctx, key).Err()
	if err == goredis.Nil {
		return true, nil // Key doesn't exist = revoked/expired
	} else if err != nil {
		return false, err
	}
	return false, nil
}

func (r *sessionRepository) BlacklistAccessToken(ctx context.Context, tokenID string, expiresAt time.Duration) error {
	key := fmt.Sprintf("blacklist:%s", tokenID)
	return r.client.Set(ctx, key, "blacklisted", expiresAt).Err()
}

func (r *sessionRepository) IsAccessTokenBlacklisted(ctx context.Context, tokenID string) (bool, error) {
	key := fmt.Sprintf("blacklist:%s", tokenID)
	err := r.client.Get(ctx, key).Err()
	if err == goredis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
