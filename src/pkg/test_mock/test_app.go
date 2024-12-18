package test_mock

import (
	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/internal/routes"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupTestApp() (*fiber.App, *gorm.DB) {
	db := SetupTestDB()
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", uint(1))
		return c.Next()
	})

	user := models.User{
		Password:     "Password",
		Email:        "test",
		Role:         "user",
		Is2FAEnabled: false,
		Wallet:       0,
	}
	db.Create(&user)

	routes.InitRoutes(app, db)
	return app, db
}
