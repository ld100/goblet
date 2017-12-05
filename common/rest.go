package common

import (
	"net/http"
)

func RootController(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("root."))
}

func PingController(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func PanicController(w http.ResponseWriter, r *http.Request) {
	panic("test")
}
