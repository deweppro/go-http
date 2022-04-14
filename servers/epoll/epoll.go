package epoll

import (
	"net"
	"sync"
	"syscall"

	"github.com/deweppro/go-http/pkg/errs"

	"github.com/deweppro/go-errors"
	"github.com/deweppro/go-http/internal"
	"github.com/deweppro/go-logger"
	"golang.org/x/sys/unix"
)

type (
	NetMap   map[int]*netItem
	NetSlice []*netItem
	Events   []unix.EpollEvent

	//Epoll ...
	Epoll struct {
		fd     int
		conn   NetMap
		events Events
		nets   NetSlice
		log    logger.Logger
		mux    sync.RWMutex
	}
)

const (
	epollEvents     = unix.POLLIN | unix.POLLRDHUP | unix.POLLERR | unix.POLLHUP | unix.POLLNVAL
	eventCount      = 100
	eventIntervalMS = 500
)

func newEpoll(log logger.Logger) (*Epoll, error) {
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
func (v *Epoll) AddOrClose(c net.Conn) error {
	fd := internal.FD(c)
	err := unix.EpollCtl(v.fd, syscall.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Events: epollEvents, Fd: int32(fd)})
	if err != nil {
		return errors.Wrap(err, c.Close())
	}
	v.mux.Lock()
	v.conn[fd] = &netItem{Conn: c, Fd: fd}
	v.mux.Unlock()
	return nil
}

func (v *Epoll) removeFD(fd int) error {
	return unix.EpollCtl(v.fd, syscall.EPOLL_CTL_DEL, fd, nil)
}

//Close ...
func (v *Epoll) Close(c *netItem) error {
	v.mux.Lock()
	defer v.mux.Unlock()
	return v.closeConn(c)
}

func (v *Epoll) closeConn(c *netItem) error {
	if err := v.removeFD(c.Fd); err != nil {
		return err
	}
	delete(v.conn, c.Fd)
	return c.Conn.Close()
}

//CloseAll ...
func (v *Epoll) CloseAll() (err error) {
	v.mux.Lock()
	defer v.mux.Unlock()

	for _, conn := range v.conn {
		if err0 := v.closeConn(conn); err0 != nil {
			err = errors.Wrap(err, err0)
		}
	}
	v.conn = make(NetMap)
	return
}

func (v *Epoll) getConn(fd int) (*netItem, bool) {
	v.mux.RLock()
	conn, ok := v.conn[fd]
	v.mux.RUnlock()
	return conn, ok
}

//Wait ...
func (v *Epoll) Wait() (NetSlice, error) {
	n, err := unix.EpollWait(v.fd, v.events, eventIntervalMS)
	if err != nil {
		return nil, err
	}
	if n <= 0 {
		return nil, errs.ErrEpollEmptyEvents
	}

	v.nets = v.nets[:0]
	for i := 0; i < n; i++ {
		fd := int(v.events[i].Fd)
		conn, ok := v.getConn(fd)
		if !ok {
			if err = v.removeFD(fd); err != nil {
				v.log.WithFields(logger.Fields{
					"err": err.Error(), "fd": fd,
				}).Errorf("close fd")
			}
			continue
		}
		if conn.IsAwait() {
			continue
		}
		conn.Await(true)

		switch v.events[i].Events {
		case unix.POLLIN:
			v.nets = append(v.nets, conn)
		default:
			if err = v.Close(conn); err != nil {
				v.log.WithFields(logger.Fields{"err": err.Error()}).Errorf("epoll close connect")
			}
		}
	}

	return v.nets, nil
}
