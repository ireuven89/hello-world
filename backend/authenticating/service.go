package authenticating

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/ireuven89/hello-world/backend/authenticating/model"
)

type Service interface {
	Register(username, password string) error
	Login(username, password string) (string, error)
	VerifyToken(tokenString string) (string, error)
}

type AuthRepo interface {
	Save(username, password string) error
	Find(username string) (model.User, error)
}

// AuthService is the core authenticating service
type AuthService struct {
	userStore AuthRepo
	logger    *zap.Logger
}

var jwtSecretKey = []byte("your_secret_key")

// NewAuthService creates a new AuthService
func NewAuthService(userStore AuthRepo, logger *zap.Logger) *AuthService {
	return &AuthService{userStore: userStore, logger: logger}
}

// Register registers a new user
func (service *AuthService) Register(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		service.logger.Error("failed to register with error", zap.Error(err))
		return err
	}
	return service.userStore.Save(username, string(hashedPassword))
}

// Login authenticates a model and returns a JWT token
func (service *AuthService) Login(username, password string) (string, error) {
	user, err := service.userStore.Find(username)
	if err != nil {
		service.logger.Error("failed to login with error", zap.Error(err))
		return "", err
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(), // Token expires in 24 hours
	})

	return token.SignedString(jwtSecretKey)
}

// VerifyToken verifies and decodes a JWT token
func (service *AuthService) VerifyToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecretKey, nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return "", errors.New("username not found in token")
	}

	return username, nil
}
