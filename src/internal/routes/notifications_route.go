package routes

import (
	"github.com/QBG-P2/Voting-System/internal/handlers"
)

func NotificationRoute(rs *RouterService) {

	app := rs.app

	notifications := app.Group("/api/notifications")

	h := handlers.NewNotificationHandler(rs.notificationService)
	notifications.Post("/", h.CreateNotification)
	notifications.Get("/", h.GetUserNotifications)
	notifications.Put("/:id/read", h.MarkAsRead)
	notifications.Delete("/:id", h.DeleteNotification)
	notifications.Get("/unread-count", h.GetUnreadCount)
}
