# Deployment Guide

## 🚀 Быстрый старт

### 1. Подготовка

```bash
# Клонировать репозиторий
git clone https://github.com/Ravierin/BudgetTracker.git
cd BudgetTracker

# Создать .env файл
cp .env.example .env
```

### 2. Настройка API ключей

Откройте `.env` и добавьте ключи от бирж:

```bash
# Bybit
BYBIT_API_KEY=your_key_here
BYBIT_SECRET_KEY=your_secret_here

# MEXC
MEXC_API_KEY=your_key_here
MEXC_SECRET_KEY=your_secret_here

# И так далее для Gate, Bitget
```

### 3. Запуск

```bash
# Запустить всё одной командой
./start.sh
```

### 4. Проверка

```bash
# Посмотреть статус
docker-compose ps

# Посмотреть логи
docker-compose logs -f

# Frontend: http://localhost:3000
# Backend: http://localhost:8080
```

### 5. Остановка

```bash
./stop.sh
```

---

## 📋 Требования

- Docker 20+
- Docker Compose 2+
- 512MB RAM минимум
- 1GB дискового пространства

---

## 🔧 Ручная настройка (без Docker)

### Backend

```bash
cd backend

# Установить зависимости
go mod download

# Настроить .env
cp .env.example .env

# Скомпилировать
go build -o budget-tracker ./cmd/main.go

# Запустить
./budget-tracker
```

### Frontend

```bash
cd frontend

# Установить зависимости
npm install

# Скомпилировать для production
npm run build

# Запустить dev сервер
npm run dev
```

### Database

```bash
# Создать БД
createdb BudgetTracker

# Применить миграции
cd backend
./budget-tracker migrate up
```

---

## 🐛 Troubleshooting

### Ошибка: "port already in use"

```bash
# Найти процесс на порту 8080
lsof -i :8080

# Убить процесс
kill -9 <PID>
```

### Ошибка: "database does not exist"

```bash
# Создать БД вручную
createdb BudgetTracker

# Применить миграции
./budget-tracker migrate up
```

### Ошибка: "API key not valid"

- Проверьте ключи в `.env`
- Убедитесь что ключи имеют права на **Read**
- Пересоздайте ключи на бирже

---

## 📊 Мониторинг

```bash
# Логи всех сервисов
docker-compose logs -f

# Логи конкретного сервиса
docker-compose logs -f backend

# Использование ресурсов
docker stats
```

---

## 🔄 Обновление

```bash
# Получить обновления
git pull

# Пересобрать контейнеры
docker-compose down
docker-compose build --no-cache
docker-compose up -d
```

---

## 💾 Бэкап базы данных

```bash
# Создать бэкап
docker exec budget-tracker-db pg_dump -U postgres BudgetTracker > backup.sql

# Восстановить из бэкапа
docker exec -i budget-tracker-db psql -U postgres BudgetTracker < backup.sql
```
