package message

import (
	"context"

	"github.com/anonychun/responder-api/internal/bootstrap"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/samber/do/v2"
)

func init() {
	do.Provide(bootstrap.Injector, NewWsHandler)
}

type IWsHandler interface {
	OnP2MessageReceiveV1(ctx context.Context, event *larkim.P2MessageReceiveV1) error
}

type WsHandler struct {
	usecase IUsecase
}

var _ IWsHandler = (*WsHandler)(nil)

func NewWsHandler(i do.Injector) (*WsHandler, error) {
	return &WsHandler{
		usecase: do.MustInvoke[*Usecase](i),
	}, nil
}

func (h *WsHandler) OnP2MessageReceiveV1(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	req := newMessageReceiveRequest(event)
	if req.MessageId == "" {
		return nil
	}

	return h.usecase.MessageReceive(ctx, req)
}

func newMessageReceiveRequest(event *larkim.P2MessageReceiveV1) MessageReceiveRequest {
	req := MessageReceiveRequest{}
	if event == nil || event.Event == nil {
		return req
	}

	if message := event.Event.Message; message != nil {
		req.MessageId = stringValue(message.MessageId)
		req.ChatId = stringValue(message.ChatId)
		req.RootId = stringValue(message.RootId)
		req.ParentId = stringValue(message.ParentId)
		req.Content = stringValue(message.Content)

		for _, mention := range message.Mentions {
			if mention == nil {
				continue
			}

			req.Mentions = append(req.Mentions, stringValue(mention.Name))
		}
	}

	if sender := event.Event.Sender; sender != nil && sender.SenderId != nil {
		req.SenderId = firstNonEmpty(stringValue(sender.SenderId.UserId), stringValue(sender.SenderId.OpenId))
	}

	return req
}
