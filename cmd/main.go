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
		Host:     os.Getenv("DB_HOST"),
		Port:     getEnvAsInt("DB_PORT", 5432),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  "disable",
	}

	// Inicializar la conexión a la base de datos
	db, err := database.NewConnection(dbConfig)
	if err != nil {
		log.Fatalf("Error al inicializar la base de datos: %v", err)
	}

	// Crear el almacenamiento PostgreSQL
	storage := database.NewPostgresStorage(db)
	defer storage.Close()

	// Inicializar el servidor HTTP
	server := api.NewServer(storage)

	// Iniciar el servidor
	port := "8080"
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
