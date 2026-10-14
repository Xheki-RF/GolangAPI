-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    email TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    last_read_id BIGINT DEFAULT -1,
    statistic JSONB NOT NULL DEFAULT '{
        "total": 0,
        "by_month": {
            "jan": 0, "feb": 0, "mar": 0, "apr": 0,
            "may": 0, "jun": 0, "jul": 0, "aug": 0,
            "sep": 0, "oct": 0, "nov": 0, "dec": 0
        },
        "by_days": {
            "mon": 0, "tue": 0, "wed": 0, "thu": 0,
            "fri": 0, "sat": 0, "sun": 0
        }
    }'::jsonb
);

INSERT INTO
    users (id, email)
VALUES (
        '00000000-0000-0000-0000-000000000001',
        'deleted-user@system.local'
    );

-- +goose Down
DROP TABLE users;