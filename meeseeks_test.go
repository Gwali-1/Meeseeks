package meeseeks

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPatternMatch(t *testing.T) {
	rt := NewMeeseeks()
	hf := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("matched"))
	}
	rt.GET("/one", hf)
	re, err := http.NewRequest(http.MethodGet, "/one", nil)
	if err != nil {
		t.Errorf("Error from NewRequest creation: %s ", err)
	}

	rr := httptest.NewRecorder()
	rt.ServeHTTP(rr, re)

	reqResponse := rr.Result()
	if reqResponse.StatusCode != http.StatusOK {
		t.Errorf("Status code %v but expected %v", reqResponse.StatusCode, http.StatusOK)
	}

	defer reqResponse.Body.Close()
	body, err := ioutil.ReadAll(reqResponse.Body)

	if err != nil {
		t.Errorf("could not read response body got: %s", err)
	}

	if string(body) != "matched" {
		t.Errorf("response body contains %v but expected %v", body, "matched")
	}

}

var ctx context.Context

func TestPathParam(t *testing.T) {
	rt := NewMeeseeks()

	hf := func(w http.ResponseWriter, r *http.Request) {
		ctx = r.Context()
	}
	rt.GET("/one/:mine", hf)
	re, err := http.NewRequest(http.MethodGet, "/one/yours", nil)
	if err != nil {
		t.Errorf("Error from NewRequest creation: %s ", err)
	}

	rr := httptest.NewRecorder()
	rt.ServeHTTP(rr, re)

	paramValue := LoadParam(ctx, "mine")
	fmt.Println(paramValue)
	if paramValue != "yours" {
		t.Errorf("expected path parameter value %v but got %v", "mine", paramValue)
	}

}
