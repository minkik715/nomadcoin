package db

import (
	"fmt"
	"github.com/rlaalsrl715/nomadcoin/utils"
	bolt "go.etcd.io/bbolt"
	"os"
)

const (
	dbName       = "blockchain"
	dataBucket   = "data"
	blocksBucket = "blocks"
	checkPoint   = "checkPoint"
)

var db *bolt.DB

func getDbName() string {
	port := os.Args[2][6:]
	return fmt.Sprintf("%s_%s.db", dbName, port)
}

func DB() *bolt.DB {
	if db == nil {
		dbPointer, err := bolt.Open(getDbName(), 0600, nil)
		db = dbPointer
		utils.HandleErr(err)
		err = db.Update(func(tx *bolt.Tx) error {
			_, err = tx.CreateBucketIfNotExists([]byte(dataBucket))
			utils.HandleErr(err)
			_, err = tx.CreateBucketIfNotExists([]byte(blocksBucket))
			utils.HandleErr(err)
			return err
		})
	}
	return db
}

func SaveBlock(hash string, data []byte) {
	err := DB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		err := bucket.Put([]byte(hash), data)
		return err
	})
	utils.HandleErr(err)
}

func SaveCheckPoint(data []byte) {
	err := DB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dataBucket))
		err := bucket.Put([]byte(checkPoint), data)
		return err
	})
	utils.HandleErr(err)
}

func Blockchain() []byte {
	var data []byte
	DB().View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dataBucket))
		data = bucket.Get([]byte(checkPoint))
		return nil
	})
	return data
}

func Block(hash string) []byte {
	var data []byte
	DB().View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		data = bucket.Get([]byte(hash))
		return nil
	})
	return data
}

func Close() {
	DB().Close()
}
