package userController

type RequestUpdateBody struct {
	Name string `json:"name" binding:"required"`
}
