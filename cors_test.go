package cors

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreation(t *testing.T) {
	emptyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	wrapped := New().Wrap(emptyHandler)
	recorder := httptest.NewRecorder()
	wrapped.ServeHTTP(recorder, nil)
	validateHeaders("", "", "", "", recorder, t)
}

func TestOrigin(t *testing.T) {
	emptyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	wrapped := New(WithOrigins("Foo")).Wrap(emptyHandler)
	recorder := httptest.NewRecorder()
	wrapped.ServeHTTP(recorder, nil)
	validateHeaders("Foo", "", "", "", recorder, t)
}

func TestOrigins(t *testing.T) {
	emptyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	wrapped := New(WithOrigins("Foo", "Bar")).Wrap(emptyHandler)
	recorder := httptest.NewRecorder()
	wrapped.ServeHTTP(recorder, nil)
	validateHeaders("Foo, Bar", "", "", "", recorder, t)
}

func TestMethod(t *testing.T) {
	emptyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	corsMw := New(WithMethods(http.MethodPut))
	wrapped := corsMw.Wrap(emptyHandler)
	recorder := httptest.NewRecorder()
	wrapped.ServeHTTP(recorder, nil)
	validateHeaders("", http.MethodPut, "", "", recorder, t)
}

func TestMethods(t *testing.T) {
	emptyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	corsMw := New(WithMethods(http.MethodDelete, http.MethodPost))
	wrapped := corsMw.Wrap(emptyHandler)
	recorder := httptest.NewRecorder()
	wrapped.ServeHTTP(recorder, nil)
	validateHeaders("", fmt.Sprintf("%s, %s", http.MethodDelete, http.MethodPost), "", "", recorder, t)
}

func TestHeader(t *testing.T) {
	emptyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	corsMw := New(WithHeaders("X-Foo"))
	wrapped := corsMw.Wrap(emptyHandler)
	recorder := httptest.NewRecorder()
	wrapped.ServeHTTP(recorder, nil)
	validateHeaders("", "", "X-Foo", "", recorder, t)
}

func TestHeaders(t *testing.T) {
	emptyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	corsMw := New(WithHeaders("X-Foo", "X-Bar"))
	wrapped := corsMw.Wrap(emptyHandler)
	recorder := httptest.NewRecorder()
	wrapped.ServeHTTP(recorder, nil)
	validateHeaders("", "", "X-Foo, X-Bar", "", recorder, t)
}

func TestAge(t *testing.T) {
	emptyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	corsMw := New(WithMaxAge(time.Hour))
	wrapped := corsMw.Wrap(emptyHandler)
	recorder := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodOptions, "", nil)
	if err != nil {
		t.Fatal(err)
	}
	wrapped.ServeHTTP(recorder, req)
	validateHeaders("", "", "", fmt.Sprint(time.Hour.Seconds()), recorder, t)
}

func validateHeaders(originVal, methodsVal, headersVal, ageVal string, recorder *httptest.ResponseRecorder, t *testing.T) {
	t.Helper()

	originHeader := recorder.Header().Get("Access-Control-Allow-Origin")
	methodsHeader := recorder.Header().Get("Access-Control-Allow-Methods")
	headersHeader := recorder.Header().Get("Access-Control-Allow-Headers")
	ageHeader := recorder.Header().Get("Access-Control-Max-Age")

	if originHeader != originVal {
		t.Fatal("unexpected header for \"Access-Control-Allow-Origin\":", originHeader)
	}

	if methodsHeader != methodsVal {
		t.Fatal("unexpected header for \"Access-Control-Allow-Methods\":", methodsHeader)
	}

	if headersHeader != headersVal {
		t.Fatal("unexpected header for \"Access-Control-Allow-Headers\":", headersHeader)
	}

	if ageHeader != ageVal {
		t.Fatal("unexpected header for \"Access-Control-Max-Age\":", ageHeader)
	}
}
