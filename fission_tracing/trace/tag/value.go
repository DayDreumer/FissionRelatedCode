package tag

import (
	"encoding/json"
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

const TYPE_LIST = "NONEBOOLINT64FLOAT64STRING"

var TYPE_ARRAY = [...]int{0, 4, 8, 13, 20, 26}

func ValueTypeToString(vtype int) string {
	if vtype < 0 || vtype > 5 {
		return "UNKNOWN"
	}
	return TYPE_LIST[TYPE_ARRAY[vtype]:TYPE_ARRAY[vtype+1]]
}

func StringToValueType(valueType string) int {
	switch valueType {
	case "BOOL":
		return 1
	case "INT64":
		return 2
	case "FLOAT64":
		return 3
	case "STRING":
		return 4
	default:
		return 0
	}
}

type NONE_TYPE struct{}

type Value struct {
	valueType   int
	valueString string
	valueNum    uint64
}

func GetNoneValue() Value {
	return Value{
		valueType: NONE,
	}
}

func GetBoolValue(v bool) Value {
	boolValue := 0
	if v {
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

func (v Value) toInt64() int64 {
	return int64(v.valueNum)
}

func (v Value) toFloat64() float64 {
	return math.Float64frombits(v.valueNum)
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

func (v Value) ToInterface() interface{} {
	switch v.valueType {
	case BOOL:
		return v.toBool()
	case INT64:
		return v.toInt64()
	case FLOAT64:
		return v.toFloat64()
	case STRING:
		return v.String()
	default:
		return NONE_TYPE{}
	}
}

func (v Value) MarshalJSON() ([]byte, error) {
	var jsonVal struct {
		Type  string
		Value interface{}
	}
	jsonVal.Type = ValueTypeToString(v.Type())
	jsonVal.Value = v.ToInterface()
	return json.Marshal(jsonVal)
}

func (v *Value) UnmarshalJSON(data []byte) error {
	var jsonVal struct {
		Type  string
		Value interface{}
	}
	err := json.Unmarshal(data, &jsonVal)
	if err != nil {
		return err
	}
	v.valueType = StringToValueType(jsonVal.Type)
	switch v.valueType {
	case BOOL:
		if jsonVal.Value == true {
			v.valueNum = 1
		} else {
			v.valueNum = 0
		}
	case INT64:
		v.valueNum = uint64(jsonVal.Value.(float64))
	case FLOAT64:
		val := jsonVal.Value.(float64)
		v.valueNum = math.Float64bits(val)
	case STRING:
		v.valueString = jsonVal.Value.(string)
	default:
	}
	return err
}
