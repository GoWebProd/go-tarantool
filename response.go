package tarantool

import (
	"github.com/GoWebProd/gip/allocator"
	"github.com/GoWebProd/msgp/msgp"
)

type Response struct {
	RequestId uint32
	Code      uint32
	Error     string // error message
	Data      []byte

	buf []byte
}

func (resp *Response) Release() {
	if resp.buf != nil {
		allocator.Free(resp.buf)
	}
}

func (resp *Response) decode() error {
	var (
		l      uint32
		err    error
		remain []byte
		cd     int
	)

	remain = resp.buf

	if l, remain, err = msgp.ReadMapHeaderBytes(remain); err != nil {
		return err
	}

	for ; l > 0; l-- {
		if cd, remain, err = msgp.ReadIntBytes(remain); err != nil {
			return err
		}

		switch cd {
		case KeySync:
			var rid uint64

			if rid, remain, err = msgp.ReadUint64Bytes(remain); err != nil {
				return err
			}

			resp.RequestId = uint32(rid)
		case KeyCode:
			var rcode uint64

			if rcode, remain, err = msgp.ReadUint64Bytes(remain); err != nil {
				return err
			}

			resp.Code = uint32(rcode)
		default:
			if remain, err = msgp.Skip(remain); err != nil {
				return err
			}
		}
	}

	if l, remain, err = msgp.ReadMapHeaderBytes(remain); err != nil {
		return err
	}

	for ; l > 0; l-- {
		if cd, remain, err = msgp.ReadIntBytes(remain); err != nil {
			return err
		}

		switch cd {
		case KeyData:
			if resp.Data, remain, err = getRawBody(remain); err != nil {
				return err
			}
		case KeyError:
			if resp.Error, remain, err = msgp.ReadStringBytes(remain); err != nil {
				return err
			}
		default:
			if remain, err = msgp.Skip(remain); err != nil {
				return err
			}
		}
	}

	if resp.Code != OkCode {
		resp.Code &^= ErrorCodeBit
	}

	return nil
}

func getRawBody(data []byte) ([]byte, []byte, error) {
	var (
		remain []byte
		err    error
	)

	if remain, err = msgp.Skip(data); err != nil {
		return nil, nil, err
	}

	return data[:len(data)-len(remain)], remain, nil
}
