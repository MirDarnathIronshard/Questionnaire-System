package handlers

import (
	"github.com/QBG-P2/Voting-System/internal/interfaces/http/request"
	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/internal/services"
	"github.com/QBG-P2/Voting-System/internal/utils"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type OptionHandler struct {
	Service *services.OptionService
}

func NewOptionHandler(service *services.OptionService) *OptionHandler {
	return &OptionHandler{Service: service}
}

func (h *OptionHandler) CreateOption(c *fiber.Ctx) error {
	var req request.CreateOptionRequest
	var option models.Option

	if err := c.BodyParser(&req); err != nil {
		errList := req.Validate()
		translatedErr, _ := utils.TranslateError(errList, &req)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", translatedErr.ErrorList)

	}

	// Validate the request
	if err := req.Validate(); err != nil {
		translatedErr, _ := utils.TranslateError(err, &req)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, translatedErr.Message, translatedErr.ErrorList)
	}

	// Map request to model
	if err := req.MapToModel(&option); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to process request", nil)
	}

	createdOption, err := h.Service.CreateOption(&option)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create option", err)
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, createdOption)
}

func (h *OptionHandler) GetOptionByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid option ID", nil)
	}

	option, err := h.Service.GetOptionByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Option not found", nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, option)
}

func (h *OptionHandler) UpdateOption(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid option ID", nil)
	}

	var req request.UpdateOptionRequest
	if err := c.BodyParser(&req); err != nil {
		errList := req.Validate()
		translatedErr, _ := utils.TranslateError(errList, &req)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", translatedErr.ErrorList)

	}

	// Validate the request
	if err := req.Validate(); err != nil {
		translatedErr, _ := utils.TranslateError(err, &req)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, translatedErr.Message, translatedErr.ErrorList)
	}

	var option models.Option
	// Map request to model
	if err := req.MapToModel(&option); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to process request", nil)
	}

	option.ID = uint(id)

	updatedOption, err := h.Service.UpdateOption(uint(id), &option)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update option", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, updatedOption)
}

func (h *OptionHandler) DeleteOption(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid option ID", nil)
	}

	err = h.Service.DeleteOption(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete option", nil)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *OptionHandler) GetOptionsByQuestionID(c *fiber.Ctx) error {
	var req request.GetOptionsRequest
	if err := c.QueryParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid query parameters", nil)
	}

	// Validate the request
	if err := req.Validate(); err != nil {
		translatedErr, _ := utils.TranslateError(err, &req)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, translatedErr.Message, translatedErr.ErrorList)
	}

	options, total, err := h.Service.OptionRepo.GetPaginatedByQuestionID(req.QuestionID, req.Pagination.Page, req.Pagination.PageSize)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve options", nil)
	}

	totalPages := int(total) / req.Pagination.PageSize
	if int(total)%req.Pagination.PageSize != 0 {
		totalPages++
	}

	return utils.PaginatedResponseWrapper(c, fiber.StatusOK, options, models.Pagination{
		Page:       req.Pagination.Page,
		PageSize:   req.Pagination.PageSize,
		Total:      int(total),
		TotalPages: totalPages,
	})
}
