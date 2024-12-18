package main

import (
	"github.com/QBG-P2/Voting-System/config"
	"github.com/QBG-P2/Voting-System/internal/routes"
	"github.com/QBG-P2/Voting-System/pkg/cache"
	"github.com/QBG-P2/Voting-System/pkg/database"
	"github.com/QBG-P2/Voting-System/pkg/rabbitmq"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {
	cfg := config.GetConfig()
	err := database.InitDb(cfg)
	if err != nil {
		return
	}
	db := database.GetDb()
	database.Migrate(db)

	err = cache.InitRedis(cfg)
	if err != nil {
		panic(err)
		return
	}

	rmq, err := rabbitmq.NewRabbitMQ(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}
	defer func(rmq *rabbitmq.RabbitMQ) {
		err = rmq.Close()
		if err != nil {
			log.Fatalf("Failed to close RabbitMQ: %v", err)
		}
	}(rmq)

	enforcer := database.SetupCasbin(db)
	database.SuperAdminInit(enforcer, db, cfg)
	app := fiber.New()
	routes.InitRoutes(app, database.GetDb(), enforcer, rmq, cfg)

	log.Fatal(app.Listen(":" + cfg.Server.ExternalPort))
}
