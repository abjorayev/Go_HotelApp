package jwt

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJwtToken(userID int, role string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  strconv.Itoa(userID),
		"role": role,
		"iss":  "Hotel_App",
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(72 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
