package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var (
	version = "dev"
	appName = "Knowledge Garden CLI - Web Dashboard"
)

func main() {
	// Get working directory for absolute paths
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}

	app := setupFiber()
	setupRoutes(app, wd)
	setupGracefulShutdown(app)

	port := getEnv("PORT", "3000")
	log.Printf("Starting %s on port %s", appName, port)

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupFiber() *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:               appName,
		DisableStartupMessage: false,
		EnablePrintRoutes:     getEnv("ENV", "development") == "development",
		ErrorHandler:          customErrorHandler,
		// Not using template engine - reading HTML files directly with os.ReadFile()
	})

	// Middleware
	app.Use(
		logger.New(logger.Config{
			// Format:     "[${time}] ${status} - ${method} ${path} (${latency})",
			// TimeFormat: "2006-01-02 15:04:05",
			// Output:     os.Stdout,
		}),
	)

	app.Use(recover.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	app.Use(cors.New(cors.Config{
		// AllowOrigins:     []string{"https://cli-notes-api.kelanach.xyz", "http://localhost:3000", "http://localhost:8080"},
		AllowMethods: "GET,POST,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
		// AllowCredentials: true,
	}))

	// Security headers
	app.Use(securityHeaders)

	// Rate limiting: 120 req/min per IP using Fiber's built-in middleware
	app.Use(limiter.New(limiter.Config{
		Max:        120,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(http.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit exceeded",
			})
		},
	}))

	return app
}

func securityHeaders(c *fiber.Ctx) error {
	c.Set("X-Content-Type-Options", "nosniff")
	c.Set("X-Frame-Options", "DENY")
	c.Set("X-XSS-Protection", "1; mode=block")
	c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
	c.Set("Permissions-Policy", "default-src 'self'")
	c.Set("X-Application-Version", version)
	return c.Next()
}

func setupRoutes(app *fiber.App, wd string) {
	// Build absolute paths
	staticPath := filepath.Join(wd, "static")
	templatesPath := filepath.Join(wd, "templates")

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"version": version,
			"app":     appName,
		})
	})

	// Static files
	app.Static("/static", staticPath)

	// Main page - uses index.html as the single template
	app.Get("/", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html; charset=utf-8")
		content, err := os.ReadFile(filepath.Join(templatesPath, "index.html"))
		if err != nil {
			return c.Status(fiber.StatusNotFound).SendString("Template not found")
		}
		return c.SendString(string(content))
	})

	// Tutorial hub
	app.Get("/tutorial", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html; charset=utf-8")
		content, err := os.ReadFile(filepath.Join(templatesPath, "tutorial.html"))
		if err != nil {
			return c.Status(fiber.StatusNotFound).SendString("Template not found")
		}
		return c.SendString(string(content))
	})

	// Self-hosting tutorial
	app.Get("/tutorial/self-hosting", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html; charset=utf-8")
		content, err := os.ReadFile(filepath.Join(templatesPath, "tutorial-self-hosting.html"))
		if err != nil {
			return c.Status(fiber.StatusNotFound).SendString("Template not found")
		}
		return c.SendString(string(content))
	})

	// CLI reference
	app.Get("/tutorial/cli-reference", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html; charset=utf-8")
		content, err := os.ReadFile(filepath.Join(templatesPath, "tutorial-cli-reference.html"))
		if err != nil {
			return c.Status(fiber.StatusNotFound).SendString("Template not found")
		}
		return c.SendString(string(content))
	})

	// 404 handler - must be last
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html; charset=utf-8")
		content, err := os.ReadFile(filepath.Join(templatesPath, "404.html"))
		if err != nil {
			return c.Status(fiber.StatusNotFound).SendString("404 - Page Not Found")
		}
		c.Status(fiber.StatusNotFound)
		return c.SendString(string(content))
	})
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	c.Set("Content-Type", "application/json")

	if code == fiber.StatusNotFound {
		c.Status(fiber.StatusNotFound)
		return c.SendString("404 - Page Not Found")
	}

	c.Status(code)
	return c.JSON(fiber.Map{
		"error": err.Error(),
	})
}

func setupGracefulShutdown(app *fiber.App) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Shutting down server...")
		if err := app.ShutdownWithTimeout(5 * time.Second); err != nil {
			log.Printf("Error during shutdown: %v", err)
		}
	}()
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
