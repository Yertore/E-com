package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/yertore/e-com/catalog-service/internal/domain"
)

// InventoryRepo — Postgres реализация repository.InventoryRepository.
// Ключевая особенность: optimistic locking через поле version.
type InventoryRepo struct {
	db *sql.DB
}

func NewInventoryRepo(db *sql.DB) *InventoryRepo {
	return &InventoryRepo{db: db}
}

func (r *InventoryRepo) Create(ctx context.Context, inv *domain.Inventory) error {
	query := `
		INSERT INTO inventory (product_id, quantity, reserved, version, updated_at)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.ExecContext(ctx, query,
		inv.ProductID, inv.Quantity, inv.Reserved, inv.Version, inv.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("InventoryRepo.Create: %w", err)
	}
	return nil
}

func (r *InventoryRepo) Get(ctx context.Context, productID uuid.UUID) (*domain.Inventory, error) {
	query := `
		SELECT product_id, quantity, reserved, version, updated_at
		FROM inventory WHERE product_id = $1`

	inv := &domain.Inventory{}
	err := r.db.QueryRowContext(ctx, query, productID).Scan(
		&inv.ProductID, &inv.Quantity, &inv.Reserved, &inv.Version, &inv.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("inventory for product %s not found", productID)
		}
		return nil, fmt.Errorf("InventoryRepo.Get: %w", err)
	}
	return inv, nil
}

// Update сохраняет изменения с проверкой version — это и есть optimistic locking.
//
// Условие WHERE version = $4 гарантирует атомарность:
// если между нашим чтением и записью кто-то другой уже изменил запись
// (version увеличился), наш UPDATE не найдёт строку (RowsAffected = 0)
// и мы вернём ошибку. Caller должен перечитать и повторить.
//
// Это аналог того, что Postgres делает внутри через MVCC —
// только на уровне приложения, без явных блокировок строк (SELECT FOR UPDATE).
func (r *InventoryRepo) Update(ctx context.Context, inv *domain.Inventory) error {
	query := `
		UPDATE inventory
		SET quantity = $1, reserved = $2, version = $3, updated_at = $4
		WHERE product_id = $5 AND version = $6`

	// version в БД = inv.Version - 1 (domain уже инкрементировал его в Reserve/Release/Confirm)
	// поэтому проверяем что в БД лежит предыдущая версия
	prevVersion := inv.Version - 1

	res, err := r.db.ExecContext(ctx, query,
		inv.Quantity, inv.Reserved, inv.Version, inv.UpdatedAt,
		inv.ProductID, prevVersion,
	)
	if err != nil {
		return fmt.Errorf("InventoryRepo.Update: %w", err)
	}

	n, _ := res.RowsAffected()
	if n == 0 {
		// Либо запись не найдена, либо version уже изменился — в обоих случаях конфликт
		return domain.ErrVersionConflict
	}
	return nil
}
