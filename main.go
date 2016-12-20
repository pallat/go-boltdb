package main

import (
	"flag"
	"fmt"
	"os"

	"log"

	"github.com/boltdb/bolt"
)

var (
	key string
	val string
)

func init() {
	flag.StringVar(&key, "k", "foo", "-k=foo")
	flag.StringVar(&val, "v", "bar", "-v=bar")
}

func main() {
	flag.Parse()

	fmt.Println("key", key, "value", val)

	db, err := bolt.Open("my.db", 0666, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(db.Path())

	insert(db, "widgets", key, val)

	v, err := query(db, "widgets", key)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(v)
}

func query(db *bolt.DB, bucket, k string) (string, error) {
	var v []byte
	if err := db.View(func(tx *bolt.Tx) error {
		v = tx.Bucket([]byte(bucket)).Get([]byte(k))
		return nil
	}); err != nil {
		return "", err
	}
	return string(v), nil
}

func insert(db *bolt.DB, bucket, k, v string) error {
	if err := db.Update(func(tx *bolt.Tx) error {
		// Create a bucket.
		b, err := tx.CreateBucket([]byte(bucket))
		if err != nil {
			return err
		}

		// Set the value "bar" for the key "foo".
		if err := b.Put([]byte(k), []byte(v)); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
