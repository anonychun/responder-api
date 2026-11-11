-- +goose Up
-- +goose StatementBegin
CREATE TABLE lark_messages (
	id UUID PRIMARY KEY DEFAULT uuidv7(),
	message_id TEXT NOT NULL UNIQUE,
	chat_id TEXT NOT NULL,
	sender_id TEXT NOT NULL,
	content TEXT NOT NULL,
	replied_at TIMESTAMPTZ,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE lark_messages;
-- +goose StatementEnd
