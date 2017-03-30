package rsc

import (
	"time"
)

type Clinic struct {
	Address           string
	Name              string
	QueuedPatients    []Patient   `json:"-"`
	NextPatientNumber uint64      `json:"-"`
	VisitTimes        []VisitTime `json:"-"`
}

type VisitTime struct {
	TimeSeconds float64
}

func (clinic *Clinic) AddPatient(priority int) Patient {
	patient := Patient{
		PatientNumber: clinic.NextPatientNumber,
		CheckInTime:   time.Now(),
		Priority:      priority}
	clinic.QueuedPatients = append(clinic.QueuedPatients, patient)
	clinic.NextPatientNumber += 1
	return patient
}

func (clinic *Clinic) ResetQueue() {
	clinic.QueuedPatients = []Patient{}
	clinic.NextPatientNumber = 1
}

func (clinic *Clinic) CalculateAvgWaitTime() float64 {
	visitTimes := clinic.VisitTimes
	if visitTimes == nil || len(visitTimes) == 0 {
		return 0
	}

	var averageVisitTime float64
	for _, visitTime := range visitTimes {
		averageVisitTime += visitTime.TimeSeconds
	}
	averageVisitTime /= float64(len(visitTimes))

	return averageVisitTime
}

func (clinic *Clinic) RemovePatient(patientNumber uint64) {
	s := clinic.QueuedPatients
	for i := 0; i < len(s); i++ {
		patient := s[i]
		if patient.PatientNumber == patientNumber {
			// put removed element at the end of the array, then return
			// array that is one element shorter
			if len(s) > 1 {
				s[len(s)-1], s[i] = s[i], s[len(s)-1]
			}
			clinic.QueuedPatients = s[:len(s)-1]
			visitTime := VisitTime{time.Since(patient.CheckInTime).Seconds()}
			clinic.VisitTimes = append(clinic.VisitTimes, visitTime)
		}
	}
}
