/*
Copyright 2014 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// A small utility to return information about the HTTP connection.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	doClose = flag.Bool("close", false, "Close connection per each HTTP request")
	port    = flag.Int("port", 80, "Port number.")
)

func main() {
	flag.Parse()

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("Error from os.Hostname(): %s", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("HTTP request from %s", r.RemoteAddr)

		fmt.Fprintf(w, "%s\n\n", time.Now().Format(time.RFC1123))
		fmt.Fprintf(w, "I am %s:%d\n\n", hostname, *port)
		fmt.Fprintf(w, "You are %s\n\n", r.RemoteAddr)

		dump, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Printf("error dumping request: %v")
			return
		}
		fmt.Fprintf(w, "%s\n", string(dump))

		if *doClose {
			// Add this header to force to close the connection after serving the request.
			w.Header().Add("Connection", "close")
		}
	})
	go func() {
		// Run in a closure so http.ListenAndServe doesn't block
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
	}()

	log.Printf("Serving on port %d.\n", *port)
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGTERM)
	sig := <-signals
	log.Printf("Shutting down after receiving signal: %s.\n", sig)
	log.Printf("Awaiting pod deletion.\n")
	time.Sleep(60 * time.Second)
}
