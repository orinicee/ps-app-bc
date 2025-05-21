package http

import (
	"github.com/gin-gonic/gin"
	"github.com/orinicee/ps-app-bc/internal/interfaces/middleware"
)

type Router struct {
	engine         *gin.Engine
	authHandler    *AuthHandler
	authMiddleware *middleware.AuthMiddleware
}

func NewRouter(authHandler *AuthHandler, authMiddleware *middleware.AuthMiddleware) *Router {
	engine := gin.Default()
	return &Router{
		engine:         engine,
		authHandler:    authHandler,
		authMiddleware: authMiddleware,
	}
}

func (r *Router) SetupRoutes() {
	// Rutas públicas
	r.engine.POST("/api/auth/register", r.authHandler.Register)
	r.engine.POST("/api/auth/login", r.authHandler.Login)

	// Rutas protegidas
	protected := r.engine.Group("/api")
	protected.Use(r.authMiddleware.RequireAuth())

	// Rutas de admin
	admin := protected.Group("/admin")
	admin.Use(r.authMiddleware.RequireAdmin())
}

func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}
