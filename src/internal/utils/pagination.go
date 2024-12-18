package utils

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type PaginationParams struct {
	Page     int
	PageSize int
}

func ParsePaginationParams(c *fiber.Ctx) PaginationParams {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.Query("page_size", "10"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
	}
}
