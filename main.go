package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const uri = "mongodb+srv://new-user:new-user@cluster0.grve526.mongodb.net/"

var client *mongo.Client

type Patient struct {
	PatientID       string    `json:"PatientID"`
	FirstName       string    `json:"FirstName"`
	LastName        string    `json:"LastName"`
	DOB             string    `json:"DOB"`
	Gender          string    `json:"Gender"`
	ContactNumber   string    `json:"ContactNumber"`
	MedicalHistory  string    `json:"MedicalHistory"`
	DateOfDischarge time.Time `json:"DateOfDischarge"`
}

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

func RemovePatient(rw http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	segments := strings.Split(path, "/")
	id := segments[len(segments)-1]
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := client.Database("crud").Collection("patients")

	result, err := collection.DeleteOne(ctx, bson.M{"patientid": id})
	if result.DeletedCount == 0 {
		http.Error(rw, "User not deleted", http.StatusNotFound)
	}
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Patient Deleted successfully"))
}

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
func main() {

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	var err error
	client, err = mongo.Connect(context.TODO(), opts)

	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		fmt.Print("Error!")
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	http.HandleFunc("/patient", Newpatient)
	http.HandleFunc("/delete/", RemovePatient)
	http.HandleFunc("/patients/", SearchPatient)
	http.HandleFunc("/update/", UpdatePatient)
	log.Fatal(http.ListenAndServe(":5000", nil))
}
