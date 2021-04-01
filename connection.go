package tarantool

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/GoWebProd/gip/fasttime"
	"github.com/GoWebProd/gip/spinlock"
)

type Connection struct {
	addr  string
	c     net.Conn
	mutex sync.Mutex
	// Schema contains schema loaded on connection.
	Schema    *Schema
	requestId uint32
	// Greeting contains first message sent by tarantool
	Greeting *Greeting

	shard []connShard
	queue chan *Future

	lenBuf [5]byte

	control chan struct{}
	opts    Opts
	state   uint32
}

// Opts is a way to configure Connection
type Opts struct {
	// Timeout is requests timeout.
	// Also used to setup net.TCPConn.Set(Read|Write)Deadline
	Timeout time.Duration
	// Reconnect is a pause between reconnection attempts.
	// If specified, then when tarantool is not reachable or disconnected,
	// new connect attempt is performed after pause.
	// By default, no reconnection attempts are performed,
	// so once disconnected, connection becomes Closed.
	Reconnect time.Duration
	// MaxReconnects is a maximum reconnect attempts.
	// After MaxReconnects attempts Connection becomes closed.
	MaxReconnects uint
	// User name for authorization
	User string
	// Pass is password for authorization
	Pass string
	// Concurrency is amount of separate mutexes for request
	// queues and buffers inside of connection.
	// It is rounded upto nearest power of 2.
	// By default it is runtime.GOMAXPROCS(-1) * 4
	Concurrency uint32
	// SkipSchema disables schema loading. Without disabling schema loading,
	// there is no way to create Connection for currently not accessible tarantool.
	SkipSchema bool
}

// Connect creates and configures new Connection
//
// Address could be specified in following ways:
//
// TCP connections:
// - tcp://192.168.1.1:3013
// - tcp://my.host:3013
// - tcp:192.168.1.1:3013
// - tcp:my.host:3013
// - 192.168.1.1:3013
// - my.host:3013
// Unix socket:
// - unix:///abs/path/tnt.sock
// - unix:path/tnt.sock
// - /abs/path/tnt.sock  - first '/' indicates unix socket
// - ./rel/path/tnt.sock - first '.' indicates unix socket
// - unix/:path/tnt.sock  - 'unix/' acts as a "host" and "/path..." as a port
//
// Note:
//
// - If opts.Reconnect is zero (default), then connection either already connected
// or error is returned.
//
// - If opts.Reconnect is non-zero, then error will be returned only if authorization// fails. But if Tarantool is not reachable, then it will attempt to reconnect later
// and will not end attempts on authorization failures.
func Connect(addr string, opts Opts) (*Connection, error) {
	conn := &Connection{
		addr:      addr,
		requestId: 0,
		Greeting:  &Greeting{},
		control:   make(chan struct{}),
		opts:      opts,
	}

	maxprocs := uint32(runtime.GOMAXPROCS(-1))
	if conn.opts.Concurrency == 0 || conn.opts.Concurrency > maxprocs*128 {
		conn.opts.Concurrency = maxprocs * 4
	}

	if c := conn.opts.Concurrency; c&(c-1) != 0 {
		for i := uint(1); i < 32; i <<= 1 {
			c |= c >> i
		}

		conn.opts.Concurrency = c + 1
	}

	conn.shard = make([]connShard, conn.opts.Concurrency)
	conn.queue = make(chan *Future, conn.opts.Concurrency*2)

	for i := range conn.shard {
		shard := &conn.shard[i]

		for j := range shard.requests {
			shard.requests[j].last = &shard.requests[j].first
		}
	}

	if err := conn.createConnection(false); err != nil {
		ter, ok := err.(Error)

		switch {
		case conn.opts.Reconnect <= 0:
			return nil, err
		case ok && (ter.Code == ErrNoSuchUser || ter.Code == ErrPasswordMismatch):
			/* reported auth errors immediatly */
			return nil, err
		default:
			// without SkipSchema it is useless
			go func(conn *Connection) {
				conn.mutex.Lock()
				defer conn.mutex.Unlock()
				if err := conn.createConnection(true); err != nil {
					conn.closeConnection(err, true)
				}
			}(conn)
		}
	}

	go conn.pinger()

	if conn.opts.Timeout > 0 {
		go conn.timeouts()
	}

	if !conn.opts.SkipSchema {
		if err := conn.loadSchema(); err != nil {
			conn.mutex.Lock()
			conn.closeConnection(err, true)
			conn.mutex.Unlock()

			return nil, err
		}
	}

	return conn, nil
}

// Close closes Connection.
// After this method called, there is no way to reopen this Connection.
func (conn *Connection) Close() error {
	err := ClientError{ErrConnectionClosed, "connection closed by client"}

	conn.mutex.Lock()
	defer conn.mutex.Unlock()

	return conn.closeConnection(err, true)
}

const (
	connDisconnected = 0
	connConnected    = 1
	connClosed       = 2
)

func (conn *Connection) createConnection(reconnect bool) error {
	var reconnects uint

	for conn.c == nil && conn.state == connDisconnected {
		now := time.Now()
		err := conn.dial()

		if err == nil {
			return nil
		}

		if !reconnect {
			return err
		}

		if conn.opts.MaxReconnects > 0 && reconnects > conn.opts.MaxReconnects {
			// mark connection as closed to avoid reopening by another goroutine
			return ClientError{ErrConnectionClosed, "last reconnect failed"}
		}

		reconnects++

		conn.mutex.Unlock()
		time.Sleep(now.Add(conn.opts.Reconnect).Sub(time.Now()))
		conn.mutex.Lock()
	}

	if conn.state == connClosed {
		return ClientError{ErrConnectionClosed, "using closed connection"}
	}

	return nil
}

