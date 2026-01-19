package config

import "log"

type ArgsJSON struct {
	ServerAddress string `json:"server_address"`
	TrustedSubnet string `json:"trusted_subnet"`
	PathToStore   string `json:"file_storage_path"`
	EnableHTTPS   string `json:"enable_https"`
	EnableGRPC    string `json:"enable_grps"`
	AuditFile     string `json:"audit_file"`
	AuditURL      string `json:"audit_url"`
	BaseURL       string `json:"base_url"`
	DB            string `json:"database_dsn"`
}

func NewArgsJSON(config string) *ArgsJSON {
	if config == "" {
		return nil
	}
	args, err := ReadConfigFile(config, &ArgsJSON{})
	if err != nil {
		log.Fatal(err)
	}
	return args
}

func (aj *ArgsJSON) GetSrvAddr() (ServerAddress string) {
	return aj.ServerAddress
}

func (aj *ArgsJSON) GetBaseURL() (BaseURL string) {
	return aj.BaseURL
}

func (aj *ArgsJSON) GetPathToStore() (PathToStore string) {
	return aj.PathToStore
}

func (aj *ArgsJSON) GetDB() (DB string) {
	return aj.DB
}

func (aj *ArgsJSON) GetAuditFile() (AuditFile string) {
	return aj.AuditFile
}

func (aj *ArgsJSON) GetAuditURL() (AuditURL string) {
	return aj.AuditURL
}

func (aj *ArgsJSON) GetEnableHTTPS() (EnableHTTPS string) {
	return aj.EnableHTTPS
}

func (aj *ArgsJSON) GetTrustedSubnet() (TrustedSubnet string) {
	return aj.TrustedSubnet
}

func (aj *ArgsJSON) GetEnableGRPC() (EnableGRPC string) {
	return aj.EnableGRPC
}
