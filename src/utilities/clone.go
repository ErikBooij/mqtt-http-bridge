package utilities

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"reflect"
)

func DeepCopy[T any](original T) (T, error) {
	val := reflect.ValueOf(original)

	if val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface || val.Kind() == reflect.Slice || val.Kind() == reflect.Map || val.Kind() == reflect.Chan || val.Kind() == reflect.Func {
		if val.IsNil() {
			var zero T
			return zero, nil
		}
	}

	var clone T

	// Create a buffer to hold the encoded data
	var buf bytes.Buffer

	// Create a new encoder and encode the original value
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(original)
	if err != nil {
		return clone, fmt.Errorf("failed to encode: %v", err)
	}

	// Create a new decoder and decode into the clone
	dec := gob.NewDecoder(&buf)
	err = dec.Decode(&clone)
	if err != nil {
		return clone, fmt.Errorf("failed to decode: %v", err)
	}

	return clone, nil
}
