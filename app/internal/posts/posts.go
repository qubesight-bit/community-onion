package posts

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Post struct {
	ID          string
	Title       string
	Body        string
	PublishedAt time.Time
}

func (p Post) PublishedAtUTC() string {
	return p.PublishedAt.UTC().Format("2006-01-02 15:04 UTC")
}

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) ListPublished(
	ctx context.Context,
	limit int,
) ([]Post, error) {
	if limit < 1 || limit > 100 {
		return nil, fmt.Errorf("limit must be between 1 and 100")
	}

	rows, err := r.pool.Query(
		ctx,
		`SELECT id::text, title, body, published_at
		 FROM published_posts
		 ORDER BY published_at DESC
		 LIMIT $1`,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("query published posts: %w", err)
	}
	defer rows.Close()

	items := make([]Post, 0, limit)

	for rows.Next() {
		var item Post

		if err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.Body,
			&item.PublishedAt,
		); err != nil {
			return nil, fmt.Errorf("scan published post: %w", err)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate published posts: %w", err)
	}

	return items, nil
}
