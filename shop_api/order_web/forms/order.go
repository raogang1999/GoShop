package forms

type ShopCartItemForm struct {
	GoodsId int32 `json:"goods" binding:"required"`
	Nums    int32 `json:"nums" binding:"required,min=1"`
}

type ShopCartUpdateItemForm struct {
	Nums    int32 `json:"nums" binding:"required,min=1"`
	Checked *bool `json:"checked"`
}

type CreateOrderForm struct {
	Address string `json:"address" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Mobile  string `json:"mobile" binding:"required,mobile"`
	Post    string `json:"post" binding:"required"`
}
