package api

import (
	"net/http"

	"github.com/orinicee/ps-app-bc/internal/domain"
)

// Server representa el servidor HTTP
type Server struct {
	storage domain.Storage
}

// NewServer crea una nueva instancia del servidor
func NewServer(storage domain.Storage) *Server {
	return &Server{
		storage: storage,
	}
}

// Start inicia el servidor HTTP
func (s *Server) Start(addr string) error {
	// Configurar rutas
	http.HandleFunc("/health", s.healthCheck)

	// Iniciar el servidor
	return http.ListenAndServe(addr, nil)
}

// healthCheck maneja las peticiones de health check
func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	if err := s.storage.HealthCheck(); err != nil {
		http.Error(w, "Database connection failed", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
