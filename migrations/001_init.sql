CREATE TABLE IF NOT EXISTS comments (
    id            SERIAL PRIMARY KEY,
    parent_id     INT REFERENCES comments(id) ON DELETE CASCADE,
    author        VARCHAR(100) NOT NULL,
    text          TEXT NOT NULL,
    created_at    TIMESTAMP DEFAULT NOW(),
    search_vector TSVECTOR GENERATED ALWAYS AS (
        to_tsvector('russian', text)
    ) STORED
);

CREATE INDEX IF NOT EXISTS idx_comments_parent_id ON comments(parent_id);
CREATE INDEX IF NOT EXISTS idx_comments_created_at ON comments(created_at);
CREATE INDEX IF NOT EXISTS idx_comments_search ON comments USING GIN(search_vector);