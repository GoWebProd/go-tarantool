package tarantool

// Code generated by github.com/GoWebProd/msgp DO NOT EDIT.

import (
	"github.com/GoWebProd/msgp/msgp"
)

// MarshalMsg implements msgp.Marshaler
func (z Field) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "id"
	o = append(o, 0x83, 0xa2, 0x69, 0x64)
	o = msgp.AppendUint32(o, z.Id)
	// string "name"
	o = append(o, 0xa4, 0x6e, 0x61, 0x6d, 0x65)
	o = msgp.AppendString(o, z.Name)
	// string "type"
	o = append(o, 0xa4, 0x74, 0x79, 0x70, 0x65)
	o = msgp.AppendString(o, z.Type)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Field) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "id":
			z.Id, bts, err = msgp.ReadUint32Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Id")
				return
			}
		case "name":
			z.Name, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Name")
				return
			}
		case "type":
			z.Type, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Type")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z Field) Msgsize() (s int) {
	s = 1 + 3 + msgp.IntSize(uint64(z.Id)) + 5 + msgp.StringSize(len(z.Name)) + 5 + msgp.StringSize(len(z.Type))
	return
}

// MarshalMsg implements msgp.Marshaler
func (z IndexField) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// array header, size 2
	o = append(o, 0x92)
	o = msgp.AppendUint32(o, z.Id)
	o = msgp.AppendString(o, z.Type)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *IndexField) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	if zb0001 != 2 {
		err = msgp.ArrayError{Wanted: 2, Got: zb0001}
		return
	}
	z.Id, bts, err = msgp.ReadUint32Bytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Id")
		return
	}
	z.Type, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Type")
		return
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z IndexField) Msgsize() (s int) {
	s = 1 + msgp.IntSize(uint64(z.Id)) + msgp.StringSize(len(z.Type))
	return
}

