package mocks

import (
"context"
"api/internal/domain"
"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
mock.Mock
}

func (m *MockUserRepository) CreateUser(ctx context.Context, email, passwordHash string) (*domain.User, error) {
args := m.Called(ctx, email, passwordHash)
var user *domain.User
if args.Get(0) != nil {
user = args.Get(0).(*domain.User)
}
return user, args.Error(1)
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
args := m.Called(ctx, email)
var user *domain.User
if args.Get(0) != nil {
user = args.Get(0).(*domain.User)
}
return user, args.Error(1)
}

func (m *MockUserRepository) VerifyUser(ctx context.Context, id int64) error {
args := m.Called(ctx, id)
return args.Error(0)
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
args := m.Called(ctx, id)
var user *domain.User
if args.Get(0) != nil {
user = args.Get(0).(*domain.User)
}
return user, args.Error(1)
}
