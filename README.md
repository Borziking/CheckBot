<h3 align="center">Telegram-bot CheckBot</h3> 

<p align="center">
  <a href="https://go.dev"><img src="https://img.shields.io/badge/Go-1.25-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go"></a>
  <a href="https://www.docker.com"><img src="https://img.shields.io/badge/Docker-2496ED?style=flat-square&logo=docker&logoColor=white" alt="Docker"></a>
</p>

Телеграм-бот для парсинга Google Таблиц (CSV-экспорт) в чат в виде картинок

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

Источник — публичная Google-таблица, привязывается из бота (⚙️ Настройки). 

Должна открываться по ссылке без авторизации — бот качает CSV-экспорт. Каждый рендерер ждёт
свой формат таблицы, иначе карточка выйдет пустой.


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
