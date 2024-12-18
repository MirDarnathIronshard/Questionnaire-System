package handlers

import (
	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/internal/services"
	"github.com/QBG-P2/Voting-System/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PermissionHandler struct {
	PermissionService services.PermissionService
}

func NewPermissionHandler(ps services.PermissionService) *PermissionHandler {
	return &PermissionHandler{
		PermissionService: ps,
	}
}

func (h *PermissionHandler) CreatePermission(c *fiber.Ctx) error {
	var permission models.Permission
	if err := c.BodyParser(&permission); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Invalid request body", nil)
	}
	if err := permission.Validate(); err != nil {
		translatedErr, _ := utils.TranslateError(err, &permission)

		return utils.ErrorResponse(c, fiber.StatusBadRequest, translatedErr.Message, translatedErr.ErrorList)
	}

	err := h.PermissionService.Create(&permission)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create permission", nil)
	}

	return c.Status(fiber.StatusCreated).JSON(permission)
}

func (h *PermissionHandler) GetPermission(c *fiber.Ctx) error {
	idParam := c.Params("id")
	permissionID, err := strconv.Atoi(idParam)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid permission ID", nil)
	}

	permission, err := h.PermissionService.GetByID(uint(permissionID))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Permission not found", nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, permission)
}

func (h *PermissionHandler) UpdatePermission(c *fiber.Ctx) error {
	idParam := c.Params("id")
	permissionID, err := strconv.Atoi(idParam)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid permission ID", nil)
	}

	var permissionUpdate models.Permission
	if err := c.BodyParser(&permissionUpdate); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}
	if err = permissionUpdate.Validate(); err != nil {
		translatedErr, _ := utils.TranslateError(err, &permissionUpdate)

		return utils.ErrorResponse(c, fiber.StatusBadRequest, translatedErr.Message, translatedErr.ErrorList)
	}

	permissionUpdate.ID = uint(permissionID)
	if err = h.PermissionService.Update(&permissionUpdate); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update permission", nil)

	}

	updatedPermission, _ := h.PermissionService.GetByID(uint(permissionID))
	return c.JSON(updatedPermission)
}

func (h *PermissionHandler) DeletePermission(c *fiber.Ctx) error {
	idParam := c.Params("id")
	permissionID, err := strconv.Atoi(idParam)
	if err != nil {

		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid permission ID", nil)

	}

	if err = h.PermissionService.Delete(uint(permissionID)); err != nil {

		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete permission", nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Permission deleted successfully")
}
