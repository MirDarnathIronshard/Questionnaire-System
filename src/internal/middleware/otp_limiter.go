package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"golang.org/x/time/rate"

	"github.com/QBG-P2/Voting-System/config"
	"github.com/QBG-P2/Voting-System/pkg/limiter"
)

func OtpLimiter(cfg *config.Config) fiber.Handler {
	var limiter = limiter.NewIPRateLimiter(rate.Every(cfg.Otp.Limiter*time.Second), 1)
	return func(c *fiber.Ctx) error {
		limiter := limiter.GetLimiter(c.Context().RemoteAddr().String())
		if !limiter.Allow() {
			c.Status(fiber.StatusTooManyRequests)
			c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized access",
			})
		} else {
			c.Next()
		}
		return nil
	}
}
