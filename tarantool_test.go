package tarantool

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/GoWebProd/msgp/msgp"
)

type Tuple struct {
	Id   uint
	Msg  string
	Name string
}

func (c *Tuple) EncodeMsg(e *msgp.Writer) error {
	e.WriteArrayHeader(3)
	e.WriteUint(c.Id)
	e.WriteString(c.Msg)
	e.WriteString(c.Name)
	return nil
}

func (c *Tuple) Msgsize() int {
	return 1 + msgp.IntSize(uint64(c.Id)) + msgp.StringSize(len(c.Msg)) + msgp.StringSize(len(c.Name))
}

func (c *Tuple) UnmarshalMsg(d []byte) ([]byte, error) {
	var err error
	var l uint32
	if l, d, err = msgp.ReadArrayHeaderBytes(d); err != nil {
		return nil, err
	}
	if l != 3 {
		return nil, fmt.Errorf("array len doesn't match: %d", l)
	}
	if c.Id, d, err = msgp.ReadUintBytes(d); err != nil {
		return nil, err
	}
	if c.Msg, d, err = msgp.ReadStringBytes(d); err != nil {
		return nil, err
	}
	if c.Name, d, err = msgp.ReadStringBytes(d); err != nil {
		return nil, err
	}
	return d, nil
}

type Tuples []Tuple

func (m Tuples) EncodeMsg(e *msgp.Writer) error {
	e.WriteArrayHeader(uint32(len(m)))
	for i := 0; i < len(m); i++ {
		if err := m[i].EncodeMsg(e); err != nil {
			return err
		}
	}
	return nil
}

func (m *Tuples) UnmarshalMsg(d []byte) ([]byte, error) {
	l, d, err := msgp.ReadArrayHeaderBytes(d)
	if err != nil {
		return nil, err
	}

	*m = make([]Tuple, l)

	for i := uint32(0); i < l; i++ {
		if d, err = (*m)[i].UnmarshalMsg(d); err != nil {
			return nil, err
		}
	}

	return d, nil
}

type Member struct {
	Name  string
	Nonce string
	Val   uint
}

type Members []Member

func (m Members) EncodeMsg(e *msgp.Writer) error {
	e.WriteArrayHeader(uint32(len(m)))
	for i := 0; i < len(m); i++ {
		if err := m[i].EncodeMsg(e); err != nil {
			return err
		}
	}
	return nil
}

func (m *Members) UnmarshalMsg(d []byte) ([]byte, error) {
	l, d, err := msgp.ReadArrayHeaderBytes(d)
	if err != nil {
		return nil, err
	}

	*m = make([]Member, l)

	for i := uint32(0); i < l; i++ {
		if d, err = (*m)[i].UnmarshalMsg(d); err != nil {
			return nil, err
		}
	}

	return d, nil
}

func (m Members) Msgsize() int {
	s := 1

	for i := range m {
		s += m[i].Msgsize()
	}

	return s
}

func (m *Member) EncodeMsg(e *msgp.Writer) error {
	e.WriteArrayHeader(2)
	e.WriteString(m.Name)
	e.WriteUint(m.Val)
	return nil
}

func (m *Member) Msgsize() int {
	return 1 + msgp.StringSize(len(m.Name)) + msgp.IntSize(uint64(m.Val))
}

func (m *Member) UnmarshalMsg(d []byte) ([]byte, error) {
	var err error
	var l uint32
	if l, d, err = msgp.ReadArrayHeaderBytes(d); err != nil {
		return nil, err
	}
	if l != 2 {
		return nil, fmt.Errorf("array len doesn't match: %d", l)
	}
	if m.Name, d, err = msgp.ReadStringBytes(d); err != nil {
		return nil, err
	}
	if m.Val, d, err = msgp.ReadUintBytes(d); err != nil {
		return nil, err
	}
	return d, nil
}

type Tuple2 struct {
	Cid     uint
	Orig    string
	Members Members
}

func (c *Tuple2) EncodeMsg(e *msgp.Writer) error {
	e.WriteArrayHeader(3)
	e.WriteUint(c.Cid)
	e.WriteString(c.Orig)
	c.Members.EncodeMsg(e)
	return nil
}

func (c *Tuple2) Msgsize() int {
	return 1 + msgp.IntSize(uint64(c.Cid)) + msgp.StringSize(len(c.Orig)) + c.Members.Msgsize()
}

