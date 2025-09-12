package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
}

// "postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable"
// "postgres://username:userpassword@localhost:5432/clickidb?sslmode=disable"

func NewConfig() *Config {
	return &Config{
		DBHost:     "postgres",
		DBPort:     5432,
		DBUser:     "postgres",
		DBPassword: "postgres",
		DBName:     "praktikum",
	}
}

func (c *Config) GetDBConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)
}

func OpenDatabase(c *Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", c.GetDBConnectionString())
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		log.Println(">>OpenDb", err)
		return nil, err
	}
	return db, nil
}
