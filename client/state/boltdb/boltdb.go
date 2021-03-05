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
	configurationBucketName      = []byte("configuration")
	dragoConfigurationObjectKey  = []byte("drago")
	nomadConfigurationObjectKey  = []byte("nomad")
	consulConfigurationObjectKey = []byte("consul")
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

		_, err := tx.CreateBucketIfNotExists(configurationBucketName)
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

// DragoConfiguration :
func (r *StateRepository) DragoConfiguration() (*structs.DragoConfiguration, error) {

	var config *structs.DragoConfiguration

	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(configurationBucketName)

		data := b.Get(dragoConfigurationObjectKey)
		if data != nil {
			config = &structs.DragoConfiguration{}
			if err := decode(data, config); err != nil {
				return err
			}
		}

		return nil
	})

	return config, err
}

// SetDragoConfiguration :
func (r *StateRepository) SetDragoConfiguration(c *structs.DragoConfiguration) error {
	err := r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(configurationBucketName)
		return b.Put(dragoConfigurationObjectKey, encode(c))
	})
	return err
}

// NomadConfiguration :
func (r *StateRepository) NomadConfiguration() (*structs.NomadConfiguration, error) {

	var config *structs.NomadConfiguration

	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(configurationBucketName)

		data := b.Get(nomadConfigurationObjectKey)
		if data != nil {
			config = &structs.NomadConfiguration{}
			if err := decode(data, config); err != nil {
				return err
			}
		}

		return nil
	})

	return config, err
}

// SetNomadConfiguration :
func (r *StateRepository) SetNomadConfiguration(c *structs.NomadConfiguration) error {
	err := r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(configurationBucketName)
		return b.Put(nomadConfigurationObjectKey, encode(c))
	})
	return err
}

// ConsulConfiguration :
func (r *StateRepository) ConsulConfiguration() (*structs.ConsulConfiguration, error) {

	var config *structs.ConsulConfiguration

	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(configurationBucketName)

		data := b.Get(consulConfigurationObjectKey)
		if data != nil {
			config = &structs.ConsulConfiguration{}
			if err := decode(data, config); err != nil {
				return err
			}
		}

		return nil
	})

	return config, err
}

// SetConsulConfiguration :
func (r *StateRepository) SetConsulConfiguration(c *structs.ConsulConfiguration) error {
	err := r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(configurationBucketName)
		return b.Put(consulConfigurationObjectKey, encode(c))
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
