package service

import (
	"errors"
	"time"

	"github.com/Aakash-Sleur/go-micro-auth/config"
	"github.com/Aakash-Sleur/go-micro-auth/models"
	"github.com/Aakash-Sleur/go-micro-auth/repository"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *repository.UserRepository
}

func NewAuthService(repo *repository.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

type CustomClaims struct {
	UserID uuid.UUID `json:"userId"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

func generateJWT(user *models.User) (string, error) {
	secret := config.Load().JWT_SECRET
	if secret == "" {
		return "", errors.New("JWT_SECRET is not set")
	}

	claims := CustomClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "go-micro-auth",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (s *AuthService) SignUp(email, name string, password string) (string, *models.User, error) {
	existingUser, err := s.repo.FindByEmail(email)

	if err != nil {
		return "", nil, err
	}

	if existingUser != nil {
		return "", nil, errors.New("email already in use")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil, errors.New("failed to hash password")
	}

	user := &models.User{
		ID:        uuid.New(),
		Name:      name,
		Email:     email,
		Password:  string(hashed),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(user); err != nil {
		return "", nil, err
	}

	token, err := generateJWT(user)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (s *AuthService) Signin(email, password string) (string, *models.User, error) {
	userDetails, err := s.repo.FindByEmail(email)

	if err != nil {
		return "", nil, err
	}

	if userDetails == nil {
		return "", nil, errors.New("user not available in the db")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userDetails.Password), []byte(password)); err != nil {
		return "", nil, err
	}

	token, err := generateJWT(userDetails)

	if err != nil {
		return "", nil, err
	}

	return token, userDetails, nil
}
