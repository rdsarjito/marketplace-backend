package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/rdsarjito/marketplace-backend/config"
	"github.com/rdsarjito/marketplace-backend/handlers"
	"github.com/rdsarjito/marketplace-backend/middleware"
	"github.com/rdsarjito/marketplace-backend/repositories"
	"github.com/rdsarjito/marketplace-backend/services"
	"github.com/rdsarjito/marketplace-backend/storage"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	db := config.InitDatabase()

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"status":  false,
				"message": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(cors.New())
	app.Use(logger.New())

	// Static files
	app.Static("/uploads", "./uploads")

	// API external for data province & city
	provinceCityApiURL := os.Getenv("API_LOCATION")

	// Initialize repositories
	provinceCityRepository := repositories.NewProvinceCityRepository(provinceCityApiURL)
	userRepository := repositories.NewUserRepository(db)
	shopRepository := repositories.NewShopRepository(db)
	categoryRepository := repositories.NewCategoryRepository(db)
	addressRepository := repositories.NewAddressRepository(db)
	productRepository := repositories.NewProductRepository(db)
	trxRepository := repositories.NewTRXRepository(db)

	// Initialize shared services
	emailService := services.NewEmailService()
	authService := services.NewAuthService(userRepository, shopRepository, provinceCityRepository, emailService)
	userService := services.NewUserService(userRepository, addressRepository)
	categoryService := services.NewCategoryService(categoryRepository)
	shopService := services.NewShopService(shopRepository)
	productService := services.NewProductService(productRepository, shopRepository, categoryRepository)
	trxService := services.NewTRXService(trxRepository, productRepository, addressRepository, shopRepository, categoryRepository)
	mediaStorage, err := storage.NewMinioStorageFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	provinceCityHandler := handlers.NewProvinceCityHandler(provinceCityRepository)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	shopHandler := handlers.NewShopHandler(shopService)
	productHandler := handlers.NewProductHandler(productService, mediaStorage)
	trxHandler := handlers.NewTRXHandler(trxService)

	// Initialize middleware
	authMiddleware := middleware.AuthMiddleware(userService)

	// API routes
	api := app.Group("/api/v1")

	// Auth routes (public)
	api.Post("/auth/register", authHandler.RegisterUser)
	api.Post("/auth/login", authHandler.LoginUser)
	api.Post("/auth/forgot-password", authHandler.ForgotPassword)
	api.Post("/auth/reset-password", authHandler.ResetPassword)
	api.Get("/auth/google", authHandler.GoogleLogin)
	api.Get("/auth/google/callback", authHandler.GoogleCallback)

	// Province & City routes (public)
	api.Get("/provcity/listprovincies", provinceCityHandler.GetListProvince)
	api.Get("/provcity/detailprovince/:prov_id", provinceCityHandler.GetDetailProvince)
	api.Get("/provcity/listcities/:prov_id", provinceCityHandler.GetListCity)
	api.Get("/provcity/detailcity/:city_id", provinceCityHandler.GetDetailCity)

	// Protected routes
	api.Use(authMiddleware)

	// User routes
	api.Get("/user", userHandler.GetMyProfile)
	api.Put("/user", userHandler.UpdateProfile)
	api.Get("/user/alamat", userHandler.GetMyAddress)
	api.Get("/user/alamat/:id", userHandler.GetDetailAddress)
	api.Post("/user/alamat", userHandler.CreateAddressUser)
	api.Put("/user/alamat/:id", userHandler.UpdateAddressUser)
	api.Delete("/user/alamat/:id", userHandler.DeleteAddressUser)

	// Category routes
	api.Get("/category", categoryHandler.GetListCategory)
	api.Get("/category/:id", categoryHandler.GetDetailCategory)
	api.Post("/category", categoryHandler.CreateCategory)
	api.Put("/category/:id", categoryHandler.UpdateCategory)
	api.Delete("/category/:id", categoryHandler.DeleteCategory)

	// Shop routes
	api.Get("/toko/my", shopHandler.MyShop)
	api.Get("/toko", shopHandler.GetListShop)
	api.Get("/toko/:id_toko", shopHandler.GetDetailShop)
	api.Put("/toko/:id_toko", shopHandler.UpdateProfileShop)

	// Product routes
	api.Get("/product", productHandler.GetListProduct)
	api.Get("/product/:id", productHandler.GetDetailProduct)
	api.Post("/product", productHandler.CreateProduct)
	api.Put("/product/:id", productHandler.UpdateProduct)
	api.Delete("/product/:id", productHandler.DeleteProduct)
	api.Post("/product/:id/photo", productHandler.UploadProductPhoto)

	// Transaction routes
	api.Get("/trx", trxHandler.GetListTRX)
	api.Get("/trx/:id", trxHandler.GetDetailTRX)
	api.Post("/trx", trxHandler.CreateTRX)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  true,
			"message": "Server is running",
		})
	})

	// Start server
	port := cfg.AppPort
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on %s:%s", cfg.AppHost, port)
	if err := app.Listen(fmt.Sprintf("%s:%s", cfg.AppHost, port)); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
