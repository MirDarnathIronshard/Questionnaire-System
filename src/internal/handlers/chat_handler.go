package handlers

import (
	"strconv"

	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/internal/services"
	"github.com/QBG-P2/Voting-System/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type ChatHandler struct {
	Service services.ChatService
}

func NewChatHandler(service services.ChatService) *ChatHandler {
	return &ChatHandler{Service: service}
}

func (h *ChatHandler) CreateChat(c *fiber.Ctx) error {
	var chat models.Chat
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid user ID", nil)
	}

	if err := c.BodyParser(&chat); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	if err := chat.Validate(); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validation failed", err)
	}

	if err := h.Service.CreateGroupChat(c.UserContext(), chat.QuestionnaireID, userID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create chat", nil)
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, chat)
}

func (h *ChatHandler) SendMessage(c *fiber.Ctx) error {
	var message models.Message

	if err := c.BodyParser(&message); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	if err := message.Validate(); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validation failed", err)
	}

	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid user ID", nil)
	}

	message.UserID = userID

	if err := h.Service.SendMessage(c.UserContext(), message.ChatID, message.UserID, message.Content, message.AttachmentURL); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to send message", nil)
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, message)
}

func (h *ChatHandler) GetChatMessages(c *fiber.Ctx) error {
	chatID, err := strconv.ParseUint(c.Params("chat_id"), 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid chat ID", nil)
	}

	offset := utils.ParseIntDefault(c.Query("offset"), 0)
	limit := utils.ParseIntDefault(c.Query("limit"), 10)

	messages, err := h.Service.GetMessages(c.UserContext(), uint(chatID), offset, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get chat messages", nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, messages)
}

func (h *ChatHandler) GetUserChats(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid user ID", nil)
	}

	chats, err := h.Service.GetUserChats(c.UserContext(), userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get user chats", nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, chats)
}
