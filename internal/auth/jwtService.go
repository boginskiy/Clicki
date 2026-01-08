package auth

import (
	"errors"
	"fmt"
	"time"

	conf "github.com/boginskiy/Clicki/cmd/config"
	"github.com/golang-jwt/jwt/v4"
)

// Claims - own statement.
type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

// JWTService - is JWT authentication.
type JWTService struct {
	Cfg conf.Config
}

func NewJWTService(config conf.Config) *JWTService {
	return &JWTService{Cfg: config}
}

// CreateToken - .
func (j *JWTService) CreateJWT(userID int) (string, error) {
	// New token.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(time.Duration(j.Cfg.GetTokenLiveTime() * int(time.Second))))},
		UserID: userID,
	})

	// Line of token.
	tokenStr, err := token.SignedString([]byte(j.Cfg.GetSecretKey()))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

// GetUserID - get id of client.
func (j *JWTService) GetIDAndValidJWT(checkingToken string) (int, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(checkingToken, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(j.Cfg.GetSecretKey()), nil
	})

	if err != nil {
		// Analize of expired token.
		var validErr *jwt.ValidationError

		if errors.As(err, &validErr) {
			// Bit И. Check flag expired token.
			if validErr.Errors&jwt.ValidationErrorExpired != 0 {
				return claims.UserID, ErrTokenIsExpired
			}
		}
		// Some errors.
		return 0, err
	}

	// Analize of invalid token.
	if !token.Valid {
		return claims.UserID, ErrTokenNotValid
	}
	return claims.UserID, nil
}
