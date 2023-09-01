# Note service

Микросервис, предназначенный для работы с заметками пользователей.
Он предоставляет API для создания, чтения, обновления и удаления заметок, а также для аутентификации
и регистрации пользователей. Особенностью этого сервиса является интеграция с сервисом Яндекс.Спеллер
для автоматической проверки орфографических ошибок в заметках перед их сохранением и изменением.

Используемые технологии:
- PostgreSQL (в качестве хранилища данных)
- Docker (для запуска сервиса)
- Swagger (для документации API)
- Chi (веб фреймворк)
- Redis (для кеширования)
- uber-go/zap (для логирования)
- golang-migrate/migrate (для миграций БД)
- pgx (драйвер для работы с PostgreSQL)
- golang/mock, testify (для тестирования)

Сервис разработан с использованием современных технологий и следует принципам Clean Architecture, 
что обеспечивает легкость расширения функционала и тестирования. Также был реализован Graceful Shutdown
для корректного завершения работы сервиса.

# Getting Started

Для запуска сервиса достаточно заполнить файл .env, расположенный в корневой директории.

# Usage

Запустить сервис можно с помощью команды `make compose-up`

Документацию после завпуска сервиса можно посмотреть по адресу `http://localhost:8080/swagger/index.html`
с портом 8080 по умолчанию

Для запуска тестов необходимо выполнить команду `make test`, для запуска тестов с покрытием `make cover` и `make cover-html` для получения отчёта в html формате

Для запуска линтера необходимо выполнить команду `make linter-golangci`

## Examples

Некоторые примеры запросов
- [Регистрация](#sign-up)
- [Аутентификация](#login)
- [Создание заметки](#create)
- [Получение заметки](#get)
- [Получение всех заметок](#get-all)
- [Обновление заметки](#update)
- [Удаление заметки](#delete)

### Регистрация <a name="sign-up"></a>

Регистрация пользователя:
```curl
curl --location --request POST 'http://localhost:8080/api/auth/sign-up' \
--header 'Content-Type: application/json' \
--data-raw '{
    "login":"exampe@mail.ru",
    "password":"Qwerty123!"
}'
```
Пример ответа:
```json
{
   "status": "success"
}
```

### Аутентификация <a name="login"></a>

Аутентификация пользователя для получения токена доступа:
```curl
curl --location --request POST 'http://localhost:8080/api/auth/login' \
--header 'Content-Type: application/json' \
--data-raw '{
    "login":"exampe@mail.ru",
    "password":"Qwerty123!"
}'
```
Пример ответа:
```json
{
   "status": "success",
   "data": {
      "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiJldVhzNzRiUnNqNHZKS2RWck4vT0tpOWxyRjZpT3NPZFNwMDNVT0U9In0.GvsuCvrPxq7EbGZE1zHMgvMZiKZymo6FF7m6xt-zIXw"
   }
}
```

### Создание заметки <a name="create"></a>

Создание заметки пользователем:
```curl
curl --location --request POST 'http://localhost:8080/api/note/create' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiJldVhzNzRiUnNqNHZKS2RWck4vT0tpOWxyRjZpT3NPZFNwMDNVT0U9In0.GvsuCvrPxq7EbGZE1zHMgvMZiKZymo6FF7m6xt-zIXw' \
--header 'Content-Type: application/json' \
--data-raw '{
    "title": "Заголовок заметки 1",
    "content": "Содержание заметки 1"
}'
```
Пример ответа:
```json
{
   "status": "success"
}
```

### Получение заметки <a name="get"></a>

Получение одной заметки пользователем:
```curl
curl --location --request GET 'http://localhost:8080/api/note/get?noteId=1' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiJldVhzNzRiUnNqNHZKS2RWck4vT0tpOWxyRjZpT3NPZFNwMDNVT0U9In0.GvsuCvrPxq7EbGZE1zHMgvMZiKZymo6FF7m6xt-zIXw'
```
Пример ответа:
```json
{
   "status": "success",
   "data": {
      "id": 1,
      "title": "Заголовок заметки 1",
      "content": "Содержание заметки 1"
   }
}
```

### Получение всех заметок <a name="get-all"></a>

Получение всех заметок пользователя:
```curl
curl --location --request GET 'http://localhost:8080/api/note/get-all' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiJldVhzNzRiUnNqNHZKS2RWck4vT0tpOWxyRjZpT3NPZFNwMDNVT0U9In0.GvsuCvrPxq7EbGZE1zHMgvMZiKZymo6FF7m6xt-zIXw'
```
Пример ответа:
```json
{
   "status": "success",
   "data": [
      {
         "id": 1,
         "title": "Заголовок заметки 1",
         "content": "Содержание заметки 1"
      },
      {
         "id": 2,
         "title": "Заголовок заметки 2",
         "content": "Содержание заметки 2"
      }
   ]
}
```

### Обновление заметки <a name="update"></a>

Обновление заметки пользователя:
```curl
curl --location --request PUT 'http://localhost:8080/api/note/update?noteId=1' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiJldVhzNzRiUnNqNHZKS2RWck4vT0tpOWxyRjZpT3NPZFNwMDNVT0U9In0.GvsuCvrPxq7EbGZE1zHMgvMZiKZymo6FF7m6xt-zIXw' \
--header 'Content-Type: application/json' \
--data-raw '{
    "title": "Новый заголовок",
    "content": "Новое содержание"
}'
```
Пример ответа:
```json
{
   "status": "success"
}
```

### Удаление заметки <a name="delete"></a>

Удаление заметки пользователя:
```curl
curl --location --request DELETE 'http://localhost:8080/api/note/delete?noteId=1' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiJldVhzNzRiUnNqNHZKS2RWck4vT0tpOWxyRjZpT3NPZFNwMDNVT0U9In0.GvsuCvrPxq7EbGZE1zHMgvMZiKZymo6FF7m6xt-zIXw' \
--header 'Content-Type: application/json' \
--data-raw '{
    "title": "Новый заголовок",
    "content": "Новое содержание"
    }'
```
Пример ответа:
```json
{
   "status": "success"
}
```

