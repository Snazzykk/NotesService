# NoteService

📓 NotesService
Описание

NotesService — это REST API для управления заметками с поддержкой JWT-аутентификации.
Позволяет пользователям создавать, получать, обновлять и удалять заметки.

## Особенности:

JWT-аутентификация

CRUD заметок

PostgreSQL база данных

Swagger UI для документации API

Docker-контейнеризация

## 🔗 Документация API

Swagger UI доступен по адресу:

http://localhost:8083/docs

## Запуск проекта с Docker



Создай файл .env.docker с переменными окружения для Docker:
```
# PostgreSQL (для Docker)
POSTGRES_USER=postgres
POSTGRES_PASSWORD=123
POSTGRES_DB=notes_service
DB_HOST=postgres
DB_PORT=5432
DB_SSLMODE=disable

# HTTP сервер
HTTP_ADDRESS=:8083
HTTP_TIMEOUT=4s
HTTP_IDLE_TIMEOUT=60s
HTTP_USER=user
HTTP_PASSWORD=user

# JWT
JWT_SECRET=xK9pL2mN7vB5cR8tQ3wZ1yA4sD6hJ0f

# Среда для Docker
ENV=local
CONFIG_PATH=./config/config.yaml

```

### Запусти контейнеры:

docker compose --env-file .env.docker up --build

docker compose down -v      
docker compose up --build 


После успешного запуска сервис будет доступен:
```
API: http://localhost:8083
Swagger: http://localhost:8083/docs
PostgreSQL: localhost:5430 (порт на хосте)
```