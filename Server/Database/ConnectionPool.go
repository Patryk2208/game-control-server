package Database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"time"
)

type ConnectionConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
	//connection pool parameters
	MaxConnections  int32
	MinConnections  int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

type DBConnectionPool struct {
	pool *pgxpool.Pool
}

func CreateConnectionString(c ConnectionConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", c.Username, c.Password, c.Host, c.Port, c.Database)
}

func InitConnectionPool(ctx context.Contex, cc ConnectionConfig) (*DBConnectionPool, error) {
	connString := CreateConnectionString(cc)
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	config.MaxConns = cc.MaxConnections
	config.MinConns = cc.MinConnections
	config.MaxConnIdleTime = cc.MaxConnIdleTime
	config.MaxConnLifetime = cc.MaxConnLifetime

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping failed: %w", err)
	}
	return &DBConnectionPool{pool}, nil
}

func CloseConnectionPool(pool *DBConnectionPool) {
	pool.pool.Close()
}

func Connection() {
	user := "patryk"
	password := "sql"
	ip := "172.17.0.2"
	port := 5432
	name := "users"
	connStr :=
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), "create table players;")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create table: %v\n", err)
	}
	defer rows.Close()
	for rows.Next() {
		fmt.Println(rows.Scan())
	}
}
