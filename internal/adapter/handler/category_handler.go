package handler

import (
	"trustnews/internal/adapter/handler/request"
	"trustnews/internal/adapter/handler/response"
	"trustnews/internal/core/domain/entity"
	"trustnews/internal/core/service"
	"trustnews/lib/conv"
	validatorLib "trustnews/lib/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

var defaultSuccessReponse response.DefaultSuccessReponse

type CategoryHandler interface {
	GetCategories(ctx *fiber.Ctx) error
	GetCategoryByID(ctx *fiber.Ctx) error
	CreateCategory(ctx *fiber.Ctx) error
	EditCategoryByID(ctx *fiber.Ctx) error
	DeleteCategory(ctx *fiber.Ctx) error

	GetCategoryFE(c *fiber.Ctx) error
}

type categoryHandler struct {
	categoryService service.CategoryService
}

// GetCategoryFE implements CategoryHandler.
func (ch *categoryHandler) GetCategoryFE(c *fiber.Ctx) error {
	results, err := ch.categoryService.GetCategories(c.Context())
	if err != nil {
		code = "[HANDLER] GetCategoryFE - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	categoryResponses := []response.SuccessCategoryResponse{}
	for _, result := range results {
		categoryResponse := response.SuccessCategoryResponse{
			ID:            result.ID,
			Title:         result.Title,
			Slug:          result.Slug,
			CreatedByName: result.User.Name,
		}

		categoryResponses = append(categoryResponses, categoryResponse)
	}

	defaultSuccessReponse.Meta.Status = true
	defaultSuccessReponse.Meta.Message = "Categories Fetched Successfully"
	defaultSuccessReponse.Pagination = nil
	defaultSuccessReponse.Data = categoryResponses

	return c.JSON(defaultSuccessReponse)
}

func (ch *categoryHandler) CreateCategory(c *fiber.Ctx) error {
	var req request.CategoryRequest
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[HANDLER] CreateCategory - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized Access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	if err = c.BodyParser(&req); err != nil {
		code = "[HANDLER] CreateCategory - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid Request Body"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	if err = validatorLib.ValidateStruct(req); err != nil {
		code = "[HANDLER] CreateCategory - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	reqEntity := entity.CategoryEntity{
		Title: req.Title,
		User: entity.UserEntity{
			ID: int64(userID),
		},
	}

	err = ch.categoryService.CreateCategory(c.Context(), reqEntity)
	if err != nil {
		code = "[HANDLER] CreateCategory - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	defaultSuccessReponse.Data = nil
	defaultSuccessReponse.Pagination = nil
	defaultSuccessReponse.Meta.Status = true
	defaultSuccessReponse.Meta.Message = "Category Created Successfully"
	return c.Status(fiber.StatusCreated).JSON(defaultSuccessReponse)
}

func (ch *categoryHandler) DeleteCategory(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[HANDLER] DeleteCategory - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized Access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	idParam := c.Params("categoryID")
	id, err := conv.StringToInt64(idParam)
	if err != nil {
		code = "[HANDLER] DeleteCategory - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	err = ch.categoryService.DeleteCategory(c.Context(), id)
	if err != nil {
		code = "[HANDLER] DeleteCategory - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessReponse.Data = nil
	defaultSuccessReponse.Pagination = nil
	defaultSuccessReponse.Meta.Status = true
	defaultSuccessReponse.Meta.Message = "Category Deleted Successfully"

	return c.JSON(defaultSuccessReponse)
}

func (ch *categoryHandler) EditCategoryByID(c *fiber.Ctx) error {
	var req request.CategoryRequest
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[HANDLER] EditCategoryByID - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized Access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	if err = c.BodyParser(&req); err != nil {
		code = "[HANDLER] EditCategoryByID - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid Request Body"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	if err = validatorLib.ValidateStruct(req); err != nil {
		code = "[HANDLER] EditCategoryByID - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	idParam := c.Params("categoryID")
	id, err := conv.StringToInt64(idParam)
	if err != nil {
		code = "[HANDLER] EditCategoryByID - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	reqEntity := entity.CategoryEntity{
		ID:    id,
		Title: req.Title,
		User: entity.UserEntity{
			ID: int64(userID),
		},
	}

	err = ch.categoryService.EditCategoryByID(c.Context(), reqEntity)
	if err != nil {
		code = "[HANDLER] EditCategoryByID - 5"
		log.Errorw(code, err)
		errorResp.Meta.Status = true
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessReponse.Data = nil
	defaultSuccessReponse.Pagination = nil
	defaultSuccessReponse.Meta.Status = true
	defaultSuccessReponse.Meta.Message = "Category Updated Successfully"
	return c.JSON(defaultSuccessReponse)
}

func (ch *categoryHandler) GetCategories(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[HANDLER] GetCategories - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized Access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	results, err := ch.categoryService.GetCategories(c.Context())
	if err != nil {
		code = "[HANDLER] GetCategories - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid Request Body"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	categoryResponses := []response.SuccessCategoryResponse{}
	for _, result := range results {
		categoryResponse := response.SuccessCategoryResponse{
			ID:            result.ID,
			Title:         result.Title,
			Slug:          result.Slug,
			CreatedByName: result.User.Name,
		}

		categoryResponses = append(categoryResponses, categoryResponse)
	}

	defaultSuccessReponse.Meta.Status = true
	defaultSuccessReponse.Meta.Message = "Categories Fetched Successfully"
	defaultSuccessReponse.Pagination = nil
	defaultSuccessReponse.Data = categoryResponses

	return c.JSON(defaultSuccessReponse)
}

func (ch *categoryHandler) GetCategoryByID(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[HANDLER] GetCategoryByID - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized Access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	idParam := c.Params("categoryID")
	id, err := conv.StringToInt64(idParam)
	if err != nil {
		code = "[HANDLER] GetCategoryByID - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	result, err := ch.categoryService.GetCategoryByID(c.Context(), id)
	if err != nil {
		code = "[HANDLER] GetCategoryByID - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	categoryResponse := response.SuccessCategoryResponse{
		ID:            id,
		Title:         result.Title,
		Slug:          result.Slug,
		CreatedByName: result.User.Name,
	}

	defaultSuccessReponse.Meta.Status = true
	defaultSuccessReponse.Meta.Message = "Categories Fetched Detail Successfully"
	defaultSuccessReponse.Pagination = nil
	defaultSuccessReponse.Data = categoryResponse

	return c.JSON(defaultSuccessReponse)
}

func NewCategoryHandler(categoryService service.CategoryService) CategoryHandler {
	return &categoryHandler{categoryService: categoryService}
}
