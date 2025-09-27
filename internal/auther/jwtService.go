package auther

import (
	"errors"
	"fmt"
	"time"

	conf "github.com/boginskiy/Clicki/cmd/config"
	repo "github.com/boginskiy/Clicki/internal/repository"
	"github.com/golang-jwt/jwt/v4"
)

// Claims - собственное утверждение
type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

type JWTService struct{}

func NewJWTService(kwargs conf.VarGetter, repoUser repo.Repository) *JWTService {
	return &JWTService{}
}

// CreateToken - создание токена
func (j *JWTService) CreateJWT(userID int) (string, error) {
	// Новый токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// Время жизни токена
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKENEXP)),
		},
		UserID: userID,
	})

	// Строка токена
	tokenStr, err := token.SignedString([]byte(SECRETKEY))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

// GetUserID - получение идентификатора клиента
func (j *JWTService) GetIDAndValidJWT(tokenStr string) (int, error) {
	// экземпляр структуры с утверждениями
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(SECRETKEY), nil
	})

	if err != nil {
		// Анализ просроченного токена
		var validErr *jwt.ValidationError
		if errors.As(err, &validErr) {
			// Битовый И. Проверка флага просроченного токена
			if validErr.Errors&jwt.ValidationErrorExpired != 0 {
				return claims.UserID, ErrTokenIsExpired
			}
		}
		// Другие ошибки
		return 0, err
	}

	// Анализ невалидного токена
	if !token.Valid {
		return claims.UserID, ErrTokenNotValid
	}
	return claims.UserID, nil
}

// Интересно, а если токен просрочился, парсер запарсит в claims.UserID ?
