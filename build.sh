#!/bin/bash

# Сборка Docker образа
docker build -t todo-app .

# Запуск контейнера
docker run -d -p 9092:9092 \
  -e TODO_PASSWORD=secret123 \
  -v todo-data:/data \
  --name todo-app \
  todo-app

echo "Приложение запущено на http://localhost:9092"
echo "Пароль: secret123"
