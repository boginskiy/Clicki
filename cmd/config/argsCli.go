package config

import "flag"

// ArgsCLI - struct for args comm line interface.
// generate:reset
type ArgsCLI struct {
	ServerAddress string // StartPort is the port for start application.
	PathToStore   string // PathToStore is the path to store URL.
	EnableHTTPS   string // Turn of/on HTTPS
	AuditFile     string // AuditFile is for turn on a file.
	AuditURL      string // AuditURL is for turn on a server.
	BaseURL       string // BaseURL is result port is the port after changing.
	DB            string // DB is data of connected.
}

func NewArgsCLI() *ArgsCLI {
	args := new(ArgsCLI)
	args.ParseFlags()
	return args
}

func (ac *ArgsCLI) ParseFlags() {
	// defaultStoreDB := "postgres://username:userpassword@localhost:5432/clickidb?sslmode=disable"
	// AuditFile - "./audit.json"
	// AuditURL -  "http://localhost:8081/"

	flag.StringVar(&ac.BaseURL, "b", "http://localhost:8080", "Result adress for application")
	flag.StringVar(&ac.ServerAddress, "a", "localhost:8080", "Start adress for application")
	flag.StringVar(&ac.AuditFile, "audit-file", "", "Path to the audit file")
	flag.StringVar(&ac.AuditURL, "audit-url", "", "URL to the audit server")
	flag.StringVar(&ac.PathToStore, "f", "", "Path to file of store URL")
	flag.StringVar(&ac.EnableHTTPS, "s", "", "Turn of/on HTTPS protocol")
	flag.StringVar(&ac.DB, "d", "", "Data of connected DB")

	flag.Parse()
}

func (ac *ArgsCLI) GetSrvAddr() (ServerAddress string) {
	return ac.ServerAddress
}

func (ac *ArgsCLI) GetBaseURL() (BaseURL string) {
	return ac.BaseURL
}

func (ac *ArgsCLI) GetPathToStore() (PathToStore string) {
	return ac.PathToStore
}

func (ac *ArgsCLI) GetDB() (DB string) {
	return ac.DB
}

func (ac *ArgsCLI) GetAuditFile() (AuditFile string) {
	return ac.AuditFile
}

func (ac *ArgsCLI) GetAuditURL() (AuditURL string) {
	return ac.AuditURL
}

func (ac *ArgsCLI) GetEnableHTTPS() (EnableHTTPS string) {
	return ac.EnableHTTPS
}
