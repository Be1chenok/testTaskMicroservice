package service

import (
	"context"
	"errors"
	"time"

	"github.com/Be1chenok/testTaskMicroservice/internal/domain"
	"github.com/Be1chenok/testTaskMicroservice/internal/repository"
	"github.com/Be1chenok/testTaskMicroservice/pkg/hash"
	"github.com/golang-jwt/jwt"
)

type AuthentificationService struct {
	postgresUserRepo repository.PostgresUser
	redisUserRepo    repository.RedisUser
	hashser          hash.PasswordHash
	tokensSigningKey string
	accessTokenTTL   time.Duration
	refreshTokenTTL  time.Duration
}

func NewAuthentificationService(
	postgresUserRepo repository.PostgresUser,
	redisUserRepo repository.RedisUser,
	hasher hash.PasswordHash,
	tokensSigningKey string,
	accessTokenTTL, refreshTokenTTL time.Duration) *AuthentificationService {

	return &AuthentificationService{
		postgresUserRepo: postgresUserRepo,
		redisUserRepo:    redisUserRepo,
		hashser:          hasher,
		tokensSigningKey: tokensSigningKey,
		accessTokenTTL:   accessTokenTTL,
		refreshTokenTTL:  refreshTokenTTL,
	}
}

func (s *AuthentificationService) SignUp(input SignUpInput) (int, error) {
	passwordHash, err := s.hashser.Hash(input.Password)
	if err != nil {
		return 0, err
	}

	user := domain.User{
		Email:    input.Email,
		Username: input.Username,
		Password: passwordHash,
	}

	return s.postgresUserRepo.CreateUser(user)
}

func (s *AuthentificationService) SignIn(ctx context.Context, input SignInInput) (Tokens, error) {
	passwordHash, err := s.hashser.Hash(input.Password)
	if err != nil {
		return Tokens{}, err
	}

	userId, err := s.postgresUserRepo.GetUserId(input.Username, passwordHash)
	if err != nil {
		return Tokens{}, err
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

func (s *AuthentificationService) RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error) {
	userId, err := s.postgresUserRepo.GetUserIdByRefreshToken(refreshToken)
	if err != nil {
		return Tokens{}, err
	}
	err = s.postgresUserRepo.DeleteUserIdByRefreshToken(refreshToken)
	if err != nil {
		return Tokens{}, err
	}
	return s.createSession(ctx, userId)
}

func (s *AuthentificationService) ParseToken(ctx context.Context, token string) (int, error) {
	userId, err := s.getUserIdByAccessToken(ctx, token)
	if err != nil {
		userId, err = s.postgresUserRepo.GetUserIdByRefreshToken(token)
		if err != nil {
			return 0, err
		}
	}

	parsedToken, err := jwt.ParseWithClaims(token, &accesTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(s.tokensSigningKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := parsedToken.Claims.(*accesTokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	if claims.UserId != userId {
		return 0, errors.New("token does not match the stored user ID")
	}

	return claims.UserId, nil
}

func (s *AuthentificationService) createSession(ctx context.Context, userId int) (Tokens, error) {
	accessToken, err := s.newToken(userId, s.accessTokenTTL)
	if err != nil {
		return Tokens{}, err
	}

	refreshToken, err := s.newToken(userId, s.refreshTokenTTL)
	if err != nil {
		return Tokens{}, err
	}

	if err = s.setTokens(ctx, userId, accessToken, refreshToken); err != nil {
		return Tokens{}, err
	}

	return Tokens{
			AccesToken:   accessToken,
			RefreshToken: refreshToken},
		nil
}

func (s *AuthentificationService) newToken(userId int, tokenTTL time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &accesTokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		userId,
	})

	return token.SignedString([]byte(s.tokensSigningKey))
}

func (s *AuthentificationService) getUserIdByAccessToken(ctx context.Context, token string) (int, error) {
	userId, err := s.redisUserRepo.GetUserIdByAccessToken(ctx, token)
	if err != nil {
		userId, err = s.postgresUserRepo.GetUserIdByAccessToken(token)
		if err != nil {
			return 0, err
		}
	}

	return userId, nil
}

func (s *AuthentificationService) setTokens(ctx context.Context, userId int, accessToken, refreshToken string) error {
	if err := s.postgresUserRepo.SetTokens(userId, accessToken, refreshToken); err != nil {
		return err
	}

	if err := s.redisUserRepo.SetAccessToken(ctx, accessToken, userId, s.accessTokenTTL); err != nil {
		return err
	}

	return nil
}

func (s *AuthentificationService) deleteAllTokensByUserId(ctx context.Context, userId int) error {
	if err := s.postgresUserRepo.DeleteAllTokensByUserId(userId); err != nil {
		return err
	}

	if err := s.redisUserRepo.DeleteAllAccessTokensByUserId(ctx, userId); err != nil {
		return err
	}

	return nil
}

func (s *AuthentificationService) deleteUserIdByAccessToken(ctx context.Context, accessToken string) error {
	if err := s.postgresUserRepo.DeleteUserIdByAccessToken(accessToken); err != nil {
		return err
	}

	if err := s.redisUserRepo.DeleteUserIdByAccessToken(ctx, accessToken); err != nil {
		return err
	}

	return nil
}
