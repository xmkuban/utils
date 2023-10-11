package utils

import (
	"fmt"
	"strconv"
	"time"
)

// ZeroSecondTimestampByTime 获取日期对应的0点时间戳，单位s
func ZeroSecondTimestampByTime(date time.Time) int64 {
	year, month, day := date.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, date.Location()).Unix()
}

// ZeroMillSecondTimestampByTime 获取日期对应的0点时间戳，单位ms
func ZeroMillSecondTimestampByTime(date time.Time) int64 {
	year, month, day := date.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, date.Location()).UnixMilli()
}

// ZeroDateByTime 获取日期对应的0点时间
func ZeroDateByTime(date time.Time) time.Time {
	year, month, day := date.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, date.Location())
}

// ZeroMillSecondTimestampByTimestamp 通过时间戳转成获取日期对应的0点时间戳
func ZeroMillSecondTimestampByTimestamp(t int64) int64 {
	date := SecondTimestampToTime(t)
	year, month, day := date.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, date.Location()).UnixMilli()
}

// ZeroDateByTimestamp 通过时间戳转成获取日期对应的0点时间
func ZeroDateByTimestamp(t int64) time.Time {
	date := SecondTimestampToTime(t)
	year, month, day := date.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, date.Location())
}

// SecondTimestampToTime 将时间戳转成 time.Time
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
	var stamp = fmt.Sprintf("%d", t.Time().UnixMilli())
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
