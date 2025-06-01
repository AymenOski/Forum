package controller

import "forum/usecase"

type PostController struct {
	service *usecase.PostService
}

func NewPostController(service *usecase.PostService) *PostController {
	return &PostController{
		service: service,
	}
}


