package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"main.go/entity"
	"main.go/service"
)

type ProductController struct {
	service *service.ProductService
}

func NewProductController(service *service.ProductService) *ProductController {
	return &ProductController{service: service}
}

// CreateCategory - Create a new category
func (pc *ProductController) CreateCategory(c *gin.Context) {
	log.Println("Controller: CreateCategory called")

	var category entity.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		log.Printf("Controller: Invalid input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Additional validation
	if category.Name == "" {
		log.Println("Controller: Category name cannot be empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category name cannot be empty"})
		return
	}

	if err := pc.service.CreateCategory(&category); err != nil {
		log.Printf("Controller: Failed to create category: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Println("Controller: Category created successfully")
	c.JSON(http.StatusCreated, gin.H{"message": "Category created successfully"})
}

// GetAllCategories - Retrieve all categories
func (pc *ProductController) GetAllCategories(c *gin.Context) {
	log.Println("Controller: GetAllCategories called")

	categories, err := pc.service.GetAllCategories()
	if err != nil {
		log.Printf("Controller: Failed to fetch categories: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	log.Printf("Controller: Fetched categories: %+v", categories)
	c.JSON(http.StatusOK, gin.H{"message": "Categories fetched successfully", "data": categories})
}

// GetCategoryByID - Retrieve a category by ID
func (pc *ProductController) GetCategoryByID(c *gin.Context) {
	log.Println("Controller: GetCategoryByID called")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		log.Println("Controller: Invalid category ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	category, err := pc.service.GetCategoryByID(uint(id))
	if err != nil {
		log.Printf("Controller: Category not found: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Controller: Fetched category: %+v", category)
	c.JSON(http.StatusOK, gin.H{"message": "Category fetched successfully", "data": category})
}

// UpdateCategory - Update an existing category
func (pc *ProductController) UpdateCategory(c *gin.Context) {
	log.Println("Controller: UpdateCategory called")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		log.Println("Controller: Invalid category ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var category entity.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		log.Printf("Controller: Invalid input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Additional validation
	if category.Name == "" {
		log.Println("Controller: Category name cannot be empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category name cannot be empty"})
		return
	}

	category.ID = uint(id)
	if err := pc.service.UpdateCategory(&category); err != nil {
		log.Printf("Controller: Failed to update category: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Println("Controller: Category updated successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Category updated successfully"})
}

// DeleteCategory - Delete a category by ID
func (pc *ProductController) DeleteCategory(c *gin.Context) {
	log.Println("Controller: DeleteCategory called")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		log.Println("Controller: Invalid category ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	if err := pc.service.DeleteCategory(uint(id)); err != nil {
		log.Printf("Controller: Failed to delete category: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Println("Controller: Category deleted successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}

// CreateProduct - Create a new product
func (pc *ProductController) CreateProduct(c *gin.Context) {
	log.Println("Controller: CreateProduct called")

	var product entity.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		log.Printf("Controller: Invalid input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if product.Name == "" || product.Price <= 0 || product.Stock < 0 {
		log.Println("Controller: Invalid product data")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product data"})
		return
	}

	if err := pc.service.CreateProduct(&product); err != nil {
		log.Printf("Controller: Failed to create product: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println("Controller: Product created successfully")
	c.JSON(http.StatusCreated, gin.H{"message": "Product created successfully"})
}

// GetAllProducts - Retrieve all products
func (pc *ProductController) GetAllProducts(c *gin.Context) {
	log.Println("Controller: GetAllProducts called")

	products, err := pc.service.GetAllProducts()
	if err != nil {
		log.Printf("Controller: Failed to fetch products: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	log.Printf("Controller: Fetched products: %+v", products)
	c.JSON(http.StatusOK, gin.H{"message": "Products fetched successfully", "data": products})
}

// GetProductByID - Retrieve a product by ID
func (pc *ProductController) GetProductByID(c *gin.Context) {
	log.Println("Controller: GetProductByID called")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		log.Println("Controller: Invalid product ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	product, err := pc.service.GetProductByID(uint(id))
	if err != nil {
		log.Printf("Controller: Product not found: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Controller: Fetched product: %+v", product)
	c.JSON(http.StatusOK, gin.H{"message": "Product fetched successfully", "data": product})
}

// UpdateProduct - Update an existing product
func (pc *ProductController) UpdateProduct(c *gin.Context) {
	log.Println("Controller: UpdateProduct called")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		log.Println("Controller: Invalid product ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product entity.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		log.Printf("Controller: Invalid input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	product.ID = uint(id)

	if err := pc.service.UpdateProduct(&product); err != nil {
		log.Printf("Controller: Failed to update product: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println("Controller: Product updated successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}

// DeleteProduct - Delete a product by ID
func (pc *ProductController) DeleteProduct(c *gin.Context) {
	log.Println("Controller: DeleteProduct called")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		log.Println("Controller: Invalid product ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	if err := pc.service.DeleteProduct(uint(id)); err != nil {
		log.Printf("Controller: Failed to delete product: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println("Controller: Product deleted successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
