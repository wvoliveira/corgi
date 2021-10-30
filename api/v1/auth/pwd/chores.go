package pwd

import (
	"time"

	l "github.com/go-kit/log"
)

func (rs *Resource) choresTicker() {
	var logger l.Logger

	ticker := time.NewTicker(time.Hour * 1)
	go func() {
		for range ticker.C {
			if err := rs.Store.PurgeExpiredToken(); err != nil {
				logger.Log("method", "choresTicker", "chore", "purgeExpiredToken", "err", err)
			}
		}
	}()
}
