package config

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/dibrinsofor/core-banking/internal/handlers"
	"github.com/dibrinsofor/core-banking/internal/middlewares"
	redistest "github.com/dibrinsofor/core-banking/internal/redis"
	"github.com/gin-gonic/gin"
)

type Server struct {
	h   *handlers.Handler
	e   *gin.Engine
	srv http.Server
}

func New(h *handlers.Handler) *Server {
	return &Server{
		h: h,
		e: gin.Default(),
	}
}

func (s *Server) SetupMiddlewares(m []gin.HandlerFunc) {
	s.e.Use(m...)
}

func (s *Server) SetupRoutes() *gin.Engine {
	mw := []gin.HandlerFunc{middlewares.Cors()}
	s.SetupMiddlewares(mw)

	s.e.GET("/healthcheck", s.h.Healthcheck)
	s.e.POST("/createAccount", s.h.CreateUser)
	s.e.POST("/deposit", s.h.Deposit).Use(redistest.VerifyIdempotencyKey())
	s.e.POST("/withdraw", s.h.Withdraw).Use(redistest.VerifyIdempotencyKey())
	s.e.POST("/transfer", s.h.Transfer).Use(redistest.VerifyIdempotencyKey())
	s.e.GET("/transHistory", s.h.TransactionHistory)

	// authenticatedRoutes := s.e.Group("/auth").Use(middlewares.AuthorizeJWT())
	// {
	// 	authenticatedRoutes.POST("/deposit", s.h.deposit)
	// 	authenticatedRoutes.GET("/withdraw", s.h.withdraw)
	// }
	return s.e
}

func (s *Server) Start() {

	s.srv = http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("PORT")),
		Handler: s.SetupRoutes(),
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		if err := s.srv.Close(); err != nil {
			log.Println("failed to shutdown server", err)
		}
	}()

	if err := s.srv.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Println("server closed after interruption")
		} else {
			log.Println("unexpected server shutdown. err:", err)
		}
	}
}

func (s *Server) Stop() error {
	return s.srv.Close()
}
