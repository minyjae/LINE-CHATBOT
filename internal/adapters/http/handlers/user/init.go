package user

import (
	util_controllers "minyjae/go-starter/utils/controller"
)

type IUserController interface {
	// GetUser(c *fiber.Ctx) error
}

type userController struct {
	util_controllers.GenericController
	// userService services.UserService
}

func NewUserController() IUserController {
	return &userController{
		GenericController: util_controllers.NewGenericController(),
		// userService:       userService,
	}
}
