package hqlschema

import (
	"reflect"
	"strings"
)

type Schema struct {
	nodes map[string]*SchemaNode
}

type SchemaProvider interface {
	Tables() []string
	Schema() *Schema
}

func NewSchema() *Schema {
	schema := &Schema{}
	schema.nodes = make(map[string]*SchemaNode)
	return schema
}

func (s *Schema) RegisterNode(fieldName string, parent *SchemaNode, structType reflect.Type) *SchemaNode {
	sn := newSchemaNode(fieldName, parent, structType)
	s.nodes[sn.ID()] = sn
	return sn
}

func (s *Schema) Nodes() map[string]*SchemaNode {
	return s.nodes
}

func (s *Schema) SchemaNode(id string) (*SchemaNode, *SchemaNodeKeys) {
	snk, path := NewSchemaNodeKeys(id)
	return s.nodes[path], snk
}

func (s *Schema) CreateAttributeID(id string) *AttributeID {
	keys, path := NewSchemaNodeKeys(id)
	lastIndex := strings.LastIndex(path, ".")
	if lastIndex != -1 {
		tablePath := path[0:lastIndex]
		fieldName := path[lastIndex+1:]
		schemaNode, _ := s.SchemaNode(tablePath)
		if schemaNode == nil {
			return nil
		}
		for i := 0; i < schemaNode.structType.NumField(); i++ {
			colName := schemaNode.structType.Field(i).Name
			if strings.ToLower(colName) == strings.ToLower(fieldName) {
				sf := &AttributeID{}
				sf.fieldName = colName
				sf.instance = schemaNode.NewStructInstance(keys)
				return sf
			}
		}
		return nil
	}
	return nil
}
