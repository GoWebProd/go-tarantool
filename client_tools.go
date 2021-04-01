package tarantool

import (
	"github.com/GoWebProd/msgp/msgp"
)

// IntKey is utility type for passing integer key to Select*, Update* and Delete*
// It serializes to array with single integer element.
type IntKey struct {
	I int
}

func (k IntKey) EncodeMsg(enc *msgp.Writer) error {
	enc.WriteArrayHeader(1)
	enc.WriteInt(k.I)
	return nil
}

func (k IntKey) Msgsize() int {
	return 1 + msgp.IntSize(uint64(k.I))
}

// UintKey is utility type for passing unsigned integer key to Select*, Update* and Delete*
// It serializes to array with single integer element.
type UintKey struct {
	I uint
}

func (k UintKey) EncodeMsg(enc *msgp.Writer) error {
	enc.WriteArrayHeader(1)
	enc.WriteUint(k.I)
	return nil
}

func (k UintKey) Msgsize() int {
	return 1 + msgp.IntSize(uint64(k.I))
}

// UintKey is utility type for passing string key to Select*, Update* and Delete*
// It serializes to array with single string element.
type StringKey struct {
	S string
}

func (k StringKey) EncodeMsg(enc *msgp.Writer) error {
	enc.WriteArrayHeader(1)
	enc.WriteString(k.S)
	return nil
}

func (k StringKey) Msgsize() int {
	return 1 + msgp.StringSize(len(k.S))
}

// Op - is update operation
type Op struct {
	Op    string
	Field int
	Arg   msgp.Encodable
}

func (o Op) EncodeMsg(enc *msgp.Writer) error {
	enc.WriteArrayHeader(3)
	enc.WriteString(o.Op)
	enc.WriteInt(o.Field)
	return o.Arg.EncodeMsg(enc)
}

type Ops []Op

func (o Ops) EncodeMsg(enc *msgp.Writer) error {
	enc.WriteArrayHeader(uint32(len(o)))

	for i := 0; i < len(o); i++ {
		if err := o[i].EncodeMsg(enc); err != nil {
			return err
		}
	}

	return nil
}

type OpSplice struct {
	Op      string
	Field   int
	Pos     int
	Len     int
	Replace string
}

func (o OpSplice) EncodeMsg(enc *msgp.Writer) error {
	enc.WriteArrayHeader(5)
	enc.WriteString(o.Op)
	enc.WriteInt(o.Field)
	enc.WriteInt(o.Pos)
	enc.WriteInt(o.Len)
	enc.WriteString(o.Replace)
	return nil
}
