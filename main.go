//go:generate rice embed-go
package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/GeertJohan/go.rice"
	"github.com/fatih/color"
	"github.com/sethgrid/multibar"
)

var (
	listen          = flag.String("listen", ":8080", "The address to run the web service on.")
	path            = flag.String("path", ".", "The path files are saved to.")
	disableColor    = flag.Bool("no-color", false, "Disable color output.")
	disableProgress = flag.Bool("no-progress", false, "Disable progress bars.")
)

// TODO: Print with bars if needed.

func main() {
	flag.Parse()
	if *disableColor {
		color.NoColor = true
	}

	// Setup progress bar
	var bars *multibar.BarContainer
	if !*disableProgress {
		var err error
		bars, err = multibar.New()
		if err != nil {
			log.Printf(color.RedString("Failed to initialize progress bars %v"), err)
			return
		}
		go bars.Listen()
	}

	// Setup Server
	mux := http.NewServeMux()

	// Serve static files with rice for portability
	staticFiles := rice.MustFindBox("frontend").HTTPBox()
	mux.Handle("/", http.FileServer(staticFiles))

	// Handle form posts
	mux.HandleFunc("/send", func(rw http.ResponseWriter, r *http.Request) {
		disableProgress := *disableProgress

		// TODO: Progress only reports the file being flushed out of
		// memory/temp files after transfer is complete.

		var length int
		if !disableProgress {
			var err error
			lengthString := r.FormValue("length")
			length, err = strconv.Atoi(lengthString)
			if err != nil {
				// Client didn't supply length. Don't use it.
				disableProgress = true
			}
		}

		inFile, header, err := r.FormFile("file")
		defer inFile.Close()
		if err != nil {
			log.Printf(color.RedString("Failed to parse file: %v"), err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		outFile, err := os.Create(header.Filename)
		defer outFile.Close()
		if err != nil {
			log.Printf(color.RedString("Failed to create file: %v"), err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		var writer io.Writer
		if !disableProgress {
			updateBar := bars.MakeBar(length, header.Filename)
			writer = &ProgressWriter{
				BytesWritten: 0,
				Callback:     updateBar,
				Writer:       outFile,
			}
		} else {
			writer = outFile
		}

		written, err := io.Copy(writer, inFile)
		if err != nil {
			log.Printf(color.RedString("Failed to copy file: %v"), err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		// log.Printf(
		// 	"Wrote "+color.CyanString("%d")+" bytes to "+color.CyanString("%s"),
		// 	written, outFile.Name())
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

type ProgressWriter struct {
	BytesWritten int
	Callback     func(progress int)
	io.Writer
}

func (writer *ProgressWriter) Write(bytes []byte) (int, error) {
	n, err := writer.Writer.Write(bytes)
	writer.BytesWritten += n
	writer.Callback(writer.BytesWritten)
	return n, err
}
