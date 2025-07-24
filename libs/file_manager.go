package libs

import "go.etcd.io/bbolt"

type Service struct {
	db DB
}

func NewService(db DB) *Service {
	return &Service{db}
}

func (s *Service) Add(key, value string) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("kv"))
		return b.Put([]byte(key), []byte(value))
	})
}

func (s *Service) Get(key string) (string, error) {
	var value []byte
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("kv"))
		value = b.Get([]byte(key))
		return nil
	})
	if err != nil {
		return "", err
	}
	return string(value), nil
}

