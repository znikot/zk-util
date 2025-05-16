package reflection

import (
	"reflect"
	"time"

	"github.com/znikot/zk-util/misc"
)

func init() {
	AddTypeConverter(Float64, reflect.TypeOf(misc.Timestamp{}), float642MiscTimestamp)
	AddTypeConverter(Int64, reflect.TypeOf(misc.Timestamp{}), float642MiscTimestamp)
	AddTypeConverter(String, reflect.TypeOf(misc.Timestamp{}), string2MiscTimestamp)
	AddTypeConverter(String, reflect.TypeOf(misc.Date{}), string2MiscDate)
	AddTypeConverter(String, reflect.TypeOf(misc.DateTime{}), string2MiscDateTime)
}

// float64 -> misc.Timestamp
func float642MiscTimestamp(src any) any {
	millis := Cast[int64](src)
	return misc.Timestamp(time.Unix(0, millis*int64(time.Millisecond)))
}

// string --> misc.Timestamp
// yyyy-MM-dd yyyy-MM-dd HH:mm:ss
var dateLayouts = []string{"2006-01-02", "2006-01-02 15:04:05", "15:04:05"}

func parseDate[T misc.Timestamp | misc.Date | misc.DateTime](src any) (t T) {
	for _, l := range dateLayouts {
		ti, err := time.ParseInLocation(l, src.(string), time.Local)
		if err == nil {
			t = T(ti)
			return
		}
	}
	return
}

func string2MiscTimestamp(src any) any {
	return parseDate[misc.Timestamp](src)
}

func string2MiscDate(src any) any {
	return parseDate[misc.Date](src)
}

func string2MiscDateTime(src any) any {
	return parseDate[misc.DateTime](src)
}
