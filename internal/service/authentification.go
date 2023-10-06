package service

import (
	"context"
	"errors"
	"time"

	"github.com/Be1chenok/testTaskMicroservice/internal/domain"
	"github.com/Be1chenok/testTaskMicroservice/internal/model"
	"github.com/Be1chenok/testTaskMicroservice/internal/repository/postgres"
	rdb "github.com/Be1chenok/testTaskMicroservice/internal/repository/redis"
	"github.com/Be1chenok/testTaskMicroservice/pkg/hash"
	"github.com/golang-jwt/jwt"
)

type Authentification interface {
	SignUp(input model.SignUpInput) (int, error)
	SignIn(ctx context.Context, input model.SignInInput) (model.Tokens, error)

	RefreshTokens(ctx context.Context, refreshToken string) (model.Tokens, error)
	ParseToken(ctx context.Context, accesToken string) (int, error)

	SignOut(ctx context.Context, accessToken string) error
	FullSignOut(ctx context.Context, accesToken string) error
}

type AuthentificationService struct {
	postgresUser     postgres.User
	redisUser        rdb.User
	hashser          *hash.SHA256Hasher
	tokensSigningKey string
	accessTokenTTL   time.Duration
	refreshTokenTTL  time.Duration
}

func NewAuthentificationService(
	postgresUser postgres.User,
	redisUser rdb.User,
	hasher *hash.SHA256Hasher,
	tokensSigningKey string,
	accessTokenTTL, refreshTokenTTL time.Duration) *AuthentificationService {

	return &AuthentificationService{
		postgresUser:     postgresUser,
		redisUser:        redisUser,
		hashser:          hasher,
		tokensSigningKey: tokensSigningKey,
		accessTokenTTL:   accessTokenTTL,
		refreshTokenTTL:  refreshTokenTTL,
	}
}

func (s *AuthentificationService) SignUp(input model.SignUpInput) (int, error) {
	passwordHash, err := s.hashser.Hash(input.Password)
	if err != nil {
		return 0, err
	}

	user := domain.User{
		Email:    input.Email,
		Username: input.Username,
		Password: passwordHash,
	}

	return s.postgresUser.CreateUser(user)
}

func (s *AuthentificationService) SignIn(ctx context.Context, input model.SignInInput) (model.Tokens, error) {
	passwordHash, err := s.hashser.Hash(input.Password)
	if err != nil {
		return model.Tokens{}, err
	}

	userId, err := s.postgresUser.GetUserId(input.Username, passwordHash)
	if err != nil {
		return model.Tokens{}, err
	}

	return s.createSession(ctx, userId)
}

func (s *AuthentificationService) SignOut(ctx context.Context, accessToken string) error {
	return s.deleteUserIdByAccessToken(ctx, accessToken)
}

func (s *AuthentificationService) FullSignOut(ctx context.Context, accesToken string) error {
	userId, err := s.getUserIdByAccessToken(ctx, accesToken)
	if err != nil {
		return err
	}

	return s.deleteAllTokensByUserId(ctx, userId)
}

func (s *AuthentificationService) RefreshTokens(ctx context.Context, refreshToken string) (model.Tokens, error) {
	userId, err := s.postgresUser.GetUserIdByRefreshToken(refreshToken)
	if err != nil {
		return model.Tokens{}, err
	}
	err = s.postgresUser.DeleteUserIdByRefreshToken(refreshToken)
	if err != nil {
		return model.Tokens{}, err
	}
	return s.createSession(ctx, userId)
}

func (s *AuthentificationService) ParseToken(ctx context.Context, token string) (int, error) {
	userId, err := s.getUserIdByAccessToken(ctx, token)
	if err != nil {
		userId, err = s.postgresUser.GetUserIdByRefreshToken(token)
		if err != nil {
			return 0, err
		}
	}

	parsedToken, err := jwt.ParseWithClaims(token, &model.AccesTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(s.tokensSigningKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := parsedToken.Claims.(*model.AccesTokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	if claims.UserId != userId {
		return 0, errors.New("token does not match the stored user ID")
	}

	return claims.UserId, nil
}

func (s *AuthentificationService) createSession(ctx context.Context, userId int) (model.Tokens, error) {
	accessToken, err := s.newToken(userId, s.accessTokenTTL)
	if err != nil {
		return model.Tokens{}, err
	}

	refreshToken, err := s.newToken(userId, s.refreshTokenTTL)
	if err != nil {
		return model.Tokens{}, err
	}

	if err = s.setTokens(ctx, userId, accessToken, refreshToken); err != nil {
		return model.Tokens{}, err
	}

	return model.Tokens{
			AccesToken:   accessToken,
			RefreshToken: refreshToken},
		nil
}

func (s *AuthentificationService) newToken(userId int, tokenTTL time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &model.AccesTokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId: userId,
	})

	return token.SignedString([]byte(s.tokensSigningKey))
}

func (s *AuthentificationService) getUserIdByAccessToken(ctx context.Context, token string) (int, error) {
	userId, err := s.redisUser.GetUserIdByAccessToken(ctx, token)
	if err != nil {
		userId, err = s.postgresUser.GetUserIdByAccessToken(token)
		if err != nil {
			return 0, err
		}
	}

	return userId, nil
}

func (s *AuthentificationService) setTokens(ctx context.Context, userId int, accessToken, refreshToken string) error {
	if err := s.postgresUser.SetTokens(userId, accessToken, refreshToken); err != nil {
		return err
	}

	if err := s.redisUser.SetAccessToken(ctx, accessToken, userId, s.accessTokenTTL); err != nil {
		return err
	}

	return nil
}

func (s *AuthentificationService) deleteAllTokensByUserId(ctx context.Context, userId int) error {
	if err := s.postgresUser.DeleteAllTokensByUserId(userId); err != nil {
		return err
	}

	if err := s.redisUser.DeleteAllAccessTokensByUserId(ctx, userId); err != nil {
		return err
	}

	return nil
}

func (s *AuthentificationService) deleteUserIdByAccessToken(ctx context.Context, accessToken string) error {
	if err := s.postgresUser.DeleteUserIdByAccessToken(accessToken); err != nil {
		return err
	}

	if err := s.redisUser.DeleteUserIdByAccessToken(ctx, accessToken); err != nil {
		return err
	}

	return nil
}
