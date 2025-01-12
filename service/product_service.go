package service

import (
	"errors"
	"log"
	"main.go/entity"
	"main.go/repository"
)

type ProductService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

// Category methods
func (s *ProductService) CreateCategory(category *entity.Category) error {
	log.Printf("Service: Creating category with data: %+v", category)

	if category.Name == "" {
		log.Println("Service: Category name cannot be empty")
		return errors.New("category name cannot be empty")
	}

	if err := s.repo.CreateCategory(category); err != nil {
		log.Printf("Service: Failed to create category: %v", err)
		return err
	}

	log.Println("Service: Category created successfully")
	return nil
}

func (s *ProductService) GetAllCategories() ([]entity.Category, error) {
	log.Println("Service: Fetching all categories")

	categories, err := s.repo.GetAllCategories()
	if err != nil {
		log.Printf("Service: Error fetching categories: %v", err)
		return nil, err
	}

	log.Printf("Service: Fetched categories: %+v", categories)
	return categories, nil
}

func (s *ProductService) GetCategoryByID(id uint) (*entity.Category, error) {
	log.Printf("Service: Fetching category with ID: %d", id)

	category, err := s.repo.GetCategoryByID(id)
	if err != nil {
		log.Printf("Service: Category not found: %v", err)
		return nil, errors.New("category not found")
	}

	log.Printf("Service: Fetched category: %+v", category)
	return category, nil
}

func (s *ProductService) UpdateCategory(category *entity.Category) error {
	log.Printf("Service: Updating category with ID: %d", category.ID)

	existingCategory, err := s.repo.GetCategoryByID(category.ID)
	if err != nil {
		log.Printf("Service: Category not found for update: %v", err)
		return errors.New("category not found")
	}

	existingCategory.Name = category.Name
	existingCategory.Description = category.Description

	if err := s.repo.UpdateCategory(existingCategory); err != nil {
		log.Printf("Service: Failed to update category: %v", err)
		return errors.New("failed to update category")
	}

	log.Println("Service: Category updated successfully")
	return nil
}

func (s *ProductService) DeleteCategory(id uint) error {
	log.Printf("Service: Deleting category with ID: %d", id)

	_, err := s.repo.GetCategoryByID(id)
	if err != nil {
		log.Printf("Service: Category not found for deletion: %v", err)
		return errors.New("category not found")
	}

	if err := s.repo.DeleteCategory(id); err != nil {
		log.Printf("Service: Failed to delete category: %v", err)
		return errors.New("failed to delete category")
	}

	log.Println("Service: Category deleted successfully")
	return nil
}

// Product methods
func (s *ProductService) CreateProduct(product *entity.Product) error {
	log.Printf("Service: Creating product with data: %+v", product)

	if product.Name == "" || product.Price <= 0 || product.Stock < 0 {
		log.Println("Service: Invalid product data")
		return errors.New("invalid product data")
	}

	// Validasi kategori
	_, err := s.repo.GetCategoryByID(product.CategoryID)
	if err != nil {
		log.Printf("Service: Category not found for product: %v", err)
		return errors.New("category not found")
	}

	if err := s.repo.CreateProduct(product); err != nil {
		log.Printf("Service: Failed to create product: %v", err)
		return err
	}

	log.Println("Service: Product created successfully")
	return nil
}

func (s *ProductService) GetAllProducts() ([]entity.Product, error) {
	log.Println("Service: Fetching all products")

	products, err := s.repo.GetAllProducts()
	if err != nil {
		log.Printf("Service: Error fetching products: %v", err)
		return nil, err
	}

	log.Printf("Service: Fetched products: %+v", products)
	return products, nil
}

func (s *ProductService) GetProductByID(id uint) (*entity.Product, error) {
	log.Printf("Service: Fetching product with ID: %d", id)

	product, err := s.repo.GetProductByID(id)
	if err != nil {
		log.Printf("Service: Product not found: %v", err)
		return nil, errors.New("product not found")
	}

	log.Printf("Service: Fetched product: %+v", product)
	return product, nil
}

func (s *ProductService) UpdateProduct(product *entity.Product) error {
	log.Printf("Service: Updating product with ID: %d", product.ID)

	existingProduct, err := s.repo.GetProductByID(product.ID)
	if err != nil {
		log.Printf("Service: Product not found for update: %v", err)
		return errors.New("product not found")
	}

	existingProduct.Name = product.Name
	existingProduct.Description = product.Description
	existingProduct.Price = product.Price
	existingProduct.Stock = product.Stock
	existingProduct.CategoryID = product.CategoryID

	if err := s.repo.UpdateProduct(existingProduct); err != nil {
		log.Printf("Service: Failed to update product: %v", err)
		return errors.New("failed to update product")
	}

	log.Println("Service: Product updated successfully")
	return nil
}

func (s *ProductService) DeleteProduct(id uint) error {
	log.Printf("Service: Deleting product with ID: %d", id)

	_, err := s.repo.GetProductByID(id)
	if err != nil {
		log.Printf("Service: Product not found for deletion: %v", err)
		return errors.New("product not found")
	}

	if err := s.repo.DeleteProduct(id); err != nil {
		log.Printf("Service: Failed to delete product: %v", err)
		return errors.New("failed to delete product")
	}

	log.Println("Service: Product deleted successfully")
	return nil
}
