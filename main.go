package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"main.go/config"
	"main.go/controller"
	"main.go/middleware"
	"main.go/repository"
	"main.go/service"
	"net/http"
)

func main() {
	// Inisialisasi koneksi ke database
	if err := config.InitDB(); err != nil {
		log.Fatalf("Gagal menginisialisasi database: %v", err)
	}
	log.Println("Database berhasil diinisialisasi")

	// Inisialisasi Repository
	userRepo := repository.NewUserRepository(config.DB)
	productRepo := repository.NewProductRepository(config.DB)

	// Inisialisasi Service
	userService := service.NewUserService(userRepo)
	productService := service.NewProductService(productRepo)

	// Inisialisasi Controller
	userController := controller.NewUserController(userService)
	productController := controller.NewProductController(productService)

	// Membuat router Gin
	r := gin.Default()
	r.SetTrustedProxies(nil) // Tidak mempercayai proxy mana pun

	// Routes untuk autentikasi
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", userController.RegisterUser) // Endpoint untuk registrasi
		authRoutes.POST("/login", userController.LoginUser)       // Endpoint untuk login
	}
	log.Println("Routes untuk autentikasi berhasil didaftarkan")

	// Routes yang dilindungi oleh JWT
	protectedRoutes := r.Group("/api")
	protectedRoutes.Use(middleware.AuthorizeJWT) // Middleware JWT diterapkan
	{
		// Rute untuk Administrator (akses penuh)
		adminRoutes := protectedRoutes.Group("/")
		adminRoutes.Use(middleware.RoleBasedAccessControl("administrator"))
		{
			// CRUD Categories
			adminRoutes.POST("/categories", productController.CreateCategory)
			adminRoutes.PUT("/categories/:id", productController.UpdateCategory)
			adminRoutes.DELETE("/categories/:id", productController.DeleteCategory)

			// CRUD Products
			adminRoutes.POST("/products", productController.CreateProduct)
			adminRoutes.PUT("/products/:id", productController.UpdateProduct)
			adminRoutes.DELETE("/products/:id", productController.DeleteProduct)
		}

		// Rute untuk User dan Administrator
		userRoutes := protectedRoutes.Group("/")
		userRoutes.Use(middleware.RoleBasedAccessControl("ANY")) // Role apapun
		{
			// Categories dan Products (GET saja)
			userRoutes.GET("/categories", productController.GetAllCategories)
			userRoutes.GET("/categories/:id", productController.GetCategoryByID)

			userRoutes.GET("/products", productController.GetAllProducts)
			userRoutes.GET("/products/:id", productController.GetProductByID)

			// Transaksi dan Laporan (tambahkan sesuai kebutuhan)
			// userRoutes.POST("/transactions", transactionController.CreateTransaction)
			// userRoutes.GET("/reports", reportController.GetReports)
		}
	}
	log.Println("Routes yang dilindungi JWT berhasil didaftarkan")

	// Endpoint debugging untuk mencetak semua rute
	r.GET("/debug/routes", func(c *gin.Context) {
		routes := r.Routes()
		for _, route := range routes {
			log.Printf("Route: %s %s", route.Method, route.Path)
		}
		c.JSON(http.StatusOK, routes)
	})

	// Menjalankan server di port 8080
	log.Println("Server dijalankan pada port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server gagal dijalankan: %v", err)
	}
}
