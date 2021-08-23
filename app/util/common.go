package util

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ToInt(data interface{}) (int, error) {
	switch data.(type) {
	case int:
		return data.(int), nil
	}

	var numStr string
	if _, ok := data.(string); ok {
		numStr = data.(string)
	} else {
		numStr = fmt.Sprintf("%v", data)
	}
	ret, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return int(ret), nil
}

func ToString(data interface{}) string {
	switch data.(type) {
	case string:
		return data.(string)
	case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", data)
	default:
		return fmt.Sprintf("%v", data)
	}
}

func Now(arg ...string) string {
	format := "Y-m-d H:i:s"
	if len(arg) > 0 {
		format = arg[0]
	}

	location := "Local"
	if len(arg) > 1 {
		location = arg[1]
	}
	return TimeStampToString(time.Now().Unix(), format, location)
}

func Today(arg ...string) string {
	format := "Y-m-d"
	if len(arg) > 0 {
		format = arg[0]
	}

	location := "Local"
	if len(arg) > 1 {
		location = arg[1]
	}
	return TimeStampToString(time.Now().Unix(), format, location)
}

func Year() int {
	return time.Now().Year()
}

func TimeStamp() int64 {
	return time.Now().Unix()
}

func TimeStampToString(sec int64, arg ...string) string {
	location := "Local"
	if len(arg) > 1 {
		location = arg[1]
	}
	local, _ := time.LoadLocation(location)

	format := "2006-01-02 15:04:05"
	if len(arg) > 0 {
		format = arg[0]
		switch format {
		case "Y-m-d H:i:s":
			format = "2006-01-02 15:04:05"
		case "Y/m/d H:i:s":
			format = "2006/01/02 15:04:05"
		case "YmdHis":
			format = "20060102150405"
		case "Y-m-d H:i":
			format = "2006-01-02 15:04"
		case "Y/m/d H:i":
			format = "2006/01/02 15:04"
		case "Y-m-d":
			format = "2006-01-02"
		case "Y/m/d":
			format = "2006/01/02"
		case "Ymd":
			format = "20060102"
		case "年/月/日":
			format = "2006年/01月/02日"
		case "年-月-日":
			format = "2006年-01月-02日"
		case "年月日":
			format = "2006年01月02日"
		}
	}

	t := time.Unix(sec, 0).In(local)
	return t.Format(format)
}

func StringToTimeStamp(str string, arg ...string) int64 {
	location := "Local"
	if len(arg) > 0 {
		location = arg[0]
	}
	local, _ := time.LoadLocation(location)
	str = strings.Replace(str, "/", "-", 0)

	var format string
	reg := regexp.MustCompile(`\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`)
	if ok := reg.MatchString(str); ok {
		format = "2006-01-02 15:04:05"
	} else {
		format = "2006-01-02"
	}

	t, _ := time.ParseInLocation(format, str, local)
	return t.Unix()
}
