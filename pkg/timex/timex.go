package timex

import "time"

func NowUnixMillis() int64 {
	return time.Now().UnixNano() / time.Millisecond.Nanoseconds()
}

func UnixMills(t time.Time) int64 {
	return t.UnixNano() / time.Millisecond.Nanoseconds()
}
