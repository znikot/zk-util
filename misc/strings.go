package misc

import (
	"math/rand"
	"strings"
	"unicode"
)

var randomChars = [][]rune{
	{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}, // 数字
	{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}, // 小写英文字
	{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}, // 大写英文字
	{'~', '`', '!', '@', '#', '$', '%', '^', '&', '*', '(', ')', '_', '+', '{', '}', '[', ']', ':', ';', '\''},                         // 特殊字符
}

type RandomString struct {
	length int
	// 0 - 数字 1 - 小写英文字 2 - 大写英文字 3 - 特殊字符
	typeIndex map[int]int
}

func NewRandomString(length int) *RandomString {
	return &RandomString{length: length, typeIndex: map[int]int{0: 0, 1: 1, 2: 2}}
}

func (me *RandomString) withType(yes bool, index int) *RandomString {
	if yes {
		me.typeIndex[index] = index
	} else {
		delete(me.typeIndex, index)
	}
	return me
}

func (me *RandomString) WithNumber(yes bool) *RandomString {
	return me.withType(yes, 0)
}

func (me *RandomString) WithSpecial(yes bool) *RandomString {
	return me.withType(yes, 3)
}

func (me *RandomString) WithLower(yes bool) *RandomString {
	return me.withType(yes, 1)
}

func (me *RandomString) WithUpper(yes bool) *RandomString {
	return me.withType(yes, 2)
}

// 随机字符串
func (me *RandomString) Build() string {
	result := make([]rune, me.length)
	types := make([]int, 0)
	for _, v := range me.typeIndex {
		types = append(types, v)
	}
	for i := 0; i < me.length; i++ {
		// 先随机类型
		t := types[rand.Intn(len(types))]
		chars := randomChars[t]
		result[i] = chars[rand.Intn(len(chars))]
	}

	return string(result)
}

// SplitAndTrim 使用指定字符串切割字符串并trim
// trim 只处理头尾的空格、换行符、制表符
func SplitAndTrim(str, sep string) []string {
	var trimFunc = func(r rune) bool {
		return r == ' ' || r == '\r' || r == '\t' || r == '\n'
	}

	var result []string

	left := str
	sepLen := len(sep)
	idx := strings.Index(str, sep)
	if idx < 0 {
		v := strings.TrimFunc(str, trimFunc)
		if v != "" {
			result = []string{v}
		} else {
			result = []string{}
		}
	} else {
		result = make([]string, 0)
		for idx >= 0 {
			v := strings.TrimFunc(left[:idx], trimFunc)
			if v != "" {
				result = append(result, v)
			}

			left = left[idx+sepLen:]
			idx = strings.Index(left, sep)
		}
		v := strings.TrimFunc(left, trimFunc)
		if v != "" {
			result = append(result, v)
		}
	}

	return result
}

// transform the first letter to upper case
//
//	Capitalize("hello world") => "Hello world"
func Capitalize(str string) string {
	runes := []rune(str)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
