package hqlschema

import (
	"bytes"
	"reflect"
)

type AttributeID struct {
	instance  *InstanceID
	fieldName string
}

func (attributeId *AttributeID) Instance() *InstanceID {
	return attributeId.instance
}

func (attributeId *AttributeID) FieldName() string {
	return attributeId.fieldName
}

func (attributeId *AttributeID) ID() string {
	buff := bytes.Buffer{}
	buff.WriteString(attributeId.instance.ID())
	buff.WriteString(".")
	buff.WriteString(attributeId.fieldName)
	return buff.String()
}

func (attributeId *AttributeID) ValueOf(root reflect.Value) []reflect.Value {
	structInstances := attributeId.instance.ValueOf(root)
	results := make([]reflect.Value, 0)
	for _, instance := range structInstances {
		instanceValue := instance.value
		if instanceValue.Kind() == reflect.Ptr {
			if instanceValue.IsNil() {
				continue
			} else {
				instanceValue = instanceValue.Elem()
			}
		}
		if instanceValue.Kind() == reflect.Slice {
			for i := 0; i < instanceValue.Len(); i++ {
				elem := instanceValue.Index(i)
				if elem.Kind() == reflect.Ptr {
					elem = elem.Elem()
				}
				results = append(results, elem.FieldByName(attributeId.fieldName))
			}
		} else {
			results = append(results, instanceValue.FieldByName(attributeId.fieldName))
		}
	}
	return results
}

func (attributeId *AttributeID) SetValue(rootAny, any interface{}) {
	root := reflect.ValueOf(rootAny)
	structInstances := attributeId.instance.valueOf(root, true)
	results := make([]reflect.Value, 0)
	for _, instance := range structInstances {
		instanceValue := instance.value
		if instanceValue.Kind() == reflect.Ptr {
			if instanceValue.IsNil() {
				continue
			} else {
				instanceValue = instanceValue.Elem()
			}
		}
		if instanceValue.Kind() == reflect.Slice {
			for i := 0; i < instanceValue.Len(); i++ {
				elem := instanceValue.Index(i)
				if elem.Kind() == reflect.Ptr {
					elem = elem.Elem()
				}
				results = append(results, elem.FieldByName(attributeId.fieldName))
			}
		} else {
			v := instanceValue.FieldByName(attributeId.fieldName)
			v.Set(reflect.ValueOf(any))
		}
	}
}
