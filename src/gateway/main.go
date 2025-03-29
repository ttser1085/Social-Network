package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
)

var paths map[string]string = map[string]string{
	"/signup": "http://auth:8091/signup",
	"/login":  "http://auth:8091/login",
	"/update": "http://auth:8091/update",
	"/whoami": "http://auth:8091/whoami",
}

func handler(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}

	rbody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read request: %v\n", err)
		http.Error(w, "Failed to read request", http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest(r.Method, paths[r.URL.Path], bytes.NewReader(rbody))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create request: %v\n", err)
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	req.Header = r.Header

	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to forward request: %v\n", err)
		http.Error(w, "Failed to forward request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read response: %v\n", err)
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func main() {
	http.HandleFunc("/", handler)

	port := 8092
	fmt.Printf("Starting server at port %d\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
