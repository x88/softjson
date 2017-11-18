package softjson

import (
	"encoding/json"
	"reflect"
	"strings"
	"strconv"
	"math"
)

type T map[string]interface{}

func UnmarshalSoft(data []byte, v interface{}) error {
	var t map[string]interface{}
	err := json.Unmarshal(data, &t)
	if err != nil {
		return err
	}
	doStruct(t, reflect.ValueOf(v).Elem(), reflect.TypeOf(v).Elem())
	return nil
}

func doStruct(m map[string]interface{}, tValue reflect.Value, tType reflect.Type) {
	for i := 0; i < tType.NumField(); i++ {
		fieldName := getFieldName(tType.Field(i))
		fieldVal := tValue.Field(i)
		if val := m[fieldName]; val != nil {
			if fieldVal.IsValid() && fieldVal.CanSet() {
				if (fieldVal.Kind() == reflect.TypeOf(val).Kind()) && (fieldVal.Kind() != reflect.Slice) && (fieldVal.Kind() != reflect.Array) {
					fieldVal.Set(reflect.ValueOf(val))
				} else {
					processType(reflect.ValueOf(val), fieldVal)
				}
			}
		}

	}

}

func processType(srcValue, fieldVal reflect.Value) {
	switch fieldVal.Kind() {
	case reflect.String:
		fieldVal.SetString(doString(srcValue))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		fieldVal.Set(reflect.ValueOf(doInt(srcValue)))
	case reflect.Int64:
		fieldVal.SetInt(doInt64(srcValue))
	case reflect.Float32:
		fieldVal.Set(reflect.ValueOf(doFloat32(srcValue)))
	case reflect.Float64:
		fieldVal.SetFloat(doFloat64(srcValue))
	case reflect.Bool:
		fieldVal.SetBool(doBool(srcValue))
	case reflect.Struct:
		doStruct(srcValue.Interface().(map[string]interface{}), fieldVal, fieldVal.Type())
	case reflect.Slice, reflect.Array:
		doArray(srcValue, fieldVal, fieldVal.Type().Elem())
	case reflect.Ptr:
		panic(`Pointer unsupported`)
	case reflect.Invalid:
		panic(`invalid`)
	}
}

func doArray(srcValue, tValue reflect.Value, tType reflect.Type) {
	for i := 0; i < srcValue.Len(); i++ {
		tmpVal := reflect.New(tValue.Type().Elem())
		processType(srcValue.Index(i), tmpVal.Elem())
		tValue.Set(reflect.Append(tValue, tmpVal.Elem()))
	}
}

func doInt(src reflect.Value) int {
	switch src.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		res := src.Int()
		if !overflowInt32(res) {
			return int(res)
		}
	case reflect.Bool:
		if src.Bool() {
			return 1
		} else {
			return 0
		}
	case reflect.Float32, reflect.Float64:
		i, f := math.Modf(src.Float())
		if (f == 0) && !overflowInt32(int64(i)) {
			return int(i)
		}
	case reflect.String:
		res, err := strconv.ParseInt(src.String(), 10, 32)
		if (err == nil) && !overflowInt32(res) {
			return int(res)
		}
	}
	panic(`Incorrect integer type`)
}

func doInt64(src reflect.Value) int64 {
	switch src.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		res := src.Int()
		if !src.OverflowInt(res) {
			return res
		}
	case reflect.Bool:
		if src.Bool() {
			return 1
		} else {
			return 0
		}
	case reflect.Float32, reflect.Float64:
		i, f := math.Modf(src.Float())
		if (f == 0) && !src.OverflowInt(int64(i)) {
			return int64(i)
		}
	case reflect.String:
		res, err := strconv.ParseInt(src.String(), 10, 64)
		if err == nil && !overflowInt64(res) {
			return res
		}
	}
	panic(`Incorrect integer type`)
}

func doFloat32(src reflect.Value) float32 {
	switch src.Kind() {
	case reflect.Float32, reflect.Float64:
		res := src.Float()
		if !overflowFloat32(res) {
			return float32(res)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		res := src.Int()
		if !overflowFloat32(float64(res)) {
			return float32(res)
		}
	case reflect.Bool:
		if src.Bool() {
			return 1
		} else {
			return 0
		}
	case reflect.String:
		res, err := strconv.ParseFloat(src.String(), 64)
		if (err == nil) && !overflowFloat32(res) {
			return float32(res)
		}
	}
	panic(`Inccorrect float32 type`)
}

func doFloat64(src reflect.Value) float64 {
	switch src.Kind() {
	case reflect.Float32, reflect.Float64:
		res := src.Float()
		if !src.OverflowFloat(res) {
			return float64(res)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		res := src.Int()
		if !src.OverflowFloat(float64(res)) {
			return float64(res)
		}
	case reflect.Bool:
		if src.Bool() {
			return 1
		} else {
			return 0
		}
	case reflect.String:
		res, err := strconv.ParseFloat(src.String(), 64)
		if err == nil {
			return res
		}
	}
	panic(`Inccorrect float64 type`)
}

func doString(src reflect.Value) string {
	switch src.Kind() {
	case reflect.String:
		return src.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(src.Int(), 10)
	case reflect.Bool:
		return strconv.FormatBool(src.Bool())
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(src.Float(), 'f', -1, 64)
	}
	panic(`Inccorrect string type`)
}

func doBool(src reflect.Value) bool {
	switch src.Kind() {
	case reflect.Bool:
		return src.Bool()
	case reflect.String:
		res := src.String()
		if res == `true` {
			return true
		} else if res == `false` {
			return false
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		res := src.Int()
		if res == 1 {
			return true
		} else if res == 0 {
			return false
		}
	case reflect.Float32, reflect.Float64:
		res := src.Float()
		if res == 1 {
			return true
		} else if res == 0 {
			return false
		}
	}
	panic(`Inccorrect bool type`)
}

func getFieldName(field reflect.StructField) string {
	tagData := field.Tag.Get(`json`)
	if tagData == `` {
		return field.Name
	} else {
		if endVal := strings.Index(tagData, `,`); endVal > -1 {
			tagData = tagData[:endVal]
		}
		return tagData
	}
}

func overflowInt32(x int64) bool {
	trunc := (x << 18) >> (18)
	return x != trunc
}

func overflowInt64(x int64) bool {
	bitSize := uint(6 * 8)
	trunc := (x << (64 - bitSize)) >> (64 - bitSize)
	return x != trunc
}

func overflowFloat32(x float64) bool {
	if x < 0 {
		x = -x
	}
	return math.MaxFloat32 < x && x <= math.MaxFloat64
}
