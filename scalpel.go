package scalpel

import (
	"errors"
	"fmt"
	"reflect"
	"scalpel/parse"
)

// SetField finds the field from data and set its value to newValue.
//
// data must be a pointer, or an error will be thrown.
func SetField(data interface{}, fieldPath []string, newValue string) error {
	typ := reflect.TypeOf(data)
	if typ.Kind() != reflect.Ptr {
		return errors.New("data must be a pointer")
	}
	val := reflect.ValueOf(data).Elem()
	return findFieldAndSet(val, fieldPath, newValue)
}

// FindFieldAndSet recursively reflects the fields of a struct according to fieldPath passed in
// and sets the field value to newVal.
//
// If success returns true, otherwise false.
func findFieldAndSet(val reflect.Value, fieldPath []string, newVal string, i ...int) error {
	var pathFlag int
	if len(i) > 0 {
		pathFlag = i[0]
	}
	param := fieldPath[pathFlag]
	nextLevelVal, err := findField(val, param)
	if err != nil {
		return err
	}
	if pathFlag >= len(fieldPath)-1 {
		err := setValue(nextLevelVal, newVal)
		if err != nil {
			return err
		}
		return nil
	}
	return findFieldAndSet(nextLevelVal, fieldPath, newVal, pathFlag+1)
}

// findField tries to find out v's field which depends on v's kind and param passed in.
//
// If fails to find, return false.
//
// If v's kind is reflect.Array or reflect.Slice, param should be the index.
//
// If v's kind is reflect.Map, param should be the key.
//
// If v's kind is reflect.Ptr, unwrap it and deal the value v points to.
func findField(v reflect.Value, param string) (reflect.Value, error) {
	t := v.Type()
	switch t.Kind() {
	case reflect.Array, reflect.Slice:
		idx, err := parse.ToIntE(param)
		if err != nil {
			return reflect.Value{}, errors.New(fmt.Sprintf("failed to parse %s to index of array", param))
		}
		if idx < 0 || idx >= v.Len() {
			return reflect.Value{}, errors.New(fmt.Sprintf("invalid index of array %d with length %d", idx, v.Len()))
		}
		return v.Index(idx), nil
	case reflect.Map:
		keyType := t.Key().Kind()
		swappedKey, err := typeSwap(keyType, param)
		if err != nil {
			return reflect.Value{}, errors.New(fmt.Sprintf("failed to find key %s in map", param))
		}
		return v.MapIndex(reflect.ValueOf(swappedKey)), nil
	case reflect.Ptr:
		ele := v.Elem()
		return findField(reflect.ValueOf(ele), param)
	case reflect.Struct:
		_, ok := t.FieldByName(param)
		if !ok {
			return reflect.Value{}, errors.New(fmt.Sprintf("failed to find field %s in %s", param, t.Name()))
		}
		return v.FieldByName(param), nil
	default:
		return reflect.Value{}, errors.New(fmt.Sprintf("unsupported kind %s, failed to find field %s", t.Kind(), param))
	}
}

// typeSwap returns the param in specific kind.
func typeSwap(kind reflect.Kind, param string) (interface{}, error) {
	switch kind {
	case reflect.Int:
		return parse.ToIntE(param)
	case reflect.Int8:
		return parse.ToInt8E(param)
	case reflect.Int32:
		return parse.ToInt32E(param)
	case reflect.Int64:
		return parse.ToInt64E(param)
	case reflect.Float32:
		return parse.ToFloat32E(param)
	case reflect.Float64:
		return parse.ToFloat64E(param)
	case reflect.Bool:
		return parse.ToBoolE(param)
	case reflect.String:
		return param, nil
	default:
		return nil, errors.New(fmt.Sprintf("unsupport kind to parse %s", kind))
	}
}

// setValue sets the value of v.
func setValue(v reflect.Value, newVal string) error {
	switch v.Type().Kind() {
	case reflect.String:
		v.SetString(newVal)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		swappedVal, err := typeSwap(reflect.Int64, newVal)
		if err != nil {
			return err
		}
		v.SetInt(swappedVal.(int64))
	case reflect.Float64, reflect.Float32:
		swappedVal, err := typeSwap(reflect.Float64, newVal)
		if err != nil {
			return err
		}
		v.SetFloat(swappedVal.(float64))
	case reflect.Bool:
		swappedVal, err := typeSwap(reflect.Bool, newVal)
		if err != nil {
			return err
		}
		v.SetBool(swappedVal.(bool))
	default:
		return errors.New(fmt.Sprintf("unsupported kind of value %s", v.Type().Kind()))
	}
	return nil
}
