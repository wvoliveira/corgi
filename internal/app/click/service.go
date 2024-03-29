package click

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
)

type Service interface {
	Find(*gin.Context, string) ([]string, error)

	NewHTTP(*gin.RouterGroup)
	HTTPFind(*gin.Context)
}

type service struct {
	db    *sql.DB
	cache *redis.Client
}

// NewService creates a new public service.
func NewService(db *sql.DB, cache *redis.Client) Service {
	return service{db, cache}
}

// FindAll get a list of links from database.
func (s service) Find(c *gin.Context, link string) (clicks []string, err error) {

	log := logger.Logger(c)

	// key := fmt.Sprintf("link_%s_click_", link)

	// err = s.cache.View(func(txn *badger.Txn) error {

	// 	it := txn.NewIterator(badger.DefaultIteratorOptions)
	// 	prefix := []byte(key)

	// 	defer it.Close()

	// 	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {

	// 		item := it.Item()
	// 		k := item.Key()

	// 		err := item.Value(func(v []byte) error {

	// 			ts := strings.Split(string(k), key)[1]

	// 			clicks = append(clicks, ts)

	// 			return nil
	// 		})

	// 		if err != nil {
	// 			return err
	// 		}

	// 	}

	// 	return nil
	// })

	if err != nil {
		log.Error().Caller().Msg(err.Error())
	}

	return
}
