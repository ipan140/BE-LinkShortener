package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"

	"link-shortener/app/modules/analytics"
	"link-shortener/app/modules/auth"
	"link-shortener/app/modules/link"
	"link-shortener/app/modules/redirect"
	"link-shortener/common/middleware"
	"link-shortener/config"
	"link-shortener/docs"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// 1. Load Configurations
	cfg := config.LoadConfig()
	log.Printf("[Bootstrap] Starting %s in %s mode...", cfg.App.Name, cfg.App.Env)

	// 2. Initialize Relational Database (GORM)
	db, err := config.InitDatabase(&cfg.Database)
	if err != nil {
		log.Fatalf("[Bootstrap] Fatal: Database connection error: %v", err)
	}

	// Run Auto Migrations for all Modular Monolith Entities
	log.Println("[Bootstrap] Running database schema auto-migrations...")
	if err := db.AutoMigrate(
		&auth.User{},
		&link.Link{},
		&analytics.Click{},
	); err != nil {
		log.Fatalf("[Bootstrap] Migration failed: %v", err)
	}

	// 3. Initialize Redis Connection Pool
	rdb, err := config.InitRedis(&cfg.Redis)
	if err != nil {
		log.Printf("[Bootstrap] Info: Redis cache is disabled or not reachable. Running with direct DB fallback.")
	}

	// 4. Initialize Dependency Injections per Module (Clean Layered Architecture)

	// Module: Auth
	authRepo := auth.NewRepository(db)
	authCache := auth.NewRedisCache(rdb)
	authService := auth.NewService(authRepo, authCache, cfg)
	authHandler := auth.NewHandler(authService)
	authRoute := auth.NewRoute(authHandler, cfg.JWT.Secret)

	// Module: Link
	linkRepo := link.NewRepository(db)
	linkCache := link.NewRedisCache(rdb)
	linkService := link.NewService(linkRepo, linkCache, cfg)
	linkHandler := link.NewHandler(linkService)
	linkRoute := link.NewRoute(linkHandler, cfg.JWT.Secret, authService)

	// Module: Redirect (Ultra-Fast Redirector & Click Collector)
	redirectRepo := redirect.NewRepository(db)
	redirectCache := redirect.NewRedisCache(rdb, cfg.Redis.QueueKey)
	redirectService := redirect.NewService(redirectRepo, redirectCache)
	redirectHandler := redirect.NewHandler(redirectService)
	redirectRoute := redirect.NewRoute(redirectHandler)

	// Module: Analytics
	analyticsRepo := analytics.NewRepository(db)
	analyticsQueue := analytics.NewRedisQueue(rdb, cfg.Redis.QueueKey)
	analyticsService := analytics.NewService(analyticsRepo, analyticsQueue, cfg)
	analyticsHandler := analytics.NewHandler(analyticsService)
	analyticsRoute := analytics.NewRoute(analyticsHandler, cfg.JWT.Secret, authService)

	// 5. Start Background Workers
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if rdb != nil {
		go analyticsService.StartClickCollectorWorker(ctx)
	}

	// 6. Setup Echo Server and Presentation Layer
	e := echo.New()
	e.HideBanner = true
	e.Use(echoMiddleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.RequestLogger())
	e.Use(middleware.PrometheusMetrics())

	// Prometheus Metrics Endpoint for Grafana monitoring
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{
			"status":    "healthy",
			"timestamp": time.Now().Format(time.RFC3339),
			"app":       cfg.App.Name,
		})
	})

	// Register Swagger / OpenAPI Docs
	docs.RegisterDocsEndpoint(e)

	// Public QR code endpoint
	e.GET("/qr/:code", linkHandler.QRCode)

	// Register Redirect Routes with Anti-Abuse Rate Limiter (max 120 redirects per 10s per IP)
	redirectGroup := e.Group("")
	if rdb != nil {
		redirectLimiter := middleware.RedisRateLimiter(rdb, middleware.RateLimiterConfig{
			Limit:  120,
			Window: 10 * time.Second,
			Prefix: "redirect",
		})
		redirectGroup.Use(redirectLimiter)
	}
	redirectRoute.RegisterRoutesGroup(redirectGroup)

	// Register REST API v1 Group with API Rate Limiter (max 60 requests per 10s per IP)
	apiV1 := e.Group("/api/v1")
	if rdb != nil {
		apiLimiter := middleware.RedisRateLimiter(rdb, middleware.RateLimiterConfig{
			Limit:  60,
			Window: 10 * time.Second,
			Prefix: "api",
		})
		apiV1.Use(apiLimiter)
	}

	authRoute.RegisterRoutes(apiV1)
	linkRoute.RegisterRoutes(apiV1)
	analyticsRoute.RegisterRoutes(apiV1)

	// 7. Start HTTP Server in Goroutine
	serverAddr := fmt.Sprintf(":%s", cfg.App.Port)
	go func() {
		log.Printf("[Server] %s listening on http://localhost%s", cfg.App.Name, serverAddr)
		if err := e.Start(serverAddr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[Server] Failed to start HTTP server: %v", err)
		}
	}()

	// 8. Graceful Shutdown on OS Signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("[Server] Received shutdown signal. Initiating graceful shutdown...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	cancel() // Stop background workers

	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Printf("[Server] Error during HTTP shutdown: %v", err)
	}

	if rdb != nil {
		_ = rdb.Close()
	}

	log.Println("[Server] Server stopped cleanly. Bye!")
}
