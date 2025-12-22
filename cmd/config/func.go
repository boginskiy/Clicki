package config

import (
	"encoding/json"
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

func ReadConfigFile(name string, args *ArgsJSON) (*ArgsJSON, error) {
	dataByte, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dataByte, &args)
	if err != nil {
		return nil, err
	}
	return args, nil
}

func CreateConfigFile(name string, args *ArgsJSON) error {
	dataByte, err := json.MarshalIndent(args, "", "    ")
	if err != nil {
		return err
	}

	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	f.Write(dataByte)
	return nil
}
