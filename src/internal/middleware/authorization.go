package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func Authorization(validRoles []string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		rolesVal := c.Locals("roles")
		fmt.Println(rolesVal)
		if rolesVal == nil {
			return fiber.ErrForbidden
		}
		roles := rolesVal.([]interface{})
		val := map[string]int{}
		for _, item := range roles {
			val[item.(string)] = 0
		}
		for _, item := range validRoles {
			if _, ok := val[item]; ok {
				c.Next()
				return nil
			}
		}
		return fiber.ErrForbidden
	}
}
