package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	//"go.mongodb.org/mongo-driver/mongo/readpref"
)

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

type event struct {
	ID          string `json:"ID"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
}

type dog struct {
	Name		string `json:"Name"`
	Breed       string `json:"Breed"`
	Details string `json:"Details"`
}


type allEvents []event

var events = allEvents{
	{
		ID:          "1",
		Title:       "golang rest",
		Description: "desc.",
	},
}

type allDogs []dog

var dogs = allDogs{}

func createEvent(w http.ResponseWriter, r *http.Request) {
	var newEvent event
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}

	json.Unmarshal(reqBody, &newEvent)
	events = append(events, newEvent)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newEvent)
}

func getOneEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]

	for _, singleEvent := range events {
		if singleEvent.ID == eventID {
			json.NewEncoder(w).Encode(singleEvent)
		}
	}
}

func getAllEvents(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(events)
}

func updateEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]
	var updatedEvent event

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}
	json.Unmarshal(reqBody, &updatedEvent)

	for i, singleEvent := range events {
		if singleEvent.ID == eventID {
			singleEvent.Title = updatedEvent.Title
			singleEvent.Description = updatedEvent.Description
			events = append(events[:i], singleEvent)
			json.NewEncoder(w).Encode(singleEvent)
		}
	}
}

func deleteEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]

	for i, singleEvent := range events {
		if singleEvent.ID == eventID {
			events = append(events[:i], events[i+1:]...)
			fmt.Fprintf(w, "The event with ID %v has been deleted successfully", eventID)
		}
	}
}

func mongoCreateEvent(w http.ResponseWriter, r *http.Request) {
	var newDog dog
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}

	json.Unmarshal(reqBody, &newDog)
	dogs = append(dogs, newDog)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newDog)
	
}

func main() {
	// router := mux.NewRouter().StrictSlash(true)
	// router.HandleFunc("/", homeLink)
	// router.HandleFunc("/event", createEvent).Methods("POST")
	// router.HandleFunc("/events", getAllEvents).Methods("GET")
	// router.HandleFunc("/events/{id}", getOneEvent).Methods("GET")
	// router.HandleFunc("/events/{id}", updateEvent).Methods("PATCH")
	// router.HandleFunc("/events/{id}", deleteEvent).Methods("DELETE")

	// var serverRunning = "Go server is running on port 8080."
	// fmt.Printf("%v\n", serverRunning)
	// log.Fatal(http.ListenAndServe(":8080", router))

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/mongo/event", mongoCreateEvent).Methods("POST")

	// Mongo DB tutorial
	client, err := mongo.NewClient(options.Client().ApplyURI(getMongoURI()))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	// err = client.Ping(ctx, readpref.Primary())
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// databases, err := client.ListDatabaseNames(ctx, bson.M{})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(databases)
	quickstartDatabase := client.Database("quickstart")
	pupsCollection := quickstartDatabase.Collection("pups")
	bbsCollection := quickstartDatabase.Collection("bbs")
	pupResult, err := pupsCollection.InsertOne(ctx, bson.D{
		{Key: "title", Value: "Mondo's Pups"},
		{Key: "author", Value: "Mondo"},
		{Key: "tags", Value: bson.A{"pups", "pooches", "bbs"}},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pupResult.InsertedID)
	bbsResult, err := bbsCollection.InsertMany(ctx, []interface{}{
		bson.D{
			{Key: "dog", Value: pupResult.InsertedID},
			{Key: "name", Value: "Jozy"},
			{Key: "details", Value: "golden girl joz miss daisy."},
		},
		bson.D{
			{Key: "dog", Value: pupResult.InsertedID},
			{Key: "name", Value: "Tonks"},
			{Key: "details", Value: "t-man jones baby boy."},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(bbsResult.InsertedIDs)

}
