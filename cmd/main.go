package main

import (
	"log"
	"os"
	"strconv"

	"github.com/orinicee/ps-app-bc/internal/infrastructure/api"
	"github.com/orinicee/ps-app-bc/internal/infrastructure/database"
)

func main() {
	// Configuración de la base de datos
	dbConfig := database.Config{
		Host:     getEnvOrDefault("DB_HOST", "127.0.0.1"),
		Port:     getEnvAsInt("DB_PORT", 5432),
		User:     getEnvOrDefault("DB_USER", "postgres"),
		Password: getEnvOrDefault("DB_PASSWORD", "postgres123"),
		DBName:   getEnvOrDefault("DB_NAME", "ps_app"),
		SSLMode:  getEnvOrDefault("DB_SSL_MODE", "disable"),
	}

	// Inicializar la conexión a la base de datos
	db, err := database.NewConnection(dbConfig)
	if err != nil {
		log.Fatalf("Error al inicializar la base de datos: %v", err)
	}

	// Crear el almacenamiento PostgreSQL
	storage := database.NewPostgresStorage(db)
	defer storage.Close()

	// Verificar la conexión
	if err := storage.HealthCheck(); err != nil {
		log.Fatalf("Error al verificar la conexión a la base de datos: %v", err)
	}

	// Inicializar el servidor HTTP
	server := api.NewServer(storage)

	// Iniciar el servidor
	port := getEnvOrDefault("PORT", "8080")
	log.Printf("Servidor iniciado en el puerto %s", port)
	if err := server.Start(":" + port); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}

func getEnvAsInt(key string, defaultVal int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultVal
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
