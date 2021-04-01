package tarantool

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"io"

	"github.com/GoWebProd/msgp/msgp"
)

func scramble(encodedSalt, pass string) (scramble []byte, err error) {
	/* ==================================================================
		According to: http://tarantool.org/doc/dev_guide/box-protocol.html

		salt = base64_decode(encodedSalt);
		step1 = sha1(password);
		step2 = sha1(step1);
		step3 = sha1(salt, step2);
		scramble = xor(step1, step3);
		return scramble;

	===================================================================== */
	scrambleSize := sha1.Size // == 20

	salt, err := base64.StdEncoding.DecodeString(encodedSalt)
	if err != nil {
		return
	}
	step1 := sha1.Sum([]byte(pass))
	step2 := sha1.Sum(step1[0:])
	hash := sha1.New() // may be create it once per connection ?
	hash.Write(salt[0:scrambleSize])
	hash.Write(step2[0:])
	step3 := hash.Sum(nil)

	return xor(step1[0:], step3[0:], scrambleSize), nil
}

func xor(left, right []byte, size int) []byte {
	result := make([]byte, size)
	for i := 0; i < size; i++ {
		result[i] = left[i] ^ right[i]
	}
	return result
}

type authRequest struct {
	UserName string
	Method   string
	Scramble []byte
}

func (a authRequest) EncodeMsg(enc *msgp.Writer) error {
	if err := enc.WriteMapHeader(2); err != nil {
		return err
	}
	if err := enc.WriteUint(KeyUserName); err != nil {
		return err
	}
	if err := enc.WriteString(a.UserName); err != nil {
		return err
	}
	if err := enc.WriteUint(KeyTuple); err != nil {
		return err
	}
	if err := enc.WriteArrayHeader(2); err != nil {
		return err
	}
	if err := enc.WriteString(a.Method); err != nil {
		return err
	}
	if err := enc.WriteStringFromBytes(a.Scramble); err != nil {
		return err
	}

	return nil
}

func (conn *Connection) writeAuthRequest(w *bufio.Writer, scramble []byte) error {
	request := &Future{
		request: request{
			requestId:   0,
			requestCode: AuthRequest,

			userName: conn.opts.User,
			method:   "chap-sha1",
			scramble: scramble,
		},
	}

	mw := msgp.NewWriter(w)

	if err := request.write(w, mw); err != nil {
		return errors.New("auth: flush error " + err.Error())
	}

	if err := w.Flush(); err != nil {
		return errors.New("auth: flush error " + err.Error())
	}

	return nil
}

func (conn *Connection) readAuthResponse(r io.Reader) error {
	respBytes, err := conn.read(r)
	if err != nil {
		return errors.New("auth: read error " + err.Error())
	}

	resp := Response{buf: respBytes}

	err = resp.decode()
	if err != nil {
		resp.Release()

		return errors.New("auth: decode response error " + err.Error())
	}

	if resp.Error != "" {
		resp.Release()

		return Error{resp.Code, resp.Error}
	}

	return nil
}
