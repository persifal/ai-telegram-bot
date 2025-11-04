package repository

import (
	"log/slog"
	"sync"
)

type DialogStorage interface {
	Exists(int64) bool
	Create(int64)
	Remove(int64)
	Get(int64) (any, bool)
	Clear(int64)
}

type InMemoryStorage struct {
	sync.RWMutex
	m map[int64]string
}

func NewInMemoryDialogStorage() DialogStorage {
	slog.Info("in-memory storage created")
	return &InMemoryStorage{
		m: make(map[int64]string),
	}
}

func (s *InMemoryStorage) Clear(chatId int64) {
	panic("unimplemented")
}

func (s *InMemoryStorage) Get(id int64) (any, bool) {
	panic("unimplemented")
}

func (s *InMemoryStorage) Remove(id int64) {
	panic("unimplemented")
}

func (s *InMemoryStorage) Create(id int64) {
	panic("unimplemented")
}

func (s *InMemoryStorage) Exists(id int64) bool {
	panic("unimplemented")
}
