package handlers

import (
	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/internal/services"
	"github.com/QBG-P2/Voting-System/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type QuestionnairePermissionHandler struct {
	Service services.QuestionnairePermissionService
}

func NewQuestionnairePermissionHandler(service services.QuestionnairePermissionService) *QuestionnairePermissionHandler {
	return &QuestionnairePermissionHandler{Service: service}
}

func (h *QuestionnairePermissionHandler) Create(c *fiber.Ctx) error {
	var permission models.QuestionnairePermission
	if err := c.BodyParser(&permission); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	if err := permission.Validate(); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validation failed", err)
	}

	if err := h.Service.Create(&permission); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create permission", nil)
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, permission)
}

func (h *QuestionnairePermissionHandler) GetByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid permission ID", nil)
	}

	permission, err := h.Service.GetByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Permission not found", nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, permission)
}

func (h *QuestionnairePermissionHandler) GetByQuestionnaireID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid questionnaire ID", nil)
	}

	permissions, err := h.Service.GetByQuestionnaireID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve permissions", nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, permissions)
}

func (h *QuestionnairePermissionHandler) GetByUserID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID", nil)
	}

	permissions, err := h.Service.GetByUserID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve permissions", nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, permissions)
}

func (h *QuestionnairePermissionHandler) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid permission ID", nil)
	}

	var permission models.QuestionnairePermission
	if err := c.BodyParser(&permission); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	permission.ID = uint(id)

	if err := h.Service.Update(&permission); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update permission", nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, permission)
}

func (h *QuestionnairePermissionHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid permission ID", nil)
	}

	if err := h.Service.Delete(uint(id)); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete permission", nil)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
