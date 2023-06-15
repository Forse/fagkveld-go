package db

import (
	"sync"
	"sync/atomic"
)

type Todo struct {
	ID          uint64
	Title       string
	Description string
	IsCompleted bool
}

var sequence uint64 = 0

type DB struct {
	db []Todo
	mu sync.Mutex
}

func NewDB() *DB {
	return &DB{
		db: make([]Todo, 0, 64),
	}
}

func (db *DB) append(todo *Todo) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.db = append(db.db, *todo)
}

func (db *DB) Create(todo *Todo) {
	todo.ID = atomic.AddUint64(&sequence, 1)

	db.append(todo)
}
