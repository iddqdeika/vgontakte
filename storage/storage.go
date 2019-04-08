package storage

type Storage interface {
	Update(path string, value string) error
	Get(path string) ([]byte, error)
}
