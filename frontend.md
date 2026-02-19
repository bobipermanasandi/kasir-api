bikinin aplikasi kasir
bisa
- get all product
- get detail product
- update product
- delete product
- checkout
api bisa di akses di https://kasir-api-production-b87e.up.railway.app/api/v1

list API sebagai berikut

a. API Product : 
- GET /products
- GET /products/id
- GET /products?name="name"
- PUT /products/id
- DELETE /products/id

b. API Categories : 
- GET /categories
- GET /categories/id
- PUT /categories/id
- DELETE /categories/id

c. API Transaction
- GET /transactions
- GET /transactions/id
- POST /checkout

d. report 
- GET /report/today
- GET /report?start_date=2026-02-06&end_date=2026-02-06


untuk model nya

type ProductRequest struct {
    Name       string `json:"name"`
    Price      int    `json:"price"`
    Stock      int    `json:"stock"`
    CategoryID int    `json:"category_id"`
}

type ProductListResponse struct {
    ID       int        `json:"id"`
    Name     string     `json:"name"`
    Price    int        `json:"price"`
    Stock    int        `json:"stock"`
    Category string     `json:"category"`
}

type ProductDetailResponse struct {
    ID       int        `json:"id"`
    Name     string     `json:"name"`
    Price    int        `json:"price"`
    Stock    int        `json:"stock"`
    Category Category   `json:"category"`
}

type Category struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Description string `json:"description"`
}

type CheckoutItem struct {
    ProductID int `json:"product_id"`
    Quantity  int `json:"quantity"`
}

type CheckoutRequest struct {
    Items   []CheckoutItem `json:"items"`
}

type BestSellerResponse struct {
    ProductName string `json:"product_name"`
    Quantity    int    `json:"quantity_sold"`
}

type DailyReportResponse struct {
    TotalRevenue   int                 `json:"total_revenue"`
    TotalOrders    int                 `json:"total_orders"`
    BestSeller     *BestSellerResponse `json:"best_seller"`
}

type Transaction struct {
    ID          int                 `json:"id"`
    TotalAmount int                 `json:"total_amount"`
    CreatedAt   time.Time           `json:"created_at"`
    Details     []TransactionDetail `json:"details"`
}




type TransactionDetail struct {
    ID            int    `json:"id"`
    TransactionID int    `json:"transaction_id"`
    ProductID     int    `json:"product_id"`
    ProductName   string `json:"product_name,omitempty"`
    Quantity      int    `json:"quantity"`
    Subtotal      int    `json:"subtotal"`
}

type TransactionResponse struct {
    ID          int                 `json:"id"`
    TotalAmount int                 `json:"total_amount"`
    CreatedAt   time.Time           `json:"created_at"`
    Details     []TransactionDetailResponse `json:"details"`
}

type TransactionDetailResponse struct {
    ID            int    `json:"id"`
    ProductName   string `json:"product_name,omitempty"`
    Quantity      int    `json:"quantity"`
    Subtotal      int    `json:"subtotal"`
}


api get detail, update, delete, checkout perlu X-API-Key=secretkeyapikey