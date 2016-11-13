package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/GeertJohan/go.rice"
	"github.com/fatih/color"
)

var (
	listen       = flag.String("listen", ":8080", "The address to run the web service on")
	path         = flag.String("path", ".", "The path files are saved to.")
	disableColor = flag.Bool("no-color", false, "Disable color output")
)

func main() {
	flag.Parse()
	if *disableColor {
		color.NoColor = true
	}

	// Setup Server
	mux := http.NewServeMux()

	// Serve static files with rice for portability
	staticFiles := rice.MustFindBox("frontend").HTTPBox()
	mux.Handle("/", http.FileServer(staticFiles))

	// Handle form posts
	mux.HandleFunc("/send", func(rw http.ResponseWriter, r *http.Request) {
		inFile, header, err := r.FormFile("file")
		if err != nil {
			log.Printf(color.RedString("Failed to parse file: %v"), err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		defer inFile.Close()

		outFile, err := os.Create(header.Filename)
		defer outFile.Close()
		if err != nil {
			log.Printf(color.RedString("Failed to create file: %v"), err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		written, err := io.Copy(outFile, inFile)
		if err != nil {
			log.Printf(color.RedString("Failed to copy file: %v"), err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf(
			"Wrote "+color.CyanString("%d")+" bytes to "+color.CyanString("%s"),
			written, outFile.Name())
	})

	s := &http.Server{
		Addr:    *listen,
		Handler: mux,
	}

	err := os.Chdir(*path)
	if err != nil {
		log.Printf(color.RedString("Failed to change to %s, %v"), *path, err)
		return
	}

	log.Printf("Starting server on %s", color.CyanString(*listen))
	log.Printf("Saving files to %s", color.CyanString(*path))

	log.Fatal(s.ListenAndServe())
}
