package hqlschema

import (
	"bytes"
	"reflect"
	"strings"
)

type SchemaNode struct {
	fieldName  string
	id         string
	parent     *SchemaNode
	structType reflect.Type
}

func newSchemaNode(fieldName string, parent *SchemaNode, structType reflect.Type) *SchemaNode {
	sn := &SchemaNode{}
	sn.parent = parent
	sn.fieldName = fieldName
	sn.id = strings.ToLower(fieldName)
	sn.structType = structType
	if sn.parent == nil {
		sn.id = strings.ToLower(sn.structType.Name())
	}
	return sn
}

func (sn *SchemaNode) FieldName() string {
	return sn.fieldName
}

func (sn *SchemaNode) Parent() *SchemaNode {
	return sn.parent
}

func (sn *SchemaNode) Type() reflect.Type {
	return sn.structType
}

func (sn *SchemaNode) NewInstance() interface{} {
	return sn.NewValue().Interface()
}

func (sn *SchemaNode) NewValue() reflect.Value {
	return reflect.New(sn.structType)
}

func (sn *SchemaNode) ID() string {
	if sn.parent == nil {
		return sn.id
	}
	buff := bytes.Buffer{}
	buff.WriteString(sn.parent.ID())
	buff.WriteString(".")
	buff.WriteString(sn.id)
	return buff.String()
}

func (sn *SchemaNode) NewStructInstance(keys *SchemaNodeKeys) *InstanceID {
	if sn.parent == nil {
		si := &InstanceID{}
		si.node = sn
		return si
	}
	parent := sn.parent.NewStructInstance(keys)
	si := &InstanceID{}
	si.node = sn
	si.parent = parent
	if keys != nil {
		key := keys.Key(si.node.ID())
		if key != "" {
			si.key = reflect.ValueOf(key)
		}
	}
	return si
}

func (sn *SchemaNode) CreateFeildID(fieldName string) string {
	buff := bytes.Buffer{}
	buff.WriteString(sn.id)
	buff.WriteString(".")
	buff.WriteString(fieldName)
	return buff.String()
}
