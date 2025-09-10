package service

import (
	"errors"
	"folder-system/internal/config"
	"folder-system/internal/entity"
	"folder-system/internal/repository"
	"folder-system/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(email, password string) error
	Login(email, password string) (accessToken, refreshToken string, err error)
	RefreshTokens(refreshToken string) (newAccessToken, newRefreshToken string, err error)
}

type authService struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) AuthService {
	return &authService{userRepo: userRepo, cfg: cfg}
}

func (s *authService) Register(email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &entity.User{
		Email:    email,
		Password: string(hashedPassword),
	}

	return s.userRepo.CreateUser(user)
}

func (s *authService) Login(email, password string) (string, string, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	accessToken, err := utils.GenerateJWT(user.ID, utils.AccessToken, s.cfg.JWT.AccessSecret, s.cfg.JWT.AccessTTL)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := utils.GenerateJWT(user.ID, utils.RefreshToken, s.cfg.JWT.RefreshSecret, s.cfg.JWT.RefreshTTL)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *authService) RefreshTokens(refreshToken string) (string, string, error) {
	claims, err := utils.ParseJWT(refreshToken, s.cfg.JWT.RefreshSecret)
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}
	if claims.Subject != string(utils.RefreshToken) {
		return "", "", errors.New("token is not a refresh token")
	}

	newAccessToken, err := utils.GenerateJWT(claims.UserID, utils.AccessToken, s.cfg.JWT.AccessSecret, s.cfg.JWT.AccessTTL)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := utils.GenerateJWT(claims.UserID, utils.RefreshToken, s.cfg.JWT.RefreshSecret, s.cfg.JWT.RefreshTTL)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}
