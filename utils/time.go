package utils

import (
	"fmt"
	"strconv"
	"time"
)

//获取时间戳 单位s
func TimeToSecondTimestamp(t time.Time) int64 {
	timestamp := t.Unix()
	if timestamp <= 0 {
		return 0
	}
	return timestamp
}

//获取时间戳 单位ms
func TimeToMillSecondTimestamp(t time.Time) int64 {
	timestamp := t.UnixNano()
	if timestamp <= 0 {
		return 0
	}
	return timestamp / 1000000
}

//获取日期对应的0点时间戳，单位s
func GetDateZeroTimestamp(date time.Time) int64 {
	t := date.Unix()
	t = t - int64(date.Hour())*3600
	t = t - int64(date.Minute())*60
	t = t - int64(date.Second())
	return t
}

//将时间戳转成 time.Time
func SecondTimestampToTime(timestamp int64) time.Time {
	if timestamp == 0 {
		return time.Time{}
	}
	var nsec int64
	if timestamp > 10000000000 {
		nsec = timestamp % 1000
		timestamp = timestamp / 1000
	}
	t := time.Unix(timestamp, nsec)
	return t
}

type ExTime time.Time

func (t ExTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("%d", t.MillSecond())
	return []byte(stamp), nil
}

func (t *ExTime) UnmarshalJSON(b []byte) error {
	str := string(b)
	sec, _ := strconv.ParseInt(str, 10, 64)
	var nsec int64
	if len(str) == 13 {
		nsec = sec % 1000
		sec = sec / 1000
	}
	*t = ExTime(time.Unix(sec, nsec))
	return nil
}

func (t ExTime) Time() time.Time {
	return time.Time(t)
}

func (t ExTime) MillSecond() int64 {
	return TimeToMillSecondTimestamp(t.Time())
}
