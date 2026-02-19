package http

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/yourusername/analytics-service/internal/port"
)

type Server struct {
	router   *mux.Router
	handler  *Handler
	port     string
}

func NewServer(useCases port.UseCases, port string) *Server {
	s := &Server{
		router:  mux.NewRouter(),
		handler: NewHandler(useCases),
		port:    port,
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	// API routes
	api := s.router.PathPrefix("/api").Subrouter()
	
	api.HandleFunc("/items", s.handler.CreateItem).Methods("POST")
	api.HandleFunc("/items", s.handler.GetItems).Methods("GET")
	api.HandleFunc("/items/{id}", s.handler.GetItem).Methods("GET")
	api.HandleFunc("/items/{id}", s.handler.UpdateItem).Methods("PUT")
	api.HandleFunc("/items/{id}", s.handler.DeleteItem).Methods("DELETE")
	api.HandleFunc("/analytics", s.handler.GetAnalytics).Methods("GET")

	// Serve static files
	s.router.PathPrefix("/").Handler(http.FileServer(http.Dir("./web")))
}

func (s *Server) Start() error {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(s.router)
	
	fmt.Printf("Server starting on port %s\n", s.port)
	return http.ListenAndServe(":"+s.port, handler)
}
