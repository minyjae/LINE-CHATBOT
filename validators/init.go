package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/gobeam/stringy"
	"github.com/gofiber/fiber/v2"
	"minyjae/go-starter/utils"
)

type ErrorResponse struct {
	FailedField string `json:"failed_field"`
	Tag         string `json:"tag"`
	Value       string `json:"value"`
	Message     string `json:"message"`
}

func ValidateStruct[T any](c *fiber.Ctx) error {
	body := new(T)
	if err := c.BodyParser(body); err != nil {
		return utils.NewResponse().BadRequest(c, "Body parser error: "+err.Error(), "BODY_PARSER_ERROR")
	}

	var errors []*ErrorResponse
	validate := validator.New()

	if err := validate.Struct(body); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			errors = append(errors, &ErrorResponse{
				FailedField: stringy.New(e.Field()).SnakeCase("?", "").ToLower(),
				Tag:         e.Tag(),
				Value:       e.Param(),
			})
		}
	}

	if len(errors) > 0 {
		return utils.NewResponse().ValidateFailed(c, "", errors)
	}

	c.Locals("request", body)
	return c.Next()
}
