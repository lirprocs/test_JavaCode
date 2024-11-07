# Test_JavaCode [![Test](https://github.com/lirprocs/test_JavaCode/actions/workflows/test.yaml/badge.svg)](https://github.com/lirprocs/test_JavaCode/actions/workflows/test.yaml)
## Реализованы следующие эндпоинты:
1. POST /api/v1/wallet -  изменениe счета в базе данных.
2. GET  /api/v1/wallets/{walletID} - проверка баланса.
## Установка
1. Клонируйте репозиторий:
```bash
git clone https://github.com/lirprocs/test_JavaCode.git
```
## Предварительная настройка
1. Перейдите в директорию проекта:
```bash
cd test_JavaCode
```
2. Создайте в корне проекта 2 файла конфигурации с именами:
```bash
config.env 
```
```bash
config_test.env
```
3. Внесите в них данные как в примере .env.example:
```bash
POSTGRES_USER=example
POSTGRES_PASSWORD=qwert
POSTGRES_DB=example
DB_HOST=db
DB_PORT=5432
APP_PORT=8080
```
## Запуск тестов
1. Перейдите в директорию проекта (Не нужно, еслу уже находитесь в ней):
```bash
cd test_JavaCode
```
2. Для запуска тестов базы данных:
```bash
docker-compose -f docker-compose.test.yml --env-file ./config_test.env up --build
```
Проверка работы методов UpdateBalance и GetBalance в условиях реального подключения к базе данных и работы в конкурентной среде. \
После вывода результатов тестов остановите выполнение (CTR + C). <br>

3. Для запуска юнит-тестов:
```bash
go test -tags !docker ./...
```
4. После тестов в дерриктории проекта выполните:
```bash
docker-compose down --rmi all --volumes --remove-orphans
```

## Запуск приложения
1. Перейдите в директорию проекта (Не нужно, еслу уже находитесь в ней):
```bash
cd test_JavaCode
```
2. Запустите Docker контейнер с приложением и базой данных:
```bash
docker-compose --env-file ./config.env  up --build
```


