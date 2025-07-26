//go:generate mockgen -destination=../mocks/db.go -package=mocks -source=db.go
package libs

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/viper"
	"go.etcd.io/bbolt"
)

type DB interface {
	Update(fn func(Tx) error) error
	View(fn func(Tx) error) error
	Close() error
}

type Tx interface {
	Bucket(name []byte) Bucket
}

type Bucket interface {
	Put(key, value []byte) error
	Get(key []byte) []byte
}

type BoltDB struct {
	db *bbolt.DB
}

func getDBPath() (string, error) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	folderName := fmt.Sprintf(".%s", viper.GetString("app_name"))
	dbPath := path.Join(dirname, folderName, "bbolt.db")
	if _, err := os.Stat(path.Dir(dbPath)); os.IsNotExist(err) {
		if err := os.MkdirAll(path.Dir(dbPath), 0755); err != nil {
			return "", err
		}
	}

	return dbPath, nil
}

func NewBoltDB(kvBucketName string) (*BoltDB, error) {
	path, err := getDBPath()
	if err != nil {
		return nil, err
	}

	db, err := bbolt.Open(path, 0666, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(kvBucketName))
		return err
	})
	if err != nil {
		return nil, err
	}

	return &BoltDB{db: db}, nil
}

func (b *BoltDB) Update(fn func(Tx) error) error {
	return b.db.Update(func(tx *bbolt.Tx) error {
		return fn(&BoltTx{tx})
	})
}

func (b *BoltDB) View(fn func(Tx) error) error {
	return b.db.View(func(tx *bbolt.Tx) error {
		return fn(&BoltTx{tx})
	})
}

func (b *BoltDB) Close() error {
	return b.db.Close()
}

type BoltTx struct {
	*bbolt.Tx
}

func (t *BoltTx) Bucket(name []byte) Bucket {
	bucket := t.Tx.Bucket(name)
	if bucket == nil {
		return nil
	}
	return &BoltBucket{bucket}
}

type BoltBucket struct {
	db *bbolt.Bucket
}

func (b *BoltBucket) Put(key, value []byte) error {
	return b.db.Put(key, value)
}

func (b *BoltBucket) Get(key []byte) []byte {
	return b.db.Get(key)
}
