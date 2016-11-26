package main

import (
	"encoding/json"
	"github.com/graphql-go/graphql"
	"log"
	"net/http"
	"time"
)

type tm struct {
	timeStr   string `json:"timeStr"`
	timestamp int64  `json:"timestamp"`
}

type dt struct {
	dateStr   string `json:"dateStr"`
	timestamp int64  `json:"timestamp"`
}

var timeType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "time",
		Fields: graphql.Fields{
			"timeStr": &graphql.Field{
				Type:    graphql.String,
				Resolve: timeStrResolver,
			},
			"timestamp": &graphql.Field{
				Type:    graphql.String,
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
				Type:    graphql.String,
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
	http.ListenAndServe(":8080", nil)
}

func graphqlHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	queryString := r.URL.Query()
	log.Printf("query values: %v", queryString)
	if queryString != nil && queryString["query"] != nil {
		result := executeQuery(queryString["query"][0], schema)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode("{\"error\":\"enter valid query string\"}")
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
	return time.Now().Format("15:04:05"), nil
}

func timestampResolver(p graphql.ResolveParams) (interface{}, error) {
	return time.Now().Unix(), nil
}

func dateStrResolver(p graphql.ResolveParams) (interface{}, error) {
	return time.Now().Format("2006-01-02"), nil
}

func timestampOfDateResolver(p graphql.ResolveParams) (interface{}, error) {
	return time.Now().Unix(), nil
}

func dateObjResolver(p graphql.ResolveParams) (interface{}, error) {
	return "dateObject", nil
}

func timeObjResolver(p graphql.ResolveParams) (interface{}, error) {
	return "timeObject", nil
}