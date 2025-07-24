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
	Update(func(*bbolt.Tx) error) error
	View(func(*bbolt.Tx) error) error
	Close() error
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

func NewBoltDB() (*BoltDB, error) {
	path, err := getDBPaht()
	if err != nil {
		return nil, err
	}

	db, err := bbolt.Open(path, 0666, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("kv"))
		return err
	})
	if err != nil {
		return nil, err
	}

	return &BoltDB{db: db}, nil
}

func (b *BoltDB) Update(fn func(*bbolt.Tx) error) error {
	return b.db.Update(fn)
}

func (b *BoltDB) View(fn func(*bbolt.Tx) error) error {
	return b.db.View(fn)
}

func (b *BoltDB) Close() error {
	return b.db.Close()
}
