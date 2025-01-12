package service

import (
	"go.uber.org/zap"
	"main.go/entity"
	"main.go/middleware"
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
	middleware.Logger.Info("Service: Creating category", zap.Any("category", category))

	if category.Name == "" {
		middleware.Logger.Warn("Service: Category name cannot be empty")
		return middleware.NewAppError(400, "Category name cannot be empty", nil)
	}

	if err := s.repo.CreateCategory(category); err != nil {
		middleware.Logger.Error("Service: Failed to create category", zap.Error(err))
		return middleware.NewAppError(500, "Failed to create category", err)
	}

	middleware.Logger.Info("Service: Category created successfully", zap.Uint("category_id", category.ID))
	return nil
}

func (s *ProductService) GetAllCategories() ([]entity.Category, error) {
	middleware.Logger.Info("Service: Fetching all categories")

	categories, err := s.repo.GetAllCategories()
	if err != nil {
		middleware.Logger.Error("Service: Error fetching categories", zap.Error(err))
		return nil, middleware.NewAppError(500, "Error fetching categories", err)
	}

	middleware.Logger.Info("Service: Fetched categories", zap.Any("categories", categories))
	return categories, nil
}

func (s *ProductService) GetCategoryByID(id uint) (*entity.Category, error) {
	middleware.Logger.Info("Service: Fetching category", zap.Uint("category_id", id))

	category, err := s.repo.GetCategoryByID(id)
	if err != nil {
		middleware.Logger.Warn("Service: Category not found", zap.Uint("category_id", id))
		return nil, middleware.NewAppError(404, "Category not found", err)
	}

	middleware.Logger.Info("Service: Fetched category", zap.Any("category", category))
	return category, nil
}

func (s *ProductService) UpdateCategory(category *entity.Category) error {
	middleware.Logger.Info("Service: Updating category", zap.Uint("category_id", category.ID))

	existingCategory, err := s.repo.GetCategoryByID(category.ID)
	if err != nil {
		middleware.Logger.Warn("Service: Category not found for update", zap.Uint("category_id", category.ID))
		return middleware.NewAppError(404, "Category not found", err)
	}

	existingCategory.Name = category.Name
	existingCategory.Description = category.Description

	if err := s.repo.UpdateCategory(existingCategory); err != nil {
		middleware.Logger.Error("Service: Failed to update category", zap.Error(err))
		return middleware.NewAppError(500, "Failed to update category", err)
	}

	middleware.Logger.Info("Service: Category updated successfully", zap.Uint("category_id", category.ID))
	return nil
}

func (s *ProductService) DeleteCategory(id uint) error {
	middleware.Logger.Info("Service: Deleting category", zap.Uint("category_id", id))

	_, err := s.repo.GetCategoryByID(id)
	if err != nil {
		middleware.Logger.Warn("Service: Category not found for deletion", zap.Uint("category_id", id))
		return middleware.NewAppError(404, "Category not found", err)
	}

	if err := s.repo.DeleteCategory(id); err != nil {
		middleware.Logger.Error("Service: Failed to delete category", zap.Error(err))
		return middleware.NewAppError(500, "Failed to delete category", err)
	}

	middleware.Logger.Info("Service: Category deleted successfully", zap.Uint("category_id", id))
	return nil
}

// Product methods
func (s *ProductService) CreateProduct(product *entity.Product) error {
	middleware.Logger.Info("Service: Creating product", zap.Any("product", product))

	if product.Name == "" || product.Price <= 0 || product.Stock < 0 {
		middleware.Logger.Warn("Service: Invalid product data", zap.Any("product", product))
		return middleware.NewAppError(400, "Invalid product data", nil)
	}

	// Validasi kategori
	_, err := s.repo.GetCategoryByID(product.CategoryID)
	if err != nil {
		middleware.Logger.Warn("Service: Category not found for product", zap.Uint("category_id", product.CategoryID))
		return middleware.NewAppError(404, "Category not found", err)
	}

	if err := s.repo.CreateProduct(product); err != nil {
		middleware.Logger.Error("Service: Failed to create product", zap.Error(err))
		return middleware.NewAppError(500, "Failed to create product", err)
	}

	middleware.Logger.Info("Service: Product created successfully", zap.Uint("product_id", product.ID))
	return nil
}

func (s *ProductService) GetAllProducts() ([]entity.Product, error) {
	middleware.Logger.Info("Service: Fetching all products")

	products, err := s.repo.GetAllProducts()
	if err != nil {
		middleware.Logger.Error("Service: Error fetching products", zap.Error(err))
		return nil, middleware.NewAppError(500, "Error fetching products", err)
	}

	middleware.Logger.Info("Service: Fetched products", zap.Any("products", products))
	return products, nil
}

func (s *ProductService) GetProductByID(id uint) (*entity.Product, error) {
	middleware.Logger.Info("Service: Fetching product", zap.Uint("product_id", id))

	product, err := s.repo.GetProductByID(id)
	if err != nil {
		middleware.Logger.Warn("Service: Product not found", zap.Uint("product_id", id))
		return nil, middleware.NewAppError(404, "Product not found", err)
	}

	middleware.Logger.Info("Service: Fetched product", zap.Any("product", product))
	return product, nil
}

func (s *ProductService) UpdateProduct(product *entity.Product) error {
	middleware.Logger.Info("Service: Updating product", zap.Uint("product_id", product.ID))

	existingProduct, err := s.repo.GetProductByID(product.ID)
	if err != nil {
		middleware.Logger.Warn("Service: Product not found for update", zap.Uint("product_id", product.ID))
		return middleware.NewAppError(404, "Product not found", err)
	}

	existingProduct.Name = product.Name
	existingProduct.Description = product.Description
	existingProduct.Price = product.Price
	existingProduct.Stock = product.Stock
	existingProduct.CategoryID = product.CategoryID

	if err := s.repo.UpdateProduct(existingProduct); err != nil {
		middleware.Logger.Error("Service: Failed to update product", zap.Error(err))
		return middleware.NewAppError(500, "Failed to update product", err)
	}

	middleware.Logger.Info("Service: Product updated successfully", zap.Uint("product_id", product.ID))
	return nil
}

func (s *ProductService) DeleteProduct(id uint) error {
	middleware.Logger.Info("Service: Deleting product", zap.Uint("product_id", id))

	_, err := s.repo.GetProductByID(id)
	if err != nil {
		middleware.Logger.Warn("Service: Product not found for deletion", zap.Uint("product_id", id))
		return middleware.NewAppError(404, "Product not found", err)
	}

	if err := s.repo.DeleteProduct(id); err != nil {
		middleware.Logger.Error("Service: Failed to delete product", zap.Error(err))
		return middleware.NewAppError(500, "Failed to delete product", err)
	}

	middleware.Logger.Info("Service: Product deleted successfully", zap.Uint("product_id", id))
	return nil
}
