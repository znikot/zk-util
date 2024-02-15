package misc

import (
	"errors"
	"math/rand"
	"reflect"
	"sort"
)

var errNotSlice = errors.New("not a slice")

// UnionSlice 求并集
func UnionSlice(slice1, slice2 interface{}) (interface{}, error) {
	// result := make([]interface{}, 0)
	var result reflect.Value
	if reflect.TypeOf(slice1).Kind() != reflect.Slice || reflect.TypeOf(slice2).Kind() != reflect.Slice {
		return nil, errNotSlice
	}

	vals1 := reflect.ValueOf(slice1)
	vals2 := reflect.ValueOf(slice2)

	if vals1.Len() == 0 && vals2.Len() == 0 {
		return vals1, nil
	}
	var ty reflect.Type
	if vals1.Len() > 0 {
		ty = vals1.Index(0).Type()
		// ty = reflect.TypeOf(vals1)
	} else if vals2.Len() > 0 {
		// ty = reflect.TypeOf(vals2)
		ty = vals2.Index(0).Type()
	}

	result = reflect.MakeSlice(reflect.SliceOf(ty), vals1.Len()+vals2.Len(), vals1.Len()+vals2.Len())

	m := make(map[interface{}]int)
	for i := 0; i < vals1.Len(); i++ {
		v := vals1.Index(i)
		m[v.Interface()]++
		//result = append(result, v)
		result.Index(i).Set(v)
	}

	idx := vals1.Len()
	for i := 0; i < vals2.Len(); i++ {
		v := vals2.Index(i)
		times := m[v.Interface()]
		if times == 0 {
			result.Index(idx).Set(v)
			idx++
		}
	}
	return result.Slice(0, idx).Interface(), nil
}

// IntersectSlice 求交集
func IntersectSlice(slice1, slice2 interface{}) (interface{}, error) {
	var result reflect.Value
	if reflect.TypeOf(slice1).Kind() != reflect.Slice || reflect.TypeOf(slice2).Kind() != reflect.Slice {
		return nil, errNotSlice
	}

	vals1 := reflect.ValueOf(slice1)
	vals2 := reflect.ValueOf(slice2)

	if vals1.Len() == 0 && vals2.Len() == 0 {
		return vals1, nil
	}
	var ty reflect.Type
	if vals1.Len() > 0 {
		ty = vals1.Index(0).Type()
	} else if vals2.Len() > 0 {
		ty = vals2.Index(0).Type()
	}

	result = reflect.MakeSlice(reflect.SliceOf(ty), vals1.Len()+vals2.Len(), vals1.Len()+vals2.Len())
	m := make(map[interface{}]int)
	for i := 0; i < vals1.Len(); i++ {
		v := vals1.Index(i)
		m[v.Interface()]++
	}

	idx := 0
	duplicate := make(map[interface{}]int)
	for i := 0; i < vals2.Len(); i++ {
		v := vals2.Index(i)
		vi := v.Interface()
		times := m[vi]
		// fmt.Printf("%v 出现次数 %d\n", v.Interface(), times)
		if times > 0 && duplicate[vi] == 0 {
			duplicate[vi]++
			result.Index(idx).Set(v)
			idx++
		}
	}

	return result.Slice(0, idx).Interface(), nil
}

