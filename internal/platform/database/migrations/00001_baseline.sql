-- +goose Up
CREATE TABLE users (
    user_id     BIGSERIAL    PRIMARY KEY,
    nickname    VARCHAR(30)  NOT NULL,
    email       VARCHAR(254) NOT NULL,
    discord_id  VARCHAR(25),
    password    VARCHAR(100) NOT NULL,
    event_count BIGINT       NOT NULL DEFAULT 0,
    guide_count BIGINT       NOT NULL DEFAULT 0,
    role        VARCHAR(15)  NOT NULL DEFAULT 'iniciante',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ,
    CONSTRAINT chk_users_role CHECK (role IN ('iniciante', 'veterano', 'admin'))
);
CREATE UNIQUE INDEX idx_users_email ON users (email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_role ON users (role);
CREATE INDEX idx_users_deleted_at ON users (deleted_at);

CREATE TABLE characters (
    character_id BIGSERIAL    PRIMARY KEY,
    name         VARCHAR(100) NOT NULL,
    description  VARCHAR(255) NOT NULL,
    image_url    VARCHAR(255),
    type         VARCHAR(30)  NOT NULL,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT now(),
    deleted_at   TIMESTAMPTZ
);
CREATE UNIQUE INDEX idx_characters_name ON characters (LOWER(name)) WHERE deleted_at IS NULL;
CREATE INDEX idx_characters_deleted_at ON characters (deleted_at);

CREATE TABLE guides (
    guide_id     BIGSERIAL     PRIMARY KEY,
    title        VARCHAR(254)  NOT NULL,
    banner_url   VARCHAR(255),
    type         VARCHAR(15),
    description  VARCHAR(2083) NOT NULL,
    link         VARCHAR(2083),
    likes        BIGINT        NOT NULL DEFAULT 0,
    dislikes     BIGINT        NOT NULL DEFAULT 0,
    created_at   TIMESTAMPTZ   NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ   NOT NULL DEFAULT now(),
    deleted_at   TIMESTAMPTZ,
    user_id      BIGINT        NOT NULL,
    character_id BIGINT        NOT NULL,
    CONSTRAINT fk_users_guides FOREIGN KEY (user_id)
        REFERENCES users (user_id) ON UPDATE CASCADE,
    CONSTRAINT fk_characters_guides FOREIGN KEY (character_id)
        REFERENCES characters (character_id) ON UPDATE CASCADE
);
CREATE INDEX idx_guides_user_id ON guides (user_id);
CREATE INDEX idx_guides_character_id ON guides (character_id);
CREATE INDEX idx_guides_deleted_at ON guides (deleted_at);

CREATE TABLE events (
    event_id    BIGSERIAL    PRIMARY KEY,
    title       VARCHAR(75)  NOT NULL,
    description VARCHAR(255) NOT NULL,
    banner_url  VARCHAR(255),
    day         DATE         NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ,
    user_id     BIGINT       NOT NULL,
    CONSTRAINT fk_users_events FOREIGN KEY (user_id)
        REFERENCES users (user_id) ON UPDATE CASCADE
);
CREATE INDEX idx_events_user_id ON events (user_id);
CREATE INDEX idx_events_day ON events (day);
CREATE INDEX idx_events_deleted_at ON events (deleted_at);

CREATE TABLE last_events (
    last_event_id BIGSERIAL   PRIMARY KEY,
    title         VARCHAR(75) NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at    TIMESTAMPTZ,
    event_id      BIGINT      NOT NULL,
    CONSTRAINT fk_events_last_events FOREIGN KEY (event_id)
        REFERENCES events (event_id) ON UPDATE CASCADE
);
CREATE INDEX idx_last_events_event_id ON last_events (event_id);
CREATE INDEX idx_last_events_deleted_at ON last_events (deleted_at);

-- +goose Down
DROP TABLE IF EXISTS last_events;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS guides;
DROP TABLE IF EXISTS characters;
DROP TABLE IF EXISTS users;