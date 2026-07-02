package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/yertore/e-com/catalog-service/internal/domain"
)

// CategoryRepository — контракт для работы с категориями.
// Domain определяет интерфейс, Postgres реализует.
// Благодаря этому можно подменить Postgres на любое другое хранилище
// (in-memory для тестов, другая БД) без изменения бизнес-логики.
type CategoryRepository interface {
	Create(ctx context.Context, category *domain.Category) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Category, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Category, error)
	List(ctx context.Context) ([]*domain.Category, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// ProductRepository — контракт для работы с товарами.
type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	GetBySKU(ctx context.Context, sku string) (*domain.Product, error)
	ListByCategory(ctx context.Context, categoryID uuid.UUID) ([]*domain.Product, error)
	Update(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// InventoryRepository — контракт для работы с остатками.
// Отдельный интерфейс от Product намеренно — разные паттерны доступа:
// инвентарь обновляется часто и под нагрузкой (optimistic locking),
// продукты почти только читаются.
type InventoryRepository interface {
	Get(ctx context.Context, productID uuid.UUID) (*domain.Inventory, error)
	Update(ctx context.Context, inv *domain.Inventory) error
	Create(ctx context.Context, inv *domain.Inventory) error
}