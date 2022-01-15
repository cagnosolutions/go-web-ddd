package memory

import (
	"errors"
	"github.com/cagnosolutions/go-web-ddd/pkg/webapp"
	"sync"
)

// MemoryDataSource is an in memory data source
// that implements the lower level DAO interface
type MemoryDataSource struct {
	data *sync.Map
	aid  *webapp.AutoID
}

func NewMemoryDataSource() *MemoryDataSource {
	mds := &MemoryDataSource{
		data: new(sync.Map),
		aid:  new(webapp.AutoID),
	}
	mds.aid.ID()
	return mds
}

// Add helps satisfy the DAO interface
// Add must only add to the underlying storage if it does not exist
func (m *MemoryDataSource) Add(e webapp.Entity) (int, error) {
	// get the id from the entity
	id := e.GetID()
	// new entry
	_, found := m.data.Load(id)
	if !found || id == 0 {
		e.SetID(m.aid.ID())
		m.data.Store(e.GetID(), e)
		return e.GetID(), nil
	}
	// otherwise, entry exists
	return 0, errors.New("entry already exists, cannot add")
}

// Get helps satisfy the DAO interface
func (m *MemoryDataSource) Get(id int) (webapp.Entity, error) {
	// attempt to get the entry by id
	v, found := m.data.Load(id)
	if !found || v == nil {
		return nil, errors.New("entry not found")
	}
	// found it, attempt cast
	e, ok := v.(webapp.Entity)
	if !ok {
		return nil, errors.New("conversion error")
	}
	return e, nil
}

// GetAll helps satisfy the DAO interface
func (m *MemoryDataSource) GetAll() ([]webapp.Entity, error) {
	// init vars
	var err error
	var ee []webapp.Entity
	// round up all the entries
	m.data.Range(func(k, v interface{}) bool {
		e, ok := v.(webapp.Entity)
		if !ok {
			err = errors.New("conversion error")
			return false
		}
		ee = append(ee, e)
		return true
	})
	return ee, err
}

// Set helps satisfy the DAO interface
func (m *MemoryDataSource) Set(e webapp.Entity) error {
	// get the id from the entity
	id := e.GetID()
	// new entry, so update id
	_, found := m.data.Load(id)
	if !found || id == 0 {
		e.SetID(m.aid.ID())
	}
	// new or existing entry, so let us
	// add or update the entry in the
	// underlying storage and we're good
	m.data.Store(e.GetID(), e)
	return nil
}

// Del helps satisfy the DAO interface
func (m *MemoryDataSource) Del(id int) error {
	// delete the entity by id
	m.data.Delete(id)
	return nil
}
