package controllers

import (
	"strconv"

	"github.com/Noviiich/todo-fiber-pgx/app/models"
	"github.com/Noviiich/todo-fiber-pgx/pkg/utils"
	"github.com/Noviiich/todo-fiber-pgx/platform/database"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// GetTasks func gets all exists tasks.
// @Description Get all exists tasks.
// @Summary get all exists tasks
// @Tags Tasks
// @Accept json
// @Produce json
// @Success 200 {array} models.Task
// @Router /v1/task [get]
func GetTasks(c *fiber.Ctx) error {
	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Get all tasks.
	tasks, err := db.GetTasks()
	if err != nil {
		// Return, if tasks not found.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "tasks were not found",
			"count": 0,
			"tasks": nil,
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"count": len(tasks),
		"tasks": tasks,
	})
}

// GetTask func gets task by given ID or 404 error.
// @Description Get task by given ID.
// @Summary get task by given ID
// @Tags Task
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} models.Task
// @Router /v1/task/{id} [get]
func GetTask(c *fiber.Ctx) error {
	// Catch task ID from URL.
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Get task by ID.
	task, err := db.GetTask(id)
	if err != nil {
		// Return, if task not found.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "task with the given ID is not found",
			"task":  nil,
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"task":  task,
	})
}

// CreateTask func for creates a new task.
// @Description Create a new task.
// @Summary create a new task
// @Tags Task
// @Accept json
// @Produce json
// @Param title body string true "Title"
// @Param description body string true "Description"
// @Success 200 {object} models.Task
// @Router /v1/task [post]
func CreateTask(c *fiber.Ctx) error {
	// Create new Task struct
	task := &models.Task{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(task); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Validate task fields.
	if err := validator.New().Struct(task); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	// Create task by given model.
	if err := db.CreateTask(task); err != nil {
		// Return status 500 and error message.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"task":  task,
	})
}

// UpdateTask func for updates task by given ID.
// @Description Update task.
// @Summary update task
// @Tags Task
// @Accept json
// @Produce json
// @Param id body string true "Task ID"
// @Param title body string true "Title"
// @Param description body string true "Description"
// @Param status body integer true "Status"
// @Success 202 {string} status "ok"
// @Router /v1/task [put]
func UpdateTask(c *fiber.Ctx) error {
	// Create new Task struct
	task := &models.Task{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(task); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Checking, if task with given ID is exists.
	foundedTask, err := db.GetTask(task.ID)
	if err != nil {
		// Return status 404 and task not found error.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "task with this ID not found",
		})
	}

	// Validate task fields.
	if err := validator.New().Struct(task); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	// Update task by given ID.
	if err := db.UpdateTask(foundedTask.ID, task); err != nil {
		// Return status 500 and error message.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Return status 201.
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error": false,
		"msg":   nil,
	})
}

// DeleteTask func for deletes task by given ID.
// @Description Delete task by given ID.
// @Summary delete task by given ID
// @Tags Task
// @Accept json
// @Produce json
// @Param id body string true "Task ID"
// @Success 204 {string} status "ok"
// @Router /v1/task [delete]
func DeleteTask(c *fiber.Ctx) error {
	// Create new Task struct
	task := &models.Task{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(task); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create a new validator for a Task model.
	if err := validator.New().Struct(task); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Checking, if task with given ID is exists.
	foundedTask, err := db.GetTask(task.ID)
	if err != nil {
		// Return status 404 and task not found error.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "task with this ID not found",
		})
	}

	// Delete task by given ID.
	if err := db.DeleteTask(foundedTask.ID); err != nil {
		// Return status 500 and error message.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Return status 204 no content.
	return c.SendStatus(fiber.StatusNoContent)
}
