package main

import (
    "testing"
    "github.com/gorilla/mux"
    "net/http"
    "net/http/httptest"
    "github.com/stretchr/testify/assert"
    "fmt"
)


/*
func TestGetCharacters(t *testing.T) {

    request, err := http.NewRequest("GET", "/characters", nil)
    
    if err != nil {
		t.Fatal(err)
	}

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(GetCharacters)
    handler.ServeHTTP(rr, request)

    if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}
*/


func Router() *mux.Router {

    router := mux.NewRouter()
    router.HandleFunc("/", SampleResponse).Methods("GET")
    router.HandleFunc("/characters", GetCharacters).Methods("GET")
    router.HandleFunc("/characters/{id}", GetCharacter).Methods("GET")
    return router

}


func TestGetCharacters(t *testing.T) {

    request, _ := http.NewRequest("GET", "/characters", nil)
    response := httptest.NewRecorder()
    fmt.Println("asdfasfasdf")
    Router().ServeHTTP(response, request)

    assert.Equal(t, 200, response.Code, "OK response is expected")

}


func TestGetCharacters(t *testing.T) {

    request, _ := http.NewRequest("GET", "/characters/", nil)
    response := httptest.NewRecorder()
    fmt.Println("asdfasfasdf")

    Router().ServeHTTP(response, request)

    assert.Equal(t, 200, response.Code, "OK response is expected")



}


