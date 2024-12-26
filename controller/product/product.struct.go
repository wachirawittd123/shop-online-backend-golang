package productController

// AddProductRequest defines the expected structure of the request body
type AddProductRequest struct {
	Name       string `json:"name" binding:"required"`
	Price      int    `json:"price" binding:"required"`
	Detail     string `json:"detail"`
	IDCategory string `json:"id_category" binding:"required"`
}
