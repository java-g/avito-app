# Avito App

К сожалению, времени хватило только на код ☺ ***

## Запуск

Установка зависимостей

```bash
go mod init
```
Запуск любым стандартным способом

```bash
go run main.go
```
или

```bash
go build main.go
```

Не забудьте поднять PostgreSQL

## .ENV

```environment
DB_HOST="localhost"
DB_PORT=5432
DB_USER="postgres"
DB_PASSWORD="postgres"
DB_NAME="jwt-auth-api"
DB_SSLMODE="disable"
DB_TIMEZONE="Europe/Moscow"
JWT_KEY="JSON Web Token Secret Key"
```

## API

Описанный API был расширен методами Регистрации и Получения токена

### Регистрация пользователя


```bash
curl -X 'POST' \
  'http://localhost:8000/register' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{"login": "user login", "password": "user password"}'
```

### Получение токена


```bash
curl -X 'POST' \
  'http://localhost:8000/login' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{"login": "user login", "password": "user password"}'
```
