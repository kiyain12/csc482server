package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gorilla/mux"
	loggly "github.com/jamespearly/loggly"
)

type JSONtime struct {
	Stamp time.Time
}

func server(w http.ResponseWriter, req *http.Request) {

	//new 10/21/2021
	//set the loggly_token env variable
	// os.Setenv("LOGGLY_TOKEN", "")

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

	client2 := loggly.New("MuxImplementation")
	client2.EchoSend("info", "Method Type: "+req.Method+", Request Path: "+req.URL.Path+", Status Code: 200")
	if err != nil {
		// client.Send("error", "This is an error message:"+err.Error())
		log.Fatal("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
	return
}

//new 10/21/2021

func badMethod(w http.ResponseWriter, req *http.Request) {
	//set the loggly_token env variable
	// os.Setenv("LOGGLY_TOKEN", "")

	//initialize loggly client
	client2 := loggly.New("MuxImplementation")

	fmt.Printf("BAD METHOD REQ: %v\n", req)
	w.WriteHeader(405)

	//Send message to loggly
	//Return 405
	client2.EchoSend("error", "Method Type: "+req.Method+", Request Path: "+req.URL.Path+", Status Code: 405")
}

//new 10/21/2021
func badPath(w http.ResponseWriter, req *http.Request) {
	//set the loggly_token env variable
	// os.Setenv("LOGGLY_TOKEN", "")

	//initialize loggly client
	client2 := loggly.New("MuxImplementation")
	fmt.Printf("BAD PATH REQ: %v\n", req)
	w.WriteHeader(404)

	client2.EchoSend("error", "Method Type: "+req.Method+", Request Path: "+req.URL.Path+", Status Code: 404")
}

//https://ndersson.me/post/capturing_status_code_in_net_http/

// func wrapHandlerWithLogging(wrappedHandler http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

// 		var tag string = "firstapplication"
// 		client := loggly.New(tag)

// 		log.Printf("--> %s %s", req.Method, req.URL.Path)

// 		lrw := NewLoggingResponseWriter(w)
// 		wrappedHandler.ServeHTTP(lrw, req)

// 		loggly_message := fmt.Sprintf("Method: " + req.Method + "\n" +
// 			"IP: " + req.RemoteAddr + "\n" +
// 			"Path: " + req.RequestURI + "\n" +
// 			"HTTP STATUS CODE: " + strconv.Itoa(lrw.statusCode))

// 		client.EchoSend("info", loggly_message)

// 		statusCode := lrw.statusCode
// 		log.Printf("<-- %d %s", statusCode, http.StatusText(statusCode))

// 	})
// }

// type loggingResponseWriter struct {
// 	http.ResponseWriter
// 	statusCode int
// }

// func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
// 	// WriteHeader(int) is not called if our response implicitly returns 200 OK, so
// 	// we default to that status code.
// 	return &loggingResponseWriter{w, http.StatusOK}
// }

// func (lrw *loggingResponseWriter) WriteHeader(code int) {
// 	lrw.statusCode = code
// 	lrw.ResponseWriter.WriteHeader(code)
// }
type weatherData struct {
	ID            int     `json:"id"`
	DATE          string  `json:"applicable_date"`
	TEMPERATURE   float64 `json:"the_temp`
	WEATHERSTATUS string  `json:"weather_state_name"`
}

type dbitem struct {
	Id   string        `json:"Id"`
	Time string        `json:"Time"`
	Data []weatherData `json:"Data"`
}

func all(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		fmt.Println("Error starting a new session")
	}
	// Create DynamoDB client
	svc := dynamodb.New(sess)
	// Name of the table
	tableName := "npayag-weather-table-csc-482"

	//Scan the DB for all items
	out, err := svc.Scan(&dynamodb.ScanInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		log.Fatal("Error scanning DB", err)
	}

	// var x []weatherData
	//unmarshal response
	var dbitem []dbitem
	// dbitem.Data = x
	err = dynamodbattribute.UnmarshalListOfMaps(out.Items, &dbitem)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record:", err))
	}

	json.NewEncoder(w).Encode(dbitem)
}

type TableStatus struct {
	Table string
	Count *int64 //*(asterisk) attached to a variable or expression (*v) indicates a pointer dereference. That is, take the value the variable is pointing at.
}

func status(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// jsonResp, err := json.Marshal()
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		fmt.Println("Error starting a new session")
	}
	// Create DynamoDB client
	svc := dynamodb.New(sess)
	// Name of the table
	tableName := "npayag-weather-table-csc-482"

	input := &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	}

	result, err := svc.DescribeTable(input)
	if err != nil {
		log.Fatalf("Got error describing the table", err)
	}

	var statusResponse TableStatus
	statusResponse.Table = "csc_482"
	statusResponse.Count = result.Table.ItemCount

	json.NewEncoder(w).Encode(statusResponse)

}

func main() {

	// os.Setenv("LOGGLY_TOKEN", "")
	// os.Setenv("AWS_ACCESS_KEY_ID", "")
	// os.Setenv("AWS_SECRET_ACCESS_KEY", "")
	// os.Setenv("AWS_SESSION_TOKEN", "")

	r := mux.NewRouter()
	r.HandleFunc("/npayag/server", server).Methods("GET")
	r.HandleFunc("/npayag/status", status).Methods("GET")
	r.HandleFunc("/npayag/all", all).Methods("GET")

	//old comment
	// handler := wrapHandlerWithLogging(r.HandleFunc("/status", status).Methods("GET"))
	// wrapHandlerWithLogging(r)
	// r.Use(wrapHandlerWithLogging)

	//new 10/21/2021
	// r := mux.NewRouter()
	// r.HandleFunc("/status", status).Methods("GET")
	r.PathPrefix("/").Methods("POST", "PUT", "DELETE").HandlerFunc(badMethod) //catch bad methods
	r.PathPrefix("/").HandlerFunc(badPath)                                    //everything else is bad path
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))

	http.Handle("/", r)
	// http.ListenAndServe(":8080", wrapHandlerWithLogging(r))

}
