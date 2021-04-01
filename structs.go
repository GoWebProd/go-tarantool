package tarantool

type Field struct {
	Id   uint32 `msg:"id"`
	Name string `msg:"name"`
	Type string `msg:"type"`
}

//msgp:tuple IndexField
type IndexField struct {
	Id   uint32
	Type string
}

//msgp:tuple SpaceResponse
type SpaceResponse struct {
	Id          uint32
	Owner       uint32
	Name        string
	Engine      string
	FieldsCount uint32
	Flags       SpaceFlags
	Fields      []*Field
}

type SpaceResponses []SpaceResponse

type SpaceFlags struct {
	Temporary bool `msg:"temporary"`
}

//msgp:tuple IndexResponse
type IndexResponse struct {
	SpaceId uint32
	IndexId uint32
	Name    string
	Type    string
	Flags   IndexFlags
	Fields  []*IndexField
}

type IndexFlags struct {
	Unique bool `msg:"unique"`
}

type IndexResponses []IndexResponse
