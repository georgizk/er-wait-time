package rsc

import (
	"time"
)

type Patient struct {
	PatientNumber uint64
	CheckInTime   time.Time
	Priority      int
}
