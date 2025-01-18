package service

import (
	"main.go/entity"
	"main.go/repository"
)

type ProductService interface {
	CreateCategory(category *entity.Category) error
	GetAllCategories() ([]entity.Category, error)
	GetCategoryByID(id uint) (*entity.Category, error)
	UpdateCategory(category *entity.Category) error
	DeleteCategory(id uint) error

	CreateProduct(product *entity.Product) error
	GetAllProducts() ([]entity.Product, error)
	GetProductByID(id uint) (*entity.Product, error)
	UpdateProduct(product *entity.Product) error
	DeleteProduct(id uint) error
	UpdateProductImage(productID string, imageURL string) error
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) CreateCategory(category *entity.Category) error {
	return s.repo.CreateCategory(category)
}

func (s *productService) GetAllCategories() ([]entity.Category, error) {
	return s.repo.GetAllCategories()
}

func (s *productService) GetCategoryByID(id uint) (*entity.Category, error) {
	return s.repo.GetCategoryByID(id)
}

func (s *productService) UpdateCategory(category *entity.Category) error {
	return s.repo.UpdateCategory(category)
}

func (s *productService) DeleteCategory(id uint) error {
	return s.repo.DeleteCategory(id)
}

func (s *productService) CreateProduct(product *entity.Product) error {
	return s.repo.CreateProduct(product)
}

func (s *productService) GetAllProducts() ([]entity.Product, error) {
	return s.repo.GetAllProducts()
}

func (s *productService) GetProductByID(id uint) (*entity.Product, error) {
	return s.repo.GetProductByID(id)
}

func (s *productService) UpdateProduct(product *entity.Product) error {
	return s.repo.UpdateProduct(product)
}

func (s *productService) DeleteProduct(id uint) error {
	return s.repo.DeleteProduct(id)
}

func (s *productService) UpdateProductImage(productID string, imageURL string) error {
	return s.repo.UpdateImage(productID, imageURL)
}
