package userController

import (
	models "github.com/wachirawittd123/shop-online-backend-golang/model"
)

type RequestUpdateBody struct {
	Name         string                 `json:"name"`
	ShippingAddr models.ShippingAddress `json:"address"`
}
