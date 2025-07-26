package libs

import "fmt"

type Service struct {
	db           DB
	kvBucketName string
}

func NewService(db DB, kvBucketName string) *Service {
	return &Service{db, kvBucketName}
}

func (s *Service) Add(key, value string) error {
	return s.db.Update(func(tx Tx) error {
		b := tx.Bucket([]byte(s.kvBucketName))
		if b == nil {
			return fmt.Errorf("bucket %s not found", s.kvBucketName)
		}
		return b.Put([]byte(key), []byte(value))
	})
}

func (s *Service) Get(key string) (string, error) {
	var value []byte
	err := s.db.View(func(tx Tx) error {
		b := tx.Bucket([]byte(s.kvBucketName))
		if b == nil {
			return fmt.Errorf("bucket %s not found", s.kvBucketName)
		}
		value = b.Get([]byte(key))
		return nil
	})
	if err != nil {
		return "", err
	}
	return string(value), nil
}

