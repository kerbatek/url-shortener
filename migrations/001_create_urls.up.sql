CREATE TABLE IF NOT EXISTS urls (
    id          BIGSERIAL    PRIMARY KEY,
    code        VARCHAR(20)  NOT NULL UNIQUE,
    original_url TEXT        NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_urls_code ON urls (code);
