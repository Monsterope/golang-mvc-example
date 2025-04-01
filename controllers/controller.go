package controllers

import (
	"monsterloveshop/databases"
	"monsterloveshop/resources"
	"monsterloveshop/store"

	"github.com/go-playground/validator/v10"
)

type Controller struct {
	DBConfig   *databases.DatabaseConfig
	RedisStore *store.RedisAuthStore
	Validator  *validator.Validate
}

func NewController(dbconfig *databases.DatabaseConfig, redisstore *store.RedisAuthStore) *Controller {
	return &Controller{
		DBConfig:   dbconfig,
		RedisStore: redisstore,
		Validator:  validator.New(),
	}
}

type (
	ResponseSuccessLogin struct {
		Status      string                 `json:"status"`
		AccessToken string                 `json:"access_token"`
		User        resources.SafeCustomer `json:"user"`
	}
	ResponseSuccess struct {
		Status string      `json:"status"`
		Data   interface{} `json:"data"`
	}
	ResponseFailure struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}
)

func ResponseSuccessLoginData(status string, accessToken string, user resources.SafeCustomer) ResponseSuccessLogin {
	return ResponseSuccessLogin{
		Status:      status,
		AccessToken: accessToken,
		User:        user,
	}
}

func ResponseSuccessData(status string, data interface{}) ResponseSuccess {
	return ResponseSuccess{
		Status: status,
		Data:   data,
	}
}

func ResponseFailureData(status string, message string) ResponseFailure {
	return ResponseFailure{
		Status:  status,
		Message: message,
	}
}
