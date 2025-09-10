# Folder System 🗂️
Система управления документами и папками с ограничением емкости папок и JWT-аутентификацией.

# О проекте 📋
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

# 📦 Установка и запуск

Клонируйте репозиторий
git clone https://github.com/Khammatoff/folder-system
cd folder-system
Настройте окружение
->
Отредактируйте .env файл
Запустите приложение
->
docker-compose up --build
Приложение будет доступно по адресу: http://localhost:8080

# ⚙️ Конфигурация
Файл .env

Database
DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=folder_system
DB_SSLMODE=disable

JWT
JWT_ACCESS_SECRET=your_secret_access_key
JWT_REFRESH_SECRET=your_secret_refresh_key
JWT_ACCESS_TTL=15
JWT_REFRESH_TTL=10080

Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

Logging
LOG_LEVEL=debug
LOG_FILE=app.log

# Примеры curl-запросов
Регистрация
curl -X POST http://localhost:8080/api/register \
-H "Content-Type: application/json" \
-d '{"email":"test@example.com", "password":"123456"}'

Логин
curl -X POST http://localhost:8080/api/login \
-H "Content-Type: application/json" \
-d '{"email":"test@example.com", "password":"123456"}'

Создать документ без папки
curl -X POST http://localhost:8080/api/protected/documents/ \
-H "Content-Type: application/json" \
-H "Authorization: Bearer <toker>" \
-d '{
  "title": "Document without folder",
  "sheets_count": 5,
  "document_type_id": 1
}'


# 🐛 Логирование
Все действия и ошибки логируются в файл app.log с указанием:
Времени события
Уровня логирования (debug, info, error)
Детальной информации о запросе
Stack trace ошибок
