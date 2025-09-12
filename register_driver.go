package main

import (
	"fmt"
	"os"
	"os/exec"

	// Импортируем драйвер libpq для PostgreSQL
	_ "github.com/lib/pq"
)

func main() {
	cmd := exec.Command("migrate", "-source", "file://./migrations", "-database", os.Getenv("DATABASE_CONN_STRING"), "up")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		panic(err)
	}
	fmt.Println("Migrations applied successfully.")
}
