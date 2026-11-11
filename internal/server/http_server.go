package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/anonychun/responder-api/internal/api"
	"github.com/anonychun/responder-api/internal/bootstrap"
	"github.com/anonychun/responder-api/internal/config"
	middlewareLogger "github.com/anonychun/responder-api/internal/middleware/logger"
	usecaseHealth "github.com/anonychun/responder-api/internal/usecase/health"
	"github.com/labstack/echo/v5"
	"github.com/samber/do/v2"
)

func init() {
	do.Provide(bootstrap.Injector, NewHttpServer)
}

type IHttpServer interface {
	Start(ctx context.Context) error
}

type HttpServer struct {
	echo   *echo.Echo
	server *http.Server

	loggerMiddleware middlewareLogger.IMiddleware

	healthHttpHandler usecaseHealth.IHttpHandler
}

var _ IHttpServer = (*HttpServer)(nil)

func NewHttpServer(i do.Injector) (*HttpServer, error) {
	e := echo.NewWithConfig(echo.Config{
		Logger:           slog.Default(),
		HTTPErrorHandler: api.HttpErrorHandler,
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Get().Http.Port),
		Handler: e,
	}

	return &HttpServer{
		echo:   e,
		server: srv,

		loggerMiddleware: do.MustInvoke[*middlewareLogger.Middleware](i),

		healthHttpHandler: do.MustInvoke[*usecaseHealth.HttpHandler](i),
	}, nil
}

func (s *HttpServer) Start(ctx context.Context) error {
	err := s.routes()
	if err != nil {
		return err
	}

	slog.Info("starting http server", slog.String("addr", s.server.Addr))
	err = s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *HttpServer) Shutdown(ctx context.Context) error {
	slog.Info("shutting down http server")
	return s.server.Shutdown(ctx)
}
