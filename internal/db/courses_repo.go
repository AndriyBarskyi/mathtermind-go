package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"mathtermind-go/internal/models"
)

// ListCourses returns a paginated list of courses.
func ListCourses(ctx context.Context, pool *pgxpool.Pool, limit, offset int) ([]models.Course, error) {
	rows, err := pool.Query(ctx, `
		SELECT id, topic, name, description, duration_min, created_at, updated_at
		FROM courses
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	courses := make([]models.Course, 0, limit)
	for rows.Next() {
		var c models.Course
		if err := rows.Scan(
			&c.ID,
			&c.Topic,
			&c.Name,
			&c.Description,
			&c.DurationMin,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		courses = append(courses, c)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return courses, nil
}
