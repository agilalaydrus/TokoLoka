package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"main.go/entity"
	"main.go/middleware"
	"main.go/service"
	"net/http"
	"strconv"
)

type ProductController struct {
	service service.ProductService
}

func NewProductController(service service.ProductService) *ProductController {
	return &ProductController{service: service}
}

// CreateCategory - Create a new category
func (pc *ProductController) CreateCategory(c *gin.Context) {
	middleware.Logger.Info("Controller: CreateCategory called")

	var category entity.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		middleware.Logger.Error("Invalid input", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if category.Name == "" {
		middleware.Logger.Warn("Category name cannot be empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category name cannot be empty"})
		return
	}

	if err := pc.service.CreateCategory(&category); err != nil {
		middleware.Logger.Error("Failed to create category", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	middleware.Logger.Info("Category created successfully", zap.String("name", category.Name))
	c.JSON(http.StatusCreated, gin.H{"message": "Category created successfully"})
}

// GetAllCategories - Retrieve all categories
func (pc *ProductController) GetAllCategories(c *gin.Context) {
	middleware.Logger.Info("Controller: GetAllCategories called")

	categories, err := pc.service.GetAllCategories()
	if err != nil {
		middleware.Logger.Error("Failed to fetch categories", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	middleware.Logger.Info("Categories fetched successfully", zap.Int("count", len(categories)))
	c.JSON(http.StatusOK, gin.H{"message": "Categories fetched successfully", "data": categories})
}

// GetCategoryByID - Retrieve a category by ID
func (pc *ProductController) GetCategoryByID(c *gin.Context) {
	middleware.Logger.Info("Controller: GetCategoryByID called")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		middleware.Logger.Error("Invalid category ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	category, err := pc.service.GetCategoryByID(uint(id))
	if err != nil {
		middleware.Logger.Error("Category not found", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	middleware.Logger.Info("Category fetched successfully", zap.String("name", category.Name))
	c.JSON(http.StatusOK, gin.H{"message": "Category fetched successfully", "data": category})
}

// UpdateCategory - Update an existing category
func (pc *ProductController) UpdateCategory(c *gin.Context) {
	middleware.Logger.Info("Controller: UpdateCategory called")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		middleware.Logger.Error("Invalid category ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var category entity.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		middleware.Logger.Error("Invalid input", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if category.Name == "" {
		middleware.Logger.Warn("Category name cannot be empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category name cannot be empty"})
		return
	}

	category.ID = uint(id)
	if err := pc.service.UpdateCategory(&category); err != nil {
		middleware.Logger.Error("Failed to update category", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	middleware.Logger.Info("Category updated successfully", zap.String("name", category.Name))
	c.JSON(http.StatusOK, gin.H{"message": "Category updated successfully"})
}

// DeleteCategory - Delete a category by ID
func (pc *ProductController) DeleteCategory(c *gin.Context) {
	middleware.Logger.Info("Controller: DeleteCategory called")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		middleware.Logger.Error("Invalid category ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	if err := pc.service.DeleteCategory(uint(id)); err != nil {
		middleware.Logger.Error("Failed to delete category", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	middleware.Logger.Info("Category deleted successfully", zap.Int("id", id))
	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}

// CreateProduct - Create a new product
func (pc *ProductController) CreateProduct(c *gin.Context) {
	middleware.Logger.Info("Controller: CreateProduct called")

	var product entity.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		middleware.Logger.Error("Invalid input", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if product.Name == "" || product.Price <= 0 || product.Stock < 0 {
		middleware.Logger.Warn("Invalid product data", zap.String("name", product.Name))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product data"})
		return
	}

	if err := pc.service.CreateProduct(&product); err != nil {
		middleware.Logger.Error("Failed to create product", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	middleware.Logger.Info("Product created successfully", zap.String("name", product.Name))
	c.JSON(http.StatusCreated, gin.H{"message": "Product created successfully"})
}

// GetAllProducts - Retrieve all products
func (pc *ProductController) GetAllProducts(c *gin.Context) {
	middleware.Logger.Info("Controller: GetAllProducts called")

	products, err := pc.service.GetAllProducts()
	if err != nil {
		middleware.Logger.Error("Failed to fetch products", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	middleware.Logger.Info("Products fetched successfully", zap.Int("count", len(products)))
	c.JSON(http.StatusOK, gin.H{"message": "Products fetched successfully", "data": products})
}

// GetProductByID - Retrieve a product by ID
func (pc *ProductController) GetProductByID(c *gin.Context) {
	middleware.Logger.Info("Controller: GetProductByID called")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		middleware.Logger.Error("Invalid product ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	product, err := pc.service.GetProductByID(uint(id))
	if err != nil {
		middleware.Logger.Error("Product not found", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	middleware.Logger.Info("Product fetched successfully", zap.String("name", product.Name))
	c.JSON(http.StatusOK, gin.H{"message": "Product fetched successfully", "data": product})
}

// UpdateProduct - Update an existing product
func (pc *ProductController) UpdateProduct(c *gin.Context) {
	middleware.Logger.Info("Controller: UpdateProduct called")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		middleware.Logger.Error("Invalid product ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product entity.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		middleware.Logger.Error("Invalid input", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	product.ID = uint(id)
	if err := pc.service.UpdateProduct(&product); err != nil {
		middleware.Logger.Error("Failed to update product", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	middleware.Logger.Info("Product updated successfully", zap.String("name", product.Name))
	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}

// DeleteProduct - Delete a product by ID
func (pc *ProductController) DeleteProduct(c *gin.Context) {
	middleware.Logger.Info("Controller: DeleteProduct called")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		middleware.Logger.Error("Invalid product ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	if err := pc.service.DeleteProduct(uint(id)); err != nil {
		middleware.Logger.Error("Failed to delete product", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	middleware.Logger.Info("Product deleted successfully", zap.Int("id", id))
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
