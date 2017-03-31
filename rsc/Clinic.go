package rsc

import (
	"log"
	"time"
)

type Clinic struct {
	Address              string
	Name                 string
	QueuedPatients       []Patient `json:"-"`
	NextPatientNumber    uint64    `json:"-"`
	NumProcessedPatients uint64    `json:"-"`
	AverageWaitTime      float64   `json:"-"`
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
	return clinic.AverageWaitTime
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
			clinic.NumProcessedPatients += 1
			visitTime := time.Since(patient.CheckInTime).Seconds()
			prevAvg := clinic.AverageWaitTime
			clinic.AverageWaitTime = prevAvg + (visitTime-prevAvg)/float64(clinic.NumProcessedPatients)
			// need to return after this, otherwise the patient will get "removed" twice
			return

		}
	}
}
