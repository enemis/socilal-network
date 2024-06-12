//package service
//
//import (
//	"crypto/sha1"
//	"fmt"
//
//	"my_treassure/internal/entity"
//	"my_treassure/internal/repository"
//
//	"github.com/sirupsen/logrus"
//)
//
//const (
//	salt = "test"
//)
//
//type AuthService struct {
//	entityManager *repository.EntityManager
//}
//
//func NewAuthService(entityManager *repository.EntityManager) *AuthService {
//	return &AuthService{
//		entityManager: entityManager,
//	}
//}
//
//func (s *AuthService) CreateUser(user *entity.User) (uint, error) {
//	return s.entityManager.CreateUser(user)
//}
//
//func (s *AuthService) GenerateToken(username, password string) (string, error) {
//	user, err := s.entityManager.GetUser(username, s.generatePasswordHash(password))
//	logrus.Errorln("user")
//	logrus.Errorln(user)
//	if err != nil {
//		return "", err
//	}
//
//	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
//	// 	jwt.StandardClaims{
//	// 		ExpiresAt: time.Now().Add(tokenTTL).Unix(),
//	// 		IssuedAt:  time.Now().Unix(),
//	// 	},
//	// 	user.Id,
//	// })
//
//	// return token.SignedString([]byte(signingKey))
//	return "test", nil
//}
//
//// func (s *AuthService) ParseToken(accessToken string) (int, error) {
//// 	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
//// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//// 			return nil, errors.New("invalid signing method")
//// 		}
//
//// 		return []byte(signingKey), nil
//// 	})
//
//// 	if err != nil {
//// 		return 0, err
//// 	}
//
//// 	claims, ok := token.Claims.(*tokenClaims)
//// 	if !ok {
//// 		return 0, errors.New("token claims are not of type tokenClaims")
//// 	}
//
//// 	return claims.UserId, nil
//// }
//
//func (s *AuthService) generatePasswordHash(password string) string {
//	hash := sha1.New()
//	hash.Write([]byte(password))
//
//	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
//}
