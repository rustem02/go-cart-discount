краткое описание архитектуры:

- расчёт корзины `/price`
- пакетная обработка нескольких корзин `/price/bulk` с конкурентным выполнением и ограничением по числу goroutine
- паттерн Strategy для разных типов скидок
- простая многослойная архитектура (api / service / cmd)

стек

- Go 1.19
- HTTP-сервер на стандартной библиотеке `net/http`
- Docker + Docker Compose

запуск

```bash
docker-compose up --build
```

примеры запросов:

POST /price

```bash
{
  "discount_type": "percent",
  "discount_value": 10,
  "items": [
    { "id": "A", "price": 100, "qty": 2 },
    { "id": "B", "price": 50,  "qty": 1 }
  ]
}
```

response:

```bash
{
  "subtotal": 250,
  "discount": 25,
  "total": 225
}
```

POST /price/bulk

```bash
[
  {
    "discount_type": "percent",
    "discount_value": 10,
    "items": [
      { "id": "A", "price": 100, "qty": 2 }
    ]
  },
  {
    "discount_type": "bulk",
    "discount_value": 20,
    "items": [
      { "id": "A", "price": 10, "qty": 5 },
      { "id": "B", "price": 10, "qty": 5 }
    ]
  }
]
```

response:

```bash
[
  {
    "subtotal": 200,
    "discount": 20,
    "total": 180
  },
  {
    "subtotal": 100,
    "discount": 20,
    "total": 80
  }
]
```
