package meeseeks

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPatternMatch(t *testing.T) {

	var matchingTest = []struct {
		RouteMethod string
		RoutePath   string

		RequestMethod string
		RequestPath   string

		ExpectedStatus int
		ExpectedBody   string
	}{

		{
			"GET",
			"/one/",

			"GET",
			"/one/",

			200,
			"matched",
		},
		{
			"GET",
			"/theManIScoming/",

			"GET",
			"/theManIScoming/",

			200,
			"matched",
		}, {
			"GET",
			"/Barbosa/",

			"POST",
			"/Barbosa/",

			405,
			"",
		},
		{
			"GET",
			"/Barbosa/:name",

			"GET",
			"/Barbosa/",

			404,
			"",
		}, {
			"GET",
			"/Pirates/love/to/pirate/",

			"GET",
			"/Pirates/love/to/pirate/",

			200,
			"matched",
		},
		{
			"GET",
			"/",

			"GET",
			"/",

			200,
			"matched",
		},
		{
			"GET",
			"/,",

			"POST",
			"/",

			404,
			"",
		},
	}
	for _, test := range matchingTest {

		rt := NewMeeseeks()
		hf := func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("matched"))
		}

		if test.RouteMethod == "GET" {
			rt.GET(test.RoutePath, hf)
		} else {
			rt.POST(test.RoutePath, hf)
		}

		re, err := http.NewRequest(test.RequestMethod, test.RequestPath, nil)
		if err != nil {
			t.Errorf("Error from NewRequest creation: %s ", err)
		}

		rr := httptest.NewRecorder()
		rt.ServeHTTP(rr, re)

		reqResponse := rr.Result()
		if reqResponse.StatusCode != test.ExpectedStatus {
			t.Errorf("Status code %v but expected %v", reqResponse.StatusCode, test.ExpectedStatus)
			continue
		}

		if reqResponse.StatusCode == http.StatusOK {
			defer reqResponse.Body.Close()
			body, err := ioutil.ReadAll(reqResponse.Body)

			if err != nil {
				t.Errorf("could not read response body got: %s", err)
			}

			if string(body) != test.ExpectedBody {
				t.Errorf("response body contains %v but expected %v", body, test.ExpectedBody)
			}

		}

	}
}

func TestPathParam(t *testing.T) {

	var matchingTest = []struct {
		RouteMethod string
		RoutePath   string

		RequestMethod string
		RequestPath   string

		ExpectedParamValue string
		ParamName          string
		HasParam           bool
	}{

		{
			"GET",
			"/one/:name",

			"GET",
			"/one/jim",

			"jim",
			"name",
			true,
		},
		{
			"GET",
			"/his/:roleName/is/coming",

			"GET",
			"/his/highness/is/coming",

			"highness",
			"roleName",
			true,
		},
	}

	for _, test := range matchingTest {
		rt := NewMeeseeks()

		var ctx context.Context
		hf := func(w http.ResponseWriter, r *http.Request) {
			ctx = r.Context()
		}

		if test.RouteMethod == "GET" {
			rt.GET(test.RoutePath, hf)
		} else {
			rt.POST(test.RoutePath, hf)
		}
		re, err := http.NewRequest(test.RequestMethod, test.RequestPath, nil)
		if err != nil {
			t.Errorf("Error from NewRequest creation: %s ", err)
		}

		rr := httptest.NewRecorder()
		rt.ServeHTTP(rr, re)

		if test.HasParam {
			paramValue := LoadParam(ctx, test.ParamName)
			if paramValue != test.ExpectedParamValue {
				t.Errorf("expected path parameter value %v but got %v", test.ExpectedParamValue, paramValue)
			}
		}

	}

}

func TestTrailingBackslash(t *testing.T) {

	var matchingTest = []struct {
		RouteMethod string
		RoutePath   string

		RequestMethod string
		RequestPath   string

		ExpectedStatus int
	}{

		{
			"GET",
			"/one/name",

			"GET",
			"/one/name/",

			404,
		},
		{
			"GET",
			"/his/:roleName/is/coming",

			"GET",
			"/his/highness/is/coming/",

			404,
		},

		{
			"GET",
			"/checking/path",

			"GET",
			"/checking/path/",

			404,
		},
	}

	for _, test := range matchingTest {
		rt := NewMeeseeks()

		hf := func(w http.ResponseWriter, r *http.Request) {
		}

		if test.RouteMethod == "GET" {
			rt.GET(test.RoutePath, hf)
		} else {
			rt.POST(test.RoutePath, hf)
		}
		re, err := http.NewRequest(test.RequestMethod, test.RequestPath, nil)
		if err != nil {
			t.Errorf("Error from NewRequest creation: %s ", err)
		}
		rr := httptest.NewRecorder()
		rt.ServeHTTP(rr, re)

		reqResponse := rr.Result()
		if reqResponse.StatusCode != test.ExpectedStatus {
			t.Errorf("Status code %v but expected %v\n", reqResponse.StatusCode, test.ExpectedStatus)
		}

	}
}
