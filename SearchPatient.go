package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func SearchPatient(rw http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	segments := strings.Split(path, "/")
	id := segments[len(segments)-1]
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := client.Database("crud").Collection("patients")
	switch id {
	case "":
		result, err := collection.Find(ctx, bson.M{})

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		var patients []Patient
		for result.Next(ctx) {
			var patient Patient

			err := result.Decode(&patient)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
			patients = append(patients, patient)
		}
		responseJSON, err := json.Marshal(patients)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(responseJSON)

	case "inpatient":
		currentDate := time.Now()
		result, err := collection.Find(ctx, bson.M{"dateofdischarge": bson.M{"$gt": currentDate}})
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		var patients []Patient
		for result.Next(ctx) {
			var patient Patient

			err := result.Decode(&patient)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
			patients = append(patients, patient)
		}

		responseJSON, err := json.Marshal(patients)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(responseJSON)

	case "outpatient":
		currentDate := time.Now()
		result, err := collection.Find(ctx, bson.M{"dateofdischarge": bson.M{"$lt": currentDate}})
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		var patients []Patient
		for result.Next(ctx) {
			var patient Patient

			err := result.Decode(&patient)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
			patients = append(patients, patient)
		}

		responseJSON, err := json.Marshal(patients)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(responseJSON)

	default:
		var patient Patient
		err := collection.FindOne(ctx, bson.M{"patientid": id}).Decode(&patient)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		responseJSON, err := json.Marshal(patient)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(responseJSON)
	}
}
