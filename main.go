package main

import (
	"os"

	_ "github.com/Noviiich/todo-fiber-pgx/docs"
	"github.com/Noviiich/todo-fiber-pgx/pkg/configs"
	"github.com/Noviiich/todo-fiber-pgx/pkg/middleware"
	"github.com/Noviiich/todo-fiber-pgx/pkg/routes"
	"github.com/Noviiich/todo-fiber-pgx/pkg/utils"
	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
)

// @title API
// @version 0.0.1
// @description This is a simple Todo App API with Fiber and PostgreSQL.
// @termsOfService http://swagger.io/terms/
// @contact.name Noviiich
// @contact.email noviiich@yandex.ru
// @in header
// @BasePath /api/v1
func main() {
	// Define Fiber config.
	config := configs.FiberConfig()

	// Define a new Fiber app with config.
	app := fiber.New(config)

	// Middlewares.
	middleware.FiberMiddleware(app) // Register Fiber's middleware for app.

	// Routes.
	routes.SwaggerRoute(app)  // Register a route for API Docs (Swagger).
	routes.PublicRoutes(app)  // Register a public routes for app.
	routes.NotFoundRoute(app) // Register route for 404 Error.

	// Start server (with or without graceful shutdown).
	if os.Getenv("STAGE_STATUS") == "dev" {
		utils.StartServer(app)
	} else {
		utils.StartServerWithGracefulShutdown(app)
	}
}
