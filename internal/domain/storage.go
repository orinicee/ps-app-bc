package domain

// Storage define la interfaz para operaciones de almacenamiento
type Storage interface {
	// HealthCheck verifica la conexión con el almacenamiento
	HealthCheck() error
	// Close cierra la conexión con el almacenamiento
	Close() error
}
