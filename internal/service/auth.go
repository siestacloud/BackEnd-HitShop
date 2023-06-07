package service

import (
	"crypto/sha1"
	"errors"
	"hitshop/internal/core"
	"hitshop/internal/repository"

	"fmt"

	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	salt       = "hjqrhjqw124617ajfhajs"        // соль добавляемая к паролю пользователей
	signingKey = "qrkjk#4#%35FSFJlja#4353KSFjH" // Набор случайных байт для подписи токена (ключ подписи) - так-же исп при расшифровке токена
	tokenTTL   = 12 * time.Hour                 // время жизни токена
)

type tokenClaims struct {
	jwt.StandardClaims
	AccountID uuid.UUID `json:"user_id"`
}

// Авторизация и аутентификация
type AuthService struct {
	repo repository.Authorization
}

// NewAuthService конструктор
func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Test() {
	logrus.Info("info in auth")
	logrus.WithFields(logrus.Fields{"tag": "a tag svc"}).Info("An info message")
	s.repo.TestDB()
}

// CreateAccount
func (s *AuthService) CreateAccount(acc *core.Account) (uuid.UUID, error) {
	acc.Password = generatePasswordHash(acc.Password)
	return s.repo.CreateAccount(acc)
}

// GetAccountByCode
func (s *AuthService) GetAccountByCode(verification_code string) (*core.Account, error) {
	return s.repo.GetAccountByCode(verification_code)
}

// GetAccountByCode
func (s *AuthService) GetAccountByEmail(email, password string) (*core.Account, error) {
	return s.repo.GetAccountByEmail(email, password)
}

// GetAccountByCode
func (s *AuthService) GetAccountByUUID(UUID uuid.UUID) (*core.Account, error) {
	return s.repo.GetAccountByUUID(UUID)
}

// GetAccountByCode
func (s *AuthService) UpdateAccount(acc *core.Account) (uuid.UUID, error) {
	return s.repo.UpdateAccount(acc)
}

// Для генерации токена нужно получить пользователя из базы
// если пользователя нет, вернуть ошибку
// в токен записывается id пользователя
func (s *AuthService) GenerateToken(email, password string) (string, error) {
	acc, err := s.repo.GetAccountByEmail(email, generatePasswordHash(password))
	if err != nil {
		return "", err
	}
	if !acc.Verified {
		return "", errors.New(fmt.Sprintf("Email %s not verify yet", acc.Email))
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(), // токен перестает быть валидным через
			IssuedAt:  time.Now().Unix(),
		},
		acc.UUID,
	})

	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) ChangePassword(UUID uuid.UUID, password, passwordNew string) (uuid.UUID, error) {

	acc, err := s.repo.GetAccountByUUID(UUID)
	if err != nil {
		return uuid.UUID{}, err
	}
	if generatePasswordHash(password) != acc.Password {
		return uuid.UUID{}, errors.New("passwords didn't match")
	}

	acc.Password = generatePasswordHash(passwordNew)
	return s.repo.UpdateAccount(acc)
}

// Используется middleware
// Достаем Id пользователя из токена
func (s *AuthService) ParseToken(accessToken string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}

			return []byte(signingKey), nil
		})
	if err != nil {
		return uuid.UUID{}, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return uuid.UUID{}, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.AccountID, nil
}

// generatePasswordHash генерирует хеш, добавляем соль, перчим
func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
