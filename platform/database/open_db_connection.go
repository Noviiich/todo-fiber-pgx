package database

import (
	"github.com/Noviiich/todo-fiber-pgx/app/queries"
)

// Queries struct for collect all app queries.
type Queries struct {
	*queries.TaskQueries // load queries from Task model
}

// OpenDBConnection func for opening database connection.
func OpenDBConnection() (*Queries, error) {
	// Define a new PostgreSQL connection.
	pool, err := PostgreSQLConnection()
	if err != nil {
		return nil, err
	}

	return &Queries{
		// Set queries from models:
		TaskQueries: &queries.TaskQueries{
			Pool: pool,
		},
	}, nil
}
