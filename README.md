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

HTTP Request → Handler → Service → Repository → PostgreSQL

<h3>API</h3>

**Users**

| Method	| Endpoint |	Description |
| -------- | -------- | -------- |
| POST | /users | Create user |
| GET	| /users/{id}	| Get user |
| GET	| /users	| Get users |
| PATCH	| /users/{id}	| Update user|
| DELETE | /users/{id}	| Delete user |

**Tasks**

| Method	| Endpoint |	Description |
| -------- | -------- | -------- |
| POST | /tasks | Create task |
| GET	| /tasks/{id}	| Get task |
| GET	| /tasks	| Get tasks |
| PATCH	| /tasks/{id}	| Update task|
| DELETE | /tasks/{id}	| Delete task |


