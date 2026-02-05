package service

import (
	"context"
	"errors"
	"log"
	"time"

	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthService interface {
	Register(ctx context.Context, req models.RegisterRequest) (*models.User, error)
	Login(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (*models.User, error)
	GenerateToken(user *models.User) (string, error)
	ValidateToken(tokenString string) (*models.User, error)
	ListUsers(ctx context.Context, page, limit, rangeDays int) ([]models.User, int, error)
	UpdateUserRole(ctx context.Context, userID uuid.UUID, role string) (*models.User, string, error)
}

type authService struct {
	userRepo  repository.UserRepository
	jwtSecret string
	jwtExpiry time.Duration
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string, jwtExpiry time.Duration) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
	}
}

func (s *authService) Register(ctx context.Context, req models.RegisterRequest) (*models.User, error) {
	log.Printf("📝 Register called with email: %s", req.Email)

	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		log.Printf("❌ GetByEmail error: %v", err)
		return nil, err
	}

	if existingUser != nil {
		log.Printf("⚠️ User already exists: %s", req.Email)
		return nil, errors.New("user already exists")
	}

	log.Printf("✅ User doesn't exist, proceeding with creation: %s", req.Email)

	// Hash password
	hashedPassword, err := models.HashPassword(req.Password)
	if err != nil {
		log.Printf("❌ Password hash error: %v", err)
		return nil, err
	}

	// Create user
	user := &models.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         "customer",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	log.Printf("🔄 Calling Create with user ID: %v", user.ID)
	err = s.userRepo.Create(ctx, user)
	if err != nil {
		log.Printf("❌ Create error: %v", err)
		return nil, err
	}

	log.Printf("✅ User created successfully: %v", user.ID)

	// Don't return password hash
	user.PasswordHash = ""
	return user, nil
}

func (s *authService) Login(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	// Check password
	if !models.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid email or password")
	}

	// Generate token
	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	// Don't return password hash
	user.PasswordHash = ""

	return &models.LoginResponse{
		User:        user,
		AccessToken: token,
	}, nil
}

func (s *authService) GetProfile(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	// Don't return password hash
	user.PasswordHash = ""
	return user, nil
}

func (s *authService) GenerateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(s.jwtExpiry).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *authService) ValidateToken(tokenString string) (*models.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			return nil, errors.New("invalid token claims")
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return nil, err
		}

		email, _ := claims["email"].(string)
		role, _ := claims["role"].(string)

		return &models.User{
			ID:    userID,
			Email: email,
			Role:  role,
		}, nil
	}

	return nil, errors.New("invalid token")
}

func (s *authService) ListUsers(ctx context.Context, page, limit, rangeDays int) ([]models.User, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	users, total, err := s.userRepo.GetAll(ctx, page, limit, rangeDays)
	if err != nil {
		return nil, 0, err
	}

	for i := range users {
		users[i].PasswordHash = ""
	}

	return users, total, nil
}

func (s *authService) UpdateUserRole(ctx context.Context, userID uuid.UUID, role string) (*models.User, string, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, "", err
	}
	if user == nil {
		return nil, "", errors.New("user not found")
	}

	if err := s.userRepo.UpdateRole(ctx, userID, role); err != nil {
		return nil, "", err
	}

	updated, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, "", err
	}
	if updated != nil {
		updated.PasswordHash = ""
	}

	token, err := s.GenerateToken(updated)
	if err != nil {
		return nil, "", err
	}

	return updated, token, nil
}
