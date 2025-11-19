package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func HardWork(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method is not supported", http.StatusBadRequest)
		return
	}

	jsonByte, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "method is not supported", http.StatusBadRequest)
		return
	}

	fmt.Fprintln(os.Stdout, string(jsonByte))

	time.Sleep(500 * time.Microsecond)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"main phrase": "Hello world!"}`))
}

func main() {
	http.HandleFunc("/", HardWork)
	log.Fatal(http.ListenAndServe("localhost:8081", nil))
}
