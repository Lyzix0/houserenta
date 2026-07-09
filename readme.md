# HouseRenta API

Backend для сервиса аренды квартир: регистрация/логин, объекты недвижимости, договоры аренды, показания счётчиков, автоматические счета, рынок свободных квартир и заявки на аренду.

- **Стек:** Go 1.26, [Fiber v3](https://github.com/gofiber/fiber), PostgreSQL, [Squirrel](https://github.com/Masterminds/squirrel) (SQL-builder), [swaggo](https://github.com/swaggo/swag) (Swagger/OpenAPI из аннотаций в коде), `golang-migrate` для миграций.
- **Аутентификация:** сессия на куке (`session_id`), без JWT.
- **Swagger UI:** `/swagger/*` на запущенном сервере; спека генерируется командой `make swagger-gen` в `docs/`.

## Быстрый старт

### Переменные окружения (`.env`)

| Переменная | Обязательна | Описание |
|---|---|---|
| `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB` | да | доступ к БД |
| `POSTGRES_HOST` | да (кроме докера, там `app-env-postgres`) | хост БД |
| `POSTGRES_PORT` | нет (`5432`) | порт БД |
| `POSTGRES_TIMEOUT` | нет (`10s`) | таймаут подключения |
| `PG_POOL_MAX` | нет (`10`) | размер пула соединений |
| `LOGGER_LEVEL` | нет (`DEBUG`) | уровень логов |
| `LOGGER_FOLDER` | да | папка для файлов логов |
| `HTTP_ADDR` | нет (`:5050`) | адрес HTTP-сервера |
| `ALLOWED_ORIGINS` | да, непусто и без `*` | список origin для CORS через запятую |

`ALLOWED_ORIGINS` не может быть пустым или содержать `*` — сессии идут через cookie (`AllowCredentials: true`), а это по спеке CORS несовместимо с wildcard-origin; при нарушении сервер падает с понятным сообщением при старте, а не глубоко в недрах `fiber/cors`.

### Запуск локально

```bash
make migrate-up      # применить миграции к БД
go run cmd/app/main.go
```

### Запуск через Docker

```bash
make app-deploy       # docker compose up -d --build renta-app
```

### Тесты

```bash
go test ./cmd/... ./internal/... ./pkg/...
```

## Модель данных

Все таблицы — в схеме `app`, первичные ключи — текстовые (UUID/сгенерированные ID), даты — ISO-строки (`TEXT`), не `TIMESTAMP`.

| Таблица | Смысл | Ключевые связи |
|---|---|---|
| `users` | пользователи: `landlord` / `tenant` / `admin` | — |
| `properties` | объекты недвижимости | `landlord_id → users.id` |
| `leases` | договор аренды (1 активный на объект) | `property_id UNIQUE → properties.id`, `tenant_user_id → users.id` |
| `readings` | показания счётчиков | `property_id → properties.id` |
| `bills` | счета (`type='rent'` для авто-счетов) | `property_id → properties.id` |
| `bill_items` | строки начислений в счёте | `bill_id → bills.id` |
| `custom_next_items` | разовые начисления в очереди на следующий счёт | `property_id → properties.id` |
| `applications` | отклики жильцов на свободные квартиры | `UNIQUE(property_id, tenant_user_id)` |

Договор аренды (`leases`) хранит только **активные** договоры: заселение делает `UPSERT` по `property_id` (замещая прежний договор с сохранением его `id`), выселение — `DELETE`. Поэтому «свободна ли квартира» = «есть ли для неё строка в `leases`».

## Роли и доступ

Роль хранится в сессии при логине/регистрации и проверяется мидлварой `RoleRequired`. Три роли: `landlord`, `tenant`, `admin` (регистрация как `admin` через публичный `/auth/register` запрещена).

Для эндпоинтов, завязанных на конкретный объект недвижимости (показания, оплата), доступ проверяется так: либо вызывающий — **владелец** объекта (`landlord_id` совпадает), либо **текущий жилец** по активному договору. Любой другой случай маскируется под `404 property not found`, а не `403` — чтобы нельзя было даже узнать, существует ли чужой объект.

## Автоматический биллинг (без cron)

Планировщика ОС нет: проверка «не пора ли выставить счёт» запускается «лениво» при каждом вызове `GET /auth/me` и `GET /properties` (см. `internal/usecase/billing`). Ошибка биллинга логируется и не влияет на ответ запроса — он не должен ломать навигацию по приложению.

Алгоритм для каждого активного договора:
1. Ищется последний счёт `type='rent'` для объекта. Если его нет или прошло ≥30 дней с даты последнего — генерируется новый.
2. Базовая строка счёта = `lease.price`.
3. Все записи `custom_next_items` для объекта прикрепляются как строки счёта и удаляются из очереди.
4. ЖКХ: берётся самое старое неучтённое показание (`is_accounted=0`) и самое свежее учтённое (`is_accounted=1`, или `0` при отсутствии) как база. Разница по `gvs/hvs/el1/el2`, умноженная на тарифы объекта, добавляется строками счёта; использованное показание помечается `is_accounted=1`.
5. `properties.balance -= total` (может уйти в минус — это долг до оплаты через `/pay`).

Запись счёта с его строками — единственное место в проекте с явной SQL-транзакцией (`BillRepo.Store`); остальные шаги (отметка показания, очистка очереди, списание баланса) идут последовательно без общей транзакции.

## Валидация

Значения, уходящие в числовые колонки Postgres (`REAL`/`INTEGER`), ограничены `validate`-тегами на уровне DTO, чтобы переполнение ловилось как `400 invalid body`, а не как необработанная ошибка БД (`500`):

- тарифы объекта — `max=100000`, суммы (аренда/платёж/разовое начисление) — `max=100000000`, показания счётчиков — `max=10000000`
- срок аренды — `max=1200` месяцев
- текстовые поля — разумные лимиты длины (100–500 символов), пароль — `max=72` (реальный лимит bcrypt)

## API

Базовый путь: `/v1`. Формат ошибки везде одинаковый: `{"error": "..."}`.

### Аутентификация (`/auth`)

| Метод | Путь | Доступ | Описание |
|---|---|---|---|
| POST | `/auth/register` | публично | регистрация + автологин (ставит cookie) |
| POST | `/auth/login` | публично | вход по email **или** телефону |
| POST | `/auth/logout` | публично | завершение сессии |
| GET | `/auth/me` | Protected | профиль текущей сессии + запуск авто-биллинга |
| POST | `/auth/profile` | Protected | частичное обновление профиля (только переданные поля) |

**POST /auth/register**
```json
// запрос
{
  "name": "Иванов Иван Иванович",
  "email": "ivanov@example.com",
  "password": "securepassword123",
  "phone": "+79991112233",
  "initialRole": "tenant",       // landlord | tenant
  "document": "Паспорт РФ 4512 № 345678",  // необязательно
  "paymentCard": "4276111122223333"        // необязательно, ровно 16 цифр
}
```
Обязательны: `name`, `email`, `password`, `phone`, `initialRole`. Ответ `200 OK`:
```json
{ "id": "user-a9b8c7d", "name": "Иванов Иван Иванович", "email": "ivanov@example.com", "role": "tenant" }
```
`400` при повторной почте: `{"error": "Пользователь с такой почтой уже зарегистрирован"}`.

**POST /auth/login** — `{"email": "...", "password": "..."}` (в `email` можно передать и телефон). `400` при неверных данных: `{"error": "Неверная почта или пароль"}`.

**POST /auth/logout** — без тела, `200 {"ok": true}`.

**GET /auth/me** — `200`:
```json
{
  "id": "user-a9b8c7d", "name": "...", "email": "...", "role": "tenant",
  "document": "...", "phone": "...",
  "paymentCard": "4276111122223333",   // или null
  "tenantPropertyId": "prop-z8y7x6w"   // или null, только для tenant
}
```

**POST /auth/profile** — тело как в register, но все поля опциональны и меняется только переданное:
```json
{ "name": "Иван Колесников", "phone": "+79997776655", "email": "newemail@example.com", "password": "newSuperSecurePassword" }
```
`200 {"ok": true}`; `400` та же ошибка про занятую почту.

### Объекты недвижимости (`/properties`)

| Метод | Путь | Доступ | Описание |
|---|---|---|---|
| GET | `/properties` | Protected | список своих объектов (landlord) / своей аренды (tenant) + авто-биллинг |
| GET | `/properties/vacant` | Protected | рынок: все свободные объекты по всей системе |
| POST | `/properties` | Landlord | создать объект |
| GET | `/properties/:id` | Landlord | один объект (только свой; сверяется `landlord_id`, а не owner/tenant — жилец получит `404` даже на объект, который снимает) |
| PUT | `/properties/:id` | Landlord | редактировать (только свой) |
| DELETE | `/properties/:id` | Landlord | удалить (каскадно удаляет договор аренды на объекте) |
| POST | `/properties/:id/lease` | Landlord | заселить жильца |
| DELETE | `/properties/:id/lease` | Landlord | выселить |
| POST | `/properties/:id/readings` | Protected (owner/tenant) | передать показания |
| POST | `/properties/:id/pay` | Protected (owner/tenant) | оплата счёта или пополнение баланса |
| POST | `/properties/:id/custom-item` | Landlord | разовое начисление в очередь |
| POST | `/properties/:id/apply` | Tenant | отклик на свободный объект |

**GET /properties** / **GET /properties/vacant** — массив объектов с автоматической вложенной сборкой:
```json
[{
  "id": "prop-z8y7x6w", "landlord_id": "user-landlord123", "name": "Квартира с видом на парк",
  "coordinates": "55.7558, 37.6173", "country": "Россия", "region": "Москва", "city": "Москва",
  "street": "Лобовая", "house": "14", "apartment": "205",
  "gvs_tariff": 220.5, "hvs_tariff": 50.2, "el1_tariff": 6.15, "el2_tariff": null, "balance": 15000,
  "readings": [{ "id": "read-123", "property_id": "prop-z8y7x6w", "date": "...", "gvs": 12.5, "hvs": 24.1, "el1": 340, "el2": null, "is_accounted": 1 }],
  "bills": [{ "id": "bill-999", "property_id": "prop-z8y7x6w", "date": "...", "due_date": "...", "status": "unpaid", "type": "rent", "total": 35000,
    "items": [{ "id": "item-1", "bill_id": "bill-999", "description": "Аренда за период", "amount": 35000 }] }],
  "customNextItems": [],
  "tenant": { "id": "lease-555", "tenant_user_id": "user-a9b8c7d", "name": "...", "document": "...", "phone": "...",
    "months_of_rent": 12, "price": 35000, "payment_day": 5, "reading_day": 25, "start_date": "...", "end_date": "..." },
  "landlordName": "Алексей Петров", "landlordPhone": "+79001234567",
  "applications": []
}]
```
- `landlord`: возвращаются **все** его объекты (свободные и занятые).
- `tenant`: возвращается только объект по его текущему договору (или пустой список).
- `/properties/vacant`: все объекты без активного договора, по всем владельцам, вместе с откликами (`applications`) на них.
- `tenant` = `null` для свободного объекта; `applications` не пустой только у свободных объектов (при заселении все отклики на объект автоматически удаляются).

**POST /properties** (Landlord; `landlord_id` берётся из сессии, не из тела):
```json
{ "name": "Студия у метро", "coordinates": "55.79, 37.60", "country": "Россия", "region": "Москва", "city": "Москва",
  "street": "Проспект Мира", "house": "45", "apartment": "12", "gvsTariff": 240, "hvsTariff": 60, "el1Tariff": 6.5, "el2Tariff": 2.5 }
```
Ответ `200`: `{ "id": "prop-k8r3f9d", "name": "Студия у метро" }`.

**PUT /properties/:id** — то же тело, что при создании. `200 {"ok": true}`.

**DELETE /properties/:id** — `200 {"ok": true}`; договор аренды на объекте удаляется каскадом на уровне БД.

**POST /properties/:id/lease** (заселение):
```json
{ "tenantUserId": "user-a9b8c7d", "price": 32000, "monthsOfRent": 11, "paymentDay": 10, "readingDay": 28 }
```
`name/document/phone` в договоре берутся из профиля жильца, не из запроса. `startDate` = сейчас, `endDate` = `startDate + monthsOfRent`. Если у объекта уже был договор — новый его замещает (тот же `id` договора). Все отклики (`applications`) на объект удаляются. `404 {"error": "Жилец не найден в системе"}`, если `tenantUserId` не существует или не имеет роль `tenant`.

**DELETE /properties/:id/lease** (выселение) — `200 {"ok": true}`, идемпотентно (не ошибка, если договора и не было).

**POST /properties/:id/readings**:
```json
{ "gvs": 15.2, "hvs": 29.8, "el1": 395.5, "el2": 110.2 }
```
`el2` необязателен. `200 {"ok": true}`.

**POST /properties/:id/pay**:
```json
{ "amount": 35000, "billId": "bill-999" }
```
Если `billId` передан — счёт помечается `paid` (без сверки суммы с `amount`). Если не передан — `amount` добавляется к `balance` объекта. `404 {"error": "bill not found"}`, если счёт не принадлежит этому объекту.

**POST /properties/:id/custom-item** (Landlord):
```json
{ "description": "Замена смесителя на кухне", "amount": 2500 }
```
`200 {"ok": true}`; попадёт в следующий авто-счёт.

**POST /properties/:id/apply** (Tenant) — без тела. `200 {"ok": true}`; `400 {"error": "Вы уже откликнулись на это предложение"}` при повторном отклике на тот же объект.

### Жильцы (`/tenants`)

| Метод | Путь | Доступ | Описание |
|---|---|---|---|
| GET | `/tenants/unlinked` | Landlord | жильцы без активного договора аренды |

Ответ — массив `{ id, name, email, document, phone }`, для выпадающего списка при ручном заселении.

## Известные ограничения

- **Заявки/тикеты за пределами `applications`.** Раздел «заявки» покрывает только отклики на свободные квартиры (`/properties/:id/apply`); другой заявочной подсистемы (тикеты в поддержку и т.п.) в проекте нет.
- **Нет атомарности между шагами авто-биллинга**, кроме самой записи счёта — если что-то упадёт после успешного `bills.Store`, счёт останется корректным, но отметка показания/очистка очереди/списание баланса могут не выполниться до следующего запуска.
- **Оплата счёта не сверяет `amount` с `total`** — `billId` просто помечает счёт оплаченным независимо от переданной суммы.
