package epoll

import (
	"io"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/deweppro/go-http/utils"
	"github.com/deweppro/go-http/web/server"
	"github.com/deweppro/go-logger"
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
)

var (
	defaultEOF = []byte("\r\n")

	ErrInvalidEOF = errors.New("invalid eof chars")
)

//Server ...
type Server struct {
	status int64
	close  chan struct{}
	wg     sync.WaitGroup

	handler ConnectionHandler
	log     logger.Logger
	config  ConfigItem
	eof     []byte

	listener net.Listener
	epoll    *Epoll
}

//NewServer ...
func New(conf *Config, handler ConnectionHandler, log logger.Logger) *Server {
	return NewCustomServer(conf.Epoll, handler, defaultEOF, log)
}

//NewCustomServer ...
func NewCustomServer(conf ConfigItem, handler ConnectionHandler, eof []byte, log logger.Logger) *Server {
	return &Server{
		status:  server.StatusOff,
		config:  conf,
		handler: handler,
		log:     log,
		close:   make(chan struct{}),
		eof:     eof,
	}
}

func (s *Server) validate() {
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

//Up ...
func (s *Server) Up() (err error) {
	if len(s.eof) == 0 {
		return ErrInvalidEOF
	}
	if !atomic.CompareAndSwapInt64(&s.status, server.StatusOff, server.StatusOn) {
		return errors.Wrapf(server.ErrServAlreadyRunning, "on %s", s.config.Addr)
	}
	s.validate()
	s.log.Infof("tcp server started on %s", s.config.Addr)

	if s.listener, err = net.Listen("tcp", s.config.Addr); err != nil {
		return
	}
	if s.epoll, err = NewEpoll(s.log); err != nil {
		return
	}
	s.wg.Add(2)
	go s.connAccept()
	go s.epollAccept()
	return
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

				if ne, ok := err.(net.Error); ok && ne.Temporary() {
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
			case ErrEpollEmptyEvents:
				continue
			case unix.EINTR:
				continue
			default:
				s.log.Errorf("epoll accept conn: %s", err.Error())
				continue
			}

			for _, c := range list {
				go func(conn *NetItem) {
					defer conn.Await(false)

					if err1 := connection(conn.Conn, s.handler, s.eof); err1 != nil {
						_ = s.epoll.Close(conn)
						if err1 != io.EOF {
							s.log.Errorf("epoll bad conn %s: %s", conn.Conn.RemoteAddr().String(), err1.Error())
						}
					}
				}(c)
			}
		}
	}
}

//Down ...
func (s *Server) Down() (err error) {
	close(s.close)
	if e := s.epoll.CloseAll(); e != nil {
		err = errors.Wrap(err, e.Error())
	}
	if e := s.listener.Close(); e != nil {
		err = errors.Wrap(err, e.Error())
	}
	s.log.Infof("tcp server stopped on %s", s.config.Addr)
	s.wg.Wait()
	return
}
