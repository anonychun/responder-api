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

func TestRepository_Create(t *testing.T) {
	t.Run("inserts the message and fills in the generated columns", func(t *testing.T) {
		ctx := t.Context()
		repository := &Repository{sqlDB: do.MustInvoke[*dbSql.PostgresDB](bootstrap.Injector)}
		larkMessage := &entity.LarkMessage{
			MessageId: uuid.New().String(),
			ChatId:    "oc_1",
			SenderId:  "ou_1",
			Content:   `{"text":"@_user_1 please help"}`,
		}

		err := repository.Create(ctx, larkMessage)

		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil(), larkMessage.Id)
		assert.False(t, larkMessage.CreatedAt.IsZero())
		assert.False(t, larkMessage.UpdatedAt.IsZero())
	})

	t.Run("rejects a second message with the same message id", func(t *testing.T) {
		ctx := t.Context()
		repository := &Repository{sqlDB: do.MustInvoke[*dbSql.PostgresDB](bootstrap.Injector)}
		messageId := uuid.New().String()
		require.NoError(t, repository.Create(ctx, &entity.LarkMessage{MessageId: messageId, ChatId: "oc_1", SenderId: "ou_1", Content: "first"}))

		err := repository.Create(ctx, &entity.LarkMessage{MessageId: messageId, ChatId: "oc_1", SenderId: "ou_1", Content: "second"})

		require.Error(t, err)
	})
}

func TestRepository_ExistsByMessageId(t *testing.T) {
	t.Run("reports whether a message has that id", func(t *testing.T) {
		ctx := t.Context()
		repository := &Repository{sqlDB: do.MustInvoke[*dbSql.PostgresDB](bootstrap.Injector)}
		larkMessage := &entity.LarkMessage{MessageId: uuid.New().String(), ChatId: "oc_1", SenderId: "ou_1", Content: "content"}
		require.NoError(t, repository.Create(ctx, larkMessage))

		exists, err := repository.ExistsByMessageId(ctx, larkMessage.MessageId)
		require.NoError(t, err)
		assert.True(t, exists)

		exists, err = repository.ExistsByMessageId(ctx, uuid.New().String())
		require.NoError(t, err)
		assert.False(t, exists)
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
		require.NoError(t, repository.Create(ctx, larkMessage))

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
		require.NoError(t, repository.Create(ctx, larkMessage))

		err := repository.MarkReplied(ctx, larkMessage.Id)

		require.NoError(t, err)
		actualLarkMessage, err := repository.FindByMessageId(ctx, larkMessage.MessageId)
		require.NoError(t, err)
		assert.True(t, actualLarkMessage.IsReplied())
	})
}
