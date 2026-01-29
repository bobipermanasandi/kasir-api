package models


type ProductRequest struct {
	Name       string `json:"name"`
	Price      int    `json:"price"`
	Stock      int    `json:"stock"`
	CategoryID int    `json:"category_id"`
}

type ProductListResponse struct {
	ID       int    	`json:"id"`
	Name     string 	`json:"name"`
	Price    int    	`json:"price"`
	Stock    int    	`json:"stock"`
	Category string		`json:"category"`
}

type ProductDetailResponse struct {
	ID       int    	`json:"id"`
	Name     string 	`json:"name"`
	Price    int    	`json:"price"`
	Stock    int    	`json:"stock"`
	Category Category	`json:"category"`
}

