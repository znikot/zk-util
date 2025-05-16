package reflection

import (
	"errors"
	"reflect"

	"github.com/iancoleman/strcase"
	"github.com/spf13/cast"
)

func init() {
	AddKindConverter(reflect.Int, reflect.String, any2string)
	AddKindConverter(reflect.Int16, reflect.String, any2string)
	AddKindConverter(reflect.Int32, reflect.String, any2string)
	AddKindConverter(reflect.Int64, reflect.String, any2string)
	AddKindConverter(reflect.Bool, reflect.String, any2string)
	AddKindConverter(reflect.Float32, reflect.String, any2string)
	AddKindConverter(reflect.Float64, reflect.String, any2string)

	AddKindConverter(reflect.String, reflect.Int, any2int)
	AddKindConverter(reflect.Float32, reflect.Int, any2int)
	AddKindConverter(reflect.Float64, reflect.Int, any2int)
	AddKindConverter(reflect.Int16, reflect.Int, any2int)
	AddKindConverter(reflect.Int32, reflect.Int, any2int)
	AddKindConverter(reflect.Int64, reflect.Int, any2int)

	AddKindConverter(reflect.String, reflect.Int64, any2int64)
	AddKindConverter(reflect.Float32, reflect.Int64, any2int64)
	AddKindConverter(reflect.Float64, reflect.Int64, any2int64)
	AddKindConverter(reflect.Int, reflect.Int64, any2int64)
	AddKindConverter(reflect.Int16, reflect.Int64, any2int64)
	AddKindConverter(reflect.Int32, reflect.Int64, any2int64)
	AddKindConverter(reflect.Int64, reflect.Int64, any2int64)

	AddKindConverter(reflect.String, reflect.Float32, any2float32)
	AddKindConverter(reflect.String, reflect.Float64, any2float64)

	AddKindConverter(reflect.Float64, reflect.Float32, any2float32)
	AddKindConverter(reflect.Float32, reflect.Float64, any2float64)

	AddKindConverter(reflect.Struct, reflect.Map, struct2map)

	AddKindConverter(reflect.Map, reflect.Struct, map2struct)

	AddKindConverter(reflect.Slice, reflect.Slice, slice2slice)

	AddKindConverter(reflect.Map, reflect.Map, map2map)

}

var (
	err_src_not_map    = errors.New("source parameter is not a map")
	err_target_not_map = errors.New("target parameter is not a map")
)

func any2string(i any, targetType reflect.Type) any {
	return cast.ToString(i)
}

func any2int(s any, targetType reflect.Type) any {
	return cast.ToInt(s)
}

func any2int64(s any, targetType reflect.Type) any {
	return cast.ToInt64(s)
}

func any2float32(s any, targetType reflect.Type) any {
	return cast.ToFloat32(s)
}
func any2float64(s any, targetType reflect.Type) any {
	return cast.ToFloat32(s)
}

func struct2struct(s any, targetType reflect.Type) any {
	// using json
	return s
}

// convert map to struct
func map2struct(src any, targetType reflect.Type) any {
	// assert src to map
	srcMap, ok := src.(map[string]any)
	if !ok {
		panic(err_src_not_map)
	}
	// new struct
	target := reflect.New(targetType)

	findMapValue := func(key string) (any, bool) {
		// match key
		if v, ok := srcMap[key]; ok {
			return v, true
		}
		// match camel case
		if v, ok := srcMap[strcase.ToLowerCamel(key)]; ok {
			return v, true
		}
		// match snake case
		if v, ok := srcMap[strcase.ToSnake(key)]; ok {
			return v, true
		}

		return nil, false
	}

	// convert all map fields to struct fields
	for i := 0; i < targetType.NumField(); i++ {
		field := targetType.Field(i)
		fieldVal := target.Elem().FieldByName(field.Name)
		if field.Type.Kind() == reflect.Struct {
			var tv any
			mapVal, found := findMapValue(field.Name)
			// this field is extends from struct
			if found {
				_, isMap := mapVal.(map[string]any)
				if isMap {
					tv = map2struct(mapVal, field.Type)
				} else {
					// tv = map2struct(src, field.Type)
					tv = CastAny(mapVal, field.Type)
				}
			} else {
				tv = map2struct(src, field.Type)
			}
			if tv != nil {
				fieldVal.Set(reflect.ValueOf(tv))
			}
		} else {
			// get value from map
			// mapVal, ok := srcMap[field.Name]
			mapVal, ok := findMapValue(field.Name)
			if ok {
				// map value found, convert and set to field value
				tv := CastAny(mapVal, fieldVal.Type())
				if tv != nil {
					fieldVal.Set(reflect.ValueOf(tv))
				}
			}
		}
	}

	return target.Elem().Interface()
}

func struct2map(src any, targetType reflect.Type) any {
	// log.Debugf("", "converting from %s to map", reflect.TypeOf(src).Name())
	t := reflect.TypeOf(src)
	v := reflect.ValueOf(src)

	var target = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		val := v.Field(i)
		if IsPrimType(val.Type()) {
			target[t.Field(i).Name] = val.Interface()
		} else if val.Type().Kind() == reflect.Struct {
			target[t.Field(i).Name] = struct2map(val.Interface(), val.Type())
		}
	}
	return target
}

func struct2pointer(src any, targetType reflect.Type) any {
	return nil
}

// convert between two slice
func slice2slice(src any, targetType reflect.Type) any {
	// assert src to slice
	srcSlice, ok := src.([]any)
	if !ok {
		panic(err_src_not_map)
	}
	// new slice
	target := reflect.MakeSlice(reflect.SliceOf(targetType.Elem()), len(srcSlice), len(srcSlice))

	for i := 0; i < len(srcSlice); i++ {
		//target.Index(i).Set(reflect.ValueOf(srcSlice[i]))
		srcVal := srcSlice[i]
		// convert it
		target.Index(i).Set(reflect.ValueOf(CastAny(srcVal, targetType.Elem())))
	}

	return target.Interface()
}

// convert between two map
func map2map(src any, targetType reflect.Type) any {
	_srcType := reflect.TypeOf(src)
	if _srcType.Kind() == reflect.Pointer {
		_srcType = _srcType.Elem()
	}
	if _srcType.Kind() != reflect.Map {
		panic(err_src_not_map)
	}
	if targetType.Kind() != reflect.Map {
		panic(err_target_not_map)
	}

	// get target key type and value type
	keyType := targetType.Key()
	valueType := targetType.Elem()

	srcVal := reflect.ValueOf(src)

	// new map
	target := reflect.MakeMap(targetType)

	// convert all map fields to struct fields
	for _, k := range srcVal.MapKeys() {
		v := srcVal.MapIndex(k)

		// convert key and value to target type
		tk := CastAny(k.Interface(), keyType)
		tv := CastAny(v.Interface(), valueType)

		// set key and value to target map
		target.SetMapIndex(reflect.ValueOf(tk), reflect.ValueOf(tv))
	}

	return target.Interface()
}
