package reflection

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type _kind_converter func(src any, _type reflect.Type) any

type _type_converter func(src any) any

var (
	_kind_converters = make(map[reflect.Kind]map[reflect.Kind]_kind_converter)
	_type_converters = make(map[reflect.Type]map[reflect.Type]_type_converter)
	// _converters = make(map[reflect.Kind]map[reflect.Kind]Converter[any, any])
	_locker = &sync.RWMutex{}

	err_src_kind_not_supported  = errors.New("convert from this kind not supported")
	err_dest_kind_not_supported = errors.New("convert to this kind not supported")
	err_converter_not_found     = errors.New("no converter found")
	err_ptr_needed              = errors.New("parameter target must be pointer")
)

func Cast[T any](src any) T {
	var t T
	_type := reflect.TypeOf(t)

	target := CastAny(src, _type)

	return target.(T)
}

func CastAny(src any, _type reflect.Type) any {
	if src == nil {
		return nil
	}
	_srcType := reflect.TypeOf(src)
	if _srcType.Kind() == reflect.Ptr {
		_srcType = _srcType.Elem()
	}

	if _srcType == _type || _type.Kind() == reflect.Interface /*|| path.Join(_srcType.PkgPath(), _srcType.Name()) == path.Join(_type.PkgPath(), _type.Name())*/ {
		return src
	}

	tts, ok := _type_converters[_srcType]
	if ok {
		_tconv, ok := tts[_type]
		if ok {
			return _tconv(src)
		}
	}

	cts, ok := _kind_converters[_srcType.Kind()]
	if !ok {
		panic(fmt.Errorf("comvert from %s to %s not supported", _srcType.Kind().String(), _type.Kind().String()))
	}
	_conv, ok := cts[_type.Kind()]
	if !ok {
		panic(fmt.Errorf("comvert from %s to %s not supported, src data: %v", _srcType.Kind().String(), _type.Kind().String(), src))
	}
	return _conv(src, _type)
}

func convertSupported(kind reflect.Kind) bool {
	exKinds := kind == reflect.Slice ||
		kind == reflect.Array ||
		kind == reflect.Map
	return !exKinds
}

func AddKindConverter(srcKind, distKind reflect.Kind, converter _kind_converter) {
	_locker.Lock()
	defer _locker.Unlock()

	c, ok := _kind_converters[srcKind]
	if !ok {
		c = make(map[reflect.Kind]_kind_converter)
		_kind_converters[srcKind] = c
	}
	c[distKind] = converter
}

func AddTypeConverter(srcType, distType reflect.Type, converter _type_converter) {
	_locker.Lock()
	defer _locker.Unlock()

	c, ok := _type_converters[srcType]
	if !ok {
		c = make(map[reflect.Type]_type_converter)
		_type_converters[srcType] = c
	}
	c[distType] = converter
}
