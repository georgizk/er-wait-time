package rsc

import (
	"log"
	"time"
)

type Clinic struct {
	Address                      string
	Name                         string
	QueuedPatients               []Patient `json:"-"`
	NextPatientNumber            uint64    `json:"-"`
	NumProcessedPatients         uint64    `json:"-"`
	AverageWaitTime              float64   `json:"avgWaitTime"`
	AverageTimeBetweenArrivals   float64   `json:"avgArrTime"`
	AverageTimeBetweenDepartures float64   `json:"avgDepTime"`
	LastArrival                  time.Time `json:"-"`
	LastDeparture                time.Time `json:"-"`
	TotalArrivals                uint64    `json:"-"`
	TotalDepartures              uint64    `json:"-"`
}

func NewClinic(name string, address string) Clinic {
	clinic := Clinic{Name: name,
		Address: address}
	clinic.ResetQueue()
	return clinic
}

func (clinic *Clinic) AddPatient(priority int) Patient {
	patient := Patient{
		PatientNumber: clinic.NextPatientNumber,
		CheckInTime:   time.Now(),
		Priority:      priority}
	clinic.QueuedPatients = append(clinic.QueuedPatients, patient)
	clinic.NextPatientNumber += 1
	clinic.TotalArrivals += 1
	arrivalInterval := time.Since(clinic.LastArrival)
	clinic.LastArrival = time.Now()
	avgArrRate := clinic.AverageTimeBetweenArrivals
	clinic.AverageTimeBetweenArrivals = avgArrRate + (arrivalInterval.Seconds()-avgArrRate)/float64(clinic.TotalArrivals)
	log.Println("Added a patient with num", patient.PatientNumber)
	return patient
}

func (clinic *Clinic) ResetQueue() {
	clinic.QueuedPatients = []Patient{}
	clinic.NextPatientNumber = 1
	clinic.LastArrival = time.Now()
	clinic.LastDeparture = time.Now()
	clinic.TotalArrivals = 0
	clinic.TotalDepartures = 0
	clinic.AverageTimeBetweenArrivals = 0
	clinic.AverageTimeBetweenDepartures = 0
	clinic.NumProcessedPatients = 0
	clinic.AverageWaitTime = 0
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

			clinic.TotalDepartures += 1
			departureInterval := time.Since(clinic.LastDeparture)
			clinic.LastDeparture = time.Now()
			depRate := clinic.AverageTimeBetweenDepartures
			clinic.AverageTimeBetweenDepartures = depRate + (departureInterval.Seconds()-depRate)/float64(clinic.TotalDepartures)

			// need to return after this, otherwise the patient will get "removed" twice
			return

		}
	}
}
