package bolt

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/maritimusj/centrum/gate/logStore"
	"sync"
	"time"

	bolt "github.com/etcd-io/bbolt"
	"github.com/sirupsen/logrus"
)

type store struct {
	db    *bolt.DB
	cache chan *logStore.Entry

	entryPool   *sync.Pool
	encoderPool *sync.Pool

	done chan struct{}

	wg sync.WaitGroup
	mu sync.RWMutex
}

func New() logStore.Store {
	return &store{
		entryPool: &sync.Pool{
			New: func() interface{} {
				return make(map[string]interface{})
			},
		},
		encoderPool: &sync.Pool{
			New: func() interface{} {
				return NewCopier()
			},
		},
		done: make(chan struct{}),
	}
}

// i2b returns an 8-byte big endian representation of v.
func i2b(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

func b2i(v []byte) uint64 {
	return binary.BigEndian.Uint64(v)
}

func (store *store) Open(option map[string]interface{}) error {
	filename, _ := option["filename"].(string)
	if filename == "" {
		return errors.New("invalid log filename")
	}

	store.cache = make(chan *logStore.Entry, 1000)

	db, err := bolt.Open(filename, 0666, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		return err
	}

	store.db = db

	store.wg.Add(1)
	go store.worker(db, store.cache)

	return nil
}

func (store *store) getEncoder() *gobCopier {
	return store.encoderPool.Get().(*gobCopier)
}

func (store *store) releaseEncoder(encoder *gobCopier) {
	store.encoderPool.Put(encoder)
}

func (store *store) Close() {
	close(store.done)

	store.mu.Lock()
	{
		close(store.cache)
		store.cache = nil
	}
	store.mu.Unlock()

	store.wg.Wait()
}

func (store *store) Delete(orgID int64, src string) error {
	store.mu.RLock()
	defer store.mu.RUnlock()

	if src == "" {
		src = logStore.SystemLog
	}

	if store.db == nil {
		return errors.New("db closed")
	}

	return store.db.Update(func(tx *bolt.Tx) error {
		orgB := tx.Bucket(i2b(uint64(orgID)))
		if orgB == nil {
			return nil
		}
		logB := orgB.Bucket([]byte("log"))
		if logB != nil {
			err := logB.DeleteBucket([]byte(src))
			if err == bolt.ErrBucketNotFound {
				return nil
			}
			return err
		}
		return nil
	})
}

func (store *store) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
		logrus.TraceLevel,
	}
}

func (store *store) Fire(entry *logrus.Entry) error {
	store.mu.RLock()
	defer store.mu.RUnlock()

	if store.cache == nil {
		return nil
	}

	encoder := store.getEncoder()
	defer store.releaseEncoder(encoder)

	err := encoder.encode(entry.Data)
	if err != nil {
		return err
	}

	e := map[string]interface{}{}
	err = encoder.decode(&e)
	if err != nil {
		return err
	}

	store.wg.Add(1)
	go func() {
		defer store.wg.Done()
		select {
		case <-store.done:
			return
		default:
			store.mu.RLock()
			if store.cache != nil {
				store.cache <- &logStore.Entry{
					Level:     entry.Level.String(),
					Message:   entry.Message,
					Fields:    e,
					CreatedAt: entry.Time,
				}
			}
			store.mu.RUnlock()
		}
	}()

	return nil
}

func (store *store) worker(db *bolt.DB, cache <-chan *logStore.Entry) {
	defer func() {
		store.wg.Done()
	}()

	write := func(entry *logStore.Entry) {
		if err := store.write(db, entry); err != nil {
			fmt.Println("fail to write log to db:", err)
		}
	}

	for {
		select {
		case <-store.done:
			if len(store.done) > 0 {
				fmt.Print("write log to db")
				for entry := range cache {
					fmt.Print(".")
					write(entry)
				}
				fmt.Println("Ok")
			}
			err := db.Close()
			if err != nil {
				fmt.Println("fail to close log db:", err)
			}
			return
		case entry := <-cache:
			if entry != nil {
				write(entry)
			}
		}
	}
}

