package tarantool

import (
	"github.com/GoWebProd/msgp/msgp"
	"github.com/pkg/errors"
)

type request struct {
	requestId   uint32
	requestCode int32

	iterator uint32
	offset   uint32
	limit    uint32
	space    uint32
	index    uint32
	key      Body
	tuple    Body
	function string

	userName string
	method   string
	scramble []byte
}

func (z *request) EncodeMsg(en *msgp.Writer) error {
	switch z.requestCode {
	case AuthRequest:
		en.WriteMapHeader(2)
		en.WriteUint(KeyUserName)
		en.WriteString(z.userName)
		en.WriteUint(KeyTuple)
		en.WriteArrayHeader(2)
		en.WriteString(z.method)
		en.WriteStringFromBytes(z.scramble)

		return nil
	case PingRequest:
		return en.WriteMapHeader(0)
	case SelectRequest:
		en.WriteMapHeader(6)
		en.WriteUint64(KeyIterator)
		en.WriteUint32(z.iterator)
		en.WriteUint64(KeyOffset)
		en.WriteUint32(z.offset)
		en.WriteUint64(KeyLimit)
		en.WriteUint32(z.limit)
		en.WriteUint64(KeySpaceNo)
		en.WriteUint32(z.space)
		en.WriteUint64(KeyIndexNo)
		en.WriteUint32(z.index)
		en.WriteUint64(KeyKey)

		if z.key == nil {
			return en.WriteArrayHeader(0)
		}

		return z.key.EncodeMsg(en)
	case ReplaceRequest, InsertRequest:
		en.WriteMapHeader(2)
		en.WriteUint64(KeySpaceNo)
		en.WriteUint32(z.space)
		en.WriteUint64(KeyTuple)

		if z.tuple == nil {
			return en.WriteArrayHeader(0)
		}

		return z.tuple.EncodeMsg(en)
	case DeleteRequest:
		en.WriteMapHeader(3)
		en.WriteUint64(KeySpaceNo)
		en.WriteUint32(z.space)
		en.WriteUint64(KeyIndexNo)
		en.WriteUint32(z.index)
		en.WriteUint64(KeyKey)

		if z.key == nil {
			return en.WriteArrayHeader(0)
		}

		return z.key.EncodeMsg(en)
	case UpdateRequest:
		en.WriteMapHeader(4)
		en.WriteUint64(KeySpaceNo)
		en.WriteUint32(z.space)
		en.WriteUint64(KeyIndexNo)
		en.WriteUint32(z.index)
		en.WriteUint64(KeyKey)

		if z.key == nil {
			en.WriteArrayHeader(0)
		} else {
			z.key.EncodeMsg(en)
		}

		en.WriteUint64(KeyTuple)

		if z.tuple == nil {
			return en.WriteArrayHeader(0)
		}

		return z.tuple.EncodeMsg(en)
	case UpsertRequest:
		en.WriteMapHeader(3)
		en.WriteUint64(KeySpaceNo)
		en.WriteUint32(z.space)
		en.WriteUint64(KeyTuple)

		if z.key == nil {
			en.WriteArrayHeader(0)
		} else {
			z.key.EncodeMsg(en)
		}

		en.WriteUint64(KeyDefTuple)

		if z.tuple == nil {
			return en.WriteArrayHeader(0)
		}

		return z.tuple.EncodeMsg(en)
	case CallRequest, Call17Request:
		en.WriteMapHeader(2)
		en.WriteUint64(KeyFunctionName)
		en.WriteString(z.function)
		en.WriteUint64(KeyTuple)

		if z.tuple == nil {
			return en.WriteArrayHeader(0)
		}

		return z.tuple.EncodeMsg(en)
	case EvalRequest:
		en.WriteMapHeader(2)
		en.WriteUint64(KeyExpression)
		en.WriteString(z.function)
		en.WriteUint64(KeyTuple)

		if z.tuple == nil {
			return en.WriteArrayHeader(0)
		}

		return z.tuple.EncodeMsg(en)
	}

	return errors.Errorf("bad request: %d", z.requestCode)
}

