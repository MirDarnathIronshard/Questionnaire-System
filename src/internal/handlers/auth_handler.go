package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/QBG-P2/Voting-System/internal/interfaces/http/request"
	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/internal/services"
	"github.com/QBG-P2/Voting-System/internal/utils"
)

type AuthHandler struct {
	Service    *services.AuthService
	otpService *services.OtpService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{Service: service}
}

func (h *AuthHandler) LoginUser(c *fiber.Ctx) error {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		models.BaseValidator
	}

	if err := c.BodyParser(&credentials); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "", nil)
	}
	if err := credentials.Validate(credentials); err != nil {
		translatedErr, _ := utils.TranslateError(err, &credentials)

		return utils.ErrorResponse(c, fiber.StatusBadRequest, translatedErr.Message, translatedErr.ErrorList)
	}
	token, err := h.Service.LoginUser(credentials.Email, credentials.Password)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error(), err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, token)
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var user request.UserCreateRequest
	if err := c.BodyParser(&user); err != nil {
		err = user.Validate()
		translatedErr, _ := utils.TranslateError(err, &user)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, translatedErr.Message, translatedErr.ErrorList)

	}
	if err := user.Validate(); err != nil {
		translatedErr, _ := utils.TranslateError(err, &user)

		return utils.ErrorResponse(c, fiber.StatusBadRequest, translatedErr.Message, translatedErr.ErrorList)
	}
	_, err := h.Service.RegisterUser(user.Email, user.Password, user.NationalID, "user", user.Email)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
	}
	return utils.SuccessResponse(c, fiber.StatusCreated, "User registered successfully")
}
func (h *AuthHandler) SendOtp(c *fiber.Ctx) error {
	req := new(request.GetOtpRequest)
	err := c.BodyParser(&req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "", nil)
	}
	err = h.otpService.SendOtp(req.Email)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "", nil)
	}
	return utils.SuccessResponse(c, fiber.StatusCreated, "Otp sent successfully")

}
