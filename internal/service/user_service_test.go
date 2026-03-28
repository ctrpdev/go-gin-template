package service

import (
	"context"
	"testing"
	"time"

	"api/internal/domain"
	"api/internal/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestRegister_Success(t *testing.T) {
	// Setup
	mockUserRepo := new(mocks.MockUserRepository)
	// Como no estamos probando la sesi�n aqu�, podemos pasar nil o un mock si lo us�ramos
	userService := NewUserService(mockUserRepo, nil, []byte("secret"))

	email := "test@example.com"
	password := "password123"

	// Retorno esperado simulado
	expectedUser := &domain.User{
		Email:    email,
		Role:     "user",
		Verified: false,
		BaseModel: domain.BaseModel{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Mockeamos la respuesta cuando CreateUser sea llamado
	// Usamos mock.AnythingOfType("string") para el passwordHash porque al ser bcrypt generado din�micamente cada vez,
	// el hash cambia, no podemos hacer match exacto de ese string.
	mockUserRepo.On("CreateUser", mock.Anything, email, mock.AnythingOfType("string")).Return(expectedUser, nil)

	// Ejecuci�n
	user, err := userService.Register(context.Background(), email, password)

	// Aserciones (Verificamos que el resultado es el esperado)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, int64(1), user.BaseModel.ID)

	// Verificamos que nuestro Mock fue llamado correctamente con los par�metros simulados.
	mockUserRepo.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	// Setup
	mockUserRepo := new(mocks.MockUserRepository)
	mockSessionRepo := new(mocks.MockSessionRepository)
	userService := NewUserService(mockUserRepo, mockSessionRepo, []byte("test_secret"))

	email := "test@example.com"
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	expectedUser := &domain.User{
		Email:     email,
		Password:  string(hashedPassword),
		Role:      "user",
		BaseModel: domain.BaseModel{ID: 1},
	}

	// Mocking
	mockUserRepo.On("GetUserByEmail", mock.Anything, email).Return(expectedUser, nil)
	mockSessionRepo.On("StoreRefreshToken", mock.Anything, int64(1), mock.AnythingOfType("string"), mock.Anything).Return(nil)

	// Execute
	accToken, refToken, err := userService.Login(context.Background(), email, password)

	// Assertions
	assert.NoError(t, err)
	assert.NotEmpty(t, accToken)
	assert.NotEmpty(t, refToken)

	mockUserRepo.AssertExpectations(t)
	mockSessionRepo.AssertExpectations(t)
}
