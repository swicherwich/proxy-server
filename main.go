package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

func main() {
	http.ListenAndServe(":8080", http.HandlerFunc(handler))
}

func handler(writer http.ResponseWriter, request *http.Request) {
	targetURL, err := url.Parse(request.RequestURI)

	if err != nil {
		log.Fatal(err)
	}

	request.Host = targetURL.Host
	request.URL.Host = targetURL.Host
	request.URL.Scheme = targetURL.Scheme
	request.RequestURI = ""

	log.Println("Connection to ", targetURL.Host)
	log.Println("Requesting ", targetURL)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Fatal(writer, err)
	}

	for key, values := range response.Header {
		for _, value := range values {
			writer.Header().Set(key, value)
		}
	}

	done := make(chan bool)
	go func() {
		for {
			select {
			case <-time.Tick(10 * time.Millisecond):
				writer.(http.Flusher).Flush()
			case <-done:
				return
			}
		}
	}()

	writer.WriteHeader(response.StatusCode)
	io.Copy(writer, response.Body)
	close(done)
}
