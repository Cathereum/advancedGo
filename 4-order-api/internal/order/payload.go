package order

type OrderCreateRequest struct {
	ProductIds []uint `json:"product_ids"`
}
