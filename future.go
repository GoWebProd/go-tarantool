package tarantool

import (
	"bufio"

	"github.com/GoWebProd/gip/allocator"
	"github.com/GoWebProd/gip/cond"
	"github.com/GoWebProd/gip/fasttime"
	"github.com/GoWebProd/msgp/msgp"
)

var ()

type Body interface {
	msgp.Encodable
	msgp.Sizer
}

// Future is a handle for asynchronous request
type Future struct {
	request request
	timeout int64
	resp    Response
	err     error
	ready   cond.Single
	header  [14]byte

	next *Future
}

func (fut *Future) Release() {
	allocator.FreeObject(fut)
}

func (fut *Future) requestLength() int {
	return 9 + fut.request.Msgsize()
}

func (fut *Future) write(w *bufio.Writer, pack *msgp.Writer) error {
	rid := fut.request.requestId
	length := fut.requestLength()
	fut.header = [14]byte{
		0xce, byte(length >> 24), byte(length >> 16), byte(length >> 8), byte(length), // length
		0x82,                                   // 2 element map
		KeyCode, byte(fut.request.requestCode), // request code
		KeySync, 0xce,
		byte(rid >> 24), byte(rid >> 16),
		byte(rid >> 8), byte(rid),
	}

	_, err := w.Write(fut.header[:])
	if err != nil {
		return err
	}

	if err := fut.request.EncodeMsg(pack); err != nil {
		return err
	}

	if err := pack.Flush(); err != nil {
		return err
	}

	return nil
}

func (fut *Future) Get() (Response, error) {
	fut.ready.Wait()

	if fut.err != nil {
		return Response{}, fut.err
	}

	resp := fut.resp

	fut.Release()

	return resp, nil
}

func (conn *Connection) newFuture(request request) *Future {
	fut := allocator.AllocObject[Future]()

	fut.request = request
	fut.request.requestId = conn.nextRequestId()

	shardn := fut.request.requestId & (conn.opts.Concurrency - 1)
	shard := &conn.shard[shardn]

	shard.rmut.Lock()
	switch conn.state {
	case connClosed:
		fut.err = ClientError{ErrConnectionClosed, "using closed connection"}

		shard.rmut.Unlock()

		return fut
	case connDisconnected:
		fut.err = ClientError{ErrConnectionNotReady, "client connection is not ready"}

		shard.rmut.Unlock()

		return fut
	}

	pos := (fut.request.requestId / conn.opts.Concurrency) & (requestsMap - 1)
	pair := &shard.requests[pos]
	*pair.last = fut
	pair.last = &fut.next

	if conn.opts.Timeout > 0 {
		fut.timeout = fasttime.NowNano() - epoch + int64(conn.opts.Timeout)
	}

	shard.rmut.Unlock()

	return fut
}

func (fut *Future) markReady(conn *Connection) {
	fut.ready.Done()
}

func (fut *Future) fail(conn *Connection, err error) *Future {
	if f := conn.fetchFuture(fut.request.requestId); f == fut {
		f.err = err

		fut.markReady(conn)
	}

	return fut
}

func (conn *Connection) fetchFuture(reqid uint32) *Future {
	shard := &conn.shard[reqid&(conn.opts.Concurrency-1)]

	shard.rmut.Lock()

	fut := conn.fetchFutureImp(reqid)

	shard.rmut.Unlock()

	return fut
}

func (conn *Connection) fetchFutureImp(reqid uint32) *Future {
	shard := &conn.shard[reqid&(conn.opts.Concurrency-1)]
	pos := (reqid / conn.opts.Concurrency) & (requestsMap - 1)
	pair := &shard.requests[pos]
	root := &pair.first

	for {
		fut := *root
		if fut == nil {
			return nil
		}

		if fut.request.requestId == reqid {
			*root = fut.next

			if fut.next == nil {
				pair.last = root
			} else {
				fut.next = nil
			}

			return fut
		}

		root = &fut.next
	}
}
