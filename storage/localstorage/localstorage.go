package localstorage

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"strconv"
	"strings"
	"time"
	"vgontakte/vgontakte"
)

const (
	messageRatePath      = "message_rates"
	peerRegistrationPath = "registered_peers"
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

func (s *boltStorage) CheckPeerRegistration(peerId int) bool {
	data, err := s.Get(peerRegistrationPath + "." + strconv.Itoa(peerId))
	if err != nil {
		return false
	}
	return string(data) == "1"
}

func (s *boltStorage) RegisterPeer(peerId int) error {
	return s.Update(peerRegistrationPath+"."+strconv.Itoa(peerId), "1")
}

func (s *boltStorage) IncrementMessageRate(peerId int, fromId int, fwdDate int, messageText string, mediaAttachmentTokens []string) error {

	var attachments string
	for _, v := range mediaAttachmentTokens {
		if len(attachments) > 0 {
			attachments += ","
		}
		attachments += v
	}

	rates, err := s.getPeerMessageRate(peerId)
	if err != nil {
		return err
	}
	rates.incrementRate(fromId, fwdDate, messageText, attachments)

	data, err := json.Marshal(rates)
	if err != nil {
		return fmt.Errorf("cannot marshal new versa of rates to json: %v", err)
	}

	err = s.Update(messageRatePath+"."+strconv.Itoa(peerId), string(data))

	if err != nil {
		return fmt.Errorf("cannot update rates in boltdb: %v", err)
	}
	return nil
}

func (s *boltStorage) getPeerMessageRate(peerId int) (*peerMessageRates, error) {
	s.db.Update(func(tx *bolt.Tx) error {
		_, err := ensureBucket(messageRatePath, tx)
		return err
	})
	data, err := s.Get(messageRatePath + "." + strconv.Itoa(peerId))
	if err != nil {
		return nil, fmt.Errorf("cannot get alter version of rates: %v", err)
	}
	rates := &peerMessageRates{}
	err = json.Unmarshal(data, rates)
	if err != nil {
		*rates = getNewPeerMessageRates(peerId)
	}
	return rates, nil
}

func (s *boltStorage) GetMessageTop(peerId int, fromId int) (map[vgontakte.RaterMessage]int, error) {
	rate, err := s.getPeerMessageRate(peerId)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve rate from boltdb: %v", err)
	}

	urate := rate.getUserRate(fromId)
	result := make(map[vgontakte.RaterMessage]int)
	for _, v := range urate.Messages {
		result[v] = v.Rate
	}
	return result, nil
}

func (s *boltStorage) Iterate(fromPath string, rule func(k, v []byte) error) error {
	err := s.db.View(func(tx *bolt.Tx) error {
		var err error
		err = s.iterate(fromPath, rule, tx)
		return err
	})
	return err
}

func (s *boltStorage) iterate(fromPath string, rule func(k, v []byte) error, tx *bolt.Tx) error {
	paths := getPath(fromPath)

	if len(paths) > 0 {
		if len(paths) == 1 {
			return fmt.Errorf("invalid path: %v", fromPath)
		} else {
			bucket, err := s.ensurePathWrite(paths, tx)
			if err != nil {
				return err
			}
			//result := bucket.Get([]byte(paths[len(paths)-1]))

			bucket.ForEach(rule)

			return nil
		}
	}
	return fmt.Errorf("void path")
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
		if tbuck == nil {
			return nil, fmt.Errorf("bucket %v not exists", tpath)
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
