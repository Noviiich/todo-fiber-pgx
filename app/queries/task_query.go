package queries

import (
	"context"
	"fmt"

	"github.com/Noviiich/todo-fiber-pgx/app/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskQueries struct {
	Pool *pgxpool.Pool
}

// GetTasks method for getting all tasks.
func (q *TaskQueries) GetTasks() ([]models.Task, error) {
	query := `
        SELECT id, title, description, status, created_at, updated_at
        FROM tasks
        ORDER BY created_at DESC
    `

	rows, err := q.Pool.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("error querying tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning task row: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating task rows: %w", err)
	}

	return tasks, nil
}

// GetTask method for getting one task by given ID.
func (q *TaskQueries) GetTask(id int) (models.Task, error) {
	query := `
		SELECT id, title, description, status, created_at, updated_at 
		FROM tasks WHERE id = $1
	`
	task := models.Task{}

	err := q.Pool.QueryRow(context.Background(), query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		return models.Task{}, err
	}

	return task, nil
}

// CreateTask method for creating task by given Task object.
func (q *TaskQueries) CreateTask(t *models.Task) (int, error) {
	query := `
        INSERT INTO tasks (title, description)
        VALUES ($1, $2)
        RETURNING id
    `

	var id int
	err := q.Pool.QueryRow(context.Background(), query, t.Title, t.Description).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error creating task: %w", err)
	}

	fmt.Printf("Created task with ID: %d\n", id)
	return id, nil
}

// UpdateTask method for updating task by given Task object.
func (q *TaskQueries) UpdateTask(id int, t *models.Task) error {
	query := `
        UPDATE tasks
        SET title = $2, description = $3, status = $4, updated_at = NOW()
        WHERE id = $1
    `

	commandTag, err := q.Pool.Exec(context.Background(), query, id, t.Title, t.Description, t.Status)
	if err != nil {
		return fmt.Errorf("error completing task: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("no task found with id %d", id)
	}

	return nil
}

// DeleteTask method for delete task by given ID.
func (q *TaskQueries) DeleteTask(id int) error {
	query := `DELETE FROM tasks WHERE id = $1`

	commandTag, err := q.Pool.Exec(context.Background(), query, id)
	if err != nil {
		return fmt.Errorf("error deleting task: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("no task found with id %d", id)
	}

	return nil
}
