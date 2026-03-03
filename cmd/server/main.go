package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	authService "github.com/Ju4n-Dieg0/Go_Points/internal/application/auth"
	companyService "github.com/Ju4n-Dieg0/Go_Points/internal/application/company"
	consumerService "github.com/Ju4n-Dieg0/Go_Points/internal/application/consumer"
	pointService "github.com/Ju4n-Dieg0/Go_Points/internal/application/point"
	productService "github.com/Ju4n-Dieg0/Go_Points/internal/application/product"
	rewardService "github.com/Ju4n-Dieg0/Go_Points/internal/application/reward"
	subscriptionService "github.com/Ju4n-Dieg0/Go_Points/internal/application/subscription"
	"github.com/Ju4n-Dieg0/Go_Points/internal/config"
	"github.com/Ju4n-Dieg0/Go_Points/internal/database"
	"github.com/Ju4n-Dieg0/Go_Points/internal/domain/auth"
	"github.com/Ju4n-Dieg0/Go_Points/internal/domain/company"
	"github.com/Ju4n-Dieg0/Go_Points/internal/domain/consumer"
	"github.com/Ju4n-Dieg0/Go_Points/internal/domain/point"
	"github.com/Ju4n-Dieg0/Go_Points/internal/domain/product"
	"github.com/Ju4n-Dieg0/Go_Points/internal/domain/reward"
	"github.com/Ju4n-Dieg0/Go_Points/internal/domain/subscription"
	"github.com/Ju4n-Dieg0/Go_Points/internal/infrastructure/http/handler"
	"github.com/Ju4n-Dieg0/Go_Points/internal/infrastructure/http/routes"
	"github.com/Ju4n-Dieg0/Go_Points/internal/infrastructure/persistence"
	"github.com/Ju4n-Dieg0/Go_Points/internal/infrastructure/service"
	"github.com/Ju4n-Dieg0/Go_Points/internal/shared/errors"
	"github.com/Ju4n-Dieg0/Go_Points/internal/shared/logger"
	"github.com/Ju4n-Dieg0/Go_Points/internal/shared/middleware"
	"github.com/Ju4n-Dieg0/Go_Points/internal/shared/notifications"
	"github.com/gofiber/fiber/v3"
)

