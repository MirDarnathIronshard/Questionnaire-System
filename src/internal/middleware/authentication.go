package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"

	"github.com/QBG-P2/Voting-System/config"
	"github.com/QBG-P2/Voting-System/internal/services"
	"github.com/QBG-P2/Voting-System/pkg/service_errors"
)

func Authentication(cfg *config.Config) func(c *fiber.Ctx) {
	var tokenService = services.NewTokenService(cfg)
	return func(c *fiber.Ctx) {
		var err error
		claimMap := map[string]interface{}{}
		auth := c.Get("Authorization")
		token := strings.Split(auth, " ")
		if auth == "" || len(token) < 2 {
			err = &service_errors.ServiceError{EndUserMessage: service_errors.TokenRequired}
		} else {
			claimMap, err = tokenService.GetClaims(token[1])
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					err = &service_errors.ServiceError{EndUserMessage: service_errors.TokenExpired}
				default:
					err = &service_errors.ServiceError{EndUserMessage: service_errors.TokenInvalid}
				}
			}
		}
		c.Locals("UserId", claimMap["UserId"])
		c.Next()
	}
}
