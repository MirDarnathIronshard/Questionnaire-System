package handlers

import (
	"github.com/QBG-P2/Voting-System/internal/services"
	"github.com/QBG-P2/Voting-System/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type EmailHandler struct {
	emailService services.EmailService
}

type SendEmailRequest struct {
	To      string `json:"to" validate:"required,email"`
	Subject string `json:"subject" validate:"required"`
	Body    string `json:"body" validate:"required"`
}

type SendInvitationRequest struct {
	Email              string `json:"email" validate:"required,email"`
	QuestionnaireID    uint   `json:"questionnaire_id" validate:"required"`
	QuestionnaireTitle string `json:"questionnaire_title" validate:"required"`
}

type SendPasswordResetRequest struct {
	Email string `json:"email" validate:"required,email"`
}

func NewEmailHandler(emailService services.EmailService) *EmailHandler {
	return &EmailHandler{
		emailService: emailService,
	}
}

func (h *EmailHandler) SendEmail(c *fiber.Ctx) error {
	var req SendEmailRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	if err := h.emailService.SendEmail(c.UserContext(), req.To, req.Subject, req.Body); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to send email", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, map[string]string{
		"message": "Email queued for delivery",
	})
}
