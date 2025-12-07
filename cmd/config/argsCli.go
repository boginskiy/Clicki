package config

import "flag"

// ArgsCLI - struct for args comm line interface.
// generate:reset
type ArgsCLI struct {
	ServerAddress string // StartPort is the port for start application.
	PathToStore   string // PathToStore is the path to store URL.
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

func (c *ArgsCLI) ParseFlags() {
	// defaultStoreDB := "postgres://username:userpassword@localhost:5432/clickidb?sslmode=disable"
	// AuditFile - "./audit.json"
	// AuditURL -  "http://localhost:8081/"

	flag.StringVar(&c.BaseURL, "b", "http://localhost:8080", "Result adress for application")
	flag.StringVar(&c.ServerAddress, "a", "localhost:8080", "Start adress for application")
	flag.StringVar(&c.AuditFile, "audit-file", "", "Path to the audit file")
	flag.StringVar(&c.AuditURL, "audit-url", "", "URL to the audit server")
	flag.StringVar(&c.PathToStore, "f", "", "Path to file of store URL")
	flag.StringVar(&c.DB, "d", "", "Data of connected DB")

	flag.Parse()
}

func (c *ArgsCLI) GetSrvAddr() (ServerAddress string) {
	return c.ServerAddress
}

func (c *ArgsCLI) GetBaseURL() (BaseURL string) {
	return c.BaseURL
}

func (c *ArgsCLI) GetPathToStore() (PathToStore string) {
	return c.PathToStore
}

func (c *ArgsCLI) GetDB() (DB string) {
	return c.DB
}

func (c *ArgsCLI) GetAuditFile() (AuditFile string) {
	return c.AuditFile
}

func (c *ArgsCLI) GetAuditURL() (AuditURL string) {
	return c.AuditURL
}
