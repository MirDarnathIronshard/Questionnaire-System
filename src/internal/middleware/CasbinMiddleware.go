package middleware

import (
	"github.com/casbin/casbin/v2"
	"github.com/gofiber/fiber/v2"
)

func CasbinMiddleware(enforcer *casbin.Enforcer) fiber.Handler {
	return func(c *fiber.Ctx) error {

		user := c.Locals("username").(string)
		obj := c.Path()
		act := c.Method()

		allowed, err := enforcer.Enforce(user, obj, act)
		if err != nil {
			return fiber.ErrInternalServerError
		}
		if !allowed {
			return fiber.ErrForbidden
		}
		return c.Next()
	}
}
