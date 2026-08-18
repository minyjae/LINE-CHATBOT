package utils

import (
	"errors"
	"log"
	"runtime"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Response Structs
type CommonResponse struct {
	Message string      `json:"message"`
	Code    string      `json:"code,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Detail  string      `json:"detail,omitempty"`
}

type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

type PaginateResponse struct {
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

// Interface
type IResponse interface {
	Item(c *fiber.Ctx, message string, data interface{}) error
	Paginate(c *fiber.Ctx, message string, data interface{}, page Pagination) error
	Created(c *fiber.Ctx, message string, data interface{}) error
	Updated(c *fiber.Ctx, message string, data interface{}) error
	Deleted(c *fiber.Ctx, message string) error
	BadRequest(c *fiber.Ctx, message string, code string) error
	NotFound(c *fiber.Ctx, message string, code string) error
	Forbidden(c *fiber.Ctx, message string) error
	Unauthorized(c *fiber.Ctx, message string, code string) error
	ValidateFailed(c *fiber.Ctx, message string, data interface{}) error
	InternalServerError(c *fiber.Ctx, message string, detail string, code string) error
	ErrorHandler(c *fiber.Ctx, err error) error
}

type Response struct{}

func NewResponse() IResponse {
	return &Response{}
}

// --- Implement Methods ---

func (root *Response) Item(c *fiber.Ctx, message string, data interface{}) error {
	resData := CommonResponse{Message: message, Data: data}
	log.Printf("[INFO] %d %s", fiber.StatusOK, message)
	return c.Status(fiber.StatusOK).JSON(resData)
}

func (root *Response) Paginate(c *fiber.Ctx, message string, data interface{}, page Pagination) error {
	resData := PaginateResponse{Message: message, Data: data, Pagination: page}
	log.Printf("[INFO] %d %s", fiber.StatusOK, message)
	return c.Status(fiber.StatusOK).JSON(resData)
}

func (root *Response) Created(c *fiber.Ctx, message string, data interface{}) error {
	resData := CommonResponse{Message: message, Data: data}
	log.Printf("[INFO] %d %s", fiber.StatusCreated, message)
	return c.Status(fiber.StatusCreated).JSON(resData)
}

func (root *Response) Updated(c *fiber.Ctx, message string, data interface{}) error {
	resData := CommonResponse{Message: message, Data: data}
	log.Printf("[INFO] %d %s", fiber.StatusOK, message)
	return c.Status(fiber.StatusOK).JSON(resData)
}

func (root *Response) Deleted(c *fiber.Ctx, message string) error {
	resData := CommonResponse{Message: message}
	log.Printf("[INFO] %d %s", fiber.StatusOK, message)
	return c.Status(fiber.StatusOK).JSON(resData)
}

func (root *Response) BadRequest(c *fiber.Ctx, message string, code string) error {
	resData := CommonResponse{Message: message, Code: code}
	log.Printf("[WARN] %d %s", fiber.StatusBadRequest, message)
	return c.Status(fiber.StatusBadRequest).JSON(resData)
}

func (root *Response) NotFound(c *fiber.Ctx, message string, code string) error {
	resData := CommonResponse{Message: message, Code: code}
	log.Printf("[WARN] %d %s", fiber.StatusNotFound, message)
	return c.Status(fiber.StatusNotFound).JSON(resData)
}

func (root *Response) Unauthorized(c *fiber.Ctx, message string, code string) error {
	resData := CommonResponse{Message: message, Code: code}
	log.Printf("[WARN] %d %s", fiber.StatusUnauthorized, message)
	return c.Status(fiber.StatusUnauthorized).JSON(resData)
}

func (root *Response) Forbidden(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Forbidden"
	}
	resData := CommonResponse{Message: message}
	log.Printf("[WARN] %d %s", fiber.StatusForbidden, message)
	return c.Status(fiber.StatusForbidden).JSON(resData)
}

func (root *Response) ValidateFailed(c *fiber.Ctx, message string, data interface{}) error {
	resData := CommonResponse{Message: message, Data: data}
	log.Printf("[WARN] %d %s", fiber.StatusUnprocessableEntity, message)
	return c.Status(fiber.StatusUnprocessableEntity).JSON(resData)
}

func (root *Response) InternalServerError(c *fiber.Ctx, message, detail, code string) error {
	resData := CommonResponse{Message: message, Detail: detail, Code: code}

	_, file, line, _ := runtime.Caller(1)
	fileLine := strings.Replace(file, "/app/", "", 1) + ":" + strconv.Itoa(line)
	errorMessage := fileLine + " " + message

	log.Printf("[ERROR] SENTRY CAPTURE: %s | Detail: %s", errorMessage, detail)

	return c.Status(fiber.StatusInternalServerError).JSON(resData)
}

func (root *Response) ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		code = fiberErr.Code
	}

	switch code {
	case fiber.StatusNotFound:
		return root.NotFound(c, err.Error(), "NOT_FOUND")
	default:
		return root.InternalServerError(c, "Internal server error", err.Error(), "INTERNAL_SERVER_ERROR")
	}
}
