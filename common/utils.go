package common

import "time"

func ToLocalTimeString(timeString string) string {
	ti, _ := time.Parse(time.RFC3339, timeString)
	localTime := ti.In(time.Now().Location())
	return localTime.Format("2006-01-02 15:04:05")
}
