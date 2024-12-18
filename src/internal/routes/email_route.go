package routes

import (
	"github.com/QBG-P2/Voting-System/internal/handlers"
	"github.com/QBG-P2/Voting-System/internal/services"
	"github.com/QBG-P2/Voting-System/pkg/rabbitmq"
)

func EmailRoute(rs *RouterService) {
	app := rs.app
	cfg := rs.cfg

	emailGroup := app.Group("/api/emails")

	rmq, _ := rabbitmq.NewRabbitMQ(cfg)
	emailService := services.NewEmailService(cfg, rmq)
	h := handlers.NewEmailHandler(emailService)

	emailGroup.Post("/send", h.SendEmail)
}
