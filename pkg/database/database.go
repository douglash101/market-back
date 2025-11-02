package database

import (
	"database/sql"
	"market/pkg/config"
	"fmt"
	"sync"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

// PostgresDB é o wrapper para sql.DB
type PostgresDB struct {
	conn *sql.DB
}

var (
	once     sync.Once
	instance *PostgresDB
)

// GetInstance retorna a instância singleton do banco de dados
func GetInstance(log *zap.SugaredLogger) *PostgresDB {
	once.Do(func() {
		connStr := fmt.Sprintf(
			"host=%s user=%s dbname=%s password=%s sslmode=disable",
			config.Get().DATABASE_HOST,
			config.Get().DATABASE_USER,
			config.Get().DATABASE_NAME,
			config.Get().DATABASE_PASSWORD,
		)

		db, err := sql.Open("postgres", connStr)
		if err != nil {
			log.Fatalf("❌ Falha ao conectar no banco: %v", err)
		}

		if err = db.Ping(); err != nil {
			log.Fatalf("❌ Falha ao pingar o banco: %v", err)
		}

		log.Info("✅ Conectado com sucesso ao PostgreSQL")
		instance = &PostgresDB{conn: db}
	})
	return instance
}

// Close fecha a conexão com o banco
func (db *PostgresDB) Close() error {
	return db.conn.Close()
}

// Query executa uma consulta que retorna múltiplas linhas
func (db *PostgresDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.conn.Query(query, args...)
}

// Exec executa um comando que não retorna linhas (INSERT, UPDATE, DELETE)
func (db *PostgresDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.conn.Exec(query, args...)
}

// Prepare prepara uma instrução para execução posterior
func (db *PostgresDB) Prepare(query string) (*sql.Stmt, error) {
	return db.conn.Prepare(query)
}

// QueryRow executa uma consulta que retorna uma única linha
func (db *PostgresDB) QueryRow(query string, args ...interface{}) *sql.Row {
	return db.conn.QueryRow(query, args...)
}
