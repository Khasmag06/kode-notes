package hasher

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type Hasher struct {
	salt string
}

func New(salt string) (*Hasher, error) {
	if salt == "" {
		return nil, errors.New("salt cannot be an empty string")
	}
	return &Hasher{salt: salt}, nil
}

func (h *Hasher) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(h.salt+password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(bytes), err
}

func (h *Hasher) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(h.salt+password))
	return err == nil
}