func (conn *Connection) closeConnection(neterr error, forever bool) error {
	var err error

	conn.lockShards()

	if forever {
		if conn.state != connClosed {
			close(conn.control)
			atomic.StoreUint32(&conn.state, connClosed)
		}
	} else {
		atomic.StoreUint32(&conn.state, connDisconnected)
	}

	if conn.c != nil {
		err = conn.c.Close()
		conn.c = nil
	}

	for i := range conn.shard {
		requests := &conn.shard[i].requests
		for pos := range requests {
			fut := requests[pos].first
			requests[pos].first = nil
			requests[pos].last = &requests[pos].first

			for fut != nil {
				fut.err = neterr

				fut.markReady(conn)

				fut, fut.next = fut.next, nil
			}
		}
	}

	conn.unlockShards()

	return err
}

func (conn *Connection) pinger() {
	to := conn.opts.Timeout
	if to == 0 {
		to = 3 * time.Second
	}

	t := time.NewTicker(to / 3)
	defer t.Stop()

	for {
		select {
		case <-conn.control:
			return
		case <-t.C:
			conn.Ping()
		}
	}
}

var epoch = fasttime.NowNano()

func (conn *Connection) timeouts() {
	timeout := conn.opts.Timeout

	t := time.NewTimer(timeout)
	defer t.Stop()

	for {
		var nowepoch int64

		select {
		case <-conn.control:
			return
		case <-t.C:
		}

		minNext := fasttime.NowNano() - epoch + int64(timeout)

		for i := range conn.shard {
			nowepoch = fasttime.NowNano() - epoch

			shard := &conn.shard[i]
			for pos := range shard.requests {
				shard.rmut.Lock()

				pair := &shard.requests[pos]
				for pair.first != nil && pair.first.timeout < nowepoch {
					fut := pair.first
					pair.first = fut.next

					if fut.next == nil {
						pair.last = &pair.first
					} else {
						fut.next = nil
					}

					fut.err = ClientError{
						Code: ErrTimeouted,
						Msg:  fmt.Sprintf("client timeout for request %d", fut.request.requestId),
					}

					fut.markReady(conn)
				}
				if pair.first != nil && pair.first.timeout < minNext {
					minNext = pair.first.timeout
				}

				shard.rmut.Unlock()
			}
		}

		nowepoch = fasttime.NowNano() - epoch
		if nowepoch+int64(time.Microsecond) < minNext {
			t.Reset(time.Duration(minNext - nowepoch))
		} else {
			t.Reset(time.Microsecond)
		}
	}
}

func (conn *Connection) dial() error {
	network := "tcp"
	address := conn.addr

	timeout := conn.opts.Reconnect / 2
	if timeout == 0 {
		timeout = 500 * time.Millisecond
	} else if timeout > 5*time.Second {
		timeout = 5 * time.Second
	}

	switch {
	case len(address) > 0 && (address[0] == '.' || address[0] == '/'):
		network = "unix"
	case len(address) >= 7 && address[:7] == "unix://":
		network = "unix"
		address = address[7:]
	case len(address) >= 5 && address[:5] == "unix:":
		network = "unix"
		address = address[5:]
	case len(address) >= 6 && address[:6] == "unix/:":
		network = "unix"
		address = address[6:]
	case len(address) >= 6 && address[:6] == "tcp://":
		address = address[6:]
	case len(address) >= 4 && address[:4] == "tcp:":
		address = address[4:]
	}

	connection, err := net.DialTimeout(network, address, timeout)
	if err != nil {
		return err
	}

	dc := &DeadlineIO{to: conn.opts.Timeout, c: connection}
	r := bufio.NewReaderSize(dc, 128*1024)
	w := bufio.NewWriterSize(dc, 128*1024)

	greeting := make([]byte, 128)

	if _, err = io.ReadFull(r, greeting); err != nil {
		connection.Close()

		return err
	}

	conn.Greeting.Version = string(greeting[:64])
	conn.Greeting.auth = string(greeting[64:108])

	// Auth
	if conn.opts.User != "" {
		scr, err := scramble(conn.Greeting.auth, conn.opts.Pass)
		if err != nil {
			connection.Close()

			return errors.New("auth: scrambling failure " + err.Error())
		}

		if err = conn.writeAuthRequest(w, scr); err != nil {
			connection.Close()

			return err
		}

		if err = conn.readAuthResponse(r); err != nil {
			connection.Close()

			return err
		}
	}

	conn.lockShards()

	conn.c = connection

	atomic.StoreUint32(&conn.state, connConnected)
	conn.unlockShards()

	go conn.writer(w, connection)
	go conn.reader(r, connection)

	return nil
}

func (conn *Connection) reconnect(neterr error, c net.Conn) {
	conn.mutex.Lock()

	if conn.opts.Reconnect > 0 {
		if c == conn.c {
			conn.closeConnection(neterr, false)

			if err := conn.createConnection(true); err != nil {
				conn.closeConnection(err, true)
			}
		}
	} else {
		conn.closeConnection(neterr, true)
	}

	conn.mutex.Unlock()
}

func (conn *Connection) lockShards() {
	for i := range conn.shard {
		conn.shard[i].rmut.Lock()
	}
}

func (conn *Connection) unlockShards() {
	for i := range conn.shard {
		conn.shard[i].rmut.Unlock()
	}
}

func (conn *Connection) nextRequestId() (requestId uint32) {
	return atomic.AddUint32(&conn.requestId, 1)
}

type Greeting struct {
	Version string
	auth    string
}

const requestsMap = 128

type connShard struct {
	rmut     spinlock.Locker
	requests [requestsMap]struct {
		first *Future
		last  **Future
	}
}
