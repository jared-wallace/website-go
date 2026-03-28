-- +goose Up
CREATE TABLE posts (
    id            BIGSERIAL    PRIMARY KEY,
    title         TEXT         NOT NULL,
    slug          TEXT         NOT NULL UNIQUE,
    body          TEXT         NOT NULL DEFAULT '',
    rendered_html TEXT         NOT NULL DEFAULT '',
    published     BOOLEAN      NOT NULL DEFAULT false,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),
    deleted_at    TIMESTAMPTZ
);

CREATE INDEX idx_posts_slug ON posts (slug);
CREATE INDEX idx_posts_published ON posts (published) WHERE deleted_at IS NULL;

-- +goose Down
DROP TABLE IF EXISTS posts;
