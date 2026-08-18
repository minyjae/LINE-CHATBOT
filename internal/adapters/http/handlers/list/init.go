package list

import (
	util_controllers "minyjae/go-starter/utils/controller"
)

type IListController interface {
}

type listController struct {
	util_controllers.GenericController
}

func NewListController() IListController {
	return &listController{
		GenericController: util_controllers.NewGenericController(),
	}
}
