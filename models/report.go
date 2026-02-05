package models

type BestSellerResponse struct {
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity_sold"`
}

type DailyReportResponse struct {
	TotalRevenue   int                 `json:"total_revenue"`
	TotalOrders    int                 `json:"total_orders"`
	BestSeller     *BestSellerResponse `json:"best_seller"`
}
