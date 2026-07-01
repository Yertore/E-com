# E-com — учебная микросервисная e-commerce платформа на Go

Pet-проект для демонстрации middle-senior Go backend навыков: микросервисная архитектура,
event-driven коммуникация, распределённые транзакции, observability.

## Статус

🚧 В разработке. Этап 1: Catalog Service.

## Архитектура

4 независимых сервиса в monorepo, каждый со своим `go.mod`, своей Postgres-базой
и собственным жизненным циклом деплоя:

| Сервис | Зона ответственности | Порт |
|---|---|---|
| `catalog-service` | товары, категории, остатки на складе | 8081 |
| `order-service` | жизненный цикл заказа, Saga-оркестрация | 8082 |
| `payment-service` | обработка платежей, идемпотентность | 8083 |
| `notification-service` | email/push уведомления (Kafka consumer) | 8084 |

Синхронная коммуникация — gRPC (контракты в `proto/`).
Асинхронная — Kafka (события `order.created`, `payment.completed`, `payment.failed`).

## Технологии

Go · PostgreSQL (индексы, ACID, транзакции, MVCC) · gRPC/REST · Kafka · Redis ·
Docker · Kubernetes · GitHub Actions CI/CD · Prometheus · Grafana · Loki ·
Jaeger/OpenTelemetry

## Паттерны

Clean architecture (domain / repository / service / transport) · Saga (оркестрация) ·
Outbox · CQRS (Order read-модель)

## Быстрый старт

​```bash
git clone https://github.com/yertore/e-com.git
cd e-com
docker compose up --build     # поднимает инфраструктуру, применяет миграции, запускает сервисы
curl http://localhost:8081/healthz
​```

## Запуск

### Все сразу:
​```bash
docker compose up --build
​```

### Только инфраструктура + миграции (без сервисов):
​```bash
docker compose up -d postgres-catalog migrate-catalog
docker compose up -d postgres-order migrate-order
docker compose up -d postgres-payment migrate-payment
​```

### Остановка:
​```bash
docker compose down -v
​```

## Структура репозитория

```
E-com/
├── catalog-service/      # независимый Go-модуль
├── order-service/        # независимый Go-модуль
├── payment-service/      # независимый Go-модуль
├── notification-service/ # независимый Go-модуль
├── pkg/                  # общие либы (logger, tracing, kafka-client)
├── proto/                # gRPC контракты между сервисами
├── deploy/               # k8s манифесты, observability конфиги
└── go.work               # связывает модули для локальной разработки
```

## Roadmap

- [x] Этап 1: monorepo skeleton, Catalog Service healthz, Postgres схема + миграции
- [ ] Этап 2: Catalog domain layer (CRUD) + clean architecture + тесты
- [ ] Этап 3: Order Service + Saga + Outbox + Kafka
- [ ] Этап 4: Redis, конкурентные резервы остатков (optimistic locking/MVCC)
- [ ] Этап 5: CQRS на Order Service
- [ ] Этап 6: Observability (Prometheus/Grafana/Loki/Jaeger)
- [ ] Этап 7: Kubernetes + cloud деплой
