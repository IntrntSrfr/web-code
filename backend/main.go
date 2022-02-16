package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(5 * time.Second))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/check", check)
	r.Post("/run", run)

	http.ListenAndServe(":8008", r)
}

type JSONResp struct {
	Value string
}

func run(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		JSONError(w, err, http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	reqBody := struct {
		Inp string `json:"inp"`
	}{}

	err = json.Unmarshal(b, &reqBody)
	if err != nil {
		JSONError(w, err, http.StatusInternalServerError)
		return
	}

	out, err := runCode(r.Context(), []byte(reqBody.Inp))
	if err != nil {
		JSONError(w, err, http.StatusBadRequest)
		return
	}

	fmt.Println(out)

	JSONResponse(w, out, http.StatusOK)
}

func JSONResponse(w http.ResponseWriter, value string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(&JSONResp{Value: value})
}

func JSONError(w http.ResponseWriter, err error, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(&JSONResp{Value: err.Error()})
}

// runCode runs a program, and returns a string with the output, or an empty string if an error occurs
func runCode(ctx context.Context, code []byte) (string, error) {
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

	cmd := exec.CommandContext(ctx, "docker", "run", "--rm", "-v",
		fmt.Sprintf("%v/%v:/go/src/app", cwd, dir),
		"-w", "/go/src/app", "golang:alpine", "go", "run", "main.go")

	stdout, err := cmd.CombinedOutput()

	fmt.Println(stdout, err)

	return string(stdout), err
}

func check(w http.ResponseWriter, r *http.Request) {
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
