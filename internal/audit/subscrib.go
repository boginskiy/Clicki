package audit

import (
	"encoding/json"
	"os"

	"github.com/boginskiy/Clicki/internal/logg"
)

// FileReceiver
type FileReceiver struct {
	Logger logg.Logger
	F      *os.File
	id     int
}

func NewFileReceiver(logger logg.Logger, pathToFile string, id int) *FileReceiver {
	file, err := os.OpenFile(pathToFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		logger.RaiseError(err, "FileReceiver>NewFileReceiver>OpenFile", nil)
		id = 0
	}
	return &FileReceiver{
		Logger: logger,
		F:      file,
		id:     id,
	}
}

func (fR *FileReceiver) serialization(event any) ([]byte, error) {
	return json.MarshalIndent(event, "", "\t")
}

func (fR *FileReceiver) GetID() int {
	return fR.id
}

func (fR *FileReceiver) Clouse() {
	fR.F.Close()
}

func (fR *FileReceiver) Update(event any) {
	go func() {
		// Serialization
		jsonData, err := fR.serialization(event)
		if err != nil {
			fR.Logger.RaiseError(err, "FileReceiver>Update>serialization", nil)
			return
		}
		// Запись в файл
		_, err = fR.F.Write(jsonData)
		if err != nil {
			fR.Logger.RaiseError(err, "FileReceiver>Update>Write", nil)
			return
		}
	}()
}

// ServerReceiver
type ServerReceiver struct {
	Logger logg.Logger
	URL    string
	id     int
}

func NewServerReceiver(logger logg.Logger, url string, id int) *ServerReceiver {
	if url == "" {
		id = 0
	}
	return &ServerReceiver{
		Logger: logger,
		URL:    url,
		id:     id,
	}
}

func (sR *ServerReceiver) serialization(event any) ([]byte, error) {
	return json.MarshalIndent(event, "", "\t")
}

func (sR *ServerReceiver) GetID() int {
	return sR.id
}

func (sR *ServerReceiver) Clouse() {

}

func (sR *ServerReceiver) Update(event any) {
	// go func() {
	// 	// Serialization
	// 	jsonData, err := sR.serialization(event)
	// 	if err != nil {
	// 		sR.Logger.RaiseError(err, "ServerReceiver>Update>serialization", nil)
	// 		return
	// 	}

	// }()
}

// TODO!!!  Реализация ServerReceiver!
