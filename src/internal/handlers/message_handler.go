package handlers

import (
	"strconv"

	"github.com/QBG-P2/Voting-System/internal/interfaces/http/request"
	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/internal/services"
	"github.com/QBG-P2/Voting-System/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type MessageHandler struct {
	Service services.MessageService
}

func NewMessageHandler(service services.MessageService) *MessageHandler {
	return &MessageHandler{Service: service}
}

func (h *MessageHandler) CreateMessage(c *fiber.Ctx) error {
	var req request.CreateMessageRequest
	var message models.Message

	if err := c.BodyParser(&req); err != nil {
		errList := req.Validate()
		translatedErr, _ := utils.TranslateError(errList, &req)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", translatedErr.ErrorList)

	}

	if err := req.Validate(); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid user ID", nil)
	}

	message.UserID = userID

	if err := h.Service.CreateMessage(c.UserContext(), &message); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create message", nil)
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, message)
}

func (h *MessageHandler) UpdateMessage(c *fiber.Ctx) error {
	idParam := c.Params("id")
	messageID, err := strconv.Atoi(idParam)
	if err != nil || messageID <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid message ID", nil)
	}

	var req request.UpdateMessageRequest
	var message models.Message

	if err := c.BodyParser(&req); err != nil {
		errList := req.Validate()
		translatedErr, _ := utils.TranslateError(errList, &req)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", translatedErr.ErrorList)

	}

	message.ID = uint(messageID)

	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid user ID", nil)
	}

	isOwner, err := h.Service.IsOwnerMessage(c.UserContext(), message.ID, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error verifying message ownership", nil)
	}
	if !isOwner {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "You do not own this message", nil)
	}

	if err := h.Service.UpdateMessage(c.UserContext(), &message); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update message", nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, message)
}

func (h *MessageHandler) GetMessageByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	messageID, err := strconv.Atoi(idParam)
	if err != nil || messageID <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid message ID", nil)
	}

	message, err := h.Service.GetMessageByID(c.UserContext(), uint(messageID))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Message not found", nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, message)
}

func (h *MessageHandler) GetMessagesByChatID(c *fiber.Ctx) error {
	chatIDParam := c.Params("chat_id")
	chatID, err := strconv.Atoi(chatIDParam)
	if err != nil || chatID <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid chat ID", nil)
	}

	messages, err := h.Service.GetMessagesByChatID(c.UserContext(), uint(chatID), 1, 10)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve messages", nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, messages)
}

func (h *MessageHandler) DeleteMessage(c *fiber.Ctx) error {
	idParam := c.Params("id")
	messageID, err := strconv.Atoi(idParam)
	if err != nil || messageID <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid message ID", nil)
	}

	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid user ID", nil)
	}

	isOwner, err := h.Service.IsOwnerMessage(c.UserContext(), uint(messageID), userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error verifying message ownership", nil)
	}
	if !isOwner {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "You do not own this message", nil)
	}

	if err := h.Service.DeleteMessage(c.UserContext(), uint(messageID)); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete message", nil)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
