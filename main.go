package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/GeertJohan/go.rice"
)

func main() {
	// TODO: Parse Command Line and Configuration Information
	mux := http.NewServeMux()

	staticFiles := rice.MustFindBox("frontend").HTTPBox()
	mux.Handle("/", http.FileServer(staticFiles))

	mux.HandleFunc("/send", func(rw http.ResponseWriter, r *http.Request) {
		inFile, header, err := r.FormFile("file")
		defer inFile.Close()
		if err != nil {
			log.Printf("Failed to parse file: %s\n", err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		outFile, err := os.Create(header.Filename)
		defer outFile.Close()
		if err != nil {
			log.Printf("Failed to create file: %s\n", err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		written, err := io.Copy(outFile, inFile)
		if err != nil {
			log.Printf("Failed to copy file: %s\n", err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("Wrote %d to %s\n", written, outFile.Name())
	})

	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Fatal(s.ListenAndServe())
}
