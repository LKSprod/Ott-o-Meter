package database

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"go.etcd.io/bbolt"
)

type GrowMedium string

const (
	GrowMediumDirt  GrowMedium = "dirt"
	GrowMediumWater GrowMedium = "water"
	GrowMediumCocos GrowMedium = "cocos"
)

var growUnitBucketName = []byte("growunits")

type GrowUnit struct {
	Id                        uint64
	Width                     uint
	Height                    uint
	Depth                     uint
	CarbonFilter              bool
	ActiveIntake              bool
	OuttakeFanThroughputInM3H uint
	WattageLamp               uint
	Ventilation               bool
	Inside                    bool
	GrowMedium                GrowMedium
}

func (gu *GrowUnit) area() uint {
	return gu.Width * gu.Depth
}

func (gu *GrowUnit) volume() uint {
	return gu.Width * gu.Depth * gu.Height
}

func (gu *GrowUnit) verify() error {
	if gu.Width <= 0 {
		return fmt.Errorf("Invalid Width")
	}

	if gu.GrowMedium == "" {
		return fmt.Errorf("unknown growth medium")
	}

	return nil
}

func (db *DB) AddGrowUnit(gu *GrowUnit) (uint64, error) {
	err := gu.verify()
	if err != nil {
		return 0, fmt.Errorf("failed to verify GrowUnit data: %e", err)
	}

	var id uint64

	err = db.bolt.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(growUnitBucketName)
		if err != nil {
			return fmt.Errorf("failed to open bucket: %e", err)
		}

		newId, err := bucket.NextSequence()
		if err != nil {
			return fmt.Errorf("failed to generate id: %e", err)
		}

		// set id in given GrowUnit
		gu.Id = newId

		// set output value for id
		id = newId

		// encode struct to json, so it can be save in the database
		json, err := json.Marshal(gu)
		if err != nil {
			return fmt.Errorf("failed to marshal: %e", err)
		}

		// encode id to bytes, so it can be used as a key
		idAsBytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(idAsBytes, uint64(id))

		//error catching
		err = bucket.Put(idAsBytes, json)
		if err != nil {
			return fmt.Errorf("failed to insert record: %e", err)
		}

		return nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed to insert into db: %e", err)
	}

	return id, nil
}

func (db *DB) ListGrowUnits() ([]*GrowUnit, error) {
	units := make([]*GrowUnit, 0)

	err := db.bolt.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(growUnitBucketName)
		if bucket == nil {
			return nil
		}

		err := bucket.ForEach(func(k, v []byte) error {
			gu := &GrowUnit{}
			err := json.Unmarshal(v, gu)
			if err != nil {
				return fmt.Errorf("failed to unmarshal value: %e", err)
			}

			units = append(units, gu)

			return nil
		})
		if err != nil {
			return fmt.Errorf("failed enumerating growunits: %e", err)
		}

		return nil
	})

	return units, err

}

// ToDo GetGrowUnit ???
func (db *DB) GetGrowUnit(id uint64) (*GrowUnit, error) {
	var gu *GrowUnit = nil

	// encode id to bytes, so it can be used as a key
	idAsBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(idAsBytes, uint64(id))

	// open read transaction to access data in database
	err := db.bolt.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(growUnitBucketName)
		if bucket == nil {
			return nil
		}

		// load json from database
		jsonBytes := bucket.Get(idAsBytes)
		if jsonBytes == nil {
			return fmt.Errorf("could not find growunit with id %v", id)
		}

		gu = &GrowUnit{}

		// decode json into struct
		err := json.Unmarshal(jsonBytes, gu)
		if err != nil {
			return fmt.Errorf("failed to unmarshal json from database")
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to load growunit from db: %e", err)
	}

	return gu, nil
}

// ToDo UpdateGrowUnit
func (db *DB) UpdateGrowUnit(gu *GrowUnit, id uint64) (uint64, error) {
	err := gu.verify()
	if err != nil {
		return 0, fmt.Errorf("failed to verify GrowUnit data: %e", err)
	}

	err = db.bolt.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(growUnitBucketName)
		if err != nil {
			return fmt.Errorf("failed to open bucket: %e", err)
		}

		gu.Id = id

		json, err := json.Marshal(gu)
		if err != nil {
			return fmt.Errorf("failed to marshal: %e", err)
		}

		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(id))

		err = bucket.Put(b, json)
		if err != nil {
			return fmt.Errorf("failed to override record: %e", err)
		}

		return nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed to insert into db: %e", err)
	}

	return id, nil
}
