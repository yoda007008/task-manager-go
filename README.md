# Task Manager на Go

RestAPI написанный на Golang с ручками для взаимодействия 
с базой данных Postgres. Реализованны все Crud операции
и linting через golangci-lint. Добавлена контейнеризация 
Docker с установкой всех зависимостей. Добавлен также 
мониторинг Prometheus + Grafana. 

Также имеются корректные моки для тестирования
их можно сгенерировать через mockgen этой командой:

```plaintext
mockery --name=TaskService --dir=./domain --output=./mocks --outpkg=mocks
```

// todo Доделать тесты (gomock lib)

## Ручки с примерами запросов
1) http://localhost:8080/tasks - получаем все задачи
2) http://localhost:8080/tasks/create - создаем новую задачу
3) http://localhost:8080/tasks/update - обновляем задачу по id
4) http://localhost:8080/tasks/delete - удаляем задачу по id
5) http://localhost:8080/tasks/{id} - получаем конкретную задачу по id 
6) http://localhost:8080/metrics - метрики текущего сервиса

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