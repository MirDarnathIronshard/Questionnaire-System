package routes

import (
	"github.com/QBG-P2/Voting-System/internal/handlers"
	"github.com/QBG-P2/Voting-System/internal/repositories"
	"github.com/QBG-P2/Voting-System/internal/services"
)

func messageRoute(rs *RouterService) {
	db := rs.db
	app := rs.app

	messageRepo := repositories.NewMessageRepository(db)
	messageServ := services.NewMessageService(messageRepo)
	controller := handlers.NewMessageHandler(messageServ)

	api := app.Group("/api/messages")

	api.Post("/", controller.CreateMessage)
	api.Get("/chat/:chat_id", controller.GetMessagesByChatID)
	api.Get("/:id", controller.GetMessageByID)
	api.Put("/:id", controller.UpdateMessage)
	api.Delete("/:id", controller.DeleteMessage)
}
