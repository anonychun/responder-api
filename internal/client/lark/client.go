package lark

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/anonychun/responder-api/internal/bootstrap"
	"github.com/anonychun/responder-api/internal/config"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	larkws "github.com/larksuite/oapi-sdk-go/v3/ws"
	"github.com/samber/do/v2"
)

func init() {
	do.Provide(bootstrap.Injector, NewClient)
}

var ErrMissingCredentials = errors.New("lark: LARK_APP_ID and LARK_APP_SECRET must be set")

type MessageReceiveHandler func(ctx context.Context, event *larkim.P2MessageReceiveV1) error

type IClient interface {
	ReplyInThread(ctx context.Context, messageId string, content string) error
	StartEventStream(ctx context.Context, handler MessageReceiveHandler) error
}

type Client struct {
	appId      string
	appSecret  string
	larkClient *lark.Client
}

var _ IClient = (*Client)(nil)

func NewClient(i do.Injector) (*Client, error) {
	cfg := config.Get()
	if cfg.Lark.AppId == "" || cfg.Lark.AppSecret == "" {
		return nil, ErrMissingCredentials
	}

	return &Client{
		appId:      cfg.Lark.AppId,
		appSecret:  cfg.Lark.AppSecret,
		larkClient: lark.NewClient(cfg.Lark.AppId, cfg.Lark.AppSecret),
	}, nil
}

func (c *Client) ReplyInThread(ctx context.Context, messageId string, content string) error {
	body, err := textMessageBody(content)
	if err != nil {
		return err
	}

	req := larkim.NewReplyMessageReqBuilder().
		MessageId(messageId).
		Body(larkim.NewReplyMessageReqBodyBuilder().
			Content(body).
			MsgType(larkim.MsgTypeText).
			ReplyInThread(true).
			Build()).
		Build()

	resp, err := c.larkClient.Im.Message.Reply(ctx, req)
	if err != nil {
		return err
	}

	if !resp.Success() {
		return resp
	}

	return nil
}

type textMessage struct {
	Text string `json:"text"`
}

// The SDK's MessageTextBuilder does not escape the text it is given, so build the body here.
func textMessageBody(text string) (string, error) {
	body, err := json.Marshal(textMessage{Text: text})
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (c *Client) StartEventStream(ctx context.Context, handler MessageReceiveHandler) error {
	eventHandler := dispatcher.NewEventDispatcher("", "").
		OnP2MessageReceiveV1(handler)

	client := larkws.NewClient(c.appId, c.appSecret,
		larkws.WithDomain(lark.LarkBaseUrl),
		larkws.WithEventHandler(eventHandler),
		larkws.WithLogLevel(larkcore.LogLevelDebug),
	)

	return client.Start(ctx)
}
