package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/v00vaa/final_task/pkg/api"
)

func Run(port, pass string) error {
	r := chi.NewRouter()
	api.Init(r, pass)
	return http.ListenAndServe(":"+port, r)
}
