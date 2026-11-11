package lark_reply

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/anonychun/responder-api/internal/bootstrap"
	clientLark "github.com/anonychun/responder-api/internal/client/lark"
	repositoryLarkMessage "github.com/anonychun/responder-api/internal/repository/lark_message"
	"github.com/riverqueue/river"
	"github.com/samber/do/v2"
)

func init() {
	do.Provide(bootstrap.Injector, NewJob)
}

type Args struct {
	MessageId string
}

func (Args) Kind() string {
	return "lark_reply"
}

type Job struct {
	river.WorkerDefaults[Args]

	larkClient            clientLark.IClient
	larkMessageRepository repositoryLarkMessage.IRepository
}

func NewJob(i do.Injector) (*Job, error) {
	return &Job{
		larkClient:            do.MustInvoke[*clientLark.Client](i),
		larkMessageRepository: do.MustInvoke[*repositoryLarkMessage.Repository](i),
	}, nil
}

func (j *Job) Work(ctx context.Context, job *river.Job[Args]) error {
	larkMessage, err := j.larkMessageRepository.FindByMessageId(ctx, job.Args.MessageId)
	if errors.Is(err, sql.ErrNoRows) {
		return river.JobCancel(fmt.Errorf("lark message %s not found", job.Args.MessageId))
	} else if err != nil {
		return err
	}

	if larkMessage.IsReplied() {
		slog.DebugContext(ctx, "lark message already replied", slog.String("message_id", larkMessage.MessageId))
		return nil
	}

	err = j.larkClient.ReplyInThread(ctx, larkMessage.MessageId, replyContent(larkMessage.SenderId))
	if err != nil {
		return err
	}

	return j.larkMessageRepository.MarkReplied(ctx, larkMessage.Id)
}
