package encoding

import (
	"encoding/binary"
	"math"
	"reflect"
)

type Type uint8

const (
	Type_null Type = iota
	Type_bool
	Type_uint8
	Type_int32
	Type_int64
	Type_float32
	Type_buffer
	Type_string
	Type_array
	Type_object
)

// Encode implements a custom encoding scheme for the wasm worker
func Encode[T bool | int | float32 | []byte | string | []any | map[string]any | any](data T) (encoded []byte) {
	value := reflect.ValueOf(data)
	switch value.Kind() {
	case reflect.Array, reflect.Slice:
		encoded = make([]byte, 5)
		encoded[0] = byte(Type_array)
		binary.LittleEndian.PutUint32(encoded[1:5], uint32(value.Len()))
		for i := 0; i < value.Len(); i++ {
			encoded = append(encoded, arrayValue(value.Index(i).Interface())...)
		}
		return encoded
	case reflect.Struct:
		objType := value.Type()
		data := make([]byte, 0)
		for i := 0; i < value.NumField(); i++ {
			info := objType.Field(i)
			if !info.IsExported() {
				continue
			}
			fieldValue := value.Field(i)
			data = append(data, objValue(info.Name, fieldValue.Interface())...)
		}
		encoded = make([]byte, 5+len(data))
		encoded[0] = byte(Type_object)
		binary.LittleEndian.PutUint32(encoded[1:5], uint32(len(data)))
		copy(encoded[5:], data)
		return encoded
	case reflect.Map:
		data := make([]byte, 0)
		for _, key := range value.MapKeys() {
			data = append(data, objValue(key.String(), value.MapIndex(key).Interface())...)
		}
		encoded = make([]byte, 5+len(data))
		encoded[0] = byte(Type_object)
		binary.LittleEndian.PutUint32(encoded[1:5], uint32(len(data)))
		copy(encoded[5:], data)
		return encoded
	}

	switch v := any(data).(type) {
	case bool:
		if v {
			encoded = []byte{byte(Type_bool), 1}
		} else {
			encoded = []byte{byte(Type_bool), 0}
		}
	case int:
		if v > math.MaxInt32 || v < math.MinInt32 {
			// encode int64
			encoded = make([]byte, 9)
			encoded[0] = byte(Type_int64)
			binary.LittleEndian.PutUint64(encoded[1:], uint64(v))
		} else if v > math.MaxUint8 || v < 0 {
			// encode int32
			encoded = make([]byte, 5)
			encoded[0] = byte(Type_int32)
			binary.LittleEndian.PutUint32(encoded[1:], uint32(v))
		} else {
			// encode uint8
			encoded = []byte{byte(Type_uint8), byte(v)}
		}
	case float32:
		encoded = make([]byte, 5)
		encoded[0] = byte(Type_float32)
		binary.LittleEndian.PutUint32(encoded[1:], math.Float32bits(v))
	case []byte:
		encoded = make([]byte, 5+len(v))
		encoded[0] = byte(Type_buffer)
		itemLen := uint32(len(v))
		binary.LittleEndian.PutUint32(encoded[1:5], itemLen)
		copy(encoded[5:], v)
	case string:
		encoded = make([]byte, 5+len(v))
		encoded[0] = byte(Type_string)
		itemLen := uint32(len(v))
		binary.LittleEndian.PutUint32(encoded[1:5], itemLen)
		copy(encoded[5:], v)
	case any:
		encoded = []byte{byte(Type_null)}
	}
	return encoded
}

func arrayValue(item any) []byte {
	itemBytes := Encode(item)
	itemLen := uint32(len(itemBytes))
	itemLenBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(itemLenBytes, itemLen)
	return append(itemLenBytes, itemBytes...)
}

func objValue(key string, value any) []byte {
	keyLen := uint32(len(key))
	keyBytes := make([]byte, 5+keyLen)
	keyBytes[0] = byte(Type_int32)
	binary.LittleEndian.PutUint32(keyBytes[1:5], keyLen)
	copy(keyBytes[5:], key)
	valueBytes := Encode(value)
	valueLen := uint32(len(valueBytes))
	valueLenBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(valueLenBytes, valueLen)
	return append(append(keyBytes, valueLenBytes...), valueBytes...)
}
