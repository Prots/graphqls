package main

import (
	"encoding/json"
	"github.com/graphql-go/graphql"
	"log"
	"net/http"
	"time"
)

const (
	dateFormat       = "2006-01-02"
	timeFormat       = "15:04:05"
	secondsPerMinute = 60
	secondsPerHour   = 3600
)

var timeType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "time",
		Fields: graphql.Fields{
			"timeStr": &graphql.Field{
				Type:    graphql.String,
				Resolve: timeStrResolver,
			},
			"timestamp": &graphql.Field{
				Type:    graphql.Int,
				Resolve: timestampResolver,
			},
		},
	},
)

var dateType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "date",
		Fields: graphql.Fields{
			"dateStr": &graphql.Field{
				Type:    graphql.String,
				Resolve: dateStrResolver,
			},
			"timestamp": &graphql.Field{
				Type:    graphql.Int,
				Resolve: timestampOfDateResolver,
			},
		},
	},
)

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"time": &graphql.Field{
				Type:    timeType,
				Resolve: timeObjResolver,
			},
			"date": &graphql.Field{
				Type:    dateType,
				Resolve: dateObjResolver,
			},
		},
	})

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query: queryType,
	},
)

func main() {
	http.HandleFunc("/graphql", graphqlHandler)
	http.HandleFunc("/*", notFoundHandler)
	http.ListenAndServe(":8080", nil)
}

func graphqlHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	queryString := r.URL.Query()
	//log.Printf("query values: %v", queryString)
	if queryString != nil && queryString["query"] != nil {
		result := executeQuery(queryString["query"][0], schema)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("{\"error\":\"enter valid query string\"}"))
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("{\"error\":\"resource not found\"}"))
}

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		log.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}

func timeStrResolver(p graphql.ResolveParams) (interface{}, error) {
	return time.Now().Format(timeFormat), nil
}

func timestampResolver(p graphql.ResolveParams) (interface{}, error) {
	return time.Now().Unix(), nil
}

func dateStrResolver(p graphql.ResolveParams) (interface{}, error) {
	return time.Now().Format(dateFormat), nil
}

func timestampOfDateResolver(p graphql.ResolveParams) (interface{}, error) {
	dt := time.Now()
	secondsNow := dt.Second() + dt.Minute()*secondsPerMinute + dt.Hour()*secondsPerHour
	return secondsNow, nil
}

func dateObjResolver(p graphql.ResolveParams) (interface{}, error) {
	return "dateObject", nil
}

func timeObjResolver(p graphql.ResolveParams) (interface{}, error) {
	return "timeObject", nil
}
