package auth

import (
	"context"
	"errors"
	"github.com/Khasmag06/kode-notes/internal/models"
	"github.com/Khasmag06/kode-notes/internal/repository/repo_errs"
	"github.com/Khasmag06/kode-notes/pkg/app_err"
	"strconv"
)

const (
	invalidLoginOrPasswordErr = "неверный логин или пароль"
	loginAlreadyExistsErr     = "пользователь с таким логином уже существует"
)

type service struct {
	repo   repository
	hasher hasher
	jwt
	decoder decoder
}

func New(repo repository, hasher hasher, jwt jwt, decoder decoder) *service {
	return &service{
		repo:    repo,
		hasher:  hasher,
		jwt:     jwt,
		decoder: decoder,
	}
}

func (s *service) SignUp(ctx context.Context, user models.User) error {
	passwordHash, err := s.hasher.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = passwordHash

	if err = s.repo.CreateUser(ctx, user); err != nil {
		if errors.Is(err, repo_errs.ErrAlreadyExists) {
			return app_err.NewConflictError(loginAlreadyExistsErr)
		}
		return err
	}
	return nil
}

func (s *service) Login(ctx context.Context, loginData models.User) (*models.TokensResponse, error) {
	user, err := s.repo.GetUserByLogin(ctx, loginData.Login)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, app_err.NewAuthorizationError(invalidLoginOrPasswordErr)
	}

	passwordMatch := s.hasher.CheckPasswordHash(loginData.Password, user.Password)
	if !passwordMatch {
		return nil, app_err.NewAuthorizationError(invalidLoginOrPasswordErr)
	}
	userID := strconv.Itoa(user.Id)
	userHash, err := s.decoder.Encrypt([]byte(userID))
	if err != nil {
		return nil, err
	}
	var tokenData models.TokensResponse

	tokenData.AccessToken, err = s.jwt.GenerateToken(userHash)
	if err != nil {
		return nil, err
	}
	return &tokenData, nil
}
