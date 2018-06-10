package ui

import (
	"time"

	"k8s.io/apimachinery/pkg/util/duration"
)

func TimeToDurationStr(t time.Time) string {
	return duration.ShortHumanDuration(time.Now().Sub(t))
}
