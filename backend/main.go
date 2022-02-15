package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

func main() {
	http.HandleFunc("/run", run)
	http.HandleFunc("/check", check)
	http.ListenAndServe(":8008", nil)
}

type CodeRespo struct {
	Ok     bool
	Output string
}

func run(w http.ResponseWriter, r *http.Request) {

	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Add("Access-Control-Allow-Headers",
			"Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
	}
	// Stop here if its Preflighted OPTIONS request
	if r.Method == "OPTIONS" {
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	reqBody := struct {
		Inp string `json:"inp"`
	}{}

	json.Unmarshal(b, &reqBody)

	out, err := runCode([]byte(reqBody.Inp))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

// runCode runs a program, and returns a string with the output, or an empty string if an error occurs
func runCode(code []byte) (string, error) {
	dir, err := os.MkdirTemp("./tmp", "")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(dir)
	//fmt.Println(dir)

	f, err := os.Create(dir + "/main.go")
	if err != nil {
		return "", err
	}
	defer f.Close()
	f.Write(code)

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	cmd := exec.Command("docker", "run", "--rm", "-v",
		fmt.Sprintf("%v/%v:/go/src/app", cwd, dir),
		"-w", "/go/src/app", "golang:alpine", "go", "run", "main.go")

	//fmt.Println(cmd.Args)

	stdout, err := cmd.CombinedOutput()
	return string(stdout), err
}

func check(w http.ResponseWriter, r *http.Request) {
	// deal with cors
	origin := r.Header.Get("Origin")
	w.Header().Set("Access-Control-Allow-Origin", origin)

	resp := struct {
		Ok   bool
		Text string
	}{
		Ok:   true,
		Text: "working fantastic",
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
