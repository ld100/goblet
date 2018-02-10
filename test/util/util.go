package util

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
)

var VERSION = "0.0.1"

func FruitsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func TestFruits(t *testing.T) {
	// create http.Handler
	handler := chi.NewRouter()
	handler.Get("/fruits", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	// run server using httptest
	server := httptest.NewServer(handler)
	defer server.Close()

	// create httpexpect instance
	//e := httpexpect.New(t, server.URL)
	//
	//// is it working?
	//e.GET("/fruits").
	//	Expect().
	//	Status(http.StatusOK).JSON().Array().Empty()
}