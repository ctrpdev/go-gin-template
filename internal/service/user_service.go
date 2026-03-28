package service

import (
	"api/internal/errors"
	"context"
	"strconv"
	"time"

	"api/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo    domain.UserRepository
	sessionRepo domain.SessionRepository
	jwtSecret   []byte
}

func NewUserService(userRepo domain.UserRepository, sessionRepo domain.SessionRepository, jwtSecret []byte) domain.UserService {
	return &userService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		jwtSecret:   jwtSecret,
	}
}

func (s *userService) Register(ctx context.Context, email, password string) (*domain.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return s.userRepo.CreateUser(ctx, email, string(hash))
}

func (s *userService) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", "", errors.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", errors.ErrInvalidCredentials
	}

	if !user.Verified {
		return "", "", errors.ErrAccountNotVerified
	}

	accessTokenID := uuid.New().String()
	refreshTokenID := uuid.New().String()

	// Generate Access Token (15m)
	accessClaims := jwt.MapClaims{
		"sub":  strconv.FormatInt(user.ID, 10),
		"role": user.Role,
		"jti":  accessTokenID,
		"exp":  time.Now().Add(15 * time.Minute).Unix(),
	}
	accessToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(s.jwtSecret)

	// Generate Refresh Token (7d)
	refreshClaims := jwt.MapClaims{
		"sub": strconv.FormatInt(user.ID, 10),
		"jti": refreshTokenID,
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	refreshToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(s.jwtSecret)

	// Persist refresh token in Redis
	err = s.sessionRepo.StoreRefreshToken(ctx, user.ID, refreshTokenID, 7*24*time.Hour)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *userService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	// Parse refresh token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.ErrInvalidToken
		}
		return s.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return "", "", errors.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.ErrInvalidToken
	}

	jti, ok := claims["jti"].(string)
	if !ok {
		return "", "", errors.ErrInvalidToken
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return "", "", errors.ErrInvalidToken
	}

	userID, err := strconv.ParseInt(sub, 10, 64)
	if err != nil {
		return "", "", errors.ErrInvalidToken
	}

	// Check if refresh token is valid in Redis
	revoked, err := s.sessionRepo.IsTokenRevoked(ctx, userID, jti)
	if err != nil || revoked {
		return "", "", errors.ErrTokenRevoked
	}

	// Invalidate old refresh token (rotation)
	_ = s.sessionRepo.RevokeRefreshToken(ctx, userID, jti)

	// User must still exist
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return "", "", errors.ErrUserNotFound
	}

	// Generate new tokens
	accessTokenID := uuid.New().String()
	newRefreshTokenID := uuid.New().String()

	accessClaims := jwt.MapClaims{
		"sub":  strconv.FormatInt(user.ID, 10),
		"role": user.Role,
		"jti":  accessTokenID,
		"exp":  time.Now().Add(15 * time.Minute).Unix(),
	}
	newAccessToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(s.jwtSecret)

	refreshClaims := jwt.MapClaims{
		"sub": strconv.FormatInt(user.ID, 10),
		"jti": newRefreshTokenID,
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	newRefreshString, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(s.jwtSecret)

	// Persist new refresh token
	err = s.sessionRepo.StoreRefreshToken(ctx, user.ID, newRefreshTokenID, 7*24*time.Hour)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshString, nil
}

func (s *userService) Logout(ctx context.Context, userID int64, accessTokenID, refreshTokenID string) error {
	// Revoke Refresh Token
	_ = s.sessionRepo.RevokeRefreshToken(ctx, userID, refreshTokenID)

	// Blacklist Access Token for its remaining lifetime (approximate max 15m)
	return s.sessionRepo.BlacklistAccessToken(ctx, accessTokenID, 15*time.Minute)
}

func (s *userService) VerifyAccount(ctx context.Context, id int64) error {
	return s.userRepo.VerifyUser(ctx, id)
}

func (s *userService) GetMe(ctx context.Context, id int64) (*domain.User, error) {
	return s.userRepo.GetUserByID(ctx, id)
}
