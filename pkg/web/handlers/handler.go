package handlers

import (
	"net/http"

	"github.com/arskydev/weatherman/pkg/service"
	"github.com/arskydev/weatherman/pkg/web/internal/middleware"
	"github.com/gorilla/mux"
)

type Handler struct {
	service *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{
		service: s,
	}

}

func (h *Handler) InitRoutes() http.Handler {
	m := middleware.New(h.service.Authorization)
	r := mux.NewRouter()
	// AUTH
	r.HandleFunc("/auth/sign-up", h.signUp).Methods("POST")
	r.HandleFunc("/auth/sign-in", h.signIn).Methods("POST")
	// API
	r.Handle("/api/get-mock", m.ValidateJWT(h.getMock)).Methods("GET")
	r.Handle("/api/get-weather", m.ValidateJWT(h.getWeather)).Methods("GET")
	return r
}
