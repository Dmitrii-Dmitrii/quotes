# Quotes

## Описание 
REST API-сервис на Go для хранения и управления цитатами.

## Инструкция по сборке
Для запуска сервера необходимо склонировать данный репозиторий и перейти в его директорию:
```
git https://github.com/Dmitrii-Dmitrii/quotes.git
cd quotes
```
Далее нужно перейти в директорию `/docker`. Внутри нее находится `docker-compose.yml`, который отвечает за развертывание базы данных. Его нужно запустить:
```
cd docker
docker-compose up -d
```
Далее нужно применить миграции из корневой директории:
```
goose -dir ./migrations postgres "<CONNECTION_STRING>" up
```
Переменная `<CONNECTION_STRING>` зависит от данных в `docker-compose.yml`, если ничего не менять, то нужно использовать `postgres://quotes_user:quotes_password@localhost:5432/quotes_database?sslmode=disable`.

После важно создать `.env` файл **в корне проекта**:
```
touch .env
```
Пример `.env` файла:
```
SERVER_PORT=8080
DB_CONNECTION_STRING=postgres://quotes_user:quotes_password@localhost:5432/quotes_database?sslmode=disable
```
Далее необходимо запустить сервер из корневой директории:
```
go run cmd/server/main.go
```

## API
1. Добавление новой цитаты (POST /quotes)
2. Получение всех цитат (GET /quotes)
3. Получение случайной цитаты (GET /quotes/random)
4. Фильтрация по автору (GET /quotes?author=Confucius)
5. Удаление цитаты по ID (DELETE /quotes/{id})