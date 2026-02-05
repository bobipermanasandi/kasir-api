package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type TransactionService struct {
	repo *repositories.TransactionRepository
}

func NewTransactionService(repo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) GetAllTransactions() ([]models.TransactionResponse, error) {
	return s.repo.GetAllTransactions()
}

func (s *TransactionService) GetTransactionByID(id int) (*models.TransactionResponse, error) {
	return s.repo.GetTransactionByID(id)
}

func (s *TransactionService) Checkout(items []models.CheckoutItem, useLock bool) (*models.Transaction, error) {
	return s.repo.CreateTransaction(items)
}