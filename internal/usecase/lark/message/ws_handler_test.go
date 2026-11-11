package message

import (
	"context"
	"errors"
	"testing"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWsHandler_OnP2MessageReceiveV1(t *testing.T) {
	t.Run("passes the message on to the usecase", func(t *testing.T) {
		ctx := context.Background()
		event := &larkim.P2MessageReceiveV1{
			Event: &larkim.P2MessageReceiveV1Data{
				Sender: &larkim.EventSender{
					SenderId: &larkim.UserId{UserId: stringPointer("on_1"), OpenId: stringPointer("ou_1")},
				},
				Message: &larkim.EventMessage{
					MessageId: stringPointer("om_1"),
					ChatId:    stringPointer("oc_1"),
					Content:   stringPointer(`{"text":"@_user_1 please help"}`),
					Mentions:  []*larkim.MentionEvent{{Name: stringPointer(supportTeamHandle)}},
				},
			},
		}
		usecase := NewMockIUsecase(t)
		handler := &WsHandler{usecase: usecase}

		usecase.EXPECT().MessageReceive(ctx, MessageReceiveRequest{
			MessageId: "om_1",
			ChatId:    "oc_1",
			SenderId:  "on_1",
			Content:   `{"text":"@_user_1 please help"}`,
			Mentions:  []string{supportTeamHandle},
		}).Return(nil).Once()

		err := handler.OnP2MessageReceiveV1(ctx, event)

		require.NoError(t, err)
	})

	t.Run("falls back to the open id when the user id is missing", func(t *testing.T) {
		ctx := context.Background()
		event := &larkim.P2MessageReceiveV1{
			Event: &larkim.P2MessageReceiveV1Data{
				Sender: &larkim.EventSender{
					SenderId: &larkim.UserId{OpenId: stringPointer("ou_1")},
				},
				Message: &larkim.EventMessage{
					MessageId: stringPointer("om_1"),
				},
			},
		}
		usecase := NewMockIUsecase(t)
		handler := &WsHandler{usecase: usecase}

		usecase.EXPECT().MessageReceive(ctx, MessageReceiveRequest{
			MessageId: "om_1",
			SenderId:  "ou_1",
		}).Return(nil).Once()

		err := handler.OnP2MessageReceiveV1(ctx, event)

		require.NoError(t, err)
	})

	t.Run("ignores an event without a message", func(t *testing.T) {
		handler := &WsHandler{}

		assert.NoError(t, handler.OnP2MessageReceiveV1(context.Background(), nil))
		assert.NoError(t, handler.OnP2MessageReceiveV1(context.Background(), &larkim.P2MessageReceiveV1{}))
	})

	t.Run("returns the usecase error so that Lark redelivers the event", func(t *testing.T) {
		ctx := context.Background()
		event := &larkim.P2MessageReceiveV1{
			Event: &larkim.P2MessageReceiveV1Data{
				Message: &larkim.EventMessage{MessageId: stringPointer("om_1")},
			},
		}
		expectedErr := errors.New("store message")
		usecase := NewMockIUsecase(t)
		handler := &WsHandler{usecase: usecase}

		usecase.EXPECT().MessageReceive(ctx, MessageReceiveRequest{MessageId: "om_1"}).Return(expectedErr).Once()

		err := handler.OnP2MessageReceiveV1(ctx, event)

		require.ErrorIs(t, err, expectedErr)
	})
}

func stringPointer(value string) *string {
	return &value
}
