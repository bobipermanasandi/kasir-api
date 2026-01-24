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


var products = []Product{
	{ID: 1, Name: "Buku Tulis", Price: 5000, Stock: 10},
	{ID: 2, Name: "Pensil", Price: 2000, Stock: 20},
	{ID: 3, Name: "Penghapus", Price: 1000, Stock: 30},
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