// DifferenceSlice 求差值
func DifferenceSlice(slice1, slice2 interface{}) (interface{}, error) {
	if reflect.TypeOf(slice1).Kind() != reflect.Slice || reflect.TypeOf(slice2).Kind() != reflect.Slice {
		return nil, errNotSlice
	}

	itr, _ := IntersectSlice(slice1, slice2)
	m := map[interface{}]int{}

	visitSlice(itr, func(v reflect.Value) {
		m[v.Interface()]++
	})

	vals1 := reflect.ValueOf(slice1)
	vals2 := reflect.ValueOf(slice2)

	var ty reflect.Type
	if vals1.Len() > 0 {
		ty = vals1.Index(0).Type()
	} else if vals2.Len() > 0 {
		ty = vals2.Index(0).Type()
	}
	result := reflect.MakeSlice(reflect.SliceOf(ty), vals1.Len()+vals2.Len(), vals1.Len()+vals2.Len())

	idx := 0

	dup := make(map[interface{}]int)

	visitor := func(v reflect.Value) {
		if m[v.Interface()] == 0 && dup[v.Interface()] == 0 {
			result.Index(idx).Set(v)
			dup[v.Interface()]++
			idx++
		}
	}

	visitSlice(slice1, visitor)
	visitSlice(slice2, visitor)

	return result.Slice(0, idx).Interface(), nil
}

// SubSlice 从 s1 中 删除与 s2 中相同的内容
func SubSlice(s1, s2 interface{}) (interface{}, error) {
	if reflect.TypeOf(s1).Kind() != reflect.Slice {
		return nil, errNotSlice
	}
	vals := reflect.ValueOf(s1)
	if vals.Len() == 0 {
		return s1, nil
	}

	result := reflect.MakeSlice(reflect.SliceOf(vals.Index(0).Type()), vals.Len(), vals.Len())
	m := make(map[interface{}]int)
	err := visitSlice(s2, func(v reflect.Value) {
		m[v.Interface()]++
	})
	if err != nil {
		return nil, err
	}
	idx := 0
	dup := make(map[interface{}]int)
	visitSlice(s1, func(v reflect.Value) {
		if m[v.Interface()] == 0 && dup[v.Interface()] == 0 {
			result.Index(idx).Set(v)
			idx++
			dup[v.Interface()]++
		}
	})

	return result.Slice(0, idx).Interface(), nil
}

// visit slice with visitor
func visitSlice(slice interface{}, visitor func(v reflect.Value)) error {
	if reflect.TypeOf(slice).Kind() != reflect.Slice {
		return errNotSlice
	}

	if visitor == nil {
		return nil
	}

	vals := reflect.ValueOf(slice)

	for i := 0; i < vals.Len(); i++ {
		visitor(vals.Index(i))
	}

	return nil
}

// shuffle slice
func ShuffleSlice(slice any) {
	if reflect.TypeOf(slice).Kind() == reflect.Pointer {
		slice = reflect.Indirect(reflect.ValueOf(slice)).Interface()
	}
	if reflect.TypeOf(slice).Kind() != reflect.Slice {
		return
	}

	vals := reflect.ValueOf(slice)
	for i := 0; i < vals.Len(); i++ {
		j := rand.Intn(i + 1)
		tmp := vals.Index(i).Interface()
		vals.Index(i).Set(reflect.ValueOf(vals.Index(j).Interface()))
		vals.Index(j).Set(reflect.ValueOf(tmp))
	}
}

// check if x exists in string slice. slice a need sorted
//
// if sorted is false, means a is not sorted outside, will sort it
func ExistsString(x string, sorted bool, a ...string) bool {
	if !sorted {
		sort.Strings(a)
	}
	i := sort.SearchStrings(a, x)
	if i < len(a) && a[i] == x {
		return true
	}
	return false
}

// check if x exists in int slice. slice a need sorted
//
// if sorted is false, means a is not sorted outside, will sort it
func ExistsInt(x int, sorted bool, a ...int) bool {
	if !sorted {
		sort.Ints(a)
	}
	i := sort.SearchInts(a, x)
	if i < len(a) && a[i] == x {
		return true
	}
	return false
}

// check if x exists in int64 slice. slice a need sorted
//
// if sorted is false, means a is not sorted outside, will sort it
func ExistsInt64(x int64, sorted bool, a ...int64) bool {
	if !sorted {
		sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
	}

	i := sort.Search(len(a), func(i int) bool { return a[i] >= x })
	if i < len(a) && a[i] == x {
		return true
	}
	return false
}
