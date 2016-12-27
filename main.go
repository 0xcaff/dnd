//go:generate rice embed-go
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/GeertJohan/go.rice"
	"github.com/fatih/color"
	"github.com/gosuri/uiprogress"
)

var (
	listen          = flag.String("listen", ":8080", "The address to run the web service on.")
	path            = flag.String("path", ".", "The path files are saved to.")
	disableColor    = flag.Bool("no-color", false, "Disable color output.")
	disableProgress = flag.Bool("no-progress", false, "Disable progress bars.")
)

var bars *uiprogress.Progress

func main() {
	flag.Parse()
	if *disableColor {
		color.NoColor = true
	}

	// Setup progress bars
	if !*disableProgress {
		bars = uiprogress.New()
		bars.Start()
	}

	// Setup Server
	mux := http.NewServeMux()

	// Serve static files with rice for portability
	staticFiles := rice.MustFindBox("frontend").HTTPBox()
	mux.Handle("/", http.FileServer(staticFiles))

	// Handle form posts
	mux.HandleFunc("/send", func(rw http.ResponseWriter, r *http.Request) {
		disableProgress := *disableProgress

		reader, err := r.MultipartReader()
		if err != nil {
			log.Printf(color.RedString("Failed to read form: %v"), err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		parsedLength := false
		parsedFile := false
		length := 0

		for {
			part, err := reader.NextPart()

			if err != nil && err == io.EOF {
				// No more parts.
				break
			}

			if err != nil {
				log.Printf(color.RedString("Failed to read part: %v"), err)
				rw.WriteHeader(http.StatusBadRequest)
				continue
			}

			// https://xhr.spec.whatwg.org/#interface-formdata
			// FormData should be sent in order.

			partName := part.FormName()

			if partName == "length" && !disableProgress && !parsedLength {
				_, err = fmt.Fscanf(part, "%d", &length)
				if err != nil {
					log.Printf(color.RedString("Failed to read length: %v"), err)
					disableProgress = true
				} else if length < 0 {
					log.Printf(color.RedString("Invalid length: %d"), err)
					disableProgress = true
				} else {
					parsedLength = true
				}
			} else if partName == "file" && !parsedFile {
				fileName := part.FileName()

				outFile, err := os.Create(fileName)
				defer outFile.Close()
				if err != nil {
					log.Printf(color.RedString("Failed to create file: %v"), err)
					rw.WriteHeader(http.StatusInternalServerError)
					return
				}

				var writer io.Writer
				if parsedLength {
					bar := bars.AddBar(length)
					progWriter := &ProgressWriter{
						Length:       length,
						FileName:     fileName,
						BytesWritten: 0,
						Bar:          bar,
						Writer:       outFile,
					}
					bar.AppendFunc(progWriter.Append())
					bar.PrependFunc(progWriter.Prepend())

					writer = progWriter
				} else {
					log.Printf("Writing "+color.CyanString("%s"), fileName)
					writer = outFile
				}

				written, err := io.Copy(writer, part)
				if err != nil {
					log.Printf(color.RedString("Failed to copy file: %v"), err)
					rw.WriteHeader(http.StatusInternalServerError)
					return
				}

				if !parsedLength {
					log.Printf(
						"Wrote "+color.CyanString("%s")+" to "+color.CyanString("%s"),
						byteUnitStr(int(written)), outFile.Name())
				}
			} else {
				log.Printf(color.CyanString("%s")+" sent invalid request", r.RemoteAddr)
			}

			part.Close()
		}
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
	Length       int
	FileName     string
	BytesWritten int
	Bar          *uiprogress.Bar
	io.Writer
}

func (writer *ProgressWriter) Write(bytes []byte) (int, error) {
	bars.RefreshInterval = time.Millisecond * 300

	n, err := writer.Writer.Write(bytes)
	writer.BytesWritten += n
	writer.Bar.Set(writer.BytesWritten)

	if err == io.EOF {
		// Slow down, end of a progress bar.
		bars.RefreshInterval = time.Second * 10
	}

	return n, err
}

func (writer *ProgressWriter) Prepend() func(*uiprogress.Bar) string {
	return func(bar *uiprogress.Bar) string {
		return writer.FileName
	}
}

func (writer *ProgressWriter) Append() func(*uiprogress.Bar) string {
	total := byteUnitStr(writer.Length)

	return func(bar *uiprogress.Bar) string {
		completed := byteUnitStr(writer.BytesWritten)
		return bar.CompletedPercentString() + " " + completed + "/" + total
	}
}

// TODO: What abotu files over 4GB?
var byteUnits = []string{"B", "KB", "MB", "GB", "TB", "PB"}

// https://github.com/mitchellh/ioprogress/blob/master/draw.go#L91
func byteUnitStr(n int) string {
	var unit string
	size := float64(n)
	for i := 1; i < len(byteUnits); i++ {
		if size < 1000 {
			unit = byteUnits[i-1]
			break
		}

		size = size / 1000
	}

	return fmt.Sprintf("%.3g %s", size, unit)
}
