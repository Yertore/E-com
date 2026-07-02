package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ErrInsufficientStock — ошибка домена, не инфраструктурная.
// Сервисный слой и транспорт смогут различить её от других ошибок
// и вернуть клиенту правильный HTTP-статус (409 Conflict).
var ErrInsufficientStock = errors.New("insufficient stock")

// ErrVersionConflict — оптимистичная блокировка: кто-то успел изменить
// запись между нашим чтением и записью. Нужно перечитать и повторить.
var ErrVersionConflict = errors.New("inventory version conflict, please retry")

// Inventory — складской учёт остатков конкретного товара.
//
// Два ключевых поля:
//   - Quantity  — физическое количество на складе
//   - Reserved  — зарезервировано (заказ оформлен, но ещё не оплачен)
//
// Доступно для продажи = Quantity - Reserved.
// Version — счётчик для optimistic locking (аналог MVCC на уровне приложения).
type Inventory struct {
	ProductID uuid.UUID
	Quantity  int
	Reserved  int
	Version   int
	UpdatedAt time.Time
}

// Available возвращает количество товара, доступного для продажи.
func (inv *Inventory) Available() int {
	return inv.Quantity - inv.Reserved
}

// Reserve резервирует qty единиц товара под заказ.
// Вызывается когда покупатель оформляет заказ (до оплаты).
//
// Optimistic locking: caller должен передать version который он прочитал.
// Если version не совпадает — значит между чтением и записью кто-то
// уже изменил запись, возвращаем ErrVersionConflict и caller повторяет.
func (inv *Inventory) Reserve(qty int, expectedVersion int) error {
	if qty <= 0 {
		return errors.New("reserve quantity must be positive")
	}
	if inv.Version != expectedVersion {
		return ErrVersionConflict
	}
	if inv.Available() < qty {
		return ErrInsufficientStock
	}

	inv.Reserved += qty
	inv.Version++
	inv.UpdatedAt = time.Now().UTC()
	return nil
}

// Release освобождает ранее зарезервированные qty единиц.
// Вызывается при отмене заказа (компенсирующая транзакция в Saga).
func (inv *Inventory) Release(qty int, expectedVersion int) error {
	if qty <= 0 {
		return errors.New("release quantity must be positive")
	}
	if inv.Version != expectedVersion {
		return ErrVersionConflict
	}
	if inv.Reserved < qty {
		return errors.New("cannot release more than reserved")
	}

	inv.Reserved -= qty
	inv.Version++
	inv.UpdatedAt = time.Now().UTC()
	return nil
}

// Confirm подтверждает резерв после успешной оплаты:
// списывает qty и из Quantity и из Reserved.
// Вызывается когда Payment Service сообщает об успешной оплате.
func (inv *Inventory) Confirm(qty int, expectedVersion int) error {
	if qty <= 0 {
		return errors.New("confirm quantity must be positive")
	}
	if inv.Version != expectedVersion {
		return ErrVersionConflict
	}
	if inv.Reserved < qty {
		return errors.New("cannot confirm more than reserved")
	}
	if inv.Quantity < qty {
		return errors.New("cannot confirm more than available quantity")
	}

	inv.Reserved -= qty
	inv.Quantity -= qty
	inv.Version++
	inv.UpdatedAt = time.Now().UTC()
	return nil
}
