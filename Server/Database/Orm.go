package Database

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Player struct {
	player_id       uuid.UUID
	username        string
	password_hash   []byte
	created_at      pgtype.Timestamptz
	updated_at      pgtype.Timestamptz
	last_login      pgtype.Timestamptz
	status          string
	suspended_until pgtype.Timestamptz
}
