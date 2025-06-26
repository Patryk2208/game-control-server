package Database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
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

type AcquiredConnection struct {
	conn *pgxpool.Conn
}

func CreateConnectionString(c ConnectionConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", c.Username, c.Password, c.Host, c.Port, c.Database)
}

func InitConnectionPool(ctx context.Context, cc ConnectionConfig) (*DBConnectionPool, error) {
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

func ReleaseConnection(c *AcquiredConnection) {
	c.conn.Release()
}
