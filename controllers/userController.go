package controllers

import (
	"monsterloveshop/middleware"
	"monsterloveshop/models"
	"monsterloveshop/request"
	"monsterloveshop/resources"
	"monsterloveshop/util"

	"github.com/gofiber/fiber/v2"
)

func (ctr *Controller) Login(c *fiber.Ctx) error {

	requestUser := new(request.LoginRequest)
	if err := c.BodyParser(requestUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ResponseFailureData("failure", "bad request"))
	}
	if validate := ctr.Validator.Struct(requestUser); validate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ResponseFailureData("failure", validate.Error()))
	}

	dbuser := new(models.User)
	result := ctr.DBConfig.DB.Where("Username = ?", requestUser.Username).First(dbuser)

	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ResponseFailureData("failure", "User not found."))
	}

	resultToken := middleware.Login(*requestUser, *dbuser, ctr.RedisStore)

	if resultToken.Status != 0 {
		return c.Status(fiber.StatusBadRequest).JSON(ResponseFailureData("failure", resultToken.Message))
	}

	responseData := ResponseSuccessLoginData("success", resultToken.Message, resources.SafeModelCustomer(dbuser))
	return c.JSON(responseData)
}

func (ctr *Controller) Register(c *fiber.Ctx) error {
	requestRegister := new(request.CreateUserRequest)

	if err := c.BodyParser(requestRegister); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ResponseFailureData("failure", "bad request"))
	}
	if validate := ctr.Validator.Struct(requestRegister); validate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ResponseFailureData("failure", validate.Error()))
	}

	user := models.User{
		Username: requestRegister.Username,
		Password: util.CreateHashPassword(requestRegister.Password),
		Name:     requestRegister.Name,
		UserType: "cust",
		Status:   1,
	}
	result := ctr.DBConfig.DB.Create(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ResponseFailureData("error", "Server error, please try again."))
	}

	return c.Status(fiber.StatusCreated).JSON(ResponseSuccessData("success", resources.ModelUser(&user)))
}

func (ctr *Controller) UserInfo(c *fiber.Ctx) error {

	claim := middleware.GetClaim(c)
	if claim == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ResponseFailureData("failure", "Unauthorzation"))
	}

	userInfo := new(models.User)
	result := ctr.DBConfig.DB.Where("id = ?", claim.ID).First(userInfo)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(ResponseFailureData("failure", result.Error.Error()))
	}

	return c.JSON(ResponseSuccessData("success", resources.ModelUser(userInfo)))
}

func (ctr *Controller) UpdateUser(c *fiber.Ctx) error {
	claim := middleware.GetClaim(c)
	if claim == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ResponseFailureData("failure", "Unauthorzation"))
	}

	requestUpd := new(request.UpdateUserRequest)
	if err := c.BodyParser(requestUpd); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ResponseFailureData("failure", "bad request"))
	}
	if validate := ctr.Validator.Struct(requestUpd); validate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ResponseFailureData("failure", validate.Error()))
	}
	userid := c.Params("userid")

	// updData := map[string]interface{}{
	// 	"name":   requestUpd.Name,
	// 	"status": requestUpd.Status,
	// }
	// valueOfReq := reflect.ValueOf(requestUpd).Elem()
	// typeOfReq := reflect.TypeOf(requestUpd).Elem()

	// for i := 0; i < valueOfReq.NumField(); i++ {
	// 	v := valueOfReq.Field(i)
	// 	k := typeOfReq.Field(i).Name
	// 	t := v.Kind()
	// 	fmt.Println(v, k, t, requestUpd)
	// }

	updData := util.CheckKeyIsHave(requestUpd)
	updResult := ctr.DBConfig.DB.Model(&models.User{}).Where("id = ?", userid).Updates(updData)

	if updResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ResponseFailureData("error", "Server error, please try again."))
	}

	user := new(models.User)
	result := ctr.DBConfig.DB.Where("id = ?", userid).First(user)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ResponseFailureData("error", result.Error.Error()))
	}

	return c.JSON(ResponseSuccessData("success", resources.ModelUser(user)))

}
