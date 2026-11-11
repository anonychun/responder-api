package server

import (
	"github.com/anonychun/responder-api/public"
	echootel "github.com/labstack/echo-otel/v5"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
)

func (s *HttpServer) routes() error {
	s.echo.Use(middleware.Recover())
	s.echo.Use(middleware.RequestID())
	s.echo.Use(s.loggerMiddleware.RequestLogger)

	s.echo.Use(echootel.NewMiddlewareWithConfig(echootel.Config{
		TracerProvider: otel.GetTracerProvider(),
		MeterProvider:  otel.GetMeterProvider(),
	}))

	s.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://127.0.0.1:3000",
			"http://localhost:3000",
			"http://127.0.0.1:5173",
			"http://localhost:5173",
		},
		AllowCredentials: true,
	}))

	s.echo.GET("/up", s.healthHttpHandler.Up)
	s.echo.StaticFS("/", public.PublicFs)

	s.echo.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	return nil
}
