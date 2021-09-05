package parse

import (
	"strconv"
)

func ToInt8(v string) int8 {
	if v == "" {
		return 0
	}
	x, _ := strconv.Atoi(v)
	return int8(x)
}

func ToInt8E(v string) (int8, error) {
	if v == "" {
		return 0, nil
	}
	x, err := strconv.Atoi(v)
	return int8(x), err
}

func ToInt(v string) int {
	if v == "" {
		return 0
	}
	x, _ := strconv.Atoi(v)
	return x
}

func ToIntE(v string) (int, error) {
	if v == "" {
		return 0, nil
	}
	return strconv.Atoi(v)
}

func ToInt32(v string) int32 {
	if v == "" {
		return 0
	}
	x, _ := strconv.Atoi(v)
	return int32(x)
}

func ToInt32E(v string) (int32, error) {
	if v == "" {
		return 0, nil
	}
	x, err := strconv.Atoi(v)
	return int32(x), err
}

func ToInt64(v string) int64 {
	if v == "" {
		return 0
	}
	x, _ := strconv.Atoi(v)
	return int64(x)
}

func ToInt64E(v string) (int64, error) {
	if v == "" {
		return 0, nil
	}
	x, err := strconv.Atoi(v)
	return int64(x), err
}

func ToFloat32(v string) float32 {
	if v == "" {
		return 0
	}
	x, _ := strconv.ParseFloat(v, 32)
	return float32(x)
}

func ToFloat64(v string) float64 {
	if v == "" {
		return 0
	}
	x, _ := strconv.ParseFloat(v, 64)
	return x
}

func ToFloat32E(v string) (float32, error) {
	if v == "" {
		return 0, nil
	}
	x, err := strconv.ParseFloat(v, 32)
	return float32(x), err
}

func ToFloat64E(v string) (float64, error) {
	if v == "" {
		return 0, nil
	}
	x, err := strconv.ParseFloat(v, 64)
	return x, err
}

func ToString(i interface{}) string {
	switch v := i.(type) {
	case int:
		return strconv.Itoa(v)
	case int8:
		return strconv.Itoa(int(v))
	case int16:
		return strconv.Itoa(int(v))
	case int32:
		return strconv.Itoa(int(v))
	case int64:
		return strconv.Itoa(int(v))
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case string:
		return v
	case bool:
		return strconv.FormatBool(v)
	default:
		panic(i)
	}
}

func ToBool(v string) bool {
	if v == "" {
		return false
	}
	val, _ := strconv.ParseBool(v)
	return val
}

func ToBoolE(v string) (bool, error) {
	if v == "" {
		return false, nil
	}
	return strconv.ParseBool(v)
}

func ToDouble(v string) float64 {
	x, _ := strconv.ParseFloat(v, 64)
	return x
}

func ToFloat(v string) float32 {
	x, _ := strconv.ParseFloat(v, 32)
	return float32(x)
}

func ToBytes(v string) []byte {
	return []byte(v)
}
