package db

// Storage описывает операции над хранилищем данных
type Storage interface {
	GetValue(key string) (value string, err error)
	PutValue(key string, value string)
}

// FileWorker описывает работу с файлом
type FileWorker interface {
	DataRecovery(record *StoreModel) map[string]string
	Save(obj any) error
	GetNextLine() int
	Close() error
}
