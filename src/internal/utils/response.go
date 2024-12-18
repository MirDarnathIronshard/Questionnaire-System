package utils

import (
	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/gofiber/fiber/v2"
	"reflect"
)

func SuccessResponse(c *fiber.Ctx, status int, data interface{}) error {
	response := fiber.Map{
		"status": "success",
		"data":   data,
	}
	return c.Status(status).JSON(response)
}

func ErrorResponse(c *fiber.Ctx, status int, message string, errors interface{}) error {
	if errors == nil || reflect.TypeOf(errors).Kind() == reflect.String {
		var errorList []ErrorDetail
		errorList = append(errorList, ErrorDetail{
			Message: message,
			Field:   "message",
		})
		errors = errorList
	}

	response := fiber.Map{
		"status":  "error",
		"message": message,
		"errors":  errors,
	}
	return c.Status(status).JSON(response)
}

func PaginatedResponseWrapper(c *fiber.Ctx, status int, data interface{}, pagination models.Pagination) error {
	response := models.PaginatedResponse{
		Data:       data,
		Pagination: pagination,
	}
	return c.Status(status).JSON(fiber.Map{
		"status":     "success",
		"data":       response.Data,
		"pagination": response.Pagination,
	})
}
