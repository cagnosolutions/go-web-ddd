package data

// https://www.tutorialspoint.com/sqlite/

type DataRepository interface {

	//Get(id int) (interface{}, error)
	//GetAll() ([]interface{}, error)

	Get(id int, v interface{}) error
	GetAll(vv ...[]interface{}) error

	Close() error
}
