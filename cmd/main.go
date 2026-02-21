package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"task-manager-go/internal/database"
	"task-manager-go/internal/handlers"
)

func main() {
	// todo run service with graceful shutdown
	databaseURL := os.Getenv("DATABASE_URL") // todo make directory config and parse config
	if databaseURL == "" {
		databaseURL = "postgresql://taskuser:taskpass@localhost:5440/taskdb?sslmode=disable"
	}

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
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

	loggerMux := loggingMiddleware(mux)
	corsHandler := corsMiddleware(loggerMux) // если фронтенд приложение работает на другом домене, то добавляем CORS

	log.Printf("Сервер запущен на порту %s", serverPort)
	slog.Info("Доступные endpoints: ")
	slog.Info("  GET 	 /tasks - получить все задачи")
	slog.Info("  GET    /tasks/{id}  - Получить задачу по ID")
	slog.Info("  POST   /tasks/create - Создать новую задачу")
	slog.Info("  PUT    /tasks/{id}  - Обновить задачу")
	slog.Info("  DELETE /tasks/{id}  - Удалить задачу")

	err = http.ListenAndServe(":"+serverPort, corsHandler)
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

// todo регистрируем middleware, который будет логировать все http запросы
func loggingMiddleware(next http.Handler) http.Handler { // todo перенести в отдельную директорию /internal/middleware
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		// Пояснение:
		// r.Method - HTTP метод (GET, POST, и т.д.)
		// r.URL.Path - путь запроса (/tasks, /tasks/1, и т.д.)
		// r.RemoteAddr - IP адрес клиента
		next.ServeHTTP(w, r)
	})
}

// эта функция, которая поможет интегрироваться с фронтом при помощи CORS
// CORS позволяет браузерам отправлять запросы к API с других доменов
func corsMiddleware(next http.Handler) http.Handler { // todo перенести в отдельную директорию
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// разрешаем запросы с других доменов
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Access-Control-Allow-Methods - разрешенные HTTP методы
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// Access-Control-Allow-Headers - разрешенные заголовки в запросе
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Обрабатываем preflight запрос
		// Браузеры отправляют OPTIONS запрос перед основным запросом для проверки CORS
		if r.Method == "OPTIONS" {
			// Отвечаем статусом 200 и завершаем обработку
			w.WriteHeader(http.StatusOK)
			return
		}

		// Вызываем следующий обработчик
		next.ServeHTTP(w, r)
	})
}
