package db

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

const SIZE = 20

type FileWorking struct {
	file    *os.File
	scanner *bufio.Scanner
	cntLine int
}

func NewFileWorking(filename string) (*FileWorking, error) {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &FileWorking{file: f, scanner: bufio.NewScanner(f)}, nil
}

func (f *FileWorking) DataRecovery(record *StoreModel) map[string]string {
	result := make(map[string]string, SIZE)

	// Проход по строкам
	for f.scanner.Scan() {
		line := f.scanner.Text()

		// Десериализация
		err := json.Unmarshal([]byte(line), &record)
		if err != nil {
			continue
		}
		// Сохранение данных с map
		result[record.ShortURL] = record.OriginalURL
		// Счетчик для UUID
		f.cntLine = max(f.cntLine, record.UUID)
	}
	return result
}

func (f *FileWorking) Close() error {
	return f.file.Close()
}

func (f *FileWorking) Save(record any) error {
	tmpB, err := json.Marshal(record)
	if err != nil {
		return err
	}
	tmpB = append(tmpB, byte('\n'))
	_, err = f.file.Write(tmpB)
	return err
}

func (f *FileWorking) ReadLastRecord(record any) (any, error) {
	// Временная переменная куда итеративно сохряняем байты
	buff := make([]byte, 1)

	// Узнаем размер файла
	info, err := f.file.Stat()
	if err != nil {
		return nil, err
	}
	size := info.Size()

	// Проверка, что файл не пуст
	if size == 0 {
		return nil, fmt.Errorf("file %s is emty", f.file.Name())
	}

	// Перемещаем указатель файла с конца по байтово
	// до первого попавшегося разделителя '\n'

	for i := int64(-1); i >= (-size); i-- {
		buff[0] = 0

		// Байтовое смещение с конца файла
		f.file.Seek(i, io.SeekEnd)
		// Чтение байта
		_, err := f.file.Read(buff)

		if err != nil {
			return nil, err
		}

		if buff[0] == '\n' {
			// Создаем буффер для последней записи
			tmpRecord := make([]byte, -i)

			// Считываем запись
			_, err := f.file.Read(tmpRecord)

			if err != nil {
				return nil, err
			}

			// Unmarshal
			err = json.Unmarshal(tmpRecord, &record)
			if err != nil {
				return nil, err
			}
			return record, nil
		}
	}
	return nil, errors.New("bad work of FileWorking.ReadLastRecord")
}

func (f *FileWorking) GetNextLine() int {
	return f.cntLine + 1
}