func main() {
	// Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Inicializar logger
	logger.Setup(cfg.App.LogLevel, cfg.App.Environment)
	logger.Info("Starting application",
		"name", cfg.App.Name,
		"environment", cfg.App.Environment,
		"version", "1.0.0",
	)

	// Conectar a la base de datos
	db, err := database.NewDatabase(&cfg.Database, cfg.App.Environment)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("Failed to close database connection", "error", err)
		}
	}()

	// Ejecutar migraciones automáticas
	logger.Info("Running database migrations...")
	if err := db.AutoMigrate(
		&auth.User{},
		&company.Company{},
		&subscription.Subscription{},
		&consumer.Consumer{},
		&product.Product{},
		&point.ConsumerCompanyPoints{},
		&point.PointTransaction{},
		&point.CompanyRankConfig{},
		&reward.Reward{},
		&reward.RewardPath{},
		&reward.RewardPathItem{},
	); err != nil {
		logger.Error("Failed to run migrations", "error", err)
		os.Exit(1)
	}

	// Inyección de dependencias manual
	// Repositories
	authRepo := persistence.NewAuthRepository(db.GetDB())
	companyRepo := persistence.NewCompanyRepository(db.GetDB())
	subscriptionRepo := persistence.NewSubscriptionRepository(db.GetDB())
	consumerRepo := persistence.NewConsumerRepository(db.GetDB())
	productRepo := persistence.NewProductRepository(db.GetDB())
	balanceRepo := persistence.NewBalanceRepository(db.GetDB())
	transactionRepo := persistence.NewTransactionRepository(db.GetDB())
	rankConfigRepo := persistence.NewRankConfigRepository(db.GetDB())
	rewardRepo := persistence.NewRewardRepository(db.GetDB())
	rewardPathRepo := persistence.NewRewardPathRepository(db.GetDB())
	rewardPathItemRepo := persistence.NewRewardPathItemRepository(db.GetDB())

	// Services
	emailService := service.NewStubEmailService()
	fileService := service.NewLocalFileService(cfg.File)
	
	// Notification services - usando composite para enviar a múltiples destinos
	logNotificationService := notifications.NewLogNotificationService(logger.GetLogger())
	emailNotificationService := notifications.NewEmailNotificationService(&cfg.Email)
	notificationService := notifications.NewCompositeNotificationService(
		logNotificationService,
		emailNotificationService,
	)
	authSvc := authService.NewService(authRepo, emailService, &cfg.JWT)
	companySvc := companyService.NewService(companyRepo, subscriptionRepo, db.GetDB())
	subscriptionSvc := subscriptionService.NewService(subscriptionRepo, companyRepo, db.GetDB())
	consumerSvc := consumerService.NewService(consumerRepo)
	productSvc := productService.NewService(productRepo, fileService, subscriptionSvc)
	pointSvc := pointService.NewService(db.GetDB(), balanceRepo, transactionRepo, rankConfigRepo, notificationService, cfg.Points)
	rewardSvc := rewardService.NewRewardService(db.GetDB(), rewardRepo, rewardPathRepo, rewardPathItemRepo, companyRepo, productRepo, subscriptionSvc)

	// Handlers
	authHandler := handler.NewAuthHandler(authSvc)
	companyHandler := handler.NewCompanyHandler(companySvc)
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionSvc)
	consumerHandler := handler.NewConsumerHandler(consumerSvc)
	productHandler := handler.NewProductHandler(productSvc)
	pointHandler := handler.NewPointHandler(pointSvc)
	rewardHandler := handler.NewRewardHandler(rewardSvc)

	// Crear aplicación Fiber
	app := createFiberApp(cfg)

	// Configurar rutas
	setupRoutes(app, db, authHandler, companyHandler, subscriptionHandler, consumerHandler, productHandler, pointHandler, rewardHandler, cfg)

	// Iniciar servidor en una goroutine
	go func() {
		addr := cfg.GetServerAddress()
		logger.Info("Server starting", "address", addr)

		if err := app.Listen(addr); err != nil {
			logger.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	gracefulShutdown(app, cfg.Server.ShutdownTimeout)
}

// createFiberApp crea y configura la aplicación Fiber
func createFiberApp(cfg *config.Config) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      cfg.App.Name,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		ErrorHandler: func(c fiber.Ctx, err error) error {
			return errors.ErrorHandler()(c)
		},
	})

	// Middlewares globales
	app.Use(errors.RecoverMiddleware())
	app.Use(middleware.RequestID())
	app.Use(middleware.Logger())
	app.Use(middleware.CORS())
	app.Use(middleware.Security())

	return app
}

// setupRoutes configura las rutas de la aplicación
func setupRoutes(
	app *fiber.App,
	db *database.Database,
	authHandler *handler.AuthHandler,
	companyHandler *handler.CompanyHandler,
	subscriptionHandler *handler.SubscriptionHandler,
	consumerHandler *handler.ConsumerHandler,
	productHandler *handler.ProductHandler,
	pointHandler *handler.PointHandler,
	rewardHandler *handler.RewardHandler,
	cfg *config.Config,
) {
	// Health check endpoint
	app.Get("/health", func(c fiber.Ctx) error {
		// Verificar estado de la base de datos
		if err := db.HealthCheck(); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status":   "error",
				"database": "unhealthy",
				"error":    err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"status":    "ok",
			"database":  "healthy",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	// API v1 routes group
	api := app.Group("/api/v1")

	// Ruta de bienvenida
	api.Get("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to Go Points API",
			"version": "1.0.0",
			"status":  "running",
		})
	})

	// Configurar rutas de módulos
	routes.SetupAuthRoutes(api, authHandler, &cfg.JWT)
	routes.SetupCompanyRoutes(api, companyHandler, &cfg.JWT)
	routes.SetupSubscriptionRoutes(api, subscriptionHandler, &cfg.JWT)
	routes.SetupConsumerRoutes(api, consumerHandler, &cfg.JWT)
	routes.SetupProductRoutes(api, productHandler, &cfg.JWT)
	routes.SetupPointRoutes(api, pointHandler, &cfg.JWT)
	routes.SetupRewardRoutes(api, rewardHandler, &cfg.JWT)

	// Aquí se agregarán más rutas cuando se implementen otros módulos
	// Ejemplo:
	// routes.SetupPointsRoutes(api, pointsHandler, &cfg.JWT)

	// Ruta 404
	app.Use(func(c fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Route not found",
			"path":    c.Path(),
		})
	})
}

// gracefulShutdown maneja el cierre elegante del servidor
func gracefulShutdown(app *fiber.App, timeout int) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	// Shutdown del servidor Fiber
	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	logger.Info("Server exited gracefully")
}
