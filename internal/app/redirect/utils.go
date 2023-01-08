package redirect

import (
	"context"
	"fmt"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
)

func increaseClick(c context.Context, db *badger.DB, link string, t time.Time) {
	log := logger.Logger(c)

	err := db.Update(func(txn *badger.Txn) (err error) {

		ts := t.UnixMilli()

		// Ex.: link_localhost:8081/globo_click_1257894000
		key := fmt.Sprintf("link_%s_click_%d", link, ts)

		err = txn.Set([]byte(key), []byte{})

		return err
	})

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}
}
