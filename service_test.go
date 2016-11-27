package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type handlerFun func(w http.ResponseWriter, r *http.Request)

func TestGraphqlHandler(t *testing.T) {
	var tests = []struct{
		qeury 		string
		responseCode 	int
		responseBody 	string
		handler 	handlerFun
		bodyAssrtType	string
	}{
		{
			"graphql/",
			http.StatusBadRequest,
			"{\"error\":\"enter valid query string\"}",
			graphqlHandler,
			"equals",
		},
		{
			"graphql?query={date{dateStr,%20timestamp}}",
			http.StatusOK,
			"{\"error\":\"enter valid query string\"}",
			graphqlHandler,
			"match",
		},
		{
			"graphql?query={time{timeStr,%20timestamp}}",
			http.StatusOK,
			"{\"error\":\"enter valid query string\"}",
			graphqlHandler,
			"match",
		},
		{
			"badpath/",
			http.StatusNotFound,
			"{\"error\":\"resource not found\"}",
			notFoundHandler,
			"equals",
		},
	}
	for _, test := range tests{
		req, err := http.NewRequest("GET", "http://localhost:8080/" + test.qeury, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(test.handler)
		handler.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != test.responseCode {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, test.responseCode)
		}

		//Check the response body is what we expect.
		if test.bodyAssrtType == "equals" {
			if rr.Body.String() != test.responseBody {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), test.responseBody)
			}
		}
	}
}