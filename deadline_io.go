package tarantool

import (
	"net"
	"time"
	"unsafe"

	"github.com/GoWebProd/gip/fasttime"
)

type DeadlineIO struct {
	to time.Duration
	c  net.Conn
}

func (d *DeadlineIO) getDeadline() time.Time {
	type internalTime struct {
		wall uint64
		ext  int64
		loc  *time.Location
	}

	t := internalTime{0, fasttime.NowNano() + int64(d.to), time.Local}

	return *(*time.Time)(unsafe.Pointer(&t))
}

func (d *DeadlineIO) Write(b []byte) (int, error) {
	if d.to > 0 {
		d.c.SetWriteDeadline(d.getDeadline())
	}

	return d.c.Write(b)
}

func (d *DeadlineIO) Read(b []byte) (int, error) {
	if d.to > 0 {
		d.c.SetReadDeadline(d.getDeadline())
	}

	return d.c.Read(b)
}
