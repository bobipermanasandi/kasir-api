package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
}


var products = []Product{
	{ID: 1, Name: "Buku Tulis", Price: 5000, Stock: 10},
	{ID: 2, Name: "Pensil", Price: 2000, Stock: 20},
	{ID: 3, Name: "Penghapus", Price: 1000, Stock: 30},
}

var categories = []Category{
	{ID: 1, Name: "Makanan", Description: "Makanan untuk di konsumsi"},
	{ID: 2, Name: "Minuman", Description: "Minuman untuk di konsumsi"},
	{ID: 3, Name: "Alat Tulis", Description: "Alat tulis untuk menulis"},
}



func getProductById(w http.ResponseWriter, r *http.Request) {
	idString := strings.TrimPrefix(r.URL.Path, "/api/v1/products/")
		id, err := strconv.Atoi(idString)
		if err != nil {
			_ = json.NewEncoder(w).Encode(map[string]string{
				"status":  "Error",
				"message": "Invalid product ID",
			})
			return
		}
		for _, product := range products {
			if product.ID == id {
				_ = json.NewEncoder(w).Encode(product)
				return
			}
		}
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status":  "Error",
			"message": "Product not found",
		})
}

func updateProductById(w http.ResponseWriter, r *http.Request) {
	idString := strings.TrimPrefix(r.URL.Path, "/api/v1/products/")
	id, err := strconv.Atoi(idString)
	if err != nil {
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status":  "Error",
			"message": "Invalid product ID",
		})
		return
	}
	for i, product := range products {
		if product.ID == id {
			var updatedProduct Product
			err := json.NewDecoder(r.Body).Decode(&updatedProduct)
			if err != nil {
				_ = json.NewEncoder(w).Encode(map[string]string{
					"status":  "Error",
					"message": "Invalid request body",
				})
				return
			}
			updatedProduct.ID = id
			products[i] = updatedProduct
			_ = json.NewEncoder(w).Encode(updatedProduct)
			return
		}
	}
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status":  "Error",
		"message": "Product not found",
	})
}

func deleteProductById(w http.ResponseWriter, r *http.Request) {
	idString := strings.TrimPrefix(r.URL.Path, "/api/v1/products/")
	id, err := strconv.Atoi(idString)
	if err != nil {
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status":  "Error",
			"message": "Invalid product ID",
		})
		return
	}
	for i, product := range products {
		if product.ID == id {
			products = append(products[:i], products[i+1:]...)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"status":  "Success",
				"message": "Product deleted successfully",
			})
			return
		}
	}
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status":  "Error",
		"message": "Product not found",
	})
}

func getCategoryById(w http.ResponseWriter, r *http.Request) {
	idString := strings.TrimPrefix(r.URL.Path, "/api/v1/categories/")
	id, err := strconv.Atoi(idString)
	if err != nil {
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status":  "Error",
			"message": "Invalid category ID",
		})
		return
	}
	for _, category := range categories {
		if category.ID == id {
			_ = json.NewEncoder(w).Encode(category)
			return
		}
	}
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status":  "Error",
		"message": "Category not found",
	})
}

func updateCategoryById(w http.ResponseWriter, r *http.Request) {
	idString := strings.TrimPrefix(r.URL.Path, "/api/v1/categories/")
	id, err := strconv.Atoi(idString)
	if err != nil {
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status":  "Error",
			"message": "Invalid category ID",
		})
		return
	}
	for i, category := range categories {
		if category.ID == id {
			var updatedCategory Category
			err := json.NewDecoder(r.Body).Decode(&updatedCategory)
			if err != nil {
				_ = json.NewEncoder(w).Encode(map[string]string{
					"status":  "Error",
					"message": "Invalid request body",
				})
				return
			}
			updatedCategory.ID = id
			categories[i] = updatedCategory
			_ = json.NewEncoder(w).Encode(updatedCategory)
			return
		}
	}
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status":  "Error",
		"message": "Category not found",
	})
}


func deleteCategoryById(w http.ResponseWriter, r *http.Request) {
	idString := strings.TrimPrefix(r.URL.Path, "/api/v1/categories/")
	id, err := strconv.Atoi(idString)
	if err != nil {
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status":  "Error",
			"message": "Invalid category ID",
		})
		return
	}
	for i, category := range categories {
		if category.ID == id {
			categories = append(categories[:i], categories[i+1:]...)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"status":  "Success",
				"message": "Category deleted successfully",
			})
			return
		}
	}
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status":  "Error",
		"message": "Category not found",
	})
}

func main() {

	// Get & Create Product -> http://localhost:8080/api/v1/products
	http.HandleFunc("/api/v1/products", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			_ = json.NewEncoder(w).Encode(products)
		case http.MethodPost:
			var product Product
			err := json.NewDecoder(r.Body).Decode(&product)
			if err != nil {
				_ = json.NewEncoder(w).Encode(map[string]string{
					"status":  "Error",
					"message": "Invalid request body",
				})
				return
			}

			product.ID = len(products) + 1
			products = append(products, product)
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"status":  "Success",
				"message": "Product added successfully",
			})
		default:
			_ = json.NewEncoder(w).Encode(map[string]string{
				"status":  "Error",
				"message": "Method not allowed",
			})
		}
	})

	// Get Product By ID -> http://localhost:8080/api/v1/products/:id
	// Update Product By ID -> http://localhost:8080/api/v1/products/:id
	// Delete Product By ID -> http://localhost:8080/api/v1/products/:id

	http.HandleFunc("/api/v1/products/{id}", func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case http.MethodGet:
			getProductById(w, r)

		case http.MethodPut:
			updateProductById(w, r)

		case http.MethodDelete:
			deleteProductById(w, r)
			
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"status":  "Error",
				"message": "Method not allowed",
			})
		}
	})

	// Get & Create Category -> http://localhost:8080/api/v1/categories
	http.HandleFunc("/api/v1/categories", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			_ = json.NewEncoder(w).Encode(categories)
		case http.MethodPost:
			var category Category
			err := json.NewDecoder(r.Body).Decode(&category)
			if err != nil {
				_ = json.NewEncoder(w).Encode(map[string]string{
					"status":  "Error",
					"message": "Invalid request body",
				})
				return
			}

			category.ID = len(categories) + 1
			categories = append(categories, category)
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"status":  "Success",
				"message": "Category added successfully",
			})
		default:
			_ = json.NewEncoder(w).Encode(map[string]string{
				"status":  "Error",
				"message": "Method not allowed",
			})
		}
	})

	// Get Category By ID -> http://localhost:8080/api/v1/categories/:id
	// Update Category By ID -> http://localhost:8080/api/v1/categories/:id
	// Delete Category By ID -> http://localhost:8080/api/v1/categories/:id

	http.HandleFunc("/api/v1/categories/{id}", func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case http.MethodGet:
			getCategoryById(w, r)

		case http.MethodPut:
			updateCategoryById(w, r)

		case http.MethodDelete:
			deleteCategoryById(w, r)
			
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"status":  "Error",
				"message": "Method not allowed",
			})
		}
	})

	// Health Check -> http://localhost:8080/health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "Service is running",
		})
	})

	fmt.Println("Server running di localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Gagal menjalankan server: ", err)
	}
}
