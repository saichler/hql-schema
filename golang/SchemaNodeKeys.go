package hqlschema

import (
	"bytes"
	"strings"
)

type SchemaNodeKeys struct {
	keys map[string]string
}

func (snk *SchemaNodeKeys) Key(path string) string {
	return snk.keys[path]
}

func (snk *SchemaNodeKeys) Strings() string {
	buff := bytes.Buffer{}
	for k, v := range snk.keys {
		buff.WriteString(k)
		buff.WriteString("\\")
		buff.WriteString(v)
		buff.WriteString("\n")
	}
	return buff.String()
}

func NewSchemaNodeKeys(id string) (*SchemaNodeKeys, string) {
	path := strings.ToLower(id)
	skm := &SchemaNodeKeys{}
	skm.keys = make(map[string]string)
	from := strings.Index(path, "[")
	for from != -1 {
		to := strings.Index(path, "]")
		prefix := path[0:from]
		suffix := path[to+1:]
		key := path[from+1 : to]
		skm.keys[prefix] = key
		buff := bytes.Buffer{}
		buff.WriteString(prefix)
		buff.WriteString(suffix)
		path = buff.String()
		from = strings.Index(path, "[")
	}
	return skm, path
}
