package reflection

import (
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/znikot/zk-util/misc"
)

type TagInfo struct {
	Name    string
	Value   string
	Values  []string
	Details map[string]string
}

func (me *TagInfo) HasDetail(name string) (ex bool) {
	ex = false
	if me.Details == nil {
		return
	}
	_, ex = me.Details[name]
	return
}

type FieldInfo struct {
	Name     string
	Index    []int
	Type     reflect.Type
	TypeName string
	Extend   bool         // extends from struct
	FromType reflect.Type // struct type extends from for this field

	// Tags map[string]map[string]string
	Tags map[string]*TagInfo
}

// stores struct's info
type TypeInfo struct {
	Entity    interface{}
	Type      reflect.Type
	FullName  string
	ShortName string
	Fields    []FieldInfo
	// ExtendFields []FieldInfo // fields extends from other struct
}

var (
	_infoCache = make(map[reflect.Type]*TypeInfo)

	_lock = &sync.RWMutex{}

	// types
	String    = reflect.TypeOf("")
	Int       = reflect.TypeOf(0)
	Bool      = reflect.TypeOf(false)
	Float     = reflect.TypeOf(0.0)
	Int8      = reflect.TypeOf(int8(0))
	Int16     = reflect.TypeOf(int16(0))
	Int32     = reflect.TypeOf(int32(0))
	Int64     = reflect.TypeOf(int64(0))
	Uint      = reflect.TypeOf(uint(0))
	Uint8     = reflect.TypeOf(uint8(0))
	Uint16    = reflect.TypeOf(uint16(0))
	Uint32    = reflect.TypeOf(uint32(0))
	Uint64    = reflect.TypeOf(uint64(0))
	Float32   = reflect.TypeOf(float32(0))
	Float64   = reflect.TypeOf(float64(0))
	BoolPtr   = reflect.TypeOf((*bool)(nil))
	Ptr       = reflect.TypeOf((*interface{})(nil))
	PtrPtr    = reflect.TypeOf((***interface{})(nil))
	PtrMap    = reflect.TypeOf((*map[string]interface{})(nil))
	Map       = reflect.TypeOf(map[string]interface{}{})
	DateTime  = reflect.TypeOf(misc.DateTime{})
	Date      = reflect.TypeOf(misc.Date{})
	Time      = reflect.TypeOf(misc.Time{})
	Timestamp = reflect.TypeOf(misc.Timestamp{})
	Error     = reflect.TypeOf((*error)(nil)).Elem()
)

func getType(entity interface{}) reflect.Type {
	_type := reflect.TypeOf(entity)
	if _type == nil {
		return nil
	}

	if _type.Kind() == reflect.Ptr || _type.Kind() == reflect.Slice {
		_type = _type.Elem()
	}

	return _type
}

// GetTypeInfo returns the type info for the given type
func GetTypeInfo(entity interface{}) *TypeInfo {
	_type := getType(entity)
	if _type == nil {
		return nil
	}
	// log.Infof("", "get type info of %s", _type.Name())
	i, ok := _infoCache[_type]
	if ok {
		return i
	}
	_lock.Lock()
	defer _lock.Unlock()

	i = resolveType(_type)
	i.Entity = entity

	_infoCache[_type] = i

	return i
}

func resolveType(_type reflect.Type) *TypeInfo {
	typeInfo := &TypeInfo{
		Type:      _type,
		FullName:  _type.PkgPath() + "." + _type.Name(),
		ShortName: _type.String(),
		Fields:    resolveStructFields(_type, false),
		// ExtendFields: resolveStructFields(_type, true),
	}

	// fc := _type.NumField()
	// for i := 0; i < fc; i++ {
	// 	field := _type.Field(i)

	// 	// if field was struct, so get it's fields
	// 	if field.Type.Kind() == reflect.Struct {

	// 	} else {
	// 		fieldInfo := FieldInfo{
	// 			Name: field.Name,
	// 			Type: field.Type,
	// 			Tags: resolveTags(field),
	// 		}
	// 		if field.Type.PkgPath() != "" {
	// 			fieldInfo.TypeName = field.Type.PkgPath() + "." + field.Type.Name()
	// 		} else {
	// 			fieldInfo.TypeName = field.Type.Name()
	// 		}

	// 		typeInfo.Fields = append(typeInfo.Fields, fieldInfo)
	// 	}
	// }

	return typeInfo
}

