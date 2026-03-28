package postgres

import (
	"context"
	"errors"

	"api/internal/domain"
	domainerr "api/internal/errors"
	"api/internal/repository/postgres/db"

	"github.com/jackc/pgx/v5/pgconn"
)

type userRepository struct {
	queries *db.Queries
}

func NewUserRepository(queries *db.Queries) domain.UserRepository {
	return &userRepository{queries: queries}
}

func (r *userRepository) CreateUser(ctx context.Context, email, passwordHash string) (*domain.User, error) {
	row, err := r.queries.CreateUser(ctx, db.CreateUserParams{
		Email:        email,
		PasswordHash: passwordHash,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, domainerr.ErrUserAlreadyExists
		}
		return nil, err
	}

	return &domain.User{
		Email:    row.Email,
		Role:     row.Role,
		Verified: row.Verified,
		BaseModel: domain.BaseModel{
			ID:        row.ID,
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
		},
	}, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	row, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return &domain.User{
		Email:    row.Email,
		Password: row.PasswordHash,
		Role:     row.Role,
		Verified: row.Verified,
		BaseModel: domain.BaseModel{
			ID:        row.ID,
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
		},
	}, nil
}

func (r *userRepository) VerifyUser(ctx context.Context, id int64) error {
	rowsAffected, err := r.queries.VerifyUser(ctx, id)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	row, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &domain.User{
		Email:    row.Email,
		Role:     row.Role,
		Verified: row.Verified,
		BaseModel: domain.BaseModel{
			ID:        row.ID,
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
		},
	}, nil
}
