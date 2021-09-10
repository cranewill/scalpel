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
	err := findFieldAndSet(val, fieldPath, newValue)
	return err
}

// FindFieldAndSet recursively reflects the fields of a struct according to fieldPath passed in
// and sets the field value to newVal.
func findFieldAndSet(val reflect.Value, fieldPath []string, newVal string, i ...int) error {
	var pathFlag int
	if len(i) > 0 {
		pathFlag = i[0]
	}
	param := fieldPath[pathFlag]
	nextLayerField, err := findField(val, param)
	if err != nil {
		return err
	}

	nextLayerFieldType := nextLayerField.Type()
	if nextLayerFieldType.Kind() == reflect.Map {
		nextParam := fieldPath[pathFlag+1]
		nextKey, err := typeParse(nextLayerFieldType.Key().Kind(), nextParam)
		if err != nil {
			return err
		}
		mapVal := nextLayerField.MapIndex(reflect.ValueOf(nextKey))
		if mapVal.IsZero() {
			return errors.New(fmt.Sprintf("failed to find key %s in map", param))
		}
		_newVal, err := createNew(mapVal, fieldPath, pathFlag+1, newVal, true)
		if err != nil {
			return err
		}
		nextLayerField.SetMapIndex(reflect.ValueOf(nextKey), _newVal)
		return nil
	} else {
		if pathFlag >= len(fieldPath)-1 {
			err := setValue(nextLayerField, newVal)
			if err != nil {
				return err
			}
			return nil
		}
		return findFieldAndSet(nextLayerField, fieldPath, newVal, pathFlag+1)
	}
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
		swappedKey, err := typeParse(keyType, param)
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

// typeParse returns the param in specific kind.
func typeParse(kind reflect.Kind, param string) (interface{}, error) {
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
		swappedVal, err := typeParse(reflect.Int64, newVal)
		if err != nil {
			return err
		}
		v.SetInt(swappedVal.(int64))
	case reflect.Float64, reflect.Float32:
		swappedVal, err := typeParse(reflect.Float64, newVal)
		if err != nil {
			return err
		}
		v.SetFloat(swappedVal.(float64))
	case reflect.Bool:
		swappedVal, err := typeParse(reflect.Bool, newVal)
		if err != nil {
			return err
		}
		v.SetBool(swappedVal.(bool))
	default:
		return errors.New(fmt.Sprintf("unsupported kind of value %s", v.Type().Kind()))
	}
	return nil
}

// createNew creates value recursively, and it will set the correct field value to newVal
func createNew(val reflect.Value, fieldPath []string, flag int, newVal string, pathCorrect bool) (reflect.Value, error) {
	var param string
	if flag+1 < len(fieldPath) {
		param = fieldPath[flag+1]
	}
	change := flag == len(fieldPath)-1 && pathCorrect
	v := reflect.Value{}
	typ := val.Type()
	switch typ.Kind() {
	case reflect.String:
		v = reflect.New(typ).Elem()
		if change {
			v.SetString(newVal)
		} else {
			v.Set(val)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v = reflect.New(typ).Elem()
		swappedVal, err := typeParse(reflect.Int64, newVal)
		if err != nil {
			return reflect.Value{}, err
		}
		if change {
			v.SetInt(swappedVal.(int64))
		} else {
			v.Set(val)
		}
	case reflect.Float64, reflect.Float32:
		v = reflect.New(typ).Elem()
		swappedVal, err := typeParse(reflect.Float64, newVal)
		if err != nil {
			return reflect.Value{}, err
		}
		if change {
			v.SetFloat(swappedVal.(float64))
		} else {
			v.Set(val)
		}
	case reflect.Bool:
		v = reflect.New(typ).Elem()
		swappedVal, err := typeParse(reflect.Bool, newVal)
		if err != nil {
			return reflect.Value{}, err
		}
		if change {
			v.SetBool(swappedVal.(bool))
		} else {
			v.Set(val)
		}
	case reflect.Struct:
		v = reflect.New(typ).Elem()
		for i := 0; i < typ.NumField(); i++ {
			field := val.Field(i)
			fieldName := typ.Field(i).Name
			nField, err := createNew(field, fieldPath, flag+1, newVal, pathCorrect && param == fieldName)
			if err != nil {
				return reflect.Value{}, nil
			}
			v.Field(i).Set(nField)
		}
	case reflect.Array, reflect.Slice:
		var idx = -1
		if pathCorrect {
			_idx, err := parse.ToIntE(param)
			if err != nil {
				return reflect.Value{}, err
			}
			idx = _idx
		}
		v = reflect.MakeSlice(typ, 0, 0)
		for i := 0; i < val.Len(); i++ {
			unit, err := createNew(val.Index(i), fieldPath, flag+1, newVal, pathCorrect && i == idx)
			if err != nil {
				return reflect.Value{}, nil
			}
			v = reflect.Append(v, unit)
		}
	case reflect.Map:
		var key interface{}
		if pathCorrect {
			_key, err := typeParse(typ.Key().Kind(), param)
			if err != nil {
				return reflect.Value{}, err
			}
			key = _key
		}
		v = reflect.MakeMap(typ)
		itr := val.MapRange()
		for itr.Next() {
			entryKey := itr.Key()
			entryValue := itr.Value()
			unit, err := createNew(entryValue, fieldPath, flag+1, newVal, pathCorrect && key == entryKey.Interface())
			if err != nil {
				return reflect.Value{}, err
			}
			v.SetMapIndex(entryKey, unit)
		}
	}
	return v, nil
}
