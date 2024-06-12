package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/LKSprod/Ott-o-Meter/db"
	"github.com/gorilla/mux"
)

func main() {
	//Init Datenbank
	db, err := db.OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/", sayHello)

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

func printPlants(w http.ResponseWriter) {

}
