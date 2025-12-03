package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/PhilAldridge/spell-api/ent"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type contextKey string

type Claims struct {
	jwt.MapClaims
	AccountType string
}

const userContextKey = contextKey("users")

func NewContextWithUser(ctx context.Context, u *ent.User) context.Context {
	return context.WithValue(ctx, userContextKey, u)
} 

func UserFromContext(ctx context.Context) (*ent.User, bool) {
	u, ok := ctx.Value(userContextKey).(*ent.User)

	return u,ok
}

func HashPassword(plain string) (string, error) {
	hash, err:= bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func ComparePassword(hash, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
}

func jwtSecret() ([]byte, error) {
	s:= os.Getenv("JWT_SECRET")
	if s=="" {
		return nil, errors.New("JWT_SECRET NOT SET")
	}

	return []byte(s), nil
}

func GenerateAccessToken(userID int, expireMinutes int, accountType string) (string, error) {
	secret, err:= jwtSecret()
	if err != nil {
		return "", err
	}

	now:= time.Now().UTC()

	claims := jwt.MapClaims{
		"sub": strconv.Itoa(userID),
		"iat": now.Unix(),
		"exp": now.Add(time.Minute*time.Duration(expireMinutes)).Unix(),
		"typ": "access",
		"accountType": accountType,
	}

	token:= jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func ParseAccessToken(tokenString string) (int, jwt.MapClaims, error) {
	secret, err:= jwtSecret()
	if err != nil {
		return 0, nil, err
	}

	token, err:= jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _,ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}

		return secret, nil
	}, jwt.WithLeeway(5*time.Second))

	if err != nil {
		return 0, nil, err
	}

	if !token.Valid {
		return 0, nil, jwt.ErrTokenInvalidClaims
	}

	claims, ok:= token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, nil, jwt.ErrTokenInvalidClaims
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return 0, nil, jwt.ErrTokenInvalidClaims
	}

	id, err := strconv.Atoi(sub)
	if err != nil {
		return 0, nil, jwt.ErrTokenInvalidId
	}

	return id, claims,nil
}

func HashRefreshToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}