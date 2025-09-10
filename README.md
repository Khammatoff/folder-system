Folder System 🗂️
Система управления документами и папками с ограничением емкости папок и JWT-аутентификацией.

📋 О проекте
Система предоставляет REST API для работы с документами и папками. Каждая папка имеет ограничение в 480 листов. Документы могут занимать различное количество листов и привязываться к папкам с проверкой доступного места.

🚀 Возможности
🔐 Аутентификация
Регистрация и авторизация пользователей

JWT токены (access + refresh)

Защищенные роуты

📁 Управление документами
Создание документов с указанием папки

Обновление документов (название, количество листов, папка)

Удаление документов

Автоматический учет занятого места

📂 Управление папками
Рекомендация подходящей папки для документа

Проверка свободного места

Автоматическое освобождение места при перемещении документов

🛠️ Технологии
Backend: Go 1.21+

Database: PostgreSQL 15+

Router: Chi Router

ORM: GORM

Authentication: JWT

Frontend: Vanilla JavaScript

Containerization: Docker + Docker Compose

📦 Установка и запуск
Требования
Docker и Docker Compose

Go 1.21+ (для разработки)

Быстрый запуск
Клонируйте репозиторий

bash
git clone <repository-url>
cd folder-system
Настройте окружение

bash
cp .env.example .env
# Отредактируйте .env файл при необходимости
Запустите приложение

bash
docker-compose up --build
Приложение будет доступно по адресу: http://localhost:8080

⚙️ Конфигурация
Файл .env
env
# Database
DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=folder_system
DB_SSLMODE=disable

# JWT
JWT_ACCESS_SECRET=your_secret_access_key
JWT_REFRESH_SECRET=your_secret_refresh_key
JWT_ACCESS_TTL=15
JWT_REFRESH_TTL=10080

# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# Logging
LOG_LEVEL=debug
LOG_FILE=app.log
📡 API Endpoints
🔓 Публичные роуты
Регистрация

http
POST /api/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
Авторизация

http
POST /api/login
Content-Type: application/json

{
  "email": "user@example.com", 
  "password": "password123"
}
🔐 Защищенные роуты (требуют JWT)
Документы

http
GET    /api/documents/{id}
POST   /api/documents
PUT    /api/documents/{id}
DELETE /api/documents/{id}
Рекомендация папки

http
GET /api/folders/recommended?document_type_id=1&sheets_count=100
🗃️ Структура базы данных
Основные таблицы
users - пользователи системы

folders - папки (емкость 480 листов)

documents - документы

document_types - типы документов

folder_types - типы папок

folder_type_assignments - связи типов документов и папок

🧪 Тестирование
Через веб-интерфейс
Откройте http://localhost:8080

Зарегистрируйте пользователя

Выполните вход

Тестируйте создание документов и папок

Через curl
bash
# Регистрация
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Авторизация
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
🐛 Логирование
Все действия и ошибки логируются в файл app.log с указанием:

Времени события

Уровня логирования (debug, info, error)

Детальной информации о запросе

Stack trace ошибок

🔧 Разработка
Локальная разработка
bash
# Установите зависимости
go mod download

# Запустите PostgreSQL
docker-compose up db

# Запустите приложение
go run ./cmd/server/main.go
Сборка
bash
# Сборка Docker образа
docker-compose build

# Запуск тестов
go test ./...