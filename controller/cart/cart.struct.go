package cartController

type RequestUpdateCart struct {
	ID          string             `json:"id"`
	SubTotal    float64            `json:"sub_total" binding:"required"`
	Total       float64            `json:"total" binding:"required"`
	Items       []RequestItemsCart `json:"items" binding:"required"`
	DeliveryFee float64            `json:"delivery_fee" binding:"required"`
	Status      string             `json:"status"`
}

type RequestItemsCart struct {
	ProductID string  `json:"product_id"`
	Qty       int     `json:"qty"`
	Total     float64 `json:"total"`
}

type RequestBuildMatchStage struct {
	CardID string `json:"id"`
	Search string `json:"search"`
}
