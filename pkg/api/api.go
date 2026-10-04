package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

var password string

func Init(r chi.Router, configuredPassword string) {
	password = configuredPassword

	r.Post("/api/signin", loginHandler)

	r.Get("/api/nextdate", nextDateHandler)

	r.Get("/api/task", auth(getTaskHandler))
	r.Post("/api/task", auth(addTaskHandler))
	r.Put("/api/task", auth(updateTaskHandler))
	r.Delete("/api/task", auth(deleteTaskHandler))

	r.Post("/api/task/done", auth(doneTaskHandler))

	r.Get("/api/tasks", auth(tasksHandler))

	r.Method(http.MethodGet, "/*", http.FileServer(http.Dir("./web")))
}
