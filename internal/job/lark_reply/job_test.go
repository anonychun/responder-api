package lark_reply

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"
	"uuid"

	clientLark "github.com/anonychun/responder-api/internal/client/lark"
	"github.com/anonychun/responder-api/internal/entity"
	repositoryLarkMessage "github.com/anonychun/responder-api/internal/repository/lark_message"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestJob_Work(t *testing.T) {
	t.Run("replies to the message and marks it as replied", func(t *testing.T) {
		ctx := context.Background()
		larkMessage := &entity.LarkMessage{
			Base:      entity.Base{Id: uuid.New()},
			MessageId: "om_1",
			SenderId:  "ou_1",
		}
		larkClient := clientLark.NewMockIClient(t)
		larkMessageRepository := repositoryLarkMessage.NewMockIRepository(t)
		job := &Job{
			larkClient:            larkClient,
			larkMessageRepository: larkMessageRepository,
		}

		larkMessageRepository.EXPECT().FindByMessageId(ctx, "om_1").Return(larkMessage, nil).Once()
		larkClient.EXPECT().ReplyInThread(ctx, "om_1", mock.AnythingOfType("string")).Return(nil).Once()
		larkMessageRepository.EXPECT().MarkReplied(ctx, larkMessage.Id).Return(nil).Once()

		err := job.Work(ctx, newRiverJob("om_1"))

		require.NoError(t, err)
	})

	t.Run("skips a message that already has a reply", func(t *testing.T) {
		ctx := context.Background()
		larkMessage := &entity.LarkMessage{
			MessageId: "om_1",
			RepliedAt: time.Now(),
		}
		larkMessageRepository := repositoryLarkMessage.NewMockIRepository(t)
		job := &Job{larkMessageRepository: larkMessageRepository}

		larkMessageRepository.EXPECT().FindByMessageId(ctx, "om_1").Return(larkMessage, nil).Once()

		err := job.Work(ctx, newRiverJob("om_1"))

		require.NoError(t, err)
	})

	t.Run("cancels the job when the message is gone", func(t *testing.T) {
		ctx := context.Background()
		larkMessageRepository := repositoryLarkMessage.NewMockIRepository(t)
		job := &Job{larkMessageRepository: larkMessageRepository}

		larkMessageRepository.EXPECT().FindByMessageId(ctx, "om_1").Return(nil, sql.ErrNoRows).Once()

		err := job.Work(ctx, newRiverJob("om_1"))

		require.Error(t, err)
		var jobCancelErr *rivertype.JobCancelError
		assert.True(t, errors.As(err, &jobCancelErr))
	})

	t.Run("returns the error from finding the message", func(t *testing.T) {
		ctx := context.Background()
		expectedErr := errors.New("find message")
		larkMessageRepository := repositoryLarkMessage.NewMockIRepository(t)
		job := &Job{larkMessageRepository: larkMessageRepository}

		larkMessageRepository.EXPECT().FindByMessageId(ctx, "om_1").Return(nil, expectedErr).Once()

		err := job.Work(ctx, newRiverJob("om_1"))

		require.ErrorIs(t, err, expectedErr)
	})

	t.Run("does not mark the message as replied when the reply fails", func(t *testing.T) {
		ctx := context.Background()
		expectedErr := errors.New("reply message")
		larkMessage := &entity.LarkMessage{MessageId: "om_1", SenderId: "ou_1"}
		larkClient := clientLark.NewMockIClient(t)
		larkMessageRepository := repositoryLarkMessage.NewMockIRepository(t)
		job := &Job{
			larkClient:            larkClient,
			larkMessageRepository: larkMessageRepository,
		}

		larkMessageRepository.EXPECT().FindByMessageId(ctx, "om_1").Return(larkMessage, nil).Once()
		larkClient.EXPECT().ReplyInThread(ctx, "om_1", mock.AnythingOfType("string")).Return(expectedErr).Once()

		err := job.Work(ctx, newRiverJob("om_1"))

		require.ErrorIs(t, err, expectedErr)
	})
}

func newRiverJob(messageId string) *river.Job[Args] {
	return &river.Job[Args]{Args: Args{MessageId: messageId}}
}
