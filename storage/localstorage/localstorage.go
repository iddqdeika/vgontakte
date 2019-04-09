package localstorage

import (
	"fmt"
	"github.com/boltdb/bolt"
	"strings"
	"time"
	"vgontakte/vgontakte"
)

func GetLocalStorage(dbFileName string) vgontakte.Storage {
	storage := &boltStorage{}
	storage.Init(dbFileName)
	return storage
}

type boltStorage struct {
	fileName string
	db       *bolt.DB
	//buckets		map[string]struct{}
}

func (s *boltStorage) Init(dbFileName string) error {
	s.fileName = dbFileName
	db, err := bolt.Open(s.fileName, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	s.db = db
	return nil
}

func (s *boltStorage) Get(path string) ([]byte, error) {
	var result []byte
	err := s.db.View(func(tx *bolt.Tx) error {
		var err error
		result, err = s.get(path, tx)
		return err
	})
	return result, err
}

func (s *boltStorage) Update(path string, value string) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		err := s.update(path, value, tx)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *boltStorage) get(path string, tx *bolt.Tx) ([]byte, error) {
	paths := getPath(path)
	if len(paths) > 0 {
		if len(paths) == 1 {
			return nil, fmt.Errorf("invalid path: %v", path)
		} else {
			bucket, err := s.ensurePathWrite(paths[:len(paths)-1], tx)
			if err != nil {
				return nil, err
			}
			result := bucket.Get([]byte(paths[len(paths)-1]))
			return result, err
		}
	}
	return nil, fmt.Errorf("void path")
}

func (s *boltStorage) update(path, value string, tx *bolt.Tx) error {
	paths := getPath(path)
	if len(paths) > 0 {
		if len(paths) == 1 {
			return fmt.Errorf("invalid path: %v", path)
		} else {
			bucket, err := s.ensurePath(paths[:len(paths)-1], tx)
			if err != nil {
				return err
			}
			err = bucket.Put([]byte(paths[len(paths)-1]), []byte(value))
			if err != nil {
				return err
			}
			return nil
		}
	}
	return fmt.Errorf("void path")
}

func (s *boltStorage) ensurePath(paths []string, tx *bolt.Tx) (*bolt.Bucket, error) {
	var tpath string
	var tbuck *bolt.Bucket
	for _, val := range paths {
		if tpath != "" {
			tpath += "."
		}
		tpath += val
		if tbuck == nil {
			var err error
			tbuck, err = ensureBucket(tpath, tx)
			if err != nil {
				return nil, err
			}
		} else {
			var err error
			tbuck, err = ensureNestedBucket(tpath, tbuck)
			if err != nil {
				return nil, err
			}
		}
	}
	return tbuck, nil
}

func (s *boltStorage) ensurePathWrite(paths []string, tx *bolt.Tx) (*bolt.Bucket, error) {
	var tpath string
	var tbuck *bolt.Bucket
	for _, val := range paths {
		if tpath != "" {
			tpath += "."
		}
		tpath += val
		if tbuck == nil {
			tbuck = tx.Bucket([]byte(tpath))
		} else {
			tbuck = tbuck.Bucket([]byte(tpath))
		}
	}
	return tbuck, nil
}

func ensureBucket(name string, tx *bolt.Tx) (*bolt.Bucket, error) {
	b, err := tx.CreateBucketIfNotExists([]byte(name))
	if err != nil {
		return nil, err
	}
	return b, nil
}

func ensureNestedBucket(name string, b *bolt.Bucket) (*bolt.Bucket, error) {
	b, err := b.CreateBucketIfNotExists([]byte(name))
	if err != nil {
		return nil, err
	}
	return b, nil
}

func getPath(path string) []string {
	return strings.Split(strings.Replace(path, " ", "", -1), ".")
}
