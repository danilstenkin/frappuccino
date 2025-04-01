☕ Frappuccino Coffee Shop Backend

Frappuccino — это backend-система управления кофейней, написанная на Go с использованием net/http и PostgreSQL. Система позволяет управлять меню, заказами, инвентарем, клиентами и отчетами, а также обрабатывать групповые заказы с учетом остатков ингредиентов.
📦 Стек технологий

    Язык: Go (golang)

    База данных: PostgreSQL

    Web-сервер: net/http (без сторонних фреймворков)

    Архитектура: MVC-подход (Handlers, Repositories, Models)

    Docker: для контейнеризации и инициализации БД

📁 Структура проекта

frappuccino/
├── db/                    // Подключение к PostgreSQL
├── handlers/              // HTTP обработчики (контроллеры)
├── models/                // Структуры моделей
├── repositories/          // SQL-логика
├── router/                // Роутинг HTTP путей
├── init.sql               // SQL-инициализация БД
├── Dockerfile             // Docker образ приложения
├── docker-compose.yml     // Конфигурация для запуска PostgreSQL и API
├── db_diagram.png         // Диаграмма БД.
└── main.go                // Точка входа

🚀 Запуск проекта

    Убедитесь, что установлен Docker и Docker Compose.

    Клонируйте проект:

git clone https://github.com/yourusername/frappuccino.git
cd frappuccino

    Запустите проект:

docker-compose up --build

API будет доступен по адресу: http://localhost:8080
📚 Основные эндпоинты
📋 Меню

    POST /menu — создать блюдо

    GET /menu — получить список блюд

    GET /menu/{id} — получить конкретное блюдо

    PUT /menu/{id} — обновить блюдо

    DELETE /menu/{id} — удалить блюдо

🛒 Заказы

    POST /orders — создать заказ

    GET /orders — получить все заказы

    GET /orders/{id} — получить заказ по ID

    PUT /orders/{id} — обновить заказ

🧾 Order Items

    POST /order-items — добавить блюдо в заказ (с проверкой остатков и списанием)

    DELETE /order-items/{id} — удалить блюдо из заказа

🧍 Клиенты

    POST /customers — создать клиента

    GET /customers — получить всех клиентов

📊 Отчёты

    GET /numberOfOrderedItems?startDate=DD.MM.YYYY&endDate=DD.MM.YYYY — количество заказанных блюд за период

    GET /reports/search?q=keyword&filter=menu,orders — полнотекстовый поиск по заказам и меню

    GET /reports/orderedItemsByPeriod?period=day|month&month=march&year=2025 — количество заказов по дням или месяцам

    GET /inventory/getLeftOvers?sortBy=quantity&page=1&pageSize=5 — остатки инвентаря с пагинацией

🔁 Массовая обработка заказов

    POST /orders/batch-process — обработка нескольких заказов одновременно с транзакционностью и проверкой остатков