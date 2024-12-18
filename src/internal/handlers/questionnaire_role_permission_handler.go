package handlers

import (
	"github.com/QBG-P2/Voting-System/internal/services"
	"github.com/QBG-P2/Voting-System/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strconv"
)

type QuestionnaireRolePermissionHandler struct {
	rolePermissionService services.QuestionnaireRolePermissionService
}

func NewQuestionnaireRolePermissionHandler(service services.QuestionnaireRolePermissionService) *QuestionnaireRolePermissionHandler {
	return &QuestionnaireRolePermissionHandler{
		rolePermissionService: service,
	}
}

func (h *QuestionnaireRolePermissionHandler) AssignPermissions(c *fiber.Ctx) error {
	var req struct {
		QuestionnaireID uint   `json:"questionnaire_id" validate:"required"`
		RoleID          uint   `json:"role_id" validate:"required"`
		PermissionIDs   []uint `json:"permission_ids" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	err := h.rolePermissionService.AssignPermissions(c.Context(), req.QuestionnaireID, req.RoleID, req.PermissionIDs)
	if err != nil {
		if errors.As(err, &gorm.ErrForeignKeyViolated) {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "role and permission id error", nil)
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Permissions assigned successfully")
}

func (h *QuestionnaireRolePermissionHandler) RemovePermissions(c *fiber.Ctx) error {
	var req struct {
		QuestionnaireID uint   `json:"questionnaire_id" validate:"required"`
		RoleID          uint   `json:"role_id" validate:"required"`
		PermissionIDs   []uint `json:"permission_ids" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	err := h.rolePermissionService.RemovePermissions(c.Context(), req.QuestionnaireID, req.RoleID, req.PermissionIDs)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Permissions removed successfully")
}

func (h *QuestionnaireRolePermissionHandler) GetRolePermissions(c *fiber.Ctx) error {
	roleID, err := strconv.ParseUint(c.Params("role_id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid role ID", nil)
	}

	questionnaireID, err := strconv.ParseUint(c.Query("questionnaire_id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid questionnaire ID", nil)
	}

	permissions, err := h.rolePermissionService.GetRolePermissions(c.Context(), uint(questionnaireID), uint(roleID))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, permissions)
}
