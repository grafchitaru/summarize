package storage

type Repositories interface {
	Ping() error
	Close()
}
