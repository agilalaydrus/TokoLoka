package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"main.go/config"
	"main.go/controller"
	"main.go/middleware"
	"main.go/repository"
	"main.go/service"
	"net/http"
	"time"
)

func main() {
	// Inisialisasi logger Zap
	middleware.InitLogger()
	defer middleware.Logger.Sync()
	middleware.Logger.Info("Logger berhasil diinisialisasi")

	// Inisialisasi koneksi ke database
	if err := config.InitDB(); err != nil {
		middleware.Logger.Fatal("Gagal menginisialisasi database", zap.Error(err))
	}
	middleware.Logger.Info("Database berhasil diinisialisasi")

	// Inisialisasi Repository
	userRepo := repository.NewUserRepository(config.DB)
	productRepo := repository.NewProductRepository(config.DB)
	transactionRepo := repository.NewTransactionsRepository(config.DB)
	activityLogRepo := repository.NewActivityLogRepository(config.DB)
	reportRepo := repository.NewReportRepository(config.DB)
	tokenRepo := repository.NewTokenRepository(config.DB)

	// Inisialisasi Service
	userService := service.NewUserService(userRepo, tokenRepo)
	productService := service.NewProductService(productRepo)
	activityLogService := service.NewActivityLogService(activityLogRepo)
	transactionService := service.NewTransactionsService(transactionRepo, productRepo, activityLogService)
	reportService := service.NewReportService(reportRepo) // Pastikan ini digunakan

	// Inisialisasi Controller
	userController := controller.NewUserController(userService)
	productController := controller.NewProductController(productService)
	transactionController := controller.NewTransactionsController(transactionService)
	reportController := controller.NewReportController(reportService) // Pastikan ini digunakan
	callbackController := controller.NewCallbackController(transactionService)

	// Membuat router Gin
	r := gin.Default()
	r.SetTrustedProxies(nil)

	// Tambahkan middleware logging request
	r.Use(middleware.RequestLogger())

	// Tambahkan middleware global untuk Error Handling
	r.Use(middleware.ErrorHandler())

	// Routes untuk Callback Simulasi
	r.POST("/callback/transaction-status", callbackController.CallbackTransactionStatus)

	// Routes untuk Autentikasi
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", middleware.LimitRequest(3, 2*time.Minute), userController.RegisterUser)
		authRoutes.POST("/login", middleware.LimitRequest(5, 1*time.Minute), userController.LoginUser)
		authRoutes.POST("/refresh", userController.RefreshToken)
		authRoutes.POST("/logout", middleware.AuthorizeJWT, userController.Logout)

	}

	middleware.Logger.Info("Routes untuk autentikasi berhasil didaftarkan")

	// Routes yang dilindungi oleh JWT
	protectedRoutes := r.Group("/api")
	protectedRoutes.Use(middleware.AuthorizeJWT)
	{
		// Rute untuk Administrator
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

			// Transactions Management
			adminRoutes.DELETE("/transactions/:id", transactionController.DeleteTransaction)
			adminRoutes.PUT("/transactions/:id/status", transactionController.UpdateTransactionStatus)
			adminRoutes.GET("/transactions", transactionController.GetAllTransactions)
		}

		// Rute untuk User dan Administrator
		userRoutes := protectedRoutes.Group("/")
		userRoutes.Use(middleware.RoleBasedAccessControl("ANY"))
		{
			// Routes untuk Categories dan Products
			userRoutes.GET("/categories", productController.GetAllCategories)
			userRoutes.GET("/categories/:id", productController.GetCategoryByID)

			userRoutes.GET("/products", productController.GetAllProducts)
			userRoutes.GET("/products/:id", productController.GetProductByID)
			userRoutes.POST("/products/:id/image", productController.UploadProductImage)

			// Routes untuk User Management
			userRoutes.GET("/user", userController.GetUserDetails)
			userRoutes.PUT("/user", userController.UpdateUser)

			// Routes untuk Transactions
			userRoutes.POST("/transactions", transactionController.CreateTransaction)
			userRoutes.GET("/transactions/:id", transactionController.GetTransactionByID)
			userRoutes.GET("/users/:user_id/transactions", transactionController.GetTransactionByUserID)

			// Reports Management
			userRoutes.POST("/reports/generate", reportController.GenerateReport)
			userRoutes.GET("/reports/download", reportController.DownloadReport)
		}
	}

	middleware.Logger.Info("Routes yang dilindungi JWT berhasil didaftarkan")

	// Endpoint debugging untuk mencetak semua rute
	r.GET("/debug/routes", func(c *gin.Context) {
		routes := r.Routes()
		for _, route := range routes {
			middleware.Logger.Info("Route", zap.String("method", route.Method), zap.String("path", route.Path))
		}
		c.JSON(http.StatusOK, routes)
	})

	// Menjalankan server di port 8080
	middleware.Logger.Info("Server dijalankan pada port 8080")
	if err := r.Run(":8080"); err != nil {
		middleware.Logger.Fatal("Server gagal dijalankan", zap.Error(err))
	}
}
