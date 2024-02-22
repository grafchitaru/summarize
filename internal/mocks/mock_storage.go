package mocks

type MockStorage struct {
	PingError error
}

func (ms *MockStorage) Ping() error {
	return ms.PingError
}

func (ms *MockStorage) Close() {

}
