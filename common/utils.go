package common

import (
	"math/rand"
	"time"
)

func ToLocalTimeString(timeString string) string {
	ti, _ := time.Parse(time.RFC3339, timeString)
	localTime := ti.In(time.Now().Location())
	return localTime.Format("2006-01-02 15:04:05")
}

func ShuffleSlice[T any](slice []T) {
	for i := range slice {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}
