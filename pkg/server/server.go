package server

import (
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/v00vaa/final_task/pkg/api"
	"github.com/v00vaa/final_task/tests"
)

func Run() {
	r := chi.NewRouter()
	
	api.Init(r)

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = strconv.Itoa(tests.Port)
	}

	r.Handle("/*", http.FileServer(http.Dir("./web")))

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		panic(err)
	}
}