func (z *request) Msgsize() int {
	switch z.requestCode {
	case AuthRequest:
		return 2 + msgp.IntSize(KeyUserName) + msgp.StringSize(len(z.userName)) + msgp.IntSize(KeyTuple) + msgp.StringSize(len(z.method)) + msgp.StringSize(len(z.scramble))
	case PingRequest:
		return 1
	case SelectRequest:
		s := 7 + msgp.IntSize(uint64(z.iterator)) + msgp.IntSize(uint64(z.offset)) + msgp.IntSize(uint64(z.limit)) + msgp.IntSize(uint64(z.space)) + msgp.IntSize(uint64(z.index))

		if z.key == nil {
			s += 1
		} else {
			s += z.key.Msgsize()
		}

		return s
	case ReplaceRequest, InsertRequest:
		s := 3 + msgp.IntSize(uint64(z.space))

		if z.tuple == nil {
			s += 1
		} else {
			s += z.tuple.Msgsize()
		}

		return s
	case DeleteRequest:
		s := 4 + msgp.IntSize(uint64(z.space)) + msgp.IntSize(uint64(z.index))

		if z.key == nil {
			s += 1
		} else {
			s += z.key.Msgsize()
		}

		return s
	case UpdateRequest:
		s := 5 + msgp.IntSize(uint64(z.space)) + msgp.IntSize(uint64(z.index))

		if z.key == nil {
			s += 1
		} else {
			s += z.key.Msgsize()
		}

		if z.tuple == nil {
			s += 1
		} else {
			s += z.tuple.Msgsize()
		}

		return s
	case UpsertRequest:
		s := 4 + msgp.IntSize(uint64(z.space))

		if z.key == nil {
			s += 1
		} else {
			s += z.key.Msgsize()
		}

		if z.tuple == nil {
			s += 1
		} else {
			s += z.tuple.Msgsize()
		}

		return s
	case CallRequest, Call17Request, EvalRequest:
		s := 3 + msgp.StringSize(len(z.function))

		if z.tuple == nil {
			s += 1
		} else {
			s += z.tuple.Msgsize()
		}

		return s
	}

	return 0
}

// Ping sends empty request to Tarantool to check connection.
func (conn *Connection) Ping() (resp Response, err error) {
	future := conn.newFuture(request{
		requestCode: PingRequest,
	})

	conn.queue <- future

	return future.Get()
}

// SelectAsync sends select request to tarantool and returns Future.
func (conn *Connection) SelectAsync(space, index, offset, limit, iterator uint32, key Body) *Future {
	future := conn.newFuture(request{
		requestCode: SelectRequest,

		space:    space,
		index:    index,
		offset:   offset,
		limit:    limit,
		iterator: iterator,
		key:      key,
	})

	conn.queue <- future

	return future
}

// Select performs select to box space.
//
// It is equal to conn.SelectAsync(...).Get()
func (conn *Connection) Select(space, index, offset, limit, iterator uint32, key Body) (resp Response, err error) {
	return conn.SelectAsync(space, index, offset, limit, iterator, key).Get()
}

// InsertAsync sends insert action to tarantool and returns Future.
// Tarantool will reject Insert when tuple with same primary key exists.
func (conn *Connection) InsertAsync(space uint32, tuple Body) *Future {
	future := conn.newFuture(request{
		requestCode: InsertRequest,

		space: space,
		tuple: tuple,
	})

	conn.queue <- future

	return future
}

// Insert performs insertion to box space.
// Tarantool will reject Insert when tuple with same primary key exists.
//
// It is equal to conn.InsertAsync(space, tuple).Get().
func (conn *Connection) Insert(space uint32, tuple Body) (resp Response, err error) {
	return conn.InsertAsync(space, tuple).Get()
}

// ReplaceAsync sends "insert or replace" action to tarantool and returns Future.
// If tuple with same primary key exists, it will be replaced.
func (conn *Connection) ReplaceAsync(space uint32, tuple Body) *Future {
	future := conn.newFuture(request{
		requestCode: ReplaceRequest,

		space: space,
		tuple: tuple,
	})

	conn.queue <- future

	return future
}

// Replace performs "insert or replace" action to box space.
// If tuple with same primary key exists, it will be replaced.
//
// It is equal to conn.ReplaceAsync(space, tuple).Get().
func (conn *Connection) Replace(space uint32, tuple Body) (resp Response, err error) {
	return conn.ReplaceAsync(space, tuple).Get()
}

