-- +goose Up
-- +goose StatementBegin
CREATE TABLE logs (
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    level      VARCHAR(16) NOT NULL,      
    message    TEXT        NOT NULL
   
);;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS logs;
-- +goose StatementEnd