func (store *store) write(db *bolt.DB, entry *logStore.Entry) error {
	var orgID uint64
	if v, ok := entry.Fields["org"].(int64); ok {
		orgID = uint64(v)
		delete(entry.Fields, "org")
	}

	var src string
	if v, ok := entry.Fields["src"].(string); ok {
		src = v
		delete(entry.Fields, "src")
	}

	if src == "" {
		src = logStore.SystemLog
	}

	return db.Batch(func(tx *bolt.Tx) error {
		orgB, err := tx.CreateBucketIfNotExists(i2b(orgID))
		if err != nil {
			return err
		}

		logB, err := orgB.CreateBucketIfNotExists([]byte("log"))
		if err != nil {
			return err
		}

		srcB, err := logB.CreateBucketIfNotExists([]byte(src))
		if err != nil {
			return err
		}

		entriesB, err := srcB.CreateBucketIfNotExists([]byte("entries"))
		if err != nil {
			return err
		}

		data, err := entry.Marshal()
		if err != nil {
			return err
		}

		entryID, _ := entriesB.NextSequence()
		err = entriesB.Put(i2b(entryID), data)
		if err != nil {
			return err
		}

		levelB, err := srcB.CreateBucketIfNotExists([]byte(entry.Level))
		if err != nil {
			return err
		}

		levelID, _ := levelB.NextSequence()
		err = levelB.Put(i2b(levelID), i2b(entryID))
		if err != nil {
			return err
		}

		return nil
	})
}

func (store *store) Stats(orgID int64) map[string]interface{} {
	store.mu.RLock()
	defer store.mu.RUnlock()

	var stats = map[string]uint64{}
	_ = store.db.View(func(tx *bolt.Tx) error {
		orgB := tx.Bucket(i2b(uint64(orgID)))
		if orgB == nil {
			return nil
		}

		logB := tx.Bucket([]byte("log"))
		if logB == nil {
			return nil
		}
		_ = logB.ForEach(func(k, v []byte) error {
			if v == nil {
				srcB := logB.Bucket(k)
				if srcB != nil {
					_ = srcB.ForEach(func(k, v []byte) error {
						if v == nil {
							b := srcB.Bucket(k)
							if b != nil {
								stats[string(k)] += b.Sequence()
							}
						}
						return nil
					})
				}
			}
			return nil
		})
		return nil
	})

	var result = map[string]interface{}{}
	for k, v := range stats {
		result[k] = v
	}
	return result
}

func (store *store) GetList(orgID int64, src, level string, start *uint64, offset, limit uint64) (result []*logStore.Data, total uint64, err error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	if store.db == nil {
		result = []*logStore.Data{}
		err = nil
		return
	}

	if start == nil {
		start = new(uint64)
	}

	if src == "" {
		src = logStore.SystemLog
	}

	var errNoResult = errors.New("no result")
	err = store.db.View(func(tx *bolt.Tx) error {
		orgB := tx.Bucket(i2b(uint64(orgID)))
		if orgB == nil {
			return errNoResult
		}
		logB := orgB.Bucket([]byte("log"))
		if logB == nil {
			return errNoResult
		}

		srcB := logB.Bucket([]byte(src))
		if srcB == nil {
			return errNoResult
		}

		entriesB := srcB.Bucket([]byte("entries"))
		if entriesB == nil {
			return errNoResult
		}

		if level != "" {
			logIDs := make([][]byte, 0, limit)
			levelB := srcB.Bucket([]byte(level))
			if levelB != nil {
				total = levelB.Sequence()
				if *start == 0 {
					*start = total
				} else if *start > total {
					*start = total
				}
				c := levelB.Cursor()
				prefix := i2b(*start - offset)
				l := uint64(0)
				for k, v := c.Seek(prefix); k != nil && l < limit; k, v = c.Prev() {
					logIDs = append(logIDs, v)
					l++
				}
			}
			for _, id := range logIDs {
				v := entriesB.Get(id)
				result = append(result, &logStore.Data{ID: b2i(id), Content: v})
			}
		} else {
			total = entriesB.Sequence()
			if *start == 0 {
				*start = total
			} else if *start > total {
				*start = total
			}

			c := entriesB.Cursor()
			prefix := i2b(*start - offset)
			l := uint64(0)
			for k, v := c.Seek(prefix); k != nil && l < limit; k, v = c.Prev() {
				result = append(result, &logStore.Data{ID: b2i(k), Content: v})
				l++
			}
		}

		return nil
	})

	if err != nil {
		result = []*logStore.Data{}
		err = nil
	}

	return
}
