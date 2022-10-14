package ddb

import (
	"context"
	"sync"
)

func (c *Client) NewTransaction() Transaction {
	return &DBTransaction{
		client: c,
	}
}

// DBTransaction writes pending items to slices in memory.
// It is goroutine-safe.
type DBTransaction struct {
	client Storage
	// mu is a mutex to prevent concurrent writes to the
	// putItems and deleteItems slices.
	mu          sync.Mutex
	putItems    []Keyer
	deleteItems []Keyer
}

func (t *DBTransaction) Put(item Keyer) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.putItems = append(t.putItems, item)
}

func (t *DBTransaction) Delete(item Keyer) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.deleteItems = append(t.deleteItems, item)
}

func (t *DBTransaction) Execute(ctx context.Context) error {
	items := t.buildTransactWriteItemsPayload()
	return t.client.TransactWriteItems(ctx, items)
}

func (t *DBTransaction) buildTransactWriteItemsPayload() []TransactWriteItem {
	items := make([]TransactWriteItem, len(t.putItems)+len(t.deleteItems))
	for i := range t.putItems {
		items[i] = TransactWriteItem{
			Put: t.putItems[i],
		}
	}
	for i := range t.deleteItems {
		items[len(t.putItems)+i] = TransactWriteItem{
			Delete: t.deleteItems[i],
		}
	}
	return items
}
