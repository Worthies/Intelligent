package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	port := flag.String("port", "80", "Port to listen on")
	wait := flag.Bool("wait", false, "Wait for file to have content for response")
	target := flag.String("target", "", "Target URL to forward requests")
	timeout := flag.Duration("timeout", 180*time.Second, "Timeout for waiting in -wait mode")
	flag.Parse()

	if *wait && *target != "" {
		log.Fatal("Cannot use -wait and -target at the same time")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if *target != "" {
			// Forward to target
			url := *target + r.URL.Path + "?" + r.URL.RawQuery
			resp, err := http.Get(url)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			key := r.Method + r.URL.Path + r.URL.RawQuery
			hash := md5.Sum([]byte(key))
			fileName := fmt.Sprintf("response-%s.txt", hex.EncodeToString(hash[:]))
			err = os.WriteFile(fileName, body, 0644)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/xml;charset=utf-8")
			w.Write([]byte("{\"Result\":\"OK\"}"))
		} else if *wait {
			// Generate file name based on request hash
			key := r.Method + r.URL.Path + r.URL.RawQuery
			hash := md5.Sum([]byte(key))
			fileName := fmt.Sprintf("response-%s.txt", hex.EncodeToString(hash[:]))

			// Build cURL command
			curl := fmt.Sprintf("curl -X %s '%s'", r.Method, r.URL.String())
			for k, v := range r.Header {
				curl += fmt.Sprintf(" -H '%s: %s'", k, strings.Join(v, ","))
			}

			fmt.Printf("Request: %s\nFile: %s\n", curl, fileName)

			// Wait for file to have content with timeout
			timeoutChan := time.After(*timeout)
			for {
				select {
				case <-timeoutChan:
					http.Error(w, "Request timeout", http.StatusRequestTimeout)
					return
				default:
					data, err := os.ReadFile(fileName)
					if err == nil && len(data) > 0 {
						w.Header().Set("Content-Length", strconv.Itoa(len(data)))
						w.Write(data)
						w.Header().Set("Date", time.Now().Format(time.RFC1123))
						w.WriteHeader(http.StatusOK)
						return
					}
					time.Sleep(100 * time.Millisecond)
				}
			}
		} else {
			// Default response
			w.Header().Set("Content-Length", "0")
			w.Header().Set("Content-Type", "application/xml;charset=utf-8")
			w.Header().Set("Date", time.Now().Format(time.RFC1123))
			w.WriteHeader(http.StatusOK)
		}
	})

	fmt.Printf("Starting server on port %s\n", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
