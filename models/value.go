package models

import (
	"fmt"
)

type Kind int

const (
	KindNil Kind = iota
	KindInt
	KindString
	KindBool
	KindField
)

type Value struct {
	K Kind
	I int
	S string
	B bool
}

func VInt(i int) Value    { return Value{K: KindInt, I: i} }
func VStr(s string) Value { return Value{K: KindString, S: s} }
func VBool(b bool) Value  { return Value{K: KindBool, B: b} }
func VNil() Value         { return Value{K: KindNil} }

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
