package web

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/yktseng/portto-assignment/internal/database"
)

type Server struct {
	r   *gin.Engine
	srv *http.Server
	db  *database.Database
	wg  *sync.WaitGroup
}

func NewServer(ctx context.Context, db *database.Database, wg *sync.WaitGroup) *Server {
	r := gin.Default()
	s := Server{
		r,
		nil,
		db,
		wg,
	}
	s.initRoutes(ctx)
	return &s
}

func (s *Server) initRoutes(ctx context.Context) {
	s.r.GET("/blocks", s.getBlocks)
	s.r.GET("/blocks/:id", s.getBlockByID)
	s.r.GET("/transaction/:tx_hash", s.getTXByHash)
}

func (s *Server) Start(port int) error {
	s.srv = &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: s.r,
	}
	if err := s.srv.ListenAndServe(); err != nil {
		s.wg.Done()
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	defer func() {
		log.Println("web server is stopped")
		s.wg.Done()
	}()
	if err := s.srv.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
