package handlers

import (
	"github.com/QBG-P2/Voting-System/internal/interfaces/http/request"
	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/internal/services"
	"github.com/QBG-P2/Voting-System/internal/utils"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"time"
)

type UserHandler struct {
	UserService services.UserService
}

func NewUserHandler(us services.UserService) *UserHandler {
	return &UserHandler{
		UserService: us,
	}
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	user, err := h.UserService.GetByID(userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found", nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, user)
}

func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	var req request.UpdateProfileRequest

	if err := c.BodyParser(&req); err != nil {

		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if err := req.Validate(); err != nil {
		translatedErr, _ := utils.TranslateError(err, &req)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, translatedErr.Message, translatedErr.ErrorList)
	}

	userID := c.Locals("userID").(uint)
	user, err := h.UserService.GetByID(userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found", nil)
	}

	if !req.BirthDate.IsZero() {
		registrationTime := user.CreatedAt
		if time.Since(registrationTime) > 24*time.Hour {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Birth date can only be changed within 24 hours of registration", nil)
		}
		user.BirthDate = req.BirthDate
	}

	if req.Email != "" {
		user.Email = req.Email
	}
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.City != "" {
		user.City = req.City
	}
	if req.Gender != "" {
		user.Gender = req.Gender
	}

	err = h.UserService.Update(user)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update user", nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, user)
}

// GetWalletBalance returns the user's current wallet balance
func (h *UserHandler) GetWalletBalance(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	user, err := h.UserService.GetByID(userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found", nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"balance": user.Wallet,
	})
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID", nil)
	}

	user, err := h.UserService.GetByID(uint(userID))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found", nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, user)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID", nil)
	}

	var userUpdate models.User
	if err := c.BodyParser(&userUpdate); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}
	if err = userUpdate.Validate(); err != nil {
		translatedErr, _ := utils.TranslateError(err, &userUpdate)

		return utils.ErrorResponse(c, fiber.StatusBadRequest, translatedErr.Message, translatedErr.ErrorList)
	}

	userUpdate.ID = uint(userID)
	err = h.UserService.Update(&userUpdate)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update user", nil)
	}

	updatedUser, _ := h.UserService.GetByID(uint(userID))
	return utils.SuccessResponse(c, fiber.StatusOK, updatedUser)
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID", nil)
	}

	err = h.UserService.Delete(uint(userID))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete user", nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"message": "User deleted successfully",
	})
}

func (h *UserHandler) AssignRole(c *fiber.Ctx) error {
	idParam := c.Params("id")
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID", nil)
	}

	var req struct {
		RoleID uint `json:"role_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	err = h.UserService.AssignRole(uint(userID), req.RoleID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to assign role to user", nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"message": "Role assigned to user successfully",
	})
}

func (h *UserHandler) RemoveRole(c *fiber.Ctx) error {
	idParam := c.Params("id")
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID", nil)
	}

	var req struct {
		RoleID uint `json:"role_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	err = h.UserService.RemoveRole(uint(userID), req.RoleID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to remove role from user", nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"message": "Role removed from user successfully",
	})
}
