package tarantool

import (
	"fmt"

	"github.com/pkg/errors"
)

// Schema contains information about spaces and indexes.
type Schema struct {
	Version uint
	// Spaces is map from space names to spaces
	Spaces map[string]*Space
	// SpacesById is map from space numbers to spaces
	SpacesById map[uint32]*Space
}

// Space contains information about tarantool space
type Space struct {
	Id        uint32
	Name      string
	Engine    string
	Temporary bool // Is this space temporaray?
	// Field configuration is not mandatory and not checked by tarantool.
	FieldsCount uint32
	Fields      map[string]*Field
	FieldsById  map[uint32]*Field
	// Indexes is map from index names to indexes
	Indexes map[string]*Index
	// IndexesById is map from index numbers to indexes
	IndexesById map[uint32]*Index
}

// Index contains information about index
type Index struct {
	Id     uint32
	Name   string
	Type   string
	Unique bool
	Fields []*IndexField
}

const (
	maxSchemas = 10000
	vspaceSpId = 281
	vindexSpId = 289
)

func (conn *Connection) loadSchema() (err error) {
	schema := new(Schema)
	schema.SpacesById = make(map[uint32]*Space)
	schema.Spaces = make(map[string]*Space)

	var response SpaceResponses

	// reload spaces
	resp, err := conn.Select(vspaceSpId, 0, 0, maxSchemas, IterAll, nil)
	if err != nil {
		resp.Release()

		return errors.Wrap(err, "can't select spaces")
	}

	if _, err = response.UnmarshalMsg(resp.Data); err != nil {
		resp.Release()

		return errors.Wrap(err, "can't unpack spaces")
	}

	resp.Release()

	for _, row := range response {
		space := new(Space)

		space.Id = row.Id
		space.Name = row.Name
		space.Engine = row.Engine
		space.FieldsCount = row.FieldsCount
		space.FieldsById = make(map[uint32]*Field)
		space.Fields = make(map[string]*Field)
		space.IndexesById = make(map[uint32]*Index)
		space.Indexes = make(map[string]*Index)
		space.Temporary = row.Flags.Temporary

		for i, field := range row.Fields {
			field.Id = uint32(i)

			space.FieldsById[uint32(i)] = field
			space.Fields[field.Name] = field
		}

		schema.SpacesById[space.Id] = space
		schema.Spaces[space.Name] = space
	}

	var indexes IndexResponses

	// reload indexes
	resp, err = conn.Select(vindexSpId, 0, 0, maxSchemas, IterAll, nil)
	if err != nil {
		resp.Release()

		return errors.Wrap(err, "can't select indexes")
	}

	if _, err = indexes.UnmarshalMsg(resp.Data); err != nil {
		resp.Release()

		return errors.Wrap(err, "can't unpack indexes")
	}

	resp.Release()

	for _, row := range indexes {
		index := new(Index)
		index.Id = row.IndexId
		index.Name = row.Name
		index.Type = row.Type
		index.Unique = row.Flags.Unique
		index.Fields = row.Fields

		schema.SpacesById[row.SpaceId].IndexesById[index.Id] = index
		schema.SpacesById[row.SpaceId].Indexes[index.Name] = index
	}

	conn.Schema = schema

	return nil
}

func (schema *Schema) ResolveSpaceIndex(s string, i string) (uint32, uint32, error) {
	var (
		spaceNo uint32
		indexNo uint32
	)

	if schema == nil {
		return 0, 0, fmt.Errorf("Schema is not loaded")
	}

	space, ok := schema.Spaces[s]
	if !ok {
		return 0, 0, fmt.Errorf("there is no space with name %s", s)
	}

	spaceNo = space.Id

	if i != "" {
		index, ok := space.Indexes[i]
		if !ok {
			return 0, 0, fmt.Errorf("space %s has not index with name %s", space.Name, i)
		}

		indexNo = index.Id
	}

	return spaceNo, indexNo, nil
}
