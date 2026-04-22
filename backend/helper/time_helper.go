package helper

import (
	"fmt"
	"time"
)

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

func FormatTanggalIndo(t time.Time) string {
	bulan := []string{
		"", "Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}
	return fmt.Sprintf("%d %s %d %02d:%02d",
		t.Day(), bulan[t.Month()], t.Year(), t.Hour(), t.Minute())
}
