package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey []byte

func init() {
	secret := os.Getenv("SENTINELOS_JWT_SECRET")
	if secret == "" {
		// fallback temporal
		secret = "CAMBIAR_ESTO_LUEGO"
	}
	jwtKey = []byte(secret)
}

// jwtkey expone la clave de forma controlada
func JwtKey() []byte {
	return jwtKey
}

type Claims struct {
	Username string `json:"sub"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT(user *User) (string, int64, error) {

	expiration := time.Now().Add(1 * time.Hour)

	claims := &Claims{
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", 0, err
	}

	return tokenString, int64(time.Until(expiration).Seconds()), nil
}
