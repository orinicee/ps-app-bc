package database

// TestConfig retorna la configuración de la base de datos para tests
func TestConfig() Config {
	return Config{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres123",
		DBName:   "ps_app_test",
		SSLMode:  "disable",
	}
}
