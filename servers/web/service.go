package web

import (
	"context"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/deweppro/go-http/servers"

	"github.com/deweppro/go-errors"
	"github.com/deweppro/go-http/internal"
	"github.com/deweppro/go-http/pkg/errs"
	"github.com/deweppro/go-logger"
)

const (
	defaultTimeout         = 10 * time.Second
	defaultShutdownTimeout = 1 * time.Second
)

type Server struct {
	status  int64
	conf    servers.Config
	serv    *http.Server
	handler http.Handler
	log     logger.Logger
	wg      sync.WaitGroup
}

//New create default http server
func New(conf servers.Config, handler http.Handler, log logger.Logger) *Server {
	srv := &Server{
		conf:    conf,
		handler: handler,
		log:     log,
		status:  servers.StatusOff,
	}
	srv.validate()
	return srv
}

func (s *Server) validate() {
	if s.conf.ReadTimeout == 0 {
		s.conf.ReadTimeout = defaultTimeout
	}
	if s.conf.WriteTimeout == 0 {
		s.conf.WriteTimeout = defaultTimeout
	}
	if s.conf.IdleTimeout == 0 {
		s.conf.IdleTimeout = defaultTimeout
	}
	if s.conf.ShutdownTimeout == 0 {
		s.conf.ShutdownTimeout = defaultShutdownTimeout
	}
	s.conf.Addr = internal.ValidateAddress(s.conf.Addr)
}

//Up start http server
func (s *Server) Up() error {
	if !atomic.CompareAndSwapInt64(&s.status, servers.StatusOff, servers.StatusOn) {
		return errors.WrapMessage(errs.ErrServAlreadyRunning, "starting server on %s", s.conf.Addr)
	}
	s.serv = &http.Server{
		Addr:         s.conf.Addr,
		ReadTimeout:  s.conf.ReadTimeout,
		WriteTimeout: s.conf.WriteTimeout,
		IdleTimeout:  s.conf.IdleTimeout,
		Handler:      s.handler,
	}

	s.wg.Add(1)
	s.log.Infof("http server started on %s", s.conf.Addr)

	go func() {
		err := s.serv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			s.log.Errorf("http server stopped on %s with error: %s", s.conf.Addr, err.Error())
		} else {
			s.log.Infof("http server stopped on %s", s.conf.Addr)
		}
		s.wg.Done()
	}()
	return nil
}

//Down stop http server
func (s *Server) Down() error {
	if !atomic.CompareAndSwapInt64(&s.status, servers.StatusOn, servers.StatusOff) {
		return errors.WrapMessage(errs.ErrServAlreadyStopped, "stopping server on %s", s.conf.Addr)
	}
	ctx, cncl := context.WithTimeout(context.TODO(), s.conf.ShutdownTimeout)
	defer cncl()
	err := s.serv.Shutdown(ctx)
	s.wg.Wait()
	return err
}
