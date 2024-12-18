package handlers

import (
	"net/http"
	"strconv"

	"github.com/QBG-P2/Voting-System/internal/interfaces/http/request"
	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/internal/services"
	"github.com/QBG-P2/Voting-System/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type ResponseHandler struct {
	Service              services.ResponseService
	QuestionnaireService services.QuestionnaireService
}

func NewResponseHandler(service services.ResponseService) *ResponseHandler {
	return &ResponseHandler{Service: service}
}

func (h *ResponseHandler) CreateResponse(c *fiber.Ctx) error {
	var req request.CreateResponseRequest
	var response models.Response

	if err := c.BodyParser(&req); err != nil {
		errList := req.Validate()
		translatedErr, _ := utils.TranslateError(errList, &req)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", translatedErr.ErrorList)
	}

	if err := req.Validate(); err != nil {
		translatedErr, _ := utils.TranslateError(err, &req)

		return utils.ErrorResponse(c, fiber.StatusBadRequest, translatedErr.Message, translatedErr.ErrorList)
	}

	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid user ID", nil)
	}

	if err := req.MapToModel(&response); err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Map data error", nil)
	}

	response.UserID = userID

	err := h.QuestionnaireService.SubmitResponse(c.UserContext(), &response)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create response", err)
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, response)
}

func (h *ResponseHandler) GetResponsesByQuestionnaireID(c *fiber.Ctx) error {
	var req request.GetResponsesRequest

	if err := c.QueryParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid query parameters", nil)
	}

	if err := req.Validate(); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validation failed", err)
	}

	page := req.Page
	size := req.PageSize

	responses, err := h.Service.GetResponsesByQuestionnaireID(c.UserContext(), req.QuestionnaireID, page, size)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve responses", nil)
	}

	total, err := h.Service.GetTotalResponsesByQuestionnaireID(c.UserContext(), req.QuestionnaireID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve total responses", nil)
	}

	totalPages := (total + size - 1) / size

	response := fiber.Map{
		"data": responses,
		"pagination": fiber.Map{
			"total":      total,
			"page":       page,
			"size":       size,
			"totalPages": totalPages,
		},
	}

	return utils.SuccessResponse(c, fiber.StatusOK, response)
}

func (h *ResponseHandler) GetResponseByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	responseID, err := strconv.Atoi(idParam)
	if err != nil || responseID <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid response ID", nil)
	}

	response, err := h.Service.GetResponseByID(c.UserContext(), uint(responseID))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Response not found", nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, response)
}

func (h *ResponseHandler) UpdateResponse(c *fiber.Ctx) error {
	idParam := c.Params("id")
	responseID, err := strconv.Atoi(idParam)
	if err != nil || responseID <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid response ID", nil)
	}

	var req request.UpdateResponseRequest
	var response models.Response

	if err := c.BodyParser(&req); err != nil {
		errList := req.Validate()
		translatedErr, _ := utils.TranslateError(errList, &req)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", translatedErr.ErrorList)

	}

	if err := req.Validate(); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validation failed", err)
	}

	response.ID = uint(responseID)
	response.Answer = req.Content

	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid user ID", nil)
	}

	isOwner, err := h.Service.IsOwnerResponse(c.UserContext(), response.ID, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error verifying response ownership", nil)
	}
	if !isOwner {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "You do not own this response", nil)
	}

	if err := h.Service.UpdateResponse(c.UserContext(), &response); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update response", nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, response)
}

func (h *ResponseHandler) DeleteResponse(c *fiber.Ctx) error {
	idParam := c.Params("id")
	responseID, err := strconv.Atoi(idParam)
	if err != nil || responseID <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid response ID", nil)
	}

	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid user ID", nil)
	}

	isOwner, err := h.Service.IsOwnerResponse(c.UserContext(), uint(responseID), userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error verifying response ownership", nil)
	}
	if !isOwner {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "You do not own this response", nil)
	}

	if err := h.Service.DeleteResponse(c.UserContext(), uint(responseID)); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete response", nil)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
