package webapp

import "net/http"

type Entity interface {
	GetID() int
	SetID(id int)
}

type DAO interface {
	Add(e Entity) (int, error)  // add a new entity, return id or error
	Get(id int) (Entity, error) // get an entity by id, return any error
	GetAll() ([]Entity, error)  // get all entities, return number found or error
	Set(e Entity) error         // update an existing entity by id, return any error
	Del(id int) error           // delete an existing entity by id
}

type Repository interface {
	WithDAO(daos ...DAO) error
}

type Service interface {
	WithRepo(repo Repository) error
}

type Controller interface {
	WithService(service Service) error
	RootHandler(w http.ResponseWriter, r *http.Request)
}