func (c *Tuple2) UnmarshalMsg(d []byte) ([]byte, error) {
	var err error
	var l uint32
	if l, d, err = msgp.ReadArrayHeaderBytes(d); err != nil {
		return nil, err
	}
	if l != 3 {
		return nil, fmt.Errorf("array len doesn't match: %d", l)
	}
	if c.Cid, d, err = msgp.ReadUintBytes(d); err != nil {
		return nil, err
	}
	if c.Orig, d, err = msgp.ReadStringBytes(d); err != nil {
		return nil, err
	}
	return c.Members.UnmarshalMsg(d)
}

type Tuples2 []Tuple2

func (m Tuples2) EncodeMsg(e *msgp.Writer) error {
	e.WriteArrayHeader(uint32(len(m)))
	for i := 0; i < len(m); i++ {
		if err := m[i].EncodeMsg(e); err != nil {
			return err
		}
	}
	return nil
}

func (m *Tuples2) UnmarshalMsg(d []byte) ([]byte, error) {
	l, d, err := msgp.ReadArrayHeaderBytes(d)
	if err != nil {
		return nil, err
	}

	*m = make([]Tuple2, l)

	for i := uint32(0); i < l; i++ {
		if d, err = (*m)[i].UnmarshalMsg(d); err != nil {
			return nil, err
		}
	}

	return d, nil
}

var server = "127.0.0.1:3013"
var spaceNo = uint32(512)
var spaceName = "test"
var indexNo = uint32(0)
var indexName = "primary"
var opts = Opts{
	Timeout: 1000 * time.Millisecond,
	User:    "test",
	Pass:    "test",
	//Concurrency: 32,
	//RateLimit: 4*1024,
}

const N = 500

type t struct {
	d []interface{}
}

func (c t) EncodeMsg(e *msgp.Writer) error {
	return e.WriteIntf(c.d)
}

func (c t) Msgsize() int {
	s := 1

	for _, i := range c.d {
		s += msgp.GuessSize(i)
	}

	return s
}

func Iface(args []interface{}) Body {
	return &t{
		d: args,
	}
}

func BenchmarkClientSerialTyped(b *testing.B) {
	var err error

	conn, err := Connect(server, opts)
	if err != nil {
		b.Fatalf("No connection available: %s", err)
		return
	}
	defer conn.Close()

	resp, err := conn.Replace(spaceNo, Iface([]interface{}{uint(1111), "hello", "world"}))
	if err != nil {
		b.Fatalf("No connection available")
	}

	resp.Release()

	b.RunParallel(func(p *testing.PB) {
		var r Tuples

		for p.Next() {
			resp, err := conn.Select(spaceNo, indexNo, 0, 1, IterEq, UintKey{1111})
			if err != nil {
				b.Fatalf("No connection available: %s", err)
			}

			r = r[:0]

			if _, err = r.UnmarshalMsg(resp.Data); err != nil {
				b.Fatalf("can't unmarshal response: %s", err)
			}

			if len(r) != 1 {
				b.Fatalf("bad length: %d", len(r))
			}

			resp.Release()
		}
	})
}

func BenchmarkClientFutureParallelTyped(b *testing.B) {
	var err error

	conn, err := Connect(server, opts)
	if err != nil {
		b.Fatalf("No connection available")
		return
	}
	defer conn.Close()

	resp, err := conn.Replace(spaceNo, Iface([]interface{}{uint(1111), "hello", "world"}))
	if err != nil {
		b.Fatalf("No connection available")
	}

	resp.Release()

	b.RunParallel(func(pb *testing.PB) {
		exit := false

		for !exit {
			var (
				fs [N]*Future
				j  int
			)

			for j = 0; j < N && pb.Next(); j++ {
				fs[j] = conn.SelectAsync(spaceNo, indexNo, 0, 1, IterEq, UintKey{1111})
			}

			exit = j < N

			var r Tuples

			for j > 0 {
				j--

				resp, err := fs[j].Get()
				if err != nil {
					b.Error(err)

					break
				}

				if _, err = r.UnmarshalMsg(resp.Data); err != nil {
					b.Error(err)

					break
				}

				if len(r) != 1 || r[0].Id != 1111 {
					b.Fatalf("Doesn't match %v", r)

					break
				}

				resp.Release()
			}
		}
	})
}

///////////////////

