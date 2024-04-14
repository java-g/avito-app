package middlewares

import (
	"avito-app/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func IsUser(c *fiber.Ctx) error {
	JwtKey := JWT
	headerToken := c.Get("token")
	token, err := jwt.ParseWithClaims(headerToken, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JwtKey), nil
	})
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	_, ok := token.Claims.(*models.Claims)
	if !ok {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	return c.Next()
}
