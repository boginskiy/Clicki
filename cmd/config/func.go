package config

import (
	"encoding/json"
	"log"
	"os"
)

var ArgsJSONDefault = &ArgsJSON{
	ServerAddress: "localhost:8080",
	PathToStore:   "",
	EnableHTTPS:   "0",
	AuditFile:     "audit.json",
	AuditURL:      "",
	BaseURL:       "http://localhost:8080",
	DB:            "",
}

func ReadConfigFile(name string, args *ArgsJSON) *ArgsJSON {
	dataByte, err := os.ReadFile(name)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(dataByte, &args)
	return args
}

func CreateConfigFile(name string, args *ArgsJSON) {
	dataByte, err := json.MarshalIndent(args, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}

	f.Write(dataByte)
}
