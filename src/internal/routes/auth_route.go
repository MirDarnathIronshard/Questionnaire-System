package routes

import (
	"github.com/QBG-P2/Voting-System/config"
	"github.com/QBG-P2/Voting-System/internal/handlers"
	"github.com/QBG-P2/Voting-System/internal/middleware"
	"github.com/QBG-P2/Voting-System/internal/repositories"
	"github.com/QBG-P2/Voting-System/internal/services"
)

func authRoute(rs *RouterService, cfg *config.Config) {

	api := rs.app.Group("/api")

	userRepo := repositories.NewUserRepository(rs.db)
	roleRepo := repositories.NewRoleRepository(rs.db)

	authService := services.NewAuthService(cfg, userRepo, roleRepo, rs.enforcer, rs.cfg.JWT.Secret)

	authHandler := handlers.NewAuthHandler(authService)

	api.Post("/register", authHandler.Register)
	api.Post("/login", authHandler.LoginUser)
	api.Post("/send-otp", authHandler.SendOtp, middleware.OtpLimiter(cfg))

}
