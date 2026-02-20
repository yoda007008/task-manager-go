package main

import (
	"net/http"
	"task-manager-go/internal/handlers"
)

func main() {
	//// todo run service
	//databaseURL := os.Getenv("DATABASE_URL") // todo make directory config and parse config
	//if databaseURL == "" {
	//	databaseURL = "postgresql://taskuser:taskpass@localhost:5432/taskdb?sslmode=disable"
	//}
	//
	//serverPort := os.Getenv("SERVER_PORT")
	//if serverPort == "" {
	//	serverPort = "8080"
	//}
	//
	//slog.Info("Запуск приложения TaskManagerAPI...")
	//
	//db, err := database.Connect(databaseURL)
	//if err != nil {
	//	slog.Error("Ошибка подключения к базе", err)
	//}
	//
	//defer db.Close()
	//
	//slog.Info("Успешное подключение к базе...")
	//
	//taskStore := database.NewTaskStore(db)
	//
	//handler := handlers.NewTaskHandler(taskStore)
	//
	//mux := http.NewServeMux()

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
