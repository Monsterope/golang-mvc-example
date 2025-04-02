package middleware

import (
	"monsterloveshop/config"
	"monsterloveshop/models"
	"monsterloveshop/request"
	"monsterloveshop/store"
	"monsterloveshop/util"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

type JWT struct {
	Secret string
}

type Claim struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	UserType  string `json:"role"`
	TokenType string `json:"token_type"`
	jwt.StandardClaims
}

type RedisAuthMiddleware struct {
	AuthRedisStore *store.RedisAuthStore
}

func NewMiddlewareAuthRedis(redisStore *store.RedisAuthStore) *RedisAuthMiddleware {
	return &RedisAuthMiddleware{
		AuthRedisStore: redisStore,
	}
}

type ReturnAction struct {
	Status   int
	Message  string
	Message2 string
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (j *JWT) CreateToken(user models.User, typetoken string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claim{
		ID:        user.Id,
		Username:  user.Username,
		UserType:  user.UserType,
		TokenType: typetoken,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.Secret))
}

func CreateTokenRefresh(claim *Claim, typetoken string, jwtSecret string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claim{
		ID:        claim.ID,
		Username:  claim.Username,
		UserType:  claim.UserType,
		TokenType: typetoken,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func Login(reqbody request.LoginRequest, dbuser models.User, redisstore *store.RedisAuthStore) ReturnAction {

	if err := util.CompareHasPassword(dbuser.Password, reqbody.Password); err != nil {
		return ReturnAction{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid Password.",
		}
	}

	jwtsecret := config.GetEnv("jwt.secret")

	jwt := &JWT{Secret: jwtsecret}
	tokenstr, err := jwt.CreateToken(dbuser, "access")
	if err != nil {
		return ReturnAction{
			Status:  fiber.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}
	tokenrefstr, err := jwt.CreateToken(dbuser, "refresh")
	if err != nil {
		return ReturnAction{
			Status:  fiber.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}

	err = redisstore.Set("token:"+strconv.FormatUint(uint64(dbuser.Id), 10), tokenstr)
	if err != nil {
		return ReturnAction{
			Status:  fiber.StatusInternalServerError,
			Message: "Error while saving token to Redis",
		}
	}

	return ReturnAction{
		Status:   0,
		Message:  tokenstr,
		Message2: tokenrefstr,
	}
}

func RefreshToken(tokenRef string, redisstore *store.RedisAuthStore) ReturnAction {

	jwtSecret := config.GetEnv("jwt.secret")

	claimUser := &Claim{}

	token, err := jwt.ParseWithClaims(tokenRef, claimUser, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return ReturnAction{
			Status:  fiber.StatusUnauthorized,
			Message: err.Error(),
		}
	}

	if claimUser.TokenType != "refresh" {
		return ReturnAction{
			Status:  fiber.StatusUnauthorized,
			Message: "Unauthorized",
		}
	}

	newAccessToken, err := CreateTokenRefresh(claimUser, "access", jwtSecret)

	if err != nil {
		return ReturnAction{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	// set token to redis || use redis because want user 1 token
	err = redisstore.Set("token:"+strconv.FormatUint(uint64(claimUser.ID), 10), newAccessToken)
	if err != nil {
		return ReturnAction{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return ReturnAction{
		Status:  0,
		Message: newAccessToken,
	}

}

func (a *RedisAuthMiddleware) AuthIsCustomer(c *fiber.Ctx) error {
	myClaim, err, tokenCur := GetTokenJWT(c.Get("Authorization"))

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "failure",
			"message": err,
		})
	}

	// get token to redis || use redis because want user 1 token
	redisToken, err := a.AuthRedisStore.Get("token:" + strconv.FormatUint(uint64(myClaim.ID), 10))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "failure",
			"message": "Unauthorizated.",
		})
	}

	if redisToken != tokenCur {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "failure",
			"message": "Unauthorizated..",
		})
	}

	if myClaim.TokenType != "access" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "failure",
			"message": "Unauthorizated",
		})
	}

	if myClaim.UserType != "cust" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "failure",
			"message": "Unauthorization.",
		})
	}

	c.Locals("claim", myClaim)
	return c.Next()
}

func (a *RedisAuthMiddleware) AuthIsAdmin(c *fiber.Ctx) error {
	myClaim, err, tokenCur := GetTokenJWT(c.Get("Authorization"))

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "failure",
			"message": err,
		})
	}

	// get token to redis || use redis because want user 1 token
	redisToken, err := a.AuthRedisStore.Get("token:" + strconv.FormatUint(uint64(myClaim.ID), 10))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "failure",
			"message": "Unauthorizated.",
		})
	}

	if redisToken != tokenCur {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "failure",
			"message": "Unauthorizated..",
		})
	}

	if myClaim.TokenType != "access" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "failure",
			"message": "Unauthorizated",
		})
	}

	if myClaim.UserType != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "failure",
			"message": "Unauthorization.",
		})
	}

	c.Locals("claim", myClaim)
	return c.Next()
}
