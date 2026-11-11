-- +goose Up
-- +goose StatementBegin
CREATE TABLE attachments (
	id UUID PRIMARY KEY DEFAULT uuidv7(),
	object_name TEXT NOT NULL UNIQUE,
	file_name TEXT NOT NULL,
	byte_size BIGINT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

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

DROP TABLE attachments;
-- +goose StatementEnd
