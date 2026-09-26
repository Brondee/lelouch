# Парсер одежды: нефункциональные требования MVP

## Аудитория
- Сейчас — один пользователь.
- Задел на несколько пользователей: пользовательские данные, включая правила поиска и уведомления, связаны с `user_id`. Жёсткой привязки к единственному пользователю нет.

## Нагрузка
- Основная нагрузка — исходящие запросы парсеров к площадкам; пользовательский трафик мал.
- Запросы распределяются равномерно во времени через rate limiting и ограничение concurrency для каждой площадки.
- Лимиты учитывают ограничения площадок; обход защит не входит в требования.

## Объём данных и БД
- Объём чтения и записи в MVP небольшой.
- Специальная оптимизация БД под read/write-нагрузку пока не требуется.

## Задержка
- Цель в штатном режиме — обнаружение объявления не позднее 5 минут после публикации на площадке.
- Желаемый ориентир — 2–3 минуты.

## Доступность
- Best-effort: простой допустим, строгий SLA не требуется.
- Остановка системы должна быстро обнаруживаться с уведомлением владельца; конкретный порог времени пока не определён.

## Консистентность и дедупликация
- Eventual consistency допустима: данные могут обновляться с задержкой.
- Дедупликация объявлений и уведомлений — строгое требование: повторная обработка не создаёт дубликаты объявлений и не приводит к повторному уведомлению одного пользователя об одном объявлении.

## Сохранность данных
- Бизнес-данные должны сохраняться при перезапусках и сбоях: правила поиска, объявления и состояние отправки уведомлений, необходимое для дедупликации.
- Хранение логов и метрик пока необязательно; это не отменяет требования обнаруживать остановку системы.

## Данные

Все данные MVP хранятся в одном PostgreSQL. Kafka, Redis и Elasticsearch на MVP не используются.

### User

Пользователь Telegram.

```text
id, telegram_user_id, telegram_chat_id, created_at
```

При `/start` пользователь создаётся либо находится по `telegram_user_id`.

### SearchRule

Поисковое правило пользователя.

```text
id, user_id, brands[], categories[], sizes[],
max_price, currency, active, created_at, updated_at
```

- `user_id NOT NULL` — внешний ключ на `User.id`; связь `User 1:N SearchRule`.
- При удалении пользователя его правила удаляются каскадно: `ON DELETE CASCADE`.
- `categories` и `sizes` необязательны; пустой список означает отсутствие ограничения.
- Размеры пока хранятся в произвольном строковом виде, без нормализации размерных сеток.
- `max_price` — обязательный `integer`.
- `currency NOT NULL DEFAULT 'RUB'`.
- `active NOT NULL DEFAULT true`.

### SearchRulePlatform

Таблица связи правил с площадками: одно правило может использовать несколько площадок, одна площадка — участвовать в нескольких правилах.

```text
search_rule_id, platform
PRIMARY KEY (search_rule_id, platform)
```

`search_rule_id` — внешний ключ на `SearchRule.id`. Отдельной таблицы `Platform` пока нет. Коды площадок задаются фиксированным набором enum/констант в Go и сохраняются в БД как стабильные строки. Во всех сущностях поле `platform` использует тот же набор кодов.

### Listing

Объявление с внешней площадки.

```text
id, platform, external_id, title, brand, category, size,
price, currency, url, image_url, published_at,
first_seen_at, last_seen_at
UNIQUE (platform, external_id)
```

- `published_at` — время публикации на площадке; `NULL`, если неизвестно.
- `first_seen_at` — время первого обнаружения; при повторном обнаружении не меняется.
- `last_seen_at` — время последнего обнаружения; обновляется при повторном получении объявления.
- `UNIQUE(platform, external_id)` исключает дубли объявления в БД.

### ListingMatch

Связь `SearchRule N:M Listing`: одно объявление может соответствовать нескольким правилам.

```text
search_rule_id, listing_id, matched_at
PRIMARY KEY (search_rule_id, listing_id)
```

`search_rule_id` и `listing_id` — внешние ключи на `SearchRule.id` и `Listing.id`.

### Notification

Состояние отправки объявления пользователю.

