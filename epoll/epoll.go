package epoll

import (
	"net"
	"sync"
	"syscall"

	"github.com/deweppro/go-http/utils"
	"github.com/deweppro/go-logger"
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
)

type (
	NetMap   map[int]*NetItem
	NetSlice []*NetItem
	Events   []unix.EpollEvent

	//Epoll ...
	Epoll struct {
		fd     int
		conn   NetMap
		events Events
		nets   NetSlice
		log    logger.Logger
		sync.RWMutex
	}
)

const (
	epollEvents     = unix.POLLIN | unix.POLLRDHUP | unix.POLLERR | unix.POLLHUP | unix.POLLNVAL
	eventCount      = 100
	eventIntervalMS = 500
)

//New init epoll
func NewEpoll(log logger.Logger) (*Epoll, error) {
	fd, err := unix.EpollCreate1(0)
	if err != nil {
		return nil, err
	}
	return &Epoll{
		fd:     fd,
		conn:   make(NetMap),
		events: make(Events, eventCount),
		nets:   make(NetSlice, eventCount),
		log:    log,
	}, nil
}

//AddOrClose ...
func (e *Epoll) AddOrClose(c net.Conn) error {
	fd := utils.FD(c)
	err := unix.EpollCtl(e.fd, syscall.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Events: epollEvents, Fd: int32(fd)})
	if err != nil {
		if er := c.Close(); er != nil {
			err = errors.Wrap(err, er.Error())
		}
		return err
	}
	e.Lock()
	e.conn[fd] = &NetItem{Conn: c, Fd: fd}
	e.Unlock()
	return nil
}

func (e *Epoll) removeFD(fd int) error {
	return unix.EpollCtl(e.fd, syscall.EPOLL_CTL_DEL, fd, nil)
}

//Close ...
func (e *Epoll) Close(c *NetItem) error {
	e.Lock()
	defer e.Unlock()
	return e.closeConn(c)
}

func (e *Epoll) closeConn(c *NetItem) error {
	if err := e.removeFD(c.Fd); err != nil {
		return err
	}
	delete(e.conn, c.Fd)
	return c.Conn.Close()
}

//CloseAll ...
func (e *Epoll) CloseAll() (er error) {
	e.Lock()
	defer e.Unlock()

	for _, conn := range e.conn {
		if err := e.closeConn(conn); err != nil {
			er = errors.Wrap(er, err.Error())
		}
	}
	e.conn = make(NetMap)
	return
}

func (e *Epoll) getConn(fd int) (*NetItem, bool) {
	e.RLock()
	conn, ok := e.conn[fd]
	e.RUnlock()
	return conn, ok
}

//Wait ...
func (e *Epoll) Wait() (NetSlice, error) {
	n, err := unix.EpollWait(e.fd, e.events, eventIntervalMS)
	if err != nil {
		return nil, err
	}
	if n <= 0 {
		return nil, ErrEpollEmptyEvents
	}

	e.nets = e.nets[:0]
	for i := 0; i < n; i++ {
		fd := int(e.events[i].Fd)
		conn, ok := e.getConn(fd)
		if !ok {
			_ = e.removeFD(fd)
			continue
		}
		if conn.IsAwait() {
			continue
		}
		conn.Await(true)

		switch e.events[i].Events {
		case unix.POLLIN:
			e.nets = append(e.nets, conn)
		default:
			if err = e.Close(conn); err != nil {
				e.log.Errorf("epoll close connect: %s", err.Error())
			}
		}
	}

	return e.nets, nil
}
