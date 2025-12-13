package staticchecker

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

func readFile(path string) []byte {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}

	dataByte, err := io.ReadAll(file)
	if err != nil {
		log.Fatalln(err)
	}
	return dataByte
}

func deserialization(data []byte, st any) {
	err := json.Unmarshal(data, st)
	if err != nil {
		log.Fatalln(err)
	}
}
