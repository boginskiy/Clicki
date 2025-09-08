package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

var (
	ApplWasStopedInfo = "INFO: the application was successfully stopped"
	TerminalCommErr   = "ERROR: error in terminal command"
	NotPidProcInfo    = "INFO: there is no PID process"
)

/*
Description:
	Testing Appl with real start web-server of Appl
*/

func commandKill(pid string) {
	if pid == "" {
		fmt.Fprintln(os.Stderr, NotPidProcInfo)
		return
	}

	cmd := exec.Command("kill", pid)
	cmd.Stdout = os.Stdout
	err := cmd.Run()

	if err != nil {
		fmt.Fprintln(os.Stderr, TerminalCommErr, err.Error())
		return
	}
	fmt.Fprintln(os.Stdout, ApplWasStopedInfo)
}

func showOpenProc(port string, out *bytes.Buffer) {
	cmd := exec.Command("lsof", "-i", port)

	// out:
	// 	   COMMAND    PID USER   FD   TYPE DEVICE SIZE/OFF NODE NAME
	//     clicki  118361  ali    3u  IPv4 447949      0t0  TCP localhost:http-alt (LISTEN)

	cmd.Stdout = out
	err := cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, TerminalCommErr, err.Error())
		return
	}
}

func takePidProc(out *bytes.Buffer) string {
	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	pid := ""

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			pid = fields[1]
		}
	}
	return pid
}

func StopApplication(port string) {
	var out bytes.Buffer

	showOpenProc(port, &out) // Starting info of Appl
	pid := takePidProc(&out) // Take PID process
	commandKill(pid)         // Kill process
}

func StartApplication(err chan<- error, path string) {
	// Запуск программы для тестирования
	cmd := exec.Command(path)
	// устанавливаем стандартные потоки
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err <- cmd.Run()
}

func CheckResBody(body string, reqular string) bool {
	re := regexp.MustCompile(reqular)
	return re.MatchString(body)
}

type req struct {
	url         string
	method      string
	body        string
	contentType string
}

type res struct {
	statusCode   int
	contentType  string
	checkingBody bool
}

func testRunningCases() {
	tests := []struct {
		name   string
		params req
		want   res
		reg    string
	}{
		{
			name: "Check POST with domen web site",
			params: req{
				url:         "http://localhost:8080",
				method:      http.MethodPost,
				body:        "https://docs.google.com/",
				contentType: "text/plain",
			},
			want: res{
				statusCode:   201,
				contentType:  "text/plain",
				checkingBody: true,
			},
			reg: `^http://localhost:8080/[a-zA-Z0-9]{8}$`,
		},
	}

	// Client
	client := http.Client{
		// Delete func of Redirect
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	for _, tt := range tests {
		fmt.Fprintf(os.Stdout, "\n=== Start test >> %s\n", tt.name)

		// Create new req
		req, err := http.NewRequest(tt.params.method, tt.params.url, strings.NewReader(tt.params.body))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}

		req.Header.Set("Content-Type", tt.params.contentType)

		// Send req
		res, err := client.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}

		// Checking StatusCode ...
		if tt.want.statusCode != res.StatusCode {
			fmt.Fprintf(os.Stderr, "expected: %v != actual: %v\n", tt.want.statusCode, res.StatusCode)
		}

		// Checking Content-Type ...
		if tt.want.contentType != res.Header.Get("Content-Type") {
			fmt.Fprintf(os.Stderr, "expected: %v != actual: %v\n", tt.want.contentType, res.Header.Get("Content-Type"))
		}

		// Checking Body ...
		bodyByte, err := io.ReadAll(res.Body)
		defer res.Body.Close()

		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}

		resultCompare := CheckResBody(string(bodyByte), tt.reg)

		if tt.want.checkingBody != resultCompare {
			fmt.Fprintf(os.Stderr, "expected: %v != actual: %v\n", tt.want.checkingBody, resultCompare)
		}

		fmt.Fprint(os.Stdout, "===================== pass =====================\n\n")
	}
}

func TestApplication(chanErr chan error, path, port string) {
	// Info
	wd, _ := os.Getwd()
	fmt.Fprintf(os.Stdout, "INFO: work catalog: %v\n", wd)

	// Chan
	done := make(chan uint8)

	// Start server
	go StartApplication(chanErr, path)

	go func(err <-chan error, done <-chan uint8) {
		select {
		case er := <-err:
			fmt.Fprintf(os.Stderr, "ERR: The program was not running: %v\n", er)
			return
		case <-done:
			fmt.Fprintln(os.Stdout, "INFO: the gorutina was successfully stopped")
		}
	}(chanErr, done)

	fmt.Fprintln(os.Stdout, "INFO: the application was successfully started")
	time.Sleep(time.Millisecond * 200)

	// Testing
	testRunningCases()

	// Stop server and gorutine
	done <- 0
	defer StopApplication(port)
}

// func main() {
// 	// Params
// 	chanErr := make(chan error)
// 	path := "./clicki"
// 	port := ":8080"

// 	//
// 	TestApplication(chanErr, path, port)
// }
