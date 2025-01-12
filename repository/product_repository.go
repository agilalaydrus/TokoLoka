package repository

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"main.go/entity"
	"main.go/middleware"
)

type ProductRepository interface {
	// Category methods
	CreateCategory(category *entity.Category) error
	GetAllCategories() ([]entity.Category, error)
	GetCategoryByID(id uint) (*entity.Category, error)
	UpdateCategory(category *entity.Category) error
	DeleteCategory(id uint) error

	// Product methods
	CreateProduct(product *entity.Product) error
	GetAllProducts() ([]entity.Product, error)
	GetProductByID(id uint) (*entity.Product, error)
	UpdateProduct(product *entity.Product) error
	DeleteProduct(id uint) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

// Category methods
func (r *productRepository) CreateCategory(category *entity.Category) error {
	middleware.Logger.Info("Repository: Creating category", zap.Any("category", category))
	if category.Name == "" {
		middleware.Logger.Warn("Repository: Category name cannot be empty")
		return errors.New("category name cannot be empty")
	}
	if err := r.db.Create(category).Error; err != nil {
		middleware.Logger.Error("Repository: Error creating category", zap.Error(err))
		return err
	}
	middleware.Logger.Info("Repository: Category created successfully", zap.Uint("category_id", category.ID))
	return nil
}

func (r *productRepository) GetAllCategories() ([]entity.Category, error) {
	middleware.Logger.Info("Repository: Fetching all categories")
	var categories []entity.Category
	if err := r.db.Find(&categories).Error; err != nil {
		middleware.Logger.Error("Repository: Error fetching categories", zap.Error(err))
		return nil, err
	}
	middleware.Logger.Info("Repository: Categories fetched successfully", zap.Int("count", len(categories)))
	return categories, nil
}

func (r *productRepository) GetCategoryByID(id uint) (*entity.Category, error) {
	middleware.Logger.Info("Repository: Fetching category by ID", zap.Uint("category_id", id))
	var category entity.Category
	if err := r.db.First(&category, id).Error; err != nil {
		middleware.Logger.Warn("Repository: Category not found", zap.Error(err))
		return nil, errors.New("category not found")
	}
	middleware.Logger.Info("Repository: Category fetched successfully", zap.Any("category", category))
	return &category, nil
}

func (r *productRepository) UpdateCategory(category *entity.Category) error {
	middleware.Logger.Info("Repository: Updating category", zap.Uint("category_id", category.ID))
	if err := r.db.Save(category).Error; err != nil {
		middleware.Logger.Error("Repository: Error updating category", zap.Error(err))
		return err
	}
	middleware.Logger.Info("Repository: Category updated successfully", zap.Uint("category_id", category.ID))
	return nil
}

func (r *productRepository) DeleteCategory(id uint) error {
	middleware.Logger.Info("Repository: Deleting category", zap.Uint("category_id", id))
	var category entity.Category
	if err := r.db.First(&category, id).Error; err != nil {
		middleware.Logger.Warn("Repository: Category not found for deletion", zap.Error(err))
		return errors.New("category not found")
	}
	if err := r.db.Delete(&category).Error; err != nil {
		middleware.Logger.Error("Repository: Error deleting category", zap.Error(err))
		return err
	}
	middleware.Logger.Info("Repository: Category deleted successfully", zap.Uint("category_id", id))
	return nil
}

// Product methods
func (r *productRepository) CreateProduct(product *entity.Product) error {
	middleware.Logger.Info("Repository: Creating product", zap.Any("product", product))
	if product.Name == "" || product.Price <= 0 || product.Stock < 0 {
		middleware.Logger.Warn("Repository: Invalid product data", zap.Any("product", product))
		return errors.New("invalid product data")
	}
	if err := r.db.Create(product).Error; err != nil {
		middleware.Logger.Error("Repository: Error creating product", zap.Error(err))
		return err
	}
	middleware.Logger.Info("Repository: Product created successfully", zap.Uint("product_id", product.ID))
	return nil
}

func (r *productRepository) GetAllProducts() ([]entity.Product, error) {
	middleware.Logger.Info("Repository: Fetching all products")
	var products []entity.Product
	if err := r.db.Preload("Category").Find(&products).Error; err != nil {
		middleware.Logger.Error("Repository: Error fetching products", zap.Error(err))
		return nil, err
	}
	middleware.Logger.Info("Repository: Products fetched successfully", zap.Int("count", len(products)))
	return products, nil
}

func (r *productRepository) GetProductByID(id uint) (*entity.Product, error) {
	middleware.Logger.Info("Repository: Fetching product by ID", zap.Uint("product_id", id))
	var product entity.Product
	if err := r.db.Preload("Category").First(&product, id).Error; err != nil {
		middleware.Logger.Warn("Repository: Product not found", zap.Error(err))
		return nil, errors.New("product not found")
	}
	middleware.Logger.Info("Repository: Product fetched successfully", zap.Any("product", product))
	return &product, nil
}

func (r *productRepository) UpdateProduct(product *entity.Product) error {
	middleware.Logger.Info("Repository: Updating product", zap.Uint("product_id", product.ID))
	if err := r.db.Save(product).Error; err != nil {
		middleware.Logger.Error("Repository: Error updating product", zap.Error(err))
		return err
	}
	middleware.Logger.Info("Repository: Product updated successfully", zap.Uint("product_id", product.ID))
	return nil
}

func (r *productRepository) DeleteProduct(id uint) error {
	middleware.Logger.Info("Repository: Deleting product", zap.Uint("product_id", id))
	var product entity.Product
	if err := r.db.First(&product, id).Error; err != nil {
		middleware.Logger.Warn("Repository: Product not found for deletion", zap.Error(err))
		return errors.New("product not found")
	}
	if err := r.db.Delete(&product).Error; err != nil {
		middleware.Logger.Error("Repository: Error deleting product", zap.Error(err))
		return err
	}
	middleware.Logger.Info("Repository: Product deleted successfully", zap.Uint("product_id", id))
	return nil
}
