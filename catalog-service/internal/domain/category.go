package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Category — бизнес-сущность категории товара.
// Не знает ничего про Postgres, HTTP или Kafka — только бизнес-правила.
type Category struct {
	ID        uuid.UUID
	Name      string
	Slug      string
	CreatedAt time.Time
}

// NewCategory создаёт новую категорию с валидацией бизнес-правил.
// Конструктор — единственный легальный способ создать Category,
// гарантирует что объект всегда в валидном состоянии.
func NewCategory(name string) (*Category, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("category name cannot be empty")
	}
	if len(name) > 255 {
		return nil, errors.New("category name too long (max 255)")
	}

	return &Category{
		ID:        uuid.New(),
		Name:      name,
		Slug:      toSlug(name),
		CreatedAt: time.Now().UTC(),
	}, nil
}

// toSlug конвертирует название в URL-friendly slug.
// Например: "Электроника и гаджеты" → "elektronika-i-gadzhety"
// Упрощённая реализация — в продакшене используют библиотеку типа gosimple/slug.
func toSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	return slug
}
