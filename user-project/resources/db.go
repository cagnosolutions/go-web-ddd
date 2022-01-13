package resources

type DB interface {
	Get(id int) (interface{}, error)
	GetAll() ([]interface{}, error)

	Close() error
}
