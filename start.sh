#!/bin/bash

start=none
if [[ -n "$1" ]]; then start="$1"; fi

. ./lib.sh

if [[ "${start}" == "depends" ]]; then
    hint="Устанавливаем зависимости:"
    printexec "dependency" sudo apt install -y docker xfce4-screenshooter curl && sudo groupadd docker && sudo usermod -aG docker ${USER}
fi

hint="Краткое описание:"
printexec "no-exec" "
Два приложения на go компилируются из исходных файлов:
1. Сервер, реализует два сервиса
2. Клиент, работает с этими двумя сервисами и имеет свой http API для внешних запросов

Логика такова, пользователь отправляет запрос на клиент, который, в зависимости от запроса,
вызывает тот или иной сервер по GRPC. Полученный ответ отдаётся пользователю.

Краткую инструкцию по использованию можно прочитать, выполнив запрос:
curl http://localhost:55000/

Программа позволяет 'хранить' секреты (текст), для сохранения секрета необходимо вызвать команду:
http://localhost:55000/register?secret=СЕКРЕТ
где вместо слова СЕКРЕТ нужно указать свой секрет. В ответ вы получите уникальный идентификатор.
Для получения секрета необходимо вызвать команду:
http://localhost:55000/secret?id=ИДЕНТИФИКАТОР
где вместо слова ИДЕНТИФИКАТОР нужно указать ранее полученный идентификатор. В случае, если
идентификатор валиден, в ответ вы получаете сохранённый секрет или ошибку, если идентификатор не валиден.
"


hint="Исходный код серверной части приложения:"
printexec "no-exec" https://github.com/mephi-learn/conveer_module7/blob/main/cmd/server/main.go

hint="Исходный код клиентской части приложения:"
printexec "no-exec" https://github.com/mephi-learn/conveer_module7/blob/main/cmd/client/main.go

hint="proto файл первого сервиса (регистрация секретов):"
printexec "show_register_service" cat proto/register.proto

hint="proto файл второго сервиса (выдача секретов):"
printexec "show_register_service" cat proto/secret.proto

hint="Содержимое Dockerfile сервера:"
printexec "show_server_dockerfile" cat docker/server/Dockerfile

hint="Содержимое Dockerfile клиента:"
printexec "show_client_dockerfile" cat docker/client/Dockerfile

hint="Собираем docker compose:"
printexec "build_images" docker-compose build --no-cache

hint="Запускаем docker compose в фоне:"
printexec "run_images" docker-compose up -d

hint="Отображаем инструкцию:"
printexec "help" curl http://localhost:55000/

hint="Сохраняем секрет:"
printexec "register_secret" 'secret_id="$(curl http://localhost:55000/register?secret=MyiImportantSecret)"'

hint="Пробуем получить секрет по некорректному идентификатору:"
printexec "incorrect_secret" curl http://localhost:55000/secret?id=invalid

hint="Пробуем получить секрет по правильному идентификатору:"
printexec "correct_secret" "curl http://localhost:55000/secret?id=${secret_id} | grep --color 'MyiImportantSecret' && echo 'Задание выполнено' || echo 'Задание провалено'"

hint="Останавливаем docker контейнеры:"
printexec "stop_docker" docker-compose down
