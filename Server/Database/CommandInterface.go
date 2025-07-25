package Database

import (
	"Server/UserAuthentication"
	"context"
	"fmt"
)

func (pool *DBConnectionPool) TryRegisterUser(username string, password string) (bool, error) {
	ctx := context.Background()
	conn, err := pool.pool.Acquire(ctx)
	if err != nil {
		return false, err
	}
	c := AcquiredConnection{conn}
	defer ReleaseConnection(&c)
	hash, err := UserAuthentication.HashPassword(password)
	if err != nil {
		return false, err
	}
	return c.TryAddUser(ctx, username, hash)
}

func (pool *DBConnectionPool) TryLogin(username string, password string) (bool, *PlayerDB, error) {
	ctx := context.Background()
	conn, err := pool.pool.Acquire(ctx)
	if err != nil {
		return false, nil, err
	}
	c := AcquiredConnection{conn}
	defer ReleaseConnection(&c)
	player, err := c.TryGetUser(ctx, username)
	if err != nil {
		return false, nil, err
	}
	success, err := UserAuthentication.VerifyPassword(password, string(player.Password_hash))
	if err != nil {
		return false, nil, err
	}
	if !success {
		return false, nil, fmt.Errorf("wrong password")
	}
	success, err = c.TryAddSession(ctx, username)
	if err != nil {
		return false, nil, err
	}
	if !success {
		return false, nil, fmt.Errorf("session error")
	}

	return true, &player, nil
}

func (pool *DBConnectionPool) TryLogout(player *PlayerDB) (bool, error) {
	ctx := context.Background()
	conn, err := pool.pool.Acquire(ctx)
	if err != nil {
		return false, err
	}
	c := AcquiredConnection{conn}
	defer ReleaseConnection(&c)
	success, err := c.TryDeleteSession(ctx, player)
	if err != nil {
		return false, err
	}
	if !success {
		return false, fmt.Errorf("session error")
	}
	return true, nil
}

func (pool *DBConnectionPool) TryStartSession(ctx context.Context, player *PlayerDB) (bool, error) {
	return true, nil //todo
}

func (pool *DBConnectionPool) TryStopSession(ctx context.Context, player *PlayerDB) (bool, error) {
	return true, nil //todo
}
