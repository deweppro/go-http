package epoll

import (
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/deweppro/go-errors"

	"github.com/deweppro/go-http/pkg/errs"

	"github.com/deweppro/go-http/internal"
	"github.com/deweppro/go-http/servers"
	"github.com/deweppro/go-logger"
	"golang.org/x/sys/unix"
)

var defaultEOF = []byte("\r\n")

//Server ...
type Server struct {
	status   int64
	close    chan struct{}
	wg       sync.WaitGroup
	handler  ConnectionHandler
	log      logger.Logger
	conf     servers.Config
	eof      []byte
	listener net.Listener
	epoll    *Epoll
}

func New(conf servers.Config, handler ConnectionHandler, eof []byte, log logger.Logger) *Server {
	return &Server{
		status:  servers.StatusOff,
		conf:    conf,
		handler: handler,
		log:     log,
		close:   make(chan struct{}),
		eof:     eof,
	}
}

func (s *Server) validate() {
	s.conf.Addr = internal.ValidateAddress(s.conf.Addr)
	if len(s.eof) == 0 {
		s.eof = defaultEOF
	}
}

//Up ...
func (s *Server) Up() (err error) {
	if !atomic.CompareAndSwapInt64(&s.status, servers.StatusOff, servers.StatusOn) {
		return errors.WrapMessage(errs.ErrServAlreadyRunning, "starting server on %s", s.conf.Addr)
	}
	s.validate()
	if s.listener, err = net.Listen("tcp", s.conf.Addr); err != nil {
		return
	}
	if s.epoll, err = newEpoll(s.log); err != nil {
		return
	}
	s.log.Infof("tcp server started on %s", s.conf.Addr)
	s.wg.Add(2)
	go s.connAccept()
	go s.epollAccept()
	return
}

//Down ...
func (s *Server) Down() error {
	close(s.close)
	err := errors.Wrap(s.epoll.CloseAll(), s.listener.Close())
	s.wg.Wait()
	s.log.Infof("tcp server stopped on %s", s.conf.Addr)
	return err
}

func (s *Server) connAccept() {
	defer func() {
		s.wg.Done()
	}()
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.close:
				return
			default:
				s.log.Errorf("epoll conn accept: %s", err.Error())
				if err0, ok := err.(net.Error); ok && err0.Temporary() {
					time.Sleep(1 * time.Second)
					continue
				}
				return
			}
		}
		if err = s.epoll.AddOrClose(conn); err != nil {
			s.log.Errorf("epoll add conn: %s", err.Error())
		}
	}
}

func (s *Server) epollAccept() {
	defer func() {
		s.wg.Done()
	}()
	for {
		select {
		case <-s.close:
			return
		default:
			list, err := s.epoll.Wait()
			switch err {
			case nil:
			case errs.ErrEpollEmptyEvents:
				continue
			case unix.EINTR:
				continue
			default:
				s.log.Errorf("epoll accept conn: %s", err.Error())
				continue
			}

			for _, c := range list {
				go func(conn *netItem) {
					defer conn.Await(false)

					if err1 := newConnection(conn.Conn, s.handler, s.eof); err1 != nil {
						if err2 := s.epoll.Close(conn); err2 != nil {
							s.log.Errorf("epoll close conn %s: %s", conn.Conn.RemoteAddr().String(), err2.Error())
						}
						if err1 != io.EOF {
							s.log.Errorf("epoll bad conn %s: %s", conn.Conn.RemoteAddr().String(), err1.Error())
						}
					}
				}(c)
			}
		}
	}
}
