CREATE TABLE players (
    player_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(32) NOT NULL UNIQUE CHECK (LENGTH(username) >= 2),
    password_hash BYTEA NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_login TIMESTAMPTZ,
    status VARCHAR(16) NOT NULL DEFAULT 'active'
        CHECK (status IN ('active', 'banned', 'suspended') AND
        (CASE WHEN status='suspended' THEN suspended_until != NULL ELSE players.suspended_until = NULL END)),
    suspended_until TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE playerProfiles (
    player_id UUID PRIMARY KEY REFERENCES players(player_id) ON DELETE CASCADE,
    nick VARCHAR(32),
    experience INT NOT NULL DEFAULT 0 CHECK ( experience >= 0 )
);

CREATE TABLE games (
    game_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID REFERENCES players(player_id) ON DELETE CASCADE,
    start_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    end_time TIMESTAMPTZ DEFAULT NULL CHECK ( end_time > start_time ),
    winner UUID REFERENCES players(player_id)
);

CREATE TABLE gamePlayed (
    game_played_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID REFERENCES games(game_id) ON DELETE CASCADE ,
    player_id UUID REFERENCES players(player_id) ON DELETE CASCADE,
    spawn_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    defeat_time TIMESTAMPTZ DEFAULT NULL CHECK ( defeat_time > spawn_time ),
    position INT CHECK ( position > 0 )
);

CREATE TABLE sessions (
    session_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(player_id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL DEFAULT NOW() + INTERVAL '5 hours' CHECK ( expires_at > created_at ),
    last_activity TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE securityEvents (
    event_id BIGSERIAL PRIMARY KEY,
    player_id UUID REFERENCES players(player_id) ON DELETE SET NULL,
    event_type VARCHAR(32) NOT NULL
     CHECK (event_type IN ('registration', 'login', 'logout', 'password_change', 'session_revoked', 'login_fail', 'sus_act')),
    time TIMESTAMPTZ NOT NULL DEFAULT NOW()
);