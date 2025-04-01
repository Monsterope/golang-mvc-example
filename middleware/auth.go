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
	ID       int64  `json:"id"`
	Username string `json:"username"`
	UserType string `json:"role"`
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
	Status  int
	Message string
}

func (j *JWT) CreateToken(user models.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claim{
		ID:       user.Id,
		Username: user.Username,
		UserType: user.UserType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.Secret))
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
	tokenstr, err := jwt.CreateToken(dbuser)

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
		Status:  0,
		Message: tokenstr,
	}
}

func (a *RedisAuthMiddleware) AuthIsCustomer(c *fiber.Ctx) error {
	myClaim, err := GetTokenJWT(c.Get("Authorization"))

	if err != "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "failure",
			"message": err,
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
	myClaim, err := GetTokenJWT(c.Get("Authorization"))

	if err != "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "failure",
			"message": err,
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
