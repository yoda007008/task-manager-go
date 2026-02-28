package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"task-manager-go/internal/config"
	"task-manager-go/internal/database"
	"task-manager-go/internal/handlers"
	"task-manager-go/internal/middleware"
	"time"
)

func main() {
	configPath := flag.String("config", "./internal/config/config.yaml", "path to config file") // config

	flag.Parse()

	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := slog.New(logHandler).With("task-manager", "database")

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		logger.Error("Ошибка парсинга конфига", "error", err)
	}

	databaseURL := cfg.Database.URL
	serverPort := cfg.Server.Port

	logger.Info("Запуск приложения TaskManagerAPI...")

	db, err := database.Connect(databaseURL)
	if err != nil {
		logger.Error("Ошибка подключения к базе", "error", err)
		os.Exit(1)
	}

	defer db.Close()

	logger.Info("Успешное подключение к базе...")

	taskStore := database.NewTaskStore(db)

	handler := handlers.NewTaskHandler(taskStore)

	mux := http.NewServeMux()

	mux.HandleFunc("/tasks", handlers.MethodHandler(handler.GetAllTask, "GET"))
	mux.HandleFunc("/tasks/create", handlers.MethodHandler(handler.CreateTask, "POST"))
	mux.HandleFunc("/tasks/", handlers.TaskIDHandler(handler))

	loggerMux := middleware.LoggingMiddleware(mux)
	corsHandler := middleware.CorsMiddleware(loggerMux) // если фронтенд приложение работает на другом домене, то добавляем CORS

	server := &http.Server{
		Addr:    serverPort,
		Handler: corsHandler,
	}

	var wg sync.WaitGroup
	upperCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err = server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Error("Ошибка запуска сервера", "error", err)
		}
	}()

	log.Printf("Сервер запущен на порту %s", serverPort)
	logger.Info("Доступные endpoints: ")
	logger.Info("  GET 	 /tasks - получить все задачи")
	logger.Info("  GET    /tasks/{id}  - Получить задачу по ID")
	logger.Info("  POST   /tasks/create - Создать новую задачу")
	logger.Info("  PUT    /tasks/{id}  - Обновить задачу")
	logger.Info("  DELETE /tasks/{id}  - Удалить задачу")

	<-upperCtx.Done()

	logger.Info("Остановка task-manager-go")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err = server.Shutdown(shutdownCtx); err != nil {
		logger.Error("Ошибка при Graceful Shutdown", "error", err)
	}

	logger.Info("Успешный Graceful Shutdown")

	wg.Wait()
}
