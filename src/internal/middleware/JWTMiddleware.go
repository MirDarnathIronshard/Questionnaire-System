package middleware

import (
	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/pkg/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
	"strings"
)

func JWTMiddleware(secretKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")
		tokenQuery := c.Query("token", "")
		if tokenQuery != "" {
			tokenString = tokenQuery
		}

		if tokenString == "" {
			return fiber.ErrUnauthorized
		}
		tokenString = strings.ReplaceAll(tokenString, "Bearer ", "")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
		if err != nil || !token.Valid {
			return fiber.ErrUnauthorized
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return fiber.ErrUnauthorized
		}

		username, ok := claims["username"].(string)
		if !ok {
			return fiber.ErrUnauthorized
		}

		userId, ok := claims["userID"].(float64)
		if !ok {
			return fiber.ErrUnauthorized
		}

		roles, ok := claims["roles"].([]interface{})
		if !ok {
			return fiber.ErrUnauthorized
		}

		var roleList []models.Role
		for _, role := range roles {
			r := role.(map[string]interface{})
			roleList = append(roleList, models.Role{
				Name: r["name"].(string),
				Model: gorm.Model{
					ID: uint(r["ID"].(float64)),
				},
			})
		}

		c.Locals("username", username)
		c.Locals("userID", uint(userId))

		if len(roleList) == 0 {
			roleList = append(roleList, models.Role{})
		}

		ctx := auth.CreateAuthContext(c.UserContext(), &models.User{
			Email: username,
			Model: gorm.Model{
				ID: uint(userId),
			},
			Roles: roleList,
		}, claims)

		c.SetUserContext(ctx)

		return c.Next()
	}
}
