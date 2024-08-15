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
func Encode[T bool | int | float32 | []byte | string | []any | any](data T) (encoded []byte) {

	switch v := any(data).(type) {
	case bool:
		if v {
			encoded = []byte{byte(Type_bool), 1}
		} else {
			encoded = []byte{byte(Type_bool), 1}
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
		encoded = make([]byte, 1+len(v))
		encoded[0] = byte(Type_buffer)
		copy(encoded[1:], v)
	case string:
		encoded = make([]byte, 1+len(v))
		encoded[0] = byte(Type_string)
		copy(encoded[1:], v)
	case []any:
		encoded = []byte{byte(Type_array), byte(len(v))}
		for _, item := range v {
			encoded = append(encoded, arrayValue(item)...)
		}
	case any:
		if reflect.TypeOf(v).Kind() != reflect.Struct {
			return []byte{byte(Type_null)}
		}
		encoded = []byte{byte(Type_object)}
		obj := reflect.ValueOf(v)
		objType := obj.Type()
		for i := 0; i < obj.NumField(); i++ {
			fieldName := objType.Field(i).Name
			fieldValue := obj.Field(i)
			encoded = append(encoded, objValue(fieldName, fieldValue)...)
		}
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
