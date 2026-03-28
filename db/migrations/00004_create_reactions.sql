-- +goose Up
CREATE TABLE reactions (
    id         BIGSERIAL   PRIMARY KEY,
    post_id    BIGINT      NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    ip_hash    TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- One reaction per IP per post enforced at DB level
CREATE UNIQUE INDEX reactions_post_ip_uidx ON reactions (post_id, ip_hash);
-- Fast count lookup by post
CREATE INDEX reactions_post_id_idx ON reactions (post_id);

-- +goose Down
DROP TABLE IF EXISTS reactions;