// DeleteAsync sends deletion action to tarantool and returns Future.
// Future's result will contain array with deleted tuple.
func (conn *Connection) DeleteAsync(space, index uint32, key Body) *Future {
	future := conn.newFuture(request{
		requestCode: DeleteRequest,

		space: space,
		index: index,
		key:   key,
	})

	conn.queue <- future

	return future
}

// Delete performs deletion of a tuple by key.
// Result will contain array with deleted tuple.
//
// It is equal to conn.DeleteAsync(space, tuple).Get().
func (conn *Connection) Delete(space, index uint32, key Body) (resp Response, err error) {
	return conn.DeleteAsync(space, index, key).Get()
}

// Update sends deletion of a tuple by key and returns Future.
// Future's result will contain array with updated tuple.
func (conn *Connection) UpdateAsync(space, index uint32, key, ops Body) *Future {
	future := conn.newFuture(request{
		requestCode: UpdateRequest,

		space: space,
		index: index,
		key:   key,
		tuple: ops,
	})

	conn.queue <- future

	return future
}

// Update performs update of a tuple by key.
// Result will contain array with updated tuple.
//
// It is equal to conn.UpdateAsync(space, tuple).Get().
func (conn *Connection) Update(space, index uint32, key, ops Body) (resp Response, err error) {
	return conn.UpdateAsync(space, index, key, ops).Get()
}

// UpsertAsync sends "update or insert" action to tarantool and returns Future.
// Future's sesult will not contain any tuple.
func (conn *Connection) UpsertAsync(space uint32, key, ops Body) *Future {
	future := conn.newFuture(request{
		requestCode: UpsertRequest,

		space: space,
		key:   key,
		tuple: ops,
	})

	conn.queue <- future

	return future
}

// Upsert performs "update or insert" action of a tuple by key.
// Result will not contain any tuple.
//
// It is equal to conn.UpsertAsync(space, tuple, ops).Get().
func (conn *Connection) Upsert(space uint32, tuple, ops Body) (resp Response, err error) {
	return conn.UpsertAsync(space, tuple, ops).Get()
}

// CallAsync sends a call to registered tarantool function and returns Future.
// It uses request code for tarantool 1.6, so future's result is always array of arrays
func (conn *Connection) CallAsync(functionName string, args Body) *Future {
	future := conn.newFuture(request{
		requestCode: CallRequest,

		function: functionName,
		tuple:    args,
	})

	conn.queue <- future

	return future
}

// Call calls registered tarantool function.
// It uses request code for tarantool 1.6, so result is converted to array of arrays
//
// It is equal to conn.CallAsync(functionName, args).Get().
func (conn *Connection) Call(functionName string, args Body) (resp Response, err error) {
	return conn.CallAsync(functionName, args).Get()
}

// Call17Async sends a call to registered tarantool function and returns Future.
// It uses request code for tarantool 1.7, so future's result will not be converted
// (though, keep in mind, result is always array)
func (conn *Connection) Call17Async(functionName string, args Body) *Future {
	future := conn.newFuture(request{
		requestCode: Call17Request,

		function: functionName,
		tuple:    args,
	})

	conn.queue <- future

	return future
}

// Call17 calls registered tarantool function.
// It uses request code for tarantool 1.7, so result is not converted
// (though, keep in mind, result is always array)
//
// It is equal to conn.Call17Async(functionName, args).Get().
func (conn *Connection) Call17(functionName string, args Body) (resp Response, err error) {
	return conn.Call17Async(functionName, args).Get()
}

// EvalAsync sends a lua expression for evaluation and returns Future.
func (conn *Connection) EvalAsync(expr string, args Body) *Future {
	future := conn.newFuture(request{
		requestCode: EvalRequest,

		function: expr,
		tuple:    args,
	})

	conn.queue <- future

	return future
}

// Eval passes lua expression for evaluation.
//
// It is equal to conn.EvalAsync(expr, tuple).Get().
func (conn *Connection) Eval(expr string, args Body) (resp Response, err error) {
	return conn.EvalAsync(expr, args).Get()
}
