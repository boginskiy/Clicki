package audit

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/boginskiy/Clicki/internal/logg"
)

// FileReceiver
type FileReceiver struct {
	Logger logg.Logger
	F      *os.File
	ID     int
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
		ID:     id,
	}
}

func (fR *FileReceiver) serialization(event any) ([]byte, error) {
	return json.MarshalIndent(event, "", "\t")
}

func (fR *FileReceiver) GetID() int {
	return fR.ID
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
	ID     int
	ctx    context.Context
	cancel context.CancelFunc
	client *http.Client
}

func NewServerReceiver(logger logg.Logger, url string, id int) *ServerReceiver {
	if url == "" {
		id = 0
	}
	// Context
	ctx, cancel := context.WithCancel(context.Background())

	return &ServerReceiver{
		Logger: logger,
		URL:    url,
		ID:     id,
		ctx:    ctx,
		cancel: cancel,
		client: &http.Client{},
	}
}

func (sR *ServerReceiver) serialization(event any) ([]byte, error) {
	return json.MarshalIndent(event, "", "\t")
}

func (sR *ServerReceiver) GetID() int {
	return sR.ID
}

func (sR *ServerReceiver) Clouse() {
	sR.cancel()
}

func (sR *ServerReceiver) Update(event any) {

	go func(ctx context.Context) {
		// Serializ
		jsonByte, err := sR.serialization(event)
		if err != nil {
			sR.Logger.RaiseError(err, "bad serialization", nil)
			return
		}

		// Request
		req, err := http.NewRequestWithContext(ctx, "POST", sR.URL, bytes.NewReader(jsonByte))
		if err != nil {
			sR.Logger.RaiseError(err, "bad prepar request", nil)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		// Client
		res, err := sR.client.Do(req)
		if err != nil {
			sR.Logger.RaiseError(err, "bad request", nil)
			return
		}

		sR.Logger.RaiseInfo("ServerReceiver: response", logg.Fields{"statusCode": res.StatusCode})
		defer res.Body.Close()

	}(sR.ctx)
}
