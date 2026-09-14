<h1>Backend Task Management API</h1>

<h2>Backend-приложение на Go для управления пользователями и задачами.</h2>

Проект реализует CRUD для пользователей и задач, а также получение статистики по задачам. В проекте используются многослойная архитектура, PostgreSQL, HTTP middleware, валидация, обработка ошибок.

<h3>Стек</h3>

- Go
- PostgreSQL
- SQL
- zap
- net/http
- pgx
- golang-migrate

<h3>Архитектура</h3>
Проект построен с использованием многослойной архитектуры:

HTTP Request → Handler → Service → Repository → PostgreSQL

<h3>API</h3>

**Users**

| Метод	| Эндпоинт |	Описание |
| -------- | -------- | -------- |
| POST | /users | Создание пользователя |
| GET	| /users/{id}	| Получение пользователя |
| GET	| /users	| Получение списка пользователей |
| PATCH	| /users/{id}	| Обновление пользователя |
| DELETE | /users/{id}	| Удаление пользователя |

**Tasks**

| Метод	| Эндпоинт |	Описание |
| -------- | -------- | -------- |
| POST | /tasks | Создание задачи |
| GET	| /tasks/{id}	| Получение задачи |
| GET	| /tasks	| Получение списка задач |
| PATCH	| /tasks/{id}	| Обновление задачи |
| DELETE | /tasks/{id}	| Удаление задачи |

**Statistics**

| Метод	| Эндпоинт |	Описание |
| -------- | -------- | -------- |
| GET	| /statistics	| Получение статистики по задачам |

<h3>Примеры запросов</h3>

Создание пользователя
```bash
curl -X POST http://localhost:5050/users \
  -H "Content-Type: application/json" \
  -d '{"full_name":"Иван Иванов","phone_number":"+79991234567"}'
```

Получение пользователя
```bash
curl http://localhost:5050/users/1
```

Получение списка пользователей
```bash
curl "http://localhost:5050/users?limit=10&offset=0"
```

Обновление пользователя
```bash
curl -X PATCH http://localhost:5050/users/1 \
  -H "Content-Type: application/json" \
  -d '{"full_name":"Пётр Петров"}'
```

Удаление пользователя
```bash
curl -X DELETE http://localhost:5050/users/1
```

Создание задачи
```bash
curl -X POST http://localhost:5050/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"title","description":"description","author_user_id":1}'
```

Получение задачи
```bash
curl http://localhost:5050/tasks/1
```

Получение списка задач
```bash
curl "http://localhost:5050/tasks?limit=10&offset=0"
```

Обновление задачи
```bash
curl -X PATCH http://localhost:5050/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"myTitle"}'
```

Удаление задачи
```bash
curl -X DELETE http://localhost:5050/tasks/1
```

Получение статистики
```bash
curl http://localhost:5050/statistics
```
