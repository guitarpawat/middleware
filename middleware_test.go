package middleware

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testDoable struct {
	doFunc func(w http.ResponseWriter, r *http.Request, v *ValueMap)
}

func (d testDoable) Do(w http.ResponseWriter, r *http.Request, v *ValueMap) {
	d.doFunc(w, r, v)
}

func makeTestDoable(f func(w http.ResponseWriter, r *http.Request, v *ValueMap)) testDoable {
	return testDoable{f}
}

var notFound = func(w http.ResponseWriter, r *http.Request, v *ValueMap) {
	w.WriteHeader(http.StatusNotFound)
}
var doNotFound = makeTestDoable(notFound)

var ok = func(w http.ResponseWriter, r *http.Request, v *ValueMap) {
	w.WriteHeader(http.StatusOK)
	(*v)["next"] = true
}
var doOK = makeTestDoable(ok)

var writeHello = func(w http.ResponseWriter, r *http.Request, v *ValueMap) {
	w.Write([]byte("hello"))
}
var doWriteHello = makeTestDoable(writeHello)

func TestMiddleware(t *testing.T) {
	mw := MakeMiddleware(nil, doNotFound)
	mockWriter := httptest.NewRecorder()
	mw.ServeHTTP(mockWriter, httptest.NewRequest("GET", "/", nil))

	result := mockWriter.Result()
	actual := result.StatusCode
	expected := http.StatusNotFound
	if expected != actual {
		t.Error("expected status:", expected, "but get:", actual)
	}
}

func TestMiddlewareChain(t *testing.T) {
	mw := MakeMiddleware(nil, doOK, doWriteHello)
	mockWriter := httptest.NewRecorder()
	mw.ServeHTTP(mockWriter, httptest.NewRequest("GET", "/", nil))

	result := mockWriter.Result()
	actual1 := result.StatusCode
	expected1 := http.StatusOK
	if expected1 != actual1 {
		t.Error("expected status:", expected1, "but get:", actual1)
	}

	actualByte2, _ := (ioutil.ReadAll(result.Body))
	actual2 := string(actualByte2)
	expected2 := "hello"
	if expected2 != actual2 {
		t.Error("expected body:", expected2, "but get:", actual2)
	}
}

func TestMiddlewareChainReject(t *testing.T) {
	mw := MakeMiddleware(nil, doNotFound, doWriteHello)
	mockWriter := httptest.NewRecorder()
	mw.ServeHTTP(mockWriter, httptest.NewRequest("GET", "/", nil))

	result := mockWriter.Result()
	actual1 := result.StatusCode
	expected1 := http.StatusNotFound
	if expected1 != actual1 {
		t.Error("expected status:", expected1, "but get:", actual1)
	}

	actualByte2, _ := (ioutil.ReadAll(result.Body))
	actual2 := string(actualByte2)
	expected2 := ""
	if expected2 != actual2 {
		t.Error("expected body:", expected2, "but get:", actual2)
	}
}

func TestEmptyMiddleware(t *testing.T) {
	mw := MakeMiddleware(nil, []Doable{}...)
	if mw.This != nil {
		t.Error("expected nil interface")
	}
	if mw.Next != nil {
		t.Error("expected nil next middleware")
	}
	if mw.ValueMap != nil {
		t.Error("expected nil ValueMap")
	}
}
