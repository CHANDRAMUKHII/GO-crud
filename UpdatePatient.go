package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func UpdatePatient(rw http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	segments := strings.Split(path, "/")
	id := segments[len(segments)-1]

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var updatedPatient Patient
	err := json.NewDecoder(r.Body).Decode(&updatedPatient)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	collection := client.Database("crud").Collection("patients")

	update := bson.M{
		"$set": bson.M{
			"firstname":      updatedPatient.FirstName,
			"lastname":       updatedPatient.LastName,
			"dob":            updatedPatient.DOB,
			"gender":         updatedPatient.Gender,
			"contactnumber":  updatedPatient.ContactNumber,
			"medicalhistory": updatedPatient.MedicalHistory,
		},
	}

	_, err = collection.UpdateOne(ctx, bson.M{"patientid": id}, update)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Patient updated successfully"))
}
