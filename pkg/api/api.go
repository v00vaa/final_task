package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Init(r chi.Router) {
	r.Post("/api/signin", loginHandler)

	r.Get("/api/nextdate", nextDateHandler)

	r.Get("/api/task", auth(getTaskHandler))
	r.Post("/api/task", auth(addTaskHandler))
	r.Put("/api/task", auth(updateTaskHandler))
	r.Delete("/api/task", auth(deleteTaskHandler))

	r.Post("/api/task/done", auth(doneTaskHandler))

	r.Get("/api/tasks", auth(tasksHandler))

	r.Handle("/*", http.FileServer(http.Dir("./web")))
}
