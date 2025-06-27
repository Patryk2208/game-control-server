package Database

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type PlayerDB struct {
	Player_id       uuid.UUID
	Username        string
	Password_hash   []byte
	Created_at      pgtype.Timestamptz
	Updated_at      pgtype.Timestamptz
	Last_login      pgtype.Timestamptz
	Status          string
	Suspended_until pgtype.Timestamptz
}

type PlayerProfileDB struct {
	Player_id  uuid.UUID
	Nick       string
	Experience int
}

type SessionDB struct {
	Session_id    uuid.UUID
	Player_id     uuid.UUID
	Created_at    pgtype.Timestamptz
	Expires_at    pgtype.Timestamptz
	Last_activity pgtype.Timestamptz
}

type GameDB struct {
	Game_id    uuid.UUID
	Player_id  uuid.UUID
	Start_time pgtype.Timestamptz
	End_time   pgtype.Timestamptz
	Winner     uuid.UUID
}
