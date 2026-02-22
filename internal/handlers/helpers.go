package handlers

import "net/http"

// функция helper для проверки http метода
func MethodHandler(handlerFunc http.HandlerFunc, allowedMethod string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != allowedMethod {
			http.Error(w, "Метод не разрешен", http.StatusBadRequest)
			return
		}
		handlerFunc(w, r)
	}
}

// специальный обработчик методов, которые имеют ID задачи
func TaskIDHandler(handler *TaskHandlers) http.HandlerFunc {
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
