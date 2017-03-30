package rsc

import (
	"time"
)

type Patient struct {
	PatientNumber int
	CheckInTime time.Time
	CheckOutTime time.Time
}
