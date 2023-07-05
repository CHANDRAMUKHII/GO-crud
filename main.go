package main

import (
	"context"

	"fmt"
	"log"
	"net/http"

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
