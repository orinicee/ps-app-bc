package api

import (
	"github.com/gin-gonic/gin"
	"github.com/orinicee/ps-app-bc/internal/domain"
	"github.com/orinicee/ps-app-bc/internal/infrastructure/database"
	"github.com/orinicee/ps-app-bc/internal/infrastructure/repository"
	"github.com/orinicee/ps-app-bc/internal/interfaces/http"
	"github.com/orinicee/ps-app-bc/internal/usecase"
)

// Server representa el servidor HTTP
type Server struct {
	storage domain.Storage
	router  *gin.Engine
}

// NewServer crea una nueva instancia del servidor
func NewServer(storage domain.Storage) *Server {
	router := gin.Default()

	// Crear repositorios
	postgresStorage := storage.(*database.PostgresStorage)
	userRepo := repository.NewUserRepository(postgresStorage.DB())

	// Crear casos de uso
	jwtKey := []byte("your-secret-key") // TODO: Mover a configuración
	authUseCase := usecase.NewAuthUseCase(userRepo, jwtKey)

	// Crear handlers
	authHandler := http.NewAuthHandler(authUseCase)

	// Configurar rutas
	api := router.Group("/api")
	{
		// Rutas de autenticación
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Ruta de health check
		api.GET("/health", func(c *gin.Context) {
			if err := storage.HealthCheck(); err != nil {
				c.JSON(503, gin.H{"error": "Database connection failed"})
				return
			}
			c.JSON(200, gin.H{"status": "OK"})
		})
	}

	return &Server{
		storage: storage,
		router:  router,
	}
}

// Start inicia el servidor HTTP
func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}
