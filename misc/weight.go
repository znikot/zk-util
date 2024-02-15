package misc

import (
	"reflect"
	"sync/atomic"
)

// WeightObject 可计算权重的对象
type WeightObject interface {
	Weight() int
}

// Weight 权重信息
type Weight struct {
	slice     reflect.Value
	lastIndex int32 //表示上一个索引
	cw        int   //表示当前调度的权值
	gcd       int   //当前所有权重的最大公约数 比如 2，4，8 的最大公约数为：2
}

// NewWeight 新建权重信息
func NewWeight(objs interface{}) *Weight {
	t := reflect.TypeOf(objs)

	if t.Kind() != reflect.Slice {
		return nil
	}

	w := &Weight{
		slice:     reflect.ValueOf(objs),
		lastIndex: -1,
		cw:        0,
	}
	w.gcd = w.calcAllGcd(reflect.ValueOf(objs))

	return w
}

// 计算所有权重的最大公约数
func (w *Weight) calcAllGcd(vals reflect.Value) int {
	if vals.Len() == 0 {
		return 0
	}
	g := vals.Index(0).Interface().(WeightObject).Weight()
	for i := 1; i < vals.Len(); i++ {
		// 使用两个数字比较的最大公约数函数进行计算, 得出当前两个数字的最大公约数
		// 循环开始后, 依次将当前最大公约数与后面的数字一一进行运算, 求最大公约数
		g = w.calcGcd(g, vals.Index(i).Interface().(WeightObject).Weight())
	}

	return g
}

// 计算最大公约数
func (w *Weight) calcGcd(a, b int) int {
	if b == 0 {
		return a
	}
	return w.calcGcd(b, a%b)
}

// getMaxWeight 获取最大权重
func (w *Weight) getMaxWeight(vals reflect.Value) int {
	max := 0
	for i := 0; i < vals.Len(); i++ {
		s := vals.Index(i).Interface().(WeightObject)
		if s.Weight() >= max {
			max = s.Weight()
		}
	}

	return max
}

// Index 按照权重信息获取索引
func (w *Weight) NextIndex() int {
	for {
		lastIndex := int(atomic.AddInt32(&w.lastIndex, 1))
		lastIndex = lastIndex % w.slice.Len()
		if lastIndex == 0 {
			w.cw = w.cw - w.gcd
			if w.cw <= 0 {
				w.cw = w.getMaxWeight(w.slice)
				if w.cw == 0 {
					return 0
				}
			}
		}
		if w.slice.Index(lastIndex).Interface().(WeightObject).Weight() >= w.cw {
			return lastIndex
		}
	}
}
