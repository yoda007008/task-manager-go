package main

import (
	"flag"
	"log"
	"log/slog"
	"net/http"
	"task-manager-go/internal/config"
	"task-manager-go/internal/database"
	"task-manager-go/internal/handlers"
	"task-manager-go/internal/middleware"
)

func main() {
	configPath := flag.String("config", "/home/kirill/GolandProjects/task-manager-go/internal/config/config.yaml", "path to config file") // config

	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		slog.Error("Error parse config", "error", err)
	}

	databaseURL := cfg.Database.Url
	serverPort := cfg.Server.Port

	slog.Info("Запуск приложения TaskManagerAPI...")

	db, err := database.Connect(databaseURL)
	if err != nil {
		slog.Error("Ошибка подключения к базе", err)
	}

	defer db.Close()

	slog.Info("Успешное подключение к базе...")

	taskStore := database.NewTaskStore(db)

	handler := handlers.NewTaskHandler(taskStore)

	mux := http.NewServeMux()

	mux.HandleFunc("/tasks", handlers.MethodHandler(handler.GetAllTask, "GET"))
	mux.HandleFunc("/tasks/create", handlers.MethodHandler(handler.CreateTask, "POST"))
	mux.HandleFunc("/tasks/", handlers.TaskIDHandler(handler))

	loggerMux := middleware.LoggingMiddleware(mux)
	corsHandler := middleware.CorsMiddleware(loggerMux) // если фронтенд приложение работает на другом домене, то добавляем CORS

	log.Printf("Сервер запущен на порту %s", serverPort)
	slog.Info("Доступные endpoints: ")
	slog.Info("  GET 	 /tasks - получить все задачи")
	slog.Info("  GET    /tasks/{id}  - Получить задачу по ID")
	slog.Info("  POST   /tasks/create - Создать новую задачу")
	slog.Info("  PUT    /tasks/{id}  - Обновить задачу")
	slog.Info("  DELETE /tasks/{id}  - Удалить задачу")

	err = http.ListenAndServe(serverPort, corsHandler)
	if err != nil {
		slog.Error("Ошибка запуска сервера", err)
	}
}