func TestClient(t *testing.T) {
	var resp Response
	var err error
	var conn *Connection

	conn, err = Connect(server, opts)
	if err != nil {
		t.Fatalf("Failed to connect: %s", err.Error())
		return
	}
	if conn == nil {
		t.Fatalf("conn is nil after Connect")
		return
	}
	defer conn.Close()

	// Ping
	resp, err = conn.Ping()
	if err != nil {
		t.Fatalf("Failed to Ping: %s", err.Error())
	}

	// Insert
	resp, err = conn.Insert(spaceNo, Iface([]interface{}{uint(1), "hello", "world"}))
	if err != nil {
		t.Fatalf("Failed to Insert: %s", err.Error())
	}
	if resp.Error != "" {
		t.Fatalf("Response Error: %s", resp.Error)
	}

	data, _, err := msgp.ReadIntfBytes(resp.Data)
	if err != nil {
		t.Fatalf("Response unpacking Error: %s", err)
	}

	if tpl, ok := data.([]interface{})[0].([]interface{}); !ok {
		t.Fatalf("Unexpected body of Insert")
	} else {
		if len(tpl) != 3 {
			t.Fatalf("Unexpected body of Insert (tuple len): %d", len(tpl))
		}
		if id, ok := tpl[0].(int64); !ok || id != 1 {
			t.Fatalf("Unexpected body of Insert (0)")
		}
		if h, ok := tpl[1].(string); !ok || h != "hello" {
			t.Fatalf("Unexpected body of Insert (1)")
		}
	}

	resp, err = conn.Insert(spaceNo, &Tuple{Id: 1, Msg: "hello", Name: "world"})
	if resp.Code != ErrTupleFound {
		t.Fatalf("Expected ErrTupleFound but got: %v", resp.Code)
	}

	// Delete
	resp, err = conn.Delete(spaceNo, indexNo, Iface([]interface{}{uint(1)}))
	if err != nil {
		t.Fatalf("Failed to Delete: %s", err.Error())
	}

	data, _, err = msgp.ReadIntfBytes(resp.Data)
	if err != nil {
		log.Println(resp.Code, resp.Data)
		t.Fatalf("Response unpacking Error: %s", err)
	}

	if tpl, ok := data.([]interface{})[0].([]interface{}); !ok {
		t.Fatalf("Unexpected body of Delete")
	} else {
		if len(tpl) != 3 {
			t.Fatalf("Unexpected body of Delete (tuple len)")
		}
		if id, ok := tpl[0].(int64); !ok || id != 1 {
			t.Fatalf("Unexpected body of Delete (0)")
		}
		if h, ok := tpl[1].(string); !ok || h != "hello" {
			t.Fatalf("Unexpected body of Delete (1)")
		}
	}
	resp, err = conn.Delete(spaceNo, indexNo, Iface([]interface{}{uint(101)}))
	if err != nil {
		t.Fatalf("Failed to Delete: %s", err.Error())
	}

	// Replace
	resp, err = conn.Replace(spaceNo, Iface([]interface{}{uint(2), "hello", "world"}))
	if err != nil {
		t.Fatalf("Failed to Replace: %s", err.Error())
	}

	resp, err = conn.Replace(spaceNo, Iface([]interface{}{uint(2), "hi", "planet"}))
	if err != nil {
		t.Fatalf("Failed to Replace (duplicate): %s", err.Error())
	}

	data, _, err = msgp.ReadIntfBytes(resp.Data)
	if err != nil {
		t.Fatalf("Response unpacking Error: %s", err)
	}

	if tpl, ok := data.([]interface{})[0].([]interface{}); !ok {
		t.Fatalf("Unexpected body of Replace")
	} else {
		if len(tpl) != 3 {
			t.Fatalf("Unexpected body of Replace (tuple len)")
		}
		if id, ok := tpl[0].(int64); !ok || id != 2 {
			t.Fatalf("Unexpected body of Replace (0)")
		}
		if h, ok := tpl[1].(string); !ok || h != "hi" {
			t.Fatalf("Unexpected body of Replace (1)")
		}
	}

	// Update
	resp, err = conn.Update(spaceNo, indexNo, Iface([]interface{}{uint(2)}), Iface([]interface{}{Iface([]interface{}{"=", 1, "bye"}), Iface([]interface{}{"#", 2, 1})}))
	if err != nil {
		t.Fatalf("Failed to Update: %s", err.Error())
	}

	data, _, err = msgp.ReadIntfBytes(resp.Data)
	if err != nil {
		t.Fatalf("Response unpacking Error: %s", err)
	}

	if tpl, ok := data.([]interface{})[0].([]interface{}); !ok {
		t.Fatalf("Unexpected body of Update")
	} else {
		if len(tpl) != 2 {
			t.Fatalf("Unexpected body of Update (tuple len)")
		}
		if id, ok := tpl[0].(int64); !ok || id != 2 {
			t.Fatalf("Unexpected body of Update (0)")
		}
		if h, ok := tpl[1].(string); !ok || h != "bye" {
			t.Fatalf("Unexpected body of Update (1)")
		}
	}

	// Upsert
	if strings.Compare(conn.Greeting.Version, "Tarantool 1.6.7") >= 0 {
		resp, err = conn.Upsert(spaceNo, Iface([]interface{}{uint(3), 1}), Iface([]interface{}{Iface([]interface{}{"+", 1, 1})}))
		if err != nil {
			t.Fatalf("Failed to Upsert (insert): %s", err.Error())
		}

		resp, err = conn.Upsert(spaceNo, Iface([]interface{}{uint(3), 1}), Iface([]interface{}{Iface([]interface{}{"+", 1, 1})}))
		if err != nil {
			t.Fatalf("Failed to Upsert (update): %s", err.Error())
		}
	}

	// Select
	for i := 10; i < 20; i++ {
		resp, err = conn.Replace(spaceNo, Iface([]interface{}{uint(i), fmt.Sprintf("val %d", i), "bla"}))
		if err != nil {
			t.Fatalf("Failed to Replace: %s", err.Error())
		}
	}
	resp, err = conn.Select(spaceNo, indexNo, 0, 1, IterEq, Iface([]interface{}{uint(10)}))
	if err != nil {
		t.Fatalf("Failed to Select: %s", err.Error())
	}

	data, _, err = msgp.ReadIntfBytes(resp.Data)
	if err != nil {
		t.Fatalf("Response unpacking Error: %s", err)
	}

	if tpl, ok := data.([]interface{})[0].([]interface{}); !ok {
		t.Fatalf("Unexpected body of Select")
	} else {
		if id, ok := tpl[0].(int64); !ok || id != 10 {
			t.Fatalf("Unexpected body of Select (0)")
		}
		if h, ok := tpl[1].(string); !ok || h != "val 10" {
			t.Fatalf("Unexpected body of Select (1)")
		}
	}

	// Select empty
	resp, err = conn.Select(spaceNo, indexNo, 0, 1, IterEq, Iface([]interface{}{uint(30)}))
	if err != nil {
		t.Fatalf("Failed to Select: %s", err.Error())
	}

	// Select Typed
	var tpl Tuples
	resp, err = conn.Select(spaceNo, indexNo, 0, 1, IterEq, Iface([]interface{}{uint(10)}))
	if err != nil {
		t.Fatalf("Failed to Select: %s", err.Error())
	}

	if _, err = tpl.UnmarshalMsg(resp.Data); err != nil {
		t.Fatalf("Failed to Select response parse: %s", err.Error())
	}

	if len(tpl) != 1 {
		t.Fatalf("Result len of SelectTyped != 1")
	} else {
		if tpl[0].Id != 10 {
			t.Fatalf("Bad value loaded from SelectTyped: %+v", tpl)
		}
	}

	// Select Typed for one tuple
	var tpl1 Tuples
	resp, err = conn.Select(spaceNo, indexNo, 0, 1, IterEq, Iface([]interface{}{uint(10)}))
	if err != nil {
		t.Fatalf("Failed to SelectTyped: %s", err.Error())
	}

	if _, err = tpl1.UnmarshalMsg(resp.Data); err != nil {
		t.Fatalf("Failed to Select response parse: %s", err.Error())
	}

	if len(tpl1) != 1 {
		t.Fatalf("Result len of SelectTyped != 1")
	} else {
		if tpl1[0].Id != 10 {
			t.Fatalf("Bad value loaded from SelectTyped")
		}
	}

	// Select Typed Empty
	var tpl2 Tuples
	resp, err = conn.Select(spaceNo, indexNo, 0, 1, IterEq, Iface([]interface{}{uint(30)}))
	if err != nil {
		t.Fatalf("Failed to SelectTyped: %s", err.Error())
	}

	if _, err = tpl2.UnmarshalMsg(resp.Data); err != nil {
		t.Fatalf("Failed to Select response parse: %s", err.Error())
	}

	if len(tpl2) != 0 {
		t.Fatalf("Result len of SelectTyped != 1")
	}

	// Call vs Call17
	resp, err = conn.Call("simple_incr", Iface([]interface{}{1}))

	data, _, err = msgp.ReadIntfBytes(resp.Data)
	if err != nil {
		t.Fatalf("Response unpacking Error: %s", err)
	}

	if data.([]interface{})[0].([]interface{})[0].(int64) != 2 {
		t.Fatalf("result is not {{1}} : %v", resp.Data)
	}

	resp, err = conn.Call17("simple_incr", Iface([]interface{}{1}))

	data, _, err = msgp.ReadIntfBytes(resp.Data)
	if err != nil {
		t.Fatalf("Response unpacking Error: %s", err)
	}

	if data.([]interface{})[0].(int64) != 2 {
		t.Fatalf("result is not {{1}} : %v", resp.Data)
	}

	// Eval
	resp, err = conn.Eval("return 5 + 6", Iface([]interface{}{}))
	if err != nil {
		t.Fatalf("Failed to Eval: %s", err.Error())
	}
	if len(resp.Data) < 1 {
		t.Fatalf("Response.Data is empty after Eval")
	}

	data, _, err = msgp.ReadIntfBytes(resp.Data)
	if err != nil {
		t.Fatalf("Response unpacking Error: %s", err)
	}

	val := data.([]interface{})[0].(int64)
	if val != 11 {
		t.Fatalf("5 + 6 == 11, but got %v", val)
	}
}