```text
id, user_id, listing_id, status, attempts,
created_at, sent_at, last_error
UNIQUE (user_id, listing_id)
```

- `user_id` и `listing_id` — внешние ключи на `User.id` и `Listing.id`.
- `status`: `PENDING`, `SENT`, `FAILED`.
- `attempts` — количество попыток отправки, изначально `0`; используется для retry.
- `sent_at` — время успешной отправки; до неё `NULL`.
- `last_error` — последняя ошибка отправки; при отсутствии ошибки `NULL`.
- `UNIQUE(user_id, listing_id)` обеспечивает одну запись уведомления на пользователя и объявление, даже при совпадении с несколькими правилами.

### ParserRun

Техническая сущность для observability каждого запуска парсера по правилу и площадке.

```text
id, search_rule_id, platform, status, started_at, finished_at,
error, fetched_count, new_listings_count, image_count
```

`search_rule_id` — внешний ключ на `SearchRule.id`.

При старте parser job создаётся запись со статусом `RUNNING` и `started_at`. При завершении она обновляется до `SUCCESS` или `FAILED`, сохраняются `finished_at` и, при ошибке, `error`.

- `fetched_count` — число полученных объявлений.
- `new_listings_count` — число новых объявлений, сохранённых в БД.
- `image_count` — число полученных или обработанных изображений.

### Access patterns и транзакции

- CRUD `SearchRule` и чтение правил пользователя по `user_id`.
- Scheduler получает правила с `active = true` и выбранные площадки из `SearchRulePlatform`.
- Повторное получение объявления обрабатывается через upsert по `UNIQUE(platform, external_id)`.
- Для объявления, соответствующего правилу, одна транзакция создаёт или обновляет `Listing`, создаёт `ListingMatch` и `Notification(PENDING)`. Конфликты уникальности связей и уведомлений не создают дубли и не сбрасывают статус существующего уведомления.

```text
BEGIN
  INSERT / UPDATE Listing
  INSERT ListingMatch ON CONFLICT DO NOTHING
  INSERT Notification(PENDING) ON CONFLICT DO NOTHING
COMMIT
```

Отдельный worker отправляет уведомления в Telegram **после commit** и обновляет состояние отправки. Retry-механизм использует ту же запись `Notification`, увеличивает `attempts` и сохраняет последнюю ошибку в `last_error`. Telegram API не вызывается внутри транзакции БД.

## API и взаимодействие системы

На MVP полноценного REST API нет. Пользователь взаимодействует с системой через Telegram. Telegram updates обрабатываются Telegram handlers, которые вызывают application/use-case слой.

Use-case слой координирует доменную логику, работу с PostgreSQL через repositories, парсерами и уведомлениями. Scheduler вызывает те же внутренние use cases. Telegram handlers и scheduler — только adapters/entry points, без бизнес-логики.

```text
Telegram -> Handlers -> Application/Use Cases -> Domain/Repositories -> PostgreSQL/Parsers
Scheduler -> Application/Use Cases -> Parsers
```

Внутренние ошибки типизированы:

- `ValidationError` — некорректные входные данные.
- `NotFound` — объект не найден.
- `Forbidden` — нет доступа к действию или объекту.
- `Conflict` — конфликт с текущим состоянием.
- `ExternalServiceError` — ошибка внешнего сервиса.
- `InternalError` — внутренняя ошибка приложения.

Telegram adapter преобразует ошибки в понятные пользовательские сообщения без внутренних технических деталей.

В будущем REST API можно добавить как ещё один adapter поверх тех же use cases без переписывания бизнес-логики.

### Конкурентность

Несколько parser jobs могут выполняться параллельно, в том числе несколько разных `SearchRule` для одной площадки.

Правила:

- разные `SearchRule` для одной платформы могут выполняться параллельно;
- все запросы к одной платформе используют общий rate limiter и общий concurrency limit;
- один и тот же `SearchRule + Platform` не должен выполняться одновременно несколько раз;
- разные платформы могут парситься независимо и параллельно;
- дедупликация при конкурентной обработке обеспечивается `UNIQUE` constraints и транзакциями PostgreSQL.