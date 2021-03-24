package server

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/deweppro/go-http/utils"
	"github.com/deweppro/go-http/web/routes"
	"github.com/deweppro/go-logger"
	"github.com/pkg/errors"
)

const (
	defaultTimeOut   = 10 * time.Second
	defaultGSTimeOut = 1 * time.Second
	StatusOn         = 1
	StatusOff        = 0
)

type (
	//Server ...
	Server struct {
		status  int64
		config  ConfigItem
		server  *http.Server
		handler http.Handler
		logger  logger.Logger
		wg      sync.WaitGroup
	}
)

//New create default http server
func New(config *Config, router *routes.Router, log logger.Logger) *Server {
	return NewCustom(config.HTTP, router, log)
}

//NewCustom create custom http server
func NewCustom(conf ConfigItem, handler http.Handler, log logger.Logger) *Server {
	srv := &Server{
		config:  conf,
		handler: handler,
		logger:  log,
		status:  StatusOff,
	}
	srv.validate()
	return srv
}

func (s *Server) validate() {
	if s.config.ReadTimeout == 0 {
		s.config.ReadTimeout = defaultTimeOut
	}
	if s.config.WriteTimeout == 0 {
		s.config.WriteTimeout = defaultTimeOut
	}
	if s.config.IdleTimeout == 0 {
		s.config.IdleTimeout = defaultTimeOut
	}
	if s.config.ShutdownTimeout == 0 {
		s.config.ShutdownTimeout = defaultGSTimeOut
	}
	hp := strings.Split(s.config.Addr, ":")
	if len(hp[0]) == 0 {
		hp[0] = "0.0.0.0"
		s.config.Addr = strings.Join(hp, ":")
	}
	if len(hp) < 2 || (len(hp) == 2 && len(hp[1]) == 0) {
		addr, err := utils.RandomPort(hp[0])
		if err != nil {
			addr = hp[0] + ":10000"
		}
		s.config.Addr = addr
	}
}

//Up start http server
func (s *Server) Up() error {
	if !atomic.CompareAndSwapInt64(&s.status, StatusOff, StatusOn) {
		return errors.Wrapf(ErrServAlreadyRunning, "on %s", s.config.Addr)
	}
	s.server = &http.Server{
		Addr:         s.config.Addr,
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
		IdleTimeout:  s.config.IdleTimeout,
		Handler:      s.handler,
	}

	s.wg.Add(1)
	s.logger.Infof("http server started on %s", s.config.Addr)

	go func() {
		err := s.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			s.logger.Errorf("http server stopped on %s with error: %s", s.config.Addr, err.Error())
		} else {
			s.logger.Infof("http server stopped on %s", s.config.Addr)
		}
		s.wg.Done()
	}()
	return nil
}

//Down stop http server
func (s *Server) Down() error {
	if !atomic.CompareAndSwapInt64(&s.status, StatusOn, StatusOff) {
		return ErrServAlreadyStopped
	}
	ctx, cncl := context.WithTimeout(context.TODO(), s.config.ShutdownTimeout)
	defer cncl()
	err := s.server.Shutdown(ctx)
	s.wg.Wait()
	return err
}