func TestSchema(t *testing.T) {
	var err error
	var conn *Connection

	conn, err = Connect(server, opts)
	if err != nil {
		t.Fatalf("Failed to connect: %s", err.Error())
		return
	}
	if conn == nil {
		t.Fatalf("conn is nil after Connect")
		return
	}
	defer conn.Close()

	// Schema
	schema := conn.Schema
	if schema.SpacesById == nil {
		t.Fatalf("schema.SpacesById is nil")
	}
	if schema.Spaces == nil {
		t.Fatalf("schema.Spaces is nil")
	}
	var space, space2 *Space
	var ok bool
	if space, ok = schema.SpacesById[514]; !ok {
		t.Fatalf("space with id = 514 was not found in schema.SpacesById")
	}
	if space2, ok = schema.Spaces["schematest"]; !ok {
		t.Fatalf("space with name 'schematest' was not found in schema.SpacesById")
	}
	if space != space2 {
		t.Fatalf("space with id = 514 and space with name schematest are different")
	}
	if space.Id != 514 {
		t.Fatalf("space 514 has incorrect Id")
	}
	if space.Name != "schematest" {
		t.Fatalf("space 514 has incorrect Name")
	}
	if !space.Temporary {
		t.Fatalf("space 514 should be temporary")
	}
	if space.Engine != "memtx" {
		t.Fatalf("space 514 engine should be memtx")
	}
	if space.FieldsCount != 7 {
		t.Fatalf("space 514 has incorrect fields count")
	}

	if space.FieldsById == nil {
		t.Fatalf("space.FieldsById is nill")
	}
	if space.Fields == nil {
		t.Fatalf("space.Fields is nill")
	}
	if len(space.FieldsById) != 6 {
		t.Fatalf("space.FieldsById len is incorrect")
	}
	if len(space.Fields) != 6 {
		t.Fatalf("space.Fields len is incorrect: %+v", space.Fields)
	}

	var field1, field2, field5, field1n, field5n *Field
	if field1, ok = space.FieldsById[1]; !ok {
		t.Fatalf("field id = 1 was not found")
	}
	if field2, ok = space.FieldsById[2]; !ok {
		t.Fatalf("field id = 2 was not found")
	}
	if field5, ok = space.FieldsById[5]; !ok {
		t.Fatalf("field id = 5 was not found")
	}

	if field1n, ok = space.Fields["name1"]; !ok {
		t.Fatalf("field name = name1 was not found")
	}
	if field5n, ok = space.Fields["name5"]; !ok {
		t.Fatalf("field name = name5 was not found")
	}
	if field1 != field1n || field5 != field5n {
		t.Fatalf("field with id = 1 and field with name 'name1' are different")
	}
	if field1.Name != "name1" {
		t.Fatalf("field 1 has incorrect Name")
	}
	if field1.Type != "unsigned" {
		t.Fatalf("field 1 has incorrect Type")
	}
	if field2.Name != "name2" {
		t.Fatalf("field 2 has incorrect Name")
	}
	if field2.Type != "string" {
		t.Fatalf("field 2 has incorrect Type")
	}

	if space.IndexesById == nil {
		t.Fatalf("space.IndexesById is nill")
	}
	if space.Indexes == nil {
		t.Fatalf("space.Indexes is nill")
	}
	if len(space.IndexesById) != 2 {
		t.Fatalf("space.IndexesById len is incorrect")
	}
	if len(space.Indexes) != 2 {
		t.Fatalf("space.Indexes len is incorrect")
	}

	var index0, index3, index0n, index3n *Index
	if index0, ok = space.IndexesById[0]; !ok {
		t.Fatalf("index id = 0 was not found")
	}
	if index3, ok = space.IndexesById[3]; !ok {
		t.Fatalf("index id = 3 was not found")
	}
	if index0n, ok = space.Indexes["primary"]; !ok {
		t.Fatalf("index name = primary was not found")
	}
	if index3n, ok = space.Indexes["secondary"]; !ok {
		t.Fatalf("index name = secondary was not found")
	}
	if index0 != index0n || index3 != index3n {
		t.Fatalf("index with id = 3 and index with name 'secondary' are different")
	}
	if index3.Id != 3 {
		t.Fatalf("index has incorrect Id")
	}
	if index0.Name != "primary" {
		t.Fatalf("index has incorrect Name")
	}
	if index0.Type != "hash" || index3.Type != "tree" {
		t.Fatalf("index has incorrect Type")
	}
	if !index0.Unique || index3.Unique {
		t.Fatalf("index has incorrect Unique")
	}
	if index3.Fields == nil {
		t.Fatalf("index.Fields is nil")
	}
	if len(index3.Fields) != 2 {
		t.Fatalf("index.Fields len is incorrect")
	}

	ifield1 := index3.Fields[0]
	ifield2 := index3.Fields[1]
	if ifield1 == nil || ifield2 == nil {
		t.Fatalf("index field is nil")
	}
	if ifield1.Id != 1 || ifield2.Id != 2 {
		t.Fatalf("index field has incorrect Id")
	}
	if (ifield1.Type != "num" && ifield1.Type != "unsigned") || (ifield2.Type != "STR" && ifield2.Type != "string") {
		t.Fatalf("index field has incorrect Type '%s'", ifield2.Type)
	}

	var rSpaceNo, rIndexNo uint32
	rSpaceNo, rIndexNo, err = schema.ResolveSpaceIndex("schematest", "secondary")
	if err != nil || rSpaceNo != 514 || rIndexNo != 3 {
		t.Fatalf("symbolic space and index params not resolved")
	}
	rSpaceNo, rIndexNo, err = schema.ResolveSpaceIndex("schematest", "")
	if err != nil || rSpaceNo != 514 || rIndexNo != 0 {
		t.Fatalf("symbolic space param not resolved")
	}
	rSpaceNo, rIndexNo, err = schema.ResolveSpaceIndex("schematest22", "secondary")
	if err == nil {
		t.Fatalf("resolveSpaceIndex didn't returned error with not existing space name")
	}
	rSpaceNo, rIndexNo, err = schema.ResolveSpaceIndex("schematest", "secondary22")
	if err == nil {
		t.Fatalf("resolveSpaceIndex didn't returned error with not existing index name")
	}
}

