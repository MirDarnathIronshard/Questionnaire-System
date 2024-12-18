package routes

import (
	"github.com/QBG-P2/Voting-System/internal/handlers"
	"github.com/QBG-P2/Voting-System/internal/repositories"
	"github.com/QBG-P2/Voting-System/internal/services"
)

func chatRoutes(rs *RouterService) {
	db := rs.db
	app := rs.app

	userRepo := repositories.NewUserRepository(db)
	messageRepo := repositories.NewMessageRepository(db)
	chatRepo := repositories.NewChatRepository(db)
	service := services.NewChatService(userRepo, messageRepo, chatRepo)
	handler := handlers.NewChatHandler(service)

	chat := app.Group("/api/chat")
	chat.Post("/", handler.CreateChat)
	chat.Post("/message", handler.SendMessage)
	chat.Get("/:chat_id/messages", handler.GetChatMessages)
	chat.Get("/user", handler.GetUserChats)
}
