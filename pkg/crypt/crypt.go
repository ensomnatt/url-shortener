package crypt

import (
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

type Crypt struct {}

func Create() *Crypt {
  return &Crypt{}
}

func (c Crypt) Password(password []byte) (string, error) {
  hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
  if err != nil {
    return "", err
  }

  slog.Debug("hashed password", "password", string(password), "hashed password", string(hashedPassword))
  return string(hashedPassword), nil
}

func (c Crypt) GetPassword(hash, userPassword []byte) bool {
  err := bcrypt.CompareHashAndPassword(hash, userPassword)
  
  slog.Debug("compared password with hash")
  return err == nil 
}