func TestComplexStructs(t *testing.T) {
	var err error
	var conn *Connection

	conn, err = Connect(server, opts)
	if err != nil {
		t.Fatalf("Failed to connect: %s", err.Error())
		return
	}
	if conn == nil {
		t.Fatalf("conn is nil after Connect")
		return
	}
	defer conn.Close()

	tuple := Tuple2{Cid: 777, Orig: "orig", Members: []Member{{"lol", "", 1}, {"wut", "", 3}}}
	_, err = conn.Replace(spaceNo, &tuple)
	if err != nil {
		t.Fatalf("Failed to insert: %s", err.Error())
		return
	}

	var tuples Tuples2
	resp, err := conn.Select(spaceNo, indexNo, 0, 1, IterEq, Iface([]interface{}{uint(777)}))
	if err != nil {
		t.Fatalf("Failed to selectTyped: %s", err.Error())
		return
	}

	if _, err = tuples.UnmarshalMsg(resp.Data); err != nil {
		t.Fatalf("Failed to selectTyped unparse: %s", err.Error())
		return
	}

	if len(tuples) != 1 {
		t.Fatalf("Failed to selectTyped: unexpected array length %d", len(tuples))
		return
	}

	if tuple.Cid != tuples[0].Cid || len(tuple.Members) != len(tuples[0].Members) || tuple.Members[1].Name != tuples[0].Members[1].Name {
		t.Fatalf("Failed to selectTyped: incorrect data")
		return
	}
}
