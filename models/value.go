package models

import (
	"fmt"
	"reflect"
)

type Kind int

const (
	KindNil Kind = iota
	KindInt
	KindIntList
	KindString
	KindStringList
	KindBool
	KindBoolList
	KindMap
)

type Value struct {
	K  Kind
	I  int
	IL []int
	S  string
	SL []string
	B  bool
	BL []bool
}

func VInt(i int) Value    { return Value{K: KindInt, I: i} }
func VStr(s string) Value { return Value{K: KindString, S: s} }
func VBool(b bool) Value  { return Value{K: KindBool, B: b} }
func VNil() Value         { return Value{K: KindNil} }

func VList[T any](xs []T) (Value, error) {
	var zero T
	et := reflect.TypeOf(zero)
	if et == nil {
		return Value{K: KindNil}, fmt.Errorf("unsupported list type nil")
	}

	switch et.Kind() {
	case reflect.Int:
		if xs == nil {
			return Value{K: KindIntList, IL: nil}, nil
		}
		il := make([]int, len(xs))
		for i, v := range xs {
			il[i] = int(reflect.ValueOf(v).Int())
		}
		return Value{K: KindIntList, IL: il}, nil

	case reflect.String:
		if xs == nil {
			return Value{K: KindStringList, SL: nil}, nil
		}
		sl := make([]string, len(xs))
		for i, v := range xs {
			sl[i] = reflect.ValueOf(v).String()
		}
		return Value{K: KindStringList, SL: sl}, nil

	case reflect.Bool:
		if xs == nil {
			return Value{K: KindBoolList, BL: nil}, nil
		}
		bl := make([]bool, len(xs))
		for i, v := range xs {
			bl[i] = reflect.ValueOf(v).Bool()
		}
		return Value{K: KindBoolList, BL: bl}, nil

	default:
		return Value{K: KindNil}, fmt.Errorf("unsupported list type")
	}
}

func FromAny(x any) (Value, error) {
	switch v := x.(type) {
	case Value:
		return v, nil
	case nil:
		return VNil(), nil
	case int:
		return VInt(v), nil
	case int8:
		return VInt(int(v)), nil
	case int16:
		return VInt(int(v)), nil
	case int32:
		return VInt(int(v)), nil
	case int64:
		return VInt(int(v)), nil
	case uint:
		return VInt(int(v)), nil
	case uint8:
		return VInt(int(v)), nil
	case uint16:
		return VInt(int(v)), nil
	case uint32:
		return VInt(int(v)), nil
	case uint64:
		return VInt(int(v)), nil
	case float32:
		return VInt(int(v)), nil
	case float64:
		return VInt(int(v)), nil
	case string:
		return VStr(v), nil
	case bool:
		return VBool(v), nil
	default:
		return Value{}, fmt.Errorf("unsupported literal type %T", x)
	}
}
