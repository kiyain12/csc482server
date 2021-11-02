package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
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
	TEMPERATURE   float64 `json:"the_temp"`
	WEATHERSTATUS string  `json:"weather_state_name"`
	TIMESTAMP     string  `json:"created"`
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
	statusResponse.Table = "npayag-weather-table-csc-482"
	statusResponse.Count = result.Table.ItemCount

	json.NewEncoder(w).Encode(statusResponse)

}

// func getDates(w http.ResponseWriter, r *http.Request) {
// 	query := r.URL.Query()
// 	filters := query.Get("applicable_date") //filters="color"
// 	w.WriteHeader(200)
// 	w.Write([]byte(filters))
// }

func search(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	urlquery := mux.Vars(req)["forecastdate"]
	format, err := regexp.MatchString("^\\d{4}\\-(0[1-9]|1[012])\\-(0[1-9]|[12][0-9]|3[01])$", urlquery)
	if err != nil {
		log.Fatal(err)
	}

	if format {
		w.WriteHeader(http.StatusOK)
	}

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

	//make the expression
	//https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/expression/#Contains
	filt := expression.Contains(expression.Name("Time"), urlquery)

	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		log.Fatalf("Got error building expression: %s", err)
	}

	// Build the query input parameters
	p := &dynamodb.ScanInput{
		TableName:                 aws.String(tableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
	}
	// Get all data for given date
	out, err := svc.Scan(p)

	if err != nil {
		log.Fatalf("Query API call failed: %s", err)
	}
	var dbitem []dbitem
	err = dynamodbattribute.UnmarshalListOfMaps(out.Items, &dbitem)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record:", err))
	}

	json.NewEncoder(w).Encode(dbitem)

	//create a new struct that just holds filtered data
	//inside DATA array, i want the applicable_date and the temperature
	//so that when http://localhost:8080/npayag/search?filter=applicable_data it just shows the temperature for that date

	// for index, value := range dbitem {
	// }

	// query := req.URL.Query()

	// filters := query.Get("Id")
	// w.Write([]byte(filters))

	// filters, present := query["filters"] //filters=["color", "price", "brand"]
	// if !present || len(filters) == 0 {
	// 	fmt.Println("filters not present")
	// }
	// w.WriteHeader(200)
	// w.Write([]byte(strings.Join(filters, ",")))
}

func main() {

	// os.Setenv("LOGGLY_TOKEN", "")
	// os.Setenv("AWS_ACCESS_KEY_ID", "")
	// os.Setenv("AWS_SECRET_ACCESS_KEY", "")
	// os.Setenv("AWS_SESSION_TOKEN", "")

	r := mux.NewRouter()

	// getDatesHandler := http.HandlerFunc(getDates)
	// r.HandleFunc("/npayag/all/date", getDates)
	// http.Handle("/npayag/all/applicable_date", getDatesHandler)

	r.HandleFunc("/npayag/server", server).Methods("GET")
	r.HandleFunc("/npayag/status", status).Methods("GET")

	r.HandleFunc("/npayag/all", all).Methods("GET")

	// getSearchHandler := http.HandlerFunc(search)
	r.HandleFunc("/npayag/search", search).Queries("forecastdate", "{forecastdate:.*}")

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

	// http.Handle("/", r)
	// http.ListenAndServe(":8080", wrapHandlerWithLogging(r))

}
