package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", getRootHandler)

	if err := http.ListenAndServe(":5000", mux); err != nil {
		fmt.Println(err)
	}

	os.Exit(0)
}

func getRootHandler(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(rw, fmt.Sprintf("Method not allowed, method [%s]\n", r.Method), http.StatusMethodNotAllowed)

		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	fmt.Printf("Got request: method [%s] uri [%s] data [%s]\n", r.Method, r.RequestURI, string(data))

	rw.WriteHeader(http.StatusOK)
}
