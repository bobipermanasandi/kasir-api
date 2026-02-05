package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type ProductService struct {
	repo *repositories.ProductRepository
}

func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetAll(name string) ([]*models.ProductListResponse, error) {
	return s.repo.GetAll(name)
}

func (s *ProductService) GetByID(id int) (*models.ProductDetailResponse, error) {
	return s.repo.GetByID(id)
}

func (s *ProductService) Create(req *models.ProductRequest) (*models.ProductDetailResponse, error) {
	return s.repo.Create(req)
}

func (s *ProductService) Update(id int, req *models.ProductRequest) (*models.ProductDetailResponse, error) {
	return s.repo.Update(id, req)
}

func (s *ProductService) Delete(id int) error {
	return s.repo.Delete(id)
}
