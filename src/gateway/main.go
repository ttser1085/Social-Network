package main

import (
	"fmt"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://localhost:8091", http.StatusSeeOther)
}

func main() {
	http.HandleFunc("/", handler)

	port := 8092
	fmt.Printf("Starting server at port %d\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
