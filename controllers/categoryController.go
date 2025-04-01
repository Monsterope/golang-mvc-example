package controllers

import (
	"monsterloveshop/models"
	"monsterloveshop/request"
	"monsterloveshop/resources"
	"monsterloveshop/util"

	"github.com/gofiber/fiber/v2"
)

func (ctr *Controller) CreateCategory(c *fiber.Ctx) error {

	requestCre := new(request.CategoryCreateRequest)
	if err := c.BodyParser(requestCre); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ResponseFailureData("failure", "bad request"))
	}
	if validate := ctr.Validator.Struct(requestCre); validate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ResponseFailureData("failure", validate.Error()))
	}

	result := ctr.DBConfig.DB.Create(requestCre.Item)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ResponseFailureData("error", "Server error, please try again."))
	}

	return c.Status(fiber.StatusCreated).JSON(ResponseSuccessData("success", "Created success."))
}

func (ctr *Controller) GetCategoryAll(c *fiber.Ctx) error {
	categories := new([]models.Category)

	result := ctr.DBConfig.DB.Find(categories)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ResponseFailureData("error", "Server error."))
	}

	var categoryResource []resources.CategoryResource
	for _, category := range *categories {
		categoryResource = append(categoryResource, resources.GetCategoryResource(&category))
	}

	return c.JSON(ResponseSuccessData("success", categoryResource))
}

func (ctr *Controller) UpdateCategory(c *fiber.Ctx) error {
	requestUpd := new(request.CategoryUpdateRequest)
	if err := c.BodyParser(requestUpd); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ResponseFailureData("failure", "bad request."))
	}
	if validate := ctr.Validator.Struct(requestUpd); validate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ResponseFailureData("failure", "bad request."))
	}

	cateid := c.Params("cateid")
	updData := util.CheckKeyIsHave(requestUpd)
	updResult := ctr.DBConfig.DB.Model(&models.Category{}).Where("id = ?", cateid).Updates(updData)
	if updResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ResponseFailureData("error", "Server error."))
	}

	category := new(models.Category)
	result := ctr.DBConfig.DB.Where("id = ?", cateid).First(category)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ResponseFailureData("error", result.Error.Error()))
	}

	return c.JSON(ResponseSuccessData("success", category))

}
