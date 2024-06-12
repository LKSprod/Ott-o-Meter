package db

import (
	"fmt"
	"log"

	"go.etcd.io/bbolt"
)

type DB struct {
	bolt *bbolt.DB
}

func OpenDB() (*DB, error) {
	db, err := bbolt.Open("ottometerDB.bolt", 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %e", err)
	}

	log.Printf("Opened Database \n")

	return &DB{
		bolt: db,
	}, nil
}

func (db *DB) Close() {
	db.bolt.Close()
	log.Printf("Closed Database \n")
}
