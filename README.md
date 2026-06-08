# CheckBot

Телеграм-бот для парсинга Google Таблиц (CSV-экспорт) в чат в виде картинок

**Стек:** Go 1.25, go-telegram-bot-api (Telegram, long polling), fogleman/gg (рендер
PNG), godotenv. Состояние хранится в JSON-файлах (`config.json`, `users.json`) в
каталоге `DATA_DIR`. Деплой — Docker

---

## Схема работы

```
Telegram ── /duty, /monitor, … ──► бот ──► CSV-export Google Таблицы
                                            │
                                            ▼
                                     FetchCSV → Parse → render (fogleman/gg)
                                            │
   фото-карточка в чат ◄────── PNG ◄────────┘
```

---

## Быстрый старт (локально)

```bash
cp .env.example .env     # вписать BOT_TOKEN (и ADMIN_ID для админ-функций)
go run ./cmd/bot         # либо: make run
```

Свой Telegram id (для `ADMIN_ID`) можно узнать, например, у [@userinfobot](https://t.me/userinfobot).
Запуск выполняется из корня проекта — оттуда подхватывается `.env`.

---

## Конфигурация

Параметры читаются из окружения (или `.env`).

| Переменная   | Обязательна | По умолчанию                  | Описание |
|--------------|:-----------:|-------------------------------|----------|
| `BOT_TOKEN`  | да          | —                             | Токен Telegram-бота от @BotFather |
| `ADMIN_ID`   | нет         | `0`                           | Telegram user id админа (доступ к ⚙️ Настройкам и рассылке) |
| `DATA_DIR`   | нет         | `.`                           | Каталог изменяемых данных (`config.json`, `users.json`) |
| `ASSETS_DIR` | нет         | рядом с бинарником / `assets` | Каталог ассетов; шрифты ищутся в `<ASSETS_DIR>/fonts` |

Источники таблиц хранятся в `config.json` (`config.example.json`), 
настраиваются прямо из бота через ⚙️ **Настройки**.

---

## Источники таблиц

Источник — это публичная Google-таблица. Привязка выполняется из чата: ⚙️ Настройки →
выбор карточки → присланная ссылка из адресной строки. Ссылка нормализуется в
`https://docs.google.com/spreadsheets/d/<id>/export?format=csv&gid=<gid>`, поэтому
таблица должна быть **доступна по ссылке без авторизации** (CSV качается без cookies).

Загрузка принимает любую таблицу, но **каждый рендерер ждёт свой формат** — иначе
выйдет пустая или нечитаемая карточка.

---

## Команды бота

| Команда                  | Назначение |
|--------------------------|------------|
| `/menu`                  | Главное меню |
| `/duty`                  | График дежурств |
| `/time_schedule`         | График учёта времени |
| `/monitor`               | Мониторинг |
| `/settings`              | Настройка источников и рассылка (только админ) |

---

## Деплой (Docker)

```bash
docker build -t checkbot .
docker run -d --env-file .env -v checkbot-data:/data checkbot
```

Изменяемые данные (`config.json`, `users.json`) живут в томе `/data`; шрифты и
сид-конфиг копируются в образ. Бот «выходит» в Telegram через long polling —
внешний HTTPS не нужен.

---

## Разработка

```bash
make run     # go run ./cmd/bot
make vet     # go vet ./...
make fmt     # gofmt -w .
```

## Структура

```
cmd/bot/            # точка входа (main, загрузка .env)
internal/
  config/           # источники таблиц, проверка админа, нормализация ссылок
  datadir/          # разрешение каталога данных (DATA_DIR)
  sheets/           # загрузка и парсинг CSV
  render/           # отрисовка PNG (fogleman/gg): duty, timesheet, monitor
  store/            # хранилище пользователей (users.json)
  telegram/         # бот: команды, callback-и, админка, рассылка
assets/fonts/       # шрифты Roboto
Dockerfile
```
