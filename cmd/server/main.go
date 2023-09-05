package main

import (
	"alerting/cmd/server/handlers"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, handlers.UpdateRequest)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}

}
