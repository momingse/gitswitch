package libs

import "fmt"

type DBService struct {
	db           DB
	kvBucketName string
}

func NewDBService(db DB, kvBucketName string) *DBService {
	return &DBService{db, kvBucketName}
}

func (s *DBService) Add(key, value string) error {
	return s.db.Update(func(tx Tx) error {
		b := tx.Bucket([]byte(s.kvBucketName))
		if b == nil {
			return fmt.Errorf("Bucket %s not found", s.kvBucketName)
		}
		return b.Put([]byte(key), []byte(value))
	})
}

func (s *DBService) Get(key string) (string, error) {
	var value []byte
	err := s.db.View(func(tx Tx) error {
		b := tx.Bucket([]byte(s.kvBucketName))
		if b == nil {
			return fmt.Errorf("Bucket %s not found", s.kvBucketName)
		}
		value = b.Get([]byte(key))
		return nil
	})
	if err != nil {
		return "", err
	}
	return string(value), nil
}
