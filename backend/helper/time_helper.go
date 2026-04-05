package helper

import "time"

func FormatTimeRFC3339(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	return t.Format(time.RFC3339)
}

func TimeNowWIB() time.Time {
	wib := time.FixedZone("WIB", 7*60*60)
	return time.Now().In(wib)
}
