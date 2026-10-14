-- +goose Up
CREATE TYPE rule_enum AS ENUM ('regexp', 'goja');

CREATE TABLE rules (
    id BIGSERIAL PRIMARY KEY,
    condition_type rule_enum NOT NULL,
    condition_value TEXT NOT NULL UNIQUE,
    action_type rule_enum NOT NULL,
    action_find TEXT,
    action_replace TEXT,
    action_goja_code TEXT,
    weight INT UNIQUE NOT NULL
);

INSERT INTO
    rules (
        id,
        condition_type,
        condition_value,
        action_type,
        action_find,
        action_replace,
        action_goja_code,
        weight
    )
VALUES (
        0,
        'regexp',
        '[бБbB]',
        'regexp',
        '[бБbB]',
        '',
        '',
        0
    );

-- +goose Down
DROP TABLE rules;

DROP TYPE rule_enum;