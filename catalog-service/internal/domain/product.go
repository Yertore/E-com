package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Product — бизнес-сущность товара.
type Product struct {
	ID          uuid.UUID
	CategoryID  uuid.UUID
	SKU         string
	Name        string
	Description string
	Price       float64
	Currency    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewProduct создаёт новый товар с валидацией бизнес-правил.
func NewProduct(categoryID uuid.UUID, sku, name, description string, price float64, currency string) (*Product, error) {
	sku = strings.TrimSpace(sku)
	if sku == "" {
		return nil, errors.New("product SKU cannot be empty")
	}
	if len(sku) > 64 {
		return nil, errors.New("product SKU too long (max 64)")
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("product name cannot be empty")
	}

	if price < 0 {
		return nil, errors.New("product price cannot be negative")
	}

	currency = strings.ToUpper(strings.TrimSpace(currency))
	if currency == "" {
		currency = "KZT"
	}
	if len(currency) != 3 {
		return nil, errors.New("currency must be a 3-letter ISO code (e.g. KZT, USD)")
	}

	if categoryID == uuid.Nil {
		return nil, errors.New("product must belong to a category")
	}

	now := time.Now().UTC()
	return &Product{
		ID:          uuid.New(),
		CategoryID:  categoryID,
		SKU:         sku,
		Name:        name,
		Description: strings.TrimSpace(description),
		Price:       price,
		Currency:    currency,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// UpdatePrice обновляет цену товара с валидацией.
// Метод на сущности — бизнес-правило живёт в domain, а не в сервисном слое.
func (p *Product) UpdatePrice(newPrice float64) error {
	if newPrice < 0 {
		return errors.New("price cannot be negative")
	}
	p.Price = newPrice
	p.UpdatedAt = time.Now().UTC()
	return nil
}
