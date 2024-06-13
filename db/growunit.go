package db

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

var growUnitBucket = []byte("growunits")

type GrowUnit struct {
	id                        uint64
	width                     uint
	height                    uint
	depth                     uint
	carbonFilter              bool
	activeIntake              bool
	outtakeFanThroughputInM3H uint
	wattageLamp               uint
	ventilation               bool
	inside                    bool
	growMedium                GrowMedium
}

func (gu *GrowUnit) area() uint {
	return gu.width * gu.depth
}

func (gu *GrowUnit) volume() uint {
	return gu.width * gu.depth * gu.height
}

func (gu *GrowUnit) verify() error {
	if gu.width <= 0 {
		return fmt.Errorf("Invalid Width")
	}

	if gu.growMedium == "" {
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
		bucket, err := tx.CreateBucketIfNotExists(growUnitBucket)
		if err != nil {
			return fmt.Errorf("failed to open bucket: %e", err)
		}

		newId, err := bucket.NextSequence()
		if err != nil {
			return fmt.Errorf("failed to generate id: %e", err)
		}

		gu.id = newId

		json, err := json.Marshal(gu)
		if err != nil {
			return fmt.Errorf("failed to marshal: %e", err)
		}

		id = newId

		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(id))

		err = bucket.Put(b, json)
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
		bucket := tx.Bucket(growUnitBucket)
		if bucket != nil {
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
func (id int) GetGrowUnit() (gu *GrowUnit, err error) {
	if id > len(growUnitBucket) {
		return nil, fmt.Errorf("Grow Unit not in list: %e", err)
	}
	return &growUnitBucket[id], err
}

// ToDo UpdateGrowUnit
func (db *DB) UpdateGrowUnit(gu *GrowUnit, id uint64) (uint64, error) {
	err := gu.verify()
	if err != nil {
		return 0, fmt.Errorf("failed to verify GrowUnit data: %e", err)
	}

	err = db.bolt.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(growUnitBucket)
		if err != nil {
			return fmt.Errorf("failed to open bucket: %e", err)
		}

		gu.id = id

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
