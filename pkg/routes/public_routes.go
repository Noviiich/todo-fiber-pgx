package routes

import (
	"github.com/Noviiich/todo-fiber-pgx/app/controllers"
	"github.com/gofiber/fiber/v2"
)

// PublicRoutes func for describe group of public routes.
func PublicRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/api/v1")

	route.Get("/tasks", controllers.GetTasks)     // get list of all tasks
	route.Get("/task/:id", controllers.GetTask)   // get one task by ID
	route.Post("/task", controllers.CreateTask)   // create a new task
	route.Put("/task", controllers.UpdateTask)    // update one task by ID
	route.Delete("/task", controllers.DeleteTask) // delete one task by ID
}
