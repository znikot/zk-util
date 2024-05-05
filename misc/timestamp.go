package misc

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Timestamp 匹配Java的 Timestamp转换
type Timestamp time.Time

type Date time.Time

type DateTime time.Time

type Time time.Time

// Now 获取当前的时间，以 Timestamp 类型返回
func Now() Timestamp {
	return ToTimestamp(time.Now())
}

// UnmarshalJSON 将毫秒数转换成 time.Time
func (j *Timestamp) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	if str == "" {
		return nil
	}
	millis, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return err
	}
	*j = Timestamp(time.Unix(0, millis*int64(time.Millisecond)))
	return nil
}

// MarshalJSON 转换JSON
func (j *Timestamp) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	origin := time.Time(*j)
	buf.WriteString(strconv.FormatInt(origin.UnixNano()/int64(time.Millisecond), 10))
	return buf.Bytes(), nil
}

func (d *Date) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	if str == "" {
		return nil
	}
	dt, err := ParseTime("yyyy-MM-dd", str)
	if err != nil {
		return err
	}
	*d = Date(dt)
	return nil
}

func (d *Date) MarshalJSON() ([]byte, error) {
	if time.Time(*d).IsZero() {
		return []byte(fmt.Sprintf("%q", "")), nil
	}
	return []byte(fmt.Sprintf("%q", time.Time(*d).Format("2006-01-02"))), nil
	// dt := FormatTime("yyyy-MM-dd", time.Time(*d))
	// return []byte(fmt.Sprintf("%q", dt)), nil
}

func (d *DateTime) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	if str == "" {
		return nil
	}
	dt, err := ParseTime("yyyy-MM-dd HH:mm:ss", str)
	if err != nil {
		return err
	}
	*d = DateTime(dt)
	return nil
}

func (d *DateTime) MarshalJSON() ([]byte, error) {
	// dt := FormatTime("yyyy-MM-dd HH:mm:ss", time.Time(*d))
	if time.Time(*d).IsZero() {
		return []byte(fmt.Sprintf("%q", "")), nil
	}
	return []byte(fmt.Sprintf("%q", time.Time(*d).Format("2006-01-02 15:04:05"))), nil
	// return []byte(fmt.Sprintf("%q", dt)), nil
}

func (t *Time) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	if str == "" {
		return nil
	}
	dt, err := ParseTime("HH:mm:ss", str)
	if err != nil {
		return err
	}
	*t = Time(dt)
	return nil
}

func (t *Time) MarshalJSON() ([]byte, error) {
	if time.Time(*t).IsZero() {
		return []byte(fmt.Sprintf("%q", "")), nil
	}
	return []byte(fmt.Sprintf("%q", time.Time(*t).Format("15:04:05"))), nil
	// dt := FormatTime("HH:mm:ss", time.Time(*t))
	// return []byte(fmt.Sprintf("%q", dt)), nil
}

// ToTime 转换成golang的Time
func (j Timestamp) ToTime() time.Time {
	return time.Time(j)
}

// FromMillis 从毫秒数转换为时间
func (j *Timestamp) FromMillis(millis int64) {
	*j = Timestamp(time.Unix(0, millis*int64(time.Millisecond)))
}

// ToMillis 将时间转换为毫秒数
func (j *Timestamp) ToMillis() int64 {
	return j.ToTime().UnixNano() / int64(time.Millisecond)
}

// Format 支持 Java 的时间格式
func (j *Timestamp) Format(layout string) string {
	return FormatTime(layout, j.ToTime())
}

func ToTimestamp(t time.Time) Timestamp {
	return Timestamp(t)
}

func FormatTime(layout string, time time.Time) string {
	return time.Format(parseLayout(layout))
}

// ParseTime 转换时间
// 鉴于golang解析毫秒的时候，不能直接用 000（这个有点恶心），需要加上点
func ParseTime(layout, value string) (time.Time, error) {
	//找第一个S
	idx := strings.Index(layout, "S")
	if idx != -1 && layout[idx-1:idx] != "." {
		//前面没有点…手动加上去
		layout = layout[0:idx] + "." + layout[idx:]
		if idx < len(value) {
			value = value[0:idx] + "." + value[idx:]
		}
	}
	return time.ParseInLocation(parseLayout(layout), value, time.Local)
}

func parseLayout(layout string) string {
	layout = strings.Replace(layout, "yyyy", "2006", 1)
	layout = strings.Replace(layout, "yy", "06", 1)
	layout = strings.Replace(layout, "MM", "01", 1)
	layout = strings.Replace(layout, "dd", "02", 1)
	layout = strings.Replace(layout, "HH", "15", 1)
	layout = strings.Replace(layout, "mm", "04", 1)
	layout = strings.Replace(layout, "ss", "05", 1)
	layout = strings.Replace(layout, "S", "0", -1)

	return layout
}
