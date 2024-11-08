# Test_JavaCode [![Test](https://github.com/lirprocs/test_JavaCode/actions/workflows/test.yaml/badge.svg)](https://github.com/lirprocs/test_JavaCode/actions/workflows/test.yaml)
***MyWalletApp*** — это REST API приложение для управления балансом кошельков, разработанное на Golang. Приложение поддерживает операции пополнения и снятия средств, а также проверку текущего баланса. Проект спроектирован для работы в Docker и включает конфигурации для тестирования в конкурентной среде.
## Реализованы следующие эндпоинты:
1. POST `/api/v1/wallet` -  Выполнение операции пополнения (`DEPOSIT`) или снятия (`WITHDRAW`) средств.

**Пример запроса в Postman:**
```bash
{
    "walletId": "{{$guid}}",
    "operationType": "DEPOSIT",
    "amount": 123 
}
````
2. GET  `/api/v1/wallets/{walletID}` - Получение текущего баланса для указанного кошелька.

**Пример ответа:**
```bash
{
    "balance": 1000
}

````
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
3. Внесите в них данные как в примере [.env.example](./.env.example).

## Запуск тестов
1. Перейдите в директорию проекта (Не нужно, если уже находитесь в ней):
```bash
cd test_JavaCode
```
2. Для выполнения тестов в Docker-среде используйте файл `docker-compose.test.yml`. Выполните команду:
```bash
docker-compose -f docker-compose.test.yml --env-file ./config_test.env up --build
```
Проверка работы методов UpdateBalance и GetBalance в условиях реального подключения к базе данных и работы в конкурентной среде. \
После вывода результатов тестов остановите выполнение (`CTR + C`). <br>

3. Для запуска юнит-тестов:
```bash
go test -tags !docker ./...
```
4. После тестов в директории проекта выполните:
```bash
docker-compose down --rmi all --volumes --remove-orphans
```

## Запуск приложения
1. Перейдите в директорию проекта (Не нужно, если уже находитесь в ней):
```bash
cd test_JavaCode
```
2. Запустите Docker контейнер с приложением и базой данных:
```bash
docker-compose --env-file ./config.env  up --build
```

Приложение будет доступно по адресу http://localhost:8080.


