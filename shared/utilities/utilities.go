package utilities

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func ArrayToString(a interface{}, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
}

func StringToStringArray(array string, delim string) []string {
	s := strings.Split(array, delim)
	return s
}

func StringToIntArray(array string, delim string) []int {
	s := strings.Split(array, delim)
	arrayInt := []int{}
	for _, v := range s {
		value, err := strconv.Atoi(v)
		if err == nil {
			arrayInt = append(arrayInt, value)
		}
	}
	return arrayInt
}

func StringToInt64Array(array string, delim string) []int64 {
	s := strings.Split(array, delim)
	arrayInt := []int64{}
	for _, v := range s {
		value, err := strconv.Atoi(v)
		if err == nil {
			arrayInt = append(arrayInt, int64(value))
		}
	}
	return arrayInt
}

func UnixToISOTimeString(input int64) string {
	ts := time.UnixMilli(input)
	return ts.Format("2006-01-02T15:04:05Z")
}
