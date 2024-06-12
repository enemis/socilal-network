package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/fx"

	"github.com/golang-jwt/jwt"

	"social-network-otus/internal/app_error"
	"social-network-otus/internal/config"
	"social-network-otus/internal/session"
	"social-network-otus/internal/token"
	"social-network-otus/internal/user"
	"social-network-otus/internal/utils"
)

type ServiceParams struct {
	fx.In
	Config         *config.Config
	Generator      token.PasswordGenerator
	SessionStorage *session.SessionStorage
	UserService    *user.Service
}

type AuthService struct {
	userService       *user.Service
	signKey           string
	config            *config.Config
	passwordGenerator token.PasswordGenerator
	sessionStorage    *session.SessionStorage
}

func NewAuthService(params ServiceParams) *AuthService {
	sign := strings.Join([]string{params.Config.SigningKey, params.Config.Salt}, "")
	return &AuthService{
		userService:       params.UserService,
		config:            params.Config,
		signKey:           sign,
		passwordGenerator: params.Generator,
		sessionStorage:    params.SessionStorage,
	}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*string, *app_error.AppError) {
	const errorMessage = "email or password invalid"
	user, appErr := s.userService.GetUserByEmail(email)

	if appErr != nil {
		return nil, app_error.New(appErr.OriginalError(), errorMessage, http.StatusBadRequest)
	}

	if !s.passwordGenerator.CompareHashAndPassword(user.Id, user.Password, password) {
		return nil, app_error.New(appErr.OriginalError(), errorMessage, http.StatusBadRequest)
	}

	expireAt := time.Now().Add(time.Duration(s.config.TokenTTL) * time.Second)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		&Claims{
			StandardClaims: jwt.StandardClaims{
				IssuedAt:  time.Now().Unix(),
				ExpiresAt: expireAt.Unix(),
				Issuer:    user.Id.String(),
			},
		},
	)

	tokenString, err := s.signToken(token)

	if appErr != nil {
		return nil, app_error.NewInternalServerError(err)
	}

	return utils.Ptr(tokenString), nil
}

func (s *AuthService) ParseToken(accessToken string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(accessToken, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(s.signKey), nil
	})

	if err = token.Claims.Valid(); err != nil {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return uuid.Nil, errors.New("token claims are not of type tokenClaims")
	}

	return uuid.FromBytes([]byte(claims.Issuer))
}

func (s *AuthService) signToken(token *jwt.Token) (string, error) {
	return token.SignedString([]byte(s.signKey))
}

//
//func (s *AuthService) AuthByToken(ctx context.Context, authTokenString string) error {
//	var tokenClaims Claims
//	token, err := jwt.ParseWithClaims(authTokenString, &tokenClaims, func(token *jwt.Token) (interface{}, error) {
//		return []byte(s.signKey), nil
//	})
//
//	if err != nil {
//		return err
//	}
//
//	if !token.Valid {
//		return s.tokenError(tokenClaims.TokenId)
//	}
//
//	authToken, err := s.repository.GetTokenById(ctx, tokenClaims.TokenId)
//	if err != nil {
//		return err
//	}
//
//	ruleSet := AdminTokenRuleSet{AdminToken: authToken}
//	validatorInstance := validator.NewValidator(ruleSet, UserTokenValidator{claim: tokenClaims})
//	validationErrors := validatorInstance.Validate(ruleSet)
//
//	if validationErrors != nil {
//		return s.tokenError(tokenClaims.TokenId)
//	}
//
//	user, err := s.sessionStorage.GetSessionUserById(ctx, tokenClaims.UserID)
//
//	if err != nil {
//		return err
//	}
//
//	s.sessionStorage.SetAuthenticatedUser(user)
//
//	return nil
//}

//func (s *AuthService) SignUp(ctx context.Context, email, password string) (string, error) {
//
//	user, err := s.repository.CreateUser(ctx, email, password)
//	if err != nil {
//		return "", err
//	}
//
//	return user.Password, nil
//}
