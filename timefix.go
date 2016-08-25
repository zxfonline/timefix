// Copyright 2016 zxfonline@sina.com. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package timefix

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var (
	//时间格式
	TimeFormate string = "2006-01-02 15:04:05"
	//时差修正 毫秒差值
	TimeFix int64
)

const (
	//一天的毫秒数
	MILLISECONDS_OF_DAY = 24 * MILLISECONDS_OF_HOUR
	//一小时的毫秒数
	MILLISECONDS_OF_HOUR = 60 * MILLISECONDS_OF_MINUTE
	//一分钟的毫秒数
	MILLISECONDS_OF_MINUTE = 60 * MILLISECONDS_OF_SECOND
	//一秒的毫秒数
	MILLISECONDS_OF_SECOND = 1000
)

//重置时间，传入utc标准毫秒时间
func ResetTime(t int64) {
	TimeFix = time.Now().UTC().UnixNano()/1e6 - t
}

//当前本地时间 毫秒 已修正
func MillisTime() int64 {
	return time.Now().Local().UnixNano()/1e6 - TimeFix
}

//当前本地时间 秒 已修正
func SecondTime() int32 {
	return int32(MillisTime() / 1e3)
}

//当前UTC时间 毫秒  已修正
func MillisUTCTime() int64 {
	return time.Now().UTC().UnixNano()/1e6 - TimeFix
}

//当前UTC时间 秒  已修正
func SecondUTCTime() int32 {
	return int32(MillisUTCTime() / 1e3)
}

//时间转成毫秒
func TimeMillis(t time.Duration) int64 {
	return t.Nanoseconds() / 1e6
}

//时间转成秒
func TimeSecond(t time.Duration) int32 {
	return int32(t.Nanoseconds() / 1e9)
}

//ms毫秒时间转成utc时间(根据传入的时间的时区在返回的时间上进行时区纠正time.Time().In(tm.Location()))
func Ms2Time(ms int64) time.Time {
	return time.Unix(ms/1e3, 0).UTC()
}

// 获得指定时间的凌晨时间
func Time2Midnight(tm time.Time) time.Time {
	year, month, day := tm.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, tm.Location())
}

//获取一个时间点的 x日后的凌晨时间
func NextMidnight(tm time.Time, day int) time.Time {
	midTime := Time2Midnight(tm)
	ms := midTime.UnixNano()/1e6 + int64(day*MILLISECONDS_OF_DAY)
	return Ms2Time(ms).In(tm.Location())
}

// 从一个时间戳获取下一个准点时间
func NextHour(tm time.Time) time.Time {
	year, month, day := tm.Date()
	hour, _, _ := tm.Clock()
	return time.Date(year, month, day, hour+1, 0, 0, 0, tm.Location())
}

// 判断两个时间是否是同一天（默认将两个时间转换成0时区的时间进行比较）
func OtherDay(t1, t2 time.Time) bool {
	year1, month1, day1 := t1.UTC().Date()
	year2, month2, day2 := t2.UTC().Date()
	return year1 == year2 && month1 == month2 && day1 == day2
}

// date format: "2006-01-02 13:04:00"
func S2UnixTime(value string, loc *time.Location) time.Time {
	re := regexp.MustCompile(`([\d]+)-([\d]+)-([\d]+) ([\d]+):([\d]+):([\d]+)`)
	slices := re.FindStringSubmatch(value)
	if slices == nil || len(slices) != 7 {
		panic(fmt.Errorf("time[%s] format error, expect format: 2006-01-02 13:04:00 \n", value))
	}
	year, _ := strconv.Atoi(slices[1])
	month, _ := strconv.Atoi(slices[2])
	day, _ := strconv.Atoi(slices[3])
	hour, _ := strconv.Atoi(slices[4])
	min, _ := strconv.Atoi(slices[5])
	sec, _ := strconv.Atoi(slices[6])
	return time.Date(year, time.Month(month), day, hour, min, sec, 0, loc)
}

//获取指定时间的下一个周末得时间,自己处理好时区的问题
func NextSundayMS(tm time.Time) time.Time {
	mt := NextMidnight(tm, 0)
	weekday := tm.Weekday()
	if weekday == time.Sunday {
		return Ms2Time(mt.UnixNano()/1e6 + 7*MILLISECONDS_OF_DAY).In(tm.Location())
	}
	return Ms2Time(mt.UnixNano()/1e6 + int64((7-int(weekday))*MILLISECONDS_OF_DAY)).In(tm.Location())
}

// 返回当前的整点时间
func SharpClock(tm time.Time) time.Time {
	year, month, day := tm.Date()
	hour, _, _ := tm.Clock()
	return time.Date(year, month, day, hour, 0, 0, 0, tm.Location())
}

//func main() {
//	t1 := S2UnixTime("2016-05-22 07:59:59", time.Local)
//	t2 := NextSundayMS(t1.UTC())
//	fmt.Println(t1.Format(TimeFormate))
//	fmt.Println(t2.UTC().Format(TimeFormate))
//}

//返回从整点到现在的差值
func NowToSharpClock(tm time.Time) time.Duration {
	return time.Duration(tm.UnixNano() - SharpClock(tm).UnixNano())
}
