package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Payload struct {
	Error      string    `json:"error"`
	ServerTime time.Time `json:"server_time"`
}

const (
	// requestLimit is the number simultaneous request we're willing to run
	requestLimit int = 10
	// requestDuration is the amount of time(in seconds) to stall a request's response
	requestDuration time.Duration = 5
)

func main() {
	requestGauge := NewRequestGauge(requestLimit) // creates a new RequestGauge instance & fills the channels with tokens
	indexHandler := getHandler(requestGauge)

	http.HandleFunc("/", indexHandler) // registers route handler with default http ServeMux
	http.ListenAndServe(":4000", nil)  // listen for incoming requests
}

// util functions

func sendJson(w http.ResponseWriter, p any) {
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func getRequestDetail(r *http.Request) string {
	return fmt.Sprintf("%s request for: %s", r.Method, r.URL.Path)
}

func getServerTimeAfter(sec time.Duration) time.Time {
	return <-time.After(sec * time.Second)
}

func getHandler(requestGauge *RequestGauge) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(getRequestDetail(r))

		// this goroutine calls `run` on the requestGauge instance
		err := requestGauge.run(func() {
			sendJson(w, Payload{Error: "", ServerTime: getServerTimeAfter(requestDuration)})
		})

		// if err returned, it indicates we'ev exceeded the number of simultaneous request allow
		if err != nil {
			w.WriteHeader(http.StatusTooManyRequests)
			sendJson(w, Payload{Error: "too many request"})
		}
	}
}
