# Documents API

## Docker

Требуется только Docker и Docker Compose.

```bash
git clone
cd golang-documents

docker compose up -d --build
```

Проверить, что всё работает:

```bash
curl http://localhost:8000/health
```

## Локальный запуск

### 1. зависимости

- Go 1.26+
- PostgreSQL 16 (локально или через Docker)

### 2. переменные окружения

```bash
cp .env.example .env
```

### 3. go-зависимости

```bash
go mod download
```

### 4. запуск миграций

```bash
go run ./cmd/migrate
```

### 5. запуск сервера

```bash
go run ./cmd/server
```

## API

### Регистрация

```bash
curl -X POST http://localhost:8000/api/register \
  -H "Content-Type: application/json" \
  -d '{"token":"secret_admin_token","login":"testlogin1","password":"Passw0rd!"}'
```

```bash
curl -X POST http://localhost:8000/api/auth \
  -H "Content-Type: application/json" \
  -d '{"login":"testlogin1","password":"Passw0rd!"}'
```

#### Загрузка JSON-документа

```bash
curl -X POST http://localhost:8000/api/docs \
  -H "Authorization: Bearer <токен>" \
  -F 'meta={"name":"note.txt","file":false,"public":false,"mime":"application/json"}' \
  -F 'json={"hello":"world"}'
```

#### Загрузка файла

```bash
curl -X POST http://localhost:8000/api/docs \
  -H "Authorization: Bearer <токен>" \
  -F 'meta={"name":"photo.png","file":true,"public":true,"mime":"image/png"}' \
  -F 'file=@/path/to/photo.png'
```

#### Список документов

```bash
curl http://localhost:8000/api/docs -H "Authorization: Bearer <токен>"
```

### Получение документа по id

```bash
curl http://localhost:8000/api/docs/<id> -H "Authorization: Bearer <токен>"
```

### Удаление документа

```bash
curl -X DELETE http://localhost:8000/api/docs/<id> -H "Authorization: Bearer <токен>"
```

#### Logout

```bash
curl -X DELETE http://localhost:8000/api/auth/logout -H "Authorization: Bearer <токен>"
```

#### Каждый час отрабатывает `token_cleaning.go` горутина, которая чистит токены, у которых истёк `TOKEN_TTL`.
