package tokener

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
  UnsupportedMethod = errors.New("unsupported token method")
  TokenNotFound = errors.New("token not found in header")
  InvalidToken = errors.New("invalid token")
  FailedToConvertClaims = errors.New("failed to convert claims to mapclaims")
  FailedToConvertUsername = errors.New("failed to convert username to string")
  FailedToParseToken = errors.New("failed to parse jwt token")
)

type Tokener struct {
  secret []byte
}

func Create(secret []byte) *Tokener {
  return &Tokener{
    secret: secret, }
}

func (t Tokener) GenToken(username string) (string, error) {
  claims := jwt.MapClaims{
    "username": username,
    "exp": time.Now().Add(2 * time.Hour).Unix(),
  }

  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

  slog.Debug("created token")
  return token.SignedString(t.secret)
}

func (t Tokener) parseToken(tokenString string) (*jwt.Token, error) {
  token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
      return nil, UnsupportedMethod
    }

    return t.secret, nil
  })

  slog.Debug("parsed token")
  return token, err
}

func (t Tokener) ValidateToken(r *http.Request) (string, error) {
  header := r.Header.Get("Authorization")
  if header == "" {
    return "", TokenNotFound
  }

  tokenString := strings.Split(header, " ")[1]

  token, err := t.parseToken(tokenString)
  if err != nil {
    return "", FailedToParseToken
  }
  if !token.Valid {
    return "", InvalidToken
  }
  
  claims, ok := token.Claims.(jwt.MapClaims)
  if !ok {
    return "", FailedToConvertClaims
  }

  username, ok := claims["username"].(string)
  if !ok {
    return "", FailedToConvertUsername
  }

  slog.Debug("got username from token")
  return username, nil
}
