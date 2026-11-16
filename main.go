package main

import (
	"log"
	"net/http"
)

func main() {
	const filepathRoot = "."
	const port = "8080"
	serverMux := http.NewServeMux()
	serverMux.Handle("/", http.FileServer(http.Dir(filepathRoot)))

	serv := http.Server{
		Handler: serverMux,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(serv.ListenAndServe())

}
