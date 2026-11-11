package message

import (
	"context"
	"log/slog"

	"github.com/anonychun/responder-api/internal/bootstrap"
	clientRiver "github.com/anonychun/responder-api/internal/client/river"
	"github.com/anonychun/responder-api/internal/entity"
	jobLarkReply "github.com/anonychun/responder-api/internal/job/lark_reply"
	repositoryLarkMessage "github.com/anonychun/responder-api/internal/repository/lark_message"
	"github.com/samber/do/v2"
)

func init() {
	do.Provide(bootstrap.Injector, NewUsecase)
}

type IUsecase interface {
	MessageReceive(ctx context.Context, req MessageReceiveRequest) error
}

type Usecase struct {
	riverClient           clientRiver.IClient
	larkMessageRepository repositoryLarkMessage.IRepository
}

var _ IUsecase = (*Usecase)(nil)

func NewUsecase(i do.Injector) (*Usecase, error) {
	return &Usecase{
		riverClient:           do.MustInvoke[*clientRiver.Client](i),
		larkMessageRepository: do.MustInvoke[*repositoryLarkMessage.Repository](i),
	}, nil
}

func (u *Usecase) MessageReceive(ctx context.Context, req MessageReceiveRequest) error {
	if !req.IsRootMessage() || !req.MentionsSupportTeam() {
		return nil
	}

	exists, err := u.larkMessageRepository.ExistsByMessageId(ctx, req.MessageId)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	larkMessage := &entity.LarkMessage{
		MessageId: req.MessageId,
		ChatId:    req.ChatId,
		SenderId:  req.SenderId,
		Content:   req.Content,
	}

	err = u.larkMessageRepository.Create(ctx, larkMessage)
	if err != nil {
		return err
	}

	slog.InfoContext(ctx, "lark message stored", slog.String("message_id", larkMessage.MessageId))

	_, err = u.riverClient.Insert(ctx, jobLarkReply.Args{MessageId: larkMessage.MessageId})

	return err
}
