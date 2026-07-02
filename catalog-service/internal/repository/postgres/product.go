package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/yertore/e-com/catalog-service/internal/domain"
)

// ProductRepo — Postgres реализация repository.ProductRepository.
type ProductRepo struct {
	db *sql.DB
}

func NewProductRepo(db *sql.DB) *ProductRepo {
	return &ProductRepo{db: db}
}

func (r *ProductRepo) Create(ctx context.Context, p *domain.Product) error {
	query := `
		INSERT INTO products (id, category_id, sku, name, description, price, currency, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.ExecContext(ctx, query,
		p.ID, p.CategoryID, p.SKU, p.Name, p.Description,
		p.Price, p.Currency, p.CreatedAt, p.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("ProductRepo.Create: %w", err)
	}
	return nil
}

func (r *ProductRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	query := `
		SELECT id, category_id, sku, name, description, price, currency, created_at, updated_at
		FROM products WHERE id = $1`

	p := &domain.Product{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.CategoryID, &p.SKU, &p.Name, &p.Description,
		&p.Price, &p.Currency, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("product %s not found", id)
		}
		return nil, fmt.Errorf("ProductRepo.GetByID: %w", err)
	}
	return p, nil
}

func (r *ProductRepo) GetBySKU(ctx context.Context, sku string) (*domain.Product, error) {
	query := `
		SELECT id, category_id, sku, name, description, price, currency, created_at, updated_at
		FROM products WHERE sku = $1`

	p := &domain.Product{}
	err := r.db.QueryRowContext(ctx, query, sku).Scan(
		&p.ID, &p.CategoryID, &p.SKU, &p.Name, &p.Description,
		&p.Price, &p.Currency, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("product with SKU %q not found", sku)
		}
		return nil, fmt.Errorf("ProductRepo.GetBySKU: %w", err)
	}
	return p, nil
}

func (r *ProductRepo) ListByCategory(ctx context.Context, categoryID uuid.UUID) ([]*domain.Product, error) {
	query := `
		SELECT id, category_id, sku, name, description, price, currency, created_at, updated_at
		FROM products WHERE category_id = $1 ORDER BY name`

	rows, err := r.db.QueryContext(ctx, query, categoryID)
	if err != nil {
		return nil, fmt.Errorf("ProductRepo.ListByCategory: %w", err)
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		p := &domain.Product{}
		if err := rows.Scan(
			&p.ID, &p.CategoryID, &p.SKU, &p.Name, &p.Description,
			&p.Price, &p.Currency, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("ProductRepo.ListByCategory scan: %w", err)
		}
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ProductRepo.ListByCategory rows: %w", err)
	}
	return products, nil
}

func (r *ProductRepo) Update(ctx context.Context, p *domain.Product) error {
	query := `
		UPDATE products
		SET name = $1, description = $2, price = $3, currency = $4, updated_at = $5
		WHERE id = $6`

	res, err := r.db.ExecContext(ctx, query,
		p.Name, p.Description, p.Price, p.Currency, p.UpdatedAt, p.ID,
	)
	if err != nil {
		return fmt.Errorf("ProductRepo.Update: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("product %s not found", p.ID)
	}
	return nil
}

func (r *ProductRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM products WHERE id = $1`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("ProductRepo.Delete: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("product %s not found", id)
	}
	return nil
}
