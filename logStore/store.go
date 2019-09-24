package logStore

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	bolt "github.com/etcd-io/bbolt"
	"github.com/sirupsen/logrus"
)

const (
	SystemLog = "system"
)

type Data struct {
	ID      uint64 `json:"id"`
	Content []byte `json:"content"`
}

type Entry struct {
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields"`
	CreatedAt time.Time              `json:"created_at"`
}

func (entry *Entry) Marshal() ([]byte, error) {
	return json.Marshal(entry)
}

type JsonCopier struct {
	encoder *json.Encoder
	decoder *json.Decoder
}

func NewJsonCopier() *JsonCopier {
	buffer := bytes.Buffer{}
	return &JsonCopier{
		encoder: json.NewEncoder(&buffer),
		decoder: json.NewDecoder(&buffer),
	}
}

func (c *JsonCopier) encode(v interface{}) error {
	return c.encoder.Encode(v)
}

func (c *JsonCopier) decode(data interface{}) error {
	return c.decoder.Decode(data)
}

type Store interface {
	Open(ctx context.Context, filename string) error
	Wait()

	Get(orgID int64, src, level string, start *uint64, offset, limit uint64) (result []*Data, total uint64, err error)
	Delete(orgID int64, src string) error
	Stats(orgID int64) map[string]interface{}

	//interface for logrus hook
	Levels() []logrus.Level
	Fire(entry *logrus.Entry) error
}

type logStore struct {
	db    *bolt.DB
	cache chan *Entry

	routine context.CancelFunc

	entryPool   *sync.Pool
	encoderPool *sync.Pool

	wg sync.WaitGroup
	mu sync.RWMutex
}

func New() Store {
	return &logStore{
		entryPool: &sync.Pool{
			New: func() interface{} {
				return make(map[string]interface{})
			},
		},
		encoderPool: &sync.Pool{
			New: func() interface{} {
				return NewJsonCopier()
			},
		},
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

func (store *logStore) Wait() {
	store.wg.Wait()
}

func (store *logStore) Open(ctx context.Context, filename string) error {
	store.Close()

	store.mu.Lock()
	defer store.mu.Unlock()

	store.cache = make(chan *Entry, 1000)

	db, err := bolt.Open(filename, 0666, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		return err
	}

	db.Stats()

	store.db = db

	routineCtx, cancel := context.WithCancel(ctx)
	store.routine = cancel

	store.wg.Add(2)
	go store.worker(routineCtx, db, store.cache)

	go func() {
		defer store.wg.Done()
		select {
		case <-ctx.Done():
			store.Close()
		}
	}()

	return nil
}

func (store *logStore) getEncoder() *JsonCopier {
	return store.encoderPool.Get().(*JsonCopier)
}

func (store *logStore) releaseEncoder(encoder *JsonCopier) {
	store.encoderPool.Put(encoder)
}

func (store *logStore) getEntry() map[string]interface{} {
	return store.entryPool.Get().(map[string]interface{})
}

func (store *logStore) releaseEntry(entry map[string]interface{}) {
	store.entryPool.Put(entry)
}

func (store *logStore) Close() {
	store.mu.Lock()
	defer store.mu.Unlock()

	if store.cache != nil {
		close(store.cache)
		store.cache = nil
	}

	if store.routine != nil {
		store.routine()
		store.routine = nil
	}
}

func (store *logStore) Delete(orgID int64, src string) error {
	store.mu.RLock()
	defer store.mu.RUnlock()

	if src == "" {
		src = SystemLog
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

func (store *logStore) Levels() []logrus.Level {
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

func (store *logStore) Fire(entry *logrus.Entry) error {
	store.mu.RLock()
	defer store.mu.RUnlock()

	if store.cache != nil {
		encoder := store.getEncoder()
		defer store.releaseEncoder(encoder)

		err := encoder.encode(entry.Data)
		if err != nil {
			return err
		}

		e := store.getEntry()
		err = encoder.decode(&e)
		if err != nil {
			store.releaseEntry(e)
			return err
		}

		store.cache <- &Entry{
			Level:     entry.Level.String(),
			Message:   entry.Message,
			Fields:    e,
			CreatedAt: entry.Time,
		}
	}

	return nil
}

func (store *logStore) worker(ctx context.Context, db *bolt.DB, cache <-chan *Entry) {
	defer func() {
		store.wg.Done()
	}()

	write := func(entry *Entry) {
		if err := store.write(db, entry); err != nil {
			fmt.Println("fail to write log to db:", err)
		}

		store.releaseEntry(entry.Fields)
	}

	for {
		select {
		case <-ctx.Done():
			for entry := range cache {
				write(entry)
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

func (store *logStore) write(db *bolt.DB, entry *Entry) error {
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
		src = SystemLog
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
		err = entriesB.Put(i2b(entryID), []byte(data))
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

func (store *logStore) Stats(orgID int64) map[string]interface{} {
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

func (store *logStore) Get(orgID int64, src, level string, start *uint64, offset, limit uint64) (result []*Data, total uint64, err error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	if store.db == nil {
		result = []*Data{}
		err = nil
		return
	}

	if start == nil {
		start = new(uint64)
	}

	if src == "" {
		src = SystemLog
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
				result = append(result, &Data{b2i(id), v})
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
				result = append(result, &Data{b2i(k), v})
				l++
			}
		}

		return nil
	})

	if err != nil {
		result = []*Data{}
		err = nil
	}

	return
}
