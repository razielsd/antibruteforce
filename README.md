# AntiBruteforce
[![codecov](https://codecov.io/gh/razielsd/antibruteforce/branch/master/graph/badge.svg)](https://codecov.io/gh/razielsd/antibruteforce)
[![Go Report Card](https://goreportcard.com/badge/github.com/razielsd/antibruteforce)](https://goreportcard.com/report/github.com/razielsd/antibruteforce)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/razielsd/antibruteforce)

## Общее описание
Сервис предназначен для борьбы с подбором паролей при авторизации в какой-либо системе.

Сервис вызывается перед авторизацией пользователя и может либо разрешить, либо заблокировать попытку.

Предполагается, что сервис используется только для server-server, т.е. скрыт от конечного пользователя.

## API
 * _/api/user/allow_ - проверка возможности авторизации
 * _/api/whitelist_ - получить белый список  IP адресов или адрес подсети
 * _/api/whitelist/add_ - добавить в белый список IP или адрес подсети
 * _/api/whitelist/remove_ - удалить из белого списка IP или адрес подсети
 * _/api/blacklist_ - получить черный список  IP адресов или адрес подсети
 * _/api/blacklist/add_ - добавить в черный список IP или адрес подсети
 * _/api/blacklist/remove_ - удалить из черного списка IP или адрес подсети
 * _/api/bucket/drop/login_ - удалить логин из бакета (сбросить счетчик авторизаций для данного логина)
 * _/api/bucket/drop/pwd_ - удалить пароль из бакета
 * _/api/bucket/drop/ip_ - удалить IP из бакета

Подробнее в [swagger](https://editor.swagger.io/?url=https://raw.githubusercontent.com/razielsd/antibruteforce/master/doc/swagger.yml)

##  Мониторинг
API содержит несколько методов для мониторинга приложения:
 * _/mertics_ - метрики приложения для prometheus
 * _/health/liveness_ - k8s liveness probe
 * _/health/readiness_ - k8s readiness probe

## CLI
Для запуска локально стоит использовать `bin/abf`, в docker - `bin/abfc`
```Использование:
abf [command]

Available Commands:
blacklist   Show/add/remove blacklist
bucket      Drop bucket by login, password or ip
completion  generate the autocompletion script for the specified shell
help        Help about any command
server      Run service
version     Show version
whitelist   Show/add/remove whitelist
```
## Сборка
 * `make build` - собирает сервис локально
 * `make build-img` - собирает образ для докера
 * `make run` - запускает сервис локально
 * `make run-img` - запускает сервис в докере
 * `make lint` - запуск линтера   
 * `make test` - запуск тестов, без интеграционных
 * `make test-int` - запуск тестов, включая интеграционные
 * `make test-e2e` - запуск e2e-тестов
 * `make test-img`, `make test-int-img`, `make test-e2e-img` - те же тесты, но запускаются в докере.

Также есть таргеты `test-int-coverage`, `test-coverage`, `test100` они используются для запуска на CI.

##  Конфигурация
Для конфигурации сервиса используются переменные окружения:
 * `ABF_ADDR` - адрес, на котором принимать подключения, значение по умолчанию: 0.0.0.0:8080
 * `ABF_RATE_LOGIN` - кол-во попыток авторизации в минуту для логина, значение по умолчанию: 10 
 * `ABF_RATE_PWD` - кол-во попыток авторизации в минуту для пароля, значение по умолчанию: 100
 * `ABF_RATE_IP` - кол-во попыток авторизации в минуту для IP, значение по умолчанию: 1000
 * `ABF_WHITELIST` - белый список ip/подсетей, пример: 192.168.1.10,10.10.1.0/24
 * `ABF_BLACKLIST` - черный список ip/подсетей, пример: 192.168.1.10,10.10.1.0/24
 * `ABF_LOG_LEVEL` - Уровень логирования, возможные значения: DEBUG, INFO, WARN, ERROR, значение по умолчанию - DEBUG

## Алгоритм работы
Сервис ограничивает частоту попыток авторизации для различных комбинаций параметров, например:
* не более N = 10 попыток в минуту для данного логина.
* не более M = 100 попыток в минуту для данного пароля (защита от обратного brute-force).
* не более K = 1000 попыток в минуту для данного IP (число большое, т.к. NAT).
 
Для работы используется пакет [time/rate](https://pkg.go.dev/golang.org/x/time/rate), в котором реализован алгоритм [token bucket](https://en.wikipedia.org/wiki/Token_bucket)

## Развертывание
Для запуска локально необходимо выполнить команды в директории с проектом:
 * `make build` - собрать сервис
 * `make run` - запуск сервиса

Для запуска в докере необходимо выполнить команды в директории с проектом:
 * `make build-img` - собрать сервис
 * `make run-img` - запуск сервиса

## Тестирование
 * Сервис покрыт тестами, core-функционал на 100%
 * Есть интеграционные тесты - в рамках теста поднимается сервер и тестируется логика работы.
 * Есть e2e тесты - проверяется, что все endpoint доступны, cli - работает. 
