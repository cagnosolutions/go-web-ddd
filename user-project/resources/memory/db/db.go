package db

import (
	"errors"
	"sync"
)

type MemoryDB struct {
	lock sync.RWMutex
	data map[int]interface{}
}

func OpenMemoryDB() (*MemoryDB, error) {
	return &MemoryDB{
		data: make(map[int]interface{}, 64),
	}, nil
}

func (m *MemoryDB) Get(id int) (interface{}, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	v, ok := m.data[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return v, nil
}

func (m *MemoryDB) GetAll() ([]interface{}, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	var records []interface{}
	for _, v := range m.data {
		records = append(records, v)
	}
	if records == nil {
		return nil, errors.New("empty")
	}
	return records, nil
}

func (m *MemoryDB) Close() error {
	m.data = nil
	return nil
}
