package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/models"
	"log"
	"time"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) GetAllTransactions() ([]models.TransactionResponse, error) {
	rows, err := repo.db.Query(`
		SELECT
			t.id,
			t.total_amount,
			t.created_at,
			td.id,
			p.name,
			td.quantity,
			td.subtotal
		FROM transactions t
		LEFT JOIN transaction_details td
			ON t.id = td.transaction_id
		LEFT JOIN products p
			ON td.product_id = p.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactionMap := make(map[int]*models.TransactionResponse)

	for rows.Next() {
		var (
			tID         int
			totalAmount int
			createdAt   time.Time
			detailID      sql.NullInt64
			productName   sql.NullString
			quantity      sql.NullInt64
			subtotal      sql.NullInt64
		)

		if err := rows.Scan(
			&tID,
			&totalAmount,
			&createdAt,
			&detailID,
			&productName,
			&quantity,
			&subtotal,
		); err != nil {
			return nil, err
		}

		if _, exists := transactionMap[tID]; !exists {
			transactionMap[tID] = &models.TransactionResponse{
				ID:          tID,
				TotalAmount: totalAmount,
				CreatedAt:   createdAt,
				Details:     []models.TransactionDetailResponse{},
			}
		}

		if detailID.Valid {
			transactionMap[tID].Details = append(
				transactionMap[tID].Details,
				models.TransactionDetailResponse{
					ID:            int(detailID.Int64),
					ProductName:   productName.String,
					Quantity:      int(quantity.Int64),
					Subtotal:      int(subtotal.Int64),
				},
			)
		}
	}

	transactions := make([]models.TransactionResponse, 0, len(transactionMap))
	for _, t := range transactionMap {
		transactions = append(transactions, *t)
	}

	return transactions, nil
}

func (repo *TransactionRepository) GetTransactionByID(id int) (*models.TransactionResponse, error) {
	rows, err := repo.db.Query(`
		SELECT
			t.id,
			t.total_amount,
			t.created_at,
			td.id,
			p.name,
			td.quantity,
			td.subtotal
		FROM transactions t
		LEFT JOIN transaction_details td
			ON t.id = td.transaction_id
		LEFT JOIN products p
			ON td.product_id = p.id
		WHERE t.id = $1
		ORDER BY td.id
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transaction *models.TransactionResponse

	for rows.Next() {
		var (
			tID         int
			totalAmount int
			createdAt   time.Time

			detailID    sql.NullInt64
			productName sql.NullString
			quantity    sql.NullInt64
			subtotal    sql.NullInt64
		)

		if err := rows.Scan(
			&tID,
			&totalAmount,
			&createdAt,
			&detailID,
			&productName,
			&quantity,
			&subtotal,
		); err != nil {
			return nil, err
		}

		if transaction == nil {
			transaction = &models.TransactionResponse{
				ID:          tID,
				TotalAmount: totalAmount,
				CreatedAt:   createdAt,
				Details:     []models.TransactionDetailResponse{},
			}
		}

		if detailID.Valid {
			transaction.Details = append(
				transaction.Details,
				models.TransactionDetailResponse{
					ID:          int(detailID.Int64),
					ProductName: productName.String,
					Quantity:    int(quantity.Int64),
					Subtotal:    int(subtotal.Int64),
				},
			)
		}
	}

	if transaction == nil {
		return nil, sql.ErrNoRows
	}

	return transaction, nil
}



func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Println("rollback error:", rbErr)
			}
		}
	}()

	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	for _, item := range items {
		var productPrice, stock int
		var productName string

		err := tx.QueryRow("SELECT name, price, stock FROM products WHERE id = $1", item.ProductID).Scan(&productName, &productPrice, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}
  
		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	var transactionID int
	var createdAt time.Time
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id, created_at", totalAmount).Scan(&transactionID, &createdAt)
	if err != nil {
		return nil, err
	}

	for i := range details {
		details[i].TransactionID = transactionID
		err = tx.QueryRow("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4) RETURNING id",
			transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal).Scan(&details[i].ID)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		CreatedAt:   createdAt,
		Details:     details,
	}, nil
}