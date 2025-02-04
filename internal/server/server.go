package server

import (
	"context"
	"log/slog"

	"github.com/akhiltak/pismo-api/config"
	"github.com/akhiltak/pismo-api/db/connection/bunorm"
	"github.com/akhiltak/pismo-api/db/connection/dbmate"
	"github.com/akhiltak/pismo-api/internal/handler"
	"github.com/akhiltak/pismo-api/internal/service"
	"github.com/akhiltak/pismo-api/internal/storage/repo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	router *echo.Echo
}

func New(ctx context.Context, cfg *config.Config) *Server {

	// auto migrate database
	dbmate.Migrate(ctx, cfg.PostgresDNS, cfg.Debug)

	// connect to database
	db := bunorm.Connect(ctx, cfg.PostgresDNS, true)

	// initialize repositories
	accountRepo := repo.NewAccountRepo(db)
	transactionRepo := repo.NewTransactionRepo(db)
	operationRepo := repo.NewOperationRepo(db)

	// initialize services
	transactionService := service.NewTransactionService(accountRepo, transactionRepo, operationRepo)

	// initialize handlers
	handler := handler.New(transactionService)

	router := echo.New()

	router.Use(middleware.Logger()) // Using default logger but can be configured to work with slog or any other (https://echo.labstack.com/docs/middleware/logger)

	// Recover Middleware recovers from panics anywhere in the chain
	router.Use(middleware.Recover())

	// Secure Middleware for protection against cross-site scripting (XSS) attack,
	// content type sniffing, Clickjacking, insecure connection and other code injection attacks.
	router.Use(middleware.Secure())

	// CORS Middleware
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			echo.HeaderAccessControlAllowHeaders,
			echo.HeaderAccessControlAllowMethods,
			echo.HeaderAccessControlAllowOrigin,
			echo.HeaderAccessControlAllowCredentials,
		},
	}))
	router.HTTPErrorHandler = customHTTPErrorHandler

	srv := &Server{router: router}
	srv.initRoutes(handler)

	return srv
}

func (s *Server) Run(addr string) error {

	// Start server
	go func() {
		if err := s.router.Start(addr); err != nil {
			slog.Error("unable to start server", "error", err)
		}
	}()

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.router.Shutdown(ctx)
}
