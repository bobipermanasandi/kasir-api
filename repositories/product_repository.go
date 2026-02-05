package repositories

import (
	"database/sql"
	"errors"
	"kasir-api/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// GetAll -> hanya category name
func (r *ProductRepository) GetAll(name string) ([]*models.ProductListResponse, error) {
	query := `
		SELECT p.id, p.name, p.price, p.stock, COALESCE(c.name,'Uncategories')
		FROM products p
		LEFT JOIN categories c ON c.id = p.category_id
		`

	var args []interface{}

	if name != "" {
		query += ` WHERE p.name ILIKE $1`
		args = append(args, "%"+name+"%")
	}
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*models.ProductListResponse
	for rows.Next() {
		var p models.ProductListResponse
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.Category); err != nil {
			return nil, err
		}
		products = append(products, &p)
	}
	return products, nil
}

// GetByID -> category lengkap
func (r *ProductRepository) GetByID(id int) (*models.ProductDetailResponse, error) {
	query := `
		SELECT p.id, p.name, p.price, p.stock,
		       COALESCE(c.id,0), COALESCE(c.name,'Uncategories'), COALESCE(c.description,'')
		FROM products p
		LEFT JOIN categories c ON c.id = p.category_id
		WHERE p.id=$1
	`
	var p models.ProductDetailResponse
	var catID int
	var catName, catDesc string
	err := r.db.QueryRow(query, id).Scan(
		&p.ID, &p.Name, &p.Price, &p.Stock,
		&catID, &catName, &catDesc,
	)
	if err != nil {
		return nil, err
	}
	p.Category = models.Category{
		ID:          catID,
		Name:        catName,
		Description: catDesc,
	}
	return &p, nil
}

// Create
func (r *ProductRepository) Create(req *models.ProductRequest) (*models.ProductDetailResponse, error) {
	query := `
		INSERT INTO products (name, price, stock, category_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var id int
	if err := r.db.QueryRow(query, req.Name, req.Price, req.Stock, req.CategoryID).Scan(&id); err != nil {
		return nil, err
	}
	return r.GetByID(id)
}

// Update
func (r *ProductRepository) Update(id int, req *models.ProductRequest) (*models.ProductDetailResponse, error) {
	query := `
		UPDATE products
		SET name=$1, price=$2, stock=$3, category_id=$4
		WHERE id=$5
	`
	res, err := r.db.Exec(query, req.Name, req.Price, req.Stock, req.CategoryID, id)
	if err != nil {
		return nil, err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return nil, errors.New("product not found")
	}
	return r.GetByID(id)
}

// Delete
func (r *ProductRepository) Delete(id int) error {
	res, err := r.db.Exec("DELETE FROM products WHERE id=$1", id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("product not found")
	}
	return nil
}
