package Database

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (c *AcquiredConnection) TryAddUser(ctx context.Context, username string, passwordHash string) (bool, error) {
	_, err := c.conn.Exec(ctx,
		"insert into players (username, password_hash) values ($1, cast($2 as bytea));",
		username, passwordHash)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (c *AcquiredConnection) TryDeleteUser(ctx context.Context, username string) (bool, error) {
	commandTag, err := c.conn.Exec(ctx,
		"delete from players where username = $1;",
		username)
	if err != nil {
		return false, err
	}
	if commandTag.RowsAffected() == 0 {
		return false, err
	}
	return true, nil
}

func (c *AcquiredConnection) TryGetUser(ctx context.Context, username string) (PlayerDB, error) {
	rows, err := c.conn.Query(ctx,
		"select * from players where username = $1;",
		username)
	if err != nil || rows.Err() != nil {
		return PlayerDB{}, fmt.Errorf("could not find player by username %s", username)
	}
	var (
		id  uuid.UUID
		u   string
		p   []byte
		cr  pgtype.Timestamptz
		upd pgtype.Timestamptz
		ll  pgtype.Timestamptz
		s   string
		su  pgtype.Timestamptz
	)
	for rows.Next() {
		err = rows.Scan(&id, &u, &p, &cr, &upd, &ll, &s, &su)
		if err != nil {
			return PlayerDB{}, fmt.Errorf("could not find player by username %s", username)
		}
	}
	player := PlayerDB{id, u, p, cr, upd, ll, s, su}
	return player, nil
}
func (c *AcquiredConnection) TryAddSession(ctx context.Context, username string) (bool, error) {
	player, err := c.TryGetUser(ctx, username)
	if err != nil {
		return false, err
	}
	_, err = c.conn.Exec(ctx,
		"insert into sessions (player_id) values ($1);",
		player.Player_id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (c *AcquiredConnection) TryDeleteSession(ctx context.Context, player *PlayerDB) (bool, error) {
	commandTag, err := c.conn.Exec(ctx,
		"delete from sessions where player_id = $1;",
		player.Player_id)
	if err != nil {
		return false, err
	}
	if commandTag.RowsAffected() == 0 {
		return false, err
	}
	return true, nil
}
