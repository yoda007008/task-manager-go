# Task Manager на Go

RestAPI написанный на Golang с ручками для взаимодействия 
с базой данных Postgres. Реализованны все Crud операции
и linting через golangci-lint. Добавлена контейнеризация 
Docker с установкой всех зависимостей

// todo Добавить мониторинг Prometeus + Grafana 

// todo Добавить тесты (testify)

// todo ручки с примерами запросов

## Запуск приложения

Чтобы запустить приложение сначала cклонируйте репозиторий

```plaintext
git clone <name-repo>
```

Далее зайдите в папку с проектом
```plaintext
cd <dir>
```

Запустите сборку образа 
```plaintext
docker compose up --build
```