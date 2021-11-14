package server

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
)

// TestInitCache test sending and retrieve items.
func TestInitCache(t *testing.T) {
	type account struct {
		id   string
		name string
	}

	c := InitCache()

	items := []account{
		{uuid.New().String(), "Jones"},
		{uuid.New().String(), "Luiz"},
	}

	// Send to cache.
	for _, item := range items {
		c.Set(fmt.Sprintf("item_id:%s", item.id), item, 0)
	}

	// Retrieve from cache.
	for _, item := range items {
		i, _ := c.Get(fmt.Sprintf("item_id:%s", item.id))
		name := i.(account).name
		if name != item.name {
			t.Errorf("name of account was incorrect, got: %s, want: %s.", name, item.name)
		}
	}
}
