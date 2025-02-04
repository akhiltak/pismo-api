package server

import (
	"github.com/akhiltak/pismo-api/internal/handler"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func (s *Server) initRoutes(h handler.Handler) {

	s.router.GET("/health", h.Health)
	s.router.GET("/swagger/*", echoSwagger.WrapHandler)

	account := s.router.Group("/accounts")
	{
		account.POST("", h.CreateAccount)
		account.GET("/:id", h.GetAccountByID)
	}
	transaction := s.router.Group("/transactions")
	{
		transaction.POST("", h.CreateTransaction)
		transaction.GET("", h.GetTransactions)
	}
}
