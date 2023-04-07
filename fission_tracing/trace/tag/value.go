package tag

import (
	"math"
	"strconv"
)

const (
	NONE int = iota
	BOOL
	INT64
	FLOAT64
	STRING
)

type Value struct {
	valueType   int
	valueString string
	valueNum    uint64
}

func GetBoolValue(v bool) Value {
	boolValue := 0
	if v == true {
		boolValue = 1
	}
	return Value{
		valueType: BOOL,
		valueNum:  uint64(boolValue),
	}
}

func GetIntValue(v int64) Value {
	return Value{
		valueType: INT64,
		valueNum:  uint64(v),
	}
}

func GetFloatValue(v float64) Value {
	return Value{
		valueType: FLOAT64,
		valueNum:  math.Float64bits(v),
	}
}

func GetStringValue(v string) Value {
	return Value{
		valueType:   STRING,
		valueString: v,
	}
}

func (v Value) Type() int {
	return v.valueType
}

func (v Value) toBool() bool {
	return v.valueNum != 0
}

func (v Value) String() string {
	switch v.Type() {
	case BOOL:
		return strconv.FormatBool(v.toBool())
	case INT64:
		return strconv.FormatInt(int64(v.valueNum), 10)
	case FLOAT64:
		return strconv.FormatFloat(math.Float64frombits(v.valueNum), 'E', -1, 64)
	case STRING:
		return v.valueString
	default:
		return "unknown"
	}
}
