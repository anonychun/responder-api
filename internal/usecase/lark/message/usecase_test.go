package message

import (
	"context"
	"errors"
	"testing"

	clientRiver "github.com/anonychun/responder-api/internal/client/river"
	"github.com/anonychun/responder-api/internal/entity"
	jobLarkReply "github.com/anonychun/responder-api/internal/job/lark_reply"
	repositoryLarkMessage "github.com/anonychun/responder-api/internal/repository/lark_message"
	"github.com/riverqueue/river/rivertype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUsecase_HandleMessageReceive(t *testing.T) {
	t.Run("ignores a reply inside a thread", func(t *testing.T) {
		ctx := context.Background()
		req := MessageReceiveRequest{
			MessageId: "om_2",
			RootId:    "om_1",
			Content:   `{"text":"@_user_1 please help"}`,
			Mentions:  []string{supportTeamHandle},
		}
		usecase := &Usecase{}

		err := usecase.HandleMessageReceive(ctx, req)

		require.NoError(t, err)
	})

	t.Run("ignores a message that does not mention the support team", func(t *testing.T) {
		ctx := context.Background()
		req := MessageReceiveRequest{
			MessageId: "om_1",
			Content:   `{"text":"hello everyone"}`,
		}
		usecase := &Usecase{}

		err := usecase.HandleMessageReceive(ctx, req)

		require.NoError(t, err)
	})

	t.Run("stores the message and queues a reply job", func(t *testing.T) {
		ctx := context.Background()
		req := MessageReceiveRequest{
			MessageId: "om_1",
			ChatId:    "oc_1",
			SenderId:  "ou_1",
			Content:   `{"text":"@_user_1 please help"}`,
			Mentions:  []string{supportTeamHandle},
		}
		riverClient := clientRiver.NewMockIClient(t)
		larkMessageRepository := repositoryLarkMessage.NewMockIRepository(t)
		usecase := &Usecase{
			riverClient:           riverClient,
			larkMessageRepository: larkMessageRepository,
		}

		var stored *entity.LarkMessage
		larkMessageRepository.EXPECT().CreateIfNotExists(ctx, mock.AnythingOfType("*entity.LarkMessage")).Run(func(ctx context.Context, actual *entity.LarkMessage) {
			stored = actual
		}).Return(true, nil).Once()
		riverClient.EXPECT().Insert(ctx, jobLarkReply.Args{MessageId: req.MessageId}).Return(&rivertype.JobInsertResult{}, nil).Once()

		err := usecase.HandleMessageReceive(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, stored)
		assert.Equal(t, req.MessageId, stored.MessageId)
		assert.Equal(t, req.ChatId, stored.ChatId)
		assert.Equal(t, req.SenderId, stored.SenderId)
		assert.Equal(t, req.Content, stored.Content)
	})

	t.Run("queues the reply again for a message that is already stored", func(t *testing.T) {
		ctx := context.Background()
		req := MessageReceiveRequest{
			MessageId: "om_1",
			Content:   `{"text":"@_user_1 please help"}`,
			Mentions:  []string{supportTeamHandle},
		}
		riverClient := clientRiver.NewMockIClient(t)
		larkMessageRepository := repositoryLarkMessage.NewMockIRepository(t)
		usecase := &Usecase{
			riverClient:           riverClient,
			larkMessageRepository: larkMessageRepository,
		}

		larkMessageRepository.EXPECT().CreateIfNotExists(ctx, mock.AnythingOfType("*entity.LarkMessage")).Return(false, nil).Once()
		riverClient.EXPECT().Insert(ctx, jobLarkReply.Args{MessageId: req.MessageId}).Return(&rivertype.JobInsertResult{UniqueSkippedAsDuplicate: true}, nil).Once()

		err := usecase.HandleMessageReceive(ctx, req)

		require.NoError(t, err)
	})

	t.Run("returns the error from storing the message", func(t *testing.T) {
		ctx := context.Background()
		req := MessageReceiveRequest{
			MessageId: "om_1",
			Content:   `{"text":"@_user_1 please help"}`,
			Mentions:  []string{supportTeamHandle},
		}
		expectedErr := errors.New("insert lark message")
		larkMessageRepository := repositoryLarkMessage.NewMockIRepository(t)
		usecase := &Usecase{larkMessageRepository: larkMessageRepository}

		larkMessageRepository.EXPECT().CreateIfNotExists(ctx, mock.AnythingOfType("*entity.LarkMessage")).Return(false, expectedErr).Once()

		err := usecase.HandleMessageReceive(ctx, req)

		require.ErrorIs(t, err, expectedErr)
	})

	t.Run("returns the error from queueing the reply", func(t *testing.T) {
		ctx := context.Background()
		req := MessageReceiveRequest{
			MessageId: "om_1",
			Content:   `{"text":"@_user_1 please help"}`,
			Mentions:  []string{supportTeamHandle},
		}
		expectedErr := errors.New("insert job")
		riverClient := clientRiver.NewMockIClient(t)
		larkMessageRepository := repositoryLarkMessage.NewMockIRepository(t)
		usecase := &Usecase{
			riverClient:           riverClient,
			larkMessageRepository: larkMessageRepository,
		}

		larkMessageRepository.EXPECT().CreateIfNotExists(ctx, mock.AnythingOfType("*entity.LarkMessage")).Return(true, nil).Once()
		riverClient.EXPECT().Insert(ctx, jobLarkReply.Args{MessageId: req.MessageId}).Return(nil, expectedErr).Once()

		err := usecase.HandleMessageReceive(ctx, req)

		require.ErrorIs(t, err, expectedErr)
	})
}
