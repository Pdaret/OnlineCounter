package bolt

import (
	"log"
	"strings"
	"time"
	"x-ui-monitor/internal/domain"

	"github.com/boltdb/bolt"
	bbolt "github.com/boltdb/bolt"
)

const ttlSeconds = 120

type BoltDB struct {
	db         *bolt.DB
	bucketName string
}

func NewBoltDB(path string, bucketName string) (*BoltDB, error) {
	db, err := bbolt.Open(path, 0666, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Ensure bucket exists
	db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		return err
	})

	return &BoltDB{
		db:         db,
		bucketName: bucketName,
	}, nil
}

func (b *BoltDB) AddIP(inbound, ip string) error {
	key := inbound + ":" + ip
	value := []byte(time.Now().Format(time.RFC3339))

	return b.db.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte(b.bucketName)).Put([]byte(key), value)
	})
}

func (b *BoltDB) GetActiveIPCount(inbound string) (int, error) {
	now := time.Now()
	count := 0

	err := b.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(b.bucketName))

		c := bucket.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			key := string(k)
			if strings.HasPrefix(key, inbound+":") {
				// Parse timestamp
				ts, err := time.Parse(time.RFC3339, string(v))
				if err != nil {
					bucket.Delete(k) // Corrupted entry
					continue
				}
				if now.Sub(ts) < ttlSeconds*time.Second {
					count++
				} else {
					bucket.Delete(k) // Expired entry
				}
			}
		}
		return nil
	})

	return count, err
}

func (b *BoltDB) ListActiveIPs(inbound string) ([]string, error) {
	now := time.Now()
	var ips []string

	err := b.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(b.bucketName))

		c := bucket.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			key := string(k)
			if strings.HasPrefix(key, inbound+":") {
				ip := strings.TrimPrefix(key, inbound+":")
				ts, err := time.Parse(time.RFC3339, string(v))
				if err != nil {
					bucket.Delete(k)
					continue
				}
				if now.Sub(ts) < ttlSeconds*time.Second {
					ips = append(ips, ip)
				} else {
					bucket.Delete(k)
				}
			}
		}
		return nil
	})

	return ips, err
}

func (b *BoltDB) GetTotalUsersCount() (*domain.TotalUsersCountResult, error) {
	stats := domain.TotalUsersCountResult{Inbounds: make(map[string]int)}
	now := time.Now()

	b.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(b.bucketName))
		c := bucket.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			key := string(k)
			parts := strings.SplitN(key, ":", 2)
			if len(parts) != 2 {
				continue
			}
			inbound := parts[0]
			ts, err := time.Parse(time.RFC3339, string(v))
			if err != nil || now.Sub(ts) > ttlSeconds*time.Second {
				continue
			}
			stats.Inbounds[inbound]++
			stats.Total++
		}
		return nil
	})
	return &stats, nil
}

func (b *BoltDB) Close() {
	b.db.Close()
}
