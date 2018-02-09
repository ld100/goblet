package util

import (
	"fmt"
	"net/http"

	"github.com/mozillazg/request"
)

func Get(url string) {
	c := new(http.Client)
	req := request.NewRequest(c)
	resp, err := req.Get("http://httpbin.org/get")
	if err != nil {

	}
	j, err := resp.Json()
	if err != nil {

	}
	defer resp.Body.Close() // Don't forget close the response body
	fmt.Println(j)
}

func Post(url string, body map[string]string) {
	c := new(http.Client)
	req := request.NewRequest(c)
	req.Data = map[string]string{
		"key": "value",
		"a":   "123",
	}
	resp, err := req.Post("http://httpbin.org/post")
	if err != nil {

	}
	fmt.Println(resp)
}
