package middleware

import (
	"monsterloveshop/config"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func GetTokenJWT(headerAuth string) (*Claim, string) {
	if headerAuth == "" || !strings.HasPrefix(headerAuth, "Bearer ") {
		return nil, "Invalid or missing authorization header"
	}

	jwtSecret := []byte(config.GetEnv("jwt.secret"))
	headerToken := strings.TrimPrefix(headerAuth, "Bearer ")
	parts := strings.Split(headerToken, ".")
	if len(parts) != 3 {
		return nil, "Invalid token"
		// return nil, c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		// 	"error": "Invalid Token",
		// })
	}

	claimUser := &Claim{}

	token, err := jwt.ParseWithClaims(headerToken, claimUser, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if token == nil {
		return nil, "Invalid token"
	} else if err != nil {
		return nil, err.Error()
	}
	// fmt.Println(token)

	return claimUser, ""

}
