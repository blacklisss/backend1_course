package server

import (
	"context"
	"gb/backend1_course/internal/usecases/app/repos/linkrepo"
	"net/http"
	"time"
)

type Server struct {
	srv http.Server
	ln  *linkrepo.Links
}

func NewServer(addr string, h http.Handler) *Server {
	s := &Server{}

	s.srv = http.Server{
		Addr:              addr,
		Handler:           h,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
	}
	return s
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	err := s.srv.Shutdown(ctx)
	cancel()
	if err != nil {
		panic(err)
	}
}

func (s *Server) Start(ln *linkrepo.Links) {
	s.ln = ln
	// TODO: migrations
	go func() {
		_ = s.srv.ListenAndServe()
	}()
}
