package repositories

import (
	"database/sql"
	"kasir-api/models"
	"time"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (repo *ReportRepository) GetTodayReport() (*models.DailyReportResponse, error) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	today := time.Now().In(loc)

	return repo.GetReportByDateRange(today, today)
}

func (repo *ReportRepository) GetReportByDateRange(startDate, endDate time.Time) (*models.DailyReportResponse, error) {
	var report models.DailyReportResponse

	err := repo.db.QueryRow(`
		SELECT
			COALESCE(SUM(total_amount), 0),
			COUNT(*)
		FROM transactions
		WHERE created_at::date BETWEEN $1 AND $2
	`, startDate, endDate).Scan(&report.TotalRevenue, &report.TotalOrders)
	if err != nil {
		return nil, err
	}

	var (
		productName sql.NullString
		totalQty    sql.NullInt64
	)

	err = repo.db.QueryRow(`
		SELECT
			p.name,
			SUM(td.quantity) AS total_qty
		FROM transaction_details td
		JOIN transactions t ON td.transaction_id = t.id
		JOIN products p ON td.product_id = p.id
		WHERE t.created_at::date BETWEEN $1 AND $2
		GROUP BY p.id, p.name
		ORDER BY total_qty DESC
		LIMIT 1
	`, startDate, endDate).Scan(&productName, &totalQty)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if productName.Valid {
		report.BestSeller = &models.BestSellerResponse{
			ProductName: productName.String,
			Quantity:    int(totalQty.Int64),
		}
	} else {
		report.BestSeller = nil
	}

	return &report, nil
}
