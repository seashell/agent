package boltdb

import (
	"context"
	"encoding/json"

	"github.com/seashell/agent/client/state"
	"github.com/seashell/agent/pkg/log"
	"github.com/seashell/agent/seashell/structs"
	bolt "go.etcd.io/bbolt"
)

var (
	rootBucketName         = []byte("/")
	configurationObjectKey = []byte("configuration")
)

// StateRepository ...
type StateRepository struct {
	db *bolt.DB
}

// Transaction :s
type Transaction struct {
	bolt.Tx
}

// Commit :
func (t *Transaction) Commit() error {
	return t.Commit()
}

// NewStateRepository creates a new BoltDB state repository
func NewStateRepository(path string, logger log.Logger) *StateRepository {
	db, err := bolt.Open(path, 0666, nil)
	if err != nil {
		panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {

		_, err := tx.CreateBucketIfNotExists(rootBucketName)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	return &StateRepository{db}

}

// Name :
func (r *StateRepository) Name() string {
	return "boltdb"
}

// Transaction :
func (r *StateRepository) Transaction(ctx context.Context) state.Transaction {
	return &Transaction{}
}

// Configuration :
func (r *StateRepository) Configuration() (*structs.Configuration, error) {

	var config *structs.Configuration

	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(rootBucketName)

		data := b.Get(configurationObjectKey)
		if data != nil {
			config = &structs.Configuration{}
			if err := decode(data, config); err != nil {
				return err
			}
		}

		return nil
	})

	return config, err
}

// SetConfiguration :
func (r *StateRepository) SetConfiguration(c *structs.Configuration) error {
	err := r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(rootBucketName)
		return b.Put(configurationObjectKey, encode(c))
	})
	return err
}

func encode(in interface{}) []byte {
	out, err := json.Marshal(in)
	if err != nil {
		panic(err)
	}
	return out
}

func decode(encoded []byte, out interface{}) error {
	return json.Unmarshal(encoded, out)
}
