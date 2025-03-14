package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// Config contém as configurações para conexão com o banco de dados
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewConnection cria uma nova conexão com o banco de dados PostgreSQL
func NewConnection(config Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir conexão com o banco de dados: %w", err)
	}

	// Configurar o pool de conexões
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verificar se a conexão está funcionando
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao verificar conexão com o banco de dados: %w", err)
	}

	return db, nil
}

// Close fecha a conexão com o banco de dados
func Close(db *sql.DB) error {
	if db != nil {
		return db.Close()
	}
	return nil
}
