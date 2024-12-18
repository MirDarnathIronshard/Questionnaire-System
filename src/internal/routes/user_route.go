package routes

import (
	"github.com/QBG-P2/Voting-System/internal/handlers"
	"github.com/QBG-P2/Voting-System/internal/repositories"
	"github.com/QBG-P2/Voting-System/internal/services"
)

func userRoutes(rs *RouterService) {
	db := rs.db
	app := rs.app
	enforcer := rs.enforcer

	userRepo := repositories.NewUserRepository(db)
	roleRepo := repositories.NewRoleRepository(db)
	userService := services.NewUserService(userRepo, roleRepo, enforcer)
	handler := handlers.NewUserHandler(userService)

	// Create users group
	users := app.Group("/api/users")

	// User profile routes
	users.Get("/profile", handler.GetProfile)
	users.Put("/profile", handler.UpdateProfile)
	users.Get("/wallet", handler.GetWalletBalance)
	// User role management routes
	users.Post("/:id/roles", handler.AssignRole)
	users.Delete("/:id/roles", handler.RemoveRole)

	// User management routes
	users.Get("/:id", handler.GetUser)
	users.Put("/:id", handler.UpdateUser)
	users.Delete("/:id", handler.DeleteUser)

}
