package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/LKSprod/Ott-o-Meter/database"
	"github.com/gorilla/mux"
)

func main() {
	//Init Datenbank
	db, err := database.OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/api/hello", sayHello)

	router.Path("/api/growunits").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		units, err := db.ListGrowUnits()
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		jsonBytes, err := json.Marshal(units)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(jsonBytes)
	})

	//API Add Grow Unit
	router.Path("/api/growunits").Methods("POST").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var gu database.GrowUnit
		err := decoder.Decode(&gu)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		id, err := db.AddGrowUnit(&gu)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		jsonBytes, err := json.Marshal(id)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(jsonBytes)
	})

	//API get Grow Unit
	router.Path("/api/growunits/{id:[0-9]+}").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		gu, err := db.GetGrowUnit(id)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		if gu == nil {
			w.WriteHeader(404)
			w.Write([]byte("Not Found"))
		}

		jsonBytes, err := json.Marshal(gu)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(jsonBytes)
	})

	port := 8080

	server := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%v", port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 15,
		Handler:      router,
	}

	go func() {
		log.Printf("Ott-O-Meter running on port %v", port)
		err := server.ListenAndServe()
		if err != nil {
			log.Printf("Error: %v\n", err)

		}
	}()
	fmt.Printf("Test")
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	server.Shutdown(ctx)

	log.Println("Shutting down Ott-O-Meter")
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!\n"))
}
