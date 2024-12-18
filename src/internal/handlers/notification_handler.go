package handlers

import (
	"context"
	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/internal/services"
	"github.com/QBG-P2/Voting-System/internal/utils"
	"github.com/gofiber/fiber/v2"
	"log"
	"strconv"
)

type NotificationHandler struct {
	service services.NotificationService
}

func NewNotificationHandler(service services.NotificationService) *NotificationHandler {
	handler := &NotificationHandler{
		service: service,
	}

	go func() {
		if err := service.StartNotificationConsumer(context.Background()); err != nil {
			log.Printf("Error starting notification consumer: %v", err)
		}
	}()

	return handler
}
func (h *NotificationHandler) CreateNotification(c *fiber.Ctx) error {
	var notification models.Notification

	userID := c.Params("user_id", "")
	if userID == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "need user_id to send notification", nil)
	}

	if err := c.BodyParser(&notification); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	notification.UserID = uint(utils.ParseIntDefault(userID, 0))

	if err := h.service.Create(c.UserContext(), &notification); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create notification", err)
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, notification)
}

func (h *NotificationHandler) GetUserNotifications(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	params := utils.ParsePaginationParams(c)

	notifications, total, err := h.service.GetUserNotifications(c.UserContext(), userID, params.Page, params.PageSize)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch notifications", err)
	}

	totalPages := (int(total) + params.PageSize - 1) / params.PageSize

	return utils.PaginatedResponseWrapper(c, fiber.StatusOK, notifications, models.Pagination{
		Page:       params.Page,
		PageSize:   params.PageSize,
		Total:      int(total),
		TotalPages: totalPages,
	})
}

func (h *NotificationHandler) MarkAsRead(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid notification ID", nil)
	}

	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	if err := h.service.MarkAsRead(c.UserContext(), uint(id), userID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to mark notification as read", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Notification marked as read")
}

func (h *NotificationHandler) DeleteNotification(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid notification ID", nil)
	}

	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	if err := h.service.Delete(c.UserContext(), uint(id), userID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete notification", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Notification deleted successfully")
}

func (h *NotificationHandler) GetUnreadCount(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	count, err := h.service.GetUnreadCount(c.UserContext(), userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get unread count", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, map[string]int64{"unread_count": count})
}
