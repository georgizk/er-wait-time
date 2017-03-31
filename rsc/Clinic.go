package rsc

import (
	"log"
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
	log.Println("Added a patient with num", patient.PatientNumber)
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
		log.Println("got to id", i, "patientNumber", patientNumber)
		if patient.PatientNumber == patientNumber {
			log.Println("got a match")
			// put removed element at the end of the array, then return
			// array that is one element shorter
			if len(s) > 1 {
				s[len(s)-1], s[i] = s[i], s[len(s)-1]
			}
			clinic.QueuedPatients = s[:len(s)-1]
			visitTime := VisitTime{time.Since(patient.CheckInTime).Seconds()}
			clinic.VisitTimes = append(clinic.VisitTimes, visitTime)
			// need to return after this, otherwise the patient will get "removed" twice
			return

		}
	}
}
