package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/flytam/filenamify"
	"github.com/google/uuid"
)

var port int
var ext = "jpg"

func init() {
	flag.IntVar(&port, "port", port, "Port to listen")
	flag.Parse()

	log.Printf("Running on port :%d\n", port)
}

func main() {
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		img := strings.TrimPrefix(r.URL.Path, "/products/")
		output, err := filenamify.Filenamify(img, filenamify.Options{
			Replacement: "",
		})

		// Is valid URL?
		if output != img || err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid ID!"))

			return
		}

		// Is valid ID?
		if _, err := uuid.Parse(img); err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte("Invalid ID format!"))

			return
		}

		// Find and serve file
		path, err := filepath.Abs("products/" + output)
		image := fmt.Sprintf("%s.%s", path, ext)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal error!"))

			return
		}

		if _, err := os.Stat(image); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("Image not found!"))

				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal error!"))

			return
		}

		http.ServeFile(w, r, image)
	})

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
