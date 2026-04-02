//# бизнес-логика (хеширование, JWT)

package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/AlexHeadGU/todolist/internal/config"
	"github.com/AlexHeadGU/todolist/internal/models"
	"github.com/AlexHeadGU/todolist/internal/repository"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

// Register регистрирует нового пользователя
func (s *AuthService) Register(email, password string) (*models.User, error) {
	// Проверяем, существует ли пользователь
	exists, err := s.userRepo.ExistsByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists {
		return nil, errors.New("user already exists")
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Создаём пользователя
	user, err := s.userRepo.Create(email, string(hashedPassword))
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// Login выполняет вход пользователя
func (s *AuthService) Login(email, password string) (string, *models.User, error) {
	// Ищем пользователя
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return "", nil, errors.New("invalid credentials")
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	// Генерируем JWT
	token, err := s.generateJWT(user)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return token, user, nil
}

// generateJWT создаёт JWT токен
func (s *AuthService) generateJWT(user *models.User) (string, error) {
	// Создаём claims (утверждения)
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // токен живёт 24 часа
		"iat":     time.Now().Unix(),                     // время выпуска
	}

	// Создаём токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен секретным ключом
	secret := []byte(config.AppConfig.JWTSecret)
	return token.SignedString(secret)
}

// ValidateToken проверяет валидность токена
func (s *AuthService) ValidateToken(tokenString string) (int, error) {
	secret := []byte(config.AppConfig.JWTSecret)

	// Парсим токен
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed to parse token: %w", err)
	}

	// Проверяем валидность и извлекаем user_id
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(float64)
		if !ok {
			return 0, errors.New("invalid token claims")
		}
		return int(userID), nil
	}

	return 0, errors.New("invalid token")
}
