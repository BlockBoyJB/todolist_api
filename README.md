# todolist_api

REST API для управления задачами (To-Do List).   
Помимо основного тестового задания (создание, получение, обновление и удаление задач) была реализована авторизация и аутентификация. 
Каждый пользователь имеет доступ только к своим задачам


### Используемый стек

* **Golang 1.22**
* **Echo** основной веб фреймворк
* **PostgreSQL** как основная БД
* **golang-jwt/jwt** для jwt
* **swaggo/swag** swagger документация API
* **golang-migrate/migrate** для миграций бд
* **logrus** для логирования
* **testify** для тестирования
* **Docker и Docker Compose** для быстрого развертывания


### Prerequisites
- Docker, Docker Compose installed

### Getting started

* Добавить репозиторий к себе
* Создать .env файл в директории с проектом и заполнить информацией из .env.example

### Usage

Запустить сервис можно с помощью `make compose-up` (или `docker-compose up -d --build`)

Документация доступна по адресу `http://localhost:8080/swagger/`

Запуск тестов доступен с помощью команды `make tests`


### Примеры запросов

* [Регистрация](#регистрация)
* [Аутентификация](#вход-аутентификация)
* [Создание задачи](#создание-задачи)
* [Получение списка задач](#получение-списка-задач)
* [Получение задачи по id](#получение-задачи-по-id)
* [Обновление задачи по id](#обновление-задачи-по-id)
* [Удаление задачи по id](#удаление-задачи-по-id)


#### Регистрация
```shell
curl -X 'POST' \
      'http://localhost:8080/auth/sign-up' \
      -H 'accept: application/json' \
      -H 'Content-Type: application/json' \
      -d '{"username": "maks", "password": "abc"}'
```
Пример ответа:  
`201`


#### Вход, аутентификация
```shell
curl -X 'POST' \
      'http://localhost:8080/auth/sign-in' \
      -H 'accept: application/json' \
      -H 'Content-Type: application/json' \
      -d '{"username": "maks", "password": "abc"}'
```

Пример ответа:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjUxOTgzMzAsImlhdCI6MTcyNDkzOTEzMCwidXNlcm5hbWUiOiJtYWtzIn0.moGyejLn8y4iUEjMr4Q1IgwDeAsyr4Q0rIIO3iycxaA"
}
```


#### Создание задачи
```shell
curl -X 'POST' \
      'http://localhost:8080/api/v1/tasks' \
      -H 'accept: application/json' \
      -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjUxOTgzMzAsImlhdCI6MTcyNDkzOTEzMCwidXNlcm5hbWUiOiJtYWtzIn0.moGyejLn8y4iUEjMr4Q1IgwDeAsyr4Q0rIIO3iycxaA' \
      -H 'Content-Type: application/json' \
      -d '{"title": "foobar", "description": "test desc", "due_date", "2024-01-01T12:00:00Z"}'
```

Пример ответа
```json
{
  "id": 1,
  "title": "foobar",
  "description": "test desc",
  "due_date": "2024-01-01T12:00:00Z",
  "created_at": "2024-08-29T15:47:17Z",
  "updated_at": "2024-08-29T15:47:17Z"
}
```


#### Получение списка задач
```shell
curl -X 'GET' \
      'http://localhost:8080/api/v1/tasks' \
      -H 'accept: application/json' \
      -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjUxOTgzMzAsImlhdCI6MTcyNDkzOTEzMCwidXNlcm5hbWUiOiJtYWtzIn0.moGyejLn8y4iUEjMr4Q1IgwDeAsyr4Q0rIIO3iycxaA'
```

Пример ответа:
```json
[
  {
    "id": 1,
    "title": "foobar",
    "description": "test desc",
    "due_date": "2024-01-01T12:00:00Z",
    "created_at": "2024-08-29T15:47:17Z",
    "updated_at": "2024-08-29T15:47:17Z"
  }
]
```


#### Получение задачи по id
```shell
curl -X 'GET' \
      'http://localhost:8080/api/v1/tasks/1' \
      -H 'accept: application/json' \
      -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjUxOTgzMzAsImlhdCI6MTcyNDkzOTEzMCwidXNlcm5hbWUiOiJtYWtzIn0.moGyejLn8y4iUEjMr4Q1IgwDeAsyr4Q0rIIO3iycxaA'
```

Пример ответа:
```json
{
  "id": 1,
  "title": "foobar",
  "description": "test desc",
  "due_date": "2024-01-01T12:00:00Z",
  "created_at": "2024-08-29T15:47:17Z",
  "updated_at": "2024-08-29T15:47:17Z"
}
```


#### Обновление задачи по id
```shell
curl -X 'PUT' \
  'http://localhost:8080/api/v1/tasks/1' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjUxOTgzMzAsImlhdCI6MTcyNDkzOTEzMCwidXNlcm5hbWUiOiJtYWtzIn0.moGyejLn8y4iUEjMr4Q1IgwDeAsyr4Q0rIIO3iycxaA' \
  -H 'Content-Type: application/json' \
  -d '{"title": "new title", "description": "new test desc", "due_date": "2024-01-01T15:00:00Z"}'
```

Пример ответа:
```json
{
  "id": 1,
  "title": "new title",
  "description": "new test desc",
  "due_date": "2024-01-01T15:00:00Z",
  "created_at": "2024-08-29T15:47:17Z",
  "updated_at": "2024-08-29T15:57:06Z"
}
```


#### Удаление задачи по id
```shell
curl -X 'DELETE' \
  'http://localhost:8080/api/v1/tasks/1' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjUxOTgzMzAsImlhdCI6MTcyNDkzOTEzMCwidXNlcm5hbWUiOiJtYWtzIn0.moGyejLn8y4iUEjMr4Q1IgwDeAsyr4Q0rIIO3iycxaA'
```

Пример ответа:  
`204`


### Тестовое задание
Разработать REST API для системы управления задачами, которая позволяет пользователям создавать, просматривать, обновлять и удалять задачи.