package db

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"time"
)

const dbBucket = "marketingClient"

type BoltDb struct {
	db *bolt.DB
}

func NewBoltDb(path string) (*BoltDb, error) {
	if path == "" {
		path = "database.db"
	}
	db, err := bolt.Open(path, 0600, nil)
	return &BoltDb{db}, err
}

func (b *BoltDb) Save(m map[string]string) error{
	enc, err := json.Marshal(m)
	if err != nil {
		return err
	}

	err = b.db.Update(func(tx *bolt.Tx) error {
		req, err := tx.CreateBucketIfNotExists([]byte(dbBucket))

		if err != nil {
			return err
		}

		err = req.Put([]byte(time.Now().Format(time.UnixDate)), enc)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
func (b *BoltDb) GetAll() {
	err := b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dbBucket))
		if b == nil {
			log.Println("file is empty")
			return nil
		}
		err := b.ForEach(func(k, v []byte) error {
			fmt.Printf("key=%s, value=%s\n", k, v)
			return nil
		})
		if err != nil {
			log.Println(err)
		}
		return nil
	})
	if err != nil {
		log.Println(err)
		return
	}
}

func (b *BoltDb) DeleteAll() {
	err := b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dbBucket))
		if b == nil {
			return nil
		}
		tx.DeleteBucket([]byte(dbBucket))
		return nil
	})
	if err != nil {
		log.Print(err)
	}
}

func (b *BoltDb) Close() {
	b.db.Close()
}
