package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"

	"scheduling-report/config"
	"scheduling-report/middlewares"
	"scheduling-report/routes"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Initialize logger
	config.InitLogger()

	// Connect to database
	config.ConnectDB()

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Scheduling Report System v1.0",
	})

	// Setup middlewares
	app.Use(middlewares.LoggingMiddleware())

	// Setup routes
	routes.SetupRoutes(app)

	// Start server
	port := fmt.Sprintf(":%s", config.Config.AppPort)
	go func() {
		if err := app.Listen(port); err != nil {
			log.Fatalf("Error starting Fiber: %v", err)
		}
	}()
	log.Printf("✅ Server is running on port %s", port)
	log.Printf("🔗 Health check: http://localhost%s/health", port)
	log.Printf("📋 API: http://localhost%s/api/report-configs", port)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("🛑 Shutting down server...")

	if err := app.Shutdown(); err != nil {
		log.Fatalf("Error shutting down Fiber: %v", err)
	}

	log.Println("✅ Server exited gracefully")
}
