package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/yertore/e-com/catalog-service/internal/domain"
)

// CategoryRepo — Postgres реализация repository.CategoryRepository.
type CategoryRepo struct {
	db *sql.DB
}

func NewCategoryRepo(db *sql.DB) *CategoryRepo {
	return &CategoryRepo{db: db}
}

func (r *CategoryRepo) Create(ctx context.Context, c *domain.Category) error {
	query := `
		INSERT INTO categories (id, name, slug, created_at)
		VALUES ($1, $2, $3, $4)`

	_, err := r.db.ExecContext(ctx, query, c.ID, c.Name, c.Slug, c.CreatedAt)
	if err != nil {
		return fmt.Errorf("CategoryRepo.Create: %w", err)
	}
	return nil
}

func (r *CategoryRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	query := `SELECT id, name, slug, created_at FROM categories WHERE id = $1`

	c := &domain.Category{}
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&c.ID, &c.Name, &c.Slug, &c.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("category %s not found", id)
		}
		return nil, fmt.Errorf("CategoryRepo.GetByID: %w", err)
	}
	return c, nil
}

func (r *CategoryRepo) GetBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	query := `SELECT id, name, slug, created_at FROM categories WHERE slug = $1`

	c := &domain.Category{}
	err := r.db.QueryRowContext(ctx, query, slug).
		Scan(&c.ID, &c.Name, &c.Slug, &c.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("category with slug %q not found", slug)
		}
		return nil, fmt.Errorf("CategoryRepo.GetBySlug: %w", err)
	}
	return c, nil
}

func (r *CategoryRepo) List(ctx context.Context) ([]*domain.Category, error) {
	query := `SELECT id, name, slug, created_at FROM categories ORDER BY name`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("CategoryRepo.List: %w", err)
	}
	defer rows.Close()

	var categories []*domain.Category
	for rows.Next() {
		c := &domain.Category{}
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("CategoryRepo.List scan: %w", err)
		}
		categories = append(categories, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("CategoryRepo.List rows: %w", err)
	}
	return categories, nil
}

func (r *CategoryRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM categories WHERE id = $1`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("CategoryRepo.Delete: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("category %s not found", id)
	}
	return nil
}
