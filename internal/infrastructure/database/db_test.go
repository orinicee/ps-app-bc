package database

import (
	"testing"
)

func TestDatabaseConnection(t *testing.T) {
	config := Config{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres123",
		DBName:   "ps_app",
		SSLMode:  "disable",
	}

	// Intentar conectar a la base de datos
	db, err := NewConnection(config)
	if err != nil {
		t.Fatalf("Error al conectar a la base de datos: %v", err)
	}
	defer db.Close()

	// Test de consulta a la tabla content_types
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM content_types").Scan(&count)
	if err != nil {
		t.Fatalf("Error al consultar content_types: %v", err)
	}

	// Verificar que existen los 4 tipos de contenido que insertamos
	if count != 4 {
		t.Errorf("Se esperaban 4 tipos de contenido, pero se encontraron %d", count)
	}

	// Test adicional: verificar los tipos específicos
	rows, err := db.Query("SELECT type_name FROM content_types ORDER BY type_name")
	if err != nil {
		t.Fatalf("Error al consultar los tipos de contenido: %v", err)
	}
	defer rows.Close()

	expectedTypes := []string{"audio", "live", "video", "vlog"}
	var types []string

	for rows.Next() {
		var typeName string
		if err := rows.Scan(&typeName); err != nil {
			t.Fatalf("Error al escanear tipo de contenido: %v", err)
		}
		types = append(types, typeName)
	}

	// Verificar que tenemos todos los tipos esperados
	for i, expectedType := range expectedTypes {
		if i >= len(types) || types[i] != expectedType {
			t.Errorf("Tipo de contenido incorrecto en posición %d. Esperado: %s, Obtenido: %s",
				i, expectedType, types[i])
		}
	}
}
