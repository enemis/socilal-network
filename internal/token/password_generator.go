package token

import (
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"social-network-otus/internal/config"
)

type PasswordGenerator interface {
	CompareHashAndPassword(userId uuid.UUID, hashedPassword, password string) bool
	PasswordHash(userId uuid.UUID, password string) (string, error)
}

type PasswordGeneratorInstance struct {
	config *config.Config
}

func NewPasswordGenerator(config *config.Config) *PasswordGeneratorInstance {
	return &PasswordGeneratorInstance{config: config}
}

func (s *PasswordGeneratorInstance) CompareHashAndPassword(userId uuid.UUID, hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(s.buildPasswordString(userId, password))) == nil
}

func (s *PasswordGeneratorInstance) PasswordHash(userId uuid.UUID, password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(s.buildPasswordString(userId, password)), 14)
	return string(bytes), err
}

func (s *PasswordGeneratorInstance) buildPasswordString(userId uuid.UUID, password string) string {
	return strings.Join([]string{password, s.config.Salt, userId.String()}, "")
}
