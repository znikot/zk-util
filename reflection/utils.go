package reflection

import (
	"fmt"
	"reflect"
	"strings"
)

// 将 val 设置到 target 的指定 path 里
// 如果 target 是 map，则会增加或者覆盖
// 只能处理指针类型，否则没意义
func SetElemByPath(target any, path string, val any) (err error) {
	targetType := reflect.TypeOf(target)
	if targetType.Kind() != reflect.Ptr {
		return fmt.Errorf("target must be a pointer")
	}

	// log.Debugf("runtime", "set %s to %s", path, target)

	// 获取指针所指的类型
	targetType = targetType.Elem()

	targetValue := reflect.ValueOf(target).Elem() // 目标值
	var currElemKey string                        // 当前层级元素名
	var isLastElem bool                           // 是否最后一个元素，只处理最后一个元素
	index := strings.Index(path, ".")
	if index != -1 {
		currElemKey = path[:index]
		isLastElem = false
	} else {
		currElemKey = path
		isLastElem = true
	}

	// 根据不同的类型进行设置
	if targetType.Kind() == reflect.Map {
		// map 的 key 必须为 string
		if targetType.Key() != String {
			return fmt.Errorf("map key must be string")
		}
		// 如果是最后一个元素，那么就设置值
		if isLastElem {
			tv, err := CastAny(val, targetType.Elem())
			if err != nil {
				return err
			}
			targetValue.SetMapIndex(reflect.ValueOf(currElemKey), reflect.ValueOf(tv))
		} else {
			// 中间元素，如果中间元素是 nil，那么需要创建一个新的
			middleValue := targetValue.MapIndex(reflect.ValueOf(currElemKey))
			if !middleValue.IsValid() || middleValue.IsZero() {
				// 创建一个新的中间元素
				if targetType.Elem().Kind() == reflect.Interface {
					// 继续创建一个 map
					// middleValue = reflect.New(reflect.TypeOf(map[string]interface{}{})).Elem()
					middleValue = reflect.ValueOf(&map[string]interface{}{}).Elem()
				} else if targetType.Elem().Kind() == reflect.Struct {
					middleValue = reflect.New(targetType.Elem())
				} else {
					return fmt.Errorf("unsupport map element type: %s", targetType.Elem().String())
				}
				targetValue.SetMapIndex(reflect.ValueOf(currElemKey), middleValue)
			} else {
				tmp := middleValue.Interface().(map[string]interface{})
				middleValue = reflect.ValueOf(&tmp)
			}
			// 继续往下
			if middleValue.CanAddr() {
				SetElemByPath(middleValue.Addr().Interface(), path[index+1:], val)
			} else {
				SetElemByPath(middleValue.Elem().Addr().Interface(), path[index+1:], val)
				// log.Warnf("runtime", "%s can't addr, path: %s", middleValue.Type().String(), path)
			}
		}
	} else if targetType.Kind() == reflect.Struct {
		return fmt.Errorf("not implement yet")
	} else {
		// 不支持其它类型的值设置
		return fmt.Errorf("target must be a pointer to a struct or map")
	}

	return nil
}

// get element from map、struct、slice through path like A.B.C
func GetElemByPath(src any, path string) (any, error) {
	// if path == "" then return src
	if len(path) == 0 {
		return src, nil
	}

	// get reflect infomation
	srcVal := reflect.ValueOf(src)
	if srcVal.Type().Kind() == reflect.Ptr {
		srcVal = srcVal.Elem()
	}

	// prim type, return error
	if IsPrimType(srcVal.Type()) {
		return nil, fmt.Errorf("can't get element from type %s through path %s", srcVal.Type().Name(), path)
	}

	// get first path element, till reach end of path
	index := strings.Index(path, ".")
	var currElemKey string
	if index != -1 {
		currElemKey = path[:index]
	} else {
		currElemKey = path
	}

	var currElem any

	// check type
	switch srcVal.Type().Kind() {
	case reflect.Map:
		// get current element from map
		currElem = srcVal.Interface().(map[string]any)[currElemKey]
	case reflect.Struct:
		refVal := reflect.ValueOf(src)
		refVal = refVal.FieldByName(currElemKey)
		if refVal.IsZero() {
			// not found, return nil
			currElem = nil
		} else {
			currElem = refVal.Interface()
		}
	case reflect.Slice:
	}

	// if path has more, get more
	if index != -1 {
		return GetElemByPath(currElem, path[index+1:])
	}

	return currElem, nil
}