func resolveStructFields(_type reflect.Type, extend bool) []FieldInfo {
	fields := make([]FieldInfo, 0)
	if _type.Kind() != reflect.Struct {
		return fields
	}
	fc := _type.NumField()
	for i := 0; i < fc; i++ {
		field := _type.Field(i)
		// private field will ignore
		if !field.IsExported() {
			continue
		}
		// if field was struct, so get it's fields, but excludes all datetime struct
		if field.Type.Kind() == reflect.Struct && field.Type != Time && field.Type != DateTime && field.Type != Date && field.Type != Timestamp {
			fields = append(fields, resolveStructFields(field.Type, true)...)
		} else {
			fieldInfo := FieldInfo{
				Name:     field.Name,
				Index:    field.Index,
				Type:     field.Type,
				Tags:     resolveTags(field),
				Extend:   extend,
				FromType: _type,
			}
			if field.Type.PkgPath() != "" {
				fieldInfo.TypeName = field.Type.PkgPath() + "." + field.Type.Name()
			} else {
				fieldInfo.TypeName = field.Type.Name()
			}
			fields = append(fields, fieldInfo)
		}
	}
	return fields
}

func resolveTags(field reflect.StructField) map[string]*TagInfo {
	tags := make(map[string]*TagInfo)
	fieldTag := field.Tag
	for fieldTag != "" {
		// Skip leading space.
		i := 0
		for i < len(fieldTag) && fieldTag[i] == ' ' {
			i++
		}
		fieldTag = fieldTag[i:]
		if fieldTag == "" {
			break
		}

		i = 0
		for i < len(fieldTag) && fieldTag[i] > ' ' && fieldTag[i] != ':' && fieldTag[i] != '"' && fieldTag[i] != 0x7f {
			i++
		}
		if i == 0 {
			break
		}
		name := string(fieldTag[:i])
		if i+1 >= len(fieldTag) || fieldTag[i] != ':' || fieldTag[i+1] != '"' {
			tags[name] = &TagInfo{Name: name, Values: make([]string, 0), Value: "", Details: make(map[string]string)}
			break
		}
		// name := string(fieldTag[:i])
		fieldTag = fieldTag[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(fieldTag) && fieldTag[i] != '"' {
			if fieldTag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(fieldTag) {
			break
		}
		qvalue := string(fieldTag[:i+1])
		fieldTag = fieldTag[i+1:]

		value, err := strconv.Unquote(qvalue)
		if err != nil {
			//
		}

		// tags[name] = strings.Split(value, " ")
		ti := &TagInfo{
			Name:    name,
			Values:  strings.Split(value, " "),
			Value:   value,
			Details: make(map[string]string),
		}
		for _, s := range ti.Values {
			if len(s) == 0 {
				continue
			}
			idx := strings.Index(s, "=")
			if idx >= 0 {
				ti.Details[s[0:idx]] = s[idx+1:]
			} else {
				ti.Details[s] = ""
			}
		}

		tags[name] = ti
	}
	return tags
}

func (t *TypeInfo) GetFieldByType(_type reflect.Type) *FieldInfo {
	if _type.Kind() == reflect.Ptr {
		_type = _type.Elem()
	}

	for _, field := range t.Fields {
		if field.Type == _type {
			return &field
		}
	}

	return nil
}

// make new ptr of this type
func (t *TypeInfo) New() interface{} {
	return reflect.New(t.Type).Interface()
}

// make new slice of this type
func (t *TypeInfo) NewSlice(len, cap int) interface{} {
	// slice := reflect.MakeSlice(reflect.SliceOf(t.Type), len, cap)
	// log.Debugf(app.LogRuntime, "slice type: %s, len: %v", slice.Type().Name(), len)
	// return slice.Interface()
	return reflect.New(reflect.SliceOf(t.Type)).Interface()
}

// make new map of this type
func (t *TypeInfo) NewMap() interface{} {
	return reflect.New(reflect.MapOf(t.Type.Key(), t.Type.Elem())).Interface()
}

// make new slice of this type with pointer
func (t *TypeInfo) NewSlicePointer(len, cap int) interface{} {
	return reflect.New(reflect.SliceOf(reflect.PointerTo(t.Type))).Interface()
}

func (t *TypeInfo) GetFieldValue(object interface{}, fieldName string) interface{} {
	vals := reflect.ValueOf(object)
	if vals.Kind() == reflect.Slice {
		if vals.Len() == 0 {
			return nil
		}
		fv := vals.Index(0)
		if fv.Kind() == reflect.Ptr {
			fv = fv.Elem()
		}

		// make a slice of values
		slice := reflect.MakeSlice(reflect.SliceOf(fv.FieldByName(fieldName).Type()), vals.Len(), vals.Len())
		for i := 0; i < vals.Len(); i++ {
			v := vals.Index(i)
			if v.Kind() == reflect.Ptr {
				v = v.Elem()
			}
			slice.Index(i).Set(v.FieldByName(fieldName))
		}
		return slice.Interface()
	} else if vals.Kind() == reflect.Ptr {
		vals = vals.Elem()
	}

	v := vals.FieldByName(fieldName)

	return v.Interface()
}

func (t *TypeInfo) SetFieldValueByName(object interface{}, fieldName string, fieldValue interface{}) {
	vals := reflect.ValueOf(object)
	if vals.Kind() == reflect.Slice {
		for i := 0; i < vals.Len(); i++ {
			v := vals.Index(i)
			if v.Kind() == reflect.Ptr {
				v = v.Elem()
			}
			fv := v.FieldByName(fieldName)
			setFieldValue(fv, fieldValue)
		}

		return
	} else if vals.Kind() == reflect.Ptr {
		vals = vals.Elem()
	}

	v := vals.FieldByName(fieldName)
	setFieldValue(v, fieldValue)
}

func (t *TypeInfo) SetFieldValueByType(object interface{}, fieldType reflect.Type, fieldValue interface{}) {
	fieldInfo := t.GetFieldByType(fieldType)
	if fieldInfo == nil {
		return
	}
	t.SetFieldValueByName(object, fieldInfo.Name, fieldValue)
}

func (t *TypeInfo) GetField(name string, ignoreCase bool) *FieldInfo {
	// base field
	for _, field := range t.Fields {
		if ignoreCase && strings.EqualFold(field.Name, name) {
			return &field
		} else if field.Name == name {
			return &field
		}
	}
	// extends field
	// for _, field := range t.ExtendFields {
	// 	if ignoreCase && strings.EqualFold(field.Name, name) {
	// 		return &field
	// 	} else if field.Name == name {
	// 		return &field
	// 	}
	// }
	return nil
}

// 获取 tag 值
//
//	type User struct {
//		Name string `validate:"required=true max=10"`
//		RoleName string `repo:"exclude=insert,update join=role.Name"`
//	}
//
// validate tag 中，validate 是 group，required、max 是 name
// 如果 name 传入空字符串 "" 则返回 "required=true max=10"
//
// repo tag 中，repo 是 group，exclude、join 是 name。insert,update、role.Name
//
//	GetTag("repo", "") // 结果为 "exclude=insert,update join=role.Name"
//	GetTag("repo", "exclude") // 结果为 "insert,update"
//	GetTag("repo", "join") // 结果为 "role.Name"
func (f *FieldInfo) GetTag(group, name string) (string, bool) {
	t, ok := f.Tags[group]
	if !ok {
		return "", false
	}
	if name == "" {
		return t.Value, true
	}
	v, ok := t.Details[name]
	return v, ok
}

// 获取字段值
func (f *FieldInfo) ValueOf(val reflect.Value) any {
	if val.Kind() == reflect.Slice {
		if val.Len() == 0 {
			return nil
		}
		fv := val.Index(0)
		if fv.Kind() == reflect.Ptr {
			fv = fv.Elem()
		}

		// make a slice of values
		slice := reflect.MakeSlice(reflect.SliceOf(fv.FieldByIndex(f.Index).Type()), val.Len(), val.Len())
		for i := 0; i < val.Len(); i++ {
			v := val.Index(i)
			if v.Kind() == reflect.Ptr {
				v = v.Elem()
			}
			slice.Index(i).Set(v.FieldByIndex(f.Index))
		}
		return slice.Interface()
	} else if val.Kind() == reflect.Ptr {
		if val.IsZero() {
			return nil
		}
		val = val.Elem()
	}

	return val.FieldByIndex(f.Index).Interface()
}

// 设置字段值
func (f *FieldInfo) SetValue(obj reflect.Value, val any) {
	if obj.Kind() == reflect.Slice {
		for i := 0; i < obj.Len(); i++ {
			v := obj.Index(i)
			if v.Kind() == reflect.Ptr {
				v = v.Elem()
			}
			fv := v.FieldByIndex(f.Index)
			setFieldValue(fv, val)
		}

		return
	} else if obj.Kind() == reflect.Ptr {
		obj = obj.Elem()
	}

	v := obj.FieldByIndex(f.Index)
	setFieldValue(v, val)
}

func setFieldValue(refVal reflect.Value, val interface{}) {
	val, err := CastAny(val, refVal.Type())
	if err != nil {
		panic(err)
	}
	refVal.Set(reflect.ValueOf(val))
}

func IsPrimType(_type reflect.Type) bool {
	if _type.Kind() == reflect.Pointer {
		_type = _type.Elem()
	}
	if _type.Kind() == reflect.Slice {
		_type = _type.Elem()
		if _type.Kind() == reflect.Pointer {
			_type = _type.Elem()
		}
	}
	return _type.Kind() == reflect.Bool ||
		_type.Kind() == reflect.Int ||
		_type.Kind() == reflect.Int8 ||
		_type.Kind() == reflect.Int16 ||
		_type.Kind() == reflect.Int32 ||
		_type.Kind() == reflect.Int64 ||
		_type.Kind() == reflect.Uint ||
		_type.Kind() == reflect.Uint8 ||
		_type.Kind() == reflect.Uint16 ||
		_type.Kind() == reflect.Uint32 ||
		_type.Kind() == reflect.Uint64 ||
		_type.Kind() == reflect.Float32 ||
		_type.Kind() == reflect.Float64 ||
		_type.Kind() == reflect.String
}
