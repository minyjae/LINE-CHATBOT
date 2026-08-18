package auth

import (
	util_controllers "minyjae/go-starter/utils/controller"
)

type IAuthController interface {
}

type authController struct {
	util_controllers.GenericController
}

func NewAuthController() IAuthController {
	return &authController{
		GenericController: util_controllers.NewGenericController(),
	}
}
