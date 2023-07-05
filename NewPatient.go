package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func Newpatient(rw http.ResponseWriter, r *http.Request) {
	var patient Patient
	err := json.NewDecoder(r.Body).Decode(&patient)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := client.Database("crud").Collection("patients")
	daysToHospitalize := r.URL.Query().Get("days")
	num, err := strconv.Atoi(daysToHospitalize)
	if err != nil {
		http.Error(rw, "Invalid hospitalization days", http.StatusBadRequest)
		return
	}
	date := time.Now().AddDate(0, 0, num)
	patient.DateOfDischarge = date
	_, err = collection.InsertOne(ctx, bson.M{
		"patientid":       patient.PatientID,
		"firstname":       patient.FirstName,
		"lastname":        patient.LastName,
		"dob":             patient.DOB,
		"gender":          patient.Gender,
		"contactnumber":   patient.ContactNumber,
		"medicalhistory":  patient.MedicalHistory,
		"dateofdischarge": patient.DateOfDischarge,
	})

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte("Patient created successfully"))
}
