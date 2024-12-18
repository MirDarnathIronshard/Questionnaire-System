package handlers

import (
	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/internal/services"
	"github.com/QBG-P2/Voting-System/internal/utils"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"time"
)

type QuestionnaireAccessHandler struct {
	service services.QuestionnaireAccessService
}

func NewQuestionnaireAccessHandler(service services.QuestionnaireAccessService) *QuestionnaireAccessHandler {
	return &QuestionnaireAccessHandler{service: service}
}

func (h *QuestionnaireAccessHandler) AssignRole(c *fiber.Ctx) error {
	var req struct {
		models.BaseValidator
		UserID          uint       `json:"user_id" validate:"required"`
		QuestionnaireID uint       `json:"questionnaire_id" validate:"required"`
		RoleID          uint       `json:"role_id" validate:"required"`
		ExpiresAt       *time.Time `json:"expires_at"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}
	err := req.Validate(req)
	if err != nil {
		translateError, _ := utils.TranslateError(err, &req)

		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", translateError.ErrorList)
	}

	err = h.service.AssignRole(c.UserContext(), req.UserID, req.QuestionnaireID, req.RoleID, req.ExpiresAt)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Role assigned successfully")
}

func (h *QuestionnaireAccessHandler) UpdateRole(c *fiber.Ctx) error {
	var req struct {
		AccessID  uint `json:"access_id" validate:"required"`
		NewRoleID uint `json:"new_role_id" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	err := h.service.UpdateRole(c.UserContext(), req.AccessID, req.NewRoleID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Role updated successfully")
}

func (h *QuestionnaireAccessHandler) RevokeAccess(c *fiber.Ctx) error {
	accessID, err := strconv.ParseUint(c.Params("access_id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid access ID", nil)
	}

	err = h.service.RevokeAccess(c.UserContext(), uint(accessID))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Access revoked successfully")
}

func (h *QuestionnaireAccessHandler) GetQuestionnaireUsers(c *fiber.Ctx) error {
	questionnaireID, err := strconv.ParseUint(c.Params("questionnaire_id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid questionnaire ID", nil)
	}

	users, err := h.service.GetQuestionnaireUsers(c.UserContext(), uint(questionnaireID))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, users)
}
