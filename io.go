package tarantool

import (
	"bufio"
	"io"
	"net"
	"runtime"
	"sync/atomic"

	"github.com/GoWebProd/gip/allocator"
	"github.com/GoWebProd/msgp/msgp"
	"github.com/pkg/errors"
)

func (conn *Connection) writer(w *bufio.Writer, c net.Conn) {
	var future *Future

	writer := msgp.NewWriter(w)

	for atomic.LoadUint32(&conn.state) != connClosed {
		select {
		case future = <-conn.queue:
		default:
			runtime.Gosched()
			if len(conn.queue) == 0 {
				if err := w.Flush(); err != nil {
					conn.reconnect(err, c)

					return
				}
			}
			select {
			case future = <-conn.queue:
			case <-conn.control:
				return
			}
		}

		if err := future.write(w, writer); err != nil {
			conn.reconnect(err, c)

			return
		}
	}
}

func (conn *Connection) reader(r *bufio.Reader, c net.Conn) {
	for atomic.LoadUint32(&conn.state) != connClosed {
		respBytes, err := conn.read(r)
		if err != nil {
			conn.reconnect(err, c)

			return
		}

		resp := Response{buf: respBytes}

		err = resp.decode()
		if err != nil {
			conn.reconnect(err, c)
			resp.Release()

			return
		}

		if fut := conn.fetchFuture(resp.RequestId); fut != nil {
			fut.resp = resp

			fut.markReady(conn)
		} else {
			resp.Release()
		}
	}
}

func (conn *Connection) read(r io.Reader) ([]byte, error) {
	var length int

	if _, err := io.ReadFull(r, conn.lenBuf[:]); err != nil {
		return nil, errors.Wrap(err, "read response header error")
	}

	if conn.lenBuf[0] != 0xce {
		return nil, errors.New("wrong reponse header")
	}

	length = (int(conn.lenBuf[1]) << 24) + (int(conn.lenBuf[2]) << 16) + (int(conn.lenBuf[3]) << 8) + int(conn.lenBuf[4])
	if length == 0 {
		return nil, errors.New("response should not be 0 length")
	}

	response := allocator.Alloc(length)

	if _, err := io.ReadFull(r, response); err != nil {
		return nil, errors.Wrap(err, "read response body error")
	}

	return response, nil
}
