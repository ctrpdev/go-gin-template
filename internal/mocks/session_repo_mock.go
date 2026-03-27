package mocks

import (
"context"
"time"

"github.com/stretchr/testify/mock"
)

type MockSessionRepository struct {
mock.Mock
}

func (m *MockSessionRepository) StoreRefreshToken(ctx context.Context, userID int64, tokenID string, expiresAt time.Duration) error {
args := m.Called(ctx, userID, tokenID, expiresAt)
return args.Error(0)
}

func (m *MockSessionRepository) RevokeRefreshToken(ctx context.Context, userID int64, tokenID string) error {
args := m.Called(ctx, userID, tokenID)
return args.Error(0)
}

func (m *MockSessionRepository) IsTokenRevoked(ctx context.Context, userID int64, tokenID string) (bool, error) {
args := m.Called(ctx, userID, tokenID)
return args.Bool(0), args.Error(1)
}

func (m *MockSessionRepository) BlacklistAccessToken(ctx context.Context, tokenID string, expiresAt time.Duration) error {
args := m.Called(ctx, tokenID, expiresAt)
return args.Error(0)
}

func (m *MockSessionRepository) IsAccessTokenBlacklisted(ctx context.Context, tokenID string) (bool, error) {
args := m.Called(ctx, tokenID)
return args.Bool(0), args.Error(1)
}
