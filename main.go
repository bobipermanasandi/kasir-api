package main

import (
	"encoding/json"
	"fmt"
	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/repositories"
	"kasir-api/services"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
}

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}


func main() {

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
}

	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	http.HandleFunc("/api/v1/products", productHandler.HandleProducts)
	http.HandleFunc("/api/v1/products/", productHandler.HandleProductByID)

	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	http.HandleFunc("/api/v1/categories", categoryHandler.HandleCategories)
	http.HandleFunc("/api/v1/categories/", categoryHandler.HandleCategoryByID)


	// Health Check -> http://localhost:8080/api/v1/health
	http.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "Service is running",
		})
	})

	addr := ":" + config.Port
	fmt.Println("Server running di " + addr)

	error := http.ListenAndServe(addr, nil)
	if error != nil {
		fmt.Println("Gagal menjalankan server: ", error)
	}

	
}
