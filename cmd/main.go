package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"task-manager-go/internal/database"
	"task-manager-go/internal/handlers"
	"task-manager-go/internal/middleware"
)

func main() {
	// todo run service with graceful shutdown
	databaseURL := os.Getenv("DATABASE_URL") // todo make directory config and parse config
	if databaseURL == "" {
		databaseURL = "postgresql://taskuser:taskpass@localhost:5440/taskdb?sslmode=disable"
	}

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = ":8080"
	}

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

	mux.HandleFunc("/tasks", methodHandler(handler.GetAllTask, "GET"))
	mux.HandleFunc("/tasks/create", methodHandler(handler.CreateTask, "POST"))
	mux.HandleFunc("/tasks/", taskIDHandler(handler))

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

// функция helper для проверки http метода
func methodHandler(handlerFunc http.HandlerFunc, allowedMethod string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != allowedMethod {
			http.Error(w, "Метод не разрешен", http.StatusBadRequest)
			return
		}
		handlerFunc(w, r)
	}
}

// специальный обработчик методов, которые имеют ID задачи
func taskIDHandler(handler *handlers.TaskHandlers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetTask(w, r)
		case http.MethodPut:
			handler.UpdateTask(w, r)
		case http.MethodDelete:
			handler.DeleteTask(w, r)
		default:
			http.Error(w, "Метод не разрешен", http.StatusBadRequest)
		}
	}
}
