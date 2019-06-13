package hqlschema

import (
	"bytes"
	"reflect"
	"strconv"
)

type InstanceID struct {
	node   *SchemaNode
	key    reflect.Value
	value  reflect.Value
	parent *InstanceID
}

func (instanceId *InstanceID) Value() reflect.Value {
	return instanceId.value
}

func (instanceId *InstanceID) Key() interface{} {
	return instanceId.key
}

func NewStructInstance(key, value reflect.Value, node *SchemaNode, parent *InstanceID) *InstanceID {
	si := &InstanceID{}
	si.key = key
	si.value = value
	si.node = node
	si.parent = parent
	return si
}

func (instanceId *InstanceID) newInstance(key, value reflect.Value) *InstanceID {
	return NewStructInstance(key, value, instanceId.node, instanceId.parent)
}

func (instanceId *InstanceID) keyStringValue() string {
	if instanceId.key.IsValid() {
		return instanceId.key.String()
	}
	return ""
}

func (instanceId *InstanceID) ID() string {
	key := instanceId.keyStringValue()
	if instanceId.parent == nil {
		if key == "" {
			return instanceId.node.id
		} else {
			buff := bytes.Buffer{}
			buff.WriteString(instanceId.node.id)
			buff.WriteString("[")
			buff.WriteString(key)
			buff.WriteString("]")
			return buff.String()
		}
	}
	buff := bytes.Buffer{}
	buff.WriteString(instanceId.parent.ID())
	buff.WriteString(".")
	buff.WriteString(instanceId.node.id)
	if key != "" {
		buff.WriteString("[")
		buff.WriteString(key)
		buff.WriteString("]")
	}
	return buff.String()
}

func (instanceId *InstanceID) newMapValue(mapValue reflect.Value) []*InstanceID {
	results := make([]*InstanceID, 0)
	newItem := reflect.New(instanceId.node.structType)
	newMap := reflect.MakeMap(reflect.MapOf(instanceId.key.Type(), newItem.Type()))
	newMap.SetMapIndex(instanceId.key, newItem)
	mapValue.Set(newMap)
	instance := NewStructInstance(instanceId.key, newItem, instanceId.node, instanceId.parent)
	results = append(results, instance)
	return results
}

func (instanceId *InstanceID) mapValue(mapValue reflect.Value, createIfNil bool) []*InstanceID {
	if createIfNil && mapValue.IsNil() && instanceId.key.IsValid() {
		return instanceId.newMapValue(mapValue)
	}

	results := make([]*InstanceID, 0)

	if instanceId.key.IsValid() {
		mapItem := mapValue.MapIndex(instanceId.key)
		if !createIfNil && !mapItem.IsValid() {
			return results
		}
		if !mapItem.IsValid() {
			mapItem = reflect.New(instanceId.node.structType)
			mapValue.SetMapIndex(instanceId.key, mapItem)
		}
		instance := NewStructInstance(instanceId.key, mapItem, instanceId.node, instanceId.parent)
		results = append(results, instance)
		return results
	}

	mapKeys := mapValue.MapKeys()
	var instance *InstanceID
	for _, mapKey := range mapKeys {
		mapItem := mapValue.MapIndex(mapKey)
		if mapItem.Kind() == reflect.Ptr {
			if mapItem.IsNil() {
				continue
			} else {
				mapItem = mapItem.Elem()
			}
		}
		if instance == nil {
			instance = NewStructInstance(mapKey, mapItem, instanceId.node, instanceId.parent)
		} else {
			instance = instance.newInstance(mapKey, mapItem)
		}
		results = append(results, instance)
	}
	return results
}

func (instanceId *InstanceID) newSliceValue(sliceValue reflect.Value) []*InstanceID {
	results := make([]*InstanceID, 0)
	newItem := reflect.New(instanceId.node.structType)
	index, err := strconv.Atoi(instanceId.key.String())
	if err != nil {
		return results
	}
	newSlice := reflect.MakeSlice(reflect.SliceOf(newItem.Type()), index+1, index+1)
	newSlice.Index(index).Set(newItem)
	sliceValue.Set(newSlice)
	instance := NewStructInstance(instanceId.key, newItem, instanceId.node, instanceId.parent)
	results = append(results, instance)
	return results
}

