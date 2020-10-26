package properties

import (
	"bytes"
	"encoding/gob"
	"errors"
	"strings"
	"time"

	bolt "github.com/etcd-io/bbolt"
)

const (
	pathSeparator  = "."
	bucketSuffix   = ".p"
	valueSuffix    = ".v"
	dataBucketName = "data"
)

var (
	marshal = func(v interface{}) ([]byte, error) {
		var b bytes.Buffer
		w := gob.NewEncoder(&b)
		if err := w.Encode(v); err != nil {
			return nil, err
		}
		return b.Bytes(), nil
	}

	unmarshal = func(data []byte, v interface{}) error {
		r := gob.NewDecoder(bytes.NewBuffer(data))
		if err := r.Decode(v); err != nil {
			return err
		}
		return nil
	}

	ErrInvalidPath = errors.New("invalid path")

	defaultStore = New()
)

type Store interface {
	Open(option map[string]interface{}) error
	Close()
	Write(path string, value interface{}) error
	Delete(path string) error
	Load(path string, v interface{}) error
	LoadAllValueString(path string) (map[string]string, error)
}

type StoreImp struct {
	db *bolt.DB
}

func New() Store {
	return &StoreImp{}
}

func Open(filename string) error {
	opts := map[string]interface{}{}
	opts["filename"] = filename
	return defaultStore.Open(opts)
}

func Close() {
	defaultStore.Close()
}

func Write(v interface{}, path ...string) error {
	return defaultStore.Write(strings.Join(path, pathSeparator), v)
}

func Delete(path ...string) error {
	return defaultStore.Delete(strings.Join(path, pathSeparator))
}

func Load(v interface{}, path ...string) error {
	return defaultStore.Load(strings.Join(path, pathSeparator), v)
}

func LoadString(path ...string) string {
	var str string
	if err := Load(&str, path...); err != nil {
		return ""
	}
	return str
}

func LoadAllString(path ...string) map[string]string {
	result, err := defaultStore.LoadAllValueString(strings.Join(path, pathSeparator))
	if err != nil {
		return map[string]string{}
	}
	return result
}

func (store *StoreImp) Open(option map[string]interface{}) error {
	filename, _ := option["filename"].(string)
	if filename == "" {
		return errors.New("invalid filename")
	}

	db, err := bolt.Open(filename, 0666, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}

	store.db = db

	return nil
}

func (store *StoreImp) Close() {
	if store.db != nil {
		_ = store.db.Close()
	}
}

func (store *StoreImp) Delete(path string) error {
	if store != nil && store.db != nil {
		return store.db.Update(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(dataBucketName))
			if bucket == nil {
				return ErrInvalidPath
			}

			paths := strings.Split(path, pathSeparator)
			key := paths[len(paths)-1:][0]
			paths = paths[:len(paths)-1]
			for _, p := range paths {
				if p != "" {
					bucket = bucket.Bucket([]byte(p + bucketSuffix))
					if bucket == nil {
						return ErrInvalidPath
					}
				}
			}

			fullKeyName := []byte(key + valueSuffix)
			v := bucket.Get(fullKeyName)
			if v != nil {
				return bucket.Delete(fullKeyName)
			}

			return bucket.DeleteBucket([]byte(key + bucketSuffix))
		})
	}
	return errors.New("not initialized")
}

func (store *StoreImp) Write(path string, value interface{}) error {
	if store != nil && store.db != nil {
		if path == "" {
			return errors.New("invalid path")
		}

		data, err := marshal(value)
		if err != nil {
			return err
		}

		return store.db.Update(func(tx *bolt.Tx) error {
			bucket, err := tx.CreateBucketIfNotExists([]byte(dataBucketName))
			if err != nil {
				return err
			}
			paths := strings.Split(path, pathSeparator)
			key := paths[len(paths)-1:][0]
			paths = paths[:len(paths)-1]
			for _, p := range paths {
				if p != "" {
					bucket, err = bucket.CreateBucketIfNotExists([]byte(p + bucketSuffix))
					if err != nil {
						return err
					}
				}
			}
			return bucket.Put([]byte(key+valueSuffix), data)
		})
	}
	return errors.New("not initialized")
}

func (store *StoreImp) LoadAllValueString(path string) (map[string]string, error) {
	if store != nil && store.db != nil {
		if path == "" {
			return nil, errors.New("invalid path")
		}
		var (
			result = map[string]string{}
		)

		err := store.db.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(dataBucketName))
			if bucket == nil {
				return errors.New("path not exists")
			}

			for _, p := range strings.Split(path, pathSeparator) {
				if p != "" {
					bucket = bucket.Bucket([]byte(p + bucketSuffix))
					if bucket == nil {
						return errors.New("path not exists")
					}
				}
			}

			return bucket.ForEach(func(k, v []byte) error {
				if v != nil {
					var str string
					if err := unmarshal(v, &str); err == nil {
						result[strings.TrimSuffix(string(k), valueSuffix)] = str
					}
				}
				return nil
			})
		})

		if err != nil {
			return nil, err
		}

		return result, nil
	}

	return nil, errors.New("not initialized")
}

func (store *StoreImp) Load(path string, v interface{}) error {
	if store != nil && store.db != nil {
		if path == "" {
			return errors.New("invalid path")
		}

		var data []byte
		err := store.db.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(dataBucketName))
			if bucket == nil {
				return errors.New("path not exists")
			}

			paths := strings.Split(path, pathSeparator)
			key := paths[len(paths)-1:][0]
			paths = paths[:len(paths)-1]

			for _, p := range paths {
				if p != "" {
					bucket = bucket.Bucket([]byte(p + bucketSuffix))
					if bucket == nil {
						return errors.New("path not exists")
					}
				}
			}

			data = bucket.Get([]byte(key + valueSuffix))
			return nil
		})

		if err != nil {
			return err
		}

		return unmarshal(data, v)
	}
	return errors.New("not initialized")
}
