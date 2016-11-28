package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type dataObj struct {
	Data struct {
		Date struct {
			DateStr       string `json:"dateStr"`
			DateTimestamp int64  `json:"timestamp"`
		} `json:"date"`
		Time struct {
			TimeStr   string `json:"timeStr"`
			Timestamp int64  `json:"timestamp"`
		} `json:"time"`
	} `json:"data"`
}

type handlerFun func(w http.ResponseWriter, r *http.Request)

func TestGraphqlService(t *testing.T) {
	var tests = []struct {
		query        string
		responseCode int
		responseBody string
		handler      handlerFun
		checkBody    bool
	}{
		{
			"graphql/",
			http.StatusBadRequest,
			"{\"error\":\"enter valid query string\"}",
			graphqlHandler,
			true,
		},
		{
			"graphql?query={date{dateStr,%20timestamp}}",
			http.StatusOK,
			" ",
			graphqlHandler,
			false,
		},
		{
			"graphql?query={time{timeStr,%20timestamp}}",
			http.StatusOK,
			" ",
			graphqlHandler,
			false,
		},
		{
			"badpath/",
			http.StatusNotFound,
			"{\"error\":\"resource not found\"}",
			notFoundHandler,
			true,
		},
	}
	for _, test := range tests {
		req, err := http.NewRequest("GET", "http://localhost:8080/"+test.query, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(test.handler)
		handler.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		assert(rr.Code, test.responseCode, t)

		//Check the response body is what we expect.
		if test.checkBody {
			assert(rr.Body.String(), test.responseBody, t)
		}
	}
}

func TestGraphqlResponse(t *testing.T) {
	var dtNow = time.Now()
	var dayTimestamp = (int64)(dtNow.Second() + dtNow.Minute()*secondsPerMinute + dtNow.Hour()*secondsPerHour)
	var tests = []struct {
		query       string
		timeStr     string
		timestamp   int64
		dateStr     string
		timestampDt int64
	}{
		{
			"graphql?query={date{dateStr,timestamp}}",
			"",
			0,
			dtNow.Format(dateFormat),
			dayTimestamp,
		},
		{
			"graphql?query={time{timeStr,timestamp}}",
			dtNow.Format(timeFormat),
			dtNow.Unix(),
			"",
			0,
		},
		{
			"graphql?query={date{timestamp}}",
			"",
			0,
			"",
			dayTimestamp,
		},
		{
			"graphql?query={time{timestamp}}",
			"",
			dtNow.Unix(),
			"",
			0,
		},
		{
			"graphql?query={date{dateStr}}",
			"",
			0,
			dtNow.Format(dateFormat),
			0,
		},
		{
			"graphql?query={time{timeStr}}",
			dtNow.Format(timeFormat),
			0,
			"",
			0,
		},
	}
	for _, test := range tests {
		req, err := http.NewRequest("GET", "http://localhost:8080/"+test.query, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(graphqlHandler)
		handler.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		assert(rr.Code, http.StatusOK, t)

		//Check the response body is what we expect.
		decBody := &dataObj{}
		body := rr.Body
		t.Logf("Response Body: %v", body)
		err = json.NewDecoder(body).Decode(decBody)
		if err != nil {
			t.Errorf("handler returned bad JSON body: got %v, err: %v", body.String(), err)
		}
		assert(decBody.Data.Date.DateStr, test.dateStr, t)
		assert(decBody.Data.Date.DateTimestamp, test.timestampDt, t)
		assert(decBody.Data.Time.Timestamp, test.timestamp, t)
		assert(decBody.Data.Time.TimeStr, test.timeStr, t)
	}
}

func BenchmarkGraphqlService(b *testing.B) {
	b.Run("Date", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			requestUrl("graphql?query={date{dateStr,timestamp}}", b)
		}
	})
	b.Run("Time", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			requestUrl("graphql?query={time{timeStr,timestamp}}", b)
		}
	})
}

func requestUrl(query string, b *testing.B) {
	req, err := http.NewRequest("GET", "http://localhost:8080/"+query, nil)
	if err != nil {
		b.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(graphqlHandler)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if rr.Code != http.StatusOK {
		b.Logf("Assertion Error: got %v want %v",
			rr.Code, http.StatusOK)
	}
}

func assert(got interface{}, want interface{}, t *testing.T) {
	if got != want {
		t.Errorf("Assertion Error: got %v want %v",
			got, want)
	}
}
