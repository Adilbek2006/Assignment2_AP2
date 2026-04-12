package http

type CreateOrderRequest struct {
	CustomerID string `json:"customer_id"`
	ItemName   string `json:"item_name"`
	Amount     int64  `json:"amount"`
}

type OrderStatsResponse struct {
	Total     int64 `json:"total"`
	Pending   int64 `json:"pending"`
	Paid      int64 `json:"paid"`
	Failed    int64 `json:"failed"`
	Cancelled int64 `json:"cancelled"`
}
