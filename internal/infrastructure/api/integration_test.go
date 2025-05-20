package api

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/orinicee/ps-app-bc/internal/infrastructure/database"
)

// setupTestDB configura la base de datos para los tests
func setupTestDB(t *testing.T) (*sql.DB, func()) {
	// Obtener configuración de test
	config := database.TestConfig()

	// Crear conexión a la base de datos
	db, err := database.NewConnection(config)
	if err != nil {
		t.Fatalf("Error al conectar a la base de datos de test: %v", err)
	}

	// Crear el storage
	storage := database.NewPostgresStorage(db)

	// Función de limpieza
	cleanup := func() {
		storage.Close()
	}

	return db, cleanup
}

func TestHealthCheckIntegration(t *testing.T) {
	// Configurar la base de datos
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Crear el servidor con la base de datos real
	storage := database.NewPostgresStorage(db)
	server := NewServer(storage)

	tests := []struct {
		name           string
		setupDB        func(*sql.DB) error
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "healthy database",
			setupDB: func(db *sql.DB) error {
				// No necesitamos hacer nada, la base de datos está saludable
				return nil
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
		},
		{
			name: "unhealthy database",
			setupDB: func(db *sql.DB) error {
				// Cerrar la conexión para simular una base de datos no saludable
				return db.Close()
			},
			expectedStatus: http.StatusServiceUnavailable,
			expectedBody:   "Database connection failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Configurar el estado de la base de datos
			if err := tt.setupDB(db); err != nil {
				t.Fatalf("Error al configurar la base de datos: %v", err)
			}

			// Crear una petición de prueba
			req := httptest.NewRequest("GET", "/health", nil)
			w := httptest.NewRecorder()

			// Llamar al handler
			server.healthCheck(w, req)

			// Verificar el código de estado
			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d; got %d", tt.expectedStatus, w.Code)
			}

			// Verificar el cuerpo de la respuesta (ignorando espacios en blanco)
			gotBody := strings.TrimSpace(w.Body.String())
			if gotBody != tt.expectedBody {
				t.Errorf("expected body %q; got %q", tt.expectedBody, gotBody)
			}
		})
	}
}
