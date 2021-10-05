package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	loggly "github.com/jamespearly/loggly"
)

type JSONtime struct {
	Stamp time.Time
}

func status(w http.ResponseWriter, req *http.Request) {

	// var tag string = "firstapplication"
	// client := loggly.New(tag)

	JSONtime := JSONtime{
		Stamp: time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResp, err := json.Marshal(JSONtime)

	// client.EchoSend("info", "Method:"+string(req.Method)+"\nPath: "+string(req.RequestURI))
	// loggly_message := fmt.Sprintf("Method: " + req.Method + "\n" +
	// 	"IP: " + req.RemoteAddr + "\n" +
	// 	"Path: " + req.RequestURI + "\n" +
	// 	"HTTP STATUS CODE: " + string(http.StatusText(200)))
	// client.EchoSend("info", loggly_message)

	if err != nil {
		// client.Send("error", "This is an error message:"+err.Error())
		log.Fatal("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
	return
}

//https://ndersson.me/post/capturing_status_code_in_net_http/

func wrapHandlerWithLogging(wrappedHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		var tag string = "firstapplication"
		client := loggly.New(tag)

		log.Printf("--> %s %s", req.Method, req.URL.Path)

		lrw := NewLoggingResponseWriter(w)
		wrappedHandler.ServeHTTP(lrw, req)

		loggly_message := fmt.Sprintf("Method: " + req.Method + "\n" +
			"IP: " + req.RemoteAddr + "\n" +
			"Path: " + req.RequestURI + "\n" +
			"HTTP STATUS CODE: " + strconv.Itoa(lrw.statusCode))

		client.EchoSend("info", loggly_message)

		statusCode := lrw.statusCode
		log.Printf("<-- %d %s", statusCode, http.StatusText(statusCode))

	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	// WriteHeader(int) is not called if our response implicitly returns 200 OK, so
	// we default to that status code.
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/status", status).Methods("GET")
	// handler := wrapHandlerWithLogging(r.HandleFunc("/status", status).Methods("GET"))
	// wrapHandlerWithLogging(r)
	// r.Use(wrapHandlerWithLogging)
	http.Handle("/", r)
	http.ListenAndServe(":8080", wrapHandlerWithLogging(r))

}
