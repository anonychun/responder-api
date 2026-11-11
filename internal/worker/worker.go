package worker

import (
	"context"
	"log/slog"

	"github.com/anonychun/responder-api/internal/bootstrap"
	clientLark "github.com/anonychun/responder-api/internal/client/lark"
	clientRiver "github.com/anonychun/responder-api/internal/client/river"
	jobHello "github.com/anonychun/responder-api/internal/job/hello"
	jobLarkReply "github.com/anonychun/responder-api/internal/job/lark_reply"
	usecaseLarkMessage "github.com/anonychun/responder-api/internal/usecase/lark/message"
	"github.com/riverqueue/river"
	"github.com/samber/do/v2"
)

func init() {
	do.Provide(bootstrap.Injector, NewWorker)
}

type IWorker interface {
	Start(ctx context.Context) error
}

type Worker struct {
	riverClient clientRiver.IClient
	larkClient  clientLark.IClient
	wsHandler   usecaseLarkMessage.IWsHandler
}

var _ IWorker = (*Worker)(nil)

func NewWorker(i do.Injector) (*Worker, error) {
	riverClient := do.MustInvoke[*clientRiver.Client](i)

	err := addWorkers(riverClient.Workers(),
		do.MustInvoke[*jobHello.Job](i),
	)
	if err != nil {
		return nil, err
	}

	err = addWorkers(riverClient.Workers(),
		do.MustInvoke[*jobLarkReply.Job](i),
	)
	if err != nil {
		return nil, err
	}

	return &Worker{
		riverClient: riverClient,
		larkClient:  do.MustInvoke[*clientLark.Client](i),
		wsHandler:   do.MustInvoke[*usecaseLarkMessage.WsHandler](i),
	}, nil
}

func (w *Worker) Start(ctx context.Context) error {
	slog.Info("starting worker")
	err := w.riverClient.Client().Start(ctx)
	if err != nil {
		return err
	}

	err = w.larkClient.StartEventStream(ctx, w.wsHandler.OnP2MessageReceiveV1)
	if err != nil {
		return err
	}

	<-w.riverClient.Client().Stopped()
	return nil
}

func (w *Worker) Shutdown(ctx context.Context) error {
	slog.Info("shutting down worker")
	return w.riverClient.Client().Stop(ctx)
}

func addWorkers[T river.JobArgs](workers *river.Workers, jobs ...river.Worker[T]) error {
	for _, job := range jobs {
		err := river.AddWorkerSafely(workers, job)
		if err != nil {
			return err
		}
	}

	return nil
}