func (instanceId *InstanceID) enlargeSliceValue(sliceValue reflect.Value, newSize int) {
	newItem := reflect.New(instanceId.node.structType)
	newSlice := reflect.MakeSlice(reflect.SliceOf(newItem.Type()), newSize, newSize)
	for i := 0; i < sliceValue.Len(); i++ {
		newSlice.Index(i).Set(sliceValue.Index(i))
	}
	sliceValue.Set(newSlice)
}

func (instanceId *InstanceID) sliceValue(sliceValue reflect.Value, createIfNil bool) []*InstanceID {
	if createIfNil && sliceValue.IsNil() && instanceId.key.IsValid() {
		return instanceId.newSliceValue(sliceValue)
	}

	results := make([]*InstanceID, 0)

	if instanceId.key.IsValid() {
		index, err := strconv.Atoi(instanceId.key.String())
		if err != nil {
			return results
		}
		if sliceValue.Len() <= index {
			instanceId.enlargeSliceValue(sliceValue, index+1)
		}
		sliceItem := sliceValue.Index(index)
		if !createIfNil && (!sliceItem.IsValid() || sliceItem.IsNil()) {
			return results
		}
		if !sliceItem.IsValid() || sliceItem.IsNil() {
			sliceItem = reflect.New(instanceId.node.structType)
			sliceValue.Index(index).Set(sliceItem)
		}
		instance := NewStructInstance(instanceId.key, sliceItem, instanceId.node, instanceId.parent)
		results = append(results, instance)
		return results
	}

	var instance *InstanceID
	for i := 0; i < sliceValue.Len(); i++ {
		sliceItem := sliceValue.Index(i)
		if sliceItem.Kind() == reflect.Ptr {
			if sliceItem.IsNil() {
				continue
			} else {
				sliceItem = sliceItem.Elem()
			}
		}
		if instance == nil {
			instance = NewStructInstance(reflect.ValueOf(i), sliceItem, instanceId.node, instanceId.parent)
		} else {
			instance = instance.newInstance(reflect.ValueOf(i), sliceItem)
		}
		results = append(results, instance)
	}
	return results
}

func (instanceId *InstanceID) valueOf(value reflect.Value, createIfNil bool) []*InstanceID {
	if instanceId.parent == nil {
		value := createIfNilOrInvalid(value, instanceId.node.structType, createIfNil)
		instance := NewStructInstance(reflect.ValueOf(nil), value, instanceId.node, nil)
		return []*InstanceID{instance}
	}
	parents := instanceId.parent.valueOf(value, createIfNil)
	results := make([]*InstanceID, 0)
	for _, parent := range parents {
		parentValue := parent.value
		if parentValue.Kind() == reflect.Ptr {
			parentValue = parentValue.Elem()
		}
		myValue := parentValue.FieldByName(instanceId.node.fieldName)
		if myValue.Kind() == reflect.Map {
			results = append(results, instanceId.mapValue(myValue, createIfNil)...)
		} else if myValue.Kind() == reflect.Slice {
			results = append(results, instanceId.sliceValue(myValue, createIfNil)...)
		} else {
			if myValue.IsNil() && createIfNil {
				v := createIfNilOrInvalid(myValue, instanceId.node.structType, createIfNil)
				myValue.Set(v)
			}
			instance := instanceId.newInstance(reflect.ValueOf(nil), myValue)
			results = append(results, instance)
		}
	}
	return results
}

func (instanceId *InstanceID) ValueOf(value reflect.Value) []*InstanceID {
	return instanceId.valueOf(value, false)
}

func (instanceId *InstanceID) ValueOfOrCreate(value reflect.Value) []*InstanceID {
	return instanceId.valueOf(value, true)
}

func createIfNilOrInvalid(value reflect.Value, structType reflect.Type, createIfNil bool) reflect.Value {
	if createIfNil && (!value.IsValid() || value.Kind() == reflect.Ptr && value.IsNil()) {
		return reflect.New(structType)
	}
	return value
}
