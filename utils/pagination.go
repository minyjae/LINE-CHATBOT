package utils

import (
	"strconv"
	"slices"

	"github.com/gofiber/fiber/v2"
)

type IPagination interface {
	GetPagination(c *fiber.Ctx, allowOrderField []string) (*RespPagination, error)
}

type RespPagination struct {
	Page           int    `json:"page"`
	PerPage        int    `json:"per_page"`
	Offset         int    `json:"offset"`
	OrderField     string `json:"order_field"`
	OrderDirection string `json:"order_direction"`
	Search         string `json:"search"`
}

type PaginationUtil struct{}

func NewPagination() IPagination {
	return &PaginationUtil{}
}

func (p *PaginationUtil) GetPagination(c *fiber.Ctx, allowOrderField []string) (*RespPagination, error) {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid page")
	}

	perPage, err := strconv.Atoi(c.Query("per_page", "10"))
	if err != nil || perPage < 1 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid per_page")
	}

	offset := (page - 1) * perPage

	orderField := c.Query("order_field", "")
	if orderField != "" && !slices.Contains(allowOrderField, orderField) {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid order field")
	}

	orderDirection := c.Query("order_direction", "desc")
	if orderDirection != "asc" && orderDirection != "desc" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid order direction")
	}

	search := c.Query("search", "")

	return &RespPagination{
		Page:           page,
		PerPage:        perPage,
		Offset:         offset,
		OrderField:     orderField,
		OrderDirection: orderDirection,
		Search:         search,
	}, nil
}