// MarshalMsg implements msgp.Marshaler
func (z IndexFlags) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 1
	// string "unique"
	o = append(o, 0x81, 0xa6, 0x75, 0x6e, 0x69, 0x71, 0x75, 0x65)
	o = msgp.AppendBool(o, z.Unique)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *IndexFlags) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "unique":
			z.Unique, bts, err = msgp.ReadBoolBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Unique")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z IndexFlags) Msgsize() (s int) {
	s = 1 + 7 + msgp.BoolSize
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *IndexResponse) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// array header, size 6
	o = append(o, 0x96)
	o = msgp.AppendUint32(o, z.SpaceId)
	o = msgp.AppendUint32(o, z.IndexId)
	o = msgp.AppendString(o, z.Name)
	o = msgp.AppendString(o, z.Type)
	// map header, size 1
	// string "unique"
	o = append(o, 0x81, 0xa6, 0x75, 0x6e, 0x69, 0x71, 0x75, 0x65)
	o = msgp.AppendBool(o, z.Flags.Unique)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Fields)))
	for za0001 := range z.Fields {
		if z.Fields[za0001] == nil {
			o = msgp.AppendNil(o)
		} else {
			// array header, size 2
			o = append(o, 0x92)
			o = msgp.AppendUint32(o, z.Fields[za0001].Id)
			o = msgp.AppendString(o, z.Fields[za0001].Type)
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *IndexResponse) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	if zb0001 != 6 {
		err = msgp.ArrayError{Wanted: 6, Got: zb0001}
		return
	}
	z.SpaceId, bts, err = msgp.ReadUint32Bytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "SpaceId")
		return
	}
	z.IndexId, bts, err = msgp.ReadUint32Bytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "IndexId")
		return
	}
	z.Name, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Name")
		return
	}
	z.Type, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Type")
		return
	}
	var field []byte
	_ = field
	var zb0002 uint32
	zb0002, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Flags")
		return
	}
	for zb0002 > 0 {
		zb0002--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err, "Flags")
			return
		}
		switch msgp.UnsafeString(field) {
		case "unique":
			z.Flags.Unique, bts, err = msgp.ReadBoolBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Flags", "Unique")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err, "Flags")
				return
			}
		}
	}
	var zb0003 uint32
	zb0003, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Fields")
		return
	}
	if cap(z.Fields) >= int(zb0003) {
		z.Fields = (z.Fields)[:zb0003]
	} else {
		z.Fields = make([]*IndexField, zb0003)
	}
	for za0001 := range z.Fields {
		if msgp.IsNil(bts) {
			bts, err = msgp.ReadNilBytes(bts)
			if err != nil {
				return
			}
			z.Fields[za0001] = nil
		} else {
			if z.Fields[za0001] == nil {
				z.Fields[za0001] = new(IndexField)
			}
			var zb0004 uint32
			zb0004, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Fields", za0001)
				return
			}
			if zb0004 != 2 {
				err = msgp.ArrayError{Wanted: 2, Got: zb0004}
				return
			}
			z.Fields[za0001].Id, bts, err = msgp.ReadUint32Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Fields", za0001, "Id")
				return
			}
			z.Fields[za0001].Type, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Fields", za0001, "Type")
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *IndexResponse) Msgsize() (s int) {
	s = 1 + msgp.IntSize(uint64(z.SpaceId)) + msgp.IntSize(uint64(z.IndexId)) + msgp.StringSize(len(z.Name)) + msgp.StringSize(len(z.Type)) + 1 + 7 + msgp.BoolSize + msgp.ArrayHeaderSize
	for za0001 := range z.Fields {
		if z.Fields[za0001] == nil {
			s += msgp.NilSize
		} else {
			s += 1 + msgp.IntSize(uint64(z.Fields[za0001].Id)) + msgp.StringSize(len(z.Fields[za0001].Type))
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z IndexResponses) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendArrayHeader(o, uint32(len(z)))
	for za0001 := range z {
		o, err = z[za0001].MarshalMsg(o)
		if err != nil {
			err = msgp.WrapError(err, za0001)
			return
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *IndexResponses) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zb0002 uint32
	zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	if cap((*z)) >= int(zb0002) {
		(*z) = (*z)[:zb0002]
	} else {
		(*z) = make(IndexResponses, zb0002)
	}
	for zb0001 := range *z {
		bts, err = (*z)[zb0001].UnmarshalMsg(bts)
		if err != nil {
			err = msgp.WrapError(err, zb0001)
			return
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z IndexResponses) Msgsize() (s int) {
	s = msgp.ArrayHeaderSize
	for zb0003 := range z {
		s += z[zb0003].Msgsize()
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z SpaceFlags) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 1
	// string "temporary"
	o = append(o, 0x81, 0xa9, 0x74, 0x65, 0x6d, 0x70, 0x6f, 0x72, 0x61, 0x72, 0x79)
	o = msgp.AppendBool(o, z.Temporary)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *SpaceFlags) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "temporary":
			z.Temporary, bts, err = msgp.ReadBoolBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Temporary")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z SpaceFlags) Msgsize() (s int) {
	s = 1 + 10 + msgp.BoolSize
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *SpaceResponse) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// array header, size 7
	o = append(o, 0x97)
	o = msgp.AppendUint32(o, z.Id)
	o = msgp.AppendUint32(o, z.Owner)
	o = msgp.AppendString(o, z.Name)
	o = msgp.AppendString(o, z.Engine)
	o = msgp.AppendUint32(o, z.FieldsCount)
	// map header, size 1
	// string "temporary"
	o = append(o, 0x81, 0xa9, 0x74, 0x65, 0x6d, 0x70, 0x6f, 0x72, 0x61, 0x72, 0x79)
	o = msgp.AppendBool(o, z.Flags.Temporary)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Fields)))
	for za0001 := range z.Fields {
		if z.Fields[za0001] == nil {
			o = msgp.AppendNil(o)
		} else {
			// map header, size 3
			// string "id"
			o = append(o, 0x83, 0xa2, 0x69, 0x64)
			o = msgp.AppendUint32(o, z.Fields[za0001].Id)
			// string "name"
			o = append(o, 0xa4, 0x6e, 0x61, 0x6d, 0x65)
			o = msgp.AppendString(o, z.Fields[za0001].Name)
			// string "type"
			o = append(o, 0xa4, 0x74, 0x79, 0x70, 0x65)
			o = msgp.AppendString(o, z.Fields[za0001].Type)
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *SpaceResponse) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	if zb0001 != 7 {
		err = msgp.ArrayError{Wanted: 7, Got: zb0001}
		return
	}
	z.Id, bts, err = msgp.ReadUint32Bytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Id")
		return
	}
	z.Owner, bts, err = msgp.ReadUint32Bytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Owner")
		return
	}
	z.Name, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Name")
		return
	}
	z.Engine, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Engine")
		return
	}
	z.FieldsCount, bts, err = msgp.ReadUint32Bytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "FieldsCount")
		return
	}
	var field []byte
	_ = field
	var zb0002 uint32
	zb0002, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Flags")
		return
	}
	for zb0002 > 0 {
		zb0002--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err, "Flags")
			return
		}
		switch msgp.UnsafeString(field) {
		case "temporary":
			z.Flags.Temporary, bts, err = msgp.ReadBoolBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Flags", "Temporary")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err, "Flags")
				return
			}
		}
	}
	var zb0003 uint32
	zb0003, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Fields")
		return
	}
	if cap(z.Fields) >= int(zb0003) {
		z.Fields = (z.Fields)[:zb0003]
	} else {
		z.Fields = make([]*Field, zb0003)
	}
	for za0001 := range z.Fields {
		if msgp.IsNil(bts) {
			bts, err = msgp.ReadNilBytes(bts)
			if err != nil {
				return
			}
			z.Fields[za0001] = nil
		} else {
			if z.Fields[za0001] == nil {
				z.Fields[za0001] = new(Field)
			}
			var zb0004 uint32
			zb0004, bts, err = msgp.ReadMapHeaderBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Fields", za0001)
				return
			}
			for zb0004 > 0 {
				zb0004--
				field, bts, err = msgp.ReadMapKeyZC(bts)
				if err != nil {
					err = msgp.WrapError(err, "Fields", za0001)
					return
				}
				switch msgp.UnsafeString(field) {
				case "id":
					z.Fields[za0001].Id, bts, err = msgp.ReadUint32Bytes(bts)
					if err != nil {
						err = msgp.WrapError(err, "Fields", za0001, "Id")
						return
					}
				case "name":
					z.Fields[za0001].Name, bts, err = msgp.ReadStringBytes(bts)
					if err != nil {
						err = msgp.WrapError(err, "Fields", za0001, "Name")
						return
					}
				case "type":
					z.Fields[za0001].Type, bts, err = msgp.ReadStringBytes(bts)
					if err != nil {
						err = msgp.WrapError(err, "Fields", za0001, "Type")
						return
					}
				default:
					bts, err = msgp.Skip(bts)
					if err != nil {
						err = msgp.WrapError(err, "Fields", za0001)
						return
					}
				}
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *SpaceResponse) Msgsize() (s int) {
	s = 1 + msgp.IntSize(uint64(z.Id)) + msgp.IntSize(uint64(z.Owner)) + msgp.StringSize(len(z.Name)) + msgp.StringSize(len(z.Engine)) + msgp.IntSize(uint64(z.FieldsCount)) + 1 + 10 + msgp.BoolSize + msgp.ArrayHeaderSize
	for za0001 := range z.Fields {
		if z.Fields[za0001] == nil {
			s += msgp.NilSize
		} else {
			s += 1 + 3 + msgp.IntSize(uint64(z.Fields[za0001].Id)) + 5 + msgp.StringSize(len(z.Fields[za0001].Name)) + 5 + msgp.StringSize(len(z.Fields[za0001].Type))
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z SpaceResponses) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendArrayHeader(o, uint32(len(z)))
	for za0001 := range z {
		o, err = z[za0001].MarshalMsg(o)
		if err != nil {
			err = msgp.WrapError(err, za0001)
			return
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *SpaceResponses) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zb0002 uint32
	zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	if cap((*z)) >= int(zb0002) {
		(*z) = (*z)[:zb0002]
	} else {
		(*z) = make(SpaceResponses, zb0002)
	}
	for zb0001 := range *z {
		bts, err = (*z)[zb0001].UnmarshalMsg(bts)
		if err != nil {
			err = msgp.WrapError(err, zb0001)
			return
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z SpaceResponses) Msgsize() (s int) {
	s = msgp.ArrayHeaderSize
	for zb0003 := range z {
		s += z[zb0003].Msgsize()
	}
	return
}
