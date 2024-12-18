package handlers

import (
	"errors"
	"fmt"
	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/internal/services"
	"github.com/QBG-P2/Voting-System/internal/utils"
	"gorm.io/gorm"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type RoleHandler struct {
	RoleService services.RoleService
}

func NewRoleHandler(rs services.RoleService) *RoleHandler {
	return &RoleHandler{
		RoleService: rs,
	}
}

func (h *RoleHandler) GetAllRoles(c *fiber.Ctx) error {
	roles, err := h.RoleService.GetAll()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get roles", err)
	}
	return utils.SuccessResponse(c, fiber.StatusOK, roles)
}

func (h *RoleHandler) GetRolePermissions(c *fiber.Ctx) error {
	roleID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid role ID", nil)
	}

	role, err := h.RoleService.GetByID(uint(roleID))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Role not found", nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, role.Permissions)
}
func (h *RoleHandler) GetUserRoles(c *fiber.Ctx) error {
	userID, err := strconv.ParseUint(c.Params("userId"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID", nil)
	}

	roles, err := h.RoleService.GetUserRoles(uint(userID))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to get user roles: %v", err), nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, roles)
}

func (h *RoleHandler) AssignRoleToUser(c *fiber.Ctx) error {
	userID, err := strconv.ParseUint(c.Params("userId"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID", nil)
	}

	var req struct {
		RoleID uint `json:"role_id" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	err = h.RoleService.AssignToUser(uint(userID), req.RoleID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to assign role", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"message": "Role assigned successfully",
	})
}
func (h *RoleHandler) RemoveRoleFromUser(c *fiber.Ctx) error {
	userID, err := strconv.ParseUint(c.Params("userId"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID", nil)
	}

	roleID, err := strconv.ParseUint(c.Params("roleId"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid role ID", nil)
	}

	err = h.RoleService.RemoveFromUser(uint(userID), uint(roleID))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to remove role", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"message": "Role removed successfully",
	})
}

func (h *RoleHandler) ValidateUserRole(c *fiber.Ctx) error {
	var req struct {
		UserID uint   `json:"user_id" validate:"required"`
		Role   string `json:"role" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	hasRole, err := h.RoleService.ValidateUserRole(req.UserID, req.Role)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to validate role", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"has_role": hasRole,
	})
}

func (h *RoleHandler) CreateRole(c *fiber.Ctx) error {
	var role models.Role
	if err := c.BodyParser(&role); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	if err := role.Validate(); err != nil {
		translatedErr, _ := utils.TranslateError(err, &role)

		return utils.ErrorResponse(c, fiber.StatusBadRequest, translatedErr.Message, translatedErr.ErrorList)
	}
	err := h.RoleService.Create(&role)
	if err != nil {
		if errors.As(err, &gorm.ErrDuplicatedKey) {
			return utils.ErrorResponse(c, fiber.StatusConflict, "Role already exists", nil)
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create role", nil)
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, role)
}

func (h *RoleHandler) GetRole(c *fiber.Ctx) error {
	idParam := c.Params("id")
	roleID, err := strconv.Atoi(idParam)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid role ID", nil)
	}

	role, err := h.RoleService.GetByID(uint(roleID))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Role not found", nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, role)
}

func (h *RoleHandler) UpdateRole(c *fiber.Ctx) error {
	idParam := c.Params("id")
	roleID, err := strconv.Atoi(idParam)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid role ID", nil)
	}

	var roleUpdate models.Role
	if err := c.BodyParser(&roleUpdate); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}
	if err := roleUpdate.Validate(); err != nil {
		translatedErr, _ := utils.TranslateError(err, &roleUpdate)

		return utils.ErrorResponse(c, fiber.StatusBadRequest, translatedErr.Message, translatedErr.ErrorList)
	}

	roleUpdate.ID = uint(roleID)
	err = h.RoleService.Update(&roleUpdate)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update role", nil)
	}

	updatedRole, _ := h.RoleService.GetByID(uint(roleID))
	return utils.SuccessResponse(c, fiber.StatusOK, updatedRole)
}

func (h *RoleHandler) DeleteRole(c *fiber.Ctx) error {
	idParam := c.Params("id")
	roleID, err := strconv.Atoi(idParam)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid role ID", nil)
	}

	err = h.RoleService.Delete(uint(roleID))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete role", nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"message": "Role deleted successfully",
	})
}

func (h *RoleHandler) AssignPermission(c *fiber.Ctx) error {
	idParam := c.Params("id")
	roleID, err := strconv.Atoi(idParam)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid role ID", nil)
	}

	var req struct {
		PermissionID         uint `json:"permission_id"`
		models.BaseValidator `json:"-"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	if err = req.Validate(req); err != nil {
		translatedErr, _ := utils.TranslateError(err, &req)

		return utils.ErrorResponse(c, fiber.StatusBadRequest, translatedErr.Message, translatedErr.ErrorList)
	}

	err = h.RoleService.AssignPermission(uint(roleID), req.PermissionID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to assign permission to role", nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"message": "Permission assigned to role successfully",
	})
}

func (h *RoleHandler) RemovePermission(c *fiber.Ctx) error {
	idParam := c.Params("id")
	roleID, err := strconv.Atoi(idParam)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid role ID", nil)
	}

	var req struct {
		PermissionID uint `json:"permission_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	err = h.RoleService.RemovePermission(uint(roleID), req.PermissionID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to remove permission from role", nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"message": "Permission removed from role successfully",
	})
}
