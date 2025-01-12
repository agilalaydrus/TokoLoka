package repository

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"main.go/entity"
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
	log.Printf("Repository: Creating category with data: %+v", category)
	if category.Name == "" {
		log.Println("Repository: Category name cannot be empty")
		return errors.New("category name cannot be empty")
	}
	if err := r.db.Create(category).Error; err != nil {
		log.Printf("Repository: Error creating category: %v", err)
		return err
	}
	log.Println("Repository: Category created successfully")
	return nil
}

func (r *productRepository) GetAllCategories() ([]entity.Category, error) {
	log.Println("Repository: Fetching all categories from database")
	var categories []entity.Category
	if err := r.db.Find(&categories).Error; err != nil {
		log.Printf("Repository: Error fetching categories: %v", err)
		return nil, err
	}
	log.Printf("Repository: Fetched categories: %+v", categories)
	return categories, nil
}

func (r *productRepository) GetCategoryByID(id uint) (*entity.Category, error) {
	log.Printf("Repository: Fetching category with ID: %d", id)
	var category entity.Category
	if err := r.db.First(&category, id).Error; err != nil {
		log.Printf("Repository: Error fetching category: %v", err)
		return nil, err
	}
	log.Printf("Repository: Fetched category: %+v", category)
	return &category, nil
}

func (r *productRepository) UpdateCategory(category *entity.Category) error {
	log.Printf("Repository: Updating category with ID: %d", category.ID)
	if err := r.db.Save(category).Error; err != nil {
		log.Printf("Repository: Error updating category: %v", err)
		return err
	}
	log.Println("Repository: Category updated successfully")
	return nil
}

func (r *productRepository) DeleteCategory(id uint) error {
	log.Printf("Repository: Deleting category with ID: %d", id)
	var category entity.Category
	if err := r.db.First(&category, id).Error; err != nil {
		log.Printf("Repository: Category not found: %v", err)
		return errors.New("category not found")
	}
	if err := r.db.Delete(&category).Error; err != nil {
		log.Printf("Repository: Error deleting category: %v", err)
		return err
	}
	log.Println("Repository: Category deleted successfully")
	return nil
}

// Product methods
func (r *productRepository) CreateProduct(product *entity.Product) error {
	log.Printf("Repository: Creating product with data: %+v", product)
	if product.Name == "" || product.Price <= 0 || product.Stock < 0 {
		log.Println("Repository: Invalid product data")
		return errors.New("invalid product data")
	}
	if err := r.db.Create(product).Error; err != nil {
		log.Printf("Repository: Error creating product: %v", err)
		return err
	}
	log.Println("Repository: Product created successfully")
	return nil
}

func (r *productRepository) GetAllProducts() ([]entity.Product, error) {
	log.Println("Repository: Fetching all products from database")
	var products []entity.Product
	if err := r.db.Preload("Category").Find(&products).Error; err != nil {
		log.Printf("Repository: Error fetching products: %v", err)
		return nil, err
	}
	log.Printf("Repository: Fetched products: %+v", products)
	return products, nil
}

func (r *productRepository) GetProductByID(id uint) (*entity.Product, error) {
	log.Printf("Repository: Fetching product with ID: %d", id)
	var product entity.Product
	if err := r.db.Preload("Category").First(&product, id).Error; err != nil {
		log.Printf("Repository: Error fetching product: %v", err)
		return nil, err
	}
	log.Printf("Repository: Fetched product: %+v", product)
	return &product, nil
}

func (r *productRepository) UpdateProduct(product *entity.Product) error {
	log.Printf("Repository: Updating product with ID: %d", product.ID)
	if err := r.db.Save(product).Error; err != nil {
		log.Printf("Repository: Error updating product: %v", err)
		return err
	}
	log.Println("Repository: Product updated successfully")
	return nil
}

func (r *productRepository) DeleteProduct(id uint) error {
	log.Printf("Repository: Deleting product with ID: %d", id)
	var product entity.Product
	if err := r.db.First(&product, id).Error; err != nil {
		log.Printf("Repository: Product not found: %v", err)
		return errors.New("product not found")
	}
	if err := r.db.Delete(&product).Error; err != nil {
		log.Printf("Repository: Error deleting product: %v", err)
		return err
	}
	log.Println("Repository: Product deleted successfully")
	return nil
}
