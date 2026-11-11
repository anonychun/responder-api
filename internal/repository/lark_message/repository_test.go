package lark_message

import (
	"database/sql"
	"testing"
	"uuid"

	"github.com/anonychun/responder-api/internal/bootstrap"
	dbSql "github.com/anonychun/responder-api/internal/db/sql"
	"github.com/anonychun/responder-api/internal/entity"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository_CreateIfNotExists(t *testing.T) {
	t.Run("inserts the message and fills in the generated columns", func(t *testing.T) {
		ctx := t.Context()
		repository := &Repository{sqlDB: do.MustInvoke[*dbSql.PostgresDB](bootstrap.Injector)}
		larkMessage := &entity.LarkMessage{
			MessageId: uuid.New().String(),
			ChatId:    "oc_1",
			SenderId:  "ou_1",
			Content:   `{"text":"@_user_1 please help"}`,
		}

		created, err := repository.CreateIfNotExists(ctx, larkMessage)

		require.NoError(t, err)
		assert.True(t, created)
		assert.NotEqual(t, uuid.Nil(), larkMessage.Id)
		assert.False(t, larkMessage.CreatedAt.IsZero())
		assert.False(t, larkMessage.UpdatedAt.IsZero())
	})

	t.Run("reports a message that is already stored as not created", func(t *testing.T) {
		ctx := t.Context()
		repository := &Repository{sqlDB: do.MustInvoke[*dbSql.PostgresDB](bootstrap.Injector)}
		messageId := uuid.New().String()
		created, err := repository.CreateIfNotExists(ctx, &entity.LarkMessage{MessageId: messageId, ChatId: "oc_1", SenderId: "ou_1", Content: "first"})
		require.NoError(t, err)
		require.True(t, created)

		created, err = repository.CreateIfNotExists(ctx, &entity.LarkMessage{MessageId: messageId, ChatId: "oc_1", SenderId: "ou_1", Content: "second"})

		require.NoError(t, err)
		assert.False(t, created)

		larkMessage, err := repository.FindByMessageId(ctx, messageId)
		require.NoError(t, err)
		assert.Equal(t, "first", larkMessage.Content)
	})
}

func TestRepository_FindByMessageId(t *testing.T) {
	t.Run("returns the message with that id", func(t *testing.T) {
		ctx := t.Context()
		repository := &Repository{sqlDB: do.MustInvoke[*dbSql.PostgresDB](bootstrap.Injector)}
		larkMessage := &entity.LarkMessage{
			MessageId: uuid.New().String(),
			ChatId:    "oc_1",
			SenderId:  "ou_1",
			Content:   `{"text":"@_user_1 please help"}`,
		}
		_, err := repository.CreateIfNotExists(ctx, larkMessage)
		require.NoError(t, err)

		actualLarkMessage, err := repository.FindByMessageId(ctx, larkMessage.MessageId)

		require.NoError(t, err)
		assert.Equal(t, larkMessage.Id, actualLarkMessage.Id)
		assert.Equal(t, larkMessage.ChatId, actualLarkMessage.ChatId)
		assert.Equal(t, larkMessage.SenderId, actualLarkMessage.SenderId)
		assert.Equal(t, larkMessage.Content, actualLarkMessage.Content)
		assert.True(t, actualLarkMessage.RepliedAt.IsZero())
	})

	t.Run("returns sql.ErrNoRows when no message has that id", func(t *testing.T) {
		ctx := t.Context()
		repository := &Repository{sqlDB: do.MustInvoke[*dbSql.PostgresDB](bootstrap.Injector)}

		actualLarkMessage, err := repository.FindByMessageId(ctx, uuid.New().String())

		require.ErrorIs(t, err, sql.ErrNoRows)
		assert.Nil(t, actualLarkMessage)
	})
}

func TestRepository_MarkReplied(t *testing.T) {
	t.Run("sets the reply timestamp", func(t *testing.T) {
		ctx := t.Context()
		repository := &Repository{sqlDB: do.MustInvoke[*dbSql.PostgresDB](bootstrap.Injector)}
		larkMessage := &entity.LarkMessage{MessageId: uuid.New().String(), ChatId: "oc_1", SenderId: "ou_1", Content: "content"}
		_, err := repository.CreateIfNotExists(ctx, larkMessage)
		require.NoError(t, err)

		err = repository.MarkReplied(ctx, larkMessage.Id)

		require.NoError(t, err)
		actualLarkMessage, err := repository.FindByMessageId(ctx, larkMessage.MessageId)
		require.NoError(t, err)
		assert.True(t, actualLarkMessage.IsReplied())
	})
}
