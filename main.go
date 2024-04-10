package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

var defaultJson = []byte(`{"key": "value"}`)
var byPath = map[string]string{}

func handler(w http.ResponseWriter, r *http.Request) {
	dump, err := httputil.DumpRequest(r, true)

	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}

	fmt.Println(string(dump))

	path := byPath[r.URL.Path]

	if path == "" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(defaultJson)

		return
	}

	file, err := os.Open(path)

	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}

	bytes := make([]byte, stat.Size())
	_, err = bufio.NewReader(file).Read(bytes)
	if err != nil && err != io.EOF {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

func main() {

	if (len(os.Args) % 2) == 1 {
		for i, v := range os.Args {
			if i == 0 {
				continue
			}

			if (i+1)%2 == 0 {
				log.Printf("define %v -> %v", v, os.Args[i+1])
				byPath[v] = os.Args[i+1]
			}
		}
	} else {
		for i, v := range os.Args {
			if i == 0 {
				continue
			}

			if i == 1 {
				log.Printf("define default -> %v", os.Args[i])
				file, err := os.Open(os.Args[i])

				if err != nil {
					log.Fatalf("failed to open default json file: %v", err)
				}
				defer file.Close()

				stat, err := file.Stat()
				if err != nil {
					log.Fatalf("failed to get stat of default json file: %v", err)
				}

				bytes := make([]byte, stat.Size())
				_, err = bufio.NewReader(file).Read(bytes)

				if err != nil && err != io.EOF {
					log.Fatalf("failed to read default json file: %v", err)
				}

				defaultJson = bytes
			}

			if (i+1)%2 == 1 {
				log.Printf("define %v -> %v", v, os.Args[i+1])
				byPath[v] = os.Args[i+1]
			}
		}
	}

	var httpServer http.Server
	http.HandleFunc("/", handler)
	log.Println("start http listening :18888")
	httpServer.Addr = ":18888"
	log.Println(httpServer.ListenAndServe())
}
