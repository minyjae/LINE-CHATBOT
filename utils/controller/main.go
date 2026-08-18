package util_controllers

import "minyjae/go-starter/utils"

type GenericController struct {
	Pagination utils.IPagination
	Response   utils.IResponse
}

func NewGenericController() GenericController {
	return GenericController{
		Pagination: utils.NewPagination(),
		Response:   utils.NewResponse(),
	}
}
