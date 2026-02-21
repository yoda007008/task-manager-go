package middleware

import (
	"log"
	"net/http"
)

// регистрируем middleware, который будет логировать все http запросы
func LoggingMiddleware(next http.Handler) http.Handler {
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
func CorsMiddleware(next http.Handler) http.Handler {
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
